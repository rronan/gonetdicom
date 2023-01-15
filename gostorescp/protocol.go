package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// struct to manage A-ASSOCIATE-RQ PDU
type AAssociateRQ struct {
	ProtocolVersion        uint16
	CalledAETitle          string
	CallingAETitle         string
	Reserved               uint32
	ApplicationContextName string
	PresentationContexts   []PresentationContext
	UserInformation        UserInformation
}

// struct to manage Presentation Context
type PresentationContext struct {
	PresentationContextID uint8
	Reserved              uint8
	AbstractSyntax        string
	TransferSyntaxes      []string
}

// struct to manage User Information
type UserInformation struct {
	UserDataItems []UserDataItem
}

// struct to manage User Data Item
type UserDataItem struct {
	Type uint8
	Data []byte
}

// struct to manage SCP/SCU role selection
type SCPSCURoleSelection struct {
	RoleSelectionUID string
	SCURole          uint8
	SCPRole          uint8
}

// parse A-ASSOCIATE-RQ PDU
func (a *AAssociateRQ) Parse(data []byte) error {
	// check PDU type
	if data[0] != 0x01 {
		return fmt.Errorf("Not an A-ASSOCIATE-RQ PDU")
	}

	// check PDU length
	pduLength := binary.BigEndian.Uint32(data[2:6])
	if uint32(len(data)) != pduLength {
		return fmt.Errorf("Invalid PDU length")
	}

	fmt.Println("PDU length:", pduLength)

	// parse PDU
	a.ProtocolVersion = binary.BigEndian.Uint16(data[6:8])
	a.CalledAETitle = string(bytes.TrimRight(data[8:18], "\x00"))
	a.CallingAETitle = string(bytes.TrimRight(data[18:28], "\x00"))
	a.Reserved = binary.BigEndian.Uint32(data[28:32])
	a.ApplicationContextName = string(bytes.TrimRight(data[32:52], "\x00"))

	// parse presentation contexts
	presentationContexts := []PresentationContext{}
	offset := 52
	for offset < len(data) {
		presentationContext := PresentationContext{}
		presentationContext.PresentationContextID = data[offset]
		presentationContext.Reserved = data[offset+1]
		presentationContext.AbstractSyntax = string(bytes.TrimRight(data[offset+2:offset+22], "\x00"))
		offset += 22

		// parse transfer syntaxes
		transferSyntaxes := []string{}
		transferSyntaxesCount := int(data[offset])
		offset++
		for i := 0; i < transferSyntaxesCount; i++ {
			transferSyntaxes = append(transferSyntaxes, string(bytes.TrimRight(data[offset:offset+20], "\x00")))
			offset += 20
		}
		presentationContext.TransferSyntaxes = transferSyntaxes
		presentationContexts = append(presentationContexts, presentationContext)
	}
	a.PresentationContexts = presentationContexts

	// parse user information
	userInformation := UserInformation{}
	userInformation.Parse(data[offset:])
	a.UserInformation = userInformation

	// print A-ASSOCIATE-RQ PDU
	a.Print()

	return nil
}

// parse User Information
func (u *UserInformation) Parse(data []byte) error {
	// parse user data items
	userDataItems := []UserDataItem{}
	offset := 0
	for offset < len(data) {
		userDataItem := UserDataItem{}
		userDataItem.Type = data[offset]
		itemLength := binary.BigEndian.Uint16(data[offset+1 : offset+3])
		userDataItem.Data = data[offset+3 : offset+3+int(itemLength)]
		userDataItems = append(userDataItems, userDataItem)
		offset += 3 + int(itemLength)
	}
	u.UserDataItems = userDataItems

	return nil
}

// print A-ASSOCIATE-RQ PDU
func (a *AAssociateRQ) Print() {
	fmt.Printf("ProtocolVersion: %d	CalledAETitle: %s	CallingAETitle: %s	Reserved: %d	ApplicationContextName: %s	PresentationContexts: %d	UserInformation: %d	\n", a.ProtocolVersion, a.CalledAETitle, a.CallingAETitle, a.Reserved, a.ApplicationContextName, len(a.PresentationContexts), len(a.UserInformation.UserDataItems))
	for _, presentationContext := range a.PresentationContexts {
		fmt.Printf("	PresentationContextID: %d	Reserved: %d	AbstractSyntax: %s	TransferSyntaxes: %d	\n", presentationContext.PresentationContextID, presentationContext.Reserved, presentationContext.AbstractSyntax, len(presentationContext.TransferSyntaxes))
		for _, transferSyntax := range presentationContext.TransferSyntaxes {
			fmt.Printf("		%s	\n", transferSyntax)
		}
	}
	for _, userDataItem := range a.UserInformation.UserDataItems {
		fmt.Printf("	UserDataItem: %d	\n", userDataItem.Type)
	}
}
