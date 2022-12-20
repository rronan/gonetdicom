package storescp

// import (
// 	"bytes"
// 	"encoding/binary"
// 	"net"
// )

// func main() {
// 	listener, err := net.Listen("tcp", ":11112")

// 	if err != nil {
// 		panic(err)
// 	}

// 	for {
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			panic(err)
// 		}
// 		go handleConnection(conn)
// 	}
// }

// func handleConnection(conn net.Conn) {
// 	defer conn.Close()
// 	buf := make([]byte, 211)
// 	buf2 := make([]byte, 512)
// 	buf3 := make([]byte, 512)
// 	buf4 := make([]byte, 512)
// 	buf5 := make([]byte, 512)
// 	buf6 := make([]byte, 512)
// 	for {
// 		n, err := conn.Read(buf)
// 		if err != nil {
// 			panic(err)
// 		}
// 		// for i := 0; i < int(buf[5]+6); i++ {
// 		// 	switch i {
// 		// 	case 0:
// 		// 		fmt.Printf("byte %d = %d This value is the PDU-type\n", i+1, buf[i])
// 		// 	case 1:
// 		// 		fmt.Printf("byte %d = %d This value is reserved and shall be 0\n", i+1, buf[i])
// 		// 	case 2, 3, 4, 5:
// 		// 		fmt.Printf("byte %d = %d This value is the length of the PDU\n", i+1, buf[i])
// 		// 	case 6, 7:
// 		// 		fmt.Printf("byte %d = %d This value is the Protocol Version\n", i+1, buf[i])
// 		// 	default:
// 		// 		fmt.Printf("byte %d = %d\n", i, buf[i])
// 		// 	}
// 		// }

// 		for {
// 			conn2, err := net.Dial("tcp", "localhost:11113")
// 			if err != nil {
// 				panic(err)
// 			}
// 			conn2.Write(buf[:n])
// 			n2, err := conn2.Read(buf2)
// 			if err != nil {
// 				panic(err)
// 			}
// 			conn.Write(buf2[:n2])
// 			n3, err := conn.Read(buf3)
// 			if err != nil {
// 				panic(err)
// 			}
// 			conn2.Write(buf3[:n3])
// 			n4, err := conn2.Read(buf4)
// 			if err != nil {
// 				panic(err)
// 			}
// 			conn.Write(buf4[:n4])
// 			n5, err := conn.Read(buf5)
// 			if err != nil {
// 				panic(err)
// 			}
// 			conn2.Write(buf5[:n5])

// 			n6, err := conn2.Read(buf6)
// 			if err != nil {
// 				panic(err)
// 			}
// 			conn.Write(buf6[:n6])

// 			println(string(buf[:n]))
// 			println(string(buf2[:n2]))
// 			println(string(buf3[:n3]))
// 			println(string(buf4[:n4]))
// 			println(string(buf5[:n5]))
// 			println(string(buf6[:n6]))

// 		}
// 	}
// }

// type A_Associate_RQ_AC_PDU struct {
// 	pduType              byte
// 	reserved             byte
// 	pduLength            uint32
// 	protocolVersion      uint16
// 	reserved2            uint16
// 	calledAETitle        Aetitle
// 	callingAETitle       Aetitle
// 	reserved3            uint32
// 	applicationContext   ApplicationContext
// 	presentationContexts []PresentationContext
// 	userInformation      UserInformation
// }

// func newAAssociateRQAC() *A_Associate_RQ_AC_PDU {
// 	return &A_Associate_RQ_AC_PDU{}
// }

// func newAAssociateRQACFromBytesBuffer(buf bytes.Buffer) *A_Associate_RQ_AC_PDU {

// 	pduType := buf.Next(1)[0]
// 	reserved := buf.Next(1)[0]
// 	pduLength := binary.LittleEndian.Uint32(buf.Next(4))
// 	protocolVersion := binary.LittleEndian.Uint16(buf.Next(2))
// 	reserved2 := binary.LittleEndian.Uint16(buf.Next(2))
// 	calledAETitle := Aetitle{buf.Next(16)}
// 	callingAETitle := Aetitle{buf.Next(16)}
// 	reserved3 := binary.LittleEndian.Uint32(buf.Next(4))

// 	applicationContext := ApplicationContext{
// 		itemType:=  buf.Next(1)[0]
// 		reserved:= buf.Next(1)[0]
// 		itemLength:= binary.LittleEndian.Uint16(buf.Next(2))

// 	}

// 	return &A_Associate_RQ_AC_PDU{
// 		pduType,
// 		reserved,
// 		pduLength,
// 		protocolVersion,
// 		reserved2,
// 		calledAETitle,
// 		callingAETitle,
// 		reserved3,
// 	}

// }

// func newApplicationContext() *ApplicationContext {
// 	return &ApplicationContext{}
// }

// func newApplicationContextFromBytesBuffer(buf bytes.Buffer) *ApplicationContext {
// 	itemType := buf.Next(1)[0]
// 	reserved := buf.Next(1)[0]
// 	itemLength := binary.LittleEndian.Uint16(buf.Next(2))
// 	applicationContextName := buf.Next(int(itemLength))

// 	return &ApplicationContext{
// 		itemType,
// 		reserved,
// 		itemLength,
// 		applicationContextName
// 	}
// }

// type Aetitle struct {
// 	aetitle []byte
// }

// type ApplicationContext struct {
// 	itemType               byte
// 	reserved               byte
// 	itemLength             uint16
// 	applicationContextName string
// }

// type PresentationContext struct {
// 	itemType              byte
// 	reserved              byte
// 	itemLength            uint16
// 	presentationContextID byte
// 	reserved2             byte
// 	reserved3             byte
// 	reserved4             byte
// 	syntaxes              []Syntax
// }

// type Syntax struct {
// 	itemType   byte
// 	reserved   byte
// 	itemLength uint16
// 	syntaxName string
// }

// type UserInformation struct {
// 	itemType   byte
// 	reserved   byte
// 	itemLength uint16
// }
