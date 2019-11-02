package server

import (
	"fmt"
	"net"
	"sync"
)

type Handler interface {
	Handle(conn net.Conn)
}

type TCPServer struct {
	ch        chan bool
	conn      net.Listener
	h         Handler
	waitGroup *sync.WaitGroup
}

func NewServer(addr string) (*TCPServer, error) {
	listener, err := net.Listen("tcp", addr)

	if err != nil {
		return nil, err
	}

	s := &TCPServer{
		ch:        make(chan bool),
		conn:      listener,
		waitGroup: &sync.WaitGroup{},
	}

	s.waitGroup.Add(1)
	return s, nil
}

func (s *TCPServer) Handle(h Handler) {
	s.h = h
}

func (s *TCPServer) Run() error {
	defer s.waitGroup.Done()

	for {
		select {
		case <-s.ch:
			fmt.Println("Stopping listening on: ", s.conn.Addr())
			return s.close()
		default:
		}

		conn, err := s.conn.Accept()

		if err != nil {
			return err
		}

		s.waitGroup.Add(1)
		go s.serveConn(s.h, conn)
	}
}

func (s *TCPServer) serveConn(h Handler, conn net.Conn) {
	defer s.waitGroup.Done()

	for {
		select {
		case <-s.ch:
			return
		default:
		}

		h.Handle(conn)
	}
}

func (s *TCPServer) close() error {
	return s.conn.Close()
}

func (s *TCPServer) stop() {
	close(s.ch)
	s.waitGroup.Wait()
}
