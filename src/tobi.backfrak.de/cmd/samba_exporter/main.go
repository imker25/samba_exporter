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
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"tobi.backfrak.de/cmd/samba_exporter/pipecomunication"
	"tobi.backfrak.de/cmd/samba_exporter/smbexporter"
	"tobi.backfrak.de/cmd/samba_exporter/smbstatusreader"
	"tobi.backfrak.de/cmd/samba_exporter/statisticsGenerator"
	"tobi.backfrak.de/internal/commonbl"
)

// Authors - Information about the authors of the program. You might want to add your name here when contributing to this software
const Authors = "tobi@backfrak.de"

// The version of this program, will be set at compile time by the gradle build script
var version = "undefined"

type SmbResponse struct {
	Data  string
	Error error
}

func main() {
	handleComandlineOptions()
	requestHandler := *commonbl.NewPipeHandler(params.Test, commonbl.RequestPipe)
	responseHandler := *commonbl.NewPipeHandler(params.Test, commonbl.ResposePipe)
	if params.Verbose {
		args := ""
		for _, arg := range os.Args {
			args = fmt.Sprintf("%s %s", args, arg)
		}
		fmt.Fprintln(os.Stdout, fmt.Sprintf("Call: %s", args))
		if !params.PrintVersion {
			printVersion()
		}
		fmt.Fprintln(os.Stdout, fmt.Sprintf("Named pipe for requests: %s", requestHandler.GetPipeFilePath()))
		fmt.Fprintln(os.Stdout, fmt.Sprintf("Named pipe for response: %s", responseHandler.GetPipeFilePath()))
	}

	if params.PrintVersion {
		printVersion()
		os.Exit(0)
	}

	if params.Help {
		flag.Usage()
		os.Exit(0)
	}

	if params.TestPipeMode {
		errTest := testPipeMode(requestHandler, responseHandler)
		if errTest != nil {
			fmt.Fprintln(os.Stderr, errTest)
			os.Exit(-2)
		}
		os.Exit(0)
	}

	// Ensure we exit clean on term and kill signals
	go waitforKillSignalAndExit()
	go waitforTermSignalAndExit()

	fmt.Fprintln(os.Stdout, fmt.Sprintf("Start %s, get metrics on %s%s", os.Args[0], params.ListenAddress, params.MetricsPath))

	exporter := smbexporter.NewSambaExporter(requestHandler, responseHandler, params.Verbose)
	prometheus.MustRegister(exporter)

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
		fmt.Fprintln(os.Stderr, errListen)
		os.Exit(-1)
	}
}

func testPipeMode(requestHandler commonbl.PipeHandler, responseHandler commonbl.PipeHandler) error {
	var processes []smbstatusreader.ProcessData
	var shares []smbstatusreader.ShareData
	var locks []smbstatusreader.LockData
	var errGet error
	if params.Verbose {
		fmt.Fprintln(os.Stdout, "Request samba_statusd to get metrics for test-pipe mode")
	}
	locks, processes, shares, errGet = pipecomunication.GetSambaStatus(requestHandler, responseHandler, params.Verbose)
	if errGet != nil {
		return errGet
	}
	if params.Verbose {
		fmt.Fprintln(os.Stdout, "Handle samba_statusd  response in test-pipe mode")
	}
	for _, share := range shares {
		fmt.Fprintln(os.Stdout, share.String())
	}
	for _, process := range processes {
		fmt.Fprintln(os.Stdout, process.String())
	}
	for _, lock := range locks {
		fmt.Fprintln(os.Stdout, lock.String())
	}

	stats := statisticsGenerator.GetSmbStatistics(locks, processes, shares)
	for _, stat := range stats {
		fmt.Fprintln(os.Stdout, fmt.Sprintf("%s_%s: %d", smbexporter.EXPORTER_LABEL_PREFIX, stat.Name, stat.Value))
	}

	return nil
}

func waitforKillSignalAndExit() {
	killSignal := make(chan os.Signal, syscall.SIGKILL)
	signal.Notify(killSignal, os.Interrupt)
	<-killSignal

	if params.Verbose {
		fmt.Fprintln(os.Stdout, fmt.Sprintf("End: %s due to kill signal", os.Args[0]))
	}
	os.Exit(0)
}

func waitforTermSignalAndExit() {
	termSignal := make(chan os.Signal, syscall.SIGTERM)
	signal.Notify(termSignal, os.Interrupt)
	<-termSignal
	if params.Verbose {
		fmt.Fprintln(os.Stdout, fmt.Sprintf("End: %s due to terminate signal", os.Args[0]))
	}
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
