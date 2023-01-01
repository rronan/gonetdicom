package storescp

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

func Storescp() {
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

conn:
	for {
		// Read data
		message := make([]byte, 1024)
		_, err := conn.Read(message)
		if err != nil {
			panic(err)
		}
		message = trimMessage(message) // trimmed because the initial byte slice is 1024 bytes long which messes up the decoding
		fmt.Println("The message:", message[:])
		pduType := message[0]

		switch pduType {
		case 0x01:
			fmt.Println("Association request received from Service User")
			AARQStruct, err := DecodeAAssociateRQ(message)
			fmt.Println(AARQStruct.ToString())
			if err != nil {
				panic(err)
			}

			AAACStruct, _ := CreateAssociateAC(AARQStruct)

			ACmessage, _ := EncodeAAssociateAC(AAACStruct)

			n, err := conn.Write(ACmessage[:])
			if err != nil {
				panic(err)
			}
			fmt.Println(n)

			buf2 := make([]byte, 512)
			n2, err := conn.Read(buf2)
			if err != nil {
				panic(err)
			}
			fmt.Println(n2)

		case 0x02:
			fmt.Println("Association accept received from Service User")
		case 0x03:
			fmt.Println("Association reject received from Service User")
		case 0x04:
			fmt.Println("Data received from Service User")
		case 0x05:
			fmt.Println("Release request received from Service User")
		case 0x06:
			fmt.Println("Release response received from Service User")
		case 0x07:
			// TODO: Handle abort
			fmt.Println("Abort received from Service User")
			break conn
		default:
			fmt.Println("Unknown PDU type received from Service User")
		}
	}
}

func trimMessage(message []byte) []byte {
	return message[:6+int(binary.BigEndian.Uint32(message[2:6]))]
}
