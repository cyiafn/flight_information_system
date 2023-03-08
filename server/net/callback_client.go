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

func makeConn(addr string) (net.Conn, error) {
	callbackConn, err := net.Dial("UDP", addr)
	if err != nil {
		logs.Error("unable to dial UDP, err: %v", err)
		return nil, err
	}
	return callbackConn, nil
}

func sendPayload(conn net.Conn, data []byte) error {
	_, err := conn.Write(data)
	if err != nil {
		logs.Error("error sending payload: %v", err)
	}
	return err

}
