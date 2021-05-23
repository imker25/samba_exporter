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
	Logger         commonbl.Logger
}

// Get a new instance of the SambaExporter
func NewSambaExporter(requestHandler commonbl.PipeHandler, responseHander commonbl.PipeHandler, logger commonbl.Logger) *SambaExporter {
	var ret SambaExporter
	ret.RequestHandler = requestHandler
	ret.ResponseHander = responseHander
	ret.Descriptions = make(map[string]prometheus.Desc)
	var err error
	ret.hostName, err = os.Hostname()
	if err != nil {
		ret.hostName = "127.0.0.1"
	}
	ret.Logger = logger

	return &ret
}

// Describe function for the Prometheus Exporter Interface
func (smbExporter *SambaExporter) Describe(ch chan<- *prometheus.Desc) {
	smbExporter.Logger.WriteVerbose("Request samba_statusd to get prometheus descriptions")
	locks, processes, shares, errGet := pipecomunication.GetSambaStatus(smbExporter.RequestHandler, smbExporter.ResponseHander, smbExporter.Logger)
	if errGet != nil {
		smbExporter.Logger.WriteError(errGet)

		// Exit with panic, since this means there are no descriptions setup for further operation
		panic(errGet)
	}

	smbExporter.Logger.WriteVerbose("Handle samba_statusd response and set prometheus descriptions")
	stats := statisticsGenerator.GetSmbStatistics(locks, processes, shares)
	if stats == nil {
		smbExporter.Logger.WriteError(pipecomunication.NewSmbStatusUnexpectedResponseError("Empty response from samba_statusd"))

		// Exit with panic, since this means there are no descriptions setup for further operation
		panic(errGet)
	}

	for _, stat := range stats {
		// Example with label
		//desc := prometheus.NewDesc(prometheus.BuildFQName(EXPORTER_LABEL_PREFIX, "", stat.Name), stat.Help, []string{"machine"}, nil)

		// Without label
		desc := prometheus.NewDesc(prometheus.BuildFQName(EXPORTER_LABEL_PREFIX, "", stat.Name), stat.Help, []string{}, nil)
		smbExporter.Descriptions[stat.Name] = *desc
		ch <- desc
	}

	descSmbUp := prometheus.NewDesc(prometheus.BuildFQName(EXPORTER_LABEL_PREFIX, "", "server_up"), "1 if the samba server seems to be running", []string{}, nil)
	smbExporter.Descriptions["server_up"] = *descSmbUp
	ch <- descSmbUp

	descStatusUp := prometheus.NewDesc(prometheus.BuildFQName(EXPORTER_LABEL_PREFIX, "", "satutsd_up"), "1 if the samba_statusd seems to be running", []string{}, nil)
	smbExporter.Descriptions["satutsd_up"] = *descStatusUp
	ch <- descStatusUp
}

// Collect function for the Prometheus Exporter Interface
func (smbExporter *SambaExporter) Collect(ch chan<- prometheus.Metric) {
	smbExporter.Logger.WriteVerbose("Request samba_statusd to get prometheus metrics")
	smbStatusUp := 1
	smbServerUp := 1
	locks, processes, shares, errGet := pipecomunication.GetSambaStatus(smbExporter.RequestHandler, smbExporter.ResponseHander, smbExporter.Logger)
	if errGet != nil {
		switch errGet.(type) {
		case *pipecomunication.SmbStatusTimeOutError:
			smbStatusUp = 0
			smbServerUp = 0
		case *pipecomunication.SmbStatusUnexpectedResponseError:
			smbServerUp = 0
		default:
			smbExporter.Logger.WriteError(errGet)
			return
		}
	}

	smbExporter.Logger.WriteVerbose("Handle samba_statusd response and set prometheus metrics")

	descSmbUp, foundDescSmbUp := smbExporter.Descriptions["server_up"]
	if foundDescSmbUp == false {
		smbExporter.Logger.WriteErrorMessage(fmt.Sprintf("No description found for %s", "server_up"))
	}
	metSmbUp := prometheus.MustNewConstMetric(&descSmbUp, prometheus.GaugeValue, float64(smbServerUp))
	ch <- metSmbUp

	descStatUp, foundDescStatUp := smbExporter.Descriptions["satutsd_up"]
	if foundDescStatUp == false {
		smbExporter.Logger.WriteErrorMessage(fmt.Sprintf("No description found for %s", "satutsd_up"))
	}
	metStatUp := prometheus.MustNewConstMetric(&descStatUp, prometheus.GaugeValue, float64(smbStatusUp))
	ch <- metStatUp

	stats := statisticsGenerator.GetSmbStatistics(locks, processes, shares)
	if stats == nil {
		smbExporter.Logger.WriteError(pipecomunication.NewSmbStatusUnexpectedResponseError("Empty response from samba_statusd"))
		return
	}

	for _, stat := range stats {
		desc, found := smbExporter.Descriptions[stat.Name]
		if found == false {
			smbExporter.Logger.WriteErrorMessage(fmt.Sprintf("No description found for %s", stat.Name))
		}
		// Example with label
		// met := prometheus.MustNewConstMetric(&desc, prometheus.GaugeValue, float64(stat.Value), smbExporter.hostName)

		// Without label
		met := prometheus.MustNewConstMetric(&desc, prometheus.GaugeValue, float64(stat.Value))
		ch <- met
	}

}
