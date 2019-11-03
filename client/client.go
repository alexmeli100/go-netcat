package client

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type Client struct {
	peerAddr string
}

func NewClient(peerAddr string) *Client {

	return &Client{peerAddr}
}

func (client *Client) Run() {
	conn, err := net.Dial("tcp", client.peerAddr)

	if err != nil {
		log.Fatal("Error connecting to server: ", err)
	}

	go client.readConnection(conn)

	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		_, err := conn.Write([]byte(text))

		if err != nil {
			fmt.Println("Error writing to stream.")
			break
		}
	}
}

func (client *Client) readConnection(conn net.Conn) {
	defer conn.Close()

	io.Copy(os.Stdout, conn)
}
