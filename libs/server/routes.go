/**
 * Copyright (C) 2015, 2016  Green Screens Ltd.
 */

// Package server contains Http server and
// routes for processing requests
package server

import (
	"fmt"
	"log"
	"strconv"
	"net/http"
	"encoding/json"
	"encoding/base64"
	"greenscreens-io/smartcard/libs/jsondef"
	"greenscreens-io/smartcard/libs/api"
)

const greeting = "Welcome to Green Screens Ltd. Smart Card!"
const mimeJSON = "application/json"
const mimeText = "text/plain"

// routeHello  - Print informational message
func routeHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", mimeText)
	fmt.Fprintf(w, greeting)
}

// routeSmartCardList - Return JSON array of localy installed smartcards
func routeSmartCardList(w http.ResponseWriter, r *http.Request) {

	devices, err := api.SmartCardList()

	if err != nil {
		log.Printf("ESR01: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		json, _ := json.Marshal(devices)
		msg := fmt.Sprintf(`{"success" : true, "data" : %s}`, string(json))
		fmt.Fprintln(w, msg)
	}

}

// routeSmartCardOpen - Open smartcard by given id
func routeSmartCardOpen(w http.ResponseWriter, r *http.Request) {

	var raw [] byte
	id, err := getID(w, r)

	if err == nil {
		raw, err = api.SmartCardConnect(id)
		sendResponseRaw(w, r, raw, err)
	}

}

// routeSmartCardClose - Close active smartcard if exist
func routeSmartCardClose(w http.ResponseWriter, r *http.Request) {
	api.SmartCardDisconnect()
	fmt.Fprintln(w, `{ "success" : true}`)
}

// routeSmartCardDOB - Get Discovery Object from smartcard
func routeSmartCardDOB(w http.ResponseWriter, r *http.Request) {
	resp, err := api.SmartCardDiscoveryObject()
	sendResponse(w, r, resp, err)
}

// routeSmartCardBIO - Get Biometric data
func routeSmartCardBIO(w http.ResponseWriter, r *http.Request) {
	resp, err := api.SmartCardBIO()
	sendResponse(w, r, resp, err)
}

// routeSmartCardVersion - Get smartcard version
func routeSmartCardVersion(w http.ResponseWriter, r *http.Request) {
	resp, err := api.SmartCardVersion()
	sendResponse(w, r, resp, err)
}

// routeSmartCardPinTrials - Get smartcard PIN remain trials
func routeSmartCardPinTrials(w http.ResponseWriter, r *http.Request) {
	resp, err := api.SmartCardPINTrials()
	sendResponse(w, r, resp, err)
}

// routeSmartCardPin - Enter smartcard PIN
func routeSmartCardPin(w http.ResponseWriter, r *http.Request) {

	pin := r.URL.Query().Get("id")
	sDec, err := base64.StdEncoding.DecodeString(pin)

	if err != nil {
		log.Printf(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := api.SmartCardPIN(sDec)
	sendResponse(w, r, resp, err)
}

// routeSmartCardOID - Get OID from smartcard
func routeSmartCardOID(w http.ResponseWriter, r *http.Request) {

	var resp api.SmartCardResponse
	id, err := getID(w, r)

	if err == nil {
		resp, err = api.SmartCardOID(id)
		sendResponse(w, r, resp, err)
	}
}

/**
* Receive file URL for downlaod and print - POST
*	{
*		"type" : 1,   - command mode 1..4
*		"cla" : 0,   - class
*		"ins" : 0,   - instruction
*		"p1" : 0,   - parameter 1
*		"p2" : 0,   - parameter 2
*		"le" : 0,   - expected return length
*		"data" : "" - base64 encoded bytes
*	}
 */
func routeSmartCardCommand(w http.ResponseWriter, r *http.Request) {

	var resp api.SmartCardResponse
	var data *jsondef.SmartCardData

	data, err := parsePostAsJSON(r)

	if err != nil || data == nil {
		log.Printf(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	raw := data.AsBinary()
	resp, err = api.SmartCardCommand(data.Type, data.Cls, data.Ins, data.P1, data.P2, raw, data.Le)

	sendResponse(w, r, resp, err)
}

// getID - Get HTTP Query ID parameter
func getID(w http.ResponseWriter, r *http.Request) (int, error) {

	var err error
	var id int

	ids := r.URL.Query().Get("id")

	if ids != "" {

		id, err = strconv.Atoi(ids)

		if err != nil {
			log.Printf(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

	}

	return id, err
}

// sendResponseRaw - Send bytes as hex encoded
func sendResponseRaw(w http.ResponseWriter, r *http.Request, resp []byte, err error) {

	if err != nil {
		log.Printf(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	msg := fmt.Sprintf(`{"success" : true, "data" : "%x"}`, resp)
	fmt.Fprintln(w, msg)

}

// sendResponse - Send smartcard object response to web
func sendResponse(w http.ResponseWriter, r *http.Request, resp api.SmartCardResponse, err error) {

	if err != nil {
		log.Printf(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	b, err := json.Marshal(resp)
	if err != nil {
		log.Printf(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	msg := fmt.Sprintf(`{"success" : true, "data" : %s}`, string(b))
	fmt.Fprintln(w, msg)

}

// parsePostAsJSON - Parse json-encoded post
func parsePostAsJSON(r *http.Request) (*jsondef.SmartCardData, error) {

	defer func() {
		if err := recover(); err != nil {
			msg := ""
			switch x := err.(type) {
			case string:
				msg = x
			case error:
				msg = x.Error()
			default:
				msg = "unknown error"
			}
			log.Printf("ESR02: %s", msg)
		}
	}()

	var err error
	var data *jsondef.SmartCardData

	contentType := r.Header.Get("Content-type")

	if mimeJSON == contentType {

		decoder := json.NewDecoder(r.Body)

		data = &jsondef.SmartCardData{}
		err = decoder.Decode(data)

		if err != nil {
			return nil, fmt.Errorf("ESR03: %s", err.Error())
		}

	} else {
		err = fmt.Errorf("ESR04: %s", "Invalid type, JSON expected")
	}

	return data, err
}
