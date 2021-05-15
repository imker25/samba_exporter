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

		if strings.HasPrefix(received, commonbl.STATUS_REQUEST) {
			handleStatusRequest(pipeHandler, received)
		}
		if strings.HasPrefix(received, commonbl.CONNECTIONS_REQUEST) {
			handleConnetionsRequest(pipeHandler, received)
		}

		time.Sleep(time.Millisecond)
	}

}

func handleStatusRequest(handler commonbl.PipeHandler, request string) {
	id := getIdFromRequest(request)
	if params.Verbose {
		fmt.Fprintln(os.Stdout, fmt.Sprintf("Handle STATUS_REQUEST %s", id))
	}

	if !params.Test {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Error: Productive code not implemented yet"))
		os.Exit(-2)
	} else {
		err := handler.WritePipeString(fmt.Sprintf("%s Test response for request %s", commonbl.STATUS_REQUEST, id))
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("Error while write \"%s\" response to pipe: %s", commonbl.STATUS_REQUEST, err))
			os.Exit(-1)
		}
	}
}

func handleConnetionsRequest(handler commonbl.PipeHandler, request string) {
	id := getIdFromRequest(request)
	if params.Verbose {
		fmt.Fprintln(os.Stdout, fmt.Sprintf("Handle CONNECTIONS_REQUEST %s", id))
	}

	if !params.Test {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Error: Productive code not implemented yet"))
		os.Exit(-2)
	} else {
		err := handler.WritePipeString(fmt.Sprintf("%s Test response for request %s", commonbl.CONNECTIONS_REQUEST, id))
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("Error while write \"%s\" response to pipe: %s", commonbl.CONNECTIONS_REQUEST, err))
			os.Exit(-1)
		}
	}
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
