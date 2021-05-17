package main

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import "fmt"

// SmbStatusTimeOutError - Error when trying to get SMbStatus data runs into a timeout
type SmbStatusTimeOutError struct {
	err string
	// Request - The request that causes this error
	Request string
}

func (e *SmbStatusTimeOutError) Error() string { // Implement the Error Interface for the SmbStatusTimeOutError struct
	return fmt.Sprintf("Error: %s", e.err)
}

// NewSmbStatusTimeOutError- Get a new SmbStatusTimeOutError struct
func NewSmbStatusTimeOutError(request string) *SmbStatusTimeOutError {
	return &SmbStatusTimeOutError{fmt.Sprintf("The \"%s\" timed out", request), request}
}
