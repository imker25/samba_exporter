package commonbl

import (
	"fmt"
	"strings"
	"testing"
)

func TestGetLoggerForFileLogger(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	ensureLogFileDirExists()
	logger, err := GetLogger(logfile_path, Information)

	if err != nil {
		t.Errorf("Got error '%s' but expected none", err.Error())
	}

	switch logger.(type) {
	case *FileLogger:
		fmt.Println("OK")
	default:
		t.Errorf("The logger is not the expected FileLogger")
	}

	logger2, err2 := GetLogger(logfile_path, Verbose)
	if err2 != nil {
		t.Errorf("Got error '%s' but expected none", err2.Error())
	}

	switch logger2.(type) {
	case *FileLogger:
		fmt.Println("OK")
	default:
		t.Errorf("The logger is not the expected FileLogger")
	}
}

func TestGetLoggerForFileLogger2(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	ensureLogFileDirExists()
	logger, err := GetLogger(logfile_path, Error)

	if err != nil {
		t.Errorf("Got error '%s' but expected none", err.Error())
	}

	switch logger.(type) {
	case *FileLogger:
		fmt.Println("OK")
	default:
		t.Errorf("The logger is not the expected FileLogger")
	}

	if logger.GetLogLevelSetting() != Error {
		t.Errorf("The loglevel setting of the logger '%d' is not the expected '%d'", logger.GetLogLevelSetting(), Error)
	}

	logger2, err2 := GetLogger(logfile_path, Verbose)
	if err2 != nil {
		t.Errorf("Got error '%s' but expected none", err2.Error())
	}

	switch logger2.(type) {
	case *FileLogger:
		fmt.Println("OK")
	default:
		t.Errorf("The logger is not the expected FileLogger")
	}

	if logger2.GetLogLevelSetting() != Verbose {
		t.Errorf("The loglevel setting of the logger '%d' is not the expected '%d'", logger2.GetLogLevelSetting(), Verbose)
	}
}

func TestGetLoggerForConsoleLogger(t *testing.T) {
	logger1, err1 := GetLogger(" ", Information)

	if err1 != nil {
		t.Errorf("Got error '%s' but expected none", err1.Error())
	}

	switch logger1.(type) {
	case *ConsoleLogger:
		fmt.Println("OK")
	default:
		t.Errorf("The logger is not the expected ConsoleLogger")
	}

	logger2, err2 := GetLogger(" ", Error)
	if err2 != nil {
		t.Errorf("Got error '%s' but expected none", err2.Error())
	}

	switch logger2.(type) {
	case *ConsoleLogger:
		fmt.Println("OK")
	default:
		t.Errorf("The logger is not the expected ConsoleLogger")
	}
}

func TestGetLoggerForConsoleLogger2(t *testing.T) {
	logger1, err1 := GetLogger(" ", Warning)

	if err1 != nil {
		t.Errorf("Got error '%s' but expected none", err1.Error())
	}

	switch logger1.(type) {
	case *ConsoleLogger:
		fmt.Println("OK")
	default:
		t.Errorf("The logger is not the expected ConsoleLogger")
	}

	if logger1.GetLogLevelSetting() != Warning {
		t.Errorf("The loglevel setting of the logger '%d' is not the expected '%d'", logger1.GetLogLevelSetting(), Warning)
	}

	logger2, err2 := GetLogger(" ", Information)
	if err2 != nil {
		t.Errorf("Got error '%s' but expected none", err2.Error())
	}

	switch logger2.(type) {
	case *ConsoleLogger:
		fmt.Println("OK")
	default:
		t.Errorf("The logger is not the expected ConsoleLogger")
	}

	if logger2.GetLogLevelSetting() != Information {
		t.Errorf("The loglevel setting of the logger '%d' is not the expected '%d'", logger1.GetLogLevelSetting(), Information)
	}
}

func TestGetValidLogLevels(t *testing.T) {
	validLevels := GetValidLogLevels()

	validLevelsStr := fmt.Sprintf("%s", validLevels)

	if len(validLevels) != 4 {
		t.Errorf("The number of valid levels %d is not the expected number %d", len(validLevels), 4)
	}

	if !strings.Contains(validLevelsStr, "Information") {
		t.Errorf("The valid levels '%s' does not contain '%s'", validLevelsStr, "Information")
	}

	if !strings.Contains(validLevelsStr, "Error") {
		t.Errorf("The valid levels '%s' does not contain '%s'", validLevelsStr, "Error")
	}

	if !strings.Contains(validLevelsStr, "Verbose") {
		t.Errorf("The valid levels '%s' does not contain '%s'", validLevelsStr, "Verbose")
	}

	if !strings.Contains(validLevelsStr, "Warning") {
		t.Errorf("The valid levels '%s' does not contain '%s'", validLevelsStr, "Warning")
	}

	if "Error" != ValidLogLevelSettings[0] {
		t.Errorf("validLogLevelSettings[0] is '%s', but '%s' is expected. ", "Error", ValidLogLevelSettings[0])
	}

	if "Error" != ValidLogLevelSettings[Error] {
		t.Errorf("validLogLevelSettings[Error] is '%s', but '%s' is expected. ", "Error", ValidLogLevelSettings[Error])
	}

	if "Warning" != ValidLogLevelSettings[1] {
		t.Errorf("validLogLevelSettings[1] is '%s', but '%s' is expected. ", "Warning", ValidLogLevelSettings[1])
	}

	if "Warning" != ValidLogLevelSettings[Warning] {
		t.Errorf("validLogLevelSettings[Warning] is '%s', but '%s' is expected. ", "Warning", ValidLogLevelSettings[Warning])
	}

	if "Information" != ValidLogLevelSettings[2] {
		t.Errorf("validLogLevelSettings[2] is '%s', but '%s' is expected. ", "Information", ValidLogLevelSettings[2])
	}

	if "Information" != ValidLogLevelSettings[Information] {
		t.Errorf("validLogLevelSettings[Information] is '%s', but '%s' is expected. ", "Information", ValidLogLevelSettings[Information])
	}

	if "Verbose" != ValidLogLevelSettings[3] {
		t.Errorf("validLogLevelSettings[3] is '%s', but '%s' is expected", "Verbose", ValidLogLevelSettings[3])
	}

	if "Verbose" != ValidLogLevelSettings[Verbose] {
		t.Errorf("validLogLevelSettings[Verbose] is '%s', but '%s' is expected", "Verbose", ValidLogLevelSettings[Verbose])
	}
}

func TestFilterMessage(t *testing.T) {
	if FilterMessage(Error, Verbose) != true {
		t.Errorf("'FilterMessage(Error, Verbose)' returns '%t', but '%t' is expected", FilterMessage(Error, Verbose), true)
	}

	if FilterMessage(Verbose, Error) != false {
		t.Errorf("'FilterMessage(Verbose, Error)' returns '%t', but '%t' is expected", FilterMessage(Verbose, Error), false)
	}

	if FilterMessage(Warning, Information) != true {
		t.Errorf("'FilterMessage(Warning, Information)' returns '%t', but '%t' is expected", FilterMessage(Warning, Information), true)
	}

	if FilterMessage(Information, Warning) != false {
		t.Errorf("'FilterMessage(Information, Warning)' returns '%t', but '%t' is expected", FilterMessage(Information, Warning), false)
	}
}

func TestStringToLogLevel(t *testing.T) {
	logLevel, result := StringToLogLevel("Error")
	if !result {
		t.Errorf("The string 'Error' could not be converted into a LogLevel")
	}
	if logLevel != Error {
		t.Errorf("The string 'Error' is converted into log Level '%d'", logLevel)
	}

	logLevel, result = StringToLogLevel("Warning")
	if !result {
		t.Errorf("The string 'Warning' could not be converted into a LogLevel")
	}
	if logLevel != Warning {
		t.Errorf("The string 'Warning' is converted into log Level '%d'", logLevel)
	}

	logLevel, result = StringToLogLevel("Information")
	if !result {
		t.Errorf("The string 'Information' could not be converted into a LogLevel")
	}
	if logLevel != Information {
		t.Errorf("The string 'Information' is converted into log Level '%d'", logLevel)
	}

	logLevel, result = StringToLogLevel("Verbose")
	if !result {
		t.Errorf("The string 'Verbose' could not be converted into a LogLevel")
	}
	if logLevel != Verbose {
		t.Errorf("The string 'Verbose' is converted into log Level '%d'", logLevel)
	}

	logLevel, result = StringToLogLevel("Bla")
	if result {
		t.Errorf("The string 'Bla' could be converted into a LogLevel")
	}
	if logLevel != -1 {
		t.Errorf("The string 'Bla' is converted into log Level '%d' but should be '-1'", logLevel)
	}
}
