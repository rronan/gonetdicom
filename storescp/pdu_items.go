package storescp

import (
	"encoding/binary"
	"fmt"
)

type Associate struct {
	pduType         uint8
	reserved        uint8
	length          [4]byte
	protocolVersion [2]byte
	reserved2       [2]byte
	calledAETitle   [16]byte
	callingAETitle  [16]byte
	reserved3       [32]byte
	variableItems   VariableItems
}

type VariableItems struct {
	applicationContext      ApplicationContext
	presentationContextList []PresentationContext
	userInfo                UserInfo
}

type ApplicationContext struct {
	itemType               uint8
	reserved               uint8
	length                 [2]byte
	applicationContextName []byte
}

type PresentationContext struct {
	itemType              uint8
	reserved              uint8
	length                [2]byte
	presentationContextID uint8
	reserved2             uint8
	resultReason          uint8
	reserved3             uint8
	abstractSyntax        AbstractSyntax
	transferSyntaxList    []TransferSyntax
}

type AbstractSyntax struct {
	itemType           uint8
	reserved           uint8
	length             [2]byte
	abstractSyntaxName []byte
}

type TransferSyntax struct {
	itemType           uint8
	reserved           uint8
	length             [2]byte
	transferSyntaxName []byte
}

type UserInfo struct {
	itemType            uint8
	reserved            uint8
	length              [2]byte
	userInfoSubItemList []UserInfoSubItem
}

type UserInfoSubItem struct {
	itemType  uint8
	reserved  uint8
	length    [2]byte
	maxLength [4]byte
}

func CreateAssociateAC(associateRQ Associate) (Associate, error) {

	ACSubItem := createUserInfoSubItem(16384)

	ACUserInfo := UserInfo{
		itemType: 0x50,
		reserved: 0x00,
		length:   [2]byte{0x00, 0x08},
		userInfoSubItemList: []UserInfoSubItem{
			ACSubItem,
		},
	}
	ACTransferSyntax := TransferSyntax{
		itemType:           0x40,
		reserved:           0x00,
		transferSyntaxName: associateRQ.variableItems.presentationContextList[0].transferSyntaxList[0].transferSyntaxName,
	}
	binary.BigEndian.PutUint16(ACTransferSyntax.length[:], uint16(len(ACTransferSyntax.transferSyntaxName)))
	ACTransferSyntaxArray := []TransferSyntax{ACTransferSyntax}
	ACAbstractSyntax := AbstractSyntax{
		itemType:           0x30,
		reserved:           0x00,
		abstractSyntaxName: associateRQ.variableItems.presentationContextList[0].abstractSyntax.abstractSyntaxName,
	}
	binary.BigEndian.PutUint16(ACAbstractSyntax.length[:], uint16(len(ACAbstractSyntax.abstractSyntaxName)))
	ACPresentationContext := PresentationContext{
		itemType:              0x21,
		reserved:              0x00,
		presentationContextID: associateRQ.variableItems.presentationContextList[0].presentationContextID,
		reserved2:             0x00,
		resultReason:          0x00,
		reserved3:             0x00,
		transferSyntaxList:    ACTransferSyntaxArray,
	}
	binary.BigEndian.PutUint16(ACPresentationContext.length[:], uint16(len(ACTransferSyntax.transferSyntaxName)+8))
	ACApplicationContext := ApplicationContext{
		itemType:               0x10,
		reserved:               0x00,
		applicationContextName: associateRQ.variableItems.applicationContext.applicationContextName,
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
		calledAETitle:   associateRQ.calledAETitle,
		callingAETitle:  associateRQ.callingAETitle,
		reserved3:       [32]byte{0x00},
		variableItems:   ACVariableItems,
	}

	mint := binary.BigEndian.Uint16(ACVariableItems.applicationContext.length[:])
	mint2 := binary.BigEndian.Uint16(ACVariableItems.presentationContextList[0].length[:])
	mint3 := binary.BigEndian.Uint16(ACVariableItems.userInfo.length[:])

	binary.BigEndian.PutUint32(AAACStruct.length[:], uint32(mint+mint2+mint3+2+2+16+16+32))
	fmt.Println(AAACStruct.ToString())

	return AAACStruct, nil
}

func createUserInfoSubItem(maxLength uint32) UserInfoSubItem {
	// get max length and convert to byte array
	var maxLengthBytes [4]byte
	tmp := make([]byte, 4)
	binary.BigEndian.PutUint32(tmp, maxLength)
	copy(maxLengthBytes[:], tmp)
	fmt.Println("max length bytes: ", maxLengthBytes)
	fmt.Println("byte max length: ", maxLengthBytes[0], maxLengthBytes[1], maxLengthBytes[2], maxLengthBytes[3])
	fmt.Println("byte max length: ", [4]byte{0x00, 0x00, 0x40, 0x00})

	fmt.Println("last", binary.BigEndian.Uint32([]byte{0x00, 0x00, 0x40, 0x00}))

	return UserInfoSubItem{
		itemType:  0x51,
		reserved:  0x00,
		length:    [2]byte{0x00, 0x04},
		maxLength: maxLengthBytes,
	}
}

func (a *Associate) ToString() string {
	var s string
	s += fmt.Sprintf("A-ASSOCIATE-RQ/AC PDU \n")
	s += fmt.Sprintf("PDU type: %x \n", a.pduType)
	s += fmt.Sprintf("Reserved: %x \n", a.reserved)
	s += fmt.Sprintf("Length: %x \n", a.length)
	s += fmt.Sprintf("Protocol version: %x \n", a.protocolVersion)
	s += fmt.Sprintf("Reserved: %x \n", a.reserved2)
	s += fmt.Sprintf("Called AE title: %s \n", a.calledAETitle)
	s += fmt.Sprintf("Calling AE title: %s \n", a.callingAETitle)
	s += fmt.Sprintf("Reserved: %x \n", a.reserved3)
	s += fmt.Sprintf("Variable items:\n")
	s += fmt.Sprintf("	Application context:\n")
	s += fmt.Sprintf("		Item type: %x \n", a.variableItems.applicationContext.itemType)
	s += fmt.Sprintf("		Reserved: %x \n", a.variableItems.applicationContext.reserved)
	s += fmt.Sprintf("		Length: %x \n", a.variableItems.applicationContext.length)
	s += fmt.Sprintf("		Application context name: %s \n", a.variableItems.applicationContext.applicationContextName)
	s += fmt.Sprintf("	Presentation context:\n")
	for i := 0; i < len(a.variableItems.presentationContextList); i++ {
		s += fmt.Sprintf("		Presentation context %d:\n", i)
		s += fmt.Sprintf("			Item type: %x \n", a.variableItems.presentationContextList[i].itemType)
		s += fmt.Sprintf("			Reserved: %x \n", a.variableItems.presentationContextList[i].reserved)
		s += fmt.Sprintf("			Length: %x \n", a.variableItems.presentationContextList[i].length)
		s += fmt.Sprintf("			Presentation context ID: %x \n", a.variableItems.presentationContextList[i].presentationContextID)
		s += fmt.Sprintf("			Reserved: %x \n", a.variableItems.presentationContextList[i].reserved2)
		s += fmt.Sprintf("			Result reason: %x \n", a.variableItems.presentationContextList[i].resultReason)
		s += fmt.Sprintf("			Reserved: %x \n", a.variableItems.presentationContextList[i].reserved3)
		s += fmt.Sprintf("			Abstract syntax:\n")
		s += fmt.Sprintf("				Item type: %x \n", a.variableItems.presentationContextList[i].abstractSyntax.itemType)
		s += fmt.Sprintf("				Reserved: %x \n", a.variableItems.presentationContextList[i].abstractSyntax.reserved)
		s += fmt.Sprintf("				Length: %x \n", a.variableItems.presentationContextList[i].abstractSyntax.length)
		s += fmt.Sprintf("				Abstract syntax name: %s \n", a.variableItems.presentationContextList[i].abstractSyntax.abstractSyntaxName)
		s += fmt.Sprintf("			Transfer syntax:\n")
		for j := 0; j < len(a.variableItems.presentationContextList[i].transferSyntaxList); j++ {
			s += fmt.Sprintf("				Transfer syntax %d:\n", j)
			s += fmt.Sprintf("					Item type: %x \n", a.variableItems.presentationContextList[i].transferSyntaxList[j].itemType)
			s += fmt.Sprintf("					Reserved: %x \n", a.variableItems.presentationContextList[i].transferSyntaxList[j].reserved)
			s += fmt.Sprintf("					Length: %x \n", a.variableItems.presentationContextList[i].transferSyntaxList[j].length)
			s += fmt.Sprintf("					Transfer syntax name: %s \n", a.variableItems.presentationContextList[i].transferSyntaxList[j].transferSyntaxName)
		}
	}
	s += fmt.Sprintf("	User info:\n")
	s += fmt.Sprintf("		Item type: %x \n", a.variableItems.userInfo.itemType)
	s += fmt.Sprintf("		Reserved: %x \n", a.variableItems.userInfo.reserved)
	s += fmt.Sprintf("		Length: %x \n", a.variableItems.userInfo.length)
	s += fmt.Sprintf("		User info sub items:\n")
	for i := 0; i < len(a.variableItems.userInfo.userInfoSubItemList); i++ {
		s += fmt.Sprintf("			User info sub item %d:\n", i)
		s += fmt.Sprintf("				Item type: %x \n", a.variableItems.userInfo.userInfoSubItemList[i].itemType)
		s += fmt.Sprintf("				Reserved: %x \n", a.variableItems.userInfo.userInfoSubItemList[i].reserved)
		s += fmt.Sprintf("				Length: %x \n", a.variableItems.userInfo.userInfoSubItemList[i].length)
		s += fmt.Sprintf("				Max length: %x \n", a.variableItems.userInfo.userInfoSubItemList[i].maxLength)
	}
	return s
}
