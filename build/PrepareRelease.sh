#!/bin/bash
# ################################################################################################################
# Copyright 2021 by tobi@backfrak.de. All
# rights reserved. Use of this source code is governed
# by a BSD-style license that can be found in the
# LICENSE file.
# ################################################################################################################
# Script to create a release branch. 
# ################################################################################################################

# ################################################################################################################
# function definition
# ################################################################################################################
function print_usage()  {
    echo "Usage: $0 options"
    echo "-help    Print this help"
    echo "This script will do all needed steps to do a new samba_exporter release localy."
    echo "And ask in the end if the result should be published on GutHub"
}

# ################################################################################################################
# variable asigenment
# ################################################################################################################
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
BRANCH_ROOT="$SCRIPT_DIR/.."
LOG_DIR="$BRANCH_ROOT/logs"
RELEASE_BASE_BRANCH="main"
RELEASE_BRANCH_NAME_PREFIX="release/V"
VERSION_MASTER_PATH="VersionMaster.txt"

# ################################################################################################################
# functional code
# ################################################################################################################

if [ "$1" == "-help" ]; then
    print_usage
    exit 0
fi  

pushd "$BRANCH_ROOT" >> /dev/null

actualBranch=$(git status -b -s)
if [ $? != 0 ]; then
    echo "Error: Can not get branch information"
    popd
    exit -1
fi 

statusLines=$(echo "$actualBranch" | wc -l)
if [ "$statusLines" != "1" ]; then
    echo "Error: There are change files that are not checked in."
    popd
    exit -1
fi 

if [[ $actualBranch == *"[ahead"* ]]; then 
    echo "Error: Local repository is ahead of remote"
    popd
    exit -1
fi 

expectedStatus="## $RELEASE_BASE_BRANCH...origin/$RELEASE_BASE_BRANCH"
if [ "$expectedStatus" != "$actualBranch" ]; then
    echo "Error: Not running on $RELEASE_BASE_BRANCH branch"
    popd
    exit -1
fi 

versionInfo=$(cat $VERSION_MASTER_PATH)
releaseBranchName="$RELEASE_BRANCH_NAME_PREFIX$versionInfo"
echo "Release Branch name with name \"$releaseBranchName\" will be created"
git checkout -b "$releaseBranchName"
if [ "$?" != "0" ]; then
    echo "Error while creating release branch \"$releaseBranchName\""
    popd
    exit -1
fi 

echo "Switch back to \"$RELEASE_BASE_BRANCH\" branch"
git checkout "$RELEASE_BASE_BRANCH"
if [ "$?" != "0" ]; then
    echo "Error while switching back to \"$RELEASE_BASE_BRANCH\""
    popd
    exit -1
fi

IFS='.' read -r -a versionNumbersArray <<< "$versionInfo"
actualMinorNumber=${versionNumbersArray[1]}
nextMinorNumber=$(expr $actualMinorNumber + 1)
echo "Actual version is: Major ${versionNumbersArray[0]}, Minor $actualMinorNumber,"
echo "Next version is: Major ${versionNumbersArray[0]}, Minor $nextMinorNumber,"
nextVersionInfo="${versionNumbersArray[0]}.$nextMinorNumber"
echo "Set new version $nextVersionInfo"
echo -n "$nextVersionInfo" > $VERSION_MASTER_PATH

echo "Commit the changed version master"
git commit -a -m "Update version number master to $nextVersionInfo"
if [ "$?" != "0" ]; then
    echo "Error while commit version number master to \"$RELEASE_BASE_BRANCH\""
    popd
    exit -1
fi
git status

echo "You want to push this changes now? (yes|no)"
read answer
if [ "$answer" != "yes" ]; then 
    echo "Push aborted"
    popd
    exit 0
fi 
git push origin "$RELEASE_BASE_BRANCH"
if [ "$?" != "0" ]; then
    echo "Error while push \"$RELEASE_BASE_BRANCH\" to origin"
    popd
    exit -1
fi

git push origin "$releaseBranchName"
if [ "$?" != "0" ]; then
    echo "Error while push \"$releaseBranchName\" to origin"
    popd
    exit -1
fi

popd  >> /dev/null
exit 0
