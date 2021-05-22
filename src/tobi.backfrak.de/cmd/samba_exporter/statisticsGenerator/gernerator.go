package statisticsGenerator

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import "tobi.backfrak.de/cmd/samba_exporter/smbstatusreader"

// Type for numeric statistic values from the samba server
type SmbStatisticsNumeric struct {
	Name  string
	Value int
	Help  string
}

// GetSmbStatistics - Get the statistic data for prometheus out of the response data arrays
func GetSmbStatistics(lockData []smbstatusreader.LockData, processData []smbstatusreader.ProcessData, shareData []smbstatusreader.ShareData) []SmbStatisticsNumeric {
	var ret []SmbStatisticsNumeric

	var users []int
	var pids []int
	var shares []string
	var clients []string

	for _, lock := range lockData {
		if !intArrContains(users, lock.UserID) {
			users = append(users, lock.UserID)
		}

		if !intArrContains(pids, lock.PID) {
			pids = append(pids, lock.PID)
		}
	}

	for _, process := range processData {
		if !intArrContains(users, process.UserID) {
			users = append(users, process.UserID)
		}

		if !intArrContains(pids, process.PID) {
			pids = append(pids, process.PID)
		}
	}

	for _, share := range shareData {
		if !intArrContains(pids, share.PID) {
			pids = append(pids, share.PID)
		}

		if !strArrContains(shares, share.Service) {
			shares = append(shares, share.Service)
		}

		if !strArrContains(clients, share.Machine) {
			clients = append(clients, share.Machine)
		}
	}

	ret = append(ret, SmbStatisticsNumeric{"individual_user_count", len(users), "The number of users connected to this samba server"})
	ret = append(ret, SmbStatisticsNumeric{"locked_file_count", len(lockData), "Number of files locked by the samba server"})
	ret = append(ret, SmbStatisticsNumeric{"pid_count", len(pids), "Number of processes running by the samba server"})
	ret = append(ret, SmbStatisticsNumeric{"share_count", len(shares), "Number of shares used by clients of the samba server"})
	ret = append(ret, SmbStatisticsNumeric{"client_count", len(clients), "Number of clients using the samba server"})

	return ret
}

func intArrContains(arr []int, value int) bool {
	for _, field := range arr {
		if field == value {
			return true
		}
	}

	return false
}

func strArrContains(arr []string, value string) bool {
	for _, field := range arr {
		if field == value {
			return true
		}
	}

	return false
}
