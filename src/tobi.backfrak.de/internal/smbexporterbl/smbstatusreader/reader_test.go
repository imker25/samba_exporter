package smbstatusreader

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"strings"
	"testing"
	"time"

	"tobi.backfrak.de/internal/commonbl"
	"tobi.backfrak.de/internal/smbstatusout"
)

func TestStringerLockData(t *testing.T) {
	logger := commonbl.NewLogger(true)
	oneLock := GetLockData(smbstatusout.LockDataOneLine, logger)[0]

	lockStr := oneLock.String()
	if strings.Contains(lockStr, "UserID: 1080;") == false {
		t.Errorf("The string does not contain the expected sub string")
	}

	if strings.Contains(lockStr, "SharePath: /usr/share/data;") == false {
		t.Errorf("The string does not contain the expected sub string")
	}

	if strings.Contains(lockStr, "ClusterNodeId: ") == true {
		t.Errorf("The string does contain the expected sub string")
	}

	oneLock = GetLockData(smbstatusout.LockDataCluster, logger)[0]
	lockStr = oneLock.String()
	if strings.Contains(lockStr, "ClusterNodeId: 1;") == false {
		t.Errorf("The string does contain the expected sub string")
	}

	if strings.Contains(lockStr, "SharePath: /lfsmnt/dst01") == false {
		t.Errorf("The string does not contain the expected sub string")
	}

}

func TestStringerShareData(t *testing.T) {
	logger := commonbl.NewLogger(true)
	oneShare := GetShareData(smbstatusout.ShareDataOneLine, logger)[0]

	shareStr := oneShare.String()
	if strings.Contains(shareStr, "PID: 1119;") == false {
		t.Errorf("The string does not contain the expected sub string")
	}

	if strings.Contains(shareStr, "Machine: 192.168.1.242;") == false {
		t.Errorf("The string does not contain the expected sub string")
	}
	if strings.Contains(shareStr, "ClusterNodeId: ") == true {
		t.Errorf("The string does contain the expected sub string")
	}

	oneShare = GetShareData(smbstatusout.ShareDataCluster, logger)[0]

	shareStr = oneShare.String()
	if strings.Contains(shareStr, "PID: 19801;") == false {
		t.Errorf("The string does not contain the expected sub string")
	}
	if strings.Contains(shareStr, "ClusterNodeId: 1;") == false {
		t.Errorf("The string does contain the expected sub string")
	}

}

func TestStringerProcessData(t *testing.T) {
	logger := commonbl.NewLogger(true)
	oneProcess := GetProcessData(smbstatusout.ProcessDataOneLine, logger)[0]

	shareStr := oneProcess.String()
	if strings.Contains(shareStr, "PID: 1117;") == false {
		t.Errorf("The string does not contain the expected sub string")
	}

	if strings.Contains(shareStr, "Machine: 192.168.1.242 (ipv4:192.168.1.242:42296);") == false {
		t.Errorf("The string does not contain the expected sub string")
	}
	if strings.Contains(shareStr, "ClusterNodeId: ") == true {
		t.Errorf("The string does contain the expected sub string")
	}

	oneProcess = GetProcessData(smbstatusout.ProcessDataCluster, logger)[0]

	shareStr = oneProcess.String()
	if strings.Contains(shareStr, "PID: 57086;") == false {
		t.Errorf("The string does not contain the expected sub string")
	}
	if strings.Contains(shareStr, "ClusterNodeId: 3;") == false {
		t.Errorf("The string does contain the expected sub string")
	}
}

func TestGetLockDataOneLine(t *testing.T) {
	logger := commonbl.NewLogger(true)
	oneEntry := GetLockData(smbstatusout.LockDataOneLine, logger)

	if len(oneEntry) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(oneEntry))
	}

	if oneEntry[0].PID != 1120 {
		t.Errorf("The PID %d is not the expected 1120", oneEntry[0].PID)
	}

	if oneEntry[0].UserID != 1080 {
		t.Errorf("The UserID %d is not the expected 1080", oneEntry[0].UserID)
	}

	if oneEntry[0].DenyMode != "DENY_NONE" {
		t.Errorf("The DenyMode %s is not the expected DENY_NONE", oneEntry[0].DenyMode)
	}

	if oneEntry[0].Access != "0x80" {
		t.Errorf("The Access %s is not the expected 0x80", oneEntry[0].Access)
	}

	if oneEntry[0].AccessMode != "RDONLY" {
		t.Errorf("The AccessMode %s is not the expected RDONLY", oneEntry[0].AccessMode)
	}

	if oneEntry[0].Oplock != "NONE" {
		t.Errorf("The Oplock %s is not the expected NONE", oneEntry[0].Oplock)
	}

	if oneEntry[0].SharePath != "/usr/share/data" {
		t.Errorf("The SharePath %s is not the expected /usr/share/data", oneEntry[0].SharePath)
	}

	if oneEntry[0].Name != "." {
		t.Errorf("The Name %s is not the expected \".\"", oneEntry[0].Name)
	}

	expectDate, _ := time.ParseInLocation(time.ANSIC, "Sun May 16 12:07:02 2021", time.Now().Location())

	if oneEntry[0].Time != expectDate {
		t.Errorf("The Time %s is not the expected Sun May 16 12:07:02 2021", oneEntry[0].Time)
	}
}

func TestGetLockData4Line(t *testing.T) {
	logger := commonbl.NewLogger(true)
	entryList := GetLockData(smbstatusout.LockData4Lines, logger)

	if len(entryList) != 4 {
		t.Errorf("Got %d entries, expected 4", len(entryList))
	}

	if entryList[0].SharePath != "/usr/share/data" {
		t.Errorf("The SharePath %s is not the expected /usr/share/data", entryList[0].SharePath)
	}

	if entryList[1].SharePath != "/usr/share/foto" {
		t.Errorf("The SharePath %s is not the expected /usr/share/foto", entryList[1].SharePath)
	}
	if entryList[2].SharePath != "/usr/share/film" {
		t.Errorf("The SharePath %s is not the expected /usr/share/film", entryList[2].SharePath)
	}
	if entryList[3].SharePath != "/usr/share/music" {
		t.Errorf("The SharePath %s is not the expected /usr/share/music", entryList[3].SharePath)
	}

	if entryList[3].ClusterNodeId != -1 {
		t.Errorf("The SharePath %d is not the expected '-1'", entryList[3].ClusterNodeId)
	}
}

func TestGetLockDataCluster(t *testing.T) {
	logger := commonbl.NewLogger(true)
	entryList := GetLockData(smbstatusout.LockDataCluster, logger)

	if len(entryList) != 7 {
		t.Errorf("Got %d entries, expected 4", len(entryList))
	}

	if entryList[0].SharePath != "/lfsmnt/dst01" {
		t.Errorf("The SharePath %s is not the expected '/lfsmnt/dst01'", entryList[0].SharePath)
	}

	if entryList[5].Name != "share/data2/CLIPS001/CC0639/CC063904.MXF" {
		t.Errorf("The Name %s is not the expected 'share/data2/CLIPS001/CC0639/CC063904.MXF'", entryList[3].Name)
	}

	if entryList[6].Name != "share/test.wav 48000.pek" {
		t.Errorf("The Name %s is not the expected 'share/test.wav 48000.pek'", entryList[0].Name)
	}

	if entryList[5].PID != 57086 {
		t.Errorf("Got %d entryList[5].PID, expected 57086", entryList[3].PID)
	}

	if entryList[5].ClusterNodeId != 3 {
		t.Errorf("Got %d entryList[5].ClusterNodeId, expected 3", entryList[3].ClusterNodeId)
	}

	if entryList[6].Time.Format(time.ANSIC) != "Tue Apr  4 14:13:28 2023" {
		t.Errorf("The time %s is not expected", entryList[6].Time.Format(time.ANSIC))
	}
}

func TestGetLockDataWrongInput(t *testing.T) {
	logger := commonbl.NewLogger(true)
	entryList := GetLockData(smbstatusout.ProcessData4Lines, logger)

	if len(entryList) != 0 {
		t.Errorf("Got entries when reading wrong input")
	}
}

func TestGetLockData0Input(t *testing.T) {
	logger := commonbl.NewLogger(true)
	entryList := GetLockData(smbstatusout.LockData0Line, logger)

	if len(entryList) != 0 {
		t.Errorf("Got entries when reading wrong input")
	}
}

func TestGetLockDataNoDta(t *testing.T) {
	logger := commonbl.NewLogger(true)
	entryList := GetLockData(smbstatusout.LockDataNoData, logger)

	if len(entryList) != 0 {
		t.Errorf("Got entries when reading wrong input")
	}
}

func TestGetShareDataDifferentTimeStampLines(t *testing.T) {
	logger := commonbl.NewLogger(true)
	entryList := GetShareData(smbstatusout.ShareDataDifferentTimeStampLines, logger)

	if len(entryList) != 3 {
		t.Errorf("Got wrong amount of entries %d", len(entryList))
	}

	if entryList[1].ConnectedAt.Format(time.ANSIC) != "Wed Jun  2 21:32:31 2021" {
		t.Errorf("The time %s is not expected", entryList[1].ConnectedAt.Format(time.ANSIC))
	}

	if entryList[2].ConnectedAt.Format(time.ANSIC) != "Mon Sep 19 18:34:17 2022" {
		t.Errorf("The time %s is not expected", entryList[2].ConnectedAt.Format(time.ANSIC))
	}
}

func TestGetShareDataOneLine(t *testing.T) {
	logger := commonbl.NewLogger(true)
	oneEntry := GetShareData(smbstatusout.ShareDataOneLine, logger)

	if len(oneEntry) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(oneEntry))
	}

	if oneEntry[0].PID != 1119 {
		t.Errorf("The PID %d is not the expected 1119", oneEntry[0].PID)
	}

	if oneEntry[0].Service != "IPC$" {
		t.Errorf("The Service %s is not the expected IPC$", oneEntry[0].Service)
	}

	if oneEntry[0].Machine != "192.168.1.242" {
		t.Errorf("The Machine %s is not the expected 192.168.1.242 ", oneEntry[0].Machine)
	}

	if oneEntry[0].Encryption != "-" {
		t.Errorf("The Encryption %s is not the expected \"-\" ", oneEntry[0].Encryption)
	}

	if oneEntry[0].Signing != "-" {
		t.Errorf("The Signing %s is not the expected \"-\" ", oneEntry[0].Signing)
	}

	if oneEntry[0].ConnectedAt.Format(time.ANSIC) != "Sun May 16 11:55:36 2021" {
		t.Errorf("The ConnectedAt %s is not the expected Sun May 16 11:55:36 2021", oneEntry[0].ConnectedAt.Format(time.ANSIC))
	}
}

func TestGetShareData4Line(t *testing.T) {
	logger := commonbl.NewLogger(true)
	entries := GetShareData(smbstatusout.ShareData4Lines, logger)

	if len(entries) != 4 {
		t.Errorf("Got %d entries, expected 4", len(entries))
	}

	if entries[0].ConnectedAt.Format(time.ANSIC) != "Sun May 16 11:55:36 2021" {
		t.Errorf("The ConnectedAt %s is not the expected Sun May 16 11:55:36 2021", entries[0].ConnectedAt.Format(time.ANSIC))
	}

	if entries[1].ConnectedAt.Format(time.ANSIC) != "Mon May 17 10:56:56 2021" {
		t.Errorf("The ConnectedAt %s is not the expected Mon May 17 10:56:56 2021", entries[1].ConnectedAt.Format(time.ANSIC))
	}

	if entries[2].ConnectedAt.Format(time.ANSIC) != "Tue May 18 09:52:38 2021" {
		t.Errorf("The ConnectedAt %s is not the expected Tue May 18 09:52:38 2021", entries[2].ConnectedAt.Format(time.ANSIC))
	}

	if entries[3].ConnectedAt.Format(time.ANSIC) != "Fri Nov  5 23:07:13 2021" {
		t.Errorf("The ConnectedAt %s is not the expected 'Fri Nov  5 23:07:13 2021'", entries[3].ConnectedAt.Format(time.ANSIC))
	}

	if entries[3].ClusterNodeId != -1 {
		t.Errorf("The ClusterNodeId %d is not the expected '-1'", entries[3].ClusterNodeId)
	}
}

func TestGetShareDataCluster(t *testing.T) {
	logger := commonbl.NewLogger(true)
	entries := GetShareData(smbstatusout.ShareDataCluster, logger)

	if len(entries) != 16 {
		t.Errorf("Got %d entries, expected 16", len(entries))
	}

	if entries[0].PID != 19801 {
		t.Errorf("Got %d entries[0].PID, expected 19801", entries[0].PID)
	}

	if entries[15].PID != 42597 {
		t.Errorf("Got %d entries[15].PID, expected 42597", entries[0].PID)
	}

	if entries[15].Encryption != "-" {
		t.Errorf("Got %s entries[15].Encryption, expected '-'", entries[0].Encryption)
	}

	if entries[0].Signing != "-" {
		t.Errorf("Got %s entries[0].Signing, expected '-'", entries[0].Signing)
	}

	if entries[3].Machine != "10.63.0.11 (ipv4:10.63.0.11:50370)" {
		t.Errorf("Got %s entries[3].Signing, expected '10.63.0.11 (ipv4:10.63.0.11:50370) '", entries[3].Machine)
	}

	if entries[3].ClusterNodeId != 1 {
		t.Errorf("Got %d entries[3].ClusterNodeId, expected '1'", entries[3].ClusterNodeId)
	}
}

func TestGetShareDataWrongData(t *testing.T) {
	logger := commonbl.NewLogger(true)
	entries := GetShareData(smbstatusout.LockData4Lines, logger)

	if len(entries) != 0 {
		t.Errorf("Got %d entries, but expected none", len(entries))
	}
}

func TestGetShareData0Input(t *testing.T) {
	logger := commonbl.NewLogger(true)
	entryList := GetShareData(smbstatusout.ShareData0Line, logger)

	if len(entryList) != 0 {
		t.Errorf("Got entries when reading wrong input")
	}
}

func TestGetProcessDataOneLine(t *testing.T) {
	logger := commonbl.NewLogger(true)
	oneProcess := GetProcessData(smbstatusout.ProcessDataOneLine, logger)

	if len(oneProcess) != 1 {
		t.Errorf("Got %d entries, expected 1", len(oneProcess))
	}

	if oneProcess[0].PID != 1117 {
		t.Errorf("The PID %d is not the expected 1117", oneProcess[0].PID)
	}

	if oneProcess[0].UserID != 1080 {
		t.Errorf("The UserID %d is not the expected 1080", oneProcess[0].UserID)
	}

	if oneProcess[0].GroupID != 117 {
		t.Errorf("The Group %d is not the expected 117", oneProcess[0].GroupID)
	}

	if oneProcess[0].Machine != "192.168.1.242 (ipv4:192.168.1.242:42296)" {
		t.Errorf("The Machine \"%s\" is not the expected \"192.168.1.242 (ipv4:192.168.1.242:42296)\"", oneProcess[0].Machine)
	}

	if oneProcess[0].ProtocolVersion != "SMB3_11" {
		t.Errorf("The ProtocolVersion \"%s\" is not the expected \"SMB3_11\"", oneProcess[0].ProtocolVersion)
	}

	if oneProcess[0].Encryption != "-" {
		t.Errorf("The Encryption \"%s\" is not the expected \"-\"", oneProcess[0].Encryption)
	}

	if oneProcess[0].Signing != "partial(AES-128-CMAC)" {
		t.Errorf("The Signing \"%s\" is not the expected \"partial(AES-128-CMAC)\"", oneProcess[0].Signing)
	}
}

func TestGetProcessData4Line(t *testing.T) {
	logger := commonbl.NewLogger(true)
	enties := GetProcessData(smbstatusout.ProcessData4Lines, logger)

	if len(enties) != 4 {
		t.Errorf("Got %d entries, expected 1", len(enties))
	}

	if enties[0].Machine != "192.168.1.242 (ipv4:192.168.1.242:42296)" {
		t.Errorf("The Machine \"%s\" is not the expected \"192.168.1.242 (ipv4:192.168.1.242:42296)\"", enties[0].Machine)
	}

	if enties[1].Machine != "192.168.1.243 (ipv4:192.168.1.243:47510)" {
		t.Errorf("The Machine \"%s\" is not the expected \"192.168.1.243 (ipv4:192.168.1.243:47510)\"", enties[1].Machine)
	}
	if enties[2].Machine != "192.168.1.244 (ipv4:192.168.1.244:47512)" {
		t.Errorf("The Machine \"%s\" is not the expected \"192.168.1.244 (ipv4:192.168.1.244:47512)\"", enties[2].Machine)
	}

	if enties[3].Machine != "192.168.1.245 (ipv4:192.168.1.245:47514)" {
		t.Errorf("The Machine \"%s\" is not the expected \"192.168.1.245 (ipv4:192.168.1.245:47514)\"", enties[3].Machine)
	}

	if enties[3].ClusterNodeId != -1 {
		t.Errorf("The ClusterNodeId \"%d\" is not the expected \"-1\"", enties[3].ClusterNodeId)
	}

	for _, entry := range enties {
		if entry.SambaVersion != "4.11.6-Ubuntu" {
			t.Errorf("The SambaVersion \"%s\" is not expected", entry.SambaVersion)
		}
	}
}

func TestGetProcessDataCluster(t *testing.T) {
	logger := commonbl.NewLogger(true)
	enties := GetProcessData(smbstatusout.ProcessDataCluster, logger)

	if len(enties) != 7 {
		t.Errorf("Got %d entries, expected 7", len(enties))
	}

	if enties[0].Machine != "10.63.0.41 (ipv4:10.63.0.41:62834)" {
		t.Errorf("The Machine \"%s\" is not the expected \"10.63.0.41 (ipv4:10.63.0.41:62834)\"", enties[0].Machine)
	}

	if enties[3].Machine != "10.63.0.28 (ipv4:10.63.0.28:58968)" {
		t.Errorf("The Machine \"%s\" is not the expected \"10.63.0.28 (ipv4:10.63.0.28:58968)\"", enties[3].Machine)
	}

	if enties[3].ClusterNodeId != 3 {
		t.Errorf("The ClusterNodeId \"%d\" is not the expected \"3\"", enties[3].ClusterNodeId)
	}

	for _, entry := range enties {
		if entry.SambaVersion != "4.9.5-Debian" {
			t.Errorf("The SambaVersion \"%s\" is not expected", entry.SambaVersion)
		}
	}
}

func TestGetProcessDataWrongData(t *testing.T) {
	logger := commonbl.NewLogger(true)
	enties := GetProcessData(smbstatusout.LockData4Lines, logger)

	if len(enties) != 0 {
		t.Errorf("Got %d entries, but expected none", len(enties))
	}
}

func TestGetProcessData0Input(t *testing.T) {
	logger := commonbl.NewLogger(true)
	entryList := GetProcessData(smbstatusout.ProcessData0Lines, logger)

	if len(entryList) != 0 {
		t.Errorf("Got entries when reading wrong input")
	}
}

func TestGetPsData0Input(t *testing.T) {
	logger := commonbl.NewLogger(true)
	entryList := GetPsData("", logger)

	if len(entryList) != 0 {
		t.Errorf("Got entries when reading wrong input")
	}
}

func TestGetPsDataEmptyInput(t *testing.T) {
	logger := commonbl.NewLogger(true)
	jsonData := commonbl.TestPsResponseEmpty()
	entryList := GetPsData(string(jsonData), logger)

	if len(entryList) != 0 {
		t.Errorf("Got entries when reading wrong input")
	}
}

func TestGetPsDataTwoPids(t *testing.T) {
	logger := commonbl.NewLogger(true)
	jsonData := commonbl.TestPsResponse()
	entryList := GetPsData(string(jsonData), logger)

	if len(entryList) != 2 {
		t.Errorf("Got %d entries but expected 2", len(entryList))
	}
}

func TestTryGetTimeStampFromStrArr(t *testing.T) {
	var suc bool
	var value time.Time
	fields := []string{"", ""}
	suc, _ = tryGetTimeStampFromStrArr(fields)
	if suc == true {
		t.Errorf("Got a time from an empty string")
	}

	fields = []string{"/my/cool/path", "RW"}
	suc, _ = tryGetTimeStampFromStrArr(fields)
	if suc == true {
		t.Errorf("Got a time from an empty string")
	}

	fields = []string{"Fri", "Nov", "5", "11:07:13", "PM", "2021", "CET"}
	suc, value = tryGetTimeStampFromStrArr(fields)
	if suc == false {
		t.Errorf("Got no time from \"Fri Nov 5 11:07:13 PM 2021 CET\"")
	}

	if value.Format(time.ANSIC) != "Fri Nov  5 23:07:13 2021" {
		t.Errorf("Time is '%s', but expected 'Fri Nov  5 23:07:13 2021'", value.Format(time.ANSIC))
	}

	fields = []string{"Fri", "Nov", "05", "11:07:13", "PM", "2021", "CET"}
	suc, value = tryGetTimeStampFromStrArr(fields)
	if suc == false {
		t.Errorf("Got no time from \"Fri Nov 5 11:07:13 PM 2021 CET\"")
	}

	if value.Format(time.ANSIC) != "Fri Nov  5 23:07:13 2021" {
		t.Errorf("Time is '%s', but expected 'Fri Nov  5 23:07:13 2021'", value.Format(time.ANSIC))
	}

	fields = []string{"Wed", "Jun", "2", "21:32:31 2021", "UTC"}
	suc, value = tryGetTimeStampFromStrArr(fields)
	if suc == false {
		t.Errorf("Got no time from \"Wed Jun  2 21:32:31 2021 UTC\"")
	}

	if value.Format(time.ANSIC) != "Wed Jun  2 21:32:31 2021" {
		t.Errorf("Time is '%s', but expected 'Wed Jun  2 21:32:31 2021'", value.Format(time.ANSIC))
	}

	fields = []string{"Wed", "Jun", " 2", "21:32:31 2021", "UTC"}
	suc, value = tryGetTimeStampFromStrArr(fields)
	if suc == false {
		t.Errorf("Got no time from \"Wed Jun  2 21:32:31 2021 UTC\"")
	}

	if value.Format(time.ANSIC) != "Wed Jun  2 21:32:31 2021" {
		t.Errorf("Time is '%s', but expected 'Wed Jun  2 21:32:31 2021'", value.Format(time.ANSIC))
	}

	fields = []string{"Wed", "Jun", "02", "21:32:31 2021", "UTC"}
	suc, value = tryGetTimeStampFromStrArr(fields)
	if suc == false {
		t.Errorf("Got no time from \"Wed Jun 02 21:32:31 2021 UTC\"")
	}

	if value.Format(time.ANSIC) != "Wed Jun  2 21:32:31 2021" {
		t.Errorf("Time is '%s', but expected 'Wed Jun  2 21:32:31 2021'", value.Format(time.ANSIC))
	}

	fields = []string{"Wed", "Jun", " 2", "21:32:31 2021"}
	suc, value = tryGetTimeStampFromStrArr(fields)
	if suc == false {
		t.Errorf("Got no time from \"Wed Jun  2 21:32:31 2021 UTC\"")
	}

	if value.Format(time.ANSIC) != "Wed Jun  2 21:32:31 2021" {
		t.Errorf("Time is '%s', but expected 'Wed Jun  2 21:32:31 2021'", value.Format(time.ANSIC))
	}
}
