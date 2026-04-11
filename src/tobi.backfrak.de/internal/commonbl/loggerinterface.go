package commonbl

import (
	"strings"
)

type LogLevel int

const (
	Error LogLevel = iota
	Warning
	Information
	Verbose
)

var ValidLogLevelSettings = map[LogLevel]string{
	Error:       "Error",
	Warning:     "Warning",
	Information: "Information",
	Verbose:     "Verbose",
}

func StringToLogLevel(logLEvelStr string) (LogLevel, bool) {
	for logLevel, str := range ValidLogLevelSettings {
		if str == logLEvelStr {
			return logLevel, true
		}
	}

	return -1, false
}

func GetValidLogLevels() []string {
	var ret []string
	for _, level := range ValidLogLevelSettings {
		ret = append(ret, level)
	}

	return ret
}

// Logger - Interface for logger implementations
type Logger interface {
	// Get the current log level setting of a logger instance
	GetLogLevelSetting() LogLevel

	// WriteInformation - Write a Info message to Stdout, will be prefixed with "Information: "
	WriteInformation(message string)

	// WriteVerbose - Write a Verbose message to Stdout
	// The message will be prefixed with "Verbose :"
	WriteVerbose(message string)

	// WriteWarning - Write a Warning message to Stdout
	// The message will be prefixed with "Warning :"
	WriteWarning(message string)

	// WriteErrorMessage - Write the message to Stderr. The Message will be prefixed with "Error: "
	WriteErrorMessage(message string)

	// WriteError - Writes the err.Error() output to Stderr
	WriteError(err error)

	// WriteError - Writes the 'err.Error() - addition' output to Stderr
	WriteErrorWithAddition(err error, addition string)
}

// Get the right logger depending on the input parameters
func GetLogger(logFilePath string, logLevelSetting LogLevel) (Logger, error) {
	trimmedPath := strings.TrimSpace(logFilePath)
	if trimmedPath != "" {
		return NewFileLogger(logLevelSetting, trimmedPath)
	}

	return NewConsoleLogger(logLevelSetting), nil

}

func FilterMessage(messageLogLevel, settingsLogLevel LogLevel) bool {

	if messageLogLevel <= settingsLogLevel {
		return true
	}

	return false
}
