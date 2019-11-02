package driver

import (
	"fmt"
	"github.com/alexmeli100/go-netcat/client"
	"github.com/alexmeli100/go-netcat/config"
	"github.com/alexmeli100/go-netcat/server"
	"log"
	"net"
	"os/exec"
	"strings"
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

	err = s.Run()
}

type ServerHandler struct {
	conn   net.Conn
	params *config.Opts
}

func (s ServerHandler) Handle(conn net.Conn) {
	s.conn = conn

	if s.params.UploadDestination != "" {
		err := s.uploadFile(s.params.UploadDestination)

		if err != nil {
			log.Fatal("Failed to upload file")
		}
	}

	if s.params.Execute != "" {

	}
}

func (s ServerHandler) runCommand(cmd string) (error, []byte) {
	cmd = strings.TrimSpace(cmd)
	shellCmd := exec.Command(cmd)

	out, err := shellCmd.CombinedOutput()

	if err != nil {
		return err, nil
	}

	return nil, out
}

func (s ServerHandler) uploadFile(dest string) error {

}
