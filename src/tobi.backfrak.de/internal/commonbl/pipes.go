package commonbl

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"bufio"
	"fmt"
	"os"
	"syscall"
)

const testPipeFileName = "samba_exporter.pipe"
const pipePath = "/run"
const testPipePath = "/dev/shm"
const pipePermission = 0660

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

// GetReaderPipe - Get a new reader for the common pipe.
// 	Remember: This is a blocking call and will return once data can be read from the pipe
func (handler *PipeHandler) GetReaderPipe() (*bufio.Reader, error) {

	if !handler.PipeExists() {
		errCreate := syscall.Mkfifo(handler.GetPipeFilePath(), pipePermission)
		if errCreate != nil {
			return nil, errCreate
		}
	}

	file, errOpen := os.OpenFile(handler.GetPipeFilePath(), os.O_CREATE, os.ModeNamedPipe)
	if errOpen != nil {
		return nil, errOpen
	}

	return bufio.NewReader(file), nil

}

// GetWriterPipe - Get a new writer for the common pipe.
func (handler *PipeHandler) GetWriterPipe() (*bufio.Writer, error) {

	if !handler.PipeExists() {
		errCreate := syscall.Mkfifo(handler.GetPipeFilePath(), pipePermission)
		if errCreate != nil {
			return nil, errCreate
		}
	}

	file, errOpen := os.OpenFile(handler.GetPipeFilePath(), os.O_RDWR|os.O_CREATE|os.O_APPEND, pipePermission)
	if errOpen != nil {
		return nil, errOpen
	}

	return bufio.NewWriter(file), nil
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
