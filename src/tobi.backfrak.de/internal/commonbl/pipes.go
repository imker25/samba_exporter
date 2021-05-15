package commonbl

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"
)

const testPipeFileName = "samba_exporter.pipe"
const pipePath = "/run"
const testPipePath = "/dev/shm"
const pipePermission = 0660
const endByte byte = 0

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
		errCreate := handler.createPipe()
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
		errCreate := handler.createPipe()
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

// WaitForPipeInputBytes - Blocking! Wait for input in the pipe and return it as byte array
// The array will be empty in case of errors
func (handler *PipeHandler) WaitForPipeInputBytes() ([]byte, error) {
	reader, errGet := handler.GetReaderPipe()
	if errGet != nil {
		return []byte{}, errGet
	}
	received, errRead := reader.ReadBytes(endByte)
	if errRead != nil {
		if errRead != io.EOF {
			return []byte{}, errRead
		}
		return []byte{}, nil
	}

	return received[0 : len(received)-1], nil
}

// WaitForPipeInputString - Blocking! Wait for input in the pipe and return it as string
// The string will be empty in case of errors
func (handler *PipeHandler) WaitForPipeInputString() (string, error) {
	data, err := handler.WaitForPipeInputBytes()

	return strings.TrimSpace(string(data)), err
}

// WritePipeBytes - Write byte data to the pipe
func (handler *PipeHandler) WritePipeBytes(data []byte) error {
	writer, errGet := handler.GetWriterPipe()
	if errGet != nil {
		return errGet
	}
	data = append(data, endByte)
	_, errWrite := writer.Write(data)
	if errWrite != nil {
		return errWrite
	}
	errFlush := writer.Flush()
	if errFlush != nil {
		return errFlush
	}

	return nil
}

// WritePipeString - Write string data to the pipe
func (handler *PipeHandler) WritePipeString(data string) error {
	return handler.WritePipeBytes([]byte(data))
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

func (handler *PipeHandler) createPipe() error {
	return syscall.Mkfifo(handler.GetPipeFilePath(), pipePermission)
}
