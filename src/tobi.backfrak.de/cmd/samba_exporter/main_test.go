package main

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"strings"
	"testing"
)

func TestGetVersion(t *testing.T) {
	version := getVersion()

	if !strings.Contains(version, "Version:") {
		t.Errorf("The version string has not the expected format")
	}
}
