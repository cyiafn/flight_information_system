package server

import (
	"time"

	"github.com/cyiafn/flight_information_system/server/dto"
	"github.com/cyiafn/flight_information_system/server/dto/status_code"
	"github.com/cyiafn/flight_information_system/server/logs"
	"github.com/cyiafn/flight_information_system/server/net"
	"github.com/cyiafn/flight_information_system/server/utils"
	"github.com/cyiafn/flight_information_system/server/utils/rpc"
)

const (
	defaultUDPPort = 8080
	udpPortKey     = "UDP_LISTENER_PORT"
)

var instance *server

func Boot(routes map[dto.RequestType]func(request any) (any, error)) {
	instance = &server{
		Routes: routes,
	}
	instance.UDPListener = net.NewUDPListener(getUDPPort(), instance.RouteRequest)
	instance.UDPListener.StartListening()
	time.Sleep(1 * time.Second) // grace time period so that closing listeners complete
}

func SpinDown() {
	logs.Info("Disabling UDP listener.")
	instance.UDPListener.StopListening()
	logs.Info("Goodbye!")
}

type server struct {
	UDPListener net.Listener
	Routes      map[dto.RequestType]func(request any) (any, error)
}

func (s *server) RouteRequest(request []byte) []byte {
	requestType := rpc.GetRequestType(request)

	logs.Info("Received Request Type: %v", requestType)
	requestDTO := dto.NewRequestDTO(requestType)
	if requestDTO != nil {
		err := rpc.Unmarshal(request[1:], &requestDTO)
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

	response, err := handler(requestDTO)

	wrappedResp := &dto.Response{
		StatusCode: status_code.GetStatusCode(err),
		Data:       response,
	}

	resp, err := rpc.Marshal(wrappedResp)
	if err != nil {
		logs.Warn("error when marshalling, skipping response, err: %v", err)
		return nil
	}

	resp = s.addHeaders(requestType, resp)

	logs.Info("response: %v", resp)
	return resp
}

func (s *server) addHeaders(requestType dto.RequestType, response []byte) []byte {
	// add requestID
	response = s.addResponseTypeHeader(response, dto.GetResponseType(requestType))
	return response
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
