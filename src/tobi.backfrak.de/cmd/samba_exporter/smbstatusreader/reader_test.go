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

	"tobi.backfrak.de/internal/smbstatusout"
)

func TestStringer(t *testing.T) {
	oneLock := GetLockData(smbstatusout.LockDataOneLine)[0]

	lockStr := oneLock.String()
	if strings.Contains(lockStr, "UserID: 1080;") == false {
		t.Errorf("The string does not contain the expected sub string")
	}

	if strings.Contains(lockStr, "SharePath: /usr/share/data;") == false {
		t.Errorf("The string does not contain the expected sub string")
	}
}

func TestGetLockDataOneLine(t *testing.T) {
	oneEntry := GetLockData(smbstatusout.LockDataOneLine)

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

	expectDate, _ := time.Parse(time.ANSIC, "Sun May 16 12:07:02 2021")

	if oneEntry[0].Time != expectDate {
		t.Errorf("The Time %s is not the expected Sun May 16 12:07:02 2021", oneEntry[0].Time)
	}
}

func TestGetLockData4Line(t *testing.T) {

	entryList := GetLockData(smbstatusout.LockData4Lines)

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
}

func TestGetLockDataWrongInput(t *testing.T) {
	entryList := GetLockData(smbstatusout.ProcessData4Lines)

	if len(entryList) != 0 {
		t.Errorf("Got entries when reading wrong input")
	}
}

func TestGetShareDataOneLine(t *testing.T) {
	oneEntry := GetShareData(smbstatusout.ShareDataOneLine)

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
	entries := GetShareData(smbstatusout.ShareData4Lines)

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

	if entries[3].ConnectedAt.Format(time.ANSIC) != "Fri May 21 18:46:29 2021" {
		t.Errorf("The ConnectedAt %s is not the expected Fri May 21 18:46:29 2021", entries[3].ConnectedAt.Format(time.ANSIC))
	}
}

func TestGetShareDataWrongData(t *testing.T) {
	entries := GetShareData(smbstatusout.LockData4Lines)

	if len(entries) != 0 {
		t.Errorf("Got %d entries, but expected none", len(entries))
	}
}
