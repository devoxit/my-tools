package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

type ConnStats struct {
	auth   bool
	id     string
	name   string
	reader *bufio.Reader
	writer *bufio.Writer
}

type Server struct {
	conns  map[*net.Conn]ConnStats
	secret string
}

func NewServer() *Server {
	return &Server{
		conns:  make(map[*net.Conn]ConnStats),
		secret: "EnCryp!e0?",
	}
}

func (s *Server) handleConnection(conn *net.Conn) {
	if conn == nil {
		return
	}
	if len(s.conns) > 0 {
		return
	}

	fmt.Println("New connection ...", conn)
	stats := ConnStats{}
	stats.auth = false
	stats.id = "executor"
	stats.name = ""
	stats.reader = bufio.NewReader(*conn)
	stats.writer = bufio.NewWriter(*conn)

	s.conns[conn] = stats
	s.handleIO(conn)
}

func (s *Server) handleIO(conn *net.Conn) {
	// buf := make([]byte, 4096)
	rw := bufio.NewReadWriter(s.conns[conn].reader, s.conns[conn].writer)
	go rw.ReadFrom(os.Stdin)
	go rw.WriteTo(os.Stdout)

}

func (s *Server) listen(proto string, port string) {
	ln, err := net.Listen(proto, ":"+port)
	if err != nil {
		// handle error
	}
	for {
		if len(s.conns) < 1 {
			fmt.Print("[[[", len(s.conns))
			conn, err := ln.Accept()
			if err != nil {
				// handle error
				fmt.Println(err)
			}
			go s.handleConnection(&conn)
		}
	}

}

func main() {
	program := os.Args[0]
	if len(os.Args[:]) < 3 {
		fmt.Println(fmt.Sprintf("usage: %s <protocol: tcp|udp> <port> ", program))
		return
	}
	proto, port := os.Args[1], os.Args[2]
	server := NewServer()

	server.listen(proto, port)
}
