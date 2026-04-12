package commonbl

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"testing"
)

func TestNewConsoleLogger(t *testing.T) {
	logger := NewConsoleLogger(Information)

	if logger.LogLevelSetting != Information {
		t.Errorf("The logger.LogLevelSetting is %s but %s is expected", ValidLogLevelSettings[logger.LogLevelSetting], ValidLogLevelSettings[Information])
	}

	logger = NewConsoleLogger(Verbose)

	if logger.LogLevelSetting != Verbose {
		t.Errorf("The logger.LogLevelSetting is %s but %s is expected", ValidLogLevelSettings[logger.LogLevelSetting], ValidLogLevelSettings[Verbose])
	}
}

func TestWriteInformation(t *testing.T) {
	logger := NewConsoleLogger(Verbose)
	logger.WriteInformation("My message")
}

func TestWriteErrorMessage(t *testing.T) {
	logger := NewConsoleLogger(Verbose)
	logger.WriteErrorMessage("My message")
}

func TestWriteWarningMessage(t *testing.T) {
	logger := NewConsoleLogger(Verbose)
	logger.WriteWarning("My message")
}

func TestWriteVerbose(t *testing.T) {
	logger := NewConsoleLogger(Information)
	logger.WriteVerbose("My message 1")

	logger = NewConsoleLogger(Verbose)
	logger.WriteVerbose("My message 2")
}

func TestWriteError(t *testing.T) {
	logger := NewConsoleLogger(Verbose)
	err := NewReaderError("my data", LOCK_REQUEST, 3)

	logger.WriteError(err)
}

func TestWriteErrorWithAddition(t *testing.T) {
	logger := NewConsoleLogger(Information)
	err := NewReaderError("my data", LOCK_REQUEST, 3)

	logger.WriteErrorWithAddition(err, "additional data")
}
