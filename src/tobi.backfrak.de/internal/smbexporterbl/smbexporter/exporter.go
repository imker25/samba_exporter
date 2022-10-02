package smbexporter

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"tobi.backfrak.de/internal/commonbl"
	"tobi.backfrak.de/internal/smbexporterbl/pipecomunication"
	"tobi.backfrak.de/internal/smbexporterbl/smbstatusreader"
	"tobi.backfrak.de/internal/smbexporterbl/statisticsGenerator"
)

// The Prefix for labels of this prometheus exporter
const EXPORTER_LABEL_PREFIX = "samba"

// SambaExporter - The class that implements the Prometheus Exporter Interface
type SambaExporter struct {
	RequestHandler              commonbl.PipeHandler
	ResponseHander              commonbl.PipeHandler
	Logger                      commonbl.Logger
	Version                     string
	RequestTimeOut              int
	StatisticsGeneratorSettings statisticsGenerator.StatisticsGeneratorSettings

	// Used to ensure that every metric is only added once
	descriptions map[string]prometheus.Desc

	// Used to ensure that the order of labels is always the same for a given metric
	metricsLabelList map[string][]string
}

// Get a new instance of the SambaExporter
func NewSambaExporter(requestHandler commonbl.PipeHandler, responseHander commonbl.PipeHandler, logger commonbl.Logger, version string, requestTimeOut int, statisticsGeneratorSettings statisticsGenerator.StatisticsGeneratorSettings) *SambaExporter {
	var ret SambaExporter
	ret.RequestHandler = requestHandler
	ret.ResponseHander = responseHander
	ret.Logger = logger
	ret.Version = version
	ret.RequestTimeOut = requestTimeOut
	ret.descriptions = make(map[string]prometheus.Desc)
	ret.StatisticsGeneratorSettings = statisticsGeneratorSettings
	ret.metricsLabelList = make(map[string][]string)

	return &ret
}

// Describe function for the Prometheus Exporter Interface
func (smbExporter *SambaExporter) Describe(ch chan<- *prometheus.Desc) {
	smbExporter.Logger.WriteVerbose("Request samba_statusd to get prometheus descriptions")
	locks, processes, shares, psData, errGet := pipecomunication.GetSambaStatus(smbExporter.RequestHandler, smbExporter.ResponseHander, smbExporter.Logger, smbExporter.RequestTimeOut)
	if errGet != nil {
		smbExporter.Logger.WriteError(errGet)

		// Exit with panic, since this means there are no descriptions setup for further operation
		panic(errGet)
	}
	smbExporter.setDescriptionsFromResponse(locks, processes, shares, psData, ch)

	return
}

// Collect function for the Prometheus Exporter Interface
func (smbExporter *SambaExporter) Collect(ch chan<- prometheus.Metric) {
	smbExporter.Logger.WriteVerbose("Request samba_statusd to get prometheus metrics")
	smbStatusUp := 1
	smbServerUp := 1
	start := time.Now()
	locks, processes, shares, psData, errGet := pipecomunication.GetSambaStatus(smbExporter.RequestHandler, smbExporter.ResponseHander, smbExporter.Logger, smbExporter.RequestTimeOut)
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
	elapsed := time.Since(start)
	elapsedFloat := float64(elapsed.Milliseconds())
	smbExporter.setMetricsFromResponse(locks, processes, shares, psData, smbStatusUp, smbServerUp, elapsedFloat, ch)

	return
}

func (smbExporter *SambaExporter) setMetricsFromResponse(locks []smbstatusreader.LockData, processes []smbstatusreader.ProcessData, shares []smbstatusreader.ShareData, psData []commonbl.PsUtilPidData, smbStatusUp int, smbServerUp int, requestTime float64, ch chan<- prometheus.Metric) {
	smbExporter.Logger.WriteVerbose("Handle samba_statusd response and set prometheus metrics")
	smbExporter.setGaugeIntMetricNoLabel("server_up", float64(smbServerUp), ch)
	smbExporter.setGaugeIntMetricNoLabel("satutsd_up", float64(smbStatusUp), ch)
	smbExporter.setGaugeIntMetricWithLabel("exporter_information", 1, map[string]string{"version": smbExporter.Version}, ch)

	stats := statisticsGenerator.GetSmbStatistics(locks, processes, shares, smbExporter.StatisticsGeneratorSettings)
	if stats == nil {
		smbExporter.Logger.WriteError(pipecomunication.NewSmbStatusUnexpectedResponseError("Empty response from samba_statusd"))
		return
	}
	stats = append(stats, statisticsGenerator.GetSmbdMetrics(psData, smbExporter.StatisticsGeneratorSettings.DoNotExportPid)...)

	for _, stat := range stats {
		if stat.Labels == nil {
			smbExporter.setGaugeIntMetricNoLabel(stat.Name, stat.Value, ch)
		} else {
			smbExporter.setGaugeIntMetricWithLabel(stat.Name, stat.Value, stat.Labels, ch)
		}
	}
	smbExporter.setGaugeIntMetricNoLabel("request_time", requestTime, ch)
}

func (smbExporter *SambaExporter) setDescriptionsFromResponse(locks []smbstatusreader.LockData, processes []smbstatusreader.ProcessData, shares []smbstatusreader.ShareData, psData []commonbl.PsUtilPidData, ch chan<- *prometheus.Desc) {
	smbExporter.Logger.WriteVerbose("Handle samba_statusd response and set prometheus descriptions")
	stats := statisticsGenerator.GetSmbStatistics(locks, processes, shares, smbExporter.StatisticsGeneratorSettings)
	if stats == nil {
		err := pipecomunication.NewSmbStatusUnexpectedResponseError("Empty response from samba_statusd")
		smbExporter.Logger.WriteError(err)

		// Exit with panic, since this means there are no descriptions setup for further operation
		panic(err)
	}
	stats = append(stats, statisticsGenerator.GetSmbdMetrics(psData, smbExporter.StatisticsGeneratorSettings.DoNotExportPid)...)

	smbExporter.setGaugeDescriptionNoLabel("server_up", "1 if the samba server seems to be running", ch)
	smbExporter.setGaugeDescriptionNoLabel("satutsd_up", "1 if the samba_statusd seems to be running", ch)
	smbExporter.setGaugeDescriptionWithLabel("exporter_information", "Information of the samba_exporter", map[string]string{"version": smbExporter.Version}, ch)

	for _, stat := range stats {
		if stat.Labels == nil {
			smbExporter.setGaugeDescriptionNoLabel(stat.Name, stat.Help, ch)
		} else {
			smbExporter.setGaugeDescriptionWithLabel(stat.Name, stat.Help, stat.Labels, ch)
		}
	}

	smbExporter.setGaugeDescriptionNoLabel("request_time", "Time it took to reqest the samba status from samba_statusd [ms]", ch)
}

func (smbExporter *SambaExporter) setGaugeIntMetricNoLabel(name string, value float64, ch chan<- prometheus.Metric) {
	desc, found := smbExporter.descriptions[name]
	if found == false {
		smbExporter.Logger.WriteErrorMessage(fmt.Sprintf("No description found for %s", name))
		return
	}

	met := prometheus.MustNewConstMetric(&desc, prometheus.GaugeValue, value)
	ch <- met
}

func (smbExporter *SambaExporter) setGaugeIntMetricWithLabel(name string, value float64, labels map[string]string, ch chan<- prometheus.Metric) {
	desc, found := smbExporter.descriptions[name]
	if !found {
		smbExporter.Logger.WriteErrorMessage(fmt.Sprintf("No description found for metric '%s'", name))
		return
	}

	// Ensure the expected order of labels is known for this metric (see bug #79)
	labelKeys, foundInLabelList := smbExporter.metricsLabelList[name]
	if !foundInLabelList {
		smbExporter.Logger.WriteErrorMessage(fmt.Sprintf("No label keys found for metric '%s'", name))
		return
	}

	// Validate that the given labels list is the same length as the expected label list for this metric (see bug #79)
	if len(labelKeys) != len(labels) {
		smbExporter.Logger.WriteErrorMessage(fmt.Sprintf(
			"The number of labels given with metric '%s' ('%d') does not match the expected number '%d'",
			name, len(labels), len(labelKeys)))
		return
	}

	labelValues := make([]string, len(labelKeys))
	// The loop over the expected label is done explicit form 0 to end, so the order can not be lost (see bug #79)
	for i := 0; i < len(labelKeys); i++ {
		key := labelKeys[i]
		value, foundValue := labels[key]
		if !foundValue {
			smbExporter.Logger.WriteErrorMessage(fmt.Sprintf("No label with key '%s' found for metric '%s'", key, name))
			return
		}
		if value != "" {
			// The set is done to a explicit field to ensure the order of labels is not lost (see bug #79)
			labelValues[i] = value
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
	smbExporter.descriptions[name] = *desc
	ch <- desc
}

func (smbExporter *SambaExporter) setGaugeDescriptionWithLabel(name string, help string, labels map[string]string, ch chan<- *prometheus.Desc) {
	// Since the a the same label can have multiple values, we need only one description
	_, found := smbExporter.descriptions[name]

	if !found {
		var labelKeys []string
		for key, _ := range labels {
			labelKeys = append(labelKeys, key)
		}

		smbExporter.metricsLabelList[name] = labelKeys
		desc := prometheus.NewDesc(prometheus.BuildFQName(EXPORTER_LABEL_PREFIX, "", name), help, labelKeys, nil)
		smbExporter.descriptions[name] = *desc
		ch <- desc
	}
}
