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

func TestGetIdFromRequest(t *testing.T) {
	id, err := getIdFromRequest("bal: 23")

	if err != nil {
		t.Errorf("Got error \"%s\" but expected none", err)
	}

	if id != 23 {
		t.Errorf("The id \"%d\" is not the expected", id)
	}

	id, err = getIdFromRequest("bal: 23: sert")
	if err == nil {
		t.Errorf("Got no error but expected one")
	}

	if id != 0 {
		t.Errorf("The id \"%d\" is not the expected", id)
	}

	id, err = getIdFromRequest("bal: 23  sert")
	if err == nil {
		t.Errorf("Got no error but expected one")
	}

	if id != 0 {
		t.Errorf("The id \"%d\" is not the expected", id)
	}
}
