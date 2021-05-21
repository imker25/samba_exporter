package main

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"fmt"

	"tobi.backfrak.de/internal/commonbl"
)

// SmbStatusTimeOutError - Error when trying to get SMbStatus data runs into a timeout
type SmbStatusTimeOutError struct {
	err string
	// Request - The request that causes this error
	Request commonbl.RequestType
}

func (e *SmbStatusTimeOutError) Error() string { // Implement the Error Interface for the SmbStatusTimeOutError struct
	return fmt.Sprintf("Error: %s", e.err)
}

// NewSmbStatusTimeOutError- Get a new SmbStatusTimeOutError struct
func NewSmbStatusTimeOutError(request commonbl.RequestType) *SmbStatusTimeOutError {
	return &SmbStatusTimeOutError{fmt.Sprintf("The \"%s\" timed out", request), request}
}

// SmbStatusUnexpectedResponseError - Error when the SMbStatus data is unexpcted
type SmbStatusUnexpectedResponseError struct {
	err string
	// Response - The Response that causes this error
	Response string
}

func (e *SmbStatusUnexpectedResponseError) Error() string { // Implement the Error Interface for the SmbStatusUnexpectedResponseError struct
	return fmt.Sprintf("Error: %s", e.err)
}

// NewSmbStatusUnexpectedResponseError - Get a new SmbStatusUnexpectedResponseError struct
func NewSmbStatusUnexpectedResponseError(response string) *SmbStatusUnexpectedResponseError {
	return &SmbStatusUnexpectedResponseError{fmt.Sprintf("The response \"%s\" was not exptected", response), response}
}
