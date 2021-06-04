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

func TestGetSmbStatisticsNoLockData(t *testing.T) {
	locks := smbstatusreader.GetLockData(smbstatusout.LockDataNoData)
	shares := smbstatusreader.GetShareData(smbstatusout.ShareDataOneLine)
	processes := smbstatusreader.GetProcessData(smbstatusout.ProcessDataOneLine)

	ret := GetSmbStatistics(locks, processes, shares)

	if len(ret) != 8 {
		t.Errorf("The number of resturn values %d was not expected", len(ret))
	}

	if ret[0].Name != "individual_user_count" && ret[0].Value != 1.0 {
		t.Errorf("The individual_user_count does not match as expected")
	}
}

func TestGetSmbStatisticsEmptyData(t *testing.T) {
	locks := smbstatusreader.GetLockData(smbstatusout.LockData0Line)
	shares := smbstatusreader.GetShareData(smbstatusout.ShareData0Line)
	processes := smbstatusreader.GetProcessData(smbstatusout.ProcessData0Lines)

	ret := GetSmbStatistics(locks, processes, shares)

	if len(ret) != 8 {
		t.Errorf("The number of resturn values %d was not expected", len(ret))
	}

	for _, field := range ret[0:6] {
		if field.Value != 0 {
			t.Errorf("The value is not 0 when reading only empty tables")
		}
	}

	if ret[6].Name != "server_information" {
		t.Errorf("The Name \"%s\" is not expected", ret[6].Name)
	}

	if ret[6].Value != 1 {
		t.Errorf("The Value %f is not expected", ret[6].Value)
	}

	if len(ret[6].Labels) != 1 {
		t.Errorf("There are more labels than expected")
	}

	value, found := ret[6].Labels["version"]
	if !found {
		t.Errorf("No label with key \"version\" found")
	}

	if value != "" {
		t.Errorf("The SambaVersion \"%s\" is not expected", value)
	}

}

func TestGetSmbStatisticsEmptyResponseLabels(t *testing.T) {
	locks := smbstatusreader.GetLockData(smbstatusout.LockData0Line)
	shares := smbstatusreader.GetShareData(smbstatusout.ShareData0Line)
	processes := smbstatusreader.GetProcessData(smbstatusout.ProcessData0Lines)

	ret := GetSmbStatistics(locks, processes, shares)
	if len(ret) != 8 {
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

	if len(ret) != 14 {
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

	if ret[9].Name != "server_information" {
		t.Errorf("The Name \"%s\" is not expected", ret[6].Name)
	}

	if ret[9].Value != 1 {
		t.Errorf("The Value %f is not expected", ret[6].Value)
	}

	if len(ret[9].Labels) != 1 {
		t.Errorf("There are more labels than expected")
	}

	value, found := ret[9].Labels["version"]
	if !found {
		t.Errorf("No label with key \"version\" found")
	}

	if value != "4.11.6-Ubuntu" {
		t.Errorf("The SambaVersion \"%s\" is not expected", value)
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
