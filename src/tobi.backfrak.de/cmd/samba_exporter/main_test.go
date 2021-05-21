package main

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"tobi.backfrak.de/internal/commonbl"
)

func TestGetVersion(t *testing.T) {
	version := getVersion()

	if !strings.Contains(version, "Version:") {
		t.Errorf("The version string has not the expected format")
	}
}

func TestPipeTestMode(t *testing.T) {
	requestHandler := *commonbl.NewPipeHandler(true, commonbl.RequestPipe)
	responseHandler := *commonbl.NewPipeHandler(true, commonbl.ResposePipe)
	err := testPipeMode(requestHandler, responseHandler)

	if err == nil {
		t.Errorf("Exptected an error but got none")
	}

	switch err.(type) {
	case *SmbStatusTimeOutError:
		fmt.Fprintln(os.Stdout, "OK")
	default:
		t.Errorf("Got error of the wrong type")
	}
}
