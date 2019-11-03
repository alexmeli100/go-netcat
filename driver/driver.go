package driver

import (
	"bufio"
	"fmt"
	"github.com/alexmeli100/go-netcat/client"
	"github.com/alexmeli100/go-netcat/config"
	"github.com/alexmeli100/go-netcat/server"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

type Driver struct {
	params *config.Opts
}

func (d *Driver) run() {
	if !d.params.Listen {
		peerAddr := fmt.Sprintf("%s:%d", d.params.Addr, d.params.Port)
		c := client.NewClient(peerAddr)
		c.Run()
	}

	serverAddr := fmt.Sprintf("%s:%d", d.params.Addr, d.params.Port)

	s, err := server.NewServer(serverAddr)

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
	params *config.Opts
}

func (s ServerHandler) Handle(conn net.Conn) {
	s.conn = conn

	if s.params.Execute != "" {
		output, _ := runCommand(s.params.Execute)
		conn.Write(output)
	}

	if s.params.Command {
		s.executeShell()
	}
}

func (s *ServerHandler) executeShell() {
	host, _ := os.Hostname()
	prompt := []byte(fmt.Sprintf("%s>", host))
	rw := bufio.NewReadWriter(bufio.NewReader(s.conn), bufio.NewWriter(s.conn))

	for {
		rw.Write(prompt)

		cmd, _ := rw.ReadString('\n')
		response, _ := runCommand(cmd)
		rw.Write(response)
	}
}

func runCommand(cmd string) ([]byte, error) {
	cmd = strings.TrimSpace(cmd)
	shellCmd := exec.Command(cmd)

	out, err := shellCmd.CombinedOutput()

	if err != nil {
		return nil, err
	}

	return out, nil
}
