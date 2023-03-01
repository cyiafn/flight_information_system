package net

import (
	"context"
	"fmt"
	"net"

	"github.com/cyiafn/flight_information_system/server/logs"
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

func NewUDPListener(port int, requestHandler func(ctx context.Context, request []byte) ([]byte, bool)) Listener {
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
	RequestHandler func(ctx context.Context, request []byte) ([]byte, bool)
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

	u.listen()
}

func (u *UDPListener) listen() {
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

		logs.Info("Received request of len %v from addr %s, data: %v", n, addr.String(), buf)

		ctx := context.WithValue(context.Background(), "addr", addr.String())
		go u.handleIncomingData(ctx, buf, addr)
	}
}

func (u *UDPListener) handleIncomingData(ctx context.Context, buf []byte, addr net.Addr) {
	resp, processed := u.RequestHandler(ctx, buf)
	if !processed {
		return
	}
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
