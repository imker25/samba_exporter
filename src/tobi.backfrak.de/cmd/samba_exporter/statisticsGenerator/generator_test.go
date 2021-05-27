package statisticsGenerator

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"testing"

	"tobi.backfrak.de/cmd/samba_exporter/smbstatusreader"
	"tobi.backfrak.de/internal/smbstatusout"
)

func TestGetSmbStatisticsEmptyData(t *testing.T) {
	locks := smbstatusreader.GetLockData(smbstatusout.LockData0Line)
	shares := smbstatusreader.GetShareData(smbstatusout.ShareData0Line)
	processes := smbstatusreader.GetProcessData(smbstatusout.ProcessData0Lines)

	ret := GetSmbStatistics(locks, processes, shares)

	if len(ret) != 6 {
		t.Errorf("The number of resturn values %d was not expected", len(ret))
	}

	for _, field := range ret {
		if field.Value != 0 {
			t.Errorf("The value is not 0 when reading only empty tables")
		}
	}

}

func TestGetSmbStatisticsEmptyResponseLabels(t *testing.T) {
	locks := smbstatusreader.GetLockData(smbstatusout.LockData0Line)
	shares := smbstatusreader.GetShareData(smbstatusout.ShareData0Line)
	processes := smbstatusreader.GetProcessData(smbstatusout.ProcessData0Lines)

	ret := GetSmbStatistics(locks, processes, shares)
	if len(ret) != 6 {
		t.Errorf("The number of resturn values %d was not expected", len(ret))
	}

	if ret[5].Name != "locks_per_share_count" {
		t.Errorf("The Name \"%s\" is not expected", ret[5].Name)
	}

	if ret[5].Labels["share"] != "" {
		t.Errorf("The Labels[\"share\"] %s is not expected", ret[5].Labels["share"])
	}
}

func TestGetSmbStatistics(t *testing.T) {
	locks := smbstatusreader.GetLockData(smbstatusout.LockData4Lines)
	shares := smbstatusreader.GetShareData(smbstatusout.ShareData4Lines)
	processes := smbstatusreader.GetProcessData(smbstatusout.ProcessData4Lines)

	ret := GetSmbStatistics(locks, processes, shares)

	if len(ret) != 9 {
		t.Errorf("The number of resturn values %d was not expected", len(ret))
	}

	if ret[0].Name != "individual_user_count" {
		t.Errorf("The individual_user_count is not at expecgted place")
	}

	if ret[0].Value != 1 {
		t.Errorf("The individual_user_count is not the expected value")
	}

	if ret[1].Name != "locked_file_count" {
		t.Errorf("The locked_file_count is not at expecgted place")
	}

	if ret[1].Value != float64(len(locks)) {
		t.Errorf("The locked_file_count is not the expected value")
	}

	if ret[2].Name != "pid_count" {
		t.Errorf("The pid_count is not at expecgted place")
	}

	if ret[2].Value != 4 {
		t.Errorf("The pid_count is not the expected value")
	}

	if ret[3].Name != "share_count" {
		t.Errorf("The share_count is not at expecgted place")
	}

	if ret[3].Value != 4 {
		t.Errorf("The share_count is not the expected value")
	}

	if ret[4].Name != "client_count" {
		t.Errorf("The client_countis not at expecgted place")
	}

	if ret[4].Value != 4 {
		t.Errorf("The client_count is not the expected value")
	}

}

func TestStringArrContains(t *testing.T) {
	arr := []string{"a", "b", "c"}

	if strArrContains(arr, "a") == false {
		t.Errorf("strArrContains returns false but should true")
	}

	if strArrContains(arr, "z") == true {
		t.Errorf("strArrContains returns true but should false")
	}
}

func TestIntArrContains(t *testing.T) {
	arr := []int{1, 2, 3}

	if intArrContains(arr, 2) == false {
		t.Errorf("strArrContains returns false but should true")
	}

	if intArrContains(arr, 100) == true {
		t.Errorf("strArrContains returns true but should false")
	}
}
