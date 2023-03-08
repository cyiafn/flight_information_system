package server

import (
	"context"
	"time"

	"github.com/cyiafn/flight_information_system/server/custom_errors"
	"github.com/cyiafn/flight_information_system/server/dto"
	"github.com/cyiafn/flight_information_system/server/dto/status_code"
	"github.com/cyiafn/flight_information_system/server/duplicate_request"
	"github.com/cyiafn/flight_information_system/server/logs"
	"github.com/cyiafn/flight_information_system/server/net"
	"github.com/cyiafn/flight_information_system/server/utils"
	"github.com/cyiafn/flight_information_system/server/utils/bytes"
	"github.com/cyiafn/flight_information_system/server/utils/rpc"
)

type serverMode int

const (
	atMostOnceServerMode serverMode = iota + 1
	atLeastOnceServerMode
)

const (
	defaultUDPPort = 8080
	udpPortKey     = "UDP_LISTENER_PORT"

	requestTypeBytesLength   = 1
	shortIDBytesLength       = 9
	currentPacketBytesLength = 8
	totalPacketsByteLength   = 8
)

var instance *server

type server struct {
	UDPListener            net.Listener
	Routes                 map[dto.RequestType]func(ctx context.Context, request any) (any, error)
	Mode                   serverMode
	DuplicateRequestFilter *duplicate_request.Filter
	RequestBuffer          *requestBuffer
}

func Boot(routes map[dto.RequestType]func(ctx context.Context, request any) (any, error), atMostOnceEnabled bool) {
	instance = &server{
		Routes: routes,
		Mode:   utils.TernaryOperator(atMostOnceEnabled, atMostOnceServerMode, atLeastOnceServerMode),
	}

	if atMostOnceEnabled {
		instance.Mode = atMostOnceServerMode
		instance.DuplicateRequestFilter = duplicate_request.NewFilter()
	} else {
		instance.Mode = atLeastOnceServerMode
	}

	instance.RequestBuffer = newRequestBuffer()

	instance.UDPListener = net.NewUDPListener(getUDPPort(), instance.RouteRequest)
	instance.UDPListener.StartListening()

	time.Sleep(1 * time.Second) // grace time period so that closing listeners complete
}

func SpinDown() {
	logs.Info("disabling duplicate request filter.")
	instance.DuplicateRequestFilter.Close()
	logs.Info("Disabling UDP listener.")
	instance.UDPListener.StopListening()
	logs.Info("Goodbye!")
}

func (s *server) RouteRequest(ctx context.Context, request []byte) ([][]byte, bool) {
	req, complete := s.RequestBuffer.ProcessRequest(ctx, request)
	if !complete {
		logs.Info("Request is not complete, waiting for all packets: %s", utils.DumpJSON(req))
	}

	if s.Mode == atMostOnceServerMode && !s.DuplicateRequestFilter.IsAllowed(req.RequestID) {
		logs.Warn("RequestID: %s was repeated, aborting request", req.RequestID)
		return nil, false
	}

	requestType, requestBody := req.CompileRequest()

	logs.Info("Received Request Type: %v", requestType)
	requestDTO := dto.NewRequestDTO(requestType)
	if requestDTO != nil {
		err := rpc.Unmarshal(requestBody, &requestDTO)
		if err != nil {
			logs.Error("Unable to marshal request, err: %v", err)
			return nil, false
		}
	}

	handler, ok := s.Routes[requestType]
	if !ok {
		logs.Error("no route for request type: %v, ignoring request", requestType)
		return nil, false
	}

	response, err := handler(ctx, requestDTO)

	wrappedResp := &dto.Response{
		StatusCode: status_code.GetStatusCode(err),
		Data:       response,
	}

	resp, err := rpc.Marshal(wrappedResp)
	if err != nil {
		logs.Warn("error when marshalling, err: %v", err)
		resp, _ = rpc.Marshal(&dto.Response{
			StatusCode: status_code.GetStatusCode(custom_errors.NewMarshallerError(err)),
			Data:       nil,
		})
	}

	res := s.splitPayloadForSending(requestType, []byte(req.RequestID), resp)

	logs.Info("response: %v", resp)
	return res, true
}

func (s *server) splitPayloadForSending(requestType dto.RequestType, requestID []byte, payload []byte) [][]byte {
	output := make([][]byte, 0)
	it := 0
	for i := 0; i < len(payload)+net.DefaultByteBufferSize-s.getTotalBytesInHeader(); i += net.DefaultByteBufferSize - s.getTotalBytesInHeader() {
		mxSize := utils.TernaryOperator(len(payload) < i+net.DefaultByteBufferSize-s.getTotalBytesInHeader(), len(payload), i+net.DefaultByteBufferSize-s.getTotalBytesInHeader())
		output[it] = payload[i:mxSize]
		it += 1
	}

	for i := range output {
		output[i] = s.addHeaders(requestType, requestID, int64(i), int64(len(output)), output[i])
	}

	return output
}

func (s *server) getTotalBytesInHeader() int {
	return requestTypeBytesLength + shortIDBytesLength + currentPacketBytesLength + totalPacketsByteLength
}

func (s *server) addHeaders(requestType dto.RequestType, requestID []byte, packetNo int64, totalPacket int64, response []byte) []byte {
	response = s.addTotalPacketToHeader(response, totalPacket)
	response = s.addPacketNoToHeader(response, packetNo)
	response = s.addRequestIDToHeader(response, requestID)
	response = s.addResponseTypeHeader(response, dto.GetResponseType(requestType))
	return response
}

func (s *server) addPacketNoToHeader(response []byte, packetNo int64) []byte {
	return append(bytes.Int64ToBytes(packetNo), response...)
}

func (s *server) addTotalPacketToHeader(response []byte, totalPacket int64) []byte {
	return append(bytes.Int64ToBytes(totalPacket), response...)
}

func (s *server) addRequestIDToHeader(response []byte, requestID []byte) []byte {
	return append(requestID, response...)
}

func (s *server) addResponseTypeHeader(response []byte, responseType dto.ResponseType) []byte {
	return append([]byte{uint8(responseType)}, response...)
}

func getUDPPort() int {
	port, ok := utils.GetEnvInt(udpPortKey)
	if !ok {
		return defaultUDPPort
	}
	return port
}

func getRequestType(request []byte) dto.RequestType {
	return dto.RequestType(request[0])
}

func getRequestID(request []byte) []byte {
	dest := make([]byte, shortIDBytesLength)
	copy(dest, request[requestTypeBytesLength:requestTypeBytesLength+shortIDBytesLength])
	return dest
}

func getCurrentPacketNumber(request []byte) []byte {
	return request[requestTypeBytesLength+shortIDBytesLength : requestTypeBytesLength+shortIDBytesLength+currentPacketBytesLength]
}

func getTotalPacketNumber(request []byte) []byte {
	return request[requestTypeBytesLength+shortIDBytesLength+currentPacketBytesLength : requestTypeBytesLength+shortIDBytesLength+currentPacketBytesLength+totalPacketsByteLength]
}

func getRequestBody(request []byte) []byte {
	return request[requestTypeBytesLength+shortIDBytesLength+totalPacketsByteLength+currentPacketBytesLength:]
}

func GetIPAddr(ctx context.Context) string {
	return ctx.Value("addr").(string)
}
