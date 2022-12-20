package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

const (
	// item-type
	ASSOCIATE_RQ      = 0x01
	ASSOCIATE_AC      = 0x02
	ASSOCIATE_RJ      = 0x03
	P_DATA_TF         = 0x04
	RELEASE_RQ        = 0x05
	RELEASE_RP        = 0x06
	ABORT             = 0x07
	APP_CONTEXT       = 0x10
	PRES_CONTEXT      = 0x20
	ABSTRACT_SYN      = 0x30
	TRANSFER_SYN      = 0x40
	USER_INFO         = 0x50
	PRES_CONTEXT_ITEM = 0x21
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

		// decode message
		aAssRQ := decodeAAssociateRQ(data)

		varItems := variableItems{
			applicationContext: applicationContext{
				itemType:               APP_CONTEXT,
				reserved:               0x00,
				length:                 [2]byte{0x00, 0x10},
				applicationContextName: []byte("1.2.840.10008."),
			},
			presentationContext: []presentationContext{
				{
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
						{
							itemType:           TRANSFER_SYN,
							reserved:           0x00,
							length:             [2]byte{0x00, 0x10},
							transferSyntaxName: []byte("1.2.840.10008.1.2.1"),
						},
					},
				},
			},
		}
		assocAc := associate{
			pduType:         ASSOCIATE_AC,
			reserved:        0x00,
			length:          [4]byte{0x00, 0x00, 0x00, 0x00},
			protocolVersion: [2]byte{0x00, 0x01},
			reserved2:       [2]byte{0x00, 0x00},
			calledAETitle:   getAETitle(data[10:26]),
			callingAETitle:  getAETitle(data[26:42]),
			reserved3:       getBigReserved(data[42:74]),
			variableItems:   varItems,
		}
		// get assocAc length
		alength := len(encodeAssociateAc(assocAc))
		bs := make([]byte, 4)
		binary.BigEndian.PutUint32(bs, uint32(alength-6))
		var bs2 [4]byte
		copy(bs, bs2[:])
		assocAc.length = bs2

		// Write data
		_, err = conn.Write(encodeAssociateAc(assocAc))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func decodeAAssociateRQ(data []byte) associate {
	var a associate
	a.pduType = data[0]
	a.reserved = data[1]
	copy(data[2:6], a.length[:])
	copy(data[6:8], a.protocolVersion[:])
	a.reserved2 = [2]byte{data[8], data[9]}
	a.calledAETitle = getAETitle(data[10:26])
	a.callingAETitle = getAETitle(data[26:42])
	a.reserved3 = getBigReserved(data[42:74])
	a.variableItems = decodeVariableItems(data[74:])
	return assocRQ
}

func encodeAssociateAc(assocAc associate) []byte {
	var data []byte
	data = append(data, assocAc.pduType)
	data = append(data, assocAc.reserved)
	data = append(data, assocAc.length[:]...)
	return data
}

// func computeMessageLength(data []byte) [4]byte {
// 	var arr [4]byte
// 	copy(data[6:10], arr[:])
// 	return arr
// }

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

type associate struct {
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

// type userInfo struct {
// 	itemType uint8
// 	reserved uint8
// 	length   [2]byte
// 	userData []userData
// }

// type userData struct {
// 	itemType  uint8
// 	reserved  uint8
// 	length    [2]byte
// 	maxLength uint32
// }
