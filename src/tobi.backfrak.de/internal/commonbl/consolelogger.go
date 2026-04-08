package commonbl

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"fmt"
	"os"
	"strings"
)

// ConsoleLogger - A "class" with log functions
type ConsoleLogger struct {
	LogLevelSetting LogLevel
}

// Get a new instance of the Logger
func NewConsoleLogger(verbose bool) *ConsoleLogger {
	ret := ConsoleLogger{}
	if verbose {
		ret.LogLevelSetting = Verbose
	} else {
		ret.LogLevelSetting = Information
	}

	return &ret
}

func NewConsoleLogger2(logLevelSetting LogLevel) *ConsoleLogger {
	ret := ConsoleLogger{}
	ret.LogLevelSetting = logLevelSetting

	return &ret
}

// GetVerbose - Tell if logger is verbose or not
func (logger *ConsoleLogger) GetVerbose() bool {
	if logger.LogLevelSetting == Verbose {
		return true
	}

	return false
}

// Get the current log level setting of a logger instance
func (logger *ConsoleLogger) GetLogLevelSetting() LogLevel {
	return logger.LogLevelSetting
}

// WriteInformation - Write a Info message to Stdout, will be prefixed with "Information: "
func (logger *ConsoleLogger) WriteInformation(message string) {
	logger.writeMessageToStream(message, Information)
}

// WriteWarning - Write a Warning message to Stdout, will be prefixed with "Information: "
func (logger *ConsoleLogger) WriteWarning(message string) {
	logger.writeMessageToStream(message, Warning)
}

// WriteVerbose - Write a Verbose message to Stdout. Message will be written only if logger.Verbose is true.
// The message will be prefixed with "Verbose :"
func (logger *ConsoleLogger) WriteVerbose(message string) {
	logger.writeMessageToStream(message, Verbose)
}

// WriteErrorMessage - Write the message to Stderr. The Message will be prefixed with "Error: "
func (logger *ConsoleLogger) WriteErrorMessage(message string) {
	logger.writeMessageToStream(message, Error)
}

// WriteError - Writes the err.Error() output to Stderr
func (logger *ConsoleLogger) WriteError(err error) {
	logger.writeMessageToStream(err.Error(), Error)
}

// WriteError - Writes the 'err.Error() - addition' output to Stderr
func (logger *ConsoleLogger) WriteErrorWithAddition(err error, addition string) {
	logger.WriteErrorMessage(fmt.Sprintf("%s - %s", err.Error(), addition))
}

func (logger *ConsoleLogger) writeMessageToStream(message string, messageLogLevel LogLevel) {
	if !FilterMessage(messageLogLevel, logger.LogLevelSetting) {
		return
	}

	var outstream *os.File
	var trimmedMsg string
	if messageLogLevel == Error {
		outstream = os.Stderr
		trimmedMsg = strings.TrimPrefix(message, "Error: ")
	} else {
		outstream = os.Stdout
		trimmedMsg = message
	}

	fmt.Fprintln(outstream, fmt.Sprintf("%s: %s", ValidLogLevelSettings[messageLogLevel], trimmedMsg))
}
