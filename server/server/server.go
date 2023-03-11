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

// serverMode is at least once or at most once mode
type serverMode int

const (
	atMostOnceServerMode serverMode = iota + 1
	atLeastOnceServerMode
)

const (
	// defaultUDPPort if env var is not set
	defaultUDPPort = 8080
	// udpPortKey for env var
	udpPortKey = "UDP_LISTENER_PORT"

	// requestType length in bytes
	requestTypeBytesLength = 1
	// requestID length in bytes
	shortIDBytesLength = 9
	// currentByteBufferArrayBytesLength no in bytes
	currentByteBufferArrayBytesLength = 8
	// totalByteBufferArrayByteLength no in bytes
	totalByteBufferArrayByteLength = 8
)

// single instance of server. This will be "singleton"
var instance *server

// the actual server orchestrating everything
type server struct {
	// UDPListener is to listen to incoming data
	UDPListener net.Listener
	// Routes routes a request to a piece of business logic
	Routes map[dto.RequestType]func(ctx context.Context, request any) (any, error)
	// Mode is at least once or at least once
	Mode serverMode
	// DuplicateRequestFilter is the filter for duplicate requests
	DuplicateRequestFilter *duplicate_request.Filter
	// RequestBuffer is the request buffer for timing out requests, processing multiple byteArrayBuffers and allowing for concurrent server access
	RequestBuffer *requestBuffer
}

// Boot initialises the server instance boots up the server
func Boot(routes map[dto.RequestType]func(ctx context.Context, request any) (any, error), atMostOnceEnabled bool) {
	instance = &server{
		Routes: routes,
		Mode:   utils.TernaryOperator(atMostOnceEnabled, atMostOnceServerMode, atLeastOnceServerMode),
	}

	// If at most once is enabled, we need the duplicate request filter to prevent duplicate requests from running multiple times
	if atMostOnceEnabled {
		instance.Mode = atMostOnceServerMode
		instance.DuplicateRequestFilter = duplicate_request.NewFilter()
	} else {
		instance.Mode = atLeastOnceServerMode
	}

	// instantiating all dependencies
	instance.RequestBuffer = newRequestBuffer()
	// take note here, that the servers route request function is passed ito the UDPListener such that all byteArrayBuffers will be received by the server, processed, routed, executed,
	// before the data is passed back the UDPListener to send back
	instance.UDPListener = net.NewUDPListener(getUDPPort(), instance.RouteRequest)
	instance.UDPListener.StartListening()

	// we sleep for 1 second on termination here as StartListening blocks and is the main thread. Upon interception of SIGINT or SIGKILL we need some grace period for all dependencies to close gracefully
	time.Sleep(1 * time.Second) // grace time period so that closing listeners complete
}

// SpinDown kills the instance and closes their deps.
func SpinDown() {
	logs.Info("disabling duplicate request filter.")
	instance.DuplicateRequestFilter.Close()
	logs.Info("Disabling UDP listener.")
	instance.UDPListener.StopListening()
	logs.Info("Goodbye!")
}

// RouteRequest is the callback function passed into the UDPListener to intercept all received data and process it accordingly
func (s *server) RouteRequest(ctx context.Context, request []byte) ([][]byte, bool) {
	// Sends the request to the request buffer to check if all byteArrayBuffers have arrived or not and whether we should process this right now.
	req, complete := s.RequestBuffer.ProcessRequest(ctx, request)
	if !complete {
		// if not all byte arrays have arrived, we do not process it
		logs.Info("Request is not complete, waiting for all byte arrays: %s", utils.DumpJSON(req))
		return nil, false
	}

	// if we decide to process it and it is at most once server mode, we need to check if it is allowed (if it was a duplicate request)
	// if the filter does not allow us to process, we discard this request.
	if s.Mode == atMostOnceServerMode && !s.DuplicateRequestFilter.IsAllowed(req.RequestID) {
		logs.Warn("RequestID: %s was repeated, aborting request", req.RequestID)
		return nil, false
	}

	// We take the request object and compile it into the necessary information
	requestType, requestBody := req.CompileRequest()

	// we generate the requestDTO object based on the requestType
	requestDTO := dto.NewRequestDTO(requestType)
	if requestDTO != nil {
		// unmarshal the request body into the DTO
		err := rpc.Unmarshal(requestBody, requestDTO)
		if err != nil {
			logs.Error("Unable to marshal request, err: %v", err)
			return nil, false
		}
	}
	logs.Info("[%s] Received Request Type: %v, Request ID: %s, Request No: %v, Total Byte Array Buffers for Request %v, Marshalled Request: %s",
		GetIPAddr(ctx),
		requestType,
		string(getRequestID(request)),
		bytes.ToInt64(getCurrentByteBufferArrayNumber(request)),
		bytes.ToInt64(getTotalByteBufferArrayNumber(request)),
		utils.DumpJSON(requestDTO),
	)

	// we route it to the correct handler based on routes provided on server boot
	handler, ok := s.Routes[requestType]
	if !ok {
		logs.Error("no route for request type: %v, ignoring request", requestType)
		return nil, false
	}

	// we execute the RPC call with the proper handler/biz logic
	response, err := handler(ctx, requestDTO)

	// we wrap the response in the response DTO wrapper such that we can properly send proper error messages to the user
	wrappedResp := &dto.Response{
		StatusCode: status_code.GetStatusCode(err),
		Data:       response,
	}

	// we marshal the wrapped response
	resp, err := rpc.Marshal(wrappedResp)
	if err != nil {
		logs.Warn("error when marshalling, err: %v", err)
		// we throw a generic marshaller error if we can't marshal for some reason
		resp, _ = rpc.Marshal(&dto.Response{
			StatusCode: status_code.GetStatusCode(custom_errors.NewMarshallerError(err)),
			Data:       nil,
		})
	}

	// our payload might be more than 512 bytes, so we might need to split it into multiple byte arrays. This will not happen in this presentation but the functionality is there
	res := s.splitPayloadForSending(requestType, []byte(req.RequestID), resp)

	for i, payload := range res {
		logs.Info("[%s] Response Payload #%v out of %v: Request Type: %v, Request ID: %s, Request No: %v, Total Byte Arrays for Request %v, Marshalled Request: %s",
			GetIPAddr(ctx),
			i+1,
			len(res),
			getRequestType(payload),
			string(getRequestID(payload)),
			bytes.ToInt64(getCurrentByteBufferArrayNumber(payload)),
			bytes.ToInt64(getTotalByteBufferArrayNumber(payload)),
			utils.DumpJSON(wrappedResp),
		)
	}

	// returns the response data to the user to the UDPListener to send back
	return res, true
}

// splitPayloadForSending splits the payload into multiple byte array buffers to send
func (s *server) splitPayloadForSending(requestType dto.RequestType, requestID []byte, payload []byte) [][]byte {
	// if the payload length == 0 we can hardcode this
	if len(payload) == 0 {
		output := make([][]byte, 1)
		output[0] = make([]byte, 0)
		output[0] = s.addHeaders(requestType, requestID, 1, 1, output[0])
		return output
	}
	output := make([][]byte, 0)
	// we split it up into array of byte arrays
	for i := 0; i < len(payload); i += net.DefaultByteBufferSize - s.getTotalBytesInHeader() {
		mxSize := utils.TernaryOperator(len(payload) < i+net.DefaultByteBufferSize-s.getTotalBytesInHeader(), len(payload), i+net.DefaultByteBufferSize-s.getTotalBytesInHeader())
		output = append(output, payload[i:mxSize])
	}

	for i := range output {
		// we add headers for each byte array
		output[i] = s.addHeaders(requestType, requestID, int64(i), int64(len(output)), output[i])
	}

	return output
}

// getTotalBytesInHeader returns header size
func (s *server) getTotalBytesInHeader() int {
	return requestTypeBytesLength + shortIDBytesLength + currentByteBufferArrayBytesLength + totalByteBufferArrayByteLength
}

// addHeaders adds headers to a payload
func (s *server) addHeaders(requestType dto.RequestType, requestID []byte, byteArrayBufferNo int64, totalByteArrayBuffer int64, response []byte) []byte {
	response = s.addTotalByteArraysToHeader(response, totalByteArrayBuffer)
	response = s.addByteArrayBufferNoToHeader(response, byteArrayBufferNo)
	response = s.addRequestIDToHeader(response, requestID)
	response = s.addResponseTypeHeader(response, dto.GetResponseType(requestType))
	return response
}

func (s *server) addByteArrayBufferNoToHeader(response []byte, byteArrayBufferNo int64) []byte {
	return append(bytes.Int64ToBytes(byteArrayBufferNo+1), response...)
}

func (s *server) addTotalByteArraysToHeader(response []byte, totalByteArrayBuffers int64) []byte {
	return append(bytes.Int64ToBytes(totalByteArrayBuffers), response...)
}

func (s *server) addRequestIDToHeader(response []byte, requestID []byte) []byte {
	return append(requestID, response...)
}

func (s *server) addResponseTypeHeader(response []byte, responseType dto.ResponseType) []byte {
	return append([]byte{uint8(responseType)}, response...)
}

// getUDPPort based on env var. Defaults to defaultUDPPort if not configured
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

func getCurrentByteBufferArrayNumber(request []byte) []byte {
	return request[requestTypeBytesLength+shortIDBytesLength : requestTypeBytesLength+shortIDBytesLength+currentByteBufferArrayBytesLength]
}

func getTotalByteBufferArrayNumber(request []byte) []byte {
	return request[requestTypeBytesLength+shortIDBytesLength+currentByteBufferArrayBytesLength : requestTypeBytesLength+shortIDBytesLength+currentByteBufferArrayBytesLength+totalByteBufferArrayByteLength]
}

func getRequestBody(request []byte) []byte {
	return request[requestTypeBytesLength+shortIDBytesLength+totalByteBufferArrayByteLength+currentByteBufferArrayBytesLength:]
}

func GetIPAddr(ctx context.Context) string {
	return ctx.Value("addr").(string)
}
