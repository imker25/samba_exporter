package main

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"tobi.backfrak.de/internal/commonbl"
	"tobi.backfrak.de/internal/smbexporterbl/pipecomunication"
	"tobi.backfrak.de/internal/smbexporterbl/smbexporter"
	"tobi.backfrak.de/internal/smbexporterbl/smbstatusreader"
	"tobi.backfrak.de/internal/smbexporterbl/statisticsGenerator"
)

// Authors - Information about the authors of the program. You might want to add your name here when contributing to this software
const Authors = "tobi@backfrak.de"

// The version of this program, will be set at compile time by the ./build.sh build script
var version = "undefined"

// The logger used in the program
var logger commonbl.Logger

func main() {
	handleComandlineOptions()
	requestHandler := *commonbl.NewPipeHandler(params.Test, commonbl.RequestPipe)
	responseHandler := *commonbl.NewPipeHandler(params.Test, commonbl.ResposePipe)
	logger = *commonbl.NewLogger(params.Verbose)

	if !strings.HasPrefix(params.MetricsPath, "/") {
		params.MetricsPath = fmt.Sprintf("/%s", params.MetricsPath)
	}

	if params.Verbose {
		args := ""
		for _, arg := range os.Args {
			args = fmt.Sprintf("%s %s", args, arg)
		}
		fmt.Fprintln(os.Stdout, fmt.Sprintf("Call: %s", args))
		if !params.PrintVersion {
			printVersion()
		}
	}

	logger.WriteVerbose(fmt.Sprintf("Named pipe for requests: %s", requestHandler.GetPipeFilePath()))
	logger.WriteVerbose(fmt.Sprintf("Named pipe for response: %s", responseHandler.GetPipeFilePath()))

	if params.PrintVersion {
		printVersion()
		os.Exit(0)
	}

	if params.Help {
		flag.Usage()
		os.Exit(0)
	}

	if params.DoNotExportUser {
		logger.WriteVerbose("-not-expose-user-data set, will not export user data")
	}

	if params.DoNotExportClient {
		logger.WriteVerbose("-not-expose-client-data set, will not export client data")
	}

	if params.DoNotExportEncryption {
		logger.WriteVerbose("-not-expose-encryption-data set, will not export encryption data")
	}

	if params.TestPipeMode {
		errTest := testPipeMode(requestHandler, responseHandler)
		if errTest != nil {
			logger.WriteError(errTest)
			os.Exit(-2)
		}
		os.Exit(0)
	}

	// Ensure we exit clean on term and kill signals
	go waitforKillSignalAndExit()
	go waitforTermSignalAndExit()

	logger.WriteVerbose("Setup prometheus exporter")

	exporter := smbexporter.NewSambaExporter(requestHandler, responseHandler, logger, version, params.RequestTimeOut, params.StatisticsGeneratorSettings)
	prometheus.MustRegister(exporter)

	logger.WriteInformation(fmt.Sprintf("Started %s, get metrics on http://%s%s", os.Args[0], params.ListenAddress, params.MetricsPath))

	http.Handle(params.MetricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<html>
			<head><title>Samba Exporter</title></head>
			<body>
			<h1>Samba Exporter</h1>
			<p><a href='` + params.MetricsPath + `'>Metrics</a></p>
			</body>
			</html>`))
	})

	errListen := http.ListenAndServe(params.ListenAddress, nil)
	if errListen != nil {
		logger.WriteError(errListen)
		os.Exit(-1)
	}
}

func testPipeMode(requestHandler commonbl.PipeHandler, responseHandler commonbl.PipeHandler) error {
	var processes []smbstatusreader.ProcessData
	var shares []smbstatusreader.ShareData
	var locks []smbstatusreader.LockData
	var psData []commonbl.PsUtilPidData
	var errGet error

	logger.WriteVerbose("Request samba_statusd to get metrics for test-pipe mode")
	locks, processes, shares, psData, errGet = pipecomunication.GetSambaStatus(requestHandler, responseHandler, logger, params.RequestTimeOut)
	if errGet != nil {
		return errGet
	}

	logger.WriteVerbose("Handle samba_statusd  response in test-pipe mode")

	for _, share := range shares {
		fmt.Fprintln(os.Stdout, share.String())
	}
	for _, process := range processes {
		fmt.Fprintln(os.Stdout, process.String())
	}
	for _, lock := range locks {
		fmt.Fprintln(os.Stdout, lock.String())
	}

	for _, ps := range psData {
		fmt.Fprintln(os.Stdout, ps.String())
	}

	stats := statisticsGenerator.GetSmbStatistics(locks, processes, shares, params.StatisticsGeneratorSettings)
	stats = append(stats, statisticsGenerator.GetSmbdMetrics(psData, params.DoNotExportPid)...)
	for _, stat := range stats {
		fmt.Fprintln(os.Stdout, fmt.Sprintf("%s_%s: %f", smbexporter.EXPORTER_LABEL_PREFIX, stat.Name, stat.Value))
	}

	return nil
}

func waitforKillSignalAndExit() {
	killSignal := make(chan os.Signal, syscall.SIGKILL)
	signal.Notify(killSignal, os.Interrupt)
	<-killSignal

	logger.WriteInformation(fmt.Sprintf("End %s due to kill signal", os.Args[0]))

	os.Exit(0)
}

func waitforTermSignalAndExit() {
	termSignal := make(chan os.Signal, syscall.SIGTERM)
	signal.Notify(termSignal, os.Interrupt)
	<-termSignal

	logger.WriteInformation(fmt.Sprintf("End %s due to terminate signal", os.Args[0]))

	os.Exit(0)
}

// Prints the version string
func printVersion() {
	fmt.Fprintln(os.Stdout, getVersion())
}

// Get the version string
func getVersion() string {
	return fmt.Sprintf("Version: %s", version)
}
