package commonbl

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

const logfile_path = "./../../../../logs/file_logger_test.log"

var mutex = sync.Mutex{}

func TestNewFileLogger(t *testing.T) {

	mutex.Lock()
	defer mutex.Unlock()
	ensureLogFileDirExists()
	if logFileExists() {
		deleteTestsLogFile(t)
	}

	sut, _ := NewFileLogger(true, logfile_path)

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

func TestNewFileLoggerNotExistingDir(t *testing.T) {
	sut, err := NewFileLogger(true, "/dev/shm/not/existing/path/file.log")

	if sut != nil {
		t.Errorf("The 'FileLogger' should be nil, but is not")
	}

	if err == nil {
		t.Errorf("The 'error' should not be nil, but is")
	}

	if !strings.Contains(err.Error(), "/dev/shm/not/existing/path") {
		t.Errorf("The error message does not contain the expected data")
	}

	switch err.(type) {
	case *DirectoryNotExistError:
		fmt.Println("OK")
	default:
		t.Errorf("Got error of type '%s', but expected type '*DirectoryNotExistError'", err)
	}
}

func TestFileLoggerWriteInformation(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	ensureLogFileDirExists()
	if logFileExists() {
		deleteTestsLogFile(t)
	}

	sut, _ := NewFileLogger(true, logfile_path)
	infoMsg1 := "Some info message - 1"
	infoMsg2 := "Some info message - 2"
	infoMsg3 := "Some info message - 3"
	sut.WriteInformation(infoMsg1)

	if sut.Verbose != true {
		t.Errorf("FileLogger is not verbose, but should")
	}

	sut.Verbose = false
	if sut.Verbose != false {
		t.Errorf("FileLogger is verbose, but should not")
	}

	sut.WriteInformation(infoMsg2)
	iut := Logger(sut)
	iut.WriteInformation(infoMsg3)

	if !logFileExists() {
		t.Errorf("Log file does not exist, but should after test was running")
	}

	fileLines := readLogFileLines()
	if len(fileLines) != 3 {
		t.Errorf("The log file has '%d' lines but '3' lines is expected", len(fileLines))
	}

	expectedMsg1 := fmt.Sprintf("Information: %s", infoMsg1)
	if strings.HasSuffix(fileLines[0], expectedMsg1) == false {
		t.Errorf("The log on index '0' is '%s', but '%s' was expected", fileLines[0], expectedMsg1)
	}

	expectedMsg2 := fmt.Sprintf("Information: %s", infoMsg2)
	if strings.HasSuffix(fileLines[1], expectedMsg2) == false {
		t.Errorf("The log on index '1' is '%s', but '%s' was expected", fileLines[1], expectedMsg2)
	}

	expectedMsg3 := fmt.Sprintf("Information: %s", infoMsg3)
	if strings.HasSuffix(fileLines[2], expectedMsg3) == false {
		t.Errorf("The log on index '2' is '%s', but '%s' was expected", fileLines[2], expectedMsg3)
	}
}

func TestFileLoggerWriteVerbose(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	ensureLogFileDirExists()
	if logFileExists() {
		deleteTestsLogFile(t)
	}

	sut, _ := NewFileLogger(true, logfile_path)
	verboseMsg1 := "Some verbose message - 1"
	verboseMsg2 := "Some verbose message - 2"
	verboseMsg3 := "Some verbose message - 3"
	sut.WriteVerbose(verboseMsg1)

	if sut.Verbose != true {
		t.Errorf("FileLogger is not verbose, but should")
	}

	sut.Verbose = false
	if sut.Verbose != false {
		t.Errorf("FileLogger is verbose, but should not")
	}

	sut.WriteVerbose(verboseMsg2)
	if !logFileExists() {
		t.Errorf("Log file does not exist, but should after test was running")
	}

	fileLines := readLogFileLines()
	if len(fileLines) != 1 {
		t.Errorf("The log file has '%d' lines but '1' lines is expected", len(fileLines))
	}

	iut := Logger(sut)
	iut.WriteVerbose(verboseMsg3)
	fileLines = readLogFileLines()
	if len(fileLines) != 1 {
		t.Errorf("The log file has '%d' lines but '1' lines is expected", len(fileLines))
	}

	sut.Verbose = true
	if iut.GetVerbose() == false {
		t.Errorf("FileLogger is not verbose, but should")
	}
	sut.WriteVerbose(verboseMsg2)
	iut.WriteVerbose(verboseMsg3)

	fileLines = readLogFileLines()
	if len(fileLines) != 3 {
		t.Errorf("The log file has '%d' lines but '3' lines is expected", len(fileLines))
	}

	expectedMsg1 := fmt.Sprintf("Verbose: %s", verboseMsg1)
	if strings.HasSuffix(fileLines[0], expectedMsg1) == false {
		t.Errorf("The log on index '0' is '%s', but '%s' was expected", fileLines[0], expectedMsg1)
	}

	expectedMsg2 := fmt.Sprintf("Verbose: %s", verboseMsg2)
	if strings.HasSuffix(fileLines[1], expectedMsg2) == false {
		t.Errorf("The log on index '1' is '%s', but '%s' was expected", fileLines[1], expectedMsg2)
	}

	expectedMsg3 := fmt.Sprintf("Verbose: %s", verboseMsg3)
	if strings.HasSuffix(fileLines[2], expectedMsg3) == false {
		t.Errorf("The log on index '2' is '%s', but '%s' was expected", fileLines[2], expectedMsg3)
	}
}

func TestFileLoggerWriteInError(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	ensureLogFileDirExists()
	if logFileExists() {
		deleteTestsLogFile(t)
	}

	sut, _ := NewFileLogger(true, logfile_path)
	errorMsg1 := "Some error message - 1"
	errorMsg2 := "Some error message - 2"
	errorMsg3 := "Some error message - 3"
	additionalMsg := "More info on error"
	sut.WriteErrorMessage(errorMsg1)

	if sut.Verbose != true {
		t.Errorf("FileLogger is not verbose, but should")
	}

	sut.Verbose = false
	if sut.Verbose != false {
		t.Errorf("FileLogger is verbose, but should not")
	}

	sut.WriteError(NewWriterError(errorMsg2))
	iut := Logger(sut)
	iut.WriteErrorWithAddition(NewWriterError(errorMsg3), additionalMsg)

	if !logFileExists() {
		t.Errorf("Log file does not exist, but should after test was running")
	}

	fileLines := readLogFileLines()
	if len(fileLines) != 3 {
		t.Errorf("The log file has '%d' lines but '3' lines is expected", len(fileLines))
	}

	expectedMsg1 := fmt.Sprintf("Error: %s", errorMsg1)
	if strings.HasSuffix(fileLines[0], expectedMsg1) == false {
		t.Errorf("The log on index '0' is '%s', but '%s' was expected", fileLines[0], expectedMsg1)
	}

	expectedMsg2 := fmt.Sprintf("Error: The data \"%s\" was not written", errorMsg2)
	if strings.HasSuffix(fileLines[1], expectedMsg2) == false {
		t.Errorf("The log on index '1' is '%s', but '%s' was expected", fileLines[1], expectedMsg2)
	}

	expectedMsg3 := fmt.Sprintf("Error: The data \"%s\" was not written - %s", errorMsg3, additionalMsg)
	if strings.HasSuffix(fileLines[2], expectedMsg3) == false {
		t.Errorf("The log on index '2' is '%s', but '%s' was expected", fileLines[2], expectedMsg3)
	}
}

func TestFileLoggerWriteMixed(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	ensureLogFileDirExists()
	if logFileExists() {
		deleteTestsLogFile(t)
	}

	sut, _ := NewFileLogger(true, logfile_path)
	infoMsg1 := "Some info message - 1"
	verboseMsg2 := "Some verbose message - 2"
	errorMsg3 := "Some error message - 3"
	sut.WriteInformation(infoMsg1)

	if sut.Verbose != true {
		t.Errorf("FileLogger is not verbose, but should")
	}

	sut.Verbose = false
	if sut.Verbose != false {
		t.Errorf("FileLogger is verbose, but should not")
	}

	sut.WriteVerbose(verboseMsg2)
	if !logFileExists() {
		t.Errorf("Log file does not exist, but should after test was running")
	}

	fileLines := readLogFileLines()
	if len(fileLines) != 1 {
		t.Errorf("The log file has '%d' lines but '1' lines is expected", len(fileLines))
	}

	iut := Logger(sut)
	iut.WriteVerbose(verboseMsg2)
	fileLines = readLogFileLines()
	if len(fileLines) != 1 {
		t.Errorf("The log file has '%d' lines but '1' lines is expected", len(fileLines))
	}

	sut.Verbose = true
	if iut.GetVerbose() == false {
		t.Errorf("FileLogger is not verbose, but should")
	}
	iut.WriteVerbose(verboseMsg2)
	iut.WriteErrorMessage(errorMsg3)

	fileLines = readLogFileLines()
	if len(fileLines) != 3 {
		t.Errorf("The log file has '%d' lines but '3' lines is expected", len(fileLines))
	}

	expectedMsg1 := fmt.Sprintf("Information: %s", infoMsg1)
	if strings.HasSuffix(fileLines[0], expectedMsg1) == false {
		t.Errorf("The log on index '0' is '%s', but '%s' was expected", fileLines[0], expectedMsg1)
	}

	expectedMsg2 := fmt.Sprintf("Verbose: %s", verboseMsg2)
	if strings.HasSuffix(fileLines[1], expectedMsg2) == false {
		t.Errorf("The log on index '1' is '%s', but '%s' was expected", fileLines[1], expectedMsg2)
	}

	expectedMsg3 := fmt.Sprintf("Error: %s", errorMsg3)
	if strings.HasSuffix(fileLines[2], expectedMsg3) == false {
		t.Errorf("The log on index '2' is '%s', but '%s' was expected", fileLines[2], expectedMsg3)
	}
}

func TestDirectoryExists(t *testing.T) {
	if directoryExists("/bin") == false {
		t.Errorf("'directoryExists' tells '/bin' does not exist")
	}

	if directoryExists("/bin/") == false {
		t.Errorf("'directoryExists' tells '/bin/' does not exist")
	}

	if directoryExists("/dev/shm/not/existing/path/") == true {
		t.Errorf("'directoryExists' tells '/dev/shm/not/existing/path/' does exist")
	}

	os.OpenFile(logfile_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if directoryExists(logfile_path) == true {
		t.Errorf("'directoryExists' tells '%s' does exist but it is a file!", logfile_path)
	}

}

func deleteTestsLogFile(t *testing.T) {
	if !logFileExists() {
		return
	}

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

func ensureLogFileDirExists() {
	logFileDir := filepath.Dir(logfile_path)
	if !directoryExists(logFileDir) {
		os.MkdirAll(logFileDir, os.ModePerm)
	}
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
