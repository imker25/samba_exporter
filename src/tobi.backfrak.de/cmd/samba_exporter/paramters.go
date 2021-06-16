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
	TestPipeMode   bool
	ListenAddress  string
	MetricsPath    string
	RequestTimeOut int
}

var params parmeters

// Setup commandline parameters  and parse them
func handleComandlineOptions() {

	// Setup the usabel parametes
	flag.BoolVar(&params.PrintVersion, "print-version", false, "With this flag the program will only print it's version and exit")
	flag.BoolVar(&params.Verbose, "verbose", false, "With this flag the program will print verbose output")
	flag.BoolVar(&params.Test, "test-mode", false,
		"Run the program in test mode. In this mode the program will always return the same test data. To work with samba_statusd both programs needs to run in test mode or not.")
	flag.BoolVar(&params.Help, "help", false, "Print this help message")
	flag.BoolVar(&params.TestPipeMode, "test-pipe", false, "Requests status from samba_statusd and exits. May be combinde with -test-mode.")
	flag.StringVar(&params.ListenAddress, "web.listen-address", ":9922", "Address to listen on for web interface and telemetry.")
	flag.StringVar(&params.MetricsPath, "web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	flag.IntVar(&params.RequestTimeOut, "request-timeout", 5, "The timeout for a request to samba_statusd")

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
