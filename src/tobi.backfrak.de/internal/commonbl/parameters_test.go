package commonbl

import (
	"fmt"
	"testing"
)

func TestGetLogLevelSettingNotExistError(t *testing.T) {
	parma := Parmeters{}
	parma.LogLevelSetting = "Bla"

	logLevel, err := parma.GetLogLevelSetting()
	if logLevel != -1 {
		t.Errorf("The logLevel for 'Bla' is expected to be '-1' but is '%d'", logLevel)
	}

	if err == nil {
		t.Errorf("Got no error, but expected one")
	}

	switch err.(type) {
	case *LogLevelNotDefinedError:
		fmt.Println("OK")
	default:
		t.Errorf("The error is not the expected LogLevelNotDefinedError")
	}

}

func TestGetLogLevelSettingVerbose(t *testing.T) {
	parma := Parmeters{}
	parma.LogLevelSetting = "Verbose"

	logLevel, err := parma.GetLogLevelSetting()
	if logLevel != Verbose {
		t.Errorf("The logLevel for 'Verbose' is expected to be 'Verbose' but is '%d'", logLevel)
	}

	if err != nil {
		t.Errorf("Got error, but expected none")
	}

}

func TestGetLogLevelSettingError(t *testing.T) {
	parma := Parmeters{}
	parma.LogLevelSetting = "Error"

	logLevel, err := parma.GetLogLevelSetting()
	if logLevel != Error {
		t.Errorf("The logLevel for 'Error' is expected to be 'Error' but is '%d'", logLevel)
	}

	if err != nil {
		t.Errorf("Got error, but expected none")
	}

}

func TestGetLogLevelSettingDefault(t *testing.T) {
	parma := Parmeters{}
	parma.LogLevelSetting = " "

	logLevel, err := parma.GetLogLevelSetting()
	if logLevel != Information {
		t.Errorf("The logLevel for 'Information' is expected to be 'Error' but is '%d'", logLevel)
	}

	if err != nil {
		t.Errorf("Got error, but expected none")
	}

}

func TestGetLogLevelSettingVerboseSetting(t *testing.T) {
	parma := Parmeters{}
	parma.Verbose = true
	parma.LogLevelSetting = "Information"

	logLevel, err := parma.GetLogLevelSetting()
	if logLevel != Verbose {
		t.Errorf("The logLevel for 'Verbose' is expected to be 'Error' but is '%d'", logLevel)
	}

	if err != nil {
		t.Errorf("Got error, but expected none")
	}

}
