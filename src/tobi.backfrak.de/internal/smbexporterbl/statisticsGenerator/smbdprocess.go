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
		vmPercentSum := float64(0)
		readCountSum := uint64(0)
		writeCountSum := uint64(0)
		readBytesSum := uint64(0)
		writeBytesSum := uint64(0)
		openFilesCountSum := 0
		threadCountSum := 0
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

			vmPercent, errVmPercent := proc.MemoryPercent()
			if errVmPercent != nil {
				return nil, errVmPercent
			}
			ret = append(ret, SmbStatisticsNumeric{"smbd_virtual_memory_usage_percent",
				float64(vmPercent), fmt.Sprintf("Virtual memory usage of the '%s' process with pid in percent", smbd_image_name),
				map[string]string{"pid": strconv.Itoa(int(pid))}})
			vmPercentSum += float64(vmPercent)

			ioCounters, errIoCounters := proc.IOCounters()
			if errIoCounters != nil {
				return nil, errIoCounters
			}
			ret = append(ret, SmbStatisticsNumeric{"smbd_io_counter_read_count",
				float64(ioCounters.ReadCount), fmt.Sprintf("IO counter read count of the process '%s'", smbd_image_name),
				map[string]string{"pid": strconv.Itoa(int(pid))}})
			readCountSum += ioCounters.ReadCount
			ret = append(ret, SmbStatisticsNumeric{"smbd_io_counter_write_count",
				float64(ioCounters.WriteCount), fmt.Sprintf("IO counter write count of the process '%s'", smbd_image_name),
				map[string]string{"pid": strconv.Itoa(int(pid))}})
			writeCountSum += ioCounters.WriteCount
			ret = append(ret, SmbStatisticsNumeric{"smbd_io_counter_read_bytes",
				float64(ioCounters.ReadBytes), fmt.Sprintf("IO counter reads of the process '%s' in byte", smbd_image_name),
				map[string]string{"pid": strconv.Itoa(int(pid))}})
			readBytesSum += ioCounters.ReadBytes
			ret = append(ret, SmbStatisticsNumeric{"smbd_io_counter_write_bytes",
				float64(ioCounters.WriteBytes), fmt.Sprintf("IO counter writes of the process '%s' in byte", smbd_image_name),
				map[string]string{"pid": strconv.Itoa(int(pid))}})
			writeBytesSum += ioCounters.WriteBytes

			openFileStats, errOpenFileStats := proc.OpenFiles()
			if errOpenFileStats != nil {
				return nil, errOpenFileStats
			}
			openFilesCount := len(openFileStats)
			ret = append(ret, SmbStatisticsNumeric{"smbd_open_file_count",
				float64(openFilesCount), fmt.Sprintf("Open file handles by process '%s'", smbd_image_name),
				map[string]string{"pid": strconv.Itoa(int(pid))}})
			openFilesCountSum += openFilesCount

			threadStats, errThreadStats := proc.Threads()
			if errThreadStats != nil {
				return nil, errThreadStats
			}
			threadCount := len(threadStats)
			ret = append(ret, SmbStatisticsNumeric{"smbd_thread_count",
				float64(threadCount), fmt.Sprintf("Threads used by process '%s'", smbd_image_name),
				map[string]string{"pid": strconv.Itoa(int(pid))}})
			threadCountSum += threadCount
		}

		// Add sum metrics (without label)
		ret = append(ret, SmbStatisticsNumeric{"smbd_sum_cpu_usage_percentage",
			cpuPercentageSum, fmt.Sprintf("Sum CPU usage of all '%s' processes in percent", smbd_image_name), nil})
		ret = append(ret, SmbStatisticsNumeric{"smbd_sum_virtual_memory_usage_bytes",
			float64(vmBytesSum), fmt.Sprintf("Virtual memory usage of all '%s' processes in bytes", smbd_image_name), nil})
		ret = append(ret, SmbStatisticsNumeric{"smbd_sum_virtual_memory_usage_percent",
			vmPercentSum, fmt.Sprintf("Virtual memory usage of all '%s' processes in percent", smbd_image_name), nil})
		ret = append(ret, SmbStatisticsNumeric{"smbd_sum_io_counter_read_count",
			float64(readCountSum), fmt.Sprintf("IO counter read count of all '%s' processes", smbd_image_name), nil})
		ret = append(ret, SmbStatisticsNumeric{"smbd_sum_io_counter_write_count",
			float64(writeCountSum), fmt.Sprintf("IO counter write count of all '%s' processes", smbd_image_name), nil})
		ret = append(ret, SmbStatisticsNumeric{"smbd_sum_io_counter_read_bytes",
			float64(readBytesSum), fmt.Sprintf("IO counter reads of all '%s' processes in bytes", smbd_image_name), nil})
		ret = append(ret, SmbStatisticsNumeric{"smbd_sum_io_counter_write_bytes",
			float64(writeBytesSum), fmt.Sprintf("IO counter writes of all '%s' processes in bytes", smbd_image_name), nil})
		ret = append(ret, SmbStatisticsNumeric{"smbd_sum_open_file_count",
			float64(openFilesCountSum), fmt.Sprintf("Open file handles of all '%s' processes", smbd_image_name), nil})
		ret = append(ret, SmbStatisticsNumeric{"smbd_sum_thread_count",
			float64(threadCountSum), fmt.Sprintf("Threads used by all '%s' processes", smbd_image_name), nil})

	} else {
		// Give back empty metrics, when smbd is not running

		// Metrics with labels
		ret = append(ret, SmbStatisticsNumeric{"smbd_cpu_usage_percentage",
			0, fmt.Sprintf("CPU usage of the '%s' process with pid in percent", smbd_image_name),
			map[string]string{"pid": ""}})
		ret = append(ret, SmbStatisticsNumeric{"smbd_virtual_memory_usage_bytes",
			0, fmt.Sprintf("Virtual memory usage of the '%s' process with pid in bytes", smbd_image_name),
			map[string]string{"pid": ""}})
		ret = append(ret, SmbStatisticsNumeric{"smbd_virtual_memory_usage_percent",
			0, fmt.Sprintf("Virtual memory usage of the '%s' process with pid in percent", smbd_image_name),
			map[string]string{"pid": ""}})
		ret = append(ret, SmbStatisticsNumeric{"smbd_io_counter_read_count",
			0, fmt.Sprintf("IO counter read count of the process '%s'", smbd_image_name),
			map[string]string{"pid": ""}})
		ret = append(ret, SmbStatisticsNumeric{"smbd_io_counter_write_count",
			0, fmt.Sprintf("IO counter write count of the process '%s'", smbd_image_name),
			map[string]string{"pid": ""}})
		ret = append(ret, SmbStatisticsNumeric{"smbd_io_counter_read_bytes",
			0, fmt.Sprintf("IO counter reads of the process '%s' in byte", smbd_image_name),
			map[string]string{"pid": ""}})
		ret = append(ret, SmbStatisticsNumeric{"smbd_io_counter_write_bytes",
			0, fmt.Sprintf("IO counter writes of the process '%s' in byte", smbd_image_name),
			map[string]string{"pid": ""}})
		ret = append(ret, SmbStatisticsNumeric{"smbd_open_file_count",
			0, fmt.Sprintf("Open file handles by process '%s'", smbd_image_name),
			map[string]string{"pid": ""}})
		ret = append(ret, SmbStatisticsNumeric{"smbd_thread_count",
			0, fmt.Sprintf("Threads used by process '%s'", smbd_image_name),
			map[string]string{"pid": ""}})

		// Metrics without labels (sum metrics)
		ret = append(ret, SmbStatisticsNumeric{"smbd_sum_cpu_usage_percentage",
			0, fmt.Sprintf("Sum CPU usage of all '%s' processes in percent", smbd_image_name), nil})
		ret = append(ret, SmbStatisticsNumeric{"smbd_sum_virtual_memory_usage_bytes",
			0, fmt.Sprintf("Virtual memory usage of all '%s' processes in bytes", smbd_image_name), nil})
		ret = append(ret, SmbStatisticsNumeric{"smbd_sum_virtual_memory_usage_percent",
			0, fmt.Sprintf("Virtual memory usage of all '%s' processes in percent", smbd_image_name), nil})
		ret = append(ret, SmbStatisticsNumeric{"smbd_sum_io_counter_read_count",
			0, fmt.Sprintf("IO counter read count of all '%s' processes", smbd_image_name), nil})
		ret = append(ret, SmbStatisticsNumeric{"smbd_sum_io_counter_write_count",
			0, fmt.Sprintf("IO counter write count of all '%s' processes", smbd_image_name), nil})
		ret = append(ret, SmbStatisticsNumeric{"smbd_sum_io_counter_read_bytes",
			0, fmt.Sprintf("IO counter reads of all '%s' processes in bytes", smbd_image_name), nil})
		ret = append(ret, SmbStatisticsNumeric{"smbd_sum_io_counter_write_bytes",
			0, fmt.Sprintf("IO counter writes of all '%s' processes in bytes", smbd_image_name), nil})
		ret = append(ret, SmbStatisticsNumeric{"smbd_sum_open_file_count",
			0, fmt.Sprintf("Open file handles of all '%s' processes", smbd_image_name), nil})
		ret = append(ret, SmbStatisticsNumeric{"smbd_sum_thread_count",
			0, fmt.Sprintf("Threads used by all '%s' processes", smbd_image_name), nil})
	}

	return ret, nil
}
