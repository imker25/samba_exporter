package main

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"flag"
	"fmt"
	"os"

	"tobi.backfrak.de/internal/commonbl"
)

// Authors - Information about the authors of the program. You might want to add your name here when contributing to this software
const Authors = "tobi@backfrak.de"

// The version of this program, will be set at compile time by the gradle build script
var version = "undefined"

func main() {
	handleComandlineOptions()
	pipeHander := *commonbl.NewPipeHandler(params.Test)
	if params.Verbose {
		args := ""
		for _, arg := range os.Args {
			args = fmt.Sprintf("%s %s", args, arg)
		}
		fmt.Fprintln(os.Stdout, fmt.Sprintf("Call: %s", args))
		if !params.PrintVersion {
			printVersion()
		}
		fmt.Fprintln(os.Stdout, fmt.Sprintf("Use named pipe: %s", pipeHander.GetPipeFilePath()))
	}

	if params.PrintVersion {
		printVersion()
		os.Exit(0)
	}

	if params.Help {
		flag.Usage()
		os.Exit(0)
	}

}

// Prints the version string
func printVersion() {
	fmt.Fprintln(os.Stdout, getVersion())
}

// Get the version string
func getVersion() string {
	return fmt.Sprintf("Version: %s", version)
}
