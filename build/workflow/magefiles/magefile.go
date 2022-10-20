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
	"strings"

	"github.com/imker25/gobuildhelpers"

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

	hash, errHash := gobuildhelpers.GetGitHash(smbExportBuildContext.WorkDir)
	if errHash != nil {
		return errHash
	}
	// fmt.Println(fmt.Sprintf("Git Hash: %s", hash))
	smbExportBuildContext.GitHash = hash

	hight, errHight := gobuildhelpers.GetGitHeight(VERSION_FILE, smbExportBuildContext.WorkDir)
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

	var errFinBuild error
	smbExportBuildContext.PackagesToBuild, errFinBuild = gobuildhelpers.FindPackagesToBuild(filepath.Join(smbExportBuildContext.SourceDir, PROJECT_NAME_SPACE, "cmd"))
	if errFinBuild != nil {
		return errFinBuild
	}

	var errFinTest error
	smbExportBuildContext.PackagesToTest, errFinTest = gobuildhelpers.FindPackagesToTest(filepath.Join(smbExportBuildContext.SourceDir, PROJECT_NAME_SPACE))
	if errFinTest != nil {
		return errFinTest
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

	ldfFlags := fmt.Sprintf("-X main.version=%s", smbExportBuildContext.ProgramVersion)

	err := gobuildhelpers.BuildFolders(smbExportBuildContext.PackagesToBuild, smbExportBuildContext.BinDir, ldfFlags)
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("# ##############################################################################################"))
	return nil
}

// Runs the tests for the project
func Test() error {
	mg.Deps(getEnvironment, Clean, GetBuildName, installTestDeps)
	fmt.Println(fmt.Sprintf("Testing samba-exporter... "))
	fmt.Println(fmt.Sprintf("# ##############################################################################################"))

	xmlResult := filepath.Join(smbExportBuildContext.LogDir, "TestsResult.xml")
	logFileName := "TestRun.log"

	testErrors := gobuildhelpers.RunTestFolders(smbExportBuildContext.PackagesToTest, smbExportBuildContext.LogDir, logFileName)

	errConv := gobuildhelpers.ConvertTestResults(filepath.Join(smbExportBuildContext.LogDir, logFileName), xmlResult, smbExportBuildContext.WorkDir)
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

	err := gobuildhelpers.CoverTestFolders(smbExportBuildContext.PackagesToTest, smbExportBuildContext.LogDir, "TestCoverage.log")
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("# ##############################################################################################"))
	return nil
}

// Remove all build output
func Clean() error {
	mg.Deps(getEnvironment)
	fmt.Println("Cleaning...")
	fmt.Println(fmt.Sprintf("# ##############################################################################################"))

	err := gobuildhelpers.RemovePaths([]string{
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

	err := gobuildhelpers.InstallTestConverter(filepath.Join(smbExportBuildContext.WorkDir, "build"))
	if err != nil {
		return err
	}

	fmt.Println("# ########################################################################################")
	return nil
}

func readVersionMaster() (string, error) {
	content, err := ioutil.ReadFile(smbExportBuildContext.VersionFilePath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(content)), nil
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
