package net

import (
	"github.com/cyiafn/flight_information_system/server/dto"
	"github.com/cyiafn/flight_information_system/server/utils/bytes"
	"github.com/cyiafn/flight_information_system/server/utils/rpc"
	"github.com/cyiafn/flight_information_system/server/utils/worker_pools"
	"github.com/teris-io/shortid"
	"testing"
)

func TestSendData(t *testing.T) {
	payload := &dto.Response{StatusCode: 1, Data: &dto.GetFlightInformationRequest{FlightIdentifier: 1}}
	payloadBytes, _ := rpc.Marshal(payload)
	resp := addHeaders(101, payloadBytes)

	users := make([]string, 20000)
	for i := range users {
		users[i] = "localhost:8080"
	}

	worker_pools.Load(func(req string) error {
		return SendData(resp, req)
	}, users, 10000)
}

func addHeaders(respType dto.ResponseType, body []byte) []byte {
	body = addTotalByteBuffersToHeader(body, 1)
	body = addByteBufferNoToHeader(body, 1)
	body = addRequestID(body)
	body = addResponseTypeHeader(respType, body)
	return body
}

// addByteBufferNoToHeader appends the byteBufferArray number to the header
func addByteBufferNoToHeader(response []byte, byteBufferArrayNo int64) []byte {
	return append(bytes.Int64ToBytes(byteBufferArrayNo), response...)
}

// addTotalByteBuffersToHeader appends the total byteBufferArray number to the header
func addTotalByteBuffersToHeader(response []byte, totalByteBufferArray int64) []byte {
	return append(bytes.Int64ToBytes(totalByteBufferArray), response...)
}

// addResponseTypeHeader appends the responseType to the header
func addResponseTypeHeader(respType dto.ResponseType, body []byte) []byte {
	return append([]byte{uint8(respType)}, body...)
}

// addRequestID appends the requestID to the header
func addRequestID(body []byte) []byte {
	// this generates a random 9 character string guaranteed for uniqueness until 2050.
	reqID := shortid.MustGenerate()
	return append([]byte(reqID), body...)
}
