package net

import (
	"net"

	"github.com/cyiafn/flight_information_system/server/logs"
)

/**
This callback_client is used to only send data back to the subscriber.
Here, we open and close connections as needed only and do not define a singular port so that we can
concurrently broadcast to all subscribers.
*/

// SendData simply sends a payload to an IP:port
func SendData(data []byte, addr string) error {
	if len(data) == 0 {
		logs.Warn("no data sent, payload had nothing")
		return nil
	}

	// creates the "connection"
	callbackConn, err := makeConn(addr)
	if err != nil {
		return err
	}

	// closes the connection once done
	defer callbackConn.Close()

	// sends payload
	err = sendPayload(callbackConn, data)
	return err
}

// creates the necessary object for sending data over UDP
func makeConn(addr string) (net.Conn, error) {
	callbackConn, err := net.Dial("udp", addr)
	if err != nil {
		logs.Error("unable to dial UDP, err: %v", err)
		return nil, err
	}
	return callbackConn, nil
}

// Simply writes the data into the "connection" object. Sends data to client
func sendPayload(conn net.Conn, data []byte) error {
	_, err := conn.Write(data)
	if err != nil {
		logs.Error("error sending payload: %v", err)
	}
	return err

}
