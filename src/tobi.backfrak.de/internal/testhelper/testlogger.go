package testhelper

import (
	"fmt"
	"strings"
	"sync"

	"tobi.backfrak.de/internal/commonbl"
)

// TestLogger - A "class" with log functions
type TestLogger struct {
	LogLevelSetting commonbl.LogLevel
	WrittenMessages []string
	WrittenErrors   []string
	mutex           sync.Mutex
}

// Get a new instance of the Logger
func NewTestLogger(verbose bool) *TestLogger {
	logLevelSetting := commonbl.Information
	if verbose {
		logLevelSetting = commonbl.Verbose
	}

	return NewTestLogger2(logLevelSetting)
}

// Get a new instance of the Logger
func NewTestLogger2(logLevelSetting commonbl.LogLevel) *TestLogger {
	writtenMessages := []string{}
	writtenErrors := []string{}
	mu := sync.Mutex{}
	ret := TestLogger{logLevelSetting, writtenMessages, writtenErrors, mu}

	return &ret
}

// GetVerbose - Tell if logger is verbose or not
func (logger *TestLogger) GetVerbose() bool {
	if logger.LogLevelSetting == commonbl.Verbose {
		return true
	}

	return false
}

// Get the current log level setting of a logger instance
func (logger *TestLogger) GetLogLevelSetting() commonbl.LogLevel {
	return logger.LogLevelSetting
}

// GetErrorCount - Get the number of error messages written to stderr
func (logger *TestLogger) GetErrorCount() int {
	return len(logger.WrittenErrors)
}

// GetMessageCount - Get the number of messages written to stdout
func (logger *TestLogger) GetMessageCount() int {
	return len(logger.WrittenMessages)
}

// Get the number of messages written to any output stream
func (logger *TestLogger) GetOutputCount() int {
	return len(logger.WrittenMessages) + len(logger.WrittenErrors)
}

// WriteInformation - Write a Info message to Stdout, will be prefixed with "Information: "
func (logger *TestLogger) WriteInformation(message string) {
	logger.writeLogMessage(message, commonbl.Information)
}

// WriteVerbose - Write a Verbose message to Stdout. Message will be written only if logger.Verbose is true.
// The message will be prefixed with "Verbose :"
func (logger *TestLogger) WriteVerbose(message string) {
	logger.writeLogMessage(message, commonbl.Verbose)
}

// WriteWarning - Write a Warning message to Stdout, will be prefixed with "Information: "
func (logger *TestLogger) WriteWarning(message string) {
	logger.writeLogMessage(message, commonbl.Warning)
}

// WriteErrorMessage - Write the message to Stderr. The Message will be prefixed with "Error: "
func (logger *TestLogger) WriteErrorMessage(message string) {
	logger.writeLogMessage(message, commonbl.Error)
}

// WriteError - Writes the err.Error() output to Stderr
func (logger *TestLogger) WriteError(err error) {
	logger.WriteErrorMessage(err.Error())
}

// WriteError - Writes the 'err.Error() - addition' output to Stderr
func (logger *TestLogger) WriteErrorWithAddition(err error, addition string) {
	logger.WriteErrorMessage(fmt.Sprintf("%s - %s", err.Error(), addition))
}

func (logger *TestLogger) writeLogMessage(message string, messageLogLevel commonbl.LogLevel) {
	if !commonbl.FilterMessage(messageLogLevel, logger.LogLevelSetting) {
		return
	}

	logger.mutex.Lock()
	defer logger.mutex.Unlock()

	if messageLogLevel == commonbl.Error {
		trimmedMsg := strings.TrimPrefix(message, "Error: ")
		logger.WrittenErrors = append(logger.WrittenErrors, fmt.Sprintf("%s: %s", commonbl.ValidLogLevelSettings[messageLogLevel], trimmedMsg))
		return
	}

	logger.WrittenMessages = append(logger.WrittenMessages, fmt.Sprintf("%s: %s", commonbl.ValidLogLevelSettings[messageLogLevel], message))
	return
}
