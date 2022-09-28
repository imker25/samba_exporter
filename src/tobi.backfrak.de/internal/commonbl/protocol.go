package commonbl

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"fmt"
	"strconv"
	"strings"
)

type RequestType string

// Request the smbd process report table
const PROCESS_REQUEST RequestType = "PROCESS_REQUEST:"

// Request the smbd services report table
const SHARE_REQUEST RequestType = "SHARE_REQUEST:"

// Request the smbd table of locked files
const LOCK_REQUEST RequestType = "LOCK_REQUEST:"

// Request the ps data of the smbd PIDs
const PS_REQUEST RequestType = "PS_REQUEST:"

// Data struct for a psutil response
type PsUtilPidData struct {
	PID                       int64
	CpuUsagePercent           float64
	VirtualMemoryUsageBytes   uint64
	VirtualMemoryUsagePercent float64
	IoCounterReadCount        uint64
	IoCounterReadBytes        uint64
	IoCounterWriteCount       uint64
	IoCounterWriteBytes       uint64
	OpenFilesCount            uint64
	ThreadCount               uint64
}

// GetIdFromRequest - Get the ID from a request telegram
func GetIdFromRequest(request string) (int, error) {
	splitted := strings.Split(request, ":")

	if len(splitted) != 2 {
		return 0, NewUnexpectedRequestFormatError(request)
	}

	idStr := strings.TrimSpace(splitted[1])
	id, errConv := strconv.Atoi(idStr)
	if errConv != nil {
		return 0, NewUnexpectedRequestFormatError(request)
	}

	return id, nil
}

// GetRequest -  Get the request string
func GetRequest(requestType RequestType, id int) string {
	return fmt.Sprintf("%s %d", requestType, id)
}

// GetTestResponseHeader - Get the header for a test response
func GetTestResponseHeader(rType RequestType, id int) string {
	return fmt.Sprintf("%s Test Response for request %d", rType, id)
}

// GetResponseHeader - Get the header for a test response
func GetResponseHeader(rType RequestType, id int) string {
	return fmt.Sprintf("%s Response for request %d", rType, id)
}

// GetResponse - Get the response string out of header and data
func GetResponse(header string, data string) string {
	return fmt.Sprintf("%s\n%s", header, data)
}

// SplitResponse - Split a response string in header and data
func SplitResponse(response string) (string, string, error) {
	splitResponse := strings.SplitN(response, "\n", 2)

	if len(splitResponse) != 2 {
		return "", "", NewUnexpectedResponseFormatError(response)
	}

	header := splitResponse[0]
	data := splitResponse[1]

	return header, data, nil
}

// CheckResponseHeader - Check if a response is for a specific request
func CheckResponseHeader(header string, rType RequestType, id int) bool {
	if !strings.HasPrefix(header, string(rType)) &&
		!strings.Contains(header, fmt.Sprintf("Response for request %d", id)) {
		return false
	}

	return true
}
