package smbstatusreader

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"tobi.backfrak.de/internal/commonbl"
)

// Type to represent a entry in the 'smbstatus -L -n' output table
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

// Implement Stringer Interface for LockData
func (lockData LockData) String() string {
	return fmt.Sprintf("PID: %d; UserID: %d; DenyMode: %s; Access: %s; AccessMode: %s; Oplock: %s; SharePath: %s; Name: %s: Time %s;",
		lockData.PID, lockData.UserID, lockData.DenyMode, lockData.Access, lockData.AccessMode, lockData.Oplock,
		lockData.SharePath, lockData.Name, lockData.Time.Format(time.RFC3339))
}

// GetLockData - Get the entries out of the 'smbstatus -L -n' output table multiline string
// Will return an empty array if the data is in unexpected format
func GetLockData(data string, logger commonbl.Logger) []LockData {
	var ret []LockData
	if strings.TrimSpace(data) == "No locked files" {
		return ret
	}

	lines := strings.Split(data, "\n")
	sepLineIndex := findSeperatorLineIndex(lines)

	if sepLineIndex < 1 {
		return ret
	}

	tableHeaderMatrix := getFieldMatrix(lines[sepLineIndex-1:sepLineIndex], "  ", 9)
	if len(tableHeaderMatrix) != 1 {
		return ret
	}
	tableHeaderFields := tableHeaderMatrix[0]

	if tableHeaderFields[0] != "Pid" || tableHeaderFields[5] != "Oplock" {
		return ret
	}

	for _, fields := range getFieldMatrix(lines[sepLineIndex+1:], " ", 13) {
		var err error
		var entry LockData
		entry.PID, err = strconv.Atoi(fields[0])
		if err != nil {
			logger.WriteError(err)
			continue
		}
		entry.UserID, err = strconv.Atoi(fields[1])
		if err != nil {
			logger.WriteError(err)
			continue
		}
		entry.DenyMode = fields[2]
		entry.Access = fields[3]
		entry.AccessMode = fields[4]
		entry.Oplock = fields[5]
		entry.SharePath = fields[6]
		entry.Name = fields[7]
		entry.Time, err = time.ParseInLocation(time.ANSIC,
			fmt.Sprintf("%s %s %s %s %s", fields[8], fields[9], fields[10], fields[11], fields[12]),
			time.Now().Location())
		if err != nil {
			logger.WriteError(err)
			continue
		}

		ret = append(ret, entry)
	}

	return ret
}

// Type to represent a entry in the 'smbstatus -S -n' output table
type ShareData struct {
	Service     string
	PID         int
	Machine     string
	ConnectedAt time.Time
	Encryption  string
	Signing     string
}

// Implement Stringer Interface for ShareData
func (shareData ShareData) String() string {
	return fmt.Sprintf("Service: %s; PID: %d; Machine: %s; ConnectedAt: %s; Encryption: %s; Signing: %s;",
		shareData.Service, shareData.PID, shareData.Machine, shareData.ConnectedAt.Format(time.RFC3339),
		shareData.Encryption, shareData.Signing)
}

// GetShareData - Get the entries out of the 'smbstatus -S -n' output table multiline string
// Will return an empty array if the data is in unexpected format
func GetShareData(data string, logger commonbl.Logger) []ShareData {
	var ret []ShareData
	lines := strings.Split(data, "\n")
	sepLineIndex := findSeperatorLineIndex(lines)

	if sepLineIndex < 1 {
		return ret
	}

	tableHeaderMatrix := getFieldMatrix(lines[sepLineIndex-1:sepLineIndex], "  ", 6)
	if len(tableHeaderMatrix) != 1 {
		return ret
	}
	tableHeaderFields := tableHeaderMatrix[0]

	if tableHeaderFields[0] != "Service" || tableHeaderFields[3] != "Connected at" {
		return ret
	}
	fieldMatrix := getFieldMatrix(lines[sepLineIndex+1:], " ", 12)
	if fieldMatrix != nil {
		for _, fields := range fieldMatrix {
			var err error
			var entry ShareData
			entry.Service = fields[0]
			entry.PID, err = strconv.Atoi(fields[1])
			if err != nil {
				logger.WriteError(err)
				continue
			}
			entry.Machine = fields[2]
			timeStr := fmt.Sprintf("%s %s %s %s %s %s %s", fields[3], fields[4], fields[5], fields[6], fields[7], fields[8], fields[9])
			entry.ConnectedAt, err = time.Parse("Mon Jan 02 03:04:05 PM 2006 MST", timeStr)
			if err != nil {
				entry.ConnectedAt, err = time.Parse("Mon Jan 2 03:04:05 PM 2006 MST", timeStr)
				if err != nil {
					logger.WriteError(err)
					continue
				}
			}
			entry.Encryption = fields[10]
			entry.Signing = fields[11]

			ret = append(ret, entry)
		}
	} else {
		fieldMatrix = getFieldMatrix(lines[sepLineIndex+1:], " ", 11)
		if fieldMatrix != nil {
			for _, fields := range fieldMatrix {
				var err error
				var entry ShareData
				entry.Service = fields[0]
				entry.PID, err = strconv.Atoi(fields[1])
				if err != nil {
					logger.WriteError(err)
					continue
				}
				entry.Machine = fields[2]
				timeStr := fmt.Sprintf("%s %s %s %s %s %s", fields[3], fields[4], fields[5], fields[6], fields[7], fields[8])
				entry.ConnectedAt, err = time.Parse("Mon Jan _2 15:04:05 2006 MST", timeStr)
				if err != nil {
					entry.ConnectedAt, err = time.Parse("Mo Jan _2 15:04:05 2006 MST", timeStr)
					if err != nil {
						logger.WriteError(err)
						continue
					}
				}
				entry.Encryption = fields[9]
				entry.Signing = fields[10]

				ret = append(ret, entry)
			}
		}
	}
	return ret
}

// Type to represent a entry in the 'smbstatus -p -n' output table
type ProcessData struct {
	PID             int
	UserID          int
	GroupID         int
	Machine         string
	ProtocolVersion string
	Encryption      string
	Signing         string
	SambaVersion    string
}

// Implement Stringer Interface for ProcessData
func (processData ProcessData) String() string {
	return fmt.Sprintf("PID: %d; UserID: %d; GroupID: %d; Machine: %s; ProtocolVersion: %s; Encryption: %s; Signing: %s;",
		processData.PID, processData.UserID, processData.GroupID, processData.Machine, processData.ProtocolVersion,
		processData.Encryption, processData.Signing)
}

// GetProcessData - Get the entries out of the 'smbstatus -p -n' output table multiline string
// Will return an empty array if the data is in unexpected format
func GetProcessData(data string, logger commonbl.Logger) []ProcessData {
	var ret []ProcessData
	lines := strings.Split(data, "\n")
	sepLineIndex := findSeperatorLineIndex(lines)

	if sepLineIndex < 2 {
		return ret
	}

	var sambaVersion string
	sambaVersionLine := lines[sepLineIndex-2 : sepLineIndex-1][0]
	if strings.HasPrefix(sambaVersionLine, "Samba version") {
		sambaVersion = strings.TrimSpace(strings.Replace(sambaVersionLine, "Samba version", "", 1))
	} else {
		return ret
	}

	tableHeaderMatrix := getFieldMatrix(lines[sepLineIndex-1:sepLineIndex], "  ", 7)
	if len(tableHeaderMatrix) != 1 {
		return ret
	}
	tableHeaderFields := tableHeaderMatrix[0]

	if tableHeaderFields[1] != "Username" || tableHeaderFields[4] != "Protocol Version" {
		return ret
	}

	for _, fields := range getFieldMatrix(lines[sepLineIndex+1:], " ", 8) {
		var err error
		var entry ProcessData
		entry.PID, err = strconv.Atoi(fields[0])
		if err != nil {
			logger.WriteError(err)
			continue
		}
		entry.UserID, err = strconv.Atoi(fields[1])
		if err != nil {
			logger.WriteError(err)
			continue
		}
		entry.GroupID, err = strconv.Atoi(fields[2])
		if err != nil {
			logger.WriteError(err)
			continue
		}
		entry.Machine = fmt.Sprintf("%s %s", fields[3], fields[4])
		entry.ProtocolVersion = fields[5]
		entry.Encryption = fields[6]
		entry.Signing = fields[7]
		entry.SambaVersion = sambaVersion

		ret = append(ret, entry)
	}
	return ret
}

func GetPsData(data string, logger commonbl.Logger) []commonbl.PsUtilPidData {
	var ret []commonbl.PsUtilPidData
	errConv := json.Unmarshal([]byte(data), &ret)
	if errConv != nil {
		logger.WriteError(errConv)
		return []commonbl.PsUtilPidData{}
	}

	return ret
}

func getFieldMatrix(dataLines []string, seperator string, lineFields int) [][]string {

	var fieldMatrix [][]string

	for _, line := range dataLines {
		fields := strings.Split(line, seperator)
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
