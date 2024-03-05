package commonbl

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// FileLogger - A "class" with log functions
type FileLogger struct {
	Verbose       bool
	FullFilePath  string
	infoLogger    *log.Logger
	verboseLogger *log.Logger
	errorLogger   *log.Logger
}

// Get a new instance of the Logger
func NewFileLogger(verbose bool, fullFilePath string) *FileLogger {

	// If the file doesn't exist, create it or append to the file
	file, err := os.OpenFile(fullFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	infoLogger := log.New(file, "Information: ", log.LstdFlags|log.Lmsgprefix /*|log.Lmicroseconds*/)
	verboseLogger := log.New(file, "Verbose: ", log.LstdFlags|log.Lmsgprefix /*|log.Lmicroseconds*/)
	errorLogger := log.New(file, "Error: ", log.LstdFlags|log.Lmsgprefix /*|log.Lmicroseconds*/)

	ret := FileLogger{verbose, fullFilePath, infoLogger, verboseLogger, errorLogger}

	return &ret
}

// GetVerbose - Tell if logger is verbose or not
func (logger *FileLogger) GetVerbose() bool {
	return logger.Verbose
}

// WriteInformation - Write a Info message to Stdout, will be prefixed with "Information: "
func (logger *FileLogger) WriteInformation(message string) {
	logger.infoLogger.Println(message)
}

// WriteVerbose - Write a Verbose message to Stdout. Message will be written only if logger.Verbose is true.
// The message will be prefixed with "Verbose :"
func (logger *FileLogger) WriteVerbose(message string) {
	if logger.Verbose {
		logger.verboseLogger.Println(message)
	}

}

// WriteErrorMessage - Write the message to Stderr. The Message will be prefixed with "Error: "
func (logger *FileLogger) WriteErrorMessage(message string) {
	trimedMsg := strings.TrimPrefix(message, "Error: ")
	logger.errorLogger.Println(trimedMsg)
}

// WriteError - Writes the err.Error() output to Stderr
func (logger *FileLogger) WriteError(err error) {
	trimedMsg := strings.TrimPrefix(err.Error(), "Error: ")
	logger.errorLogger.Println(trimedMsg)
}

// WriteError - Writes the 'err.Error() - addition' output to Stderr
func (logger *FileLogger) WriteErrorWithAddition(err error, addition string) {
	logger.WriteErrorMessage(fmt.Sprintf("%s - %s", err.Error(), addition))
}
