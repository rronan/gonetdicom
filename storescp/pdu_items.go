package storescp

import "fmt"

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
	itemType uint8
	reserved uint8
	length   [2]byte
	subItem  SubItem
}

type SubItem struct {
	itemType  uint8
	reserved  uint8
	length    [2]byte
	maxLength [4]byte
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
	s += fmt.Sprintf("		Sub item:\n")
	s += fmt.Sprintf("			Item type: %x \n", a.variableItems.userInfo.subItem.itemType)
	s += fmt.Sprintf("			Reserved: %x \n", a.variableItems.userInfo.subItem.reserved)
	s += fmt.Sprintf("			Length: %x \n", a.variableItems.userInfo.subItem.length)
	s += fmt.Sprintf("			Max length: %x \n", a.variableItems.userInfo.subItem.maxLength)

	return s
}
