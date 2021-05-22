package smbexporter

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"github.com/prometheus/client_golang/prometheus"
	"tobi.backfrak.de/internal/commonbl"
)

// SambaExporter - The class that implements the Prometheus Exporter Interface
type SambaExporter struct {
	RequestHandler commonbl.PipeHandler
	ResponseHander commonbl.PipeHandler
}

// Get a new instance of the SambaExporter
func NewSambaExporter(requestHandler commonbl.PipeHandler, responseHander commonbl.PipeHandler) *SambaExporter {
	var ret SambaExporter
	ret.RequestHandler = requestHandler
	ret.ResponseHander = responseHander

	return &ret
}

// Describe function for the Prometheus Exporter Interface
func (smbExporter *SambaExporter) Describe(ch chan<- *prometheus.Desc) {

}

// Collect function for the Prometheus Exporter Interface
func (smbExporter *SambaExporter) Collect(ch chan<- prometheus.Metric) {

}
