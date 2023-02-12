package net

import (
	"fmt"
	"net"

	"github.com/cyiafn/flight_information_system/server/logs"
	"github.com/cyiafn/flight_information_system/server/utils/bytes"
)

var _ Listener = (*UDPListener)(nil)

const (
	udpAddress = "localhost"

	defaultByteBufferSize = 512
)

type Listener interface {
	StartListening()
	StopListening()
}

func NewUDPListener(port int, requestHandler func(request []byte) []byte) Listener {
	return &UDPListener{
		listener:       nil,
		Port:           port,
		RequestHandler: requestHandler,
	}
}

type UDPListener struct {
	listener      net.PacketConn
	terminateChan chan struct{}

	Port           int
	RequestHandler func(request []byte) []byte
}

func (u *UDPListener) StartListening() {
	logs.Info("Booting up listener...")
	udpServer, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.ParseIP(udpAddress),
		Port: u.Port,
		Zone: "",
	})
	if err != nil {
		logs.Fatal("unable to start udp listener, err: %v", err)
	}
	u.listener = udpServer

	u.terminateChan = make(chan struct{})
	logs.Info("Good day, listener booted up.")

	var prevAddress net.Addr
	var currentRequest []byte

	for {
		buf := make([]byte, defaultByteBufferSize)
		n, addr, err := u.listener.ReadFrom(buf)
		if err != nil {
			if err.Error() == fmt.Sprintf("read udp [::]:%v: use of closed network connection", u.Port) {
				return
			}
			logs.Warn("unable to read from buffer, err: %v", err)
			continue
		}

		if prevAddress == nil {
			prevAddress = addr
		}

		logs.Info("Received request of len %v from addr %s, data: %v", n, addr.String(), buf)

		if bytes.IsEmpty(buf) {
			currentRequest = append(currentRequest, buf...)
			go u.handleIncomingData(currentRequest, addr)
			currentRequest = nil
			prevAddress = nil
		} else {
			if prevAddress.String() != addr.String() {
				logs.Error("previous IP address did not match current IP address, discarding data due to potential corruption")
				prevAddress = nil
				currentRequest = nil
			} else {
				currentRequest = append(currentRequest, buf...)
			}
		}
	}
}

func (u *UDPListener) handleIncomingData(buf []byte, addr net.Addr) {
	resp := u.RequestHandler(buf)

	_, err := u.listener.WriteTo(resp, addr)
	if err != nil {
		logs.Error("unable to reply ")
	}
}

func (u *UDPListener) StopListening() {
	err := u.listener.Close()
	if err != nil {
		logs.Warn("unable to close listener, err: %v, you might need to restart your computer")
	}
	logs.Info("Listener stopped")
}
