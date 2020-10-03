/**
 * Copyright (C) 2015, 2016  Green Screens Ltd.
 */

// Package main
// Desktop smartcard linker between browser and locally instaleld smartcard
// Browser JavaScript might call service to communicate directly with smartcard
package main

import (
	"github.com/greenscreens-io/smartcard/applib"
)

// Main program entry point
func main() {
	applib.Startup()
}
