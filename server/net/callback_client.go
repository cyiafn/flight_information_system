package net

import (
	"net"

	"github.com/cyiafn/flight_information_system/server/logs"
)

func SendData(data []byte, addr string) error {
	if len(data) == 0 {
		logs.Warn("no data sent, payload had nothing")
		return nil
	}

	callbackConn, err := makeConn(addr)
	if err != nil {
		return err
	}

	defer callbackConn.Close()

	err = sendPayload(callbackConn, data)
	return err
}

func makeConn(addr string) (*net.UDPConn, error) {
	callbackClientAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		logs.Error("unable to resolve UDP address, err: %v", err)
		return nil, err
	}

	callbackConn, err := net.DialUDP("UDP", nil, callbackClientAddr)
	if err != nil {
		logs.Error("unable to dial UDP, err: %v", err)
		return nil, err
	}
	return callbackConn, nil
}

func sendPayload(conn *net.UDPConn, data []byte) error {
	for i := 0; i < len(data); i += defaultByteBufferSize {
		endIndex := i + defaultByteBufferSize
		if endIndex > len(data) {
			endIndex = len(data)
		}

		_, err := conn.Write(data[i:endIndex])
		if err != nil {
			logs.Error("error in sending data, aborting, err: %v", err)
			return err
		}
	}

	_, err := conn.Write(make([]byte, defaultByteBufferSize))
	if err != nil {
		logs.Error("error in sending data, aborting, err: %v", err)
		return err
	}
	return nil
}
