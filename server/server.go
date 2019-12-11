package server

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type Handler interface {
	Handle(conn net.Conn)
}

type TCPServer struct {
	ch        chan bool
	conn      *net.TCPListener
	h         Handler
	waitGroup *sync.WaitGroup
}

func NewServer(addr string) (*TCPServer, error) {
	laddr, err := net.ResolveTCPAddr("tcp", addr)

	if err != nil {
		return nil, err
	}

	listener, err := net.ListenTCP("tcp", laddr)

	if err != nil {
		return nil, err
	}

	s := &TCPServer{
		ch:        make(chan bool),
		conn:      listener,
		waitGroup: &sync.WaitGroup{},
	}

	s.waitGroup.Add(1)
	go s.Run()
	return s, nil
}

func (s *TCPServer) Handle(h Handler) {
	s.h = h
}

func (s *TCPServer) Run() {
	defer s.waitGroup.Done()

	for {
		select {
		case <-s.ch:
			fmt.Println("Stopping listening on: ", s.conn.Addr())
			s.close()
			return
		default:
		}

		s.conn.SetDeadline(time.Now().Add(1e9))
		conn, err := s.conn.Accept()

		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			log.Println("Failed to accept connection:", err.Error())
		}

		log.Println("Received connection from: ", s.conn.Addr())

		s.waitGroup.Add(1)
		go s.serveConn(s.h, conn)
	}
}

func (s *TCPServer) serveConn(h Handler, conn net.Conn) {
	defer conn.Close()
	defer s.waitGroup.Done()

	h.Handle(conn)
}

func (s *TCPServer) close() error {
	return s.conn.Close()
}

func (s *TCPServer) Stop() {
	close(s.ch)
	s.waitGroup.Wait()
}
