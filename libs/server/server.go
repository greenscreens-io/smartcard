/**
 * Copyright (C) 2015, 2016  Green Screens Ltd.
 */

 // Package server Generic simple http server
package server

import (
	"log"
	"net/http"
	"golang.org/x/time/rate"
)

// limit to one request / second, no bursts
// to prevent DoS
var limiter = rate.NewLimiter(1, 1)

// Initialize URL routes for handling printer
func initRoutes(security bool) {

	isLogging := false
	http.Handle("/", http.FileServer(http.Dir("./assets")))
	//http.HandleFunc("/", Chain(routeHello, Cors("GET"), Logging(isLogging)))
	http.HandleFunc("/list", Chain(routeSmartCardList, Cors("GET"), Logging(isLogging)))
	http.HandleFunc("/connect", Chain(routeSmartCardOpen, Cors("GET"), Logging(isLogging)))
	http.HandleFunc("/disconnect", Chain(routeSmartCardClose, Cors("GET"), Logging(isLogging)))
	http.HandleFunc("/request", Chain(routeSmartCardCommand, Cors("PUT"), Logging(isLogging)))

	http.HandleFunc("/valid", Chain(routeSmartCardPinTrials, Cors("GET"), Logging(isLogging)))
	http.HandleFunc("/pin", Chain(routeSmartCardPin, Cors("GET"), Logging(isLogging)))

	http.HandleFunc("/bio", Chain(routeSmartCardBIO, Cors("GET"), Logging(isLogging)))
	http.HandleFunc("/dob", Chain(routeSmartCardDOB, Cors("GET"), Logging(isLogging)))
	http.HandleFunc("/version", Chain(routeSmartCardVersion, Cors("GET"), Logging(isLogging)))
	http.HandleFunc("/oid", Chain(routeSmartCardOID, Cors("GET"), Logging(isLogging)))
}

// StartServer Start Http server on given port
// If bindAll is true, will bind to all interfaces
// Otherwise, will bind to localhost only
func StartServer(address string, security bool) {
	initRoutes(security)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Printf("ES001 : %s", err.Error())
		log.Fatal("Error while starting : ", err)
	}
}
