package commonbl

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"io"
	"os"
	"sync"
	"testing"
)

var testData = []byte{1, 3, 5, 7, 9, 0}
var flushed bool
var mux sync.Mutex

func TestNewPipeHandler(t *testing.T) {
	handler := NewPipeHandler(false)

	path := handler.GetPipeFilePath()
	if path != "/run/samba_exporter.pipe" {
		t.Errorf("GetPipeFilePath() has not the expected value")
	}
}

func TestGetPipeFilePath(t *testing.T) {
	handler := NewPipeHandler(true)
	path := handler.GetPipeFilePath()

	if path == "" {
		t.Errorf("GetPipeFilePath is empty")
	}
}

func TestPipeFileExists(t *testing.T) {
	handler := NewPipeHandler(true)

	os.Remove(handler.GetPipeFilePath())
	if handler.PipeExists() == true {
		t.Errorf("PipeExists is true but should not")
	}

	os.Create(handler.GetPipeFilePath())
	if handler.PipeExists() == false {
		t.Errorf("PipeExists is false but should not")
	}

}

func TestGetWriterPipe(t *testing.T) {
	handler := NewPipeHandler(true)
	defer os.Remove(handler.GetPipeFilePath())

	writer, err := handler.GetWriterPipe()
	if err != nil {
		t.Fatalf("Got error \"%s\" but expected none", err)
	}

	if writer == nil {
		t.Errorf("The writer is nil, but should but")
	}

}

func TestGetWriterPipeTwoTimes(t *testing.T) {
	handler := NewPipeHandler(true)
	defer os.Remove(handler.GetPipeFilePath())

	writer1, err1 := handler.GetWriterPipe()
	if err1 != nil {
		t.Fatalf("Got error \"%s\" but expected none", err1)
	}

	if writer1 == nil {
		t.Errorf("The writer is nil, but should but")
	}

	writer2, err2 := handler.GetWriterPipe()
	if err2 != nil {
		t.Fatalf("Got error \"%s\" but expected none", err2)
	}

	if writer2 == nil {
		t.Errorf("The writer is nil, but should but")
	}

}

func TestGetReaderPipe(t *testing.T) {
	handler := NewPipeHandler(true)
	defer os.Remove(handler.GetPipeFilePath())

	flushed = false
	mux.Lock()
	go scheduleWrite(t)

	reader, errGet := handler.GetReaderPipe()
	if errGet != nil {
		t.Fatalf("Got error \"%s\" but expected none", errGet)
	}

	mux.Lock()
	defer mux.Unlock()

	if flushed == false {
		t.Errorf("Data was not flushed yet")
	}

	data, errRead := reader.ReadBytes(0)
	if errRead != nil && errRead != io.EOF {
		t.Fatalf("Got error \"%s\" but expected none", errRead)
	}

	if len(data) != len(testData) {
		t.Errorf("The len of the received data %d does not match the send data %d", len(data), len(testData))
	}

	for i, _ := range data {
		if data[i] != testData[i] {
			t.Errorf("The received byte does not match the send byte at position %d", i)
		}
	}
}

func TestGetReaderPipeReaderReuse(t *testing.T) {
	handler := NewPipeHandler(true)
	defer os.Remove(handler.GetPipeFilePath())

	flushed = false
	mux.Lock()
	go scheduleWrite(t)

	reader, errGet := handler.GetReaderPipe()
	if errGet != nil {
		t.Fatalf("Got error \"%s\" but expected none", errGet)
	}

	mux.Lock()
	mux.Unlock()
	if flushed == false {
		t.Errorf("Data was not flushed yet")
	}

	data, errRead := reader.ReadBytes(0)
	if errRead != nil && errRead != io.EOF {
		t.Fatalf("Got error \"%s\" but expected none", errRead)
	}

	if len(data) != len(testData) {
		t.Errorf("The len of the received data %d does not match the send data %d", len(data), len(testData))
	}

	for i, _ := range data {
		if data[i] != testData[i] {
			t.Errorf("The received byte does not match the send byte at position %d", i)
		}
	}

	handler = NewPipeHandler(true)
	flushed = false
	mux.Lock()
	go scheduleWrite(t)
	mux.Lock()
	defer mux.Unlock()
	if flushed == false {
		t.Errorf("Data was not flushed yet")
	}

	data, errRead = reader.ReadBytes(0)
	if errRead != nil && errRead != io.EOF {
		t.Fatalf("Got error \"%s\" but expected none", errRead)
	}

	if len(data) != len(testData) {
		t.Errorf("The len of the received data %d does not match the send data %d", len(data), len(testData))
	}

	for i, _ := range data {
		if data[i] != testData[i] {
			t.Errorf("The received byte does not match the send byte at position %d", i)
		}
	}
}

func scheduleWrite(t *testing.T) {
	defer mux.Unlock()

	handler := NewPipeHandler(true)
	writer, errGet := handler.GetWriterPipe()
	if errGet != nil {
		t.Fatalf("Got error \"%s\" but expected none", errGet)
	}

	count, errWrite := writer.Write(testData)
	if errWrite != nil {
		t.Fatalf("Got error \"%s\" but expected none", errWrite)
	}

	if count != len(testData) {
		t.Errorf("Did not write the expected amount of data")
	}

	errFlush := writer.Flush()
	if errFlush != nil {
		t.Fatalf("Got error \"%s\" but expected none", errFlush)
	}
	flushed = true
	return
}
