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
	"tobi.backfrak.de/internal/smbstatusout"
)

func TestNewSambaExporter(t *testing.T) {
	requestHandler := *commonbl.NewPipeHandler(true, commonbl.RequestPipe)
	responseHandler := *commonbl.NewPipeHandler(true, commonbl.ResposePipe)
	logger := *commonbl.NewLogger(true)
	exporter := NewSambaExporter(requestHandler, responseHandler, logger, "0.0.0", 5)

	if exporter.RequestHandler.PipeType != commonbl.RequestPipe {
		t.Errorf("The exporter.RequestHandler is not of the expected type")
	}

	if exporter.ResponseHander.PipeType != commonbl.ResposePipe {
		t.Errorf("The exporter.RequestHandler is not of the expected type")
	}

	if exporter.Descriptions == nil {
		t.Errorf("exporter.Descriptions are nil")
	}

	if logger.Verbose != exporter.Logger.Verbose {
		t.Errorf("The exporter uses the wrong logger")
	}

	if exporter.Version != "0.0.0" {
		t.Errorf("The Version \"%s\" is not expected", exporter.Version)
	}
}

func TestSetDescriptionsFromResponse(t *testing.T) {
	expectedChanels := 14
	requestHandler := *commonbl.NewPipeHandler(true, commonbl.RequestPipe)
	responseHandler := *commonbl.NewPipeHandler(true, commonbl.ResposePipe)
	logger := *commonbl.NewLogger(true)
	locks := smbstatusreader.GetLockData(smbstatusout.LockDataNoData, logger)
	shares := smbstatusreader.GetShareData(smbstatusout.ShareDataOneLine, logger)
	processes := smbstatusreader.GetProcessData(smbstatusout.ProcessDataOneLine, logger)
	ch := make(chan *prometheus.Desc, expectedChanels)
	exporter := NewSambaExporter(requestHandler, responseHandler, logger, "0.0.0", 5)
	exporter.setDescriptionsFromResponse(locks, processes, shares, ch)

	if len(ch) != expectedChanels {
		t.Errorf("The number of descriptions is not expected")
	}
}

func TestSetMetricsFromResponse(t *testing.T) {
	expectedDescChanels := 20
	expectedMetChanels := 20
	requestHandler := *commonbl.NewPipeHandler(true, commonbl.RequestPipe)
	responseHandler := *commonbl.NewPipeHandler(true, commonbl.ResposePipe)
	logger := *commonbl.NewLogger(true)
	locks := smbstatusreader.GetLockData(smbstatusout.LockData4Lines, logger)
	shares := smbstatusreader.GetShareData(smbstatusout.ShareData4Lines, logger)
	processes := smbstatusreader.GetProcessData(smbstatusout.ProcessData4Lines, logger)
	chDesc := make(chan *prometheus.Desc, expectedDescChanels)
	exporter := NewSambaExporter(requestHandler, responseHandler, logger, "0.0.0", 5)
	exporter.setDescriptionsFromResponse(locks, processes, shares, chDesc)
	chMet := make(chan prometheus.Metric, expectedMetChanels)
	exporter.setMetricsFromResponse(locks, processes, shares, 1, 1, chMet)

	if len(chMet) != expectedMetChanels {
		t.Errorf("Got %d metric chanels, but expected %d", len(chMet), expectedMetChanels)
	}

}

func TestSetMetricsFromEmptyResponse(t *testing.T) {
	expectedDescChanels := 14
	expectedMetChanels := 8
	requestHandler := *commonbl.NewPipeHandler(true, commonbl.RequestPipe)
	responseHandler := *commonbl.NewPipeHandler(true, commonbl.ResposePipe)
	logger := *commonbl.NewLogger(true)
	locks := smbstatusreader.GetLockData(smbstatusout.LockData0Line, logger)
	shares := smbstatusreader.GetShareData(smbstatusout.LockData0Line, logger)
	processes := smbstatusreader.GetProcessData(smbstatusout.LockData0Line, logger)
	chDesc := make(chan *prometheus.Desc, expectedDescChanels)
	exporter := NewSambaExporter(requestHandler, responseHandler, logger, "0.0.0", 5)
	exporter.setDescriptionsFromResponse(locks, processes, shares, chDesc)
	chMet := make(chan prometheus.Metric, expectedMetChanels)
	exporter.setMetricsFromResponse(locks, processes, shares, 1, 1, chMet)

	if len(chMet) != expectedMetChanels {
		t.Errorf("Got %d metric chanels, but expected %d", len(chMet), expectedMetChanels)
	}

}

func TestSetGaugeDescriptionNoLabel(t *testing.T) {
	requestHandler := *commonbl.NewPipeHandler(true, commonbl.RequestPipe)
	responseHandler := *commonbl.NewPipeHandler(true, commonbl.ResposePipe)
	logger := *commonbl.NewLogger(true)
	help := "My help"
	name := "my_name"
	ch := make(chan *prometheus.Desc, 1)
	exporter := NewSambaExporter(requestHandler, responseHandler, logger, "0.0.0", 5)

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

}

func TestSetGaugeDescriptionWithLabel(t *testing.T) {
	requestHandler := *commonbl.NewPipeHandler(true, commonbl.RequestPipe)
	responseHandler := *commonbl.NewPipeHandler(true, commonbl.ResposePipe)
	logger := *commonbl.NewLogger(true)
	help := "My help"
	name := "my_name"
	labels := map[string]string{"key1": "value1", "key2": "value2"}
	ch := make(chan *prometheus.Desc, 1)
	exporter := NewSambaExporter(requestHandler, responseHandler, logger, "0.0.0", 5)

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

}

func TestSetGaugeIntMetricNoLabel(t *testing.T) {
	requestHandler := *commonbl.NewPipeHandler(true, commonbl.RequestPipe)
	responseHandler := *commonbl.NewPipeHandler(true, commonbl.ResposePipe)
	logger := *commonbl.NewLogger(true)
	help := "My help"
	name := "my_name"
	chDesc := make(chan *prometheus.Desc, 1)
	exporter := NewSambaExporter(requestHandler, responseHandler, logger, "0.0.0", 5)
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
}

func TestSetGaugeIntMetricNoDescription(t *testing.T) {
	requestHandler := *commonbl.NewPipeHandler(true, commonbl.RequestPipe)
	responseHandler := *commonbl.NewPipeHandler(true, commonbl.ResposePipe)
	logger := *commonbl.NewLogger(true)
	exporter := NewSambaExporter(requestHandler, responseHandler, logger, "0.0.0", 5)
	name := "my_name"
	chMet := make(chan prometheus.Metric, 1)
	exporter.setGaugeIntMetricNoLabel(name, 42.0, chMet)

	if len(chMet) != 0 {
		t.Errorf("Got metric from the chanel but expected none")
	}

}

func TestSetGaugeIntMetricWithLabel(t *testing.T) {
	requestHandler := *commonbl.NewPipeHandler(true, commonbl.RequestPipe)
	responseHandler := *commonbl.NewPipeHandler(true, commonbl.ResposePipe)
	logger := *commonbl.NewLogger(true)
	help := "My help"
	name := "my_name"
	labels := map[string]string{"key1": "value1", "key2": "value2"}
	chDesc := make(chan *prometheus.Desc, 1)
	exporter := NewSambaExporter(requestHandler, responseHandler, logger, "0.0.0", 5)
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
}

func TestSetGaugeIntMetricWithLabelNoDescription(t *testing.T) {
	requestHandler := *commonbl.NewPipeHandler(true, commonbl.RequestPipe)
	responseHandler := *commonbl.NewPipeHandler(true, commonbl.ResposePipe)
	logger := *commonbl.NewLogger(true)
	labels := map[string]string{"key1": "value1", "key2": "value2"}
	exporter := NewSambaExporter(requestHandler, responseHandler, logger, "0.0.0", 5)
	name := "my_name"
	chMet := make(chan prometheus.Metric, 1)
	exporter.setGaugeIntMetricWithLabel(name, 42.0, labels, chMet)

	if len(chMet) != 0 {
		t.Errorf("Got metric from the chanel but expected none")
	}

}
