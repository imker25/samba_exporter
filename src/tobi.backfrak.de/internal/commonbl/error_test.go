package commonbl

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"strings"
	"testing"
)

func TestReaderError(t *testing.T) {
	path := "/some/sample/path"
	err := NewReaderError(path)

	if err.Data != path {
		t.Errorf("The File was %s, but %s was expected", err.Data, path)
	}

	if strings.Contains(err.Error(), path) == false {
		t.Errorf("The error message of ReaderError does not contain the expected data")
	}
}

func TestWriterError(t *testing.T) {
	path := "/some/sample/path"
	err := NewWriterError(path)

	if err.Data != path {
		t.Errorf("The File was %s, but %s was expected", err.Data, path)
	}

	if strings.Contains(err.Error(), path) == false {
		t.Errorf("The error message of WriterError does not contain the expected data")
	}
}

func TestUnexpectedRequestFormatError(t *testing.T) {
	path := "/some/sample/path"
	err := NewUnexpectedRequestFormatError(path)

	if err.Request != path {
		t.Errorf("The File was %s, but %s was expected", err.Request, path)
	}

	if strings.Contains(err.Error(), path) == false {
		t.Errorf("The error message of WriterError does not contain the expected data")
	}
}
