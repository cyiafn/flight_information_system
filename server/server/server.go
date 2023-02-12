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
		rpc.Unmarshal(request, &requestDTO)
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

	resp, err := rpc.Marshal(wrappedResp, dto.GetResponseType(requestType))
	if err != nil {
		logs.Warn("error when marshalling, skipping response, err: %v", err)
		return nil
	}
	logs.Info("response: %v", resp)
	return resp
}

func getUDPPort() int {
	port, ok := utils.GetEnvInt(udpPortKey)
	if !ok {
		return defaultUDPPort
	}
	return port
}
