/**
 * Copyright (C) 2015, 2016  Green Screens Ltd.
 */

 // Package api contains functions to cummunicate with
 // smartcards as conencting, issuing commands etc.
package api

import (
	"fmt"
	"log"
	"bytes"
	"github.com/sf1/go-card/smartcard"
)

//**************************************************************************
//   G E N E R I C  C O M M A N D S
//**************************************************************************

// readData - read all data and marege into single byte array
func readData(card *smartcard.Card, data smartcard.ResponseAPDU) ([] byte) {

	var res = data
	var bb = bytes.NewBuffer(data.Data());

	for (res.SW1() == 0x61) {
		res, _ = commandRESPLe(card, res.SW2())
		bb.Write(res.Data())
	}

	return bb.Bytes()
}

// commandGETLe - call GET function wth expected length response
func commandGETLe(card *smartcard.Card, data []byte, len byte) (smartcard.ResponseAPDU, error) {
	var apdu = smartcard.Command4(0x00, 0xcb, 0x3f, 0xff, data, len)
	return command(card, apdu)
}

// commandGET - call GET function
func commandGET(card *smartcard.Card, data []byte) (smartcard.ResponseAPDU, error) {
	var apdu = smartcard.Command3(0x00, 0xcb, 0x3f, 0xff, data)
	return command(card, apdu)
}

// commandRESPLe - command to receive remaining data
func commandRESPLe(card *smartcard.Card, len byte) (smartcard.ResponseAPDU, error) {
	var apdu = smartcard.Command2(0x00, 0xc0, 0x00, 0x00, len)
	return command(card, apdu)
}

// command - generic smartcard comamnd call
func command(card *smartcard.Card, command smartcard.CommandAPDU) (smartcard.ResponseAPDU, error) {

	fmt.Printf("Request: %s\n", command)

	response, err := card.TransmitAPDU(command)

	if err != nil {
		log.Printf("SMC01: %s", err.Error());
		log.Fatal(err)
	}

	fmt.Printf("Response: %s\n", response)
	return response, err
}
