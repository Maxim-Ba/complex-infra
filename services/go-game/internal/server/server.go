package server

import (
	"fmt"
	"log/slog"
	"net"
)

type Instance struct {
	conn *net.UDPConn
}

var Inst *Instance

func Start(a string) error {
	addr, err := net.ResolveUDPAddr("udp", a)
	if err != nil {
		return fmt.Errorf("server ResolveUDPAddr: %w", err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("server ListenUDP: " + err.Error())
	}
	Inst = &Instance{
		conn: conn,
	}
	buffer := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			slog.Error("ReadFromUDP:" + err.Error())
			continue
		}

		response := []byte("test: " + string(buffer[:n]))
		_, err = conn.WriteToUDP(response, clientAddr)
		if err != nil {
			slog.Error("WriteToUDP:" + err.Error())
		}
	}
}

func (s *Instance) Close() error {
	return s.conn.Close()
}
