package commonbl

import "testing"

func TestNewFileLogger(t *testing.T) {
	logFileName := "mý_log_file.log"
	sut := NewFileLogger(true, logFileName)

	if sut.Verbose != true {
		t.Errorf("FileLogger is not verbose, but should")
	}

	if sut.GetVerbose() != true {
		t.Errorf("FileLogger is not verbose, but should")
	}

	if sut.FullFilePath != logFileName {
		t.Errorf("The FileLoggers FullFilePath is '%s' but should be '%s'", sut.FullFilePath, logFileName)
	}

	iut := Logger(sut)
	if iut.GetVerbose() != true {
		t.Errorf("FileLogger is not verbose, but should")
	}
}
