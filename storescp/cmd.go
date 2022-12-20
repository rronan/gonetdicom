package storescp

import (
	"encoding/binary"
	"fmt"
	"net"
)

func Storescp() {
	listener, err := net.Listen("tcp", ":11112")

	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {

	defer conn.Close()

	echoBiteSlice := make([]byte, 1024)

	for {
		_, err := conn.Read(echoBiteSlice)
		if err != nil {
			panic(err)
		}
		fmt.Println("Handling buffer:")
		pduType := echoBiteSlice[0]
		fmt.Println("PDU type:", pduType)
		reserved := echoBiteSlice[1]
		fmt.Println("Reserved:", reserved)
		pduLength := echoBiteSlice[2:6]
		fmt.Println("PDU length:", binary.BigEndian.Uint32(pduLength))
		remainingLength := binary.BigEndian.Uint32(pduLength)
		fmt.Println("Remaining length:", remainingLength)
		protocolVersion := echoBiteSlice[6:8]
		remainingLength -= 2
		fmt.Println("Protocol version:", binary.BigEndian.Uint16(protocolVersion))
		reserved2 := echoBiteSlice[8:10]
		remainingLength -= 2
		fmt.Println("Reserved2:", binary.BigEndian.Uint16(reserved2))
		calledAETitle := echoBiteSlice[10:26]
		remainingLength -= 16
		fmt.Println("Called AE title:", string(calledAETitle))
		callingAETitle := echoBiteSlice[26:42]
		remainingLength -= 16
		fmt.Println("Calling AE title:", string(callingAETitle))
		reserved3 := echoBiteSlice[42:74]
		remainingLength -= 32
		fmt.Println("Reserved3:", string(reserved3))
		variableItems := echoBiteSlice[74 : 74+remainingLength]
		varItemCursor := 74
		fmt.Println("Variable items:", string(variableItems))
		fmt.Println("===============================")
		fmt.Println("Application context subItems:")
		applicationContextLength := applicationContextHandler(variableItems)
		varItemCursor += applicationContextLength
		remainingLength -= uint32(applicationContextLength)
		fmt.Println("===============================")
		fmt.Println("Presentation context subItems:")
		presentationContextLength := presentationContextHandler(variableItems[applicationContextLength:])
		varItemCursor += presentationContextLength
		remainingLength -= uint32(presentationContextLength)
		fmt.Println("===============================")
		fmt.Println("User information subItems:")
		userInformationLength := userInformationHandler(variableItems[applicationContextLength+presentationContextLength:])
		varItemCursor += userInformationLength
		remainingLength -= uint32(userInformationLength)
		implementationClassUIDLength := implementationClassUIDHandler(variableItems[applicationContextLength+presentationContextLength+userInformationLength:])
		varItemCursor += implementationClassUIDLength
		remainingLength -= uint32(implementationClassUIDLength)
		implementationVersionNameLength := implementationVersionNameHandler(variableItems[applicationContextLength+presentationContextLength+userInformationLength+implementationClassUIDLength:])
		varItemCursor += implementationVersionNameLength
		remainingLength -= uint32(implementationVersionNameLength)

		// fmt.Printf("%s", echoBiteSlice[:n])
	}
}

func applicationContextHandler(byteSlice []byte) int {
	itemType := byteSlice[0]
	fmt.Println("	Application context item type:", itemType)
	reserved := byteSlice[1]
	fmt.Println("	Application context reserved:", reserved)
	applicationContextItemLength := byteSlice[2:4]
	fmt.Println("	Application context item length:", binary.BigEndian.Uint16(applicationContextItemLength))
	applicationContextRemainingLength := binary.BigEndian.Uint16(applicationContextItemLength)
	fmt.Println("	Application context remaining length:", applicationContextRemainingLength)
	applicationContextName := byteSlice[4 : 4+applicationContextRemainingLength]
	fmt.Println("	Application context name:", string(applicationContextName))
	return 4 + int(binary.BigEndian.Uint16(applicationContextItemLength))
}

func presentationContextHandler(byteSlice []byte) int {
	itemType := byteSlice[0]
	fmt.Println("	Presentation context item type:", itemType)
	reserved := byteSlice[1]
	fmt.Println("	Presentation context reserved:", reserved)
	presentationContextItemLength := byteSlice[2:4]
	fmt.Println("	Presentation context item length:", binary.BigEndian.Uint16(presentationContextItemLength))
	presentationContextRemainingLength := binary.BigEndian.Uint16(presentationContextItemLength)
	fmt.Println("	Presentation context remaining length:", presentationContextRemainingLength)
	presentationContextID := byteSlice[4]
	fmt.Println("	Presentation context ID:", presentationContextID)
	presentationContextRemainingLength -= 1
	reserved2 := byteSlice[5]
	fmt.Println("	Presentation context reserved2:", reserved2)
	reserved3 := byteSlice[6]
	fmt.Println("	Presentation context reserved3:", reserved3)
	reserved4 := byteSlice[7]
	fmt.Println("	Presentation context reserved4:", reserved4)
	presentationContextRemainingLength -= 3

	syntaxSubItemLength := syntaxHandler(byteSlice[8:])
	presentationContextRemainingLength -= uint16(syntaxSubItemLength)
	syntaxSubItemLength = syntaxHandler(byteSlice[8+syntaxSubItemLength:])
	presentationContextRemainingLength -= uint16(syntaxSubItemLength)

	return 8 + int(binary.BigEndian.Uint16(presentationContextItemLength))
}

func syntaxHandler(byteSlice []byte) int {
	itemType := byteSlice[0]
	if itemType == 0x30 {
		fmt.Println("		==> Abstract syntax <==")
	} else if itemType == 0x40 {
		fmt.Println("		==> Transfer syntax <==")
	}
	fmt.Println("		Syntax subItem item type:", itemType)
	reserved := byteSlice[1]
	fmt.Println("		Syntax subItem reserved:", reserved)
	syntaxSubItemLength := byteSlice[2:4]
	fmt.Println("		Syntax subItem item length:", binary.BigEndian.Uint16(syntaxSubItemLength))
	syntaxSubItemRemainingLength := binary.BigEndian.Uint16(syntaxSubItemLength)
	fmt.Println("		Syntax subItem remaining length:", syntaxSubItemRemainingLength)
	syntaxSubItemName := byteSlice[4 : 4+syntaxSubItemRemainingLength]
	fmt.Println("		Syntax subItem name:", string(syntaxSubItemName))
	return 4 + int(binary.BigEndian.Uint16(syntaxSubItemLength))
}

func userInformationHandler(byteSlice []byte) int {
	itemType := byteSlice[0]
	fmt.Println("	User information item type:", itemType)
	reserved := byteSlice[1]
	fmt.Println("	User information reserved:", reserved)
	userInformationItemLength := byteSlice[2:4]
	fmt.Println("	User information item length:", binary.BigEndian.Uint16(userInformationItemLength))
	userInformationRemainingLength := binary.BigEndian.Uint16(userInformationItemLength)
	fmt.Println("	User information remaining length:", userInformationRemainingLength)
	maximumLengthSubItemLength := binary.BigEndian.Uint16(byteSlice[4:8])
	fmt.Println("	Maximum length subItem length:", maximumLengthSubItemLength)
	return 4 + int(binary.BigEndian.Uint16(userInformationItemLength))
}

func implementationClassUIDHandler(byteSlice []byte) int {
	itemType := byteSlice[0]
	fmt.Println("	Implementation class UID item type:", itemType)
	reserved := byteSlice[1]
	fmt.Println("	Implementation class UID reserved:", reserved)
	implementationClassUIDItemLength := byteSlice[2:4]
	fmt.Println("	Implementation class UID item length:", binary.BigEndian.Uint16(implementationClassUIDItemLength))
	implementationClassUIDRemainingLength := binary.BigEndian.Uint16(implementationClassUIDItemLength)
	fmt.Println("	Implementation class UID remaining length:", implementationClassUIDRemainingLength)
	implementationClassUID := byteSlice[4 : 4+implementationClassUIDRemainingLength]
	fmt.Println("	Implementation class UID:", string(implementationClassUID))
	return 4 + int(binary.BigEndian.Uint16(implementationClassUIDItemLength))
}

func implementationVersionNameHandler(byteSlice []byte) int {
	itemType := byteSlice[0]
	fmt.Println("	Implementation version name item type:", itemType)
	reserved := byteSlice[1]
	fmt.Println("	Implementation version name reserved:", reserved)
	implementationVersionNameItemLength := byteSlice[2:4]
	fmt.Println("	Implementation version name item length:", binary.BigEndian.Uint16(implementationVersionNameItemLength))
	implementationVersionNameRemainingLength := binary.BigEndian.Uint16(implementationVersionNameItemLength)
	fmt.Println("	Implementation version name remaining length:", implementationVersionNameRemainingLength)
	implementationVersionName := byteSlice[4 : 4+implementationVersionNameRemainingLength]
	fmt.Println("	Implementation version name:", string(implementationVersionName))
	return 4 + int(binary.BigEndian.Uint16(implementationVersionNameItemLength))
}
