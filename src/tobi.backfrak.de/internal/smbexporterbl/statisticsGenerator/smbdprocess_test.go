package statisticsGenerator

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import "testing"

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
	smbd_image_name = "smbd"
}
