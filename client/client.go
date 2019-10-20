package client

import (
	"bufio"
	"fmt"
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
		fmt.Fprintf(os.Stderr, "Error connecting to server: %v", err)
		os.Exit(1)
	}

	go client.readConnection(conn)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">> ")
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
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		text := scanner.Text()

		fmt.Printf("%s", text)
	}
}
