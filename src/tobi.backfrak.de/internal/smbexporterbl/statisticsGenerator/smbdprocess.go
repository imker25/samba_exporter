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

	"github.com/shirou/gopsutil/v3/process"
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
		return getMetricsFromPidList([]int32{})
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

	return getMetricsFromPidList(pidList)
}

func getMetricsFromPidList(pidList []int32) ([]SmbStatisticsNumeric, error) {
	var ret []SmbStatisticsNumeric

	ret = append(ret, SmbStatisticsNumeric{"smbd_unique_process_id_count", float64(len(pidList)), fmt.Sprintf("Count of unique process IDs for '%s'", smbd_image_name), nil})

	if len(pidList) > 0 {
		cpuPercentageSum := float64(0)
		vmBytesSum := uint64(0)
		for _, pid := range pidList {
			proc, errProc := process.NewProcess(pid)
			if errProc != nil {
				return nil, errProc
			}

			cpuPercent, errPer := proc.CPUPercent()
			if errPer != nil {
				return nil, errPer
			}
			ret = append(ret, SmbStatisticsNumeric{"smbd_cpu_usage_percentage",
				cpuPercent, fmt.Sprintf("CPU usage of the '%s' process with pid in percent", smbd_image_name),
				map[string]string{"pid": strconv.Itoa(int(pid))}})
			cpuPercentageSum += cpuPercent

			vmBytes, errVmBytes := proc.MemoryInfo()
			if errVmBytes != nil {
				return nil, errVmBytes
			}
			ret = append(ret, SmbStatisticsNumeric{"smbd_virtual_memory_usage_bytes",
				float64(vmBytes.VMS), fmt.Sprintf("Virtual memory usage of the '%s' process with pid in bytes", smbd_image_name),
				map[string]string{"pid": strconv.Itoa(int(pid))}})
			vmBytesSum += vmBytes.VMS

		}

		// Add sum metrics (without label)
		ret = append(ret, SmbStatisticsNumeric{"smbd_sum_cpu_usage_percentage",
			cpuPercentageSum, fmt.Sprintf("Sum CPU usage of all '%s' processes in percent", smbd_image_name), nil})
		ret = append(ret, SmbStatisticsNumeric{"smbd_sum_virtual_memory_usage_bytes",
			float64(vmBytesSum), fmt.Sprintf("Virtual memory usage of all '%s' processes in bytes", smbd_image_name), nil})

	} else {
		// Give back empty metrics, when smbd is not running

		// Metrics with labels
		ret = append(ret, SmbStatisticsNumeric{"smbd_cpu_usage_percentage",
			0, fmt.Sprintf("CPU usage of the '%s' process with pid in percent", smbd_image_name),
			map[string]string{"pid": ""}})
		ret = append(ret, SmbStatisticsNumeric{"smbd_virtual_memory_usage_bytes",
			0, fmt.Sprintf("Virtual memory usage of the '%s' process with pid in bytes", smbd_image_name),
			map[string]string{"pid": ""}})

		// Metrics without labels (sum metrics)
		ret = append(ret, SmbStatisticsNumeric{"smbd_sum_cpu_usage_percentage",
			0, fmt.Sprintf("Sum CPU usage of all '%s' processes in percent", smbd_image_name), nil})
		ret = append(ret, SmbStatisticsNumeric{"smbd_sum_virtual_memory_usage_bytes",
			0, fmt.Sprintf("Virtual memory usage of all '%s' processes in bytes", smbd_image_name), nil})
	}

	return ret, nil
}
