package smbstatusdbl

import "testing"

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

func TestNewPsDataGenerator(t *testing.T) {
	processImage := "my_pid"
	sut, err := NewPsDataGenerator(processImage)
	if err != nil {
		t.Errorf("Error when getting a new PsDataGenerator")
	}

	if sut.ProcessToRequest != processImage {
		t.Errorf("The sut.ProcessToRequest '%s' is not the expected '%s'", sut.ProcessToRequest, processImage)
	}
}

func TestPsDataGeneratorNotRunningProcess(t *testing.T) {
	processImage := "my__not_existing_pid"
	sut, errNew := NewPsDataGenerator(processImage)
	if errNew != nil {
		t.Errorf("Error when getting a new PsDataGenerator: %s", errNew.Error())
	}

	if sut.ProcessToRequest != processImage {
		t.Errorf("The sut.ProcessToRequest '%s' is not the expected '%s'", sut.ProcessToRequest, processImage)
	}

	pidList, errPidList := sut.getPidList()
	if errPidList == nil {
		t.Errorf("No error when getting a pid list for not running process image")
	}

	if len(pidList) != 0 {
		t.Errorf("Expected an empty list but got %d entries", len(pidList))
	}

	pidData, errData := sut.GetPsUtilPidData()
	if errData != nil {
		t.Errorf("Error when getting a pid data: %s", errData.Error())
	}

	if len(pidData) != 0 {
		t.Errorf("Expected an empty list but got %d entries", len(pidData))
	}

}

func TestPsDataGeneratorRunningProcess(t *testing.T) {
	processImage := "go"
	sut, errNew := NewPsDataGenerator(processImage)
	if errNew != nil {
		t.Errorf("Error when getting a new PsDataGenerator: %s", errNew.Error())
	}

	if sut.ProcessToRequest != processImage {
		t.Errorf("The sut.ProcessToRequest '%s' is not the expected '%s'", sut.ProcessToRequest, processImage)
	}

	pidList, errPidList := sut.getPidList()
	if errPidList != nil {
		t.Errorf("Error when getting a pid list for running process image '%s': %s", processImage, errPidList.Error())
	}

	if len(pidList) == 0 {
		t.Errorf("Expected not an empty list")
	}

	// Disable this tests, since the latest ubuntu kernel seems to prevent normal users
	// from accessing /proc/<PID>/io:
	// pidData, errData := sut.GetPsUtilPidData()
	// if errData != nil {
	// 	t.Errorf("Error when getting a pid data: %s", errData.Error())
	// }

	// if len(pidData) == 0 {
	// 	t.Errorf("Expected not an empty list")
	// }

	// if len(pidData) != len(pidList) {
	// 	t.Errorf("Got '%d' data entries but '%d' pids", len(pidData), len(pidList))
	// }

}
