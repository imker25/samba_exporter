package main

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"tobi.backfrak.de/internal/commonbl"
	"tobi.backfrak.de/internal/smbexporterbl/pipecomunication"
	"tobi.backfrak.de/internal/smbexporterbl/smbstatusreader"
)

var mMutext sync.Mutex

func TestGetVersion(t *testing.T) {
	version := getVersion()

	if !strings.Contains(version, "Version:") {
		t.Errorf("The version string has not the expected format")
	}

	printVersion()
}

func TestTestPipeMode(t *testing.T) {
	mMutext.Lock()
	defer mMutext.Unlock()

	oldParmas := params
	defer func() { params = oldParmas }()

	params.RequestTimeOut = 1
	requestHandler := commonbl.NewPipeHandler(true, commonbl.RequestPipe)
	responseHandler := commonbl.NewPipeHandler(true, commonbl.ResposePipe)
	err := testPipeMode(requestHandler, responseHandler)
	if err == nil {
		t.Errorf("Got no error, but expected one")
	}

	switch err.(type) {
	case *pipecomunication.SmbStatusTimeOutError:
		fmt.Println("OK")
	default:
		t.Errorf("Got error of type '%s', but expected type '*pipecomunication.SmbStatusTimeOutError'", err)
	}

}

func TestHandleTestResponse(t *testing.T) {
	mMutext.Lock()
	defer mMutext.Unlock()

	oldParmas := params
	defer func() { params = oldParmas }()

	logger := commonbl.NewLogger(true)
	shares := smbstatusreader.GetShareData(commonbl.TestShareResponse, logger)
	processes := smbstatusreader.GetProcessData(commonbl.TestProcessResponse, logger)
	locks := smbstatusreader.GetLockData(commonbl.TestLockResponse, logger)
	psData := smbstatusreader.GetPsData(commonbl.TestPsResponse(), logger)

	handleTestResponse(processes, shares, locks, psData)

}

func TestMainWithHelp(t *testing.T) {
	mMutext.Lock()
	defer mMutext.Unlock()

	oldParmas := params
	defer func() { params = oldParmas }()

	params.Test = true
	params.Help = true

	res := realMain()
	if res != 0 {
		t.Errorf("Got %d from main, but expected 0", res)
	}

}

func TestMainWithVerbose(t *testing.T) {
	mMutext.Lock()
	defer mMutext.Unlock()

	oldParmas := params
	defer func() { params = oldParmas }()

	params.Test = true
	params.Help = true
	params.Verbose = true

	res := realMain()
	if res != 0 {
		t.Errorf("Got %d from main, but expected 0", res)
	}

}

func TestMainWithPrintVersion(t *testing.T) {
	mMutext.Lock()
	defer mMutext.Unlock()

	oldParmas := params
	defer func() { params = oldParmas }()

	params.Test = true
	params.PrintVersion = true

	res := realMain()
	if res != 0 {
		t.Errorf("Got %d from main, but expected 0", res)
	}

}

func TestMainWithTestPipeMode(t *testing.T) {
	mMutext.Lock()
	defer mMutext.Unlock()

	oldParmas := params
	defer func() { params = oldParmas }()

	params.Test = true
	params.TestPipeMode = true

	res := realMain()
	if res != -2 {
		t.Errorf("Got %d from main, but expected -2", res)
	}

}
