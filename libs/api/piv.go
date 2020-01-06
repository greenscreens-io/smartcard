/**
 * Copyright (C) 2015, 2016  Green Screens Ltd.
 */

 // Package api contains functions to cummunicate with
 // smartcards as conencting, issuing commands etc.
package api

import (
	"github.com/sf1/go-card/smartcard"
)

//**************************************************************************
//   PIV FUNCTIONS
//**************************************************************************

// commandSelect - select PIV Applet
func commandSelect(card *smartcard.Card) ([] byte) {
	arg := []byte{0xa0, 0x00, 0x00, 0x03, 0x08, 0x00, 0x00, 0x10, 0x00}
	apdu := smartcard.SelectCommand(arg...)
	res, _ := command(card, apdu)
	return readData(card, res)
}

// commandOID - NIST.SP.800-73-4.pdf (Table 3, pg. 30)
func commandOID(card *smartcard.Card, code byte) (smartcard.ResponseAPDU, error) {
	var arg = []byte{0x5c, 0x03, 0x5f, 0xc1, code}
	return commandGET(card, arg)
}

// commandBIO - Get biometric data
func commandBIO(card *smartcard.Card) (smartcard.ResponseAPDU, error) {
	var arg = []byte{0x7f,0x61}
	return commandGET(card, arg)
}

// commandVersion - yubikey custom
func commandVersion(card *smartcard.Card) (smartcard.ResponseAPDU, error) {
	var apdu = smartcard.Command1(0x00, 0xfd, 0x00, 0x00)
	return command(card, apdu)
}

// commandDOB - get discovery object
func commandDOB(card *smartcard.Card) (smartcard.ResponseAPDU, error) {
	var arg = []byte{0x5c, 0x01, 0x7e}
	return commandGET(card, arg)
}

// commandPinRetry - get number of pin retrys
func commandPinRetry(card *smartcard.Card) (smartcard.ResponseAPDU, error) {
	var apdu = smartcard.Command1(0x00, 0x20, 0x00, 0x80)
	return command(card, apdu)
}

// commandPin - enter card pin
func commandPin(card *smartcard.Card, data []byte) (smartcard.ResponseAPDU, error) {
	//var arg = []byte{0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0xff, 0xff }
	var apdu = smartcard.Command3(0x00, 0x20, 0x00, 0x80, data)
	return command(card, apdu)
}

// smartCardCommand4 - call commad option 4
func smartCardCommand4(cls byte, ins byte, p1 byte, p2 byte, data []byte, le byte) (smartcard.ResponseAPDU, error) {
	var apdu = smartcard.Command4(cls, ins, p1, p2, data, le)
	return command(activeCard, apdu)
}

// smartCardCommand3 - call commad option 3
func smartCardCommand3(cls byte, ins byte, p1 byte, p2 byte, data []byte) (smartcard.ResponseAPDU, error) {
	var apdu = smartcard.Command3(cls, ins, p1, p2, data)
	return command(activeCard, apdu)
}

// smartCardCommand2 - call commad option 2
func smartCardCommand2(cls byte, ins byte, p1 byte, p2 byte, le byte) (smartcard.ResponseAPDU, error) {
	var apdu = smartcard.Command2(cls, ins, p1, p2, le)
	return command(activeCard, apdu)
}

// smartCardCommand1 - call commad option 1
func smartCardCommand1(cls byte, ins byte, p1 byte, p2 byte) (smartcard.ResponseAPDU, error) {
	var apdu = smartcard.Command1(cls, ins, p1, p2)
	return command(activeCard, apdu)
}
