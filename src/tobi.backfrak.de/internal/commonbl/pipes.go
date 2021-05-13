package commonbl

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"os"
)

const testPipeFilePath = "./samba_exporter.pipe"

// Get the path to the named pipe files for this application
func GetPipeFilePath() string {
	return testPipeFilePath
}

// PipeExists - Check if the named pipe files for this application exists
func PipeExists() bool {
	return FileExists(GetPipeFilePath())
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
