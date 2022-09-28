package statisticsGenerator

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"testing"

	"tobi.backfrak.de/internal/commonbl"
)

func TestGetSmbdMetricsNotRunningProcess(t *testing.T) {

	metrics := GetSmbdMetrics([]commonbl.PsUtilPidData{})

	if len(metrics) != 19 {
		t.Errorf("Got %d lines but expected %d", len(metrics), 7)
	}

	if metrics[0].Name != "smbd_unique_process_id_count" {
		t.Errorf("The metric at index '0' name '%s' is not expected", metrics[0].Name)
	}

	if metrics[0].Value != 0 {
		t.Errorf("Found '%f' processes, but 0 expected", metrics[0].Value)
	}

	if metricArrContainsItemWithName(metrics, "smbd_cpu_usage_percentage") == false {
		t.Errorf("Can not find a metric named 'smbd_cpu_usage_percentage'")
	}

	if metricArrContainsItemWithName(metrics, "smbd_sum_cpu_usage_percentage") == false {
		t.Errorf("Can not find a metric named 'smbd_sum_cpu_usage_percentage'")
	}

	if metricArrContainsItemWithName(metrics, "smbd_virtual_memory_usage_bytes") == false {
		t.Errorf("Can not find a metric named 'smbd_virtual_memory_usage_bytes'")
	}

	if metricArrContainsItemWithName(metrics, "smbd_sum_virtual_memory_usage_bytes") == false {
		t.Errorf("Can not find a metric named 'smbd_sum_virtual_memory_usage_bytes'")
	}

	if metricArrContainsItemWithName(metrics, "smbd_cpu_usage_percentage") == false {
		t.Errorf("Can not find a metric named 'smbd_cpu_usage_percentage'")
	}

	if metricArrContainsItemWithName(metrics, "smbd_sum_cpu_usage_percentage") == false {
		t.Errorf("Can not find a metric named 'smbd_sum_cpu_usage_percentage'")
	}

	if metricArrContainsItemWithName(metrics, "smbd_io_counter_read_count") == false {
		t.Errorf("Can not find a metric named 'smbd_io_counter_read_count'")
	}

	if metricArrContainsItemWithName(metrics, "smbd_sum_io_counter_read_count") == false {
		t.Errorf("Can not find a metric named 'smbd_sum_io_counter_read_count'")
	}

	if metricArrContainsItemWithName(metrics, "smbd_io_counter_write_count") == false {
		t.Errorf("Can not find a metric named 'smbd_io_counter_write_count'")
	}

	if metricArrContainsItemWithName(metrics, "smbd_sum_io_counter_write_count") == false {
		t.Errorf("Can not find a metric named 'smbd_sum_io_counter_write_count'")
	}

	if metricArrContainsItemWithName(metrics, "smbd_io_counter_read_bytes") == false {
		t.Errorf("Can not find a metric named 'smbd_io_counter_read_bytes'")
	}

	if metricArrContainsItemWithName(metrics, "smbd_sum_io_counter_read_bytes") == false {
		t.Errorf("Can not find a metric named 'smbd_sum_io_counter_read_bytes'")
	}

	if metricArrContainsItemWithName(metrics, "smbd_io_counter_write_bytes") == false {
		t.Errorf("Can not find a metric named 'smbd_io_counter_write_bytes'")
	}

	if metricArrContainsItemWithName(metrics, "smbd_sum_io_counter_write_bytes") == false {
		t.Errorf("Can not find a metric named 'smbd_sum_io_counter_write_bytes'")
	}

	if metricArrContainsItemWithName(metrics, "smbd_open_file_count") == false {
		t.Errorf("Can not find a metric named 'smbd_open_file_count'")
	}

	if metricArrContainsItemWithName(metrics, "smbd_sum_open_file_count") == false {
		t.Errorf("Can not find a metric named 'smbd_sum_open_file_count'")
	}

	if metricArrContainsItemWithName(metrics, "smbd_thread_count") == false {
		t.Errorf("Can not find a metric named 'smbd_thread_count'")
	}

	if metricArrContainsItemWithName(metrics, "smbd_sum_thread_count") == false {
		t.Errorf("Can not find a metric named 'smbd_sum_thread_count'")
	}
}

func TestGetSmbdMetricsRunningProcess(t *testing.T) {

	pidData := []commonbl.PsUtilPidData{}
	pidData = append(pidData, commonbl.PsUtilPidData{
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

	pidData = append(pidData, commonbl.PsUtilPidData{
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

	metrics := GetSmbdMetrics(pidData)

	if len(metrics) < 1 {
		t.Errorf("Got less then one metric")
	}

	if metrics[0].Name != "smbd_unique_process_id_count" {
		t.Errorf("The metric at index '0' name '%s' is not expected", metrics[0].Name)
	}

	if metrics[0].Value != 2 {
		t.Errorf("Found '0' processes, but at two expected")
	}

	numUnqueMetrics := 9
	numSumMetrics := numUnqueMetrics
	expectedMetricCount := 1 + (int(metrics[0].Value) * numUnqueMetrics) + numSumMetrics
	if len(metrics) != expectedMetricCount {
		t.Errorf("Got '%d' metrics but expected '%d'", len(metrics), expectedMetricCount)
	}

	if metricArrContainsItemWithName(metrics, "smbd_cpu_usage_percentage") == false {
		t.Errorf("Can not find a metric named 'smbd_cpu_usage_percentage'")
	}

	if metricArrContainsItemWithName(metrics, "smbd_sum_cpu_usage_percentage") == false {
		t.Errorf("Can not find a metric named 'smbd_sum_cpu_usage_percentage'")
	}

	if metricArrContainsItemWithName(metrics, "smbd_virtual_memory_usage_bytes") == false {
		t.Errorf("Can not find a metric named 'smbd_virtual_memory_usage_bytes'")
	}

	if metricArrContainsItemWithName(metrics, "smbd_sum_virtual_memory_usage_bytes") == false {
		t.Errorf("Can not find a metric named 'smbd_sum_virtual_memory_usage_bytes'")
	}

	if metricArrCountItemWithName(metrics, "smbd_virtual_memory_usage_bytes") != int(metrics[0].Value) {
		t.Errorf("The metric 'smbd_virtual_memory_usage_bytes' is not exported as often as expected")
	}

	if metricArrCountItemWithName(metrics, "smbd_cpu_usage_percentage") != int(metrics[0].Value) {
		t.Errorf("The metric 'smbd_cpu_usage_percentage' is not exported as often as expected")
	}

	if metricArrSumItemWithName(metrics, "smbd_virtual_memory_usage_bytes") !=
		metricArrGetValueithName(metrics, "smbd_sum_virtual_memory_usage_bytes") {

		t.Errorf("The metrics 'smbd_virtual_memory_usage_bytes' sum is not equal 'smbd_sum_virtual_memory_usage_bytes'")
	}

	if metricArrSumItemWithName(metrics, "smbd_cpu_usage_percentage") !=
		metricArrGetValueithName(metrics, "smbd_sum_cpu_usage_percentage") {

		t.Errorf("The metrics 'smbd_cpu_usage_percentage' sum is not equal 'smbd_sum_cpu_usage_percentage'")
	}

	if metricArrContainsItemWithName(metrics, "smbd_virtual_memory_usage_percent") == false {
		t.Errorf("Can not find a metric named 'smbd_virtual_memory_usage_percent'")
	}

	if metricArrContainsItemWithName(metrics, "smbd_sum_virtual_memory_usage_percent") == false {
		t.Errorf("Can not find a metric named 'smbd_sum_virtual_memory_usage_percent'")
	}

	if metricArrCountItemWithName(metrics, "smbd_virtual_memory_usage_percent") != int(metrics[0].Value) {
		t.Errorf("The metric 'smbd_virtual_memory_usage_percent' is not exported as often as expected")
	}

	if metricArrSumItemWithName(metrics, "smbd_virtual_memory_usage_percent") !=
		metricArrGetValueithName(metrics, "smbd_sum_virtual_memory_usage_percent") {

		t.Errorf("The metrics 'smbd_virtual_memory_usage_percent' sum is not equal 'smbd_sum_virtual_memory_usage_percent'")
	}

	if metricArrContainsItemWithName(metrics, "smbd_io_counter_read_count") == false {
		t.Errorf("Can not find a metric named 'smbd_io_counter_read_count'")
	}

	if metricArrContainsItemWithName(metrics, "smbd_sum_io_counter_read_count") == false {
		t.Errorf("Can not find a metric named 'smbd_sum_io_counter_read_count'")
	}

	if metricArrCountItemWithName(metrics, "smbd_io_counter_read_count") != int(metrics[0].Value) {
		t.Errorf("The metric 'smbd_io_counter_read_count' is not exported as often as expected")
	}

	if metricArrSumItemWithName(metrics, "smbd_io_counter_read_count") !=
		metricArrGetValueithName(metrics, "smbd_sum_io_counter_read_count") {

		t.Errorf("The metrics 'smbd_io_counter_read_count' (%f) sum is not equal 'smbd_sum_io_counter_read_count' (%f)",
			metricArrSumItemWithName(metrics, "smbd_io_counter_read_count"),
			metricArrGetValueithName(metrics, "smbd_sum_io_counter_read_count"))
	}

	if metricArrContainsItemWithName(metrics, "smbd_io_counter_write_count") == false {
		t.Errorf("Can not find a metric named 'smbd_io_counter_write_count'")
	}

	if metricArrContainsItemWithName(metrics, "smbd_sum_io_counter_write_count") == false {
		t.Errorf("Can not find a metric named 'smbd_sum_io_counter_write_count'")
	}

	if metricArrCountItemWithName(metrics, "smbd_io_counter_write_count") != int(metrics[0].Value) {
		t.Errorf("The metric 'smbd_io_counter_write_count' is not exported as often as expected")
	}

	if metricArrSumItemWithName(metrics, "smbd_io_counter_write_count") !=
		metricArrGetValueithName(metrics, "smbd_sum_io_counter_write_count") {

		t.Errorf("The metrics 'smbd_io_counter_write_count' (%f) sum is not equal 'smbd_sum_io_counter_write_count' (%f)",
			metricArrSumItemWithName(metrics, "smbd_io_counter_write_count"),
			metricArrGetValueithName(metrics, "smbd_sum_io_counter_write_count"))
	}

	if metricArrContainsItemWithName(metrics, "smbd_io_counter_read_bytes") == false {
		t.Errorf("Can not find a metric named 'smbd_io_counter_read_bytes'")
	}

	if metricArrContainsItemWithName(metrics, "smbd_sum_io_counter_read_bytes") == false {
		t.Errorf("Can not find a metric named 'smbd_sum_io_counter_read_bytes'")
	}

	if metricArrCountItemWithName(metrics, "smbd_io_counter_read_bytes") != int(metrics[0].Value) {
		t.Errorf("The metric 'smbd_io_counter_read_bytes' is not exported as often as expected")
	}

	if metricArrSumItemWithName(metrics, "smbd_io_counter_read_bytes") !=
		metricArrGetValueithName(metrics, "smbd_sum_io_counter_read_bytes") {

		t.Errorf("The metrics 'smbd_io_counter_read_bytes' (%f) sum is not equal 'smbd_sum_io_counter_read_bytes' (%f)",
			metricArrSumItemWithName(metrics, "smbd_io_counter_read_bytes"),
			metricArrGetValueithName(metrics, "smbd_sum_io_counter_read_bytes"))
	}

	if metricArrContainsItemWithName(metrics, "smbd_io_counter_write_bytes") == false {
		t.Errorf("Can not find a metric named 'smbd_io_counter_write_bytes'")
	}

	if metricArrContainsItemWithName(metrics, "smbd_sum_io_counter_write_bytes") == false {
		t.Errorf("Can not find a metric named 'smbd_sum_io_counter_write_bytes'")
	}

	if metricArrCountItemWithName(metrics, "smbd_io_counter_write_bytes") != int(metrics[0].Value) {
		t.Errorf("The metric 'smbd_io_counter_write_bytes' is not exported as often as expected")
	}

	if metricArrSumItemWithName(metrics, "smbd_io_counter_write_bytes") !=
		metricArrGetValueithName(metrics, "smbd_sum_io_counter_write_bytes") {

		t.Errorf("The metrics 'smbd_io_counter_write_bytes' (%f) sum is not equal 'smbd_sum_io_counter_write_bytes' (%f)",
			metricArrSumItemWithName(metrics, "smbd_io_counter_write_bytes"),
			metricArrGetValueithName(metrics, "smbd_sum_io_counter_write_bytes"))
	}

	if metricArrContainsItemWithName(metrics, "smbd_open_file_count") == false {
		t.Errorf("Can not find a metric named 'smbd_open_file_count'")
	}

	if metricArrContainsItemWithName(metrics, "smbd_sum_open_file_count") == false {
		t.Errorf("Can not find a metric named 'smbd_sum_open_file_count'")
	}

	if metricArrCountItemWithName(metrics, "smbd_open_file_count") != int(metrics[0].Value) {
		t.Errorf("The metric 'smbd_open_file_count' is not exported as often as expected")
	}

	if metricArrSumItemWithName(metrics, "smbd_open_file_count") !=
		metricArrGetValueithName(metrics, "smbd_sum_open_file_count") {

		t.Errorf("The metrics 'smbd_open_file_count' (%f) sum is not equal 'smbd_sum_open_file_count' (%f)",
			metricArrSumItemWithName(metrics, "smbd_open_file_count"),
			metricArrGetValueithName(metrics, "smbd_sum_open_file_count"))
	}

	if metricArrContainsItemWithName(metrics, "smbd_thread_count") == false {
		t.Errorf("Can not find a metric named 'smbd_thread_count'")
	}

	if metricArrContainsItemWithName(metrics, "smbd_sum_thread_count") == false {
		t.Errorf("Can not find a metric named 'smbd_sum_thread_count'")
	}

	if metricArrCountItemWithName(metrics, "smbd_thread_count") != int(metrics[0].Value) {
		t.Errorf("The metric 'smbd_thread_count' is not exported as often as expected")
	}

	if metricArrSumItemWithName(metrics, "smbd_thread_count") !=
		metricArrGetValueithName(metrics, "smbd_sum_thread_count") {

		t.Errorf("The metrics 'smbd_thread_count' (%f) sum is not equal 'smbd_sum_thread_count' (%f)",
			metricArrSumItemWithName(metrics, "smbd_thread_count"),
			metricArrGetValueithName(metrics, "smbd_sum_thread_count"))
	}

}

func metricArrContainsItemWithName(arr []SmbStatisticsNumeric, name string) bool {
	for _, item := range arr {
		if item.Name == name {
			return true
		}
	}

	return false
}

func metricArrCountItemWithName(arr []SmbStatisticsNumeric, name string) int {
	ret := 0
	for _, item := range arr {
		if item.Name == name {
			ret++
		}
	}

	return ret
}

func metricArrSumItemWithName(arr []SmbStatisticsNumeric, name string) float64 {
	ret := float64(0)
	for _, item := range arr {
		if item.Name == name {
			ret += item.Value
		}
	}

	return ret
}

func metricArrGetValueithName(arr []SmbStatisticsNumeric, name string) float64 {
	ret := float64(0)
	for _, item := range arr {
		if item.Name == name {
			return item.Value
		}
	}

	return ret
}
