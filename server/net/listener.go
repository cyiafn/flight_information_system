package net

import (
	"fmt"
	"net"
	"time"

	"github.com/cyiafn/flight_information_system/server/logs"
	"github.com/cyiafn/flight_information_system/server/utils/bytes"
)

var _ Listener = (*UDPListener)(nil)

const (
	udpAddress = "localhost"

	defaultByteBufferSize = 512

	defaultServerTimeout = 5 * time.Second
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
	firstRequestStart := time.Now()

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
			firstRequestStart = time.Now()
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
			} else if firstRequestStart.Add(defaultServerTimeout).Before(time.Now()) {
				logs.Error("Received another packet after default server timeout, discarding data due to potential corruption")
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
	if resp == nil {
		logs.Warn("no reply to user as response is nil")
		return
	}

	_, err := u.listener.WriteTo(resp, addr)
	if err != nil {
		logs.Error("unable to reply, err: %v", err)
		return
	}
	_, err = u.listener.WriteTo(make([]byte, defaultByteBufferSize), addr)
	if err != nil {
		logs.Error("unable to end reply, err: %v", err)
	}

}

func (u *UDPListener) StopListening() {
	err := u.listener.Close()
	if err != nil {
		logs.Warn("unable to close listener, err: %v, you might need to restart your computer")
	}
	logs.Info("Listener stopped")
}
