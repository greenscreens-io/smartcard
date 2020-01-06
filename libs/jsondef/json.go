/**
 * Copyright (C) 2015, 2016  Green Screens Ltd.
 */

 // Package jsondef  contains http requests definitions
package jsondef

import (
	"log"
	"time"
	"strconv"
	"encoding/base64"
)

// SmartCardData Request data structure
type SmartCardData struct {
	Type  int `json:"type"`		//
	Cls  int `json:"cls"`		//
	Ins  int `json:"ins"`		//
	P1  int `json:"p1"`			//
	P2  int `json:"p2"`			//
	Data string `json:"data"`	//
	Le  int `json:"le"`			//
}

// PrintLog Log printing request
func (d SmartCardData) PrintLog() {
	log.Printf("Request: Class:%d, Instruction:%d, Param1:%d, Param2:%d, data:%s", d.Cls, d.Ins, d.P1, d.P2, d.Data)
}

//AsText Decode data from base64
 func (d SmartCardData) AsText() string {
	return string(d.AsBinary())
}

// AsBinary Decode data from base64
func (d SmartCardData) AsBinary() []byte {

	if (d.Data == "") {
		return []byte{}
	}

	sDec, _ := base64.StdEncoding.DecodeString(d.Data)
	return sDec
}

// Contains Search an element in given array
func Contains(list []string, val string) bool {

	for _, el := range list {

		if el == val {
			return true
		}

	}

	return false
}

// getRandomString Generate random file name
func getRandomString() (string) {
	now := time.Now()
	sec := now.Unix()
	return strconv.FormatInt(sec, 10)
}
