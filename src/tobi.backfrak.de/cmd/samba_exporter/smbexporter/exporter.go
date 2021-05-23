package smbexporter

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"fmt"

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
	Logger         commonbl.Logger
}

// Get a new instance of the SambaExporter
func NewSambaExporter(requestHandler commonbl.PipeHandler, responseHander commonbl.PipeHandler, logger commonbl.Logger) *SambaExporter {
	var ret SambaExporter
	ret.RequestHandler = requestHandler
	ret.ResponseHander = responseHander
	ret.Logger = logger
	ret.Descriptions = make(map[string]prometheus.Desc)

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
	smbExporter.setGaugeDescriptionNoLabel("server_up", "1 if the samba server seems to be running", ch)
	smbExporter.setGaugeDescriptionNoLabel("satutsd_up", "1 if the samba_statusd seems to be running", ch)

	for _, stat := range stats {
		if stat.Labels == nil {
			smbExporter.setGaugeDescriptionNoLabel(stat.Name, stat.Help, ch)
		} else {
			smbExporter.setGaugeDescriptionWithLabel(stat.Name, stat.Help, stat.Labels, ch)
		}
	}
}

// Collect function for the Prometheus Exporter Interface
func (smbExporter *SambaExporter) Collect(ch chan<- prometheus.Metric) {
	smbExporter.Logger.WriteVerbose("Request samba_statusd to get prometheus metrics")
	smbStatusUp := 1
	smbServerUp := 1
	locks, processes, shares, errGet := pipecomunication.GetSambaStatus(smbExporter.RequestHandler, smbExporter.ResponseHander, smbExporter.Logger)
	if errGet != nil {
		smbExporter.Logger.WriteError(errGet)
		switch errGet.(type) {
		case *pipecomunication.SmbStatusTimeOutError:
			smbStatusUp = 0
			smbServerUp = 0
		case *pipecomunication.SmbStatusUnexpectedResponseError:
			smbServerUp = 0
		default:
			return
		}
	}

	smbExporter.Logger.WriteVerbose("Handle samba_statusd response and set prometheus metrics")
	smbExporter.setGaugeIntMetricNoLabel("server_up", float64(smbServerUp), ch)
	smbExporter.setGaugeIntMetricNoLabel("satutsd_up", float64(smbStatusUp), ch)

	stats := statisticsGenerator.GetSmbStatistics(locks, processes, shares)
	if stats == nil {
		smbExporter.Logger.WriteError(pipecomunication.NewSmbStatusUnexpectedResponseError("Empty response from samba_statusd"))
		return
	}

	for _, stat := range stats {
		if stat.Labels == nil {
			smbExporter.setGaugeIntMetricNoLabel(stat.Name, stat.Value, ch)
		} else {
			smbExporter.setGaugeIntMetricWithLabel(stat.Name, stat.Value, stat.Labels, ch)
		}
	}

}

func (smbExporter *SambaExporter) setGaugeIntMetricNoLabel(name string, value float64, ch chan<- prometheus.Metric) {
	desc, found := smbExporter.Descriptions[name]
	if found == false {
		smbExporter.Logger.WriteErrorMessage(fmt.Sprintf("No description found for %s", name))
		return
	}

	met := prometheus.MustNewConstMetric(&desc, prometheus.GaugeValue, value)
	ch <- met
}

func (smbExporter *SambaExporter) setGaugeIntMetricWithLabel(name string, value float64, labels map[string]string, ch chan<- prometheus.Metric) {
	desc, found := smbExporter.Descriptions[name]
	if found == false {
		smbExporter.Logger.WriteErrorMessage(fmt.Sprintf("No description found for %s", name))
		return
	}
	var labelValues []string
	for _, value := range labels {
		if value != "" {
			labelValues = append(labelValues, value)
		} else {
			// if a labels value is "", we don't add the value at all
			return
		}
	}

	met := prometheus.MustNewConstMetric(&desc, prometheus.GaugeValue, value, labelValues...)
	ch <- met
}

func (smbExporter *SambaExporter) setGaugeDescriptionNoLabel(name string, help string, ch chan<- *prometheus.Desc) {
	desc := prometheus.NewDesc(prometheus.BuildFQName(EXPORTER_LABEL_PREFIX, "", name), help, []string{}, nil)
	smbExporter.Descriptions[name] = *desc
	ch <- desc
}

func (smbExporter *SambaExporter) setGaugeDescriptionWithLabel(name string, help string, labels map[string]string, ch chan<- *prometheus.Desc) {
	// Since the a the same label can have multiple values, we need only one description
	_, found := smbExporter.Descriptions[name]
	if found == false {
		var labelKeys []string
		for key, _ := range labels {
			labelKeys = append(labelKeys, key)
		}

		desc := prometheus.NewDesc(prometheus.BuildFQName(EXPORTER_LABEL_PREFIX, "", name), help, labelKeys, nil)
		smbExporter.Descriptions[name] = *desc
		ch <- desc
	}
}
