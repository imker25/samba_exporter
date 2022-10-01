package statisticsGenerator

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"fmt"
	"strconv"

	"tobi.backfrak.de/internal/commonbl"
)

const smbd_image_name = "smbd"

// Get the SmbStatisticsNumeric metrics out of a list of []commonbl.PsUtilPidData)
func GetSmbdMetrics(pidDataList []commonbl.PsUtilPidData, notExportPid bool) []SmbStatisticsNumeric {

	var ret []SmbStatisticsNumeric

	ret = append(ret, SmbStatisticsNumeric{"smbd_unique_process_id_count", float64(len(pidDataList)), fmt.Sprintf("Count of unique process IDs for '%s'", smbd_image_name), nil})

	if len(pidDataList) > 0 {
		cpuPercentageSum := float64(0)
		vmBytesSum := uint64(0)
		vmPercentSum := float64(0)
		readCountSum := uint64(0)
		writeCountSum := uint64(0)
		readBytesSum := uint64(0)
		writeBytesSum := uint64(0)
		openFilesCountSum := uint64(0)
		threadCountSum := uint64(0)
		for _, pidData := range pidDataList {

			cpuPercentageSum += pidData.CpuUsagePercent
			vmBytesSum += pidData.VirtualMemoryUsageBytes
			vmPercentSum += pidData.VirtualMemoryUsagePercent
			readCountSum += pidData.IoCounterReadCount
			writeCountSum += pidData.IoCounterWriteCount
			readBytesSum += pidData.IoCounterReadBytes
			writeBytesSum += pidData.IoCounterWriteBytes
			openFilesCountSum += pidData.OpenFilesCount
			threadCountSum += pidData.ThreadCount

			if !notExportPid {
				// Metrics with PID label
				ret = append(ret, SmbStatisticsNumeric{"smbd_cpu_usage_percentage",
					pidData.CpuUsagePercent, fmt.Sprintf("CPU usage of the '%s' process with pid in percent", smbd_image_name),
					map[string]string{"pid": strconv.Itoa(int(pidData.PID))}})
				ret = append(ret, SmbStatisticsNumeric{"smbd_virtual_memory_usage_bytes",
					float64(pidData.VirtualMemoryUsageBytes), fmt.Sprintf("Virtual memory usage of the '%s' process with pid in bytes", smbd_image_name),
					map[string]string{"pid": strconv.Itoa(int(pidData.PID))}})
				ret = append(ret, SmbStatisticsNumeric{"smbd_virtual_memory_usage_percent",
					pidData.VirtualMemoryUsagePercent, fmt.Sprintf("Virtual memory usage of the '%s' process with pid in percent", smbd_image_name),
					map[string]string{"pid": strconv.Itoa(int(pidData.PID))}})
				ret = append(ret, SmbStatisticsNumeric{"smbd_io_counter_read_count",
					float64(pidData.IoCounterReadCount), fmt.Sprintf("IO counter read count of the process '%s'", smbd_image_name),
					map[string]string{"pid": strconv.Itoa(int(pidData.PID))}})
				ret = append(ret, SmbStatisticsNumeric{"smbd_io_counter_write_count",
					float64(pidData.IoCounterWriteCount), fmt.Sprintf("IO counter write count of the process '%s'", smbd_image_name),
					map[string]string{"pid": strconv.Itoa(int(pidData.PID))}})
				ret = append(ret, SmbStatisticsNumeric{"smbd_io_counter_read_bytes",
					float64(pidData.IoCounterReadBytes), fmt.Sprintf("IO counter reads of the process '%s' in byte", smbd_image_name),
					map[string]string{"pid": strconv.Itoa(int(pidData.PID))}})
				ret = append(ret, SmbStatisticsNumeric{"smbd_io_counter_write_bytes",
					float64(pidData.IoCounterWriteBytes), fmt.Sprintf("IO counter writes of the process '%s' in byte", smbd_image_name),
					map[string]string{"pid": strconv.Itoa(int(pidData.PID))}})
				ret = append(ret, SmbStatisticsNumeric{"smbd_open_file_count",
					float64(pidData.OpenFilesCount), fmt.Sprintf("Open file handles by process '%s'", smbd_image_name),
					map[string]string{"pid": strconv.Itoa(int(pidData.PID))}})
				ret = append(ret, SmbStatisticsNumeric{"smbd_thread_count",
					float64(pidData.ThreadCount), fmt.Sprintf("Threads used by process '%s'", smbd_image_name),
					map[string]string{"pid": strconv.Itoa(int(pidData.PID))}})
			}
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
		if !notExportPid {
			// Metrics with PID labels
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
		}

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

	return ret
}
