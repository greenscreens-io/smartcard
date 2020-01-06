/**
 * Copyright (C) 2015, 2016  Green Screens Ltd.
 */

package applib

import (
	"io"
	"os"
	"fmt"
	"log"
	"greenscreens-io/smartcard/libs/server"
)

var arg *ProgramArgs

// Startup - Main program entry point
func Startup() {

	arg = getArguments()

	initLogger();
	printInfo()

	addr, ip := getAddress(arg.Port, arg.Bind)

	printAddress(addr, ip, arg.Port, arg.Bind)
	server.StartServer(addr, false)
}

// initLogger Initialize file and console logging
func initLogger() {

	f, err := os.OpenFile("application.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("ELF01: %s", err.Error());
		log.Fatal(err)
	}

	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)
}

// printInfo Print Product Information
func printInfo() {

	fmt.Println("**********************************************")
	fmt.Println("* Copyright: Green Screens Ltd. 2016 - 2019  *")
	fmt.Println("* Contact: info@greenscreens.io              *")
	fmt.Println("*                                            *")
	fmt.Println("* Browser To Smart Card                      *")
	fmt.Println("* Version : 1.0.0.                           *")
	fmt.Println("**********************************************")

	getHostname()
}


// printAddress Print Service Information
func printAddress(addr string, ip string, port int, bind bool) {

	fmt.Println("Bind interfaces : ", bind)
	fmt.Println("Service Port    : ", port)
	fmt.Println("Service IP      : ", ip)
	fmt.Println("Listening at    : ", addr)
	fmt.Println("")

}

// getHostname Print and return computer hostname
func getHostname() (string, error) {

	hostname, err := os.Hostname()

	if err != nil {
		log.Printf("ELF02: %s", err.Error());
		log.Fatal(err)
	}

	fmt.Println("Hostname        : ", hostname)
	return hostname, nil
}

// getAddress Generate Server Listening string
func getAddress(port int, bindAll bool) (string, string) {

	ip := ""

	if bindAll {
		ip = "0.0.0.0"
	} else {
		ip = "127.0.0.1"
	}

	addr := fmt.Sprintf("%s:%d", ip, port)
	return addr, ip
}
