package commonbl

import (
	"encoding/json"
)

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

// Contains the test data for a Lock Response
const TestLockResponse = `
Locked files:
Pid          User(ID)   DenyMode   Access      R/W        Oplock           SharePath   Name   Time
--------------------------------------------------------------------------------------------------
1120         1080       DENY_NONE  0x80        RDONLY     NONE             /usr/share/data   .   Sun May 16 12:07:02 2021`

// Contains the test data for a Service Response
const TestShareResponse = `
Service      pid     Machine       Connected at                      Encryption   Signing     
---------------------------------------------------------------------------------------------
IPC$         1119    192.168.1.242  Sun May 16 11:55:36 AM 2021 CEST -            -           `

// Contains the test data for a Process Response
const TestProcessResponse = `
Samba version 4.11.6-Ubuntu
PID     Username     Group        Machine                                   Protocol Version  Encryption           Signing              
----------------------------------------------------------------------------------------------------------------------------------------
1117    1080    117     192.168.1.242 (ipv4:192.168.1.242:42296)  SMB3_11           -                    partial(AES-128-CMAC)`

func TestPsResponse() string {

	jsonData, _ := json.MarshalIndent(GetTestPsUtilPidData(), "", " ")

	return string(jsonData)
}

func TestPsResponseEmpty() string {

	jsonData, _ := json.MarshalIndent([]PsUtilPidData{}, "", " ")

	return string(jsonData)
}

// Always returns the same PsUtilPidData for test propose
func GetTestPsUtilPidData() []PsUtilPidData {
	pidData := []PsUtilPidData{}
	pidData = append(pidData, PsUtilPidData{
		1234,
		0.023,
		456789,
		0.0034,
		123456,
		789123,
		2345,
		6789,
		1467,
		8765,
	})

	pidData = append(pidData, PsUtilPidData{
		4234,
		0.123,
		8789,
		0.5034,
		23456,
		912378,
		34576,
		789543,
		467123,
		765853,
	})

	return pidData
}
