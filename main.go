package main

import (
	"fmt"
	"log"
	"net"

	"github.com/rronan/gonetdicom/storescp"
)

func main() {
	// Create a listener
	listener, err := net.Listen("tcp", ":11112")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	// Listen for connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		// Read data
		data := make([]byte, 1024)
		n, err := conn.Read(data)
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(string(data[:n]))

		assocRQ, _ := storescp.Decode(data[:n])
		fmt.Println(assocRQ)
	}
}
