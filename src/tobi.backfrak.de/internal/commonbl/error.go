package commonbl

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import "fmt"

// ReaderError - Error when trying to read from the buffer
type ReaderError struct {
	err string
	// File - The path to the dir that caused this error
	Data string
}

func (e *ReaderError) Error() string { // Implement the Error Interface for the ReaderError struct
	return fmt.Sprintf("Error: %s", e.err)
}

// NewReaderError- Get a new OutFileIsDirError struct
func NewReaderError(data string) *ReaderError {
	return &ReaderError{fmt.Sprintf("The received data \"%s\" was not expected", data), data}
}

// WriterError - Error when trying to write to the buffer
type WriterError struct {
	err string
	// File - The path to the dir that caused this error
	Data string
}

func (e *WriterError) Error() string { // Implement the Error Interface for the WriterError struct
	return fmt.Sprintf("Error: %s", e.err)
}

// NewWriterError- Get a new OutFileIsDirError struct
func NewWriterError(data string) *WriterError {
	return &WriterError{fmt.Sprintf("The data \"%s\" was not written", data), data}
}
