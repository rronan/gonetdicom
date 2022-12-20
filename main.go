package main

import (
	"fmt"
	"log"
	"net"
)

const (
	// item-type
	ASSOCIATE_RQ      byte = 0x01
	ASSOCIATE_AC           = 0x02
	ASSOCIATE_RJ           = 0x03
	P_DATA_TF              = 0x04
	RELEASE_RQ             = 0x05
	RELEASE_RP             = 0x06
	ABORT                  = 0x07
	APP_CONTEXT            = 0x10
	PRES_CONTEXT           = 0x20
	ABSTRACT_SYN           = 0x30
	TRANSFER_SYN           = 0x40
	USER_INFO              = 0x50
	PRES_CONTEXT_ITEM      = 0x21
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
		fmt.Println(string(data[:n]))
		varItems := variableItems{
			applicationContext: applicationContext{
				itemType:               APP_CONTEXT,
				reserved:               0x00,
				length:                 [2]byte{0x00, 0x10},
				applicationContextName: []byte("1.2.840.10008."),
			},
			presentationContext: []presentationContext{
				presentationContext{
					itemType:              PRES_CONTEXT,
					reserved:              0x00,
					length:                [2]byte{0x00, 0x00},
					presentationContextID: 0x01,
					reserved2:             0x00,
					reserved3:             0x00,
					reserved4:             0x00,
					abstractSyntax: abstractSyntax{
						itemType:           ABSTRACT_SYN,
						reserved:           0x00,
						length:             [2]byte{0x00, 0x10},
						abstractSyntaxName: []byte("1.2.840.10008.1.1"),
					},
					transferSyntax: []transferSyntax{
						transferSyntax{
							itemType:           TRANSFER_SYN,
							reserved:           0x00,
							length:             [2]byte{0x00, 0x10},
							transferSyntaxName: []byte("1.2.840.10008.1.2.1"),
						},
					},
				},
			},
		}
		assocAc := associate_ac_struct{
			pduType:         ASSOCIATE_AC,
			reserved:        0x00,
			length:          [4]byte{0x00, 0x00, 0x00, 0x00},
			protocolVersion: [2]byte(0x0001),
			reserved2:       0x0000,
			calledAETitle:   getAETitle(data[10:26]),
			callingAETitle:  getAETitle(data[26:42]),
			reserved3:       getBigReserved(data[42:74]),
			variableItems:   varItems,
		}
		// Write data
		_, err = conn.Write(assocAc)
	}
}

func encodeAssociateAc(assocAc associate_ac_struct) []byte {
	var data []byte
	return data
}

func getAETitle(data []byte) [16]byte {
	var arr [16]byte
	copy(data[10:26], arr[:])
	return arr
}

func getBigReserved(data []byte) [32]byte {
	var arr [32]byte
	copy(data[42:74], arr[:])
	return arr
}

type associate_ac_struct struct {
	pduType         uint8
	reserved        uint8
	length          [4]byte
	protocolVersion [2]byte
	reserved2       [2]byte
	calledAETitle   [16]byte
	callingAETitle  [16]byte
	reserved3       [32]byte
	variableItems   variableItems
}

type variableItems struct {
	applicationContext  applicationContext
	presentationContext []presentationContext
}

type applicationContext struct {
	itemType               uint8
	reserved               uint8
	length                 [2]byte
	applicationContextName []byte
}

type presentationContext struct {
	itemType              uint8
	reserved              uint8
	length                [2]byte
	presentationContextID uint8
	reserved2             uint8
	reserved3             uint8
	reserved4             uint8
	abstractSyntax        abstractSyntax
	transferSyntax        []transferSyntax
}

type abstractSyntax struct {
	itemType           uint8
	reserved           uint8
	length             [2]byte
	abstractSyntaxName []byte
}

type transferSyntax struct {
	itemType           uint8
	reserved           uint8
	length             [2]byte
	transferSyntaxName []byte
}

type userInfo struct {
	itemType uint8
	reserved uint8
	length   [2]byte
	userData []userData
}

type userData struct {
	itemType  uint8
	reserved  uint8
	length    [2]byte
	maxLength uint32
}
