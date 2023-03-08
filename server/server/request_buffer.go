package server

import (
	"context"
	"fmt"
	"time"

	"github.com/cyiafn/flight_information_system/server/dto"
	"github.com/cyiafn/flight_information_system/server/logs"
	"github.com/cyiafn/flight_information_system/server/utils/bytes"
)

const (
	cleanUpDuration = 5 * time.Second
)

func newRequestBuffer() *requestBuffer {
	reqBuf := &requestBuffer{
		Buffer: make(map[string]*request),
	}
	reqBuf.StartCleanUp()
	return reqBuf
}

type requestBuffer struct {
	Buffer map[string]*request
}

func (r *requestBuffer) ProcessRequest(ctx context.Context, payload []byte) (*request, bool) {
	key := makeBufferKey(GetIPAddr(ctx), string(getRequestID(payload)))
	request, ok := r.Buffer[key]
	if ok {
		request.Body[bytes.ToInt64(getCurrentPacketNumber(payload))] = getRequestBody(payload)
	} else {
		r.Buffer[key] = newRequest(ctx, payload)
	}

	if r.Buffer[key].IsComplete() {
		delete(r.Buffer, key)
		return r.Buffer[key], true
	}
	return nil, false
}

func (r *requestBuffer) StartCleanUp() {
	ticker := time.NewTicker(2 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				for key, value := range r.Buffer {
					if !value.TimedOut() {
						continue
					}
					if _, ok := r.Buffer[key]; ok {
						logs.Info("Timing out requestID: %s as it has exceeded 5 seconds to deliver all packets.", value.RequestID)
						delete(r.Buffer, key)
					}
				}
			}
		}
	}()
}

func newRequest(ctx context.Context, payload []byte) *request {
	currentPacket := bytes.ToInt64(getCurrentPacketNumber(payload))
	totalPackets := bytes.ToInt64(getTotalPacketNumber(payload))
	req := &request{
		IPAddr:        GetIPAddr(ctx),
		RequestID:     string(getRequestID(payload)),
		Type:          getRequestType(payload),
		TimeCreated:   time.Now(),
		CurrentPacket: currentPacket,
		TotalPackets:  totalPackets,
		Body:          make([][]byte, totalPackets),
	}
	req.Body[currentPacket-1] = getRequestBody(payload)
	return req
}

type request struct {
	IPAddr    string
	RequestID string
	Type      dto.RequestType

	TimeCreated   time.Time
	CurrentPacket int64
	TotalPackets  int64
	Body          [][]byte
}

func (r *request) TimedOut() bool {
	return r.TimeCreated.Add(cleanUpDuration).After(time.Now())
}

func (r *request) IsComplete() bool {
	count := int64(0)
	for _, bodyPayload := range r.Body {
		if len(bodyPayload) != 0 {
			count += 1
		}
	}
	return count == r.TotalPackets
}

func (r *request) CompileRequest() (dto.RequestType, []byte) {
	var response []byte

	for _, part := range r.Body {
		part := part
		response = append(response, part...)
	}

	return r.Type, response
}

func makeBufferKey(ipAddr, requestID string) string {
	return fmt.Sprintf("%s_%s", ipAddr, requestID)
}
