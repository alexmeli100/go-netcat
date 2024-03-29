package driver

import (
	"bufio"
	"fmt"
	"github.com/alexmeli100/go-netcat/client"
	"github.com/alexmeli100/go-netcat/config"
	"github.com/alexmeli100/go-netcat/server"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

type Driver struct {
	Params *config.Opts
}

func (d *Driver) Run() {
	if !d.Params.Listen {
		peerAddr := fmt.Sprintf("%s:%d", d.Params.Addr, d.Params.Port)
		c := client.NewClient(peerAddr)
		c.Run()

		return
	}

	serverAddr := fmt.Sprintf("%s:%d", d.Params.Addr, d.Params.Port)

	s, err := server.NewServer(serverAddr)
	handler := ServerHandler{Params: d.Params}
	s.Handle(handler)

	if err != nil {
		log.Fatal("Error creating server")
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	s.Stop()
}

type ServerHandler struct {
	conn   net.Conn
	Params *config.Opts
}

func (s ServerHandler) Handle(conn net.Conn) {
	s.conn = conn

	if s.Params.Execute != "" {
		output, _ := runCommand(s.Params.Execute)
		conn.Write(output)
	}

	if s.Params.Command {
		s.executeShell()
	}
}

func (s *ServerHandler) executeShell() {
	host, _ := os.Hostname()
	prompt := fmt.Sprintf("%s>", host)
	rw := bufio.NewReadWriter(bufio.NewReader(s.conn), bufio.NewWriter(s.conn))

	for {
		_, err := rw.WriteString(prompt)

		if err != nil {
			break
		}
		rw.Flush()

		cmd, err := rw.ReadString('\n')

		if err == io.EOF {
			break
		}

		response, _ := runCommand(cmd)
		log.Println("loop")

		_, err = rw.Write(response)

		if err != nil {
			break
		}
		rw.Flush()
	}
}

func runCommand(cmd string) ([]byte, error) {
	cmd = strings.TrimSpace(cmd)
	shellCmd := exec.Command("cmd", "/C", cmd)

	out, err := shellCmd.CombinedOutput()

	if err != nil {
		return nil, err
	}

	return out, nil
}
