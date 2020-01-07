/**
 * Copyright (C) 2015, 2016  Green Screens Ltd.
 */

// Package api contains functions to cummunicate with
// smartcards as connecting, issuing commands etc.
package api

import (
	"fmt"
	"encoding/json"
)

// SmartCardTlv decoded tags
type SmartCardTlv struct {
	Struct bool
	Cls string
	Tag string
	Data string
	Child []SmartCardTlv
}

// SmartCardResponse contains data and sw codes
type SmartCardResponse struct {
	Sw1  string
	Sw2  string
	Data string
	Tlv []SmartCardTlv
}

// Device is list of availabe devices to connect
type Device struct {
	ID int
	Name string
}

func (tv SmartCardTlv) String() string {
	b, err := json.Marshal(tv)
	if err != nil {
		return err.Error()
	}

	return string(b)
}


func decodeBERTLV(data []byte) (result []SmartCardTlv) {

	tlvs, err := Decode(data)
	if (err != nil) {
		return result
	}

	for _, tlv := range tlvs {

		stlv := mapBERTLV(tlv);

		deep := tlv.T == 0x53 || tlv.T == 0x30

		if deep || tlv.IsConstructed() {
			list := decodeBERTLV(tlv.V)
			stlv.Child = append(stlv.Child, list...)
		}

		result = append(result, stlv)
	}


	return result
}

func mapBERTLV(tlv TagValue) (tag SmartCardTlv) {
	tag.Struct = tlv.IsConstructed()
	tag.Cls = fmt.Sprintf("%d", tlv.Class())
	tag.Tag = fmt.Sprintf("%x", tlv.T)
	tag.Data = fmt.Sprintf("%x", tlv.V)
	return tag
}
