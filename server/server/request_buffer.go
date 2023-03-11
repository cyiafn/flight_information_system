package server

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cyiafn/flight_information_system/server/dto"
	"github.com/cyiafn/flight_information_system/server/logs"
	"github.com/cyiafn/flight_information_system/server/utils/bytes"
)

/*
A request may be split into multiple buffers... E.g. a 1kb request might be split into 2 * 512 byte buffers.
This request buffer allows us to keep track of all requests coming in, puts them together and as such allow us to
truly support concurrent connections of variable length (even though not required in report).

This fully supports concurrent requests from a single client.
*/

const (
	// cleanUpDuration is the timing out of a request
	cleanUpDuration = 5 * time.Second
)

// newRequestBuffer instantiates a new requestBuffer
func newRequestBuffer() *requestBuffer {
	reqBuf := &requestBuffer{
		Buffer: make(map[string]*request),
	}
	reqBuf.StartCleanUp()
	return reqBuf
}

// requestBuffer simply stores the map of IP addresses + requestID to a request
type requestBuffer struct {
	sync.RWMutex
	Buffer map[string]*request
}

// ProcessRequest checks if all the byte arrays for a request have arrived, if not, it will not release the request for processing
func (r *requestBuffer) ProcessRequest(ctx context.Context, payload []byte) (*request, bool) {
	// generates the key for the buffer key
	key := makeBufferKey(GetIPAddr(ctx), string(getRequestID(payload)))
	r.RLock()
	request, ok := r.Buffer[key]
	if ok {
		// if there is such a request, we can simply use the byte array number in the request
		request.Body[bytes.ToInt64(getCurrentByteBufferArrayNumber(payload))] = getRequestBody(payload)
		r.RUnlock()
	} else {
		r.RUnlock()
		r.Lock()
		// creates a new request with that payload
		r.Buffer[key] = newRequest(ctx, payload)
		r.Unlock()
	}

	if r.Buffer[key].IsComplete() {
		r.Lock()
		defer func() {
			delete(r.Buffer, key)
			r.Unlock()
		}()
		// if all the byteBufferArrays are here, return the request for processing
		return r.Buffer[key], true
	}
	return nil, false
}

// StartCleanUp ticks every 2 seconds to clean up timed out requests
func (r *requestBuffer) StartCleanUp() {
	ticker := time.NewTicker(2 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				r.RLock()
				for key, value := range r.Buffer {
					if !value.TimedOut() {
						continue
					}
					if _, ok := r.Buffer[key]; ok {
						logs.Info("Timing out requestID: %s as it has exceeded 5 seconds to deliver all byte arrays.", value.RequestID)
						r.RUnlock()
						r.Lock()
						delete(r.Buffer, key)
						r.Unlock()
						r.RLock()
					}
				}
				r.RUnlock()
			}
		}
	}()
}

// newRequest creates a new request
func newRequest(ctx context.Context, payload []byte) *request {
	currentByteArrayBuffer := bytes.ToInt64(getCurrentByteBufferArrayNumber(payload))
	totalByteArrayBuffer := bytes.ToInt64(getTotalByteBufferArrayNumber(payload))
	req := &request{
		IPAddr:               GetIPAddr(ctx),
		RequestID:            string(getRequestID(payload)),
		Type:                 getRequestType(payload),
		TimeCreated:          time.Now(),
		TotalByteArrayBuffer: totalByteArrayBuffer,
		Body:                 make([][]byte, totalByteArrayBuffer),
	}
	req.Body[currentByteArrayBuffer-1] = getRequestBody(payload)
	return req
}

// request represents a request
type request struct {
	IPAddr    string
	RequestID string
	Type      dto.RequestType

	TimeCreated          time.Time
	TotalByteArrayBuffer int64
	Body                 [][]byte
}

// TimedOut checks if a request is timed out or not
func (r *request) TimedOut() bool {
	return r.TimeCreated.Add(cleanUpDuration).After(time.Now())
}

// IsComplete checks if all the byte arrays are here
func (r *request) IsComplete() bool {
	count := int64(0)
	for _, bodyPayload := range r.Body {
		if len(bodyPayload) != 0 {
			count += 1
		}
	}
	return count == r.TotalByteArrayBuffer
}

// CompileRequest gets all the compiled bodies of the different byte arrays.
func (r *request) CompileRequest() (dto.RequestType, []byte) {
	var response []byte

	for _, part := range r.Body {
		part := part
		response = append(response, part...)
	}

	return r.Type, response
}

// makeBufferKey makes the key for the buffer
func makeBufferKey(ipAddr, requestID string) string {
	return fmt.Sprintf("%s_%s", ipAddr, requestID)
}
