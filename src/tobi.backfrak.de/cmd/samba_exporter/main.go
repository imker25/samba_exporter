package main

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"tobi.backfrak.de/internal/commonbl"
)

// Authors - Information about the authors of the program. You might want to add your name here when contributing to this software
const Authors = "tobi@backfrak.de"

// The version of this program, will be set at compile time by the gradle build script
var version = "undefined"
var requestCount = 0

func main() {
	handleComandlineOptions()
	pipeHander := *commonbl.NewPipeHandler(params.Test)
	if params.Verbose {
		args := ""
		for _, arg := range os.Args {
			args = fmt.Sprintf("%s %s", args, arg)
		}
		fmt.Fprintln(os.Stdout, fmt.Sprintf("Call: %s", args))
		if !params.PrintVersion {
			printVersion()
		}
		fmt.Fprintln(os.Stdout, fmt.Sprintf("Use named pipe: %s", pipeHander.GetPipeFilePath()))
	}

	if params.PrintVersion {
		printVersion()
		os.Exit(0)
	}

	if params.Help {
		flag.Usage()
		os.Exit(0)
	}

	fmt.Fprintln(os.Stdout, getStatusData(pipeHander))
	fmt.Fprintln(os.Stdout, getConnectionsData(pipeHander))

	os.Exit(0)
}

func getStatusData(handler commonbl.PipeHandler) string {
	requestCount++
	requestString := fmt.Sprintf("%s %d", commonbl.STATUS_REQUEST, requestCount)
	errWrite := handler.WritePipeString(requestString)
	if errWrite != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Error while write STATUS_REQUEST to the pipe: %s", errWrite))
		os.Exit(-1)
	}

	var errRead error
	response := requestString
	for response == requestString && errRead == nil {
		time.Sleep(time.Millisecond)
		response, errRead = handler.WaitForPipeInputString()
	}
	if errRead != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Error while read STATUS_REQUEST response from the pipe: %s", errRead))
		os.Exit(-1)
	}

	if !strings.Contains(response, commonbl.STATUS_REQUEST) &&
		!strings.Contains(response, fmt.Sprintf("esponse for request %d", requestCount)) {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Error: Got unexpected response: \"%s\"", response))
		os.Exit(-1)
	}

	return response
}

func getConnectionsData(handler commonbl.PipeHandler) string {
	requestCount++
	requestString := fmt.Sprintf("%s %d", commonbl.CONNECTIONS_REQUEST, requestCount)
	errWrite := handler.WritePipeString(requestString)
	if errWrite != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Error while write CONNECTIONS_REQUEST to the pipe: %s", errWrite))
		os.Exit(-1)
	}

	var errRead error
	response := requestString
	for response == requestString && errRead == nil {
		time.Sleep(time.Millisecond)
		response, errRead = handler.WaitForPipeInputString()
	}
	if errRead != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Error while read CONNECTIONS_REQUEST response from the pipe: %s", errRead))
		os.Exit(-1)
	}

	if !strings.Contains(response, commonbl.CONNECTIONS_REQUEST) &&
		!strings.Contains(response, fmt.Sprintf("esponse for request %d", requestCount)) {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Error: Got unexpected response: \"%s\"", response))
		os.Exit(-1)
	}

	return response
}

// Prints the version string
func printVersion() {
	fmt.Fprintln(os.Stdout, getVersion())
}

// Get the version string
func getVersion() string {
	return fmt.Sprintf("Version: %s", version)
}
