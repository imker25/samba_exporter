package main

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"strings"
	"testing"
)

func TestSmbStatusTimeOutError(t *testing.T) {
	path := "/some/sample/path"
	err := NewSmbStatusTimeOutError(path)

	if err.Request != path {
		t.Errorf("The File was %s, but %s was expected", err.Request, path)
	}

	if strings.Contains(err.Error(), path) == false {
		t.Errorf("The error message of SmbStatusTimeOutError does not contain the expected request")
	}
}
