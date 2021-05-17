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
	requestHandler := *commonbl.NewPipeHandler(params.Test, commonbl.RequestPipe)
	responseHandler := *commonbl.NewPipeHandler(params.Test, commonbl.ResposePipe)
	if params.Verbose {
		args := ""
		for _, arg := range os.Args {
			args = fmt.Sprintf("%s %s", args, arg)
		}
		fmt.Fprintln(os.Stdout, fmt.Sprintf("Call: %s", args))
		if !params.PrintVersion {
			printVersion()
		}
		fmt.Fprintln(os.Stdout, fmt.Sprintf("Named pipe for requests: %s", requestHandler.GetPipeFilePath()))
		fmt.Fprintln(os.Stdout, fmt.Sprintf("Named pipe for response: %s", responseHandler.GetPipeFilePath()))
	}

	if params.PrintVersion {
		printVersion()
		os.Exit(0)
	}

	if params.Help {
		flag.Usage()
		os.Exit(0)
	}

	res, errGet := getSmbStatusDataTimeOut(requestHandler, responseHandler, commonbl.PROCESS_REQUEST)
	if errGet != nil {
		fmt.Fprintln(os.Stderr, errGet)
	} else {
		fmt.Fprintln(os.Stdout, res)
	}
	res, errGet = getSmbStatusDataTimeOut(requestHandler, responseHandler, commonbl.SERVICE_REQUEST)
	if errGet != nil {
		fmt.Fprintln(os.Stderr, errGet)
	} else {
		fmt.Fprintln(os.Stdout, res)
	}
	res, errGet = getSmbStatusDataTimeOut(requestHandler, responseHandler, commonbl.LOCK_REQUEST)
	if errGet != nil {
		fmt.Fprintln(os.Stderr, errGet)
	} else {
		fmt.Fprintln(os.Stdout, res)
	}

	os.Exit(0)
}

func getSmbStatusDataTimeOut(requestHandler commonbl.PipeHandler, responseHandler commonbl.PipeHandler, request string) (string, error) {
	c := make(chan SmbResponse, 1)
	var data string
	go goGetSmbStatusData(requestHandler, responseHandler, request, c)
	select {
	case res := <-c:
		if res.Error == nil {
			data = res.Data
		} else {
			return "", res.Error
		}
	case <-time.After(requestTimeOut * time.Second):
		return "", NewSmbStatusTimeOutError(request)
	}

	return data, nil
}

func goGetSmbStatusData(requestHandler commonbl.PipeHandler, responseHandler commonbl.PipeHandler, request string, c chan SmbResponse) {
	retStr, err := getSmbStatusData(requestHandler, responseHandler, request)

	ret := SmbResponse{retStr, err}

	c <- ret
}

func getSmbStatusData(requestHandler commonbl.PipeHandler, responseHandler commonbl.PipeHandler, request string) (string, error) {
	requestCount++
	requestString := fmt.Sprintf("%s %d", request, requestCount)
	errWrite := requestHandler.WritePipeString(requestString)
	if errWrite != nil {
		return "", errWrite
	}

	var errRead error
	response := requestString
	for response == requestString && errRead == nil {
		time.Sleep(time.Millisecond)
		response, errRead = responseHandler.WaitForPipeInputString()
	}
	if errRead != nil {
		return "", errRead
	}

	splitResponse := strings.SplitN(response, "\n", 2)
	header := splitResponse[0]
	response = splitResponse[1]

	if !strings.Contains(header, request) &&
		!strings.Contains(header, fmt.Sprintf("Response for request %d", requestCount)) {
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
