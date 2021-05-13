package commonbl

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"fmt"
	"os"
)

const testPipeFileName = "samba_exporter.pipe"
const pipePath = "/run"
const testPipePath = "/dev/shm"

// PipeHandler - Type to handle the pipe for comunication between samba_exporter and samba_statusd
type PipeHandler struct {
	TestMode bool
}

// NewPipeHandler - Get a new instance of the PipeHandler type
func NewPipeHandler(testMode bool) *PipeHandler {
	retVal := PipeHandler{}
	retVal.TestMode = testMode

	return &retVal
}

// GetPipeFilePath -  Get the path to the named pipe files for this application
func (handler *PipeHandler) GetPipeFilePath() string {
	var dirname string
	if handler.TestMode {
		dirname = testPipePath
	} else {
		dirname = pipePath
	}

	return fmt.Sprintf("%s/%s", dirname, testPipeFileName)
}

// PipeExists - Check if the named pipe files for this application exists
func (handler *PipeHandler) PipeExists() bool {
	return FileExists(handler.GetPipeFilePath())
}

// FileExists - Check if a file exists. Return false in case the path does not exist or is a directory
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	if info.IsDir() {
		return false
	}
	return true
}
