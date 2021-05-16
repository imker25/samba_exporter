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

// The timeout for a request to samba_statusd in seconds
const requestTimeOut = 2

// The version of this program, will be set at compile time by the gradle build script
var version = "undefined"
var requestCount = 0

type SmbResponse struct {
	Data  string
	Error error
}

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

	res, _ := getSmbStatusDataTimeOut(pipeHander, commonbl.PROCESS_REQUEST)
	fmt.Fprintln(os.Stdout, res)
	res, _ = getSmbStatusDataTimeOut(pipeHander, commonbl.SERVICE_REQUEST)
	fmt.Fprintln(os.Stdout, res)
	res, _ = getSmbStatusDataTimeOut(pipeHander, commonbl.LOCK_REQUEST)
	fmt.Fprintln(os.Stdout, res)

	os.Exit(0)
}

func getSmbStatusDataTimeOut(handler commonbl.PipeHandler, request string) (string, error) {
	c := make(chan SmbResponse, 1)
	var data string
	go goGetSmbStatusData(handler, request, c)
	select {
	case res := <-c:
		if res.Error == nil {
			data = res.Data
		} else {
			return "", res.Error
		}
	case <-time.After(requestTimeOut * time.Second):
		return "", &time.ParseError{} // ToDo: Write a own error type
	}

	return data, nil
}

func goGetSmbStatusData(handler commonbl.PipeHandler, request string, c chan SmbResponse) {
	retStr, err := getSmbStatusData(handler, request)

	ret := SmbResponse{retStr, err}

	c <- ret
}

func getSmbStatusData(handler commonbl.PipeHandler, request string) (string, error) {
	requestCount++
	requestString := fmt.Sprintf("%s %d", request, requestCount)
	errWrite := handler.WritePipeString(requestString)
	if errWrite != nil {
		return "", errWrite
	}

	var errRead error
	response := requestString
	for response == requestString && errRead == nil {
		time.Sleep(time.Millisecond)
		response, errRead = handler.WaitForPipeInputString()
	}
	if errRead != nil {
		return "", errRead
	}

	if !strings.Contains(response, request) &&
		!strings.Contains(response, fmt.Sprintf("Response for request %d", requestCount)) {
		return "", commonbl.NewReaderError(response)
	}

	return response, nil
}

// Prints the version string
func printVersion() {
	fmt.Fprintln(os.Stdout, getVersion())
}

// Get the version string
func getVersion() string {
	return fmt.Sprintf("Version: %s", version)
}
