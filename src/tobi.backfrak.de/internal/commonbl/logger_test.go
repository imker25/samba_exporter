package commonbl

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"testing"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger(false)
	if logger.Verbose == true {
		t.Errorf("Logger is verbose but should not")
	}

	logger = NewLogger(true)
	if logger.Verbose == false {
		t.Errorf("Logger is not verbose but should")
	}
}

func TestWriteInformation(t *testing.T) {
	logger := NewLogger(false)
	logger.WriteInformation("My message")
}

func TestWriteErrorMessage(t *testing.T) {
	logger := NewLogger(false)
	logger.WriteErrorMessage("My message")
}

func TestWriteVerbose(t *testing.T) {
	logger := NewLogger(false)
	logger.WriteVerbose("My message 1")

	logger = NewLogger(true)
	logger.WriteVerbose("My message 2")
}

func TestWriteError(t *testing.T) {
	logger := NewLogger(false)
	err := NewReaderError("my data")

	logger.WriteError(err)
}
