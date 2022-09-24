package statisticsGenerator

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"strconv"
	"time"

	"tobi.backfrak.de/internal/smbexporterbl/smbstatusreader"
)

// Type for numeric statistic values from the samba server
type SmbStatisticsNumeric struct {
	Name   string
	Value  float64
	Help   string
	Labels map[string]string
}

type StatisticsGeneratorSettings struct {
	DoNotExportClient     bool
	DoNotExportUser       bool
	DoNotExportEncryption bool
}

type lockCreationEntry struct {
	UserID       int
	CreationTime time.Time
	Share        string
}

// GetSmbStatistics - Get the statistic data for prometheus out of the response data arrays
func GetSmbStatistics(lockData []smbstatusreader.LockData, processData []smbstatusreader.ProcessData, shareData []smbstatusreader.ShareData, settings StatisticsGeneratorSettings) []SmbStatisticsNumeric {
	var ret []SmbStatisticsNumeric

	var users []int
	var pids []int
	var shares []string
	var clients []string
	var sambaVersion string
	var lockCreationEntries []lockCreationEntry
	locksPerShare := make(map[string]int, 0)
	processPerClient := make(map[string]int, 0)
	protocolVersionCount := make(map[string]int, 0)
	signingMethodCount := make(map[string]int, 0)
	encryptionMethodCount := make(map[string]int, 0)
	clientConnectionTime := make(map[string]int64, 0)

	for _, lock := range lockData {
		if !intArrContains(users, lock.UserID) {
			users = append(users, lock.UserID)
		}

		if !intArrContains(pids, lock.PID) {
			pids = append(pids, lock.PID)
		}

		locksOfShare, found := locksPerShare[lock.SharePath]
		if !found {
			locksPerShare[lock.SharePath] = 1
		} else {
			locksPerShare[lock.SharePath] = locksOfShare + 1
		}

		newEntry := lockCreationEntry{lock.UserID, lock.Time, lock.SharePath}
		if !lockArrContainsEntry(lockCreationEntries, newEntry) {
			lockCreationEntries = append(lockCreationEntries, newEntry)
		}
	}

	for _, process := range processData {
		if !intArrContains(users, process.UserID) {
			users = append(users, process.UserID)
		}

		if !intArrContains(pids, process.PID) {
			pids = append(pids, process.PID)
		}
		sambaVersion = process.SambaVersion

		processOnShare, foundC := processPerClient[process.Machine]
		if !foundC {
			processPerClient[process.Machine] = 1
		} else {
			processPerClient[process.Machine] = processOnShare + 1
		}

		versionCount, foundV := protocolVersionCount[process.ProtocolVersion]
		if !foundV {
			protocolVersionCount[process.ProtocolVersion] = 1
		} else {
			protocolVersionCount[process.ProtocolVersion] = versionCount + 1
		}

		signingCount, foundS := signingMethodCount[process.Signing]
		if !foundS {
			signingMethodCount[process.Signing] = 1
		} else {
			signingMethodCount[process.Signing] = signingCount + 1
		}

		encryptionCount, foundE := encryptionMethodCount[process.Encryption]
		if !foundE {
			encryptionMethodCount[process.Encryption] = 1
		} else {
			encryptionMethodCount[process.Encryption] = encryptionCount + 1
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

		_, foundC := clientConnectionTime[share.Machine]
		if !foundC {
			clientConnectionTime[share.Machine] = share.ConnectedAt.Unix()
		}
	}
	// TODO: Generate more metrics
	ret = append(ret, SmbStatisticsNumeric{"individual_user_count", float64(len(users)), "The number of users connected to this samba server", nil})
	ret = append(ret, SmbStatisticsNumeric{"locked_file_count", float64(len(lockData)), "Number of files locked by the samba server", nil})
	ret = append(ret, SmbStatisticsNumeric{"pid_count", float64(len(pids)), "Number of processes running by the samba server", nil})
	ret = append(ret, SmbStatisticsNumeric{"share_count", float64(len(shares)), "Number of shares servered by the samba server", nil})
	ret = append(ret, SmbStatisticsNumeric{"client_count", float64(len(clients)), "Number of clients using the samba server", nil})

	if len(locksPerShare) > 0 {
		for share, locks := range locksPerShare {
			ret = append(ret, SmbStatisticsNumeric{"locks_per_share_count", float64(locks), "Number of locks on share", map[string]string{"share": share}})
		}
	} else {
		// Add this value even if no locks found, so prometheus description will be created
		ret = append(ret, SmbStatisticsNumeric{"locks_per_share_count", float64(0), "Number of locks on share", map[string]string{"share": ""}})
	}

	ret = append(ret, SmbStatisticsNumeric{"server_information", 1, "Version of the samba server", map[string]string{"version": sambaVersion}})

	if !settings.DoNotExportEncryption {
		if len(protocolVersionCount) > 0 {
			for version, count := range protocolVersionCount {
				ret = append(ret, SmbStatisticsNumeric{"protocol_version_count", float64(count), "Number of processes on the server using the protocol", map[string]string{"protocol_version": version}})
			}
		} else {
			ret = append(ret, SmbStatisticsNumeric{"protocol_version_count", float64(0), "Number of processes on the server using the protocol", map[string]string{"protocol_version": ""}})
		}

		if len(signingMethodCount) > 0 {
			for method, count := range signingMethodCount {
				ret = append(ret, SmbStatisticsNumeric{"signing_method_count", float64(count), "Number of processes on the server using the signing", map[string]string{"signing": method}})
			}
		} else {
			ret = append(ret, SmbStatisticsNumeric{"signing_method_count", float64(0), "Number of processes on the server using the signing", map[string]string{"signing": ""}})
		}

		if len(encryptionMethodCount) > 0 {
			for method, count := range encryptionMethodCount {
				ret = append(ret, SmbStatisticsNumeric{"encryption_method_count", float64(count), "Number of processes on the server using the encryption", map[string]string{"encryption": method}})
			}
		} else {
			ret = append(ret, SmbStatisticsNumeric{"encryption_method_count", float64(0), "Number of processes on the server using the encryption", map[string]string{"encryption": ""}})
		}
	}

	if !settings.DoNotExportClient {
		if len(processPerClient) > 0 {
			for client, count := range processPerClient {
				ret = append(ret, SmbStatisticsNumeric{"process_per_client_count", float64(count), "Number of processes on the server used by one client", map[string]string{"client": client}})
			}
		} else {
			ret = append(ret, SmbStatisticsNumeric{"process_per_client_count", float64(0), "Number of processes on the server used by one client", map[string]string{"client": ""}})
		}

		if len(clientConnectionTime) > 0 {
			for client, connectTime := range clientConnectionTime {
				ret = append(ret, SmbStatisticsNumeric{"client_connected_at", float64(connectTime), "Unix time stamp a client connected", map[string]string{"client": client}})
				now := time.Now()
				connected_since := now.Sub(time.Unix(connectTime, 0))
				ret = append(ret, SmbStatisticsNumeric{"client_connected_since_seconds", connected_since.Seconds(), "Seconds since a client connected", map[string]string{"client": client}})
			}
		} else {
			// Add this values even if no locks found, so prometheus description will be created
			ret = append(ret, SmbStatisticsNumeric{"client_connected_at", float64(0), "Unix time stamp a client connected", map[string]string{"client": ""}})
			ret = append(ret, SmbStatisticsNumeric{"client_connected_since_seconds", float64(0), "Seconds since a client connected", map[string]string{"client": ""}})
		}
	}

	if !settings.DoNotExportUser {
		if len(lockCreationEntries) > 0 {
			for _, lockEntry := range lockCreationEntries {
				ret = append(ret, SmbStatisticsNumeric{"lock_created_at", float64(lockEntry.CreationTime.Unix()),
					"Unix time stamp a lock was created",
					map[string]string{"user": strconv.Itoa(lockEntry.UserID), "share": lockEntry.Share}})

				ret = append(ret, SmbStatisticsNumeric{"lock_created_since_seconds", float64(time.Since(lockEntry.CreationTime).Seconds()),
					"Seconds since a lock was created",
					map[string]string{"user": strconv.Itoa(lockEntry.UserID), "share": lockEntry.Share}})
			}
		} else {
			ret = append(ret, SmbStatisticsNumeric{"lock_created_at", float64(0), "Unix time stamp a lock was created", map[string]string{"user": "", "share": ""}})
			ret = append(ret, SmbStatisticsNumeric{"lock_created_since_seconds", float64(0), "Seconds since a lock was created", map[string]string{"user": "", "share": ""}})
		}
	}

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

func lockArrContainsEntry(arr []lockCreationEntry, value lockCreationEntry) bool {
	for _, field := range arr {
		if field.Share == value.Share && field.UserID == value.UserID {
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
