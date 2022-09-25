package statisticsGenerator

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

var smbd_image_name = "smbd"

func GetSmbdMetrics() ([]SmbStatisticsNumeric, error) {

	var pidList []int32

	pgrepPath, errLookPath := exec.LookPath("pgrep")
	if errLookPath != nil {
		return nil, errLookPath
	}
	pidListInByte, errGetPids := exec.Command(pgrepPath, smbd_image_name).Output()
	if errGetPids != nil {
		return getMetricsFromPidList([]int32{}), nil
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

	return getMetricsFromPidList(pidList), nil
}

func getMetricsFromPidList(pidList []int32) []SmbStatisticsNumeric {
	var ret []SmbStatisticsNumeric
	ret = append(ret, SmbStatisticsNumeric{"smbd_unique_process_id_count", float64(len(pidList)), fmt.Sprintf("Count of unique process IDs for '%s'", smbd_image_name), nil})

	return ret
}
