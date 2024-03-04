package commonbl

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

const logfile_path = "./../../../../logs/file_logger_test.log"

func TestNewFileLogger(t *testing.T) {

	if logFileExists() {
		deleteTestsLogFile(t)
	}

	sut := NewFileLogger(true, logfile_path)

	if sut.Verbose != true {
		t.Errorf("FileLogger is not verbose, but should")
	}

	if sut.GetVerbose() != true {
		t.Errorf("FileLogger is not verbose, but should")
	}

	if sut.FullFilePath != logfile_path {
		t.Errorf("The FileLoggers FullFilePath is '%s' but should be '%s'", sut.FullFilePath, logfile_path)
	}

	iut := Logger(sut)
	if iut.GetVerbose() != true {
		t.Errorf("FileLogger is not verbose, but should")
	}

	if sut.infoLogger.Prefix() != "Information: " {
		t.Errorf("Infologger has prefix '%s', but 'Information: ' is expected", sut.infoLogger.Prefix())
	}

	if sut.errorLogger.Prefix() != "Error: " {
		t.Errorf("Infologger has prefix '%s', but 'Error: ' is expected", sut.errorLogger.Prefix())
	}

	if sut.verboseLogger.Prefix() != "Verbose: " {
		t.Errorf("Infologger has prefix '%s', but 'Verbose: ' is expected", sut.verboseLogger.Prefix())
	}

	if !logFileExists() {
		t.Errorf("Log file does not exist, but should after test was running")
	}

	fileLines := readLogFileLines()
	if len(fileLines) != 0 {
		t.Errorf("The log file has '%d' lines but '0' lines is expected", len(fileLines))
	}
}

func TestFileLoggerWriteInformation(t *testing.T) {
	if logFileExists() {
		deleteTestsLogFile(t)
	}

	sut := NewFileLogger(true, logfile_path)

	sut.WriteInformation("Some info message - 1")

	if sut.Verbose != true {
		t.Errorf("FileLogger is not verbose, but should")
	}

	sut.Verbose = false
	if sut.Verbose != false {
		t.Errorf("FileLogger is verbose, but should not")
	}

	sut.WriteInformation("Some info message - 2")

	if !logFileExists() {
		t.Errorf("Log file does not exist, but should after test was running")
	}

	fileLines := readLogFileLines()
	if len(fileLines) != 2 {
		t.Errorf("The log file has '%d' lines but '2' lines is expected", len(fileLines))
	}

}

func deleteTestsLogFile(t *testing.T) {
	err := os.Remove(logfile_path)
	if err != nil {
		t.Fatal(fmt.Println(fmt.Sprintf("Error '%s' when deleting file '%s'", err.Error(), logfile_path)))
	}
}

func logFileExists() bool {
	info, err := os.Stat(logfile_path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func readLogFileLines() []string {
	readFile, err := os.Open(logfile_path)

	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	readFile.Close()

	return fileLines
}
