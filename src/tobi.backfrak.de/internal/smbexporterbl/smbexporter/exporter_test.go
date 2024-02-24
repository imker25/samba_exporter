package smbexporter

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"fmt"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"tobi.backfrak.de/internal/commonbl"
	"tobi.backfrak.de/internal/smbexporterbl/smbstatusreader"
	"tobi.backfrak.de/internal/smbexporterbl/statisticsGenerator"
	"tobi.backfrak.de/internal/smbstatusout"
	"tobi.backfrak.de/internal/testhelper"
)

func getNewStatisticGenSettings() statisticsGenerator.StatisticsGeneratorSettings {
	return statisticsGenerator.StatisticsGeneratorSettings{}
}

func TestNewSambaExporter(t *testing.T) {
	requestHandler := *commonbl.NewPipeHandler(true, commonbl.RequestPipe)
	responseHandler := *commonbl.NewPipeHandler(true, commonbl.ResposePipe)
	logger := *testhelper.NewTestLogger(true)
	exporter := NewSambaExporter(&requestHandler, &responseHandler, &logger, "0.0.0", 5, getNewStatisticGenSettings())

	if exporter.RequestHandler.PipeType != commonbl.RequestPipe {
		t.Errorf("The exporter.RequestHandler is not of the expected type")
	}

	if exporter.ResponseHander.PipeType != commonbl.ResposePipe {
		t.Errorf("The exporter.RequestHandler is not of the expected type")
	}

	if exporter.descriptions == nil {
		t.Errorf("exporter.Descriptions are nil")
	}

	if logger.Verbose != exporter.Logger.GetVerbose() {
		t.Errorf("The exporter uses the wrong logger")
	}

	if exporter.Version != "0.0.0" {
		t.Errorf("The Version \"%s\" is not expected", exporter.Version)
	}

	if logger.GetOutputCount() != 0 {
		t.Errorf("The OutputCount '%d' is not the expected '0'", logger.GetOutputCount())
	}
}

func TestSetDescriptionsFromResponse(t *testing.T) {
	expectedChanels := 38
	requestHandler := *commonbl.NewPipeHandler(true, commonbl.RequestPipe)
	responseHandler := *commonbl.NewPipeHandler(true, commonbl.ResposePipe)
	logger := *testhelper.NewTestLogger(true)
	locks := smbstatusreader.GetLockData(smbstatusout.LockDataNoData, &logger)
	shares := smbstatusreader.GetShareData(smbstatusout.ShareDataOneLine, &logger)
	processes := smbstatusreader.GetProcessData(smbstatusout.ProcessDataOneLine, &logger)
	psData := smbstatusreader.GetPsData(commonbl.TestPsResponseEmpty(), &logger)
	ch := make(chan *prometheus.Desc, expectedChanels)
	exporter := NewSambaExporter(&requestHandler, &responseHandler, &logger, "0.0.0", 5, getNewStatisticGenSettings())
	exporter.setDescriptionsFromResponse(locks, processes, shares, psData, ch)

	if len(ch) != expectedChanels {
		t.Errorf("The number of descriptions is not expected")
	}

	for i := 0; i < expectedChanels; i++ {
		desc := <-ch
		if desc == nil {
			t.Errorf("Got a nil description for a metric")
		}
	}

	if logger.GetErrorCount() != 0 {
		t.Errorf("The ErrorCount '%d' is not the expected '0'", logger.GetErrorCount())
	}
}

func TestSetMetricsFromResponse(t *testing.T) {
	expectedDescChanels := 38
	expectedMetChanels := 65
	requestHandler := commonbl.NewPipeHandler(true, commonbl.RequestPipe)
	responseHandler := commonbl.NewPipeHandler(true, commonbl.ResposePipe)
	logger := testhelper.NewTestLogger(true)
	locks := smbstatusreader.GetLockData(smbstatusout.LockData4Lines, logger)
	shares := smbstatusreader.GetShareData(smbstatusout.ShareData4Lines, logger)
	processes := smbstatusreader.GetProcessData(smbstatusout.ProcessData4Lines, logger)
	psData := smbstatusreader.GetPsData(commonbl.TestPsResponse(), logger)
	chDesc := make(chan *prometheus.Desc, expectedDescChanels)
	exporter := NewSambaExporter(requestHandler, responseHandler, logger, "0.0.0", 5, getNewStatisticGenSettings())
	exporter.setDescriptionsFromResponse(locks, processes, shares, psData, chDesc)
	chMet := make(chan prometheus.Metric, expectedMetChanels)
	exporter.setMetricsFromResponse(locks, processes, shares, psData, 1, 1, 31, chMet)

	if len(chMet) != expectedMetChanels {
		t.Errorf("Got %d metric channels, but expected %d", len(chMet), expectedMetChanels)
	}

	for i := 0; i < expectedMetChanels; i++ {
		metric := <-chMet
		desc := metric.Desc()
		if desc == nil {
			t.Errorf("Got a nil description for a metric")
		}
	}

	if logger.GetErrorCount() != 0 {
		t.Errorf("The ErrorCount '%d' is not the expected '0'", logger.GetErrorCount())
	}
}

func TestSetMetricsFromResponseNameWithSpaces(t *testing.T) {
	expectedDescChanels := 38
	expectedMetChanels := 61
	requestHandler := commonbl.NewPipeHandler(true, commonbl.RequestPipe)
	responseHandler := commonbl.NewPipeHandler(true, commonbl.ResposePipe)
	logger := testhelper.NewTestLogger(true)
	locks := smbstatusreader.GetLockData(smbstatusout.LockData4Lines, logger)
	shares := smbstatusreader.GetShareData(smbstatusout.ShareData4LinesWithSpacesInName, logger)
	processes := smbstatusreader.GetProcessData(smbstatusout.ProcessData4Lines, logger)
	psData := smbstatusreader.GetPsData(commonbl.TestPsResponse(), logger)
	chDesc := make(chan *prometheus.Desc, expectedDescChanels)
	exporter := NewSambaExporter(requestHandler, responseHandler, logger, "0.0.0", 5, getNewStatisticGenSettings())
	exporter.setDescriptionsFromResponse(locks, processes, shares, psData, chDesc)
	chMet := make(chan prometheus.Metric, expectedMetChanels)
	exporter.setMetricsFromResponse(locks, processes, shares, psData, 1, 1, 31, chMet)

	if len(chMet) != expectedMetChanels {
		t.Errorf("Got %d metric channels, but expected %d", len(chMet), expectedMetChanels)
	}

	var metrics []prometheus.Metric
	for i := 0; i < expectedMetChanels; i++ {
		metric := <-chMet
		desc := metric.Desc()
		if desc == nil {
			t.Errorf("Got a nil description for a metric")
		}
		metrics = append(metrics, metric)
	}

	if len(metrics) != expectedMetChanels {
		t.Errorf("Got '%d' metrics but expected '%d'", len(metrics), expectedMetChanels)
	}

	if logger.GetErrorCount() != 0 {
		t.Errorf("The ErrorCount '%d' is not the expected '0'", logger.GetErrorCount())
	}
}

func TestSetMetricsFromResponseNoPid(t *testing.T) {
	exportSettings := statisticsGenerator.StatisticsGeneratorSettings{false, false, false, true, false}
	expectedDescChanels := 38
	expectedMetChanels := 47
	requestHandler := commonbl.NewPipeHandler(true, commonbl.RequestPipe)
	responseHandler := commonbl.NewPipeHandler(true, commonbl.ResposePipe)
	logger := testhelper.NewTestLogger(true)
	locks := smbstatusreader.GetLockData(smbstatusout.LockData4Lines, logger)
	shares := smbstatusreader.GetShareData(smbstatusout.ShareData4Lines, logger)
	processes := smbstatusreader.GetProcessData(smbstatusout.ProcessData4Lines, logger)
	psData := smbstatusreader.GetPsData(commonbl.TestPsResponse(), logger)
	chDesc := make(chan *prometheus.Desc, expectedDescChanels)
	exporter := NewSambaExporter(requestHandler, responseHandler, logger, "0.0.0", 5, exportSettings)
	exporter.setDescriptionsFromResponse(locks, processes, shares, psData, chDesc)
	chMet := make(chan prometheus.Metric, expectedMetChanels)
	exporter.setMetricsFromResponse(locks, processes, shares, psData, 1, 1, 31, chMet)

	if len(chMet) != expectedMetChanels {
		t.Errorf("Got %d metric channels, but expected %d", len(chMet), expectedMetChanels)
	}

	if logger.GetErrorCount() != 0 {
		t.Errorf("The ErrorCount '%d' is not the expected '0'", logger.GetErrorCount())
	}

}

func TestSetMetricsFromResponseNoUser(t *testing.T) {
	exportSettings := statisticsGenerator.StatisticsGeneratorSettings{false, true, false, false, false}
	expectedDescChanels := 38
	expectedMetChanels := 57
	requestHandler := commonbl.NewPipeHandler(true, commonbl.RequestPipe)
	responseHandler := commonbl.NewPipeHandler(true, commonbl.ResposePipe)
	logger := testhelper.NewTestLogger(true)
	locks := smbstatusreader.GetLockData(smbstatusout.LockData4Lines, logger)
	shares := smbstatusreader.GetShareData(smbstatusout.ShareData4Lines, logger)
	processes := smbstatusreader.GetProcessData(smbstatusout.ProcessData4Lines, logger)
	psData := smbstatusreader.GetPsData(commonbl.TestPsResponse(), logger)
	chDesc := make(chan *prometheus.Desc, expectedDescChanels)
	exporter := NewSambaExporter(requestHandler, responseHandler, logger, "0.0.0", 5, exportSettings)
	exporter.setDescriptionsFromResponse(locks, processes, shares, psData, chDesc)
	chMet := make(chan prometheus.Metric, expectedMetChanels)
	exporter.setMetricsFromResponse(locks, processes, shares, psData, 1, 1, 31, chMet)

	if len(chMet) != expectedMetChanels {
		t.Errorf("Got %d metric channels, but expected %d", len(chMet), expectedMetChanels)
	}

	if logger.GetErrorCount() != 0 {
		t.Errorf("The ErrorCount '%d' is not the expected '0'", logger.GetErrorCount())
	}

}

func TestSetMetricsFromResponseNoShareDetails(t *testing.T) {
	exportSettings := statisticsGenerator.StatisticsGeneratorSettings{false, false, false, false, true}
	expectedDescChanels := 38
	expectedMetChanels := 53
	requestHandler := commonbl.NewPipeHandler(true, commonbl.RequestPipe)
	responseHandler := commonbl.NewPipeHandler(true, commonbl.ResposePipe)
	logger := testhelper.NewTestLogger(true)
	locks := smbstatusreader.GetLockData(smbstatusout.LockData4Lines, logger)
	shares := smbstatusreader.GetShareData(smbstatusout.ShareData4Lines, logger)
	processes := smbstatusreader.GetProcessData(smbstatusout.ProcessData4Lines, logger)
	psData := smbstatusreader.GetPsData(commonbl.TestPsResponse(), logger)
	chDesc := make(chan *prometheus.Desc, expectedDescChanels)
	exporter := NewSambaExporter(requestHandler, responseHandler, logger, "0.0.0", 5, exportSettings)
	exporter.setDescriptionsFromResponse(locks, processes, shares, psData, chDesc)
	chMet := make(chan prometheus.Metric, expectedMetChanels)
	exporter.setMetricsFromResponse(locks, processes, shares, psData, 1, 1, 31, chMet)

	if len(chMet) != expectedMetChanels {
		t.Errorf("Got %d metric channels, but expected %d", len(chMet), expectedMetChanels)
	}

	if logger.GetErrorCount() != 0 {
		t.Errorf("The ErrorCount '%d' is not the expected '0'", logger.GetErrorCount())
	}

}

func TestSetMetricsFromResponseNoClient(t *testing.T) {
	exportSettings := statisticsGenerator.StatisticsGeneratorSettings{true, false, false, false, false}
	expectedDescChanels := 38
	expectedMetChanels := 53
	requestHandler := commonbl.NewPipeHandler(true, commonbl.RequestPipe)
	responseHandler := commonbl.NewPipeHandler(true, commonbl.ResposePipe)
	logger := testhelper.NewTestLogger(true)
	locks := smbstatusreader.GetLockData(smbstatusout.LockData4Lines, logger)
	shares := smbstatusreader.GetShareData(smbstatusout.ShareData4Lines, logger)
	processes := smbstatusreader.GetProcessData(smbstatusout.ProcessData4Lines, logger)
	psData := smbstatusreader.GetPsData(commonbl.TestPsResponse(), logger)
	chDesc := make(chan *prometheus.Desc, expectedDescChanels)
	exporter := NewSambaExporter(requestHandler, responseHandler, logger, "0.0.0", 5, exportSettings)
	exporter.setDescriptionsFromResponse(locks, processes, shares, psData, chDesc)
	chMet := make(chan prometheus.Metric, expectedMetChanels)
	exporter.setMetricsFromResponse(locks, processes, shares, psData, 1, 1, 31, chMet)

	if len(chMet) != expectedMetChanels {
		t.Errorf("Got %d metric channels, but expected %d", len(chMet), expectedMetChanels)
	}

	if logger.GetErrorCount() != 0 {
		t.Errorf("The ErrorCount '%d' is not the expected '0'", logger.GetErrorCount())
	}
}

func TestSetMetricsFromResponseCluster(t *testing.T) {
	exportSettings := statisticsGenerator.StatisticsGeneratorSettings{true, false, false, false, false}
	expectedDescChanels := 42
	expectedMetChanels := 53
	requestHandler := commonbl.NewPipeHandler(true, commonbl.RequestPipe)
	responseHandler := commonbl.NewPipeHandler(true, commonbl.ResposePipe)
	logger := testhelper.NewTestLogger(true)
	locks := smbstatusreader.GetLockData(smbstatusout.LockDataCluster, logger)
	shares := smbstatusreader.GetShareData(smbstatusout.ShareDataCluster, logger)
	processes := smbstatusreader.GetProcessData(smbstatusout.ProcessDataCluster, logger)
	psData := smbstatusreader.GetPsData(commonbl.TestPsResponse(), logger)
	chDesc := make(chan *prometheus.Desc, expectedDescChanels)
	exporter := NewSambaExporter(requestHandler, responseHandler, logger, "0.0.0", 5, exportSettings)
	exporter.setDescriptionsFromResponse(locks, processes, shares, psData, chDesc)
	chMet := make(chan prometheus.Metric, expectedMetChanels)
	exporter.setMetricsFromResponse(locks, processes, shares, psData, 1, 1, 31, chMet)

	if len(chMet) != expectedMetChanels {
		t.Errorf("Got %d metric channels, but expected %d", len(chMet), expectedMetChanels)
	}

	if logger.GetErrorCount() != 0 {
		t.Errorf("The ErrorCount '%d' is not the expected '0'", logger.GetErrorCount())
	}
}

func TestSetMetricsFromResponseNoShare(t *testing.T) {
	exportSettings := statisticsGenerator.StatisticsGeneratorSettings{false, false, true, false, false}
	expectedDescChanels := 38
	expectedMetChanels := 62
	requestHandler := commonbl.NewPipeHandler(true, commonbl.RequestPipe)
	responseHandler := commonbl.NewPipeHandler(true, commonbl.ResposePipe)
	logger := testhelper.NewTestLogger(true)
	locks := smbstatusreader.GetLockData(smbstatusout.LockData4Lines, logger)
	shares := smbstatusreader.GetShareData(smbstatusout.ShareData4Lines, logger)
	processes := smbstatusreader.GetProcessData(smbstatusout.ProcessData4Lines, logger)
	psData := smbstatusreader.GetPsData(commonbl.TestPsResponse(), logger)
	chDesc := make(chan *prometheus.Desc, expectedDescChanels)
	exporter := NewSambaExporter(requestHandler, responseHandler, logger, "0.0.0", 5, exportSettings)
	exporter.setDescriptionsFromResponse(locks, processes, shares, psData, chDesc)
	chMet := make(chan prometheus.Metric, expectedMetChanels)
	exporter.setMetricsFromResponse(locks, processes, shares, psData, 1, 1, 31, chMet)

	if len(chMet) != expectedMetChanels {
		t.Errorf("Got %d metric channels, but expected %d", len(chMet), expectedMetChanels)
	}

	if logger.GetErrorCount() != 0 {
		t.Errorf("The ErrorCount '%d' is not the expected '0'", logger.GetErrorCount())
	}
}

func TestSetMetricsFromEmptyResponse1(t *testing.T) {
	expectedDescChanels := 38
	expectedMetChanels := 19
	requestHandler := commonbl.NewPipeHandler(true, commonbl.RequestPipe)
	responseHandler := commonbl.NewPipeHandler(true, commonbl.ResposePipe)
	logger := testhelper.NewTestLogger(true)
	locks := smbstatusreader.GetLockData(smbstatusout.LockData0Line, logger)
	shares := smbstatusreader.GetShareData(smbstatusout.ShareData0Line, logger)
	processes := smbstatusreader.GetProcessData(smbstatusout.ProcessData0Lines, logger)
	psData := smbstatusreader.GetPsData(commonbl.TestPsResponseEmpty(), logger)
	chDesc := make(chan *prometheus.Desc, expectedDescChanels)
	exporter := NewSambaExporter(requestHandler, responseHandler, logger, "0.0.0", 5, getNewStatisticGenSettings())
	exporter.setDescriptionsFromResponse(locks, processes, shares, psData, chDesc)
	chMet := make(chan prometheus.Metric, expectedMetChanels)
	exporter.setMetricsFromResponse(locks, processes, shares, psData, 1, 1, 32, chMet)

	if len(chMet) != expectedMetChanels {
		t.Errorf("Got %d metric chanels, but expected %d", len(chMet), expectedMetChanels)
	}

	if logger.GetErrorCount() != 0 {
		t.Errorf("The ErrorCount '%d' is not the expected '0'", logger.GetErrorCount())
	}
}

func TestSetMetricsFromEmptyResponse2(t *testing.T) {
	expectedDescChanels := 38
	expectedMetChanels := 19
	requestHandler := commonbl.NewPipeHandler(true, commonbl.RequestPipe)
	responseHandler := commonbl.NewPipeHandler(true, commonbl.ResposePipe)
	logger := testhelper.NewTestLogger(true)
	locks := smbstatusreader.GetLockData(smbstatusout.LockDataEmpty, logger)
	shares := smbstatusreader.GetShareData(smbstatusout.ShareDataEmpty, logger)
	processes := smbstatusreader.GetProcessData(smbstatusout.ProcessDataEmpty, logger)
	psData := smbstatusreader.GetPsData(commonbl.TestPsResponseEmpty(), logger)
	chDesc := make(chan *prometheus.Desc, expectedDescChanels)
	exporter := NewSambaExporter(requestHandler, responseHandler, logger, "0.0.0", 5, getNewStatisticGenSettings())
	exporter.setDescriptionsFromResponse(locks, processes, shares, psData, chDesc)
	chMet := make(chan prometheus.Metric, expectedMetChanels)
	exporter.setMetricsFromResponse(locks, processes, shares, psData, 1, 1, 32, chMet)

	if len(chMet) != expectedMetChanels {
		t.Errorf("Got %d metric chanels, but expected %d", len(chMet), expectedMetChanels)
	}

	if logger.GetErrorCount() != 0 {
		t.Errorf("The ErrorCount '%d' is not the expected '0'", logger.GetErrorCount())
	}
}

func TestSetGaugeDescriptionNoLabel(t *testing.T) {
	requestHandler := commonbl.NewPipeHandler(true, commonbl.RequestPipe)
	responseHandler := commonbl.NewPipeHandler(true, commonbl.ResposePipe)
	logger := testhelper.NewTestLogger(true)
	help := "My help"
	name := "my_name"
	ch := make(chan *prometheus.Desc, 1)
	exporter := NewSambaExporter(requestHandler, responseHandler, logger, "0.0.0", 5, getNewStatisticGenSettings())

	exporter.setGaugeDescriptionNoLabel(name, help, ch)

	desc := <-ch

	if desc == nil {
		t.Errorf("There was no description added to the chanel")
	}

	descString := desc.String()
	if !strings.Contains(descString, help) {
		t.Errorf("The description does not contain the given help")
	}

	if !strings.Contains(descString, fmt.Sprintf("samba_%s", name)) {
		t.Errorf("The description does not contain the name")
	}

	if logger.GetErrorCount() != 0 {
		t.Errorf("The ErrorCount '%d' is not the expected '0'", logger.GetErrorCount())
	}
}

func TestSetGaugeDescriptionWithLabel(t *testing.T) {
	requestHandler := commonbl.NewPipeHandler(true, commonbl.RequestPipe)
	responseHandler := commonbl.NewPipeHandler(true, commonbl.ResposePipe)
	logger := testhelper.NewTestLogger(true)
	help := "My help"
	name := "my_name"
	labels := map[string]string{"key1": "value1", "key2": "value2"}
	ch := make(chan *prometheus.Desc, 1)
	exporter := NewSambaExporter(requestHandler, responseHandler, logger, "0.0.0", 5, getNewStatisticGenSettings())

	exporter.setGaugeDescriptionWithLabel(name, help, labels, ch)

	desc := <-ch

	if desc == nil {
		t.Errorf("There was no description added to the chanel")
	}

	descString := desc.String()
	if !strings.Contains(descString, help) {
		t.Errorf("The description does not contain the given help")
	}

	if !strings.Contains(descString, fmt.Sprintf("samba_%s", name)) {
		t.Errorf("The description does not contain the name")
	}

	for key, _ := range labels {
		if !strings.Contains(descString, key) {
			t.Errorf("The Description does not contain the expected label")
		}
	}

	if logger.GetErrorCount() != 0 {
		t.Errorf("The ErrorCount '%d' is not the expected '0'", logger.GetErrorCount())
	}
}

func TestSetGaugeIntMetricNoLabel(t *testing.T) {
	requestHandler := commonbl.NewPipeHandler(true, commonbl.RequestPipe)
	responseHandler := commonbl.NewPipeHandler(true, commonbl.ResposePipe)
	logger := testhelper.NewTestLogger(true)
	help := "My help"
	name := "my_name"
	chDesc := make(chan *prometheus.Desc, 1)
	exporter := NewSambaExporter(requestHandler, responseHandler, logger, "0.0.0", 5, getNewStatisticGenSettings())
	exporter.setGaugeDescriptionNoLabel(name, help, chDesc)
	desc := <-chDesc
	if desc == nil {
		t.Errorf("There was no description added to the chanel")
	}
	chMet := make(chan prometheus.Metric, 1)
	exporter.setGaugeIntMetricNoLabel(name, 42.0, chMet)

	met := <-chMet

	if met == nil {
		t.Errorf("Got no metric from the chanel")
	}

	if met.Desc().String() != desc.String() {
		t.Errorf("The metrics description is not the expected")
	}

	if logger.GetErrorCount() != 0 {
		t.Errorf("The ErrorCount '%d' is not the expected '0'", logger.GetErrorCount())
	}
}

func TestSetGaugeIntMetricNoDescription(t *testing.T) {
	requestHandler := commonbl.NewPipeHandler(true, commonbl.RequestPipe)
	responseHandler := commonbl.NewPipeHandler(true, commonbl.ResposePipe)
	logger := testhelper.NewTestLogger(true)
	exporter := NewSambaExporter(requestHandler, responseHandler, logger, "0.0.0", 5, getNewStatisticGenSettings())
	name := "my_name"
	chMet := make(chan prometheus.Metric, 1)
	exporter.setGaugeIntMetricNoLabel(name, 42.0, chMet)

	if len(chMet) != 0 {
		t.Errorf("Got metric from the chanel but expected none")
	}

	if logger.GetErrorCount() != 1 {
		t.Errorf("The ErrorCount '%d' is not the expected '1'", logger.GetErrorCount())
	}

	if logger.WrittenErrors[0] != "Error: No description found for my_name" {
		t.Errorf("The error message '%s' is not the expected 'Error: No description found for my_name'", logger.WrittenErrors[0])
	}
}

func TestSetGaugeIntMetricWithLabel(t *testing.T) {
	requestHandler := commonbl.NewPipeHandler(true, commonbl.RequestPipe)
	responseHandler := commonbl.NewPipeHandler(true, commonbl.ResposePipe)
	logger := testhelper.NewTestLogger(true)
	help := "My help"
	name := "my_name"
	labels := map[string]string{"key1": "value1", "key2": "value2"}
	chDesc := make(chan *prometheus.Desc, 1)
	exporter := NewSambaExporter(requestHandler, responseHandler, logger, "0.0.0", 5, getNewStatisticGenSettings())
	exporter.setGaugeDescriptionWithLabel(name, help, labels, chDesc)
	desc := <-chDesc
	if desc == nil {
		t.Errorf("There was no description added to the chanel")
	}
	chMet := make(chan prometheus.Metric, 1)
	exporter.setGaugeIntMetricWithLabel(name, 42.0, labels, chMet)

	met := <-chMet

	if met == nil {
		t.Errorf("Got no metric from the chanel")
	}

	if met.Desc().String() != desc.String() {
		t.Errorf("The metrics description is not the expected")
	}

	if logger.GetErrorCount() != 0 {
		t.Errorf("The ErrorCount '%d' is not the expected '0'", logger.GetErrorCount())
	}
}

func TestSetGaugeIntMetricWithLabelNoDescription(t *testing.T) {
	requestHandler := commonbl.NewPipeHandler(true, commonbl.RequestPipe)
	responseHandler := commonbl.NewPipeHandler(true, commonbl.ResposePipe)
	logger := testhelper.NewTestLogger(true)
	labels := map[string]string{"key1": "value1", "key2": "value2"}
	exporter := NewSambaExporter(requestHandler, responseHandler, logger, "0.0.0", 5, getNewStatisticGenSettings())
	name := "my_name"
	chMet := make(chan prometheus.Metric, 1)
	exporter.setGaugeIntMetricWithLabel(name, 42.0, labels, chMet)

	if len(chMet) != 0 {
		t.Errorf("Got metric from the chanel but expected none")
	}

	if logger.GetErrorCount() != 1 {
		t.Errorf("The ErrorCount '%d' is not the expected '1'", logger.GetErrorCount())
	}

	if logger.WrittenErrors[0] != "Error: No description found for metric 'my_name'" {
		t.Errorf("The error message '%s' is not the expected 'Error: No description found for metric 'my_name''", logger.WrittenErrors[0])
	}
}
