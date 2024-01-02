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
	PID           int
	ClusterNodeId int // In case smaba is running in cluster mode, otherwise -1
	UserID        int
	DenyMode      string
	Access        string
	AccessMode    string
	Oplock        string
	SharePath     string
	Name          string
	Time          time.Time
}

// Implement Stringer Interface for LockData
func (lockData LockData) String() string {
	if lockData.ClusterNodeId > -1 {
		return fmt.Sprintf("ClusterNodeId: %d; PID: %d; UserID: %d; DenyMode: %s; Access: %s; AccessMode: %s; Oplock: %s; SharePath: %s; Name: %s: Time %s;",
			lockData.ClusterNodeId, lockData.PID, lockData.UserID, lockData.DenyMode, lockData.Access, lockData.AccessMode, lockData.Oplock,
			lockData.SharePath, lockData.Name, lockData.Time.Format(time.RFC3339))
	}
	return fmt.Sprintf("PID: %d; UserID: %d; DenyMode: %s; Access: %s; AccessMode: %s; Oplock: %s; SharePath: %s; Name: %s: Time %s;",
		lockData.PID, lockData.UserID, lockData.DenyMode, lockData.Access, lockData.AccessMode, lockData.Oplock,
		lockData.SharePath, lockData.Name, lockData.Time.Format(time.RFC3339))
}

// GetLockData - Get the entries out of the 'smbstatus -L -n' output table multiline string
// Will return an empty array if the data is in unexpected format
func GetLockData(data string, logger *commonbl.Logger) []LockData {
	var ret []LockData
	if strings.HasPrefix(strings.TrimSpace(data), commonbl.NO_LOCKED_FILES) {
		return ret
	}

	if strings.TrimSpace(data) == "" {
		logger.WriteInformation("Got an empty string from 'smbstatus -L -n'")
		return ret
	}

	lines := strings.Split(data, "\n")
	sepLineIndex := findSeperatorLineIndex(lines)

	if sepLineIndex < 1 {
		return ret
	}

	tableHeaderMatrix := getFieldMatrixFixLength(lines[sepLineIndex-1:sepLineIndex], "  ", 9)
	if len(tableHeaderMatrix) != 1 {
		return ret
	}
	tableHeaderFields := tableHeaderMatrix[0]

	if tableHeaderFields[0] != "Pid" || tableHeaderFields[5] != "Oplock" {
		return ret
	}

	i := -1
	for _, oneLineFields := range getFieldMatrix(lines[sepLineIndex+1:], " ") {
		i++
		var err error
		var entry LockData
		fieldLength := len(oneLineFields)
		if strings.Contains(oneLineFields[0], ":") {
			pidFields := strings.Split(oneLineFields[0], ":")
			entry.ClusterNodeId, err = strconv.Atoi(pidFields[0])
			if err != nil {
				logger.WriteErrorWithAddition(err, "while getting LockData ClusterNodeId")
				continue
			}
			entry.PID, err = strconv.Atoi(pidFields[1])
			if err != nil {
				logger.WriteErrorWithAddition(err, "while getting LockData PID (ClusterNodeId)")
				continue
			}
		} else {
			entry.ClusterNodeId = -1
			entry.PID, err = strconv.Atoi(oneLineFields[0])
			if err != nil {
				logger.WriteErrorWithAddition(err, "while getting LockData PID")
				continue
			}
		}
		entry.UserID, err = strconv.Atoi(oneLineFields[1])
		if err != nil {
			logger.WriteErrorWithAddition(err, "while getting LockData UserID")
			continue
		}
		entry.DenyMode = oneLineFields[2]
		entry.Access = oneLineFields[3]
		entry.AccessMode = oneLineFields[4]
		entry.Oplock = oneLineFields[5]
		entry.SharePath = oneLineFields[6]
		timeConvSuc := false
		var connectTime time.Time
		var lastNameIndex = -1
		timeConvSuc, connectTime = tryGetTimeStampFromStrArr(oneLineFields[fieldLength-5 : fieldLength])
		if timeConvSuc {
			entry.Time = connectTime
			lastNameIndex = fieldLength - 5
		} else {
			timeConvSuc, connectTime = tryGetTimeStampFromStrArr(oneLineFields[fieldLength-6 : fieldLength])
			if timeConvSuc {
				entry.Time = connectTime
				lastNameIndex = fieldLength - 6
			}
		}

		if lastNameIndex == -1 {
			logger.WriteErrorMessage(fmt.Sprintf("Not able to parse the time stamp in following LockData line: \"%s\"", lines[sepLineIndex+1+i]))
			continue
		}

		if lastNameIndex <= 7 {
			logger.WriteErrorMessage(fmt.Sprintf("Not able to find the name in following LockData line: \"%s\"", lines[sepLineIndex+1+i]))
			continue
		}

		name := ""
		for _, namePart := range oneLineFields[7:lastNameIndex] {
			name = fmt.Sprintf("%s %s", name, namePart)
		}
		entry.Name = strings.TrimSpace(name)

		ret = append(ret, entry)
	}
	return ret
}

// Type to represent a entry in the 'smbstatus -S -n' output table
type ShareData struct {
	Service       string
	PID           int
	ClusterNodeId int // In case smaba is running in cluster mode, otherwise -1
	Machine       string
	ConnectedAt   time.Time
	Encryption    string
	Signing       string
}

// Implement Stringer Interface for ShareData
func (shareData ShareData) String() string {
	if shareData.ClusterNodeId > -1 {
		return fmt.Sprintf("Service: %s; ClusterNodeId: %d; PID: %d; Machine: %s; ConnectedAt: %s; Encryption: %s; Signing: %s;",
			shareData.Service, shareData.ClusterNodeId, shareData.PID, shareData.Machine, shareData.ConnectedAt.Format(time.RFC3339),
			shareData.Encryption, shareData.Signing)
	}
	return fmt.Sprintf("Service: %s; PID: %d; Machine: %s; ConnectedAt: %s; Encryption: %s; Signing: %s;",
		shareData.Service, shareData.PID, shareData.Machine, shareData.ConnectedAt.Format(time.RFC3339),
		shareData.Encryption, shareData.Signing)
}

// GetShareData - Get the entries out of the 'smbstatus -S -n' output table multiline string
// Will return an empty array if the data is in unexpected format
func GetShareData(data string, logger *commonbl.Logger) []ShareData {
	var ret []ShareData

	if strings.TrimSpace(data) == "" {
		logger.WriteInformation("Got an empty string from 'smbstatus -S -n'")
		return ret
	}

	lines := strings.Split(data, "\n")
	sepLineIndex := findSeperatorLineIndex(lines)

	if sepLineIndex < 1 {
		return ret
	}

	// Normal setup gives 6 fields in this line
	tableHeaderMatrix := getFieldMatrixFixLength(lines[sepLineIndex-1:sepLineIndex], "  ", 6)

	if len(tableHeaderMatrix) != 1 {
		// Cluster setup gives 7 fields in this line
		tableHeaderMatrix = getFieldMatrixFixLength(lines[sepLineIndex-1:sepLineIndex], "  ", 7)

		if len(tableHeaderMatrix) != 1 {
			return ret
		}
	}
	tableHeaderFields := tableHeaderMatrix[0]
	runningMode := "none"
	if tableHeaderFields[0] == "Service" && tableHeaderFields[3] == "Connected at" {
		runningMode = "normal"
	}

	if tableHeaderFields[0] == "PID" && tableHeaderFields[4] == "Protocol Version" {
		runningMode = "cluster"
	}

	if runningMode == "normal" {
		i := -1
		for _, oneLineFields := range getFieldMatrix(lines[sepLineIndex+1:], " ") {
			i++
			lastNameField := -1
			var err error
			var entry ShareData
			fieldLength := len(oneLineFields)
			if strings.Contains(oneLineFields[1], ":") {
				pidFields := strings.Split(oneLineFields[1], ":")
				entry.ClusterNodeId, err = strconv.Atoi(pidFields[0])
				if err != nil {
					logger.WriteErrorWithAddition(err, "while getting ShareData ClusterNodeId (normal with :)")
					continue
				}
				entry.PID, err = strconv.Atoi(pidFields[1])
				if err != nil {
					logger.WriteErrorWithAddition(err, "while getting ShareData PID (normal with :)")
					continue
				}
			} else {

				entry.ClusterNodeId = -1

				pidFound := true
				for {
					lastNameField++
					entry.PID, err = strconv.Atoi(oneLineFields[lastNameField+1])
					if err == nil {
						break
					}
					if len(oneLineFields)-11 <= lastNameField {
						logger.WriteErrorWithAddition(err, "while getting ShareData PID (normal without :)")
						pidFound = false
						break
					}
				}

				if !pidFound {
					continue
				}
				entry.Service = concatStrFromArr(oneLineFields[0 : lastNameField+1])
			}
			entry.Machine = oneLineFields[lastNameField+2]
			timeConvSuc := false
			var connectTime time.Time
			var lastTimeIndex = -1
			timeConvSuc, connectTime = tryGetTimeStampFromStrArr(oneLineFields[lastNameField+3 : lastNameField+10])
			if timeConvSuc {
				entry.ConnectedAt = connectTime
				lastTimeIndex = lastNameField + 9
			} else {
				timeConvSuc, connectTime = tryGetTimeStampFromStrArr(oneLineFields[lastNameField+3 : lastNameField+9])
				if timeConvSuc {
					entry.ConnectedAt = connectTime
					lastTimeIndex = lastNameField + 8
				}
			}

			if lastTimeIndex == -1 {
				logger.WriteErrorMessage(fmt.Sprintf("Not able to parse the time stamp in following LockData line: \"%s\"", lines[sepLineIndex+1+i]))
				continue
			}
			if lastTimeIndex != fieldLength-3 {
				logger.WriteErrorMessage(fmt.Sprintf("Can not find end of time stamp in following LockData line: \"%s\"", lines[sepLineIndex+1+i]))
				continue
			}
			entry.Encryption = oneLineFields[lastTimeIndex+1]
			entry.Signing = oneLineFields[lastTimeIndex+2]

			ret = append(ret, entry)
		}

	} else if runningMode == "cluster" {
		i := -1
		for _, oneLineFields := range getFieldMatrix(lines[sepLineIndex+1:], " ") {
			i++
			var err error
			var entry ShareData
			fieldLength := len(oneLineFields)
			if strings.Contains(oneLineFields[0], ":") {
				pidFields := strings.Split(oneLineFields[0], ":")
				entry.ClusterNodeId, err = strconv.Atoi(pidFields[0])
				if err != nil {
					logger.WriteErrorWithAddition(err, "while getting ShareData ClusterNodeId (cluster - with :)")
					continue
				}
				entry.PID, err = strconv.Atoi(pidFields[1])
				if err != nil {
					logger.WriteErrorWithAddition(err, "while getting ShareData PID (cluster - with :)")
					continue
				}
			} else {
				entry.ClusterNodeId = -1
				entry.PID, err = strconv.Atoi(oneLineFields[0])
				if err != nil {
					logger.WriteErrorWithAddition(err, "while getting ShareData PID (cluster - without :)")
					continue
				}
			}
			if fieldLength == 8 {
				entry.Machine = fmt.Sprintf("%s %s", oneLineFields[3], oneLineFields[4])
				entry.Encryption = oneLineFields[6]
				entry.Signing = oneLineFields[7]

			} else if fieldLength == 7 {
				entry.Machine = oneLineFields[3]
				entry.Encryption = oneLineFields[5]
				entry.Signing = oneLineFields[6]
			} else {
				logger.WriteErrorMessage(fmt.Sprintf("Can not parse the following ShareData line: \"%s\"", lines[i]))
				continue
			}

			ret = append(ret, entry)
		}
	}

	return ret
}

// Type to represent a entry in the 'smbstatus -p -n' output table
type ProcessData struct {
	PID             int
	ClusterNodeId   int // In case smaba is running in cluster mode, otherwise -1
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
	if processData.ClusterNodeId > -1 {
		return fmt.Sprintf("ClusterNodeId: %d; PID: %d; UserID: %d; GroupID: %d; Machine: %s; ProtocolVersion: %s; Encryption: %s; Signing: %s;",
			processData.ClusterNodeId, processData.PID, processData.UserID, processData.GroupID, processData.Machine, processData.ProtocolVersion,
			processData.Encryption, processData.Signing)
	}
	return fmt.Sprintf("PID: %d; UserID: %d; GroupID: %d; Machine: %s; ProtocolVersion: %s; Encryption: %s; Signing: %s;",
		processData.PID, processData.UserID, processData.GroupID, processData.Machine, processData.ProtocolVersion,
		processData.Encryption, processData.Signing)
}

// GetProcessData - Get the entries out of the 'smbstatus -p -n' output table multiline string
// Will return an empty array if the data is in unexpected format
func GetProcessData(data string, logger *commonbl.Logger) []ProcessData {
	var ret []ProcessData

	if strings.TrimSpace(data) == "" {
		logger.WriteInformation("Got an empty string from 'smbstatus -p -n'")
		return ret
	}

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

	tableHeaderMatrix := getFieldMatrixFixLength(lines[sepLineIndex-1:sepLineIndex], "  ", 7)
	if len(tableHeaderMatrix) != 1 {
		return ret
	}
	tableHeaderFields := tableHeaderMatrix[0]

	if tableHeaderFields[1] != "Username" || tableHeaderFields[4] != "Protocol Version" {
		return ret
	}

	i := -1
	for _, oneLineFields := range getFieldMatrix(lines[sepLineIndex+1:], " ") {
		i++
		var err error
		var entry ProcessData
		fieldLength := len(oneLineFields)
		// In cluster versions samba adds an extra id separated by ':'
		if strings.Contains(oneLineFields[0], ":") {
			pidFields := strings.Split(oneLineFields[0], ":")
			entry.ClusterNodeId, err = strconv.Atoi(pidFields[0])
			if err != nil {
				logger.WriteErrorWithAddition(err, "while getting ProcessData ClusterNodeId")
				continue
			}
			entry.PID, err = strconv.Atoi(pidFields[1])
			if err != nil {
				logger.WriteErrorWithAddition(err, "while getting ProcessData PID (with :)")
				continue
			}
		} else {
			entry.ClusterNodeId = -1
			entry.PID, err = strconv.Atoi(oneLineFields[0])
			if err != nil {
				logger.WriteErrorWithAddition(err, "while getting ProcessData PID (without :)")
				continue
			}
		}
		// In cluster versions samba does not print the users id, but nobody
		if oneLineFields[1] == "nobody" {
			entry.UserID = -1
		} else {
			entry.UserID, err = strconv.Atoi(oneLineFields[1])
			if err != nil {
				logger.WriteErrorWithAddition(err, "while getting ProcessData UserID")
				continue
			}
		}
		// In cluster versions samba does not print the group id, but nogroup
		if oneLineFields[2] == "nogroup" {
			entry.GroupID = -1
		} else {
			entry.GroupID, err = strconv.Atoi(oneLineFields[2])
			if err != nil {
				logger.WriteErrorWithAddition(err, "while getting ProcessData GroupID")
				continue
			}
		}
		if fieldLength == 8 {
			entry.Machine = fmt.Sprintf("%s %s", oneLineFields[3], oneLineFields[4])
			entry.ProtocolVersion = oneLineFields[5]
			entry.Encryption = oneLineFields[6]
			entry.Signing = oneLineFields[7]
		} else if fieldLength == 7 {
			entry.Machine = oneLineFields[3]
			entry.ProtocolVersion = oneLineFields[4]
			entry.Encryption = oneLineFields[5]
			entry.Signing = oneLineFields[6]
		} else {
			logger.WriteErrorMessage(fmt.Sprintf("Can not parse the following ProcessData line: \"%s\"", lines[i]))
			continue
		}
		entry.SambaVersion = sambaVersion

		ret = append(ret, entry)
	}
	return ret
}

func GetPsData(data string, logger *commonbl.Logger) []commonbl.PsUtilPidData {
	var ret []commonbl.PsUtilPidData
	errConv := json.Unmarshal([]byte(data), &ret)
	if errConv != nil {
		logger.WriteErrorWithAddition(errConv, "while converting PsData json")
		return []commonbl.PsUtilPidData{}
	}

	return ret
}

func getFieldMatrixFixLength(dataLines []string, separator string, lineFields int) [][]string {

	var fieldMatrix [][]string

	for _, matrixLine := range getFieldMatrix(dataLines, separator) {
		if len(matrixLine) == lineFields {
			fieldMatrix = append(fieldMatrix, matrixLine)
		}
	}

	return fieldMatrix
}

func getFieldMatrix(dataLines []string, separator string) [][]string {

	var fieldMatrix [][]string

	for _, line := range dataLines {
		fields := strings.Split(line, separator)
		var matrixLine []string
		for _, field := range fields {
			trimmedField := strings.TrimSpace(field)
			if trimmedField != "" {
				matrixLine = append(matrixLine, trimmedField)
			}
		}
		fieldMatrix = append(fieldMatrix, matrixLine)
	}

	return fieldMatrix
}

func concatStrFromArr(fields []string) string {
	ret := ""
	for i, field := range fields {
		if i == 0 {
			ret = field
		} else {
			ret = ret + " " + field
		}
	}

	return ret
}

func tryGetTimeStampFromStrArr(fields []string) (bool, time.Time) {
	timeStr := ""
	var ret time.Time
	var err error
	for _, sec := range fields {
		timeStr = fmt.Sprintf("%s %s", timeStr, sec)
	}
	timeStr = strings.TrimSpace(timeStr)
	ret, err = time.ParseInLocation(time.ANSIC, timeStr, time.Now().Location())
	if err == nil {
		return true, ret
	}
	ret, err = time.Parse(time.ANSIC, timeStr)
	if err == nil {
		return true, ret
	}
	ret, err = time.Parse("Mon Jan 02 03:04:05 PM 2006 MST", timeStr)
	if err == nil {
		return true, ret
	}
	ret, err = time.Parse("Mon Jan 2 03:04:05 PM 2006 MST", timeStr)
	if err == nil {
		return true, ret
	}
	ret, err = time.Parse("Mon Jan _2 15:04:05 2006 MST", timeStr)
	if err == nil {
		return true, ret
	}
	ret, err = time.Parse("Mo Jan _2 15:04:05 2006 MST", timeStr)
	if err == nil {
		return true, ret
	}

	return false, time.Now()
}

func findSeperatorLineIndex(lines []string) int {

	for i, line := range lines {
		if strings.HasPrefix(line, "-----------------------------------------") {
			return i
		}
	}

	return -1
}
