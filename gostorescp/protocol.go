package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

// struct for A-ASSOCIATE-RQ PDU (AARQ)
type AARQ struct {
	pduType              uint8
	reserved1            uint8
	MaxPDULength         uint32
	protocolVersion      uint8
	reserved2            uint8
	calledAETitle        string
	callingAETitle       string
	reserved3            uint8
	applicationContext   applicationContext
	presentationContexts []presentationContext
	userInformation      userInformation
}

// struct for A-ASSOCIATE-AC PDU (A-ASSOCIATE-AC)
type AARE struct {
	pduType              uint8
	reserved1            uint8
	MaxPDULength         uint32
	protocolVersion      uint8
	reserved2            uint8
	calledAETitle        string
	callingAETitle       string
	reserved3            uint8
	applicationContext   applicationContext
	presentationContexts []presentationContext
	userInformation      userInformation
}

// struct for A-ASSOCIATE-RJ PDU (A-ASSOCIATE-RJ)
type AARJ struct {
	pduType   uint8
	reserved1 uint8
	result    uint8
	source    uint8
	reason    uint8
}

// struct for P-DATA-TF PDU (P-DATA-TF)
type PDATA struct {
	pduType    uint8
	reserved1  uint8
	dataValues []presentationDataValue
}

// struct for A-RELEASE-RQ PDU (A-RELEASE-RQ)
type RLRQ struct {
	pduType   uint8
	reserved1 uint8
	pduLength uint32
	reserved2 uint8
}

// struct for A-RELEASE-RP PDU (A-RELEASE-RP)
type RLRE struct {
	pduType   uint8
	reserved1 uint8
	pduLength uint32
	reserved2 uint8
}

// struct for A-ABORT PDU (A-ABORT)
type ABRT struct {
	pduType   uint8
	reserved1 uint8
	pduLength uint32
	reserved2 uint8
	reserved3 uint8
	source    uint8
	reason    uint8
}

// struct for ApplicationContext
type applicationContext struct {
	itemType               uint8
	reserved1              uint8
	itemLength             uint8
	applicationContextName string
}

// struct for PresentationContext
type presentationContext struct {
	itemType              uint8
	reserved1             uint8
	itemLength            uint8
	presentationContextID uint8
	reserved2             uint8
	reserved3             uint8
	reserved4             uint8
	abstractSyntax        abstractSyntax
	transferSyntaxes      []transferSyntax
}

// struct for AbstractSyntax
type abstractSyntax struct {
	itemType           uint8
	reserved1          uint8
	itemLength         uint8
	abstractSyntaxName string
}

// struct for TransferSyntax
type transferSyntax struct {
	itemType           uint8
	reserved1          uint8
	itemLength         uint8
	transferSyntaxName string
}

// struct for UserInformation
type userInformation struct {
	itemType      uint8
	reserved1     uint8
	itemLength    uint8
	userDataItems []userDataItem
}

// struct for UserDataItem
type userDataItem struct {
	itemType   uint8
	reserved1  uint8
	itemLength uint8
	data       []byte
}

// struct for PresentationDataValue
type presentationDataValue struct {
	presentationContextID uint8
	reserved1             uint8
	presentationDataValue uint8
	data                  []byte
}

// parse the DICOM PDU
func parsePDU(r io.Reader) (err error) {
	// read the PDU header
	var pduHeader [6]byte
	_, err = io.ReadFull(r, pduHeader[:])
	if err != nil {
		return
	}

	// read the PDU type
	pduType := pduHeader[0]
	fmt.Printf("PDU type: %02X	", pduType)

	// read the reserved byte
	reserved1 := pduHeader[1]

	// read the PDU length
	pduLength := binary.BigEndian.Uint32(pduHeader[2:])
	fmt.Printf("PDU length: %d	", pduLength)

	// read the PDU body
	var pduBody = make([]byte, pduLength)
	_, err = io.ReadFull(r, pduBody)
	if err != nil {
		return
	}

	// parse the PDU body
	switch pduType {
	case 0x01:
		err = parseAAssociateRQ(pduBody)
	// case 0x02:
	// 	err = parseAAssociateAC(pduBody)
	// case 0x03:
	// 	err = parseAAssociateRJ(pduBody)
	// case 0x04:
	// 	err = parsePDataTF(pduBody)
	// case 0x05:
	// 	err = parseAReleaseRQ(pduBody)
	// case 0x06:
	// 	err = parseAReleaseRP(pduBody)
	// case 0x07:
	// 	err = parseAAAbort(pduBody)
	default:
		err = fmt.Errorf("unknown PDU type: %02X", pduType)
	}

	return
}

// parse the A-ASSOCIATE-RQ PDU
func parseAAssociateRQ(pduBody []byte) (err error) {

	// read the protocol version
	protocolVersion := binary.BigEndian.Uint16(pduBody[0:])
	fmt.Printf("	Protocol version: %04X", protocolVersion)

	// read the called AE title
	calledAETitle := string(bytes.Trim(pduBody[4:19], "\x20"))
	fmt.Printf("	Called AE title: %s", calledAETitle)

	// read the calling AE title
	callingAETitle := string(bytes.Trim(pduBody[20:35], "\x20"))
	fmt.Printf("	Calling AE title: %s", callingAETitle)

	// parse the variable items
	variableItems := pduBody[68:]
	err = parseVariableItems(variableItems)

	return
}

// parse the variable items
func parseVariableItems(variableItems []byte) (err error) {
	// create a reader for the variable items
	r := bytes.NewReader(variableItems)

	// read the items
	for {
		// read the item header
		var itemHeader [4]byte
		_, err = io.ReadFull(r, itemHeader[:])
		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}

		// read the item type
		itemType := itemHeader[0]
		fmt.Printf("	Item type: %02X	", itemType)

		// read the item length
		itemLength := binary.BigEndian.Uint16(itemHeader[2:])
		fmt.Printf("	Item length: %d", itemLength)

		// read the item body
		var itemBody = make([]byte, itemLength)
		_, err = io.ReadFull(r, itemBody)
		if err != nil {
			return
		}

		// parse the item body
		switch itemType {
		case 0x10:
			err = parseApplicationContextItem(itemBody)
		case 0x20:
			err = parsePresentationContextItem(itemBody)
		case 0x50:
			err = parseUserInformationItem(itemBody)
		default:
			err = fmt.Errorf("unknown item type: %02X", itemType)
		}
	}

	return
}

// parse the Application Context Item
func parseApplicationContextItem(itemBody []byte) (err error) {
	// read the application context name
	applicationContextName := string(itemBody)
	fmt.Printf("	Application context name: %s", applicationContextName)

	return
}

// parse the Presentation Context Item
func parsePresentationContextItem(itemBody []byte) (err error) {
	// read the presentation context ID
	presentationContextID := itemBody[0]
	fmt.Printf("	Presentation context ID: %02X", presentationContextID)

	// read the abstract syntax name
	contextItems := itemBody[4:]
	r := bytes.NewReader(contextItems)

	// read sub items
	for {
		// read the item header
		var itemHeader [4]byte
		_, err = io.ReadFull(r, itemHeader[:])
		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}

		// read the item type
		itemType := itemHeader[0]
		fmt.Printf("	Item type: %02X	", itemType)

		// read the item length
		itemLength := binary.BigEndian.Uint16(itemHeader[2:])
		fmt.Printf("	Item length: %d", itemLength)

		// read the item body
		var itemBody = make([]byte, itemLength)
		_, err = io.ReadFull(r, itemBody)
		if err != nil {
			return
		}

		// parse the item body
		switch itemType {
		case 0x30:
			err = parseAbstractSyntax(itemBody)
		case 0x40:
			err = parseTransferSyntax(itemBody)
		default:
			err = fmt.Errorf("unknown item type: %02X", itemType)
		}
	}

	return
}

// parse the Abstract Syntax
func parseAbstractSyntax(itemBody []byte) (err error) {
	// read the abstract syntax name
	abstractSyntaxName := string(itemBody)
	fmt.Printf("	Abstract syntax name: %s", abstractSyntaxName)

	return
}

// parse the Transfer Syntax
func parseTransferSyntax(itemBody []byte) (err error) {
	// read the transfer syntax name
	transferSyntaxName := string(itemBody)
	fmt.Printf("	Transfer syntax name: %s", transferSyntaxName)

	return
}

// parse the User Information Item
func parseUserInformationItem(itemBody []byte) (err error) {
	// parse the sub-items
	userInformation := itemBody
	r := bytes.NewReader(userInformation)

	for {
		// read the item header
		var itemHeader [4]byte
		_, err = io.ReadFull(r, itemHeader[:])
		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}

		// read the item type
		itemType := itemHeader[0]
		fmt.Printf("	Item type: %02X	", itemType)

		// read the item length
		itemLength := binary.BigEndian.Uint16(itemHeader[2:])
		fmt.Printf("	Item length: %d", itemLength)

		// read the item body
		var itemBody = make([]byte, itemLength)
		_, err = io.ReadFull(r, itemBody)
		if err != nil {
			return
		}

		// parse the item body
		switch itemType {
		case 0x51:
			err = parseMaxLengthReceived(itemBody)
		case 0x52:
			err = parseImplementationClassUID(itemBody)
		case 0x53:
			err = parseAsynchronousOperationsWindowSubItem(itemBody)
		case 0x54:
			err = parseSCPSCURoleSelectionSubItem(itemBody)
		case 0x55:
			err = parseImplementationVersionName(itemBody)
		case 0x56:
			err = parseSOPClassExtendedNegotiationSubItem(itemBody)
		case 0x57:
			err = parseSOPClassCommonExtendedNegotiationSubItem(itemBody)
		case 0x58:
			err = parseUserIdentityNegociationSubItem(itemBody)
		default:
			err = fmt.Errorf("unknown item type: %02X", itemType)
		}
	}

	return
}

// parse the Maximum Length Sub-Item
func parseMaxLengthReceived(itemBody []byte) (err error) {
	// read the maximum length received
	maximumLengthReceived := binary.BigEndian.Uint32(itemBody)
	fmt.Printf("	Maximum length received: %d", maximumLengthReceived)

	return
}

// parse the Implementation Class UID Sub-Item
func parseImplementationClassUID(itemBody []byte) (err error) {
	// read the implementation class UID
	implementationClassUID := string(itemBody)
	fmt.Printf("	Implementation class UID: %s", implementationClassUID)

	return
}

// parse the Asynchronous Operations Window Sub-Item
func parseAsynchronousOperationsWindowSubItem(itemBody []byte) (err error) {
	// read the maximum number operations invoked
	maximumNumberOperationsInvoked := binary.BigEndian.Uint16(itemBody)
	fmt.Printf("	Maximum number operations invoked: %d", maximumNumberOperationsInvoked)

	// read the maximum number operations performed
	maximumNumberOperationsPerformed := binary.BigEndian.Uint16(itemBody[2:])
	fmt.Printf("	Maximum number operations performed: %d", maximumNumberOperationsPerformed)

	return
}

// parse the SCP/SCU Role Selection Sub-Item
func parseSCPSCURoleSelectionSubItem(itemBody []byte) (err error) {
	// read the SOP class UID
	sopClassUID := string(itemBody[0:16])
	fmt.Printf("	SOP class UID: %s", sopClassUID)

	// read the SCP/SCU role
	scpSCURole := itemBody[16]
	fmt.Printf("	SCP/SCU role: %02X", scpSCURole)

	return
}

// parse the SOP Class Extended Negotiation Sub-Item
func parseSOPClassExtendedNegotiationSubItem(itemBody []byte) (err error) {
	// read the SOP class UID length
	sopClassUIDLength := binary.BigEndian.Uint16(itemBody[0:1])
	fmt.Printf("	SOP class UID length:	%d", sopClassUIDLength)

	// read the SOP class UID
	sopClassUID := string(itemBody[2:sopClassUIDLength])
	fmt.Printf("	SOP class UID: %s", sopClassUID)

	// read the service class application
	serviceClassApplication := itemBody[sopClassUIDLength+1:]
	fmt.Printf("	Service class application: %02X", serviceClassApplication)

	return
}

// parse the SOP Class Common Extended Negotiation Sub-Item
func parseSOPClassCommonExtendedNegotiationSubItem(itemBody []byte) (err error) {
	// read the SOP class UID length
	sopClassUIDLength := binary.BigEndian.Uint16(itemBody[0:1])
	fmt.Printf("	SOP class UID length:	%d", sopClassUIDLength)

	// read the SOP class UID
	sopClassUID := string(itemBody[2:sopClassUIDLength])
	fmt.Printf("	SOP class UID: %s", sopClassUID)

	// read the service class uid length
	serviceClassUIDLength := binary.BigEndian.Uint16(itemBody[sopClassUIDLength+1 : sopClassUIDLength+2])
	fmt.Printf("	Service class UID length:	%d", serviceClassUIDLength)

	// read the service class UID
	serviceClassUID := string(itemBody[sopClassUIDLength+3 : sopClassUIDLength+serviceClassUIDLength])
	fmt.Printf("	Service class UID: %s", serviceClassUID)

	// read related general SOP class identification length
	relatedGeneralSOPClassIdentificationLength := binary.BigEndian.Uint16(itemBody[serviceClassUIDLength+1 : serviceClassUIDLength+2])
	fmt.Printf("	Related general SOP class identification length:	%d", relatedGeneralSOPClassIdentificationLength)

	// read the related general SOP class identification
	relatedGeneralSOPClassIdentification := string(itemBody[serviceClassUIDLength+3 : relatedGeneralSOPClassIdentificationLength])
	fmt.Printf("	Related general SOP class identification: %s", relatedGeneralSOPClassIdentification)

	return
}

func parseImplementationVersionName(itemBody []byte) (err error) {
	// read the implementation version name
	implementationVersionName := string(itemBody)
	fmt.Printf("	Implementation version name: %s", implementationVersionName)

	return
}

func parseUserIdentityNegociationSubItem(itemBody []byte) (err error) {
	// read the user identity type
	userIdentityType := itemBody[0]
	fmt.Printf("	User identity type: %02X", userIdentityType)

	// read the primary field
	primaryField := string(itemBody[1:])
	fmt.Printf("	Primary field: %s", primaryField)

	return
}

// generate a A-ASSOCIATE-RJ PDU
func generateAAssociateRJPDU(connection net.Conn) {
	// create a buffer to write the PDU into
	buffer := new(bytes.Buffer)

	// write the PDU type
	pduType := byte(0x03)
	binary.Write(buffer, binary.BigEndian, pduType)

	// write the reserved byte
	reserved := byte(0x00)
	binary.Write(buffer, binary.BigEndian, reserved)

	// write the length
	length := uint32(0x00000004)
	binary.Write(buffer, binary.BigEndian, length)

	// write the result
	result := uint8(0x01)
	binary.Write(buffer, binary.BigEndian, result)

	// write the source
	source := uint8(0x02)
	binary.Write(buffer, binary.BigEndian, source)

	// write the reason
	reason := uint8(0x01)
	binary.Write(buffer, binary.BigEndian, reason)

	// write the PDU to the connection
	connection.Write(buffer.Bytes())

	return
}

// create A-ASSOCIATE-AC PDU from A-ASSOCIATE-RQ PDU
func createAARE(aarq AARQ) AARE {
	aare := AARE{}
	aare.pduType = 2
	aare.reserved1 = 0
	aare.MaxPDULength = aarq.MaxPDULength
	aare.protocolVersion = aarq.protocolVersion
	aare.reserved2 = 0
	aare.calledAETitle = aarq.calledAETitle
	aare.callingAETitle = aarq.callingAETitle
	aare.reserved3 = 0
	aare.applicationContext = aarq.applicationContext
	aare.presentationContexts = aarq.presentationContexts
	aare.userInformation = aarq.userInformation
	return aare
}
