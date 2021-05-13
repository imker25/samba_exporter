package commonbl

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"os"
	"testing"
)

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
