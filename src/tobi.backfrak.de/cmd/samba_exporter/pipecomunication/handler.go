package pipecomunication

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"fmt"
	"os"
	"sync"
	"time"

	"tobi.backfrak.de/cmd/samba_exporter/smbstatusreader"
	"tobi.backfrak.de/internal/commonbl"
)

// The timeout for a request to samba_statusd in seconds
const requestTimeOut = 2

var requestCount = 0
var mux sync.Mutex

type smbResponse struct {
	Data  string
	Error error
}

// GetSambaStatus - Get the output of all data tables from samba_statusd
func GetSambaStatus(requestHandler commonbl.PipeHandler, responseHandler commonbl.PipeHandler, verbose bool) ([]smbstatusreader.LockData, []smbstatusreader.ProcessData, []smbstatusreader.ShareData, error) {
	var processes []smbstatusreader.ProcessData
	var shares []smbstatusreader.ShareData
	var locks []smbstatusreader.LockData

	res, errGet := getSmbStatusDataTimeOut(requestHandler, responseHandler, commonbl.PROCESS_REQUEST, verbose)
	if errGet != nil {
		return nil, nil, nil, errGet
	} else {
		processes = smbstatusreader.GetProcessData(res)
		if len(processes) != 1 {
			return nil, nil, nil, NewSmbStatusUnexpectedResponseError(res)
		}
	}
	res, errGet = getSmbStatusDataTimeOut(requestHandler, responseHandler, commonbl.SHARE_REQUEST, verbose)
	if errGet != nil {
		return nil, nil, nil, errGet
	} else {
		shares = smbstatusreader.GetShareData(res)
		if len(shares) != 1 {
			return nil, nil, nil, NewSmbStatusUnexpectedResponseError(res)
		}
	}
	res, errGet = getSmbStatusDataTimeOut(requestHandler, responseHandler, commonbl.LOCK_REQUEST, verbose)
	if errGet != nil {
		return nil, nil, nil, errGet
	} else {
		locks = smbstatusreader.GetLockData(res)
		if len(locks) != 1 {
			return nil, nil, nil, NewSmbStatusUnexpectedResponseError(res)
		}
	}

	return locks, processes, shares, nil
}

func getSmbStatusDataTimeOut(requestHandler commonbl.PipeHandler, responseHandler commonbl.PipeHandler, request commonbl.RequestType, verbose bool) (string, error) {
	c := make(chan smbResponse, 1)
	var data string
	go goGetSmbStatusData(requestHandler, responseHandler, request, verbose, c)
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

func goGetSmbStatusData(requestHandler commonbl.PipeHandler, responseHandler commonbl.PipeHandler, request commonbl.RequestType, verbose bool, c chan smbResponse) {
	retStr, err := getSmbStatusData(requestHandler, responseHandler, request, verbose)

	ret := smbResponse{retStr, err}

	c <- ret
}

func getSmbStatusData(requestHandler commonbl.PipeHandler, responseHandler commonbl.PipeHandler, request commonbl.RequestType, verbose bool) (string, error) {
	requestCount++
	requestString := commonbl.GetRequest(request, requestCount)

	// Ensure we run only one request per time on the pipes
	mux.Lock()
	defer mux.Unlock()
	if verbose {
		fmt.Fprintln(os.Stdout, fmt.Sprintf("Send \"%s\" request on pipe", request))
	}
	errWrite := requestHandler.WritePipeString(requestString)
	if errWrite != nil {
		return "", errWrite
	}

	if verbose {
		fmt.Fprintln(os.Stdout, fmt.Sprintf("Wait for \"%s\" response on pipe", request))
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
	if verbose {
		fmt.Fprintln(os.Stdout, fmt.Sprintf("Handle \"%s\" response from pipe", request))
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
