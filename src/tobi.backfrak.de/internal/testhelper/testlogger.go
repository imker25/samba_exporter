package testhelper

import (
	"fmt"
)

// TestLogger - A "class" with log functions
type TestLogger struct {
	Verbose         bool
	WrittenMessages []string
	WrittenErrors   []string
}

// Get a new instance of the Logger
func NewTestLogger(verbose bool) *TestLogger {
	writtenMessages := []string{}
	writtenErrors := []string{}
	ret := TestLogger{verbose, writtenMessages, writtenErrors}

	return &ret
}

// GetVerbose - Tell if logger is verbose or not
func (logger *TestLogger) GetVerbose() bool {
	return logger.Verbose
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
	logger.WrittenMessages = append(logger.WrittenMessages, fmt.Sprintf("Information: %s", message))

	return
}

// WriteVerbose - Write a Verbose message to Stdout. Message will be written only if logger.Verbose is true.
// The message will be prefixed with "Verbose :"
func (logger *TestLogger) WriteVerbose(message string) {
	if logger.Verbose {
		logger.WrittenMessages = append(logger.WrittenMessages, fmt.Sprintf("Verbose: %s", message))
	}

	return
}

// WriteErrorMessage - Write the message to Stderr. The Message will be prefixed with "Error: "
func (logger *TestLogger) WriteErrorMessage(message string) {
	logger.WrittenErrors = append(logger.WrittenErrors, fmt.Sprintf("Error: %s", message))
}

// WriteError - Writes the err.Error() output to Stderr
func (logger *TestLogger) WriteError(err error) {
	logger.WrittenErrors = append(logger.WrittenErrors, err.Error())
}

// WriteError - Writes the 'err.Error() - addition' output to Stderr
func (logger *TestLogger) WriteErrorWithAddition(err error, addition string) {
	logger.WrittenErrors = append(logger.WrittenErrors, fmt.Sprintf("%s - %s", err.Error(), addition))
}
