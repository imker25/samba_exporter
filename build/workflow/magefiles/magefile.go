//go:build mage
// +build mage

package main

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
)

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = Build

const VERSION_FILE = "VersionMaster.txt"
const PROJECT_NAME_SPACE = "tobi.backfrak.de"

type buildContext struct {
	GitHight             int
	GitHash              string
	ProgramVersion       string
	ProgramVersionNumber string
	DebPackageVersion    string
	DebPackageName       string
	ShortVersion         string
	BinDir               string
	PackageDir           string
	LogDir               string
	WorkDir              string
	SourceDir            string
	TmpDir               string
	PackagesToBuild      []string
	PackagesToTest       []string
	VersionFilePath      string
}

var smbExportBuildContext buildContext

func getEnvironment() error {
	fmt.Println(fmt.Sprintf("Get the environment for the build..."))
	fmt.Println(fmt.Sprintf("# ##############################################################################################"))

	workDir, errWorkDir := os.Getwd()
	if errWorkDir != nil {
		return errWorkDir
	}

	if strings.HasSuffix(workDir, "build/workflow/") || strings.HasSuffix(workDir, "build/workflow") {
		workDir = filepath.Join(workDir, "..", "..")
	}

	smbExportBuildContext.WorkDir = workDir
	smbExportBuildContext.BinDir = filepath.Join(workDir, "bin")
	smbExportBuildContext.LogDir = filepath.Join(workDir, "logs")
	smbExportBuildContext.PackageDir = filepath.Join(workDir, "pkg")
	smbExportBuildContext.TmpDir = filepath.Join(workDir, "tmp")
	smbExportBuildContext.SourceDir = filepath.Join(workDir, "src")
	smbExportBuildContext.VersionFilePath = filepath.Join(smbExportBuildContext.WorkDir, VERSION_FILE)

	hash, errHash := getGitHash()
	if errHash != nil {
		return errHash
	}
	// fmt.Println(fmt.Sprintf("Git Hash: %s", hash))
	smbExportBuildContext.GitHash = hash

	hight, errHight := getGitHight()
	if errHight != nil {
		return errHight
	}
	// fmt.Println(fmt.Sprintf("Git Hight: %d", hight))
	smbExportBuildContext.GitHight = hight

	givenVersion, errVersion := readVersionMaster()
	if errVersion != nil {
		return errVersion
	}

	smbExportBuildContext.ShortVersion = givenVersion
	smbExportBuildContext.ProgramVersionNumber = fmt.Sprintf("%s.%d", givenVersion, hight)
	smbExportBuildContext.ProgramVersion = fmt.Sprintf("%s.%d-%s", givenVersion, hight, hash)
	debPackVersion := smbExportBuildContext.ProgramVersion
	if os.Getenv("GITHUB_RUNNER_OS") != "" {
		debPackVersion = fmt.Sprintf("%s+%s", smbExportBuildContext.ProgramVersion, os.Getenv("GITHUB_RUNNER_OS"))
	}
	smbExportBuildContext.DebPackageVersion = debPackVersion
	smbExportBuildContext.DebPackageName = fmt.Sprintf("samba-exporter_%s", debPackVersion)
	fmt.Println(fmt.Sprintf("Run samba-exporter build workflow for V%s", smbExportBuildContext.ProgramVersion))

	errFindBuild := filepath.Walk(filepath.Join(smbExportBuildContext.SourceDir, PROJECT_NAME_SPACE, "cmd"), func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return nil
		}

		packToBuild := filepath.Dir(path)
		if !info.IsDir() && filepath.Base(path) == "go.mod" && !listContains(smbExportBuildContext.PackagesToBuild, packToBuild) {
			smbExportBuildContext.PackagesToBuild = append(smbExportBuildContext.PackagesToBuild, packToBuild)
		}

		return nil
	})
	if errFindBuild != nil {
		return errFindBuild
	}

	errFindTest := filepath.Walk(filepath.Join(smbExportBuildContext.SourceDir, PROJECT_NAME_SPACE), func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return nil
		}

		if !info.IsDir() && filepath.Ext(path) == ".go" {
			packToTest := filepath.Dir(path)
			if strings.HasSuffix(path, "_test.go") && !listContains(smbExportBuildContext.PackagesToTest, packToTest) {
				smbExportBuildContext.PackagesToTest = append(smbExportBuildContext.PackagesToTest, packToTest)
			}
		}

		return nil
	})
	if errFindTest != nil {
		return errFindTest
	}

	fmt.Println(fmt.Sprintf("# ##############################################################################################"))
	return nil
}

// Get the build name files
func GetBuildName() error {
	mg.Deps(getEnvironment, Clean)
	fmt.Println(fmt.Sprintf("Create samba-exporter Version files..."))
	fmt.Println(fmt.Sprintf("# ##############################################################################################"))

	if _, err := os.Stat(smbExportBuildContext.LogDir); os.IsNotExist(err) {
		errCreate := os.Mkdir(smbExportBuildContext.LogDir, 0755)
		if errCreate != nil {
			return errCreate
		}
	}

	errNr := ioutil.WriteFile(filepath.Join(smbExportBuildContext.LogDir, "PackageName.txt"), []byte(smbExportBuildContext.DebPackageName), 0644)
	if errNr != nil {
		return errNr
	}

	errVersion := ioutil.WriteFile(filepath.Join(smbExportBuildContext.LogDir, "Version.txt"), []byte(smbExportBuildContext.DebPackageVersion), 0644)
	if errVersion != nil {
		return errVersion
	}

	errShort := ioutil.WriteFile(filepath.Join(smbExportBuildContext.LogDir, "ShortVersion.txt"), []byte(smbExportBuildContext.ProgramVersionNumber), 0644)
	if errShort != nil {
		return errShort
	}

	fmt.Println(fmt.Sprintf("# ##############################################################################################"))
	return nil
}

// Compiles the project
func Build() error {
	mg.Deps(getEnvironment, Clean, GetBuildName)
	fmt.Println(fmt.Sprintf("Building samba-exporter V%s ...", smbExportBuildContext.ProgramVersion))
	fmt.Println(fmt.Sprintf("# ##############################################################################################"))

	if _, err := os.Stat(smbExportBuildContext.BinDir); os.IsNotExist(err) {
		errCreate := os.Mkdir(smbExportBuildContext.BinDir, 0755)
		if errCreate != nil {
			return errCreate
		}
	}

	for _, packToBuild := range smbExportBuildContext.PackagesToBuild {
		outPutPath := filepath.Join(smbExportBuildContext.BinDir, filepath.Base(packToBuild))
		fmt.Println(fmt.Sprintf("Compile package '%s' to '%s'", packToBuild, outPutPath))

		ldfFlags := fmt.Sprintf("-X main.version=%s", smbExportBuildContext.ProgramVersion)
		fmt.Println(fmt.Sprintf("Run in %s: %s %s %s %s %s -ldflags=\"%s\"", packToBuild, "go", "build", "-o", outPutPath, "-v", ldfFlags))
		cmd := exec.Command("go", "build", "-o", outPutPath, "-v", "-ldflags", ldfFlags)
		cmd.Dir = packToBuild
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		errBuild := cmd.Run()

		if errBuild != nil {
			return errBuild
		}
	}

	fmt.Println(fmt.Sprintf("# ##############################################################################################"))
	return nil
}

// Runs the tests for the project
func Test() error {
	mg.Deps(getEnvironment, Clean, GetBuildName, installTestDeps)
	fmt.Println(fmt.Sprintf("Testing samba-exporter... "))
	fmt.Println(fmt.Sprintf("# ##############################################################################################"))

	if _, err := os.Stat(smbExportBuildContext.LogDir); os.IsNotExist(err) {
		errCreate := os.Mkdir(smbExportBuildContext.LogDir, 0755)
		if errCreate != nil {
			return errCreate
		}
	}

	logPath := filepath.Join(smbExportBuildContext.LogDir, "TestsRun.log")
	xmlResult := filepath.Join(smbExportBuildContext.LogDir, "TestsResult.xml")
	logFile, errOpen := os.Create(logPath)
	if errOpen != nil {
		return errOpen
	}
	defer logFile.Close()

	testErrors := []error{}
	for _, packToTest := range smbExportBuildContext.PackagesToTest {

		fmt.Println(fmt.Sprintf("Test package '%s', logging to '%s'", packToTest, logPath))
		fmt.Println(fmt.Sprintf("Run in %s: %s %s %s %s >> %s", packToTest, "go", "test", "-v", "-race", logPath))
		cmd := exec.Command("go", "test", "-v", "-race")
		cmd.Dir = packToTest
		cmd.Stderr = logFile
		cmd.Stdout = logFile
		errTest := cmd.Run()
		if errTest != nil {
			fmt.Println(errTest.Error())
			testErrors = append(testErrors, errTest)
		}
	}

	errConv := convertTestResults(logPath, xmlResult)
	if errConv != nil {
		return errConv
	}

	if len(testErrors) > 0 {
		return testErrors[0]
	}

	fmt.Println(fmt.Sprintf("# ##############################################################################################"))
	return nil
}

// Runs test coverage for the project
func Cover() error {
	mg.Deps(getEnvironment, Clean, GetBuildName, installTestDeps)
	fmt.Println(fmt.Sprintf("Testing samba-exporter... "))
	fmt.Println(fmt.Sprintf("# ##############################################################################################"))

	if _, err := os.Stat(smbExportBuildContext.LogDir); os.IsNotExist(err) {
		errCreate := os.Mkdir(smbExportBuildContext.LogDir, 0755)
		if errCreate != nil {
			return errCreate
		}
	}

	logPath := filepath.Join(smbExportBuildContext.LogDir, "TestsCoverRun.log")
	logFile, errOpen := os.Create(logPath)
	if errOpen != nil {
		return errOpen
	}

	for _, packToTest := range smbExportBuildContext.PackagesToTest {

		fmt.Println(fmt.Sprintf("Test package '%s', logging to '%s'", packToTest, logPath))
		fmt.Println(fmt.Sprintf("Run in %s: %s %s %s %s >> %s", packToTest, "go", "test", "-v", "-cover", logPath))
		cmd := exec.Command("go", "test", "-v", "-cover")

		cmd.Dir = packToTest
		cmd.Stderr = logFile
		cmd.Stdout = logFile
		errTest := cmd.Run()
		if errTest != nil {
			logFile.Close()
			fmt.Println(errTest.Error())
			return errTest
		}
	}
	logFile.Close()

	fmt.Println(fmt.Sprintf("# ##############################################################################################"))
	return nil
}

// Remove all build output
func Clean() error {
	mg.Deps(getEnvironment)
	fmt.Println("Cleaning...")
	fmt.Println(fmt.Sprintf("# ##############################################################################################"))

	err := removePaths([]string{
		smbExportBuildContext.BinDir,
		smbExportBuildContext.PackageDir,
		smbExportBuildContext.LogDir,
		smbExportBuildContext.TmpDir,
	})
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("# ##############################################################################################"))
	return nil
}

// Creates a folder that includes all files needed within a distribution specific install package.
func PreparePack() error {
	mg.Deps(Build)
	fmt.Println("Prepare packing...")
	fmt.Println(fmt.Sprintf("# ##############################################################################################"))
	debPackageDir := filepath.Join(smbExportBuildContext.TmpDir, smbExportBuildContext.DebPackageName)

	dist, errDist := readOSDistribution()
	if errDist != nil {
		return errDist
	}

	fmt.Println(fmt.Sprintf("Copy the files needed by the package to '%s'", debPackageDir))
	if _, err := os.Stat(smbExportBuildContext.TmpDir); os.IsNotExist(err) {
		errCreate := os.Mkdir(smbExportBuildContext.TmpDir, 0755)
		if errCreate != nil {
			return errCreate
		}
	}
	if _, err := os.Stat(debPackageDir); os.IsNotExist(err) {
		errCreate := os.Mkdir(debPackageDir, 0755)
		if errCreate != nil {
			return errCreate
		}
	}

	if dist == "ubuntu" || dist == "debian" {
		fmt.Println("Copy deb package specific files")
		debFilesPath := filepath.Join(debPackageDir, "DEBIAN")
		if _, err := os.Stat(debFilesPath); os.IsNotExist(err) {
			errCreate := os.Mkdir(debFilesPath, 0755)
			if errCreate != nil {
				return errCreate
			}
		}

		cpDebCmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("cp -rf \"%s\"/* \"%s\"",
			filepath.Join(smbExportBuildContext.WorkDir, "install", "debian"),
			fmt.Sprintf("%s/", debFilesPath)))
		cpDebCmd.Stdout = os.Stdout
		cpDebCmd.Stderr = os.Stderr
		errCopyDeb := cpDebCmd.Run()
		if errCopyDeb != nil {
			return errCopyDeb
		}

	} else if dist == "fedora" {
		fmt.Println("Copy rpm package specific files")
		cpSpecCmd := exec.Command("cp", filepath.Join(smbExportBuildContext.WorkDir, "install", "fedora", "samba-exporter.from_gradle.spec"),
			filepath.Join(debPackageDir, "samba-exporter.spec"))
		cpSpecCmd.Stdout = os.Stdout
		cpSpecCmd.Stderr = os.Stderr
		errCopySpec := cpSpecCmd.Run()
		if errCopySpec != nil {
			return errCopySpec
		}
	} else {
		fmt.Println("Warning: Distribution unknown, no distribution specific files copied to the package tmp folder")
	}

	fmt.Println("Copy package files")
	installCmd := exec.Command("./build/InstallProgram.sh", smbExportBuildContext.WorkDir, smbExportBuildContext.BinDir, debPackageDir, smbExportBuildContext.ShortVersion)
	installCmd.Dir = smbExportBuildContext.WorkDir
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	errInstall := installCmd.Run()
	if errInstall != nil {
		fmt.Println(errInstall.Error())
		return errInstall
	}

	cpLicenseCmd := exec.Command("cp", filepath.Join(smbExportBuildContext.WorkDir, "LICENSE"),
		filepath.Join(debPackageDir, "usr", "share", "doc", "samba-exporter"))
	cpLicenseCmd.Stdout = os.Stdout
	cpLicenseCmd.Stderr = os.Stderr
	errCopyLicense := cpLicenseCmd.Run()
	if errCopyLicense != nil {
		return errCopyLicense
	}

	fmt.Println(fmt.Sprintf("# ##############################################################################################"))
	return nil
}

func installTestDeps() error {
	mg.Deps(Clean)
	fmt.Println("Installing Test Dependencies...")
	fmt.Println(fmt.Sprintf("# ##############################################################################################"))
	cmd := exec.Command("go", "install", "-v", "github.com/tebeka/go2xunit@v1.4.10")
	cmd.Dir = filepath.Join(smbExportBuildContext.WorkDir, "build")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdin
	err := cmd.Run()
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("# ##############################################################################################"))
	return nil
}

func getGitHash() (string, error) {
	cmd := exec.Command("git", "describe", "--always", "--long", "--dirty")
	cmd.Dir = smbExportBuildContext.WorkDir
	cmd.Stderr = os.Stderr
	hash, err := cmd.Output()
	if err != nil {
		return "", err
	}
	hashStr := strings.TrimSpace(string(hash))
	return hashStr, nil
}

func getGitHight() (int, error) {
	cmd := exec.Command("git", "log", "--pretty=format:\"%H\"", "-n 1", "--follow", VERSION_FILE)
	cmd.Dir = smbExportBuildContext.WorkDir
	cmd.Stderr = os.Stderr
	lastChange, errLast := cmd.Output()
	if errLast != nil {
		return -1, errLast
	}
	lastChangeStr := strings.ReplaceAll(strings.TrimSpace(string(lastChange)), "\"", "")

	cmd = exec.Command("git", "log", "--pretty=format:\"%H\"", "-n 1")
	cmd.Dir = smbExportBuildContext.WorkDir
	cmd.Stderr = os.Stderr
	head, errHead := cmd.Output()
	if errHead != nil {
		return -1, errHead
	}

	headStr := strings.ReplaceAll(strings.TrimSpace(string(head)), "\"", "")

	cmd = exec.Command("git", "rev-list", "--count", lastChangeStr+".."+headStr)
	cmd.Dir = smbExportBuildContext.WorkDir
	cmd.Stderr = os.Stderr
	hight, hightErr := cmd.Output()
	if hightErr != nil {
		return -1, hightErr
	}

	hightStr := strings.TrimSpace(string(hight))
	hightInt, errCon := strconv.Atoi(hightStr)
	if errCon != nil {
		return -1, nil
	}

	return hightInt, nil
}

func readVersionMaster() (string, error) {
	content, err := ioutil.ReadFile(smbExportBuildContext.VersionFilePath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(content)), nil
}

func listContains(list []string, value string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}

	return false
}

func removePaths(paths []string) error {
	for _, path := range paths {
		err := os.RemoveAll(path)
		if err != nil {
			return err
		}
	}
	return nil
}

func readOSDistribution() (string, error) {
	ret := ""
	byteContent, err := ioutil.ReadFile("/etc/os-release")
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(byteContent), "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "ID=") {
			ret = strings.Replace(line, "ID=", "", 1)
		}
	}

	return ret, nil
}

func convertTestResults(logPath, xmlResult string) error {
	fmt.Println(fmt.Sprintf("Convert the test results %s to %s", logPath, xmlResult))
	cmd := exec.Command("go", "run", "github.com/tebeka/go2xunit", "-input", logPath, "-output", xmlResult)
	cmd.Dir = filepath.Join(smbExportBuildContext.WorkDir, "build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	errConvert := cmd.Run()
	if errConvert != nil {
		return errConvert
	}

	return nil
}
