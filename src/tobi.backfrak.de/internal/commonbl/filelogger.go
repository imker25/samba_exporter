package commonbl

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// FileLogger - A "class" with log functions
type FileLogger struct {
	LogLevelSetting LogLevel
	FullFilePath    string
	internalLoggers map[LogLevel]*log.Logger
}

// Get a new instance of the Logger
func NewFileLogger(verbose bool, fullFilePath string) (*FileLogger, error) {

	logLevelSetting := Information
	if verbose {
		logLevelSetting = Verbose
	}

	ret, err := NewFileLogger2(logLevelSetting, fullFilePath)

	return ret, err
}

// Get a new instance of the Logger
func NewFileLogger2(logLevelSetting LogLevel, fullFilePath string) (*FileLogger, error) {
	logFileDir := filepath.Dir(fullFilePath)
	if !directoryExists(logFileDir) {
		return nil, NewDirectoryNotExistError(logFileDir)
	}

	// If the file doesn't exist, create it or append to the file
	file, err := os.OpenFile(fullFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	internalLoggers := make(map[LogLevel]*log.Logger)
	internalLoggers[Information] = log.New(file, fmt.Sprintf("%s: ", ValidLogLevelSettings[Information]), log.LstdFlags|log.Lmsgprefix)
	internalLoggers[Verbose] = log.New(file, fmt.Sprintf("%s: ", ValidLogLevelSettings[Verbose]), log.LstdFlags|log.Lmsgprefix)
	internalLoggers[Error] = log.New(file, fmt.Sprintf("%s: ", ValidLogLevelSettings[Error]), log.LstdFlags|log.Lmsgprefix)
	internalLoggers[Warning] = log.New(file, fmt.Sprintf("%s: ", ValidLogLevelSettings[Warning]), log.LstdFlags|log.Lmsgprefix)

	ret := FileLogger{logLevelSetting, fullFilePath, internalLoggers}

	return &ret, nil
}

// GetVerbose - Tell if logger is verbose or not
func (logger *FileLogger) GetVerbose() bool {
	if logger.LogLevelSetting == Verbose {
		return true
	}

	return false
}

// Get the current log level setting of a logger instance
func (logger *FileLogger) GetLogLevelSetting() LogLevel {
	return logger.LogLevelSetting
}

// WriteInformation - Write a Info message to Stdout, will be prefixed with "Information: "
func (logger *FileLogger) WriteInformation(message string) {
	logger.writeLogMessage(message, Information)
}

// WriteWarning - Write a Warning message to Stdout, will be prefixed with "Information: "
func (logger *FileLogger) WriteWarning(message string) {
	logger.writeLogMessage(message, Warning)
}

// WriteVerbose - Write a Verbose message to Stdout. Message will be written only if logger.Verbose is true.
// The message will be prefixed with "Verbose :"
func (logger *FileLogger) WriteVerbose(message string) {
	logger.writeLogMessage(message, Verbose)
}

// WriteErrorMessage - Write the message to Stderr. The Message will be prefixed with "Error: "
func (logger *FileLogger) WriteErrorMessage(message string) {
	logger.writeLogMessage(message, Error)
}

// WriteError - Writes the err.Error() output to Stderr
func (logger *FileLogger) WriteError(err error) {
	logger.WriteErrorMessage(err.Error())
}

// WriteError - Writes the 'err.Error() - addition' output to Stderr
func (logger *FileLogger) WriteErrorWithAddition(err error, addition string) {
	logger.WriteErrorMessage(fmt.Sprintf("%s - %s", err.Error(), addition))
}

func directoryExists(path string) bool {
	if stat, err := os.Stat(path); err == nil && stat.IsDir() {
		return true
	}

	return false
}

func (logger *FileLogger) writeLogMessage(message string, messageLogLevel LogLevel) {
	if !FilterMessage(messageLogLevel, logger.LogLevelSetting) {
		return
	}

	var trimmedMsg string
	if messageLogLevel == Error {
		trimmedMsg = strings.TrimPrefix(message, "Error: ")
	} else {
		trimmedMsg = message
	}

	logWriter := logger.internalLoggers[messageLogLevel]
	logWriter.Println(trimmedMsg)
}
