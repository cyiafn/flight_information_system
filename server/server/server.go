package server

import (
	"context"
	"time"

	"github.com/cyiafn/flight_information_system/server/dto"
	"github.com/cyiafn/flight_information_system/server/dto/status_code"
	"github.com/cyiafn/flight_information_system/server/duplicate_request"
	"github.com/cyiafn/flight_information_system/server/logs"
	"github.com/cyiafn/flight_information_system/server/net"
	"github.com/cyiafn/flight_information_system/server/utils"
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

	requestTypeBytesLength = 1
	shortIDBytesLength     = 9
)

var instance *server

type server struct {
	UDPListener            net.Listener
	Routes                 map[dto.RequestType]func(ctx context.Context, request any) (any, error)
	Mode                   serverMode
	DuplicateRequestFilter *duplicate_request.Filter
}

func Boot(routes map[dto.RequestType]func(ctx context.Context, request any) (any, error), atMostOnceEnabled bool) {
	instance = &server{
		Routes: routes,
		Mode:   utils.TernaryOperator(atMostOnceEnabled, atMostOnceServerMode, atLeastOnceServerMode),
	}

	if atMostOnceEnabled {
		instance.Mode = atMostOnceServerMode
		duplicate_request.NewFilter()
	} else {
		instance.Mode = atLeastOnceServerMode
	}

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

func (s *server) RouteRequest(ctx context.Context, request []byte) []byte {
	requestID := getRequestID(request)
	if s.Mode == atMostOnceServerMode && !s.DuplicateRequestFilter.IsAllowed(string(requestID)) {
		logs.Warn("RequestID: %s was repeated, aborting request", requestID)
		return nil
	}

	requestType := getRequestType(request)
	requestBody := getRequestBody(request)

	logs.Info("Received Request Type: %v", requestType)
	requestDTO := dto.NewRequestDTO(requestType)
	if requestDTO != nil {
		err := rpc.Unmarshal(requestBody, &requestDTO)
		if err != nil {
			logs.Error("Unable to marshal request, err: %v", err)
			return nil
		}
	}

	handler, ok := s.Routes[requestType]
	if !ok {
		logs.Error("no route for request type: %v, ignoring request", requestType)
		return nil
	}

	response, err := handler(ctx, requestDTO)

	wrappedResp := &dto.Response{
		StatusCode: status_code.GetStatusCode(err),
		Data:       response,
	}

	resp, err := rpc.Marshal(wrappedResp)
	if err != nil {
		logs.Warn("error when marshalling, skipping response, err: %v", err)
		return nil
	}

	resp = s.addHeaders(requestType, requestID, resp)

	logs.Info("response: %v", resp)
	return resp
}

func (s *server) addHeaders(requestType dto.RequestType, requestID []byte, response []byte) []byte {
	response = s.addRequestIDToHeader(requestID, response)
	response = s.addResponseTypeHeader(response, dto.GetResponseType(requestType))
	return response
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

func getRequestBody(request []byte) []byte {
	return request[requestTypeBytesLength+shortIDBytesLength:]
}

func GetIPAddr(ctx context.Context) string {
	return ctx.Value("addr").(string)
}
