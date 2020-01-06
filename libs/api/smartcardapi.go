/**
 * Copyright (C) 2015, 2016  Green Screens Ltd.
 */

// Package api contains functions to cummunicate with
// smartcards as connecting, issuing commands etc.
package api

import (
	"errors"
	"fmt"
	"github.com/sf1/go-card/smartcard"
)

const noCardDetected = "No smartcards detected"
const noCardInit = "card not initialized"

var cardContext *smartcard.Context
var activeCard *smartcard.Card

// SmartCardResponse contains data and sw codes
type SmartCardResponse struct {
	Sw1  string
	Sw2  string
	Data []byte
}

// Device is list of availabe devices to connect
type Device struct {
	ID int
	Name string
}

// getContext - retrieve smartcard API context
func getContext() *smartcard.Context {

	if cardContext == nil {

		ctx, err := smartcard.EstablishContext()

		if err != nil {
			panic(err)
		}

		cardContext = ctx

	}

	return cardContext
}

// freeContext - release smartcard API context
func freeContext() {
	if cardContext != nil {
		cardContext.Release()
		cardContext = nil
	}
}

// SmartCardList list available readers
func SmartCardList() ([]Device, error) {

	var ctx = getContext()
	var list []Device

	readers, err := ctx.ListReadersWithCard()

	if len(readers) == 0 {
		return nil, errors.New(noCardDetected)
	}

	for i, reader := range readers {
		var dev Device
		dev.ID = i
		dev.Name = reader.Name()
		list = append(list, dev)
	}

	return list, err
}

// SmartCardConnect connect to reader by index opening card
func SmartCardConnect(id int) ([]byte, error) {

	var ctx = getContext()
	var rdrs, err = ctx.ListReadersWithCard()

	if err != nil {
		return nil, err
	}

	if len(rdrs) == 0 {
		return nil, errors.New(noCardDetected)
	}

	card, err := rdrs[id].Connect()

	if err != nil {
		return nil, err
	}
	activeCard = card

	var raw = commandSelect(activeCard)
	return raw, nil
}

// SmartCardDisconnect close active card
func SmartCardDisconnect() error {

	if activeCard == nil {
		activeCard.Disconnect()
		activeCard = nil
	}

	return nil
}

// SmartCardVersion return data based on OID code
func SmartCardVersion() (SmartCardResponse, error) {

	var err error
	var resp smartcard.ResponseAPDU
	var result SmartCardResponse

	if activeCard == nil {
		return result, errors.New(noCardInit)
	}

	resp, err = commandVersion(activeCard)

	if err == nil {
		result.Sw1 = fmt.Sprintf("%x", resp.SW1())
		result.Sw2 = fmt.Sprintf("%x", resp.SW2())
		result.Data = readData(activeCard, resp)
	}

	return result, err
}

// SmartCardBIO return Biometric data
func SmartCardBIO() (SmartCardResponse, error) {

	var err error
	var resp smartcard.ResponseAPDU
	var result SmartCardResponse

	if activeCard == nil {
		return result, errors.New(noCardInit)
	}

	resp, err = commandBIO(activeCard)

	if err == nil {
		result.Sw1 = fmt.Sprintf("%x", resp.SW1())
		result.Sw2 = fmt.Sprintf("%x", resp.SW2())
		result.Data = readData(activeCard, resp)
	}

	return result, err

}

// SmartCardDiscoveryObject return data based on OID code
func SmartCardDiscoveryObject() (SmartCardResponse, error) {

	var err error
	var resp smartcard.ResponseAPDU
	var result SmartCardResponse

	if activeCard == nil {
		return result, errors.New(noCardInit)
	}

	resp, err = commandDOB(activeCard)

	if err == nil {
		result.Sw1 = fmt.Sprintf("%x", resp.SW1())
		result.Sw2 = fmt.Sprintf("%x", resp.SW2())
		result.Data = readData(activeCard, resp)
	}

	return result, err
}

// SmartCardPINTrials return number of remain pin entries
func SmartCardPINTrials() (SmartCardResponse, error) {

	var err error
	var resp smartcard.ResponseAPDU
	var result SmartCardResponse

	if activeCard == nil {
		return result, errors.New(noCardInit)
	}

	resp, err = commandPinRetry(activeCard)

	if err == nil {
		result.Sw1 = fmt.Sprintf("%x", resp.SW1())
		result.Sw2 = fmt.Sprintf("%x", resp.SW2())
		result.Data = readData(activeCard, resp)
	}

	return result, err
}

// SmartCardPIN enter card pin, must be 8 bytes, if shorter, fill with 0xff
func SmartCardPIN(data []byte) (SmartCardResponse, error) {

	var err error
	var resp smartcard.ResponseAPDU
	var result SmartCardResponse

	if activeCard == nil {
		return result, errors.New(noCardInit)
	}

	resp, err = commandPin(activeCard, data)

	if err == nil {
		result.Sw1 = fmt.Sprintf("%x", resp.SW1())
		result.Sw2 = fmt.Sprintf("%x", resp.SW2())
		result.Data = readData(activeCard, resp)
	}

	return result, err
}

// SmartCardOID return data based on OID code
func SmartCardOID(data int) (SmartCardResponse, error) {

	var err error
	var resp smartcard.ResponseAPDU
	var result SmartCardResponse

	if activeCard == nil {
		return result, errors.New(noCardInit)
	}

	resp, err = commandOID(activeCard, byte(data))

	if err == nil {
		result.Sw1 = fmt.Sprintf("%x", resp.SW1())
		result.Sw2 = fmt.Sprintf("%x", resp.SW2())
		result.Data = readData(activeCard, resp)
	}

	return result, err
}

// SmartCardCommand send commadn to card
func SmartCardCommand(typ int, cls int, ins int, p1 int, p2 int, data []byte, le int) (SmartCardResponse, error) {

	var err error
	var resp smartcard.ResponseAPDU
	var result SmartCardResponse

	if activeCard == nil {
		return result, errors.New(noCardInit)
	}

	switch typ {
	case 1:
		resp, err = smartCardCommand1(byte(cls), byte(ins), byte(p1), byte(p2))
	case 2:
		resp, err = smartCardCommand2(byte(cls), byte(ins), byte(p1), byte(p2), byte(le))
	case 3:
		resp, err = smartCardCommand3(byte(cls), byte(ins), byte(p1), byte(p2), data)
	case 4:
		resp, err = smartCardCommand4(byte(cls), byte(ins), byte(p1), byte(p2), data, byte(le))
	default:
		err = errors.New("invalid command type")
	}

	if err == nil {
		result.Sw1 = fmt.Sprintf("%x", resp.SW1())
		result.Sw2 = fmt.Sprintf("%x", resp.SW2())
		result.Data = readData(activeCard, resp)
	}

	return result, err
}
