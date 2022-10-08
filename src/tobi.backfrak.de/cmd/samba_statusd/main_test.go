package main

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"strings"
	"sync"
	"testing"

	"tobi.backfrak.de/internal/commonbl"
)

var mMutext sync.Mutex

func TestGetVersion(t *testing.T) {
	version := getVersion()

	if !strings.Contains(version, "Version:") {
		t.Errorf("The version string has not the expected format")
	}

	printVersion()
}

func TestTestPsResponse(t *testing.T) {
	mMutext.Lock()
	defer mMutext.Unlock()

	oldParmas := params
	defer func() { params = oldParmas }()
	responseHandler := commonbl.NewPipeHandler(true, commonbl.ResposePipe)

	err := testPsResponse(responseHandler, 0)
	if err != nil {
		t.Errorf("Get error '%s' but expected none", err.Error())
	}
}

func TestTestProcessResponse(t *testing.T) {
	mMutext.Lock()
	defer mMutext.Unlock()

	oldParmas := params
	defer func() { params = oldParmas }()
	responseHandler := commonbl.NewPipeHandler(true, commonbl.ResposePipe)

	err := testProcessResponse(responseHandler, 10)
	if err != nil {
		t.Errorf("Get error '%s' but expected none", err.Error())
	}
}

func TestTestLockResponse(t *testing.T) {
	mMutext.Lock()
	defer mMutext.Unlock()

	oldParmas := params
	defer func() { params = oldParmas }()
	responseHandler := commonbl.NewPipeHandler(true, commonbl.ResposePipe)

	err := testLockResponse(responseHandler, 20)
	if err != nil {
		t.Errorf("Get error '%s' but expected none", err.Error())
	}
}

func TestTestShareResponse(t *testing.T) {
	mMutext.Lock()
	defer mMutext.Unlock()

	oldParmas := params
	defer func() { params = oldParmas }()
	responseHandler := commonbl.NewPipeHandler(true, commonbl.ResposePipe)

	err := testShareResponse(responseHandler, 30)
	if err != nil {
		t.Errorf("Get error '%s' but expected none", err.Error())
	}
}

func TestHandleRequest(t *testing.T) {
	mMutext.Lock()
	defer mMutext.Unlock()

	oldParmas := params
	defer func() { params = oldParmas }()
	responseHandler := commonbl.NewPipeHandler(true, commonbl.ResposePipe)

	errNil := handleRequest(responseHandler,
		commonbl.GetRequest(commonbl.LOCK_REQUEST, 12),
		commonbl.LOCK_REQUEST,
		func(ph *commonbl.PipeHandler, i int) error { return nil },
		func(ph *commonbl.PipeHandler, i int) error { return nil },
	)

	if errNil != nil {
		t.Errorf("Get error '%s' but expected none", errNil.Error())
	}

	params.Test = true
	errNil = handleRequest(responseHandler,
		commonbl.GetRequest(commonbl.LOCK_REQUEST, 12),
		commonbl.LOCK_REQUEST,
		func(ph *commonbl.PipeHandler, i int) error { return nil },
		func(ph *commonbl.PipeHandler, i int) error { return nil },
	)

	if errNil != nil {
		t.Errorf("Get error '%s' but expected none", errNil.Error())
	}

	errRequest := commonbl.NewEmptyStringQueueError()
	errHandle := handleRequest(responseHandler,
		commonbl.GetRequest(commonbl.LOCK_REQUEST, 12),
		commonbl.LOCK_REQUEST,
		func(ph *commonbl.PipeHandler, i int) error { return nil },
		func(ph *commonbl.PipeHandler, i int) error { return errRequest },
	)

	if errHandle != errRequest {
		t.Errorf("Got error '%s', but expected error '%s'", errHandle.Error(), errRequest.Error())
	}

	params.Test = false
	errHandle = handleRequest(responseHandler,
		commonbl.GetRequest(commonbl.LOCK_REQUEST, 12),
		commonbl.LOCK_REQUEST,
		func(ph *commonbl.PipeHandler, i int) error { return errRequest },
		func(ph *commonbl.PipeHandler, i int) error { return nil },
	)

	if errHandle != errRequest {
		t.Errorf("Got error '%s', but expected error '%s'", errHandle.Error(), errRequest.Error())
	}
}

func TestGoHandleRequestQueue(t *testing.T) {
	mMutext.Lock()
	defer mMutext.Unlock()

	oldParmas := params
	defer func() { params = oldParmas }()
	responseHandler := commonbl.NewPipeHandler(true, commonbl.ResposePipe)
	requestQueue = *commonbl.NewStringQueue()
	params.Test = true

	requestQueue.Push(commonbl.GetRequest(commonbl.LOCK_REQUEST, 0))
	goHandleRequestQueue(responseHandler)

	requestQueue.Push(commonbl.GetRequest(commonbl.SHARE_REQUEST, 1))
	goHandleRequestQueue(responseHandler)

	requestQueue.Push(commonbl.GetRequest(commonbl.PROCESS_REQUEST, 2))
	goHandleRequestQueue(responseHandler)

	requestQueue.Push(commonbl.GetRequest(commonbl.PS_REQUEST, 3))
	goHandleRequestQueue(responseHandler)

	requestQueue.Push(commonbl.GetRequest("NO_REQUEST", 3))
	goHandleRequestQueue(responseHandler)

	requestQueue.Push("")
	goHandleRequestQueue(responseHandler)
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
