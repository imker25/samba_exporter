package main

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"flag"
	"fmt"
	"os"

	"tobi.backfrak.de/internal/commonbl"
)

// The paramters for this executable
type parmeters struct {
	commonbl.Parmeters
}

var params parmeters

// Setup commandline parameters  and parse them
func handleComandlineOptions() {

	// Setup the usabel parametes
	flag.BoolVar(&params.PrintVersion, "print-version", false, "With this flag the program will only print it's version and exit")
	flag.BoolVar(&params.Verbose, "verbose", false, "With this flag the program will print verbose output")
	flag.BoolVar(&params.Help, "help", false, "Print this help message")

	// Overwrite the std Usage function with some custom stuff
	flag.Usage = customHelpMessage

	// Read the given flags
	flag.Parse()
}

// customHelpMessage - Print he customized help message
func customHelpMessage() {
	fmt.Fprintln(os.Stdout, fmt.Sprintf("%s: prometheus exporter for the samba file server. Collects data using the samba_statusd service.", os.Args[0]))
	fmt.Fprintln(os.Stdout, fmt.Sprintf("Program %s", getVersion()))
	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, fmt.Sprintf("Usage: %s [options]", os.Args[0]))
	fmt.Fprintln(os.Stdout, "Options:")
	flag.PrintDefaults()
}
