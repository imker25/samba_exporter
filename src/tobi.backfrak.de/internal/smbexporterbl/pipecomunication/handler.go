package pipecomunication

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"fmt"
	"sync"
	"time"

	"tobi.backfrak.de/internal/commonbl"
	"tobi.backfrak.de/internal/smbexporterbl/smbstatusreader"
)

var requestCount = 0
var requestMux sync.Mutex
var collectMux sync.Mutex

type smbResponse struct {
	Data  string
	Error error
}

// GetSambaStatus - Get the output of all data tables from samba_statusd
func GetSambaStatus(requestHandler *commonbl.PipeHandler, responseHandler *commonbl.PipeHandler, logger commonbl.Logger, requestTimeOut int) ([]smbstatusreader.LockData, []smbstatusreader.ProcessData, []smbstatusreader.ShareData, []commonbl.PsUtilPidData, error) {
	var processes []smbstatusreader.ProcessData
	var shares []smbstatusreader.ShareData
	var locks []smbstatusreader.LockData
	var psdata []commonbl.PsUtilPidData
	sharesChan := make(chan []smbstatusreader.ShareData, 1)
	processesChan := make(chan []smbstatusreader.ProcessData, 1)
	locksChan := make(chan []smbstatusreader.LockData, 1)
	psdataChan := make(chan []commonbl.PsUtilPidData, 1)
	collectMux.Lock()
	defer collectMux.Unlock()

	res, errGet := getSmbStatusDataTimeOut(requestHandler, responseHandler, commonbl.PROCESS_REQUEST, logger, requestTimeOut)
	if errGet != nil {
		return nil, nil, nil, nil, errGet
	}
	go goGetProcessData(res, logger, processesChan)

	res, errGet = getSmbStatusDataTimeOut(requestHandler, responseHandler, commonbl.SHARE_REQUEST, logger, requestTimeOut)
	if errGet != nil {
		return nil, nil, nil, nil, errGet
	}
	go goGetShareData(res, logger, sharesChan)

	res, errGet = getSmbStatusDataTimeOut(requestHandler, responseHandler, commonbl.LOCK_REQUEST, logger, requestTimeOut)
	if errGet != nil {
		return nil, nil, nil, nil, errGet
	}
	go goGetLockData(res, logger, locksChan)

	res, errGet = getSmbStatusDataTimeOut(requestHandler, responseHandler, commonbl.PS_REQUEST, logger, requestTimeOut)
	if errGet != nil {
		return nil, nil, nil, nil, errGet
	}
	go goGetPsData(res, logger, psdataChan)

	processes = <-processesChan
	shares = <-sharesChan
	locks = <-locksChan
	psdata = <-psdataChan

	if len(shares) < 1 {
		logger.WriteVerbose("Got an empty share table when requesting \"smbstatus -S -n\" from samba_statusd")
	}

	if len(processes) < 1 {
		logger.WriteVerbose("Got an empty process table when requesting \"smbstatus -p -n\" from samba_statusd")
	}

	return locks, processes, shares, psdata, nil
}

func goGetProcessData(res string, logger commonbl.Logger, c chan []smbstatusreader.ProcessData) {
	processes := smbstatusreader.GetProcessData(res, logger)

	c <- processes
}

func goGetShareData(res string, logger commonbl.Logger, c chan []smbstatusreader.ShareData) {
	shares := smbstatusreader.GetShareData(res, logger)

	c <- shares
}

func goGetLockData(res string, logger commonbl.Logger, c chan []smbstatusreader.LockData) {
	locks := smbstatusreader.GetLockData(res, logger)

	c <- locks
}

func goGetPsData(res string, logger commonbl.Logger, c chan []commonbl.PsUtilPidData) {
	locks := smbstatusreader.GetPsData(res, logger)

	c <- locks
}

func getSmbStatusDataTimeOut(requestHandler *commonbl.PipeHandler, responseHandler *commonbl.PipeHandler, request commonbl.RequestType, logger commonbl.Logger, requestTimeOut int) (string, error) {
	c := make(chan smbResponse, 1)
	var data string

	go goGetSmbStatusData(requestHandler, responseHandler, request, logger, c)
	select {
	case res := <-c:
		if res.Error == nil {
			data = res.Data
		} else {
			return "", res.Error
		}
	case <-time.After(time.Second * time.Duration(requestTimeOut)):
		logger.WriteVerbose("Clear request pipe after request time out")
		errClear := requestHandler.WritePipeString("")
		if errClear != nil {
			panic(errClear)
		}
		return "", NewSmbStatusTimeOutError(request)
	}

	return data, nil
}

func goGetSmbStatusData(requestHandler *commonbl.PipeHandler, responseHandler *commonbl.PipeHandler, request commonbl.RequestType, logger commonbl.Logger, c chan smbResponse) {
	retStr, err := getSmbStatusData(requestHandler, responseHandler, request, logger)

	ret := smbResponse{retStr, err}

	c <- ret
}

func getSmbStatusData(requestHandler *commonbl.PipeHandler, responseHandler *commonbl.PipeHandler, request commonbl.RequestType, logger commonbl.Logger) (string, error) {
	// Ensure we run only one request per time on the pipes
	requestMux.Lock()
	defer requestMux.Unlock()
	requestCount++
	requestString := commonbl.GetRequest(request, requestCount)

	logger.WriteVerbose(fmt.Sprintf("Send \"%s\" request with ID %d on pipe", request, requestCount))

	errWrite := requestHandler.WritePipeString(requestString)
	if errWrite != nil {
		return "", errWrite
	}

	logger.WriteVerbose(fmt.Sprintf("Wait for \"%s\" response with ID %d on pipe", request, requestCount))

	response := ""
	var errRead error

	// as long as the response pipe is empty, wait for response
	for response == "" {
		response, errRead = responseHandler.WaitForPipeInputString()
		if errRead != nil {
			return "", errRead
		}
	}

	logger.WriteVerbose(fmt.Sprintf("Handle \"%s\" response with ID %d from pipe", request, requestCount))

	header, data, errSplit := commonbl.SplitResponse(response)
	if errSplit != nil {
		return "", errSplit
	}

	if !commonbl.CheckResponseHeader(header, request, requestCount) {
		return "", commonbl.NewReaderError(response, request, requestCount)
	}

	return data, nil
}
