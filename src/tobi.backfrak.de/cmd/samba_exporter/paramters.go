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
	"tobi.backfrak.de/internal/smbexporterbl/statisticsGenerator"
)

// The paramters for this executable
type parmeters struct {
	commonbl.Parmeters
	statisticsGenerator.StatisticsGeneratorSettings
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
	flag.BoolVar(&params.TestPipeMode, "test-pipe", false, "Requests status from samba_statusd and exits. May be combined with -test-mode.")
	flag.StringVar(&params.ListenAddress, "web.listen-address", ":9922", "Address to listen on for web interface and telemetry.")
	flag.StringVar(&params.MetricsPath, "web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	flag.IntVar(&params.RequestTimeOut, "request-timeout", 5, "The timeout for a request to samba_statusd in seconds")
	flag.BoolVar(&params.DoNotExportEncryption, "not-expose-encryption-data", false, "Set to 'true', no details about the used encryption or signing will be exported")
	flag.BoolVar(&params.DoNotExportClient, "not-expose-client-data", false, "Set to 'true', no details about the connected clients will be exported")
	flag.BoolVar(&params.DoNotExportUser, "not-expose-user-data", false, "Set to 'true', no details about the connected users will be exported")
	flag.BoolVar(&params.DoNotExportPid, "not-expose-pid-data", false, "Set to 'true', no process IDs will be exported")
	flag.BoolVar(&params.DoNotExportShareDetails, "not-expose-share-details", false, "Set to 'true', no details about the shares will be exported")
	flag.StringVar(&params.LogFilePath, "log-file-path", " ",
		"Give the full file path for a log file. When parameter is not set (as by default), logs will be written to stdout and stderr")

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
	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, "This program is used to run as a service. To change the service behavior edit '/etc/default/samba_exporter' according to your needs.")
}
