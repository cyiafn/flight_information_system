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
	body = addTotalPacketToHeader(body, 1)
	body = addPacketNoToHeader(body, 1)
	body = addRequestID(body)
	body = addResponseTypeHeader(respType, body)
	return body
}

// addPacketNoToHeader appends the packet number to the header
func addPacketNoToHeader(response []byte, packetNo int64) []byte {
	return append(bytes.Int64ToBytes(packetNo), response...)
}

// addTotalPacketToHeader appends the total packet number to the header
func addTotalPacketToHeader(response []byte, totalPacket int64) []byte {
	return append(bytes.Int64ToBytes(totalPacket), response...)
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
