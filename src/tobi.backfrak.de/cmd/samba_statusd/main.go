package main

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"tobi.backfrak.de/internal/commonbl"
)

// Authors - Information about the authors of the program. You might want to add your name here when contributing to this software
const Authors = "tobi@backfrak.de"

// The version of this program, will be set at compile time by the gradle build script
var version = "undefined"

// Type for functions that can create a response string
type response func(commonbl.PipeHandler, string) error

func main() {
	handleComandlineOptions()
	pipeHandler := *commonbl.NewPipeHandler(params.Test)
	if params.Verbose {
		args := ""
		for _, arg := range os.Args {
			args = fmt.Sprintf("%s %s", args, arg)
		}
		fmt.Fprintln(os.Stdout, fmt.Sprintf("Call: %s", args))
		if !params.PrintVersion {
			printVersion()
		}
		fmt.Fprintln(os.Stdout, fmt.Sprintf("Use named pipe: %s", pipeHandler.GetPipeFilePath()))
	}

	if params.PrintVersion {
		printVersion()
		os.Exit(0)
	}

	if params.Help {
		flag.Usage()
		os.Exit(0)
	}

	// Ensure we exit clean on term and kill signals
	go waitforKillSignalAndExit()
	go waitforTermSignalAndExit()

	// Wait for pipe input and process it in an infinite loop
	for {
		received, errRecv := pipeHandler.WaitForPipeInputString()
		if errRecv != nil {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("Error while receive data from the pipe: %s", errRecv))
			os.Exit(-1)
		}

		if strings.HasPrefix(received, commonbl.PROCESS_REQUEST) {
			handleRequest(pipeHandler, received, commonbl.PROCESS_REQUEST, processResponse, testProcessResponse)
		} else if strings.HasPrefix(received, commonbl.SERVICE_REQUEST) {
			handleRequest(pipeHandler, received, commonbl.SERVICE_REQUEST, serviceResponse, testServiceResponse)
		} else if strings.HasPrefix(received, commonbl.LOCK_REQUEST) {
			handleRequest(pipeHandler, received, commonbl.LOCK_REQUEST, lockResponse, testLockResponse)
		}

		time.Sleep(time.Millisecond)
	}

}

func handleRequest(handler commonbl.PipeHandler, request string, requestType string, productiveFunc response, testFunc response) {
	id := getIdFromRequest(request)
	if params.Verbose {
		fmt.Fprintln(os.Stdout, fmt.Sprintf("Handle \"%s\" with id %s", requestType, id))
	}

	var writeErr error
	if !params.Test {
		writeErr = productiveFunc(handler, id)
	} else {
		writeErr = testFunc(handler, id)
	}
	if writeErr != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Error while write \"%s\" response to pipe: %s", requestType, writeErr))
		os.Exit(-1)
	}
}

func lockResponse(handler commonbl.PipeHandler, id string) error {
	fmt.Fprintln(os.Stderr, fmt.Sprintf("Error: Productive code not implemented yet"))
	os.Exit(-2)

	return nil
}

func serviceResponse(handler commonbl.PipeHandler, id string) error {
	fmt.Fprintln(os.Stderr, fmt.Sprintf("Error: Productive code not implemented yet"))
	os.Exit(-2)

	return nil
}

func processResponse(handler commonbl.PipeHandler, id string) error {
	fmt.Fprintln(os.Stderr, fmt.Sprintf("Error: Productive code not implemented yet"))
	os.Exit(-2)

	return nil
}

func testProcessResponse(handler commonbl.PipeHandler, id string) error {
	return handler.WritePipeString(fmt.Sprintf("%s Test Response for request %s", commonbl.PROCESS_REQUEST, id))
}

func testServiceResponse(handler commonbl.PipeHandler, id string) error {
	return handler.WritePipeString(fmt.Sprintf("%s Test Response for request %s", commonbl.SERVICE_REQUEST, id))
}

func testLockResponse(handler commonbl.PipeHandler, id string) error {
	return handler.WritePipeString(fmt.Sprintf("%s Test Response for request %s", commonbl.LOCK_REQUEST, id))
}

func getIdFromRequest(request string) string {
	splitted := strings.Split(request, ":")

	if len(splitted) != 2 {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Error: Got invalid request: \"%s\"", request))
		os.Exit(-1)
	}

	return strings.TrimSpace(splitted[1])
}

func waitforKillSignalAndExit() {
	killSignal := make(chan os.Signal, syscall.SIGKILL)
	signal.Notify(killSignal, os.Interrupt)
	<-killSignal

	if params.Verbose {
		fmt.Fprintln(os.Stdout, fmt.Sprintf("End: %s due to kill signal", os.Args[0]))
	}
	os.Exit(0)
}

func waitforTermSignalAndExit() {
	termSignal := make(chan os.Signal, syscall.SIGTERM)
	signal.Notify(termSignal, os.Interrupt)
	<-termSignal
	if params.Verbose {
		fmt.Fprintln(os.Stdout, fmt.Sprintf("End: %s due to terminate signal", os.Args[0]))
	}
	os.Exit(0)
}

// Prints the version string
func printVersion() {
	fmt.Fprintln(os.Stdout, getVersion())
}

// Get the version string
func getVersion() string {
	return fmt.Sprintf("Version: %s", version)
}
