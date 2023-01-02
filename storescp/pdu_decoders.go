package storescp

import (
	"encoding/binary"
)

func DecodeAAssociateRQ(b []byte) (a Associate, e error) {
	a.pduType = b[0]
	a.reserved = b[1]
	copy(a.length[:], b[2:6])
	copy(a.protocolVersion[:], b[6:8])
	a.reserved2 = [2]byte{b[8], b[9]}
	copy(a.calledAETitle[:], b[10:26])
	copy(a.callingAETitle[:], b[26:42])
	copy(a.reserved3[:], b[42:74])
	a.variableItems, _ = decodeVariableItems(b[74:])
	return a, nil
}

func EncodeAAssociateAC(a Associate) ([]byte, error) {
	var b []byte
	b = append(b, a.pduType)
	b = append(b, a.reserved)
	b = append(b, a.length[:]...)
	b = append(b, a.protocolVersion[:]...)
	b = append(b, a.reserved2[:]...)
	b = append(b, a.calledAETitle[:]...)
	b = append(b, a.callingAETitle[:]...)
	b = append(b, a.reserved3[:]...)
	b = append(b, a.variableItems.applicationContext.itemType)
	b = append(b, a.variableItems.applicationContext.reserved)
	b = append(b, a.variableItems.applicationContext.length[:]...)
	b = append(b, a.variableItems.applicationContext.applicationContextName[:]...)
	for _, p := range a.variableItems.presentationContextList {
		b = append(b, p.itemType)
		b = append(b, p.reserved)
		b = append(b, p.length[:]...)
		b = append(b, p.presentationContextID)
		b = append(b, p.reserved2)
		b = append(b, p.resultReason)
		b = append(b, p.reserved3)
		for _, t := range p.transferSyntaxList {
			b = append(b, t.itemType)
			b = append(b, t.reserved)
			b = append(b, t.length[:]...)
			b = append(b, t.transferSyntaxName[:]...)
		}
	}
	b = append(b, a.variableItems.userInfo.itemType)
	b = append(b, a.variableItems.userInfo.reserved)
	b = append(b, a.variableItems.userInfo.length[:]...)
	b = append(b, a.variableItems.userInfo.userInfoSubItemList[0].itemType)
	b = append(b, a.variableItems.userInfo.userInfoSubItemList[0].reserved)
	b = append(b, a.variableItems.userInfo.userInfoSubItemList[0].length[:]...)
	b = append(b, a.variableItems.userInfo.userInfoSubItemList[0].maxLength[:]...)
	return b, nil
}

func decodeVariableItems(b []byte) (v VariableItems, lenght int) {
	var appContextLenght int
	var presContextLenght int
	v.applicationContext, appContextLenght = decodeApplicationContext(b)
	v.presentationContextList, presContextLenght = decodePresentationContext(b[appContextLenght:])
	v.userInfo, _ = decodeUserInfo(b[appContextLenght+presContextLenght:])
	return v, len(b)
}

func decodeUserInfo(b []byte) (u UserInfo, lenght int) {
	var subItemLenght int
	u.itemType = b[0]
	u.reserved = b[1]
	copy(u.length[:], b[2:4])
	u.userInfoSubItemList, subItemLenght = decodeSubItem(b[4:])
	return u, 4 + subItemLenght
}

func decodeSubItem(b []byte) (u []UserInfoSubItem, lenght int) {
	u = make([]UserInfoSubItem, 1)
	var totalLenght int
	u[0].itemType = b[0]
	u[0].reserved = b[1]
	copy(u[0].length[:], b[2:4])
	copy(u[0].maxLength[:], b[4:6])
	totalLenght = 6 + int(binary.BigEndian.Uint16(u[0].length[:]))
	// recursive call to decode more sub items
	if totalLenght < len(b) {
		u2, u2Lenght := decodeSubItem(b[totalLenght:])
		u = append(u, u2...)
		totalLenght += u2Lenght
	}
	return u, totalLenght
}

func decodeApplicationContext(b []byte) (a ApplicationContext, lenght int) {
	a.itemType = b[0]
	a.reserved = b[1]
	copy(a.length[:], b[2:4])
	a.applicationContextName = b[4 : 4+binary.BigEndian.Uint16(b[2:4])]
	return a, 4 + int(binary.BigEndian.Uint16(a.length[:]))
}

func decodePresentationContext(b []byte) (p []PresentationContext, lenght int) {
	p = make([]PresentationContext, 1)
	var abstractSyntaxLenght int
	var transferSyntaxLenght int
	var totalLenght int
	p[0].itemType = b[0]
	p[0].reserved = b[1]
	copy(p[0].length[:], b[2:4])
	totalLenght = 4 + int(binary.BigEndian.Uint16(p[0].length[:]))
	p[0].presentationContextID = b[4]
	p[0].reserved2 = b[5]
	p[0].resultReason = b[6]
	p[0].reserved3 = b[7]
	p[0].abstractSyntax, abstractSyntaxLenght = decodeAbstractSyntax(b[8:])
	p[0].transferSyntaxList, transferSyntaxLenght = decodeTransferSyntax(b[8+abstractSyntaxLenght:])

	// recursive call to decode more than one presentation context
	if len(b) > 8+abstractSyntaxLenght+transferSyntaxLenght {
		if b[8+abstractSyntaxLenght+transferSyntaxLenght] == 0x20 {
			p2, p2Lenght := decodePresentationContext(b[8+abstractSyntaxLenght+transferSyntaxLenght:])
			p = append(p, p2...)
			totalLenght += p2Lenght
		}
	}
	return p, totalLenght
}

func decodeAbstractSyntax(b []byte) (a AbstractSyntax, lenght int) {
	a.itemType = b[0]
	a.reserved = b[1]
	copy(a.length[:], b[2:4])
	a.abstractSyntaxName = b[4 : 4+binary.BigEndian.Uint16(a.length[:])]
	return a, 4 + int(binary.BigEndian.Uint16(a.length[:]))
}

func decodeTransferSyntax(b []byte) (t []TransferSyntax, lenght int) {
	t = make([]TransferSyntax, 1)
	var totalLenght int
	t[0].itemType = b[0]
	t[0].reserved = b[1]
	copy(t[0].length[:], b[2:4])
	totalLenght = 4 + int(binary.BigEndian.Uint16(t[0].length[:]))
	t[0].transferSyntaxName = b[4 : 4+binary.BigEndian.Uint16(t[0].length[:])]

	// recursive call to decode more than one transfer syntax
	if len(b) > 4+int(binary.BigEndian.Uint16(t[0].length[:])) {
		if b[4+int(binary.BigEndian.Uint16(t[0].length[:]))] == 0x40 {
			t2, t2Lenght := decodeTransferSyntax(b[4+int(binary.BigEndian.Uint16(t[0].length[:])):])
			t = append(t, t2...)
			totalLenght += t2Lenght
		}
	}
	return t, totalLenght
}
