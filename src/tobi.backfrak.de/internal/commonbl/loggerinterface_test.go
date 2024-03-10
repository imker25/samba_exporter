package commonbl

import (
	"fmt"
	"testing"
)

func TestGetLoggerForFileLogger(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	ensureLogFileDirExists()
	logger, err := GetLogger(logfile_path, false)

	if err != nil {
		t.Errorf("Got error '%s' but expected none", err.Error())
	}

	if logger.GetVerbose() == true {
		t.Errorf("The logger is verbose, but should not")
	}

	switch logger.(type) {
	case *FileLogger:
		fmt.Println("OK")
	default:
		t.Errorf("The logger is not the expected FileLogger")
	}

	logger2, err2 := GetLogger(logfile_path, true)
	if err2 != nil {
		t.Errorf("Got error '%s' but expected none", err2.Error())
	}

	if logger2.GetVerbose() == false {
		t.Errorf("The logger is not verbose, but should ")
	}

	switch logger2.(type) {
	case *FileLogger:
		fmt.Println("OK")
	default:
		t.Errorf("The logger is not the expected FileLogger")
	}
}

func TestGetLoggerForConsoleLogger(t *testing.T) {
	logger1, err1 := GetLogger(" ", false)

	if err1 != nil {
		t.Errorf("Got error '%s' but expected none", err1.Error())
	}

	if logger1.GetVerbose() == true {
		t.Errorf("The logger is verbose, but should not")
	}

	switch logger1.(type) {
	case *ConsoleLogger:
		fmt.Println("OK")
	default:
		t.Errorf("The logger is not the expected ConsoleLogger")
	}

	logger2, err2 := GetLogger(" ", true)
	if err2 != nil {
		t.Errorf("Got error '%s' but expected none", err2.Error())
	}

	if logger2.GetVerbose() == false {
		t.Errorf("The logger is not verbose, but should ")
	}

	switch logger2.(type) {
	case *ConsoleLogger:
		fmt.Println("OK")
	default:
		t.Errorf("The logger is not the expected ConsoleLogger")
	}
}
