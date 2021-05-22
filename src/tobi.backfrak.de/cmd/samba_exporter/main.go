package main

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"flag"
	"fmt"
	"os"
	"time"

	"tobi.backfrak.de/cmd/samba_exporter/smbstatusreader"
	"tobi.backfrak.de/cmd/samba_exporter/statisticsGenerator"
	"tobi.backfrak.de/internal/commonbl"
)

// Authors - Information about the authors of the program. You might want to add your name here when contributing to this software
const Authors = "tobi@backfrak.de"

// The timeout for a request to samba_statusd in seconds
const requestTimeOut = 2
const exporter_label_prefix = "samba"

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

	if params.TestPipeMode {
		errTest := testPipeMode(requestHandler, responseHandler)
		if errTest != nil {
			fmt.Fprintln(os.Stderr, errTest)
		}
		os.Exit(0)
	}
	os.Exit(0)
}

func testPipeMode(requestHandler commonbl.PipeHandler, responseHandler commonbl.PipeHandler) error {
	var processes []smbstatusreader.ProcessData
	var shares []smbstatusreader.ShareData
	var locks []smbstatusreader.LockData
	res, errGet := getSmbStatusDataTimeOut(requestHandler, responseHandler, commonbl.PROCESS_REQUEST)
	if errGet != nil {
		return errGet
	} else {
		processes = smbstatusreader.GetProcessData(res)
		if len(processes) != 1 {
			return NewSmbStatusUnexpectedResponseError(res)
		}
		fmt.Fprintln(os.Stdout, processes[0].String())
	}
	res, errGet = getSmbStatusDataTimeOut(requestHandler, responseHandler, commonbl.SHARE_REQUEST)
	if errGet != nil {
		return errGet
	} else {
		shares = smbstatusreader.GetShareData(res)
		if len(shares) != 1 {
			return NewSmbStatusUnexpectedResponseError(res)
		}
		fmt.Fprintln(os.Stdout, shares[0].String())
	}
	res, errGet = getSmbStatusDataTimeOut(requestHandler, responseHandler, commonbl.LOCK_REQUEST)
	if errGet != nil {
		return errGet
	} else {
		locks = smbstatusreader.GetLockData(res)
		if len(locks) != 1 {
			return NewSmbStatusUnexpectedResponseError(res)
		}
		fmt.Fprintln(os.Stdout, locks[0].String())
	}

	stats := statisticsGenerator.GetSmbStatistics(locks, processes, shares)
	for _, stat := range stats {
		fmt.Fprintln(os.Stdout, fmt.Sprintf("%s_%s: %d", exporter_label_prefix, stat.Name, stat.Value))
	}

	return nil
}

func getSmbStatusDataTimeOut(requestHandler commonbl.PipeHandler, responseHandler commonbl.PipeHandler, request commonbl.RequestType) (string, error) {
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

func goGetSmbStatusData(requestHandler commonbl.PipeHandler, responseHandler commonbl.PipeHandler, request commonbl.RequestType, c chan SmbResponse) {
	retStr, err := getSmbStatusData(requestHandler, responseHandler, request)

	ret := SmbResponse{retStr, err}

	c <- ret
}

func getSmbStatusData(requestHandler commonbl.PipeHandler, responseHandler commonbl.PipeHandler, request commonbl.RequestType) (string, error) {
	requestCount++
	requestString := commonbl.GetRequest(request, requestCount)
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

	header, data, errSplit := commonbl.SplitResponse(response)
	if errSplit != nil {
		return "", errSplit
	}

	if !commonbl.CheckResponseHeader(header, request, requestCount) {
		return "", commonbl.NewReaderError(response)
	}

	return data, nil
}

// Prints the version string
func printVersion() {
	fmt.Fprintln(os.Stdout, getVersion())
}

// Get the version string
func getVersion() string {
	return fmt.Sprintf("Version: %s", version)
}
