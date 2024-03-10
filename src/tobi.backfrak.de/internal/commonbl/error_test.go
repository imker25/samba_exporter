package commonbl

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"fmt"
	"strings"
	"testing"
)

func TestReaderError(t *testing.T) {
	path := "/some/sample/path"
	id := 123
	rType := LOCK_REQUEST
	err := NewReaderError(path, rType, id)

	if err.Data != path {
		t.Errorf("The Data was %s, but %s was expected", err.Data, path)
	}

	if string(err.Request) != string(LOCK_REQUEST) {
		t.Errorf("The Request was %s, but %s was expected", err.Request, path)
	}

	if err.ID != id {
		t.Errorf("The ID was %d, but %d was expected", err.ID, id)
	}

	if strings.Contains(err.Error(), path) == false {
		t.Errorf("The error message of ReaderError does not contain the expected data")
	}

	if strings.Contains(err.Error(), string(rType)) == false {
		t.Errorf("The error message of ReaderError does not contain the expected request type")
	}

	if strings.Contains(err.Error(), fmt.Sprintf("%d", id)) == false {
		t.Errorf("The error message of ReaderError does not contain the expected request id")
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

func TestDirectoryNotExistError(t *testing.T) {
	path := "/some/sample/path"
	err := NewDirectoryNotExistError(path)

	if err.DirectoryPath != path {
		t.Errorf("DirectoryNotExistError DirectoryPath path value is '%s', but '%s' is expected", err.DirectoryPath, path)
	}

	if strings.Contains(err.Error(), path) == false {
		t.Errorf("The error message of DirectoryNotExistError does not contain the expected data")
	}
}
