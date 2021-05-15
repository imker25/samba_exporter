package commonbl

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"os"
	"sync"
	"testing"
)

var testDataString = "Hello Word"
var testData = []byte(testDataString)

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

func TestReadWriteData(t *testing.T) {
	handler := NewPipeHandler(true)
	defer os.Remove(handler.GetPipeFilePath())
	mux.Lock()
	go scheduleWriter(t)
	mux.Lock()
	defer mux.Unlock()

	data, err := handler.WaitForPipeInputBytes()
	if err != nil {
		t.Fatalf("Got error \"%s\" but expected none", err)
	}

	for i, _ := range data {
		if data[i] != testData[i] {
			t.Errorf("The received byte does not match the send byte at position %d", i)
		}
	}
}

func TestReadWriteStringData(t *testing.T) {
	handler := NewPipeHandler(true)
	defer os.Remove(handler.GetPipeFilePath())
	mux.Lock()
	go scheduleStringWriter(t)
	mux.Lock()
	defer mux.Unlock()

	data, err := handler.WaitForPipeInputString()
	if err != nil {
		t.Fatalf("Got error \"%s\" but expected none", err)
	}

	if data != testDataString {
		t.Errorf("The received string \"%s\" does not match the send string \"%s\"", data, testDataString)
	}
}

func TestReadWriteStringDataReuse(t *testing.T) {
	handler := NewPipeHandler(true)
	defer os.Remove(handler.GetPipeFilePath())
	mux.Lock()
	go scheduleStringWriter(t)
	mux.Lock()
	mux.Unlock()

	data, err := handler.WaitForPipeInputString()
	if err != nil {
		t.Fatalf("Got error \"%s\" but expected none", err)
	}

	if data != testDataString {
		t.Errorf("The received string \"%s\" does not match the send string \"%s\"", data, testDataString)
	}

	mux.Lock()
	go scheduleStringWriter(t)
	mux.Lock()
	defer mux.Unlock()

	data, err = handler.WaitForPipeInputString()
	if err != nil {
		t.Fatalf("Got error \"%s\" but expected none", err)
	}

	if data != testDataString {
		t.Errorf("The received string \"%s\" does not match the send string \"%s\"", data, testDataString)
	}
}

func scheduleWriter(t *testing.T) {
	defer mux.Unlock()

	handler := NewPipeHandler(true)
	err := handler.WritePipeBytes(testData)
	if err != nil {
		t.Fatalf("Got error \"%s\" but expected none", err)
	}
}

func scheduleStringWriter(t *testing.T) {
	defer mux.Unlock()

	handler := NewPipeHandler(true)
	err := handler.WritePipeString(testDataString)
	if err != nil {
		t.Fatalf("Got error \"%s\" but expected none", err)
	}
}
