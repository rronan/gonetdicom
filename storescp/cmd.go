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
		// fmt.Println(string(message[:n]))
		pduType := message[0]

		switch pduType {
		case 0x01:
			fmt.Println("Association request received from Service User")
			message = trimMessage(message) // trimmed because the initial byte slice is 1024 bytes long which messes up the decoding
			AARQStruct, err := DecodeAAssociateRQ(message)
			fmt.Println(AARQStruct.ToString())
			if err != nil {
				panic(err)
			}

			ACSubItem := SubItem{
				itemType:  0x51,
				reserved:  0x00,
				length:    [2]byte{0x00, 0x04},
				maxLength: [4]byte{0x00, 0x00, 0x40, 0x00},
			}
			ACUserInfo := UserInfo{
				itemType: 0x50,
				reserved: 0x00,
				length:   [2]byte{0x00, 0x08},
				subItem:  ACSubItem,
			}
			ACTransferSyntax := TransferSyntax{
				itemType:           0x40,
				reserved:           0x00,
				transferSyntaxName: AARQStruct.variableItems.presentationContextList[0].transferSyntaxList[0].transferSyntaxName,
			}
			binary.BigEndian.PutUint16(ACTransferSyntax.length[:], uint16(len(ACTransferSyntax.transferSyntaxName)))
			ACTransferSyntaxArray := []TransferSyntax{ACTransferSyntax}
			ACAbstractSyntax := AbstractSyntax{
				itemType:           0x30,
				reserved:           0x00,
				abstractSyntaxName: AARQStruct.variableItems.presentationContextList[0].abstractSyntax.abstractSyntaxName,
			}
			binary.BigEndian.PutUint16(ACAbstractSyntax.length[:], uint16(len(ACAbstractSyntax.abstractSyntaxName)))
			ACPresentationContext := PresentationContext{
				itemType:              0x21,
				reserved:              0x00,
				presentationContextID: AARQStruct.variableItems.presentationContextList[0].presentationContextID,
				reserved2:             0x00,
				resultReason:          0x00,
				reserved3:             0x00,
				transferSyntaxList:    ACTransferSyntaxArray,
			}
			binary.BigEndian.PutUint16(ACPresentationContext.length[:], uint16(len(ACTransferSyntax.transferSyntaxName)+8))
			ACApplicationContext := ApplicationContext{
				itemType:               0x10,
				reserved:               0x00,
				applicationContextName: AARQStruct.variableItems.applicationContext.applicationContextName,
			}
			binary.BigEndian.PutUint16(ACApplicationContext.length[:], uint16(len(ACApplicationContext.applicationContextName)))
			ACVariableItems := VariableItems{
				applicationContext: ACApplicationContext,
				presentationContextList: []PresentationContext{
					ACPresentationContext,
				},
				userInfo: ACUserInfo,
			}

			AAACStruct := Associate{
				pduType:         0x02,
				reserved:        0x00,
				protocolVersion: [2]byte{0x00, 0x01},
				reserved2:       [2]byte{0x00, 0x00},
				calledAETitle:   AARQStruct.calledAETitle,
				callingAETitle:  AARQStruct.callingAETitle,
				reserved3:       [32]byte{0x00},
				variableItems:   ACVariableItems,
			}

			mint := binary.BigEndian.Uint16(ACVariableItems.applicationContext.length[:])
			mint2 := binary.BigEndian.Uint16(ACVariableItems.presentationContextList[0].length[:])
			mint3 := binary.BigEndian.Uint16(ACVariableItems.userInfo.length[:])

			binary.BigEndian.PutUint32(AAACStruct.length[:], uint32(mint+mint2+mint3+2+2+16+16+32))
			fmt.Println(AAACStruct.ToString())

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
