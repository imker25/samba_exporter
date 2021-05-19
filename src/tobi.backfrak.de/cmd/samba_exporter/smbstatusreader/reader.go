package smbstatusreader

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"strconv"
	"strings"
	"time"
)

// Type to represent a entry in the 'smbstatus -L' output table
type LockData struct {
	PID        int
	UserID     int
	DenyMode   string
	Access     string
	AccessMode string
	Oplock     string
	SharePath  string
	Name       string
	Time       time.Time
}

// GetLockData - Get the entries out of the 'smbstatus -L' output table multiline string
func GetLockData(data string) []LockData {
	var ret []LockData

	lines := strings.Split(data, "\n")
	sepLineIndex := findSeperatorLineIndex(lines)

	if sepLineIndex < 0 {
		return ret
	}

	tableHeaderMatrix := getFieldMatrix(lines[sepLineIndex-1:sepLineIndex], 9)
	if len(tableHeaderMatrix) != 1 {
		return ret
	}
	tableHeaderFields := tableHeaderMatrix[0]

	if tableHeaderFields[0] != "Pid" && tableHeaderFields[5] != "Oplock" {
		return ret
	}

	for _, fields := range getFieldMatrix(lines[sepLineIndex+1:], 9) {

		var entry LockData
		entry.PID, _ = strconv.Atoi(fields[0])
		entry.UserID, _ = strconv.Atoi(fields[1])
		entry.DenyMode = fields[2]
		entry.Access = fields[3]
		entry.AccessMode = fields[4]
		entry.Oplock = fields[5]
		entry.SharePath = fields[6]
		entry.Name = fields[7]
		entry.Time, _ = time.Parse(time.ANSIC, fields[8])

		ret = append(ret, entry)
	}

	return ret
}

func getFieldMatrix(dataLines []string, lineFields int) [][]string {

	var fieldMatrix [][]string

	for _, line := range dataLines {
		fields := strings.Split(line, "  ")
		var matrixLine []string
		for _, field := range fields {
			trimmedField := strings.TrimSpace(field)
			if trimmedField != "" {
				matrixLine = append(matrixLine, trimmedField)
			}
		}
		if len(matrixLine) == lineFields {
			fieldMatrix = append(fieldMatrix, matrixLine)
		}
	}

	return fieldMatrix
}

func findSeperatorLineIndex(lines []string) int {

	for i, line := range lines {
		if strings.HasPrefix(line, "-----------------------------------------") {
			return i
		}
	}

	return -1
}
