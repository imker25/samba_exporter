package commonbl

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"os"
	"testing"
)

func TestGetPipeFilePath(t *testing.T) {
	path := GetPipeFilePath()

	if path == "" {
		t.Errorf("GetPipeFilePath is empty")
	}
}

func TestPipeFileExists(t *testing.T) {
	os.Remove(GetPipeFilePath())
	if PipeExists() == true {
		t.Errorf("PipeExists is true but should not")
	}

	os.Create(GetPipeFilePath())
	if PipeExists() == false {
		t.Errorf("PipeExists is false but should not")
	}

}
