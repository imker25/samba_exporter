package smbexporter

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"fmt"
	"os"

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
	Descriptions   map[string]prometheus.Desc
	hostName       string
}

// Get a new instance of the SambaExporter
func NewSambaExporter(requestHandler commonbl.PipeHandler, responseHander commonbl.PipeHandler) *SambaExporter {
	var ret SambaExporter
	ret.RequestHandler = requestHandler
	ret.ResponseHander = responseHander
	ret.Descriptions = make(map[string]prometheus.Desc)
	var err error
	ret.hostName, err = os.Hostname()
	if err != nil {
		ret.hostName = "127.0.0.1"
	}

	return &ret
}

// Describe function for the Prometheus Exporter Interface
func (smbExporter *SambaExporter) Describe(ch chan<- *prometheus.Desc) {
	locks, processes, shares, errGet := pipecomunication.GetSambaStatus(smbExporter.RequestHandler, smbExporter.ResponseHander)
	if errGet != nil {
		fmt.Fprintln(os.Stderr, errGet)
		return
	}
	stats := statisticsGenerator.GetSmbStatistics(locks, processes, shares)

	if stats == nil {
		fmt.Fprintln(os.Stderr, pipecomunication.NewSmbStatusUnexpectedResponseError("Empty response from samba_statusd"))
		return
	}

	for _, stat := range stats {
		desc := prometheus.NewDesc(prometheus.BuildFQName(EXPORTER_LABEL_PREFIX, "", stat.Name), stat.Help, []string{"machine"}, nil)
		smbExporter.Descriptions[stat.Name] = *desc
		ch <- desc
	}
}

// Collect function for the Prometheus Exporter Interface
func (smbExporter *SambaExporter) Collect(ch chan<- prometheus.Metric) {
	locks, processes, shares, errGet := pipecomunication.GetSambaStatus(smbExporter.RequestHandler, smbExporter.ResponseHander)
	if errGet != nil {
		fmt.Fprintln(os.Stderr, errGet)
		return
	}
	stats := statisticsGenerator.GetSmbStatistics(locks, processes, shares)

	if stats == nil {
		fmt.Fprintln(os.Stderr, pipecomunication.NewSmbStatusUnexpectedResponseError("Empty response from samba_statusd"))
		return
	}

	for _, stat := range stats {
		desc, found := smbExporter.Descriptions[stat.Name]
		if found == false {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("No description found for %s", stat.Name))
		}
		met := prometheus.MustNewConstMetric(&desc, prometheus.GaugeValue, float64(stat.Value), smbExporter.hostName)
		ch <- met
	}

}
