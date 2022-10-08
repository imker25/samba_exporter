package main

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"testing"
)

func TestHandleComandlineOptions(t *testing.T) {
	mMutext.Lock()
	defer mMutext.Unlock()

	handleComandlineOptions()
	if params.PrintVersion {
		t.Errorf("params.PrintVersion is true, but should not")
	}

	if params.Verbose {
		t.Errorf("params.Verbose is true, but should not")
	}
}

func TestCustomHelpMessage(t *testing.T) {
	customHelpMessage()
}
