package storescp

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
	applicationContext  ApplicationContext
	presentationContext []PresentationContext
	userInfo            UserInfo
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
	transferSyntax        []TransferSyntax
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
	maxLength uint32
}
