package testhelper

import (
	"testing"

	"tobi.backfrak.de/internal/commonbl"
)

type TestError struct {
	message string
}

func (e *TestError) Error() string { // Implement the Error Interface for the TestError struct 
	return e.message
}

func TestNewTestLogger(t *testing.T) {
	logger := NewTestLogger(true)

	if !logger.Verbose {
		t.Errorf("Logger is not verbose but should")
	}

	if logger.GetErrorCount() != 0 {
		t.Errorf("Logger has ErrorCount '%d', but '0' is expected", logger.GetErrorCount())
	}

	if logger.GetMessageCount() != 0 {
		t.Errorf("Logger has MessageCount '%d', but '0' is expected", logger.GetMessageCount())
	}

	if logger.GetOutputCount() != 0 {
		t.Errorf("Logger has OutputCount '%d', but '0' is expected", logger.GetOutputCount())
	}

	iLogger := commonbl.Logger(logger)
	if iLogger.GetVerbose() == false {
		t.Errorf("Logger is not verbose but should")
	}

	logger.Verbose = false
	if iLogger.GetVerbose() == true {
		t.Errorf("Logger is verbose but not should")
	}
}

func TestWriteVerbose(t *testing.T) {
	logger := NewTestLogger(true)

	if !logger.Verbose {
		t.Errorf("Logger is not verbose but should")
	}

	iLogger := commonbl.Logger(logger)

	logger.WriteVerbose("just a message - 1")
	if logger.GetErrorCount() != 0 {
		t.Errorf("Logger has ErrorCount '%d', but '0' is expected", logger.GetErrorCount())
	}

	if logger.GetMessageCount() != 1 {
		t.Errorf("Logger has MessageCount '%d', but '1' is expected", logger.GetMessageCount())
	}

	if logger.GetOutputCount() != 1 {
		t.Errorf("Logger has OutputCount '%d', but '1' is expected", logger.GetOutputCount())
	}

	iLogger.WriteVerbose("just a message - 2")
	if logger.GetErrorCount() != 0 {
		t.Errorf("Logger has ErrorCount '%d', but '0' is expected", logger.GetErrorCount())
	}

	if logger.GetMessageCount() != 2 {
		t.Errorf("Logger has MessageCount '%d', but '2' is expected", logger.GetMessageCount())
	}

	if logger.GetOutputCount() != 2 {
		t.Errorf("Logger has OutputCount '%d', but '2' is expected", logger.GetOutputCount())
	}

	logger.Verbose = false
	iLogger.WriteVerbose("just a message - 3")
	if logger.GetErrorCount() != 0 {
		t.Errorf("Logger has ErrorCount '%d', but '0' is expected", logger.GetErrorCount())
	}

	if logger.GetMessageCount() != 2 {
		t.Errorf("Logger has MessageCount '%d', but '2' is expected", logger.GetMessageCount())
	}

	if logger.GetOutputCount() != 2 {
		t.Errorf("Logger has OutputCount '%d', but '2' is expected", logger.GetOutputCount())
	}

	if logger.WrittenMessages[1] != "Verbose: just a message - 2" {
		t.Errorf("The message '%s' is not the expected 'Verbose: just a message - 2'", logger.WrittenMessages[1])
	}
}

func TestWriteInformation(t *testing.T) {
	logger := NewTestLogger(true)

	if !logger.Verbose {
		t.Errorf("Logger is not verbose but should")
	}
	iLogger := commonbl.Logger(logger)

	iLogger.WriteInformation("just a message - 1")
	if logger.GetErrorCount() != 0 {
		t.Errorf("Logger has ErrorCount '%d', but '0' is expected", logger.GetErrorCount())
	}

	if logger.GetMessageCount() != 1 {
		t.Errorf("Logger has MessageCount '%d', but '1' is expected", logger.GetMessageCount())
	}

	if logger.GetOutputCount() != 1 {
		t.Errorf("Logger has OutputCount '%d', but '1' is expected", logger.GetOutputCount())
	}

	iLogger.WriteInformation("just a message - 2")
	if logger.GetErrorCount() != 0 {
		t.Errorf("Logger has ErrorCount '%d', but '0' is expected", logger.GetErrorCount())
	}

	if logger.GetMessageCount() != 2 {
		t.Errorf("Logger has MessageCount '%d', but '2' is expected", logger.GetMessageCount())
	}

	if logger.GetOutputCount() != 2 {
		t.Errorf("Logger has OutputCount '%d', but '2' is expected", logger.GetOutputCount())
	}

	if logger.WrittenMessages[1] != "Information: just a message - 2" {
		t.Errorf("The message '%s' is not the expected 'Information: just a message - 2'", logger.WrittenMessages[1])
	}
}

func TestWriteErrorMessage(t *testing.T) {
	logger := NewTestLogger(true)

	if !logger.Verbose {
		t.Errorf("Logger is not verbose but should")
	}
	iLogger := commonbl.Logger(logger)

	iLogger.WriteErrorMessage("just error 1")
	if logger.GetErrorCount() != 1 {
		t.Errorf("Logger has ErrorCount '%d', but '1' is expected", logger.GetErrorCount())
	}

	if logger.GetMessageCount() != 0 {
		t.Errorf("Logger has MessageCount '%d', but '0' is expected", logger.GetMessageCount())
	}

	if logger.GetOutputCount() != 1 {
		t.Errorf("Logger has OutputCount '%d', but '1' is expected", logger.GetOutputCount())
	}

	iLogger.WriteErrorMessage("just error 2")
	if logger.GetErrorCount() != 2 {
		t.Errorf("Logger has ErrorCount '%d', but '2' is expected", logger.GetErrorCount())
	}

	if logger.GetMessageCount() != 0 {
		t.Errorf("Logger has MessageCount '%d', but '0' is expected", logger.GetMessageCount())
	}

	if logger.GetOutputCount() != 2 {
		t.Errorf("Logger has OutputCount '%d', but '2' is expected", logger.GetOutputCount())
	}

	if logger.WrittenErrors[1] != "Error: just error 2" {
		t.Errorf("The message '%s' is not the expected 'Error: just error 2'", logger.WrittenErrors[1])
	}
}

func TestWriteError(t *testing.T) {
	logger := NewTestLogger(true)

	if !logger.Verbose {
		t.Errorf("Logger is not verbose but should")
	}
	iLogger := commonbl.Logger(logger)

	err1 := TestError{"just error 1"}
	iLogger.WriteError(&err1)
	if logger.GetErrorCount() != 1 {
		t.Errorf("Logger has ErrorCount '%d', but '1' is expected", logger.GetErrorCount())
	}

	if logger.GetMessageCount() != 0 {
		t.Errorf("Logger has MessageCount '%d', but '0' is expected", logger.GetMessageCount())
	}

	if logger.GetOutputCount() != 1 {
		t.Errorf("Logger has OutputCount '%d', but '1' is expected", logger.GetOutputCount())
	}

	err2 := TestError{"just error 2"}
	iLogger.WriteError(&err2)
	if logger.GetErrorCount() != 2 {
		t.Errorf("Logger has ErrorCount '%d', but '2' is expected", logger.GetErrorCount())
	}

	if logger.GetMessageCount() != 0 {
		t.Errorf("Logger has MessageCount '%d', but '0' is expected", logger.GetMessageCount())
	}

	if logger.GetOutputCount() != 2 {
		t.Errorf("Logger has OutputCount '%d', but '2' is expected", logger.GetOutputCount())
	}

	if logger.WrittenErrors[1] != "just error 2" {
		t.Errorf("The message '%s' is not the expected 'just error 2'", logger.WrittenErrors[1])
	}
}

func TestWriteErrorWithAddition(t *testing.T) {
	logger := NewTestLogger(true)

	if !logger.Verbose {
		t.Errorf("Logger is not verbose but should")
	}
	iLogger := commonbl.Logger(logger)

	err1 := TestError{"just error 1"}
	iLogger.WriteErrorWithAddition(&err1, "additional message 1")
	if logger.GetErrorCount() != 1 {
		t.Errorf("Logger has ErrorCount '%d', but '1' is expected", logger.GetErrorCount())
	}

	if logger.GetMessageCount() != 0 {
		t.Errorf("Logger has MessageCount '%d', but '0' is expected", logger.GetMessageCount())
	}

	if logger.GetOutputCount() != 1 {
		t.Errorf("Logger has OutputCount '%d', but '1' is expected", logger.GetOutputCount())
	}

	err2 := TestError{"just error 2"}
	iLogger.WriteErrorWithAddition(&err2, "additional message 2")
	if logger.GetErrorCount() != 2 {
		t.Errorf("Logger has ErrorCount '%d', but '2' is expected", logger.GetErrorCount())
	}

	if logger.GetMessageCount() != 0 {
		t.Errorf("Logger has MessageCount '%d', but '0' is expected", logger.GetMessageCount())
	}

	if logger.GetOutputCount() != 2 {
		t.Errorf("Logger has OutputCount '%d', but '2' is expected", logger.GetOutputCount())
	}

	if logger.WrittenErrors[1] != "just error 2 - additional message 2" {
		t.Errorf("The message '%s' is not the expected 'just error 2 - additional message 2'", logger.WrittenErrors[1])
	}
}

func TestOutPutCount(t *testing.T) {
	logger := NewTestLogger(true)

	if !logger.Verbose {
		t.Errorf("Logger is not verbose but should")
	}
	iLogger := commonbl.Logger(logger)
	iLogger.WriteErrorMessage("first error")
	iLogger.WriteVerbose("verbose message")
	iLogger.WriteInformation("info message")

	if logger.GetOutputCount() != 3 {
		t.Errorf("The OutputCount '%d' is not the expected '3'", logger.GetOutputCount())
	}
}
