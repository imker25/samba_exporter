package statisticsGenerator

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"testing"
)

func TestGetSmbdMetricsNotRunningProcess(t *testing.T) {

	if smbd_image_name != "smbd" {
		t.Errorf("The variable 'smbd_image_name' has the value '%s' but 'smbd' is expected", smbd_image_name)
	}

	smbd_image_name = "not_existing_process"
	metrics, err := GetSmbdMetrics()
	if err != nil {
		t.Errorf("Got the unexpected error: %s", err.Error())
	}

	if len(metrics) != 5 {
		t.Errorf("Got %d lines but expected %d", len(metrics), 3)
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

	smbd_image_name = "smbd"
}

func TestGetSmbdMetricsGoProcess(t *testing.T) {

	if smbd_image_name != "smbd" {
		t.Errorf("The variable 'smbd_image_name' has the value '%s' but 'smbd' is expected", smbd_image_name)
	}

	smbd_image_name = "go"
	metrics, err := GetSmbdMetrics()
	if err != nil {
		t.Errorf("Got the unexpected error: %s", err.Error())
	}

	if len(metrics) < 1 {
		t.Errorf("Got less then one metric")
	}

	if metrics[0].Name != "smbd_unique_process_id_count" {
		t.Errorf("The metric at index '0' name '%s' is not expected", metrics[0].Name)
	}

	if metrics[0].Value < 1 {
		t.Errorf("Found '0' processes, but at least one expected")
	}

	numUnqueMetrics := 2
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

	smbd_image_name = "smbd"
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
