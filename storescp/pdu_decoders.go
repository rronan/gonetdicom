package storescp

import (
	"encoding/binary"
	"errors"
)

// type Decoder interface {
// 	Decode() (interface{}, error)
// }

func Decode(b []byte) (interface{}, error) {
	switch b[0] {
	case 0x01:
		return decodeAAssociateRQ(b)
	case 0x02:
		return decodeAAssociateRQ(b)
	// case 0x03:
	// 	return decodeAssociateRJ(b)
	// case 0x04:
	// 	return decodePDataTF(b)
	// case 0x05:
	// 	return decodeReleaseRQ(b)
	// case 0x06:
	// 	return decodeReleaseRP(b)
	// case 0x07:
	// 	return decodeAbort(b)
	// case 0x10:
	// 	return decodeApplicationContext(b)
	// case 0x20:
	// 	return decodePresentationContext(b)
	// case 0x30:
	// 	return decodeAbstractSyntax(b)
	// case 0x40:
	// 	return decodeTransferSyntax(b)
	// case 0x50:
	// 	return decodeUserInfo(b)
	// case 0x21:
	// 	return decodePresentationContextItem(b)
	default:
		return nil, errors.New("unknown pdu type")
	}
}

func decodeAAssociateRQ(b []byte) (Associate, error) {
	var a Associate
	a.pduType = b[0]
	a.reserved = b[1]
	copy(b[2:6], a.length[:])
	copy(b[6:8], a.protocolVersion[:])
	a.reserved2 = [2]byte{b[8], b[9]}
	a.calledAETitle = getAETitle(b[10:26])
	a.callingAETitle = getAETitle(b[26:42])
	a.reserved3 = getBigReserved(b[42:74])
	a.variableItems, _ = decodeVariableItems(b[74:])
	return a, nil
}

func (a Associate) Decode(b []byte) (Associate, error) {
	a.pduType = b[0]
	a.reserved = b[1]
	copy(b[2:6], a.length[:])
	copy(b[6:8], a.protocolVersion[:])
	a.reserved2 = [2]byte{b[8], b[9]}
	a.calledAETitle = getAETitle(b[10:26])
	a.callingAETitle = getAETitle(b[26:42])
	a.reserved3 = getBigReserved(b[42:74])
	// a.variableItems, _ = decodeVariableItems(b[74:])
	return a, nil
}

func decodeVariableItems(b []byte) (VariableItems, int) {
	var v VariableItems
	var appContextLenght int
	// var presContextLenght [2]byte
	v.applicationContext, appContextLenght = decodeApplicationContext(b)
	v.presentationContext, _ = decodePresentationContext(b[len(b)-appContextLenght:])
	return v, len(b)
}

func decodeApplicationContext(b []byte) (ApplicationContext, int) {
	var a ApplicationContext
	a.itemType = b[0]
	a.reserved = b[1]
	copy(b[2:4], a.length[:])
	a.applicationContextName = b[4 : 4+binary.BigEndian.Uint16(b[2:4])]
	return a, 4 + int(binary.BigEndian.Uint16(a.length[:]))
}

func decodePresentationContext(b []byte) ([]PresentationContext, int) {
	var p []PresentationContext
	p[0].itemType = b[0]
	p[0].reserved = b[1]
	copy(b[2:4], p[0].length[:])
	p[0].presentationContextID = b[4]
	p[0].reserved2 = b[5]
	p[0].reserved3 = b[6]
	p[0].reserved4 = b[7]
	// p[0].abstractSyntax, _ = decodeAbstractSyntax(b[8:])
	return p, 0
}

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

type PDU interface {
	Associate | VariableItems | ApplicationContext | PresentationContext | AbstractSyntax | TransferSyntax | UserInfo | UserData
}

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
	reserved3             uint8
	reserved4             uint8
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
	userData []UserData
}

type UserData struct {
	itemType  uint8
	reserved  uint8
	length    [2]byte
	maxLength uint32
}
