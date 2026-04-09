package commonbl

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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

	if sut.LogLevelSetting != Verbose {
		t.Errorf("FileLogger is not verbose, but should")
	}

	if sut.GetVerbose() != true {
		t.Errorf("FileLogger is not verbose, but should")
	}

	if sut.GetLogLevelSetting() != Verbose {
		t.Errorf("FileLogger is not verbose, but should")
	}

	if sut.FullFilePath != logfile_path {
		t.Errorf("The FileLoggers FullFilePath is '%s' but should be '%s'", sut.FullFilePath, logfile_path)
	}

	iut := Logger(sut)
	if iut.GetVerbose() != true {
		t.Errorf("FileLogger is not verbose, but should")
	}

	if sut.internalLoggers[Information].Prefix() != "Information: " {
		t.Errorf("Infologger has prefix '%s', but 'Information: ' is expected", sut.internalLoggers[Information].Prefix())
	}

	if sut.internalLoggers[Error].Prefix() != "Error: " {
		t.Errorf("Infologger has prefix '%s', but 'Error: ' is expected", sut.internalLoggers[Error].Prefix())
	}

	if sut.internalLoggers[Verbose].Prefix() != "Verbose: " {
		t.Errorf("Infologger has prefix '%s', but 'Verbose: ' is expected", sut.internalLoggers[Verbose].Prefix())
	}

	if sut.internalLoggers[Warning].Prefix() != "Warning: " {
		t.Errorf("Infologger has prefix '%s', but 'Warning: ' is expected", sut.internalLoggers[Warning].Prefix())
	}

	if !logFileExists() {
		t.Errorf("Log file does not exist, but should after test was running")
	}

	fileLines := readLogFileLines()
	if len(fileLines) != 0 {
		t.Errorf("The log file has '%d' lines but '0' lines is expected", len(fileLines))
	}
}

func TestNewFileLogger2Verbose(t *testing.T) {

	mutex.Lock()
	defer mutex.Unlock()
	ensureLogFileDirExists()
	if logFileExists() {
		deleteTestsLogFile(t)
	}

	sut, _ := NewFileLogger2(Verbose, logfile_path)

	if sut.LogLevelSetting != Verbose {
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

	if sut.internalLoggers[Information].Prefix() != "Information: " {
		t.Errorf("Infologger has prefix '%s', but 'Information: ' is expected", sut.internalLoggers[Information].Prefix())
	}

	if sut.internalLoggers[Error].Prefix() != "Error: " {
		t.Errorf("Infologger has prefix '%s', but 'Error: ' is expected", sut.internalLoggers[Error].Prefix())
	}

	if sut.internalLoggers[Verbose].Prefix() != "Verbose: " {
		t.Errorf("Infologger has prefix '%s', but 'Verbose: ' is expected", sut.internalLoggers[Verbose].Prefix())
	}

	if sut.internalLoggers[Warning].Prefix() != "Warning: " {
		t.Errorf("Infologger has prefix '%s', but 'Warning: ' is expected", sut.internalLoggers[Warning].Prefix())
	}

	if !logFileExists() {
		t.Errorf("Log file does not exist, but should after test was running")
	}

	fileLines := readLogFileLines()
	if len(fileLines) != 0 {
		t.Errorf("The log file has '%d' lines but '0' lines is expected", len(fileLines))
	}
}

func TestNewFileLogger2Information(t *testing.T) {

	mutex.Lock()
	defer mutex.Unlock()
	ensureLogFileDirExists()
	if logFileExists() {
		deleteTestsLogFile(t)
	}

	sut, _ := NewFileLogger2(Information, logfile_path)

	if sut.LogLevelSetting != Information {
		t.Errorf("FileLogger is not verbose, but should")
	}

	if sut.GetLogLevelSetting() != Information {
		t.Errorf("FileLogger is not verbose, but should")
	}

	if sut.GetVerbose() != false {
		t.Errorf("FileLogger is not verbose, but should")
	}

	if sut.FullFilePath != logfile_path {
		t.Errorf("The FileLoggers FullFilePath is '%s' but should be '%s'", sut.FullFilePath, logfile_path)
	}

	iut := Logger(sut)
	if iut.GetVerbose() != false {
		t.Errorf("FileLogger is not verbose, but should")
	}

	if sut.internalLoggers[Information].Prefix() != "Information: " {
		t.Errorf("Infologger has prefix '%s', but 'Information: ' is expected", sut.internalLoggers[Information].Prefix())
	}

	if sut.internalLoggers[Error].Prefix() != "Error: " {
		t.Errorf("Infologger has prefix '%s', but 'Error: ' is expected", sut.internalLoggers[Error].Prefix())
	}

	if sut.internalLoggers[Verbose].Prefix() != "Verbose: " {
		t.Errorf("Infologger has prefix '%s', but 'Verbose: ' is expected", sut.internalLoggers[Verbose].Prefix())
	}

	if sut.internalLoggers[Warning].Prefix() != "Warning: " {
		t.Errorf("Infologger has prefix '%s', but 'Warning: ' is expected", sut.internalLoggers[Warning].Prefix())
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

	if sut.GetVerbose() != true {
		t.Errorf("FileLogger is not verbose, but should")
	}

	sut.LogLevelSetting = Information
	if sut.GetVerbose() != false {
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

	if !runningOnUbuntuFocal() {
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
}

func TestFileLoggerWriteWarning(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	ensureLogFileDirExists()
	if logFileExists() {
		deleteTestsLogFile(t)
	}

	sut, _ := NewFileLogger(true, logfile_path)
	warnMsg1 := "Some warning message - 1"
	warnMsg2 := "Some warning message - 2"
	// warnMsg3 := "Some warning message - 3"
	sut.WriteWarning(warnMsg1)

	if sut.GetVerbose() != true {
		t.Errorf("FileLogger is not verbose, but should")
	}

	sut.LogLevelSetting = Information
	if sut.GetVerbose() != false {
		t.Errorf("FileLogger is verbose, but should not")
	}

	sut.WriteWarning(warnMsg2)
	// iut := Logger(sut)
	// iut.WriteInformation(warnMsg3)

	if !logFileExists() {
		t.Errorf("Log file does not exist, but should after test was running")
	}

	fileLines := readLogFileLines()
	if len(fileLines) != 2 {
		t.Errorf("The log file has '%d' lines but '2' lines is expected", len(fileLines))
	}

	if !runningOnUbuntuFocal() {
		expectedMsg1 := fmt.Sprintf("Warning: %s", warnMsg1)
		if strings.HasSuffix(fileLines[0], expectedMsg1) == false {
			t.Errorf("The log on index '0' is '%s', but '%s' was expected", fileLines[0], expectedMsg1)
		}

		expectedMsg2 := fmt.Sprintf("Warning: %s", warnMsg2)
		if strings.HasSuffix(fileLines[1], expectedMsg2) == false {
			t.Errorf("The log on index '1' is '%s', but '%s' was expected", fileLines[1], expectedMsg2)
		}

		// expectedMsg3 := fmt.Sprintf("Warning: %s", warnMsg3)
		// if strings.HasSuffix(fileLines[2], expectedMsg3) == false {
		// 	t.Errorf("The log on index '2' is '%s', but '%s' was expected", fileLines[2], expectedMsg3)
		// }
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

	if sut.GetVerbose() != true {
		t.Errorf("FileLogger is not verbose, but should")
	}

	sut.LogLevelSetting = Information
	if sut.GetVerbose() != false {
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

	sut.LogLevelSetting = Verbose
	if iut.GetVerbose() == false {
		t.Errorf("FileLogger is not verbose, but should")
	}
	sut.WriteVerbose(verboseMsg2)
	iut.WriteVerbose(verboseMsg3)

	fileLines = readLogFileLines()
	if len(fileLines) != 3 {
		t.Errorf("The log file has '%d' lines but '3' lines is expected", len(fileLines))
	}

	if !runningOnUbuntuFocal() {
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

	if sut.GetVerbose() != true {
		t.Errorf("FileLogger is not verbose, but should")
	}

	sut.LogLevelSetting = Information
	if sut.GetVerbose() != false {
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

	if !runningOnUbuntuFocal() {
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

	if sut.GetVerbose() != true {
		t.Errorf("FileLogger is not verbose, but should")
	}

	sut.LogLevelSetting = Information
	if sut.GetVerbose() != false {
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

	sut.LogLevelSetting = Verbose
	if iut.GetVerbose() == false {
		t.Errorf("FileLogger is not verbose, but should")
	}
	iut.WriteVerbose(verboseMsg2)
	iut.WriteErrorMessage(errorMsg3)

	fileLines = readLogFileLines()
	if len(fileLines) != 3 {
		t.Errorf("The log file has '%d' lines but '3' lines is expected", len(fileLines))
	}

	if !runningOnUbuntuFocal() {
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

func TestErrorLogLevel(t *testing.T) {
	testLogLevels(t, Error)
}

func TestWarningLogLevel(t *testing.T) {
	testLogLevels(t, Warning)
}

func TestInformationLogLevel(t *testing.T) {
	testLogLevels(t, Information)
}

func TestVerboseLogLevel(t *testing.T) {
	testLogLevels(t, Verbose)
}

func testLogLevels(t *testing.T, logLevelSetting LogLevel) {

	mutex.Lock()
	defer mutex.Unlock()
	ensureLogFileDirExists()
	if logFileExists() {
		deleteTestsLogFile(t)
	}

	sut, _ := NewFileLogger2(logLevelSetting, logfile_path)
	infoMsg1 := "Some info message - 1"
	verboseMsg2 := "Some verbose message - 2"
	errorMsg3 := "Some error message - 3"
	warningMsg4 := "Some Warning message - 4"
	sut.WriteInformation(infoMsg1)
	sut.WriteVerbose(verboseMsg2)
	sut.WriteErrorMessage(errorMsg3)
	sut.WriteWarning(warningMsg4)

	fileLines := readLogFileLines()
	if len(fileLines) != int(sut.GetLogLevelSetting())+1 {
		t.Errorf("The number of logged messages '%d' is not the expected '%d'", len(fileLines), sut.GetLogLevelSetting())
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
	return readFileLines(logfile_path)
}

func runningOnUbuntuFocal() bool {
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
		return false
	}

	lines := readFileLines("/etc/os-release")
	for _, line := range lines {
		if line == "VERSION_CODENAME=focal" {
			return true
		}
	}

	return false
}

func readFileLines(path string) []string {
	readFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	return fileLines
}
