package commonbl

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"testing"
)

func TestNewConsoleLogger(t *testing.T) {
	logger := NewConsoleLogger(false)
	if logger.Verbose == true {
		t.Errorf("Logger is verbose but should not")
	}

	logger = NewConsoleLogger(true)
	if logger.Verbose == false {
		t.Errorf("Logger is not verbose but should")
	}

	iLogger := Logger(logger)

	if iLogger.GetVerbose() == false {
		t.Errorf("Logger is not verbose but should")
	}

	logger.Verbose = false
	if iLogger.GetVerbose() == true {
		t.Errorf("Logger is verbose but not should")
	}
}

func TestWriteInformation(t *testing.T) {
	logger := NewConsoleLogger(false)
	logger.WriteInformation("My message")
}

func TestWriteErrorMessage(t *testing.T) {
	logger := NewConsoleLogger(false)
	logger.WriteErrorMessage("My message")
}

func TestWriteVerbose(t *testing.T) {
	logger := NewConsoleLogger(false)
	logger.WriteVerbose("My message 1")

	logger = NewConsoleLogger(true)
	logger.WriteVerbose("My message 2")
}

func TestWriteError(t *testing.T) {
	logger := NewConsoleLogger(false)
	err := NewReaderError("my data", LOCK_REQUEST, 3)

	logger.WriteError(err)
}

func TestWriteErrorWithAddition(t *testing.T) {
	logger := NewConsoleLogger(false)
	err := NewReaderError("my data", LOCK_REQUEST, 3)

	logger.WriteErrorWithAddition(err, "additional data")
}
