package smbexporter

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"tobi.backfrak.de/cmd/samba_exporter/pipecomunication"
	"tobi.backfrak.de/cmd/samba_exporter/statisticsGenerator"
	"tobi.backfrak.de/internal/commonbl"
)

// The Prefix for labels of this prometheus exporter
const EXPORTER_LABEL_PREFIX = "samba"

// SambaExporter - The class that implements the Prometheus Exporter Interface
type SambaExporter struct {
	RequestHandler commonbl.PipeHandler
	ResponseHander commonbl.PipeHandler
	mux            sync.Mutex
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
	smbExporter.mux.Lock()
	defer smbExporter.mux.Unlock()
	locks, processes, shares, errGet := pipecomunication.GetSambaStatus(smbExporter.RequestHandler, smbExporter.ResponseHander)
	if errGet != nil {
		return
	}
	stats := statisticsGenerator.GetSmbStatistics(locks, processes, shares)

	if stats == nil {
		return
	}
}

// Collect function for the Prometheus Exporter Interface
func (smbExporter *SambaExporter) Collect(ch chan<- prometheus.Metric) {
	smbExporter.mux.Lock()
	defer smbExporter.mux.Unlock()
	locks, processes, shares, errGet := pipecomunication.GetSambaStatus(smbExporter.RequestHandler, smbExporter.ResponseHander)
	if errGet != nil {
		return
	}
	stats := statisticsGenerator.GetSmbStatistics(locks, processes, shares)

	if stats == nil {
		return
	}

}
