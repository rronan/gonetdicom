package main

// https://dicom.nema.org/dicom/2013/output/chtml/part08/sect_9.3.html

import (
	"fmt"
	"net"
	"os"
)

const (
	HOST = "localhost"
	PORT = "3333"
	TYPE = "tcp"
)

func main() {
	// Listen for incoming connections.
	l, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + HOST + ":" + PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	parsePDU(conn)

	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	// create AARQ var
	aarq := AARQ{}

	// parse AARQ
	aarq = parseAARQ(buf)

	aare := createAARE(aarq)

	// generateAAssociateRJPDU(conn)

	conn.Write([]byte("Message received."))
	// Close the connection when you're done with it.
	conn.Close()
}
