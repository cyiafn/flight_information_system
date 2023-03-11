package net

import (
	"context"
	"fmt"
	"github.com/cyiafn/flight_information_system/server/utils"
	"net"

	"github.com/cyiafn/flight_information_system/server/logs"
)

// Validate interface compliance for listener at compile time.
// Useful to swap for TCP listener in the future
var _ Listener = (*UDPListener)(nil)

const (
	// address of listener of server
	udpAddress = "localhost"

	udpAddressEnvKey = "IP_ADDRESS"
	// DefaultByteBufferSize of each request
	DefaultByteBufferSize = 512
)

// Listener interface to listen to requests
type Listener interface {
	StartListening()
	StopListening()
}

// NewUDPListener instantiates a listener.
func NewUDPListener(port int, requestHandler func(ctx context.Context, request []byte) ([][]byte, bool)) Listener {
	return &UDPListener{
		listener:       nil,
		Port:           port,
		RequestHandler: requestHandler,
	}
}

// UDPListener is a UDP listener.
type UDPListener struct {
	// listener stores the actual listener object
	listener net.PacketConn
	// Port is the port of the listener
	Port int
	// RequestHandler is the callback handler for all incoming data to the listener. This will be provided by the server.
	RequestHandler func(ctx context.Context, request []byte) ([][]byte, bool)
}

// StartListening starts the listener
func (u *UDPListener) StartListening() {
	logs.Info("Booting up listener...")
	// starts listener
	udpServer, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.ParseIP(getIPAddress()),
		Port: u.Port,
		Zone: "",
	})
	if err != nil {
		logs.Fatal("unable to start udp listener, err: %v", err)
	}
	u.listener = udpServer

	logs.Info("Good day, listener booted up.")

	// event loop for processing requests
	u.listen()
}

func (u *UDPListener) listen() {
	for {
		buf := make([]byte, DefaultByteBufferSize)
		// blocks until there is data being read from buffer
		n, addr, err := u.listener.ReadFrom(buf)
		if err != nil {
			if err.Error() == fmt.Sprintf("read udp [::]:%v: use of closed network connection", u.Port) {
				return
			}
			logs.Warn("unable to read from buffer, err: %v", err)
			continue
		}

		logs.Info("Received request of len %v from addr %s, data: %v", n, addr.String(), buf)

		// we add the IP address:port of the request to the context object
		ctx := context.WithValue(context.Background(), "addr", addr.String())
		// spawn a go routine to process each incoming data
		go u.handleIncomingData(ctx, buf, addr)
	}
}

// handleIncomingData handles all incoming data and processes data to return
func (u *UDPListener) handleIncomingData(ctx context.Context, buf []byte, addr net.Addr) {
	// passes the request to the requestHandler (server callback function) outlined during instantiation of this object
	// this will return a response and whether the request was processed or not
	resp, processed := u.RequestHandler(ctx, buf)
	// request might be unprocessed
	if !processed {
		return
	}
	// Only replies if there is a response, but there should generally be one as server wraps any response payload in the generic response type.
	if resp == nil {
		logs.Warn("no reply to user as response is nil")
		return
	}

	// for each "packet" buffer, we send it back to the client
	for _, packet := range resp {
		packet := packet
		_, err := u.listener.WriteTo(packet, addr)
		if err != nil {
			logs.Error("unable to reply, err: %v", err)
			return
		}
	}

}

// StopListening gracefully closes the listener, freeing up the port and terminating the listeners services
func (u *UDPListener) StopListening() {
	err := u.listener.Close()
	if err != nil {
		logs.Warn("unable to close listener, err: %v, you might need to restart your computer")
	}
	logs.Info("Listener stopped")
}

func getIPAddress() string {
	port, ok := utils.GetEnvStr(udpAddressEnvKey)
	if !ok {
		return udpAddress
	}
	return port
}
