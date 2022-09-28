package smbstatusdbl

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v3/process"

	"tobi.backfrak.de/internal/commonbl"
)

// Class to get commonbl.PsUtilPidData for a Process
type PsDataGenerator struct {
	ProcessToRequest string
	pgrepPath        string
}

// Get a new instance of PsDataGenerator
func NewPsDataGenerator(processToRequest string) (*PsDataGenerator, error) {
	ret := PsDataGenerator{}
	var errLookPath error
	ret.ProcessToRequest = processToRequest
	ret.pgrepPath, errLookPath = exec.LookPath("pgrep")
	if errLookPath != nil {
		return nil, errLookPath
	}

	return &ret, nil
}

// Get the commonbl.PsUtilPidData data of the ProcessToRequest.
// - In case this process is not running an empty list is returned
// - In case an error, other then not finding the process, occurs during gathering data it is returned
func (generator *PsDataGenerator) GetPsUtilPidData() ([]commonbl.PsUtilPidData, error) {
	ret := []commonbl.PsUtilPidData{}

	pidList, pidListErr := generator.getPidList()
	if (pidListErr != nil) || (len(pidList) == 0) {
		return ret, nil
	}

	for _, pid := range pidList {

		proc, errProc := process.NewProcess(pid)
		if errProc != nil {
			return nil, errProc
		}

		cpuPercent, errPer := proc.CPUPercent()
		if errPer != nil {
			return nil, errPer
		}
		vmBytes, errVmBytes := proc.MemoryInfo()
		if errVmBytes != nil {
			return nil, errVmBytes
		}
		vmPercent, errVmPercent := proc.MemoryPercent()
		if errVmPercent != nil {
			return nil, errVmPercent
		}
		ioCounters, errIoCounters := proc.IOCounters()
		if errIoCounters != nil {
			return nil, errIoCounters
		}
		openFileStats, errOpenFileStats := proc.OpenFiles()
		if errOpenFileStats != nil {
			return nil, errOpenFileStats
		}
		threadStats, errThreadStats := proc.Threads()
		if errThreadStats != nil {
			return nil, errThreadStats
		}

		entry := commonbl.PsUtilPidData{
			int64(pid),
			cpuPercent,
			vmBytes.VMS,
			float64(vmPercent),
			ioCounters.ReadCount,
			ioCounters.ReadBytes,
			ioCounters.WriteCount,
			ioCounters.WriteBytes,
			uint64(len(openFileStats)),
			uint64(len(threadStats)),
		}

		ret = append(ret, entry)
	}

	return ret, nil
}

func (generator *PsDataGenerator) getPidList() ([]int32, error) {
	var pidList []int32

	pidListInByte, errGetPids := exec.Command(generator.pgrepPath, generator.ProcessToRequest).Output()
	if errGetPids != nil {
		return nil, errGetPids
	}

	pidListLines := strings.Split(string(pidListInByte), "\n")
	for _, line := range pidListLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		pid, errConv := strconv.Atoi(line)
		if errConv != nil {
			return nil, errConv
		} else {
			pidList = append(pidList, int32(pid))
		}
	}

	return pidList, nil
}
