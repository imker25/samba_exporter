package commonbl

import "strings"

type LogLevel int

const (
	ErrorOnly LogLevel = iota
	Warning
	Information
	Verbose
)

var validLogLevelSettings = map[LogLevel]string{
	ErrorOnly:   "ErrorOnly",
	Warning:     "Warning",
	Information: "Information",
	Verbose:     "Verbose",
}

func GetValidLogLevels() []string {
	var ret []string
	for _, level := range validLogLevelSettings {
		ret = append(ret, level)
	}

	return ret
}

// Logger - Interface for logger implementations
type Logger interface {
	// GetVerbose - Tell if logger is verbose or not
	GetVerbose() bool

	// WriteInformation - Write a Info message to Stdout, will be prefixed with "Information: "
	WriteInformation(message string)

	// WriteVerbose - Write a Verbose message to Stdout. Message will be written only if logger.Verbose is true.
	// The message will be prefixed with "Verbose :"
	WriteVerbose(message string)

	// WriteErrorMessage - Write the message to Stderr. The Message will be prefixed with "Error: "
	WriteErrorMessage(message string)

	// WriteError - Writes the err.Error() output to Stderr
	WriteError(err error)

	// WriteError - Writes the 'err.Error() - addition' output to Stderr
	WriteErrorWithAddition(err error, addition string)
}

// Get the right logger depending on the input parameters
func GetLogger(logFilePath string, verbose bool) (Logger, error) {
	trimmedPath := strings.TrimSpace(logFilePath)
	if trimmedPath != "" {
		return NewFileLogger(verbose, trimmedPath)
	}

	return NewConsoleLogger(verbose), nil

}
