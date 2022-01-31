#!/bin/bash
# ######################################################################################
# Copyright 2022 by tobi@backfrak.de. All
# rights reserved. Use of this source code is governed
# by a BSD-style license that can be found in the
# LICENSE file.
# ######################################################################################
# Script to build and run a docker container and do the following inside the container
# * import the given samba_exporter github sources to the rpm build area
# * do the needed conversation steps, so rpm package build can run
# * run rpm binary package  build
# * run rpm source package  build 
# ######################################################################################

# ################################################################################################################
# function definition
# ################################################################################################################
function print_usage()  {
    echo "Script to create the release RPM by building a suitable container and run rpm-publish/RpmPublish.sh"
    echo ""
    echo "Usage: $0 tag <dry>"
    echo "-help     Print this help"
    echo "tag       The tag on the github repo to import, e. g. 0.7.5"
    echo "dry       Optional: Do not publish the RPM"
    echo ""
    echo "The script expect the following environment variables to be set"
    echo "  COPR_SSH_ID_PUB        Public SSH key for the launchapd git repo"
    echo "  COPR_SSH_ID_PRV        Private SSH key for the launchapd git repo"
    echo "  COPR_GPG_KEY_PUB       Public GPG Key for the copr ppa"
    echo "  COPR_GPG_KEY_PRV       Private GPG Key for the copr ppa"
}

function buildAndRunDocker() {
    distVersion="$1"

    echo "Build the needed container from '$WORK_DIR/Dockerfile.${distVersion}', logging to $BRANCH_ROOT/logs/docker-build-${distVersion}.log"
    docker build --file "$WORK_DIR/Dockerfile.${distVersion}" --tag rpm-publish-container-$distVersion . > $BRANCH_ROOT/logs/docker-build-${distVersion}.log 2>&1
    if [ "$?" != "0" ]; then 
        echo "Error during docker build"
        return 1
    fi
    echo "# ###################################################################"
    echo "Run the container"

    if [ "$dryRun" == "false" ]; then
        docker run --env COPR_SSH_ID_PUB="$COPR_SSH_ID_PUB" \
            --env COPR_SSH_ID_PRV="$COPR_SSH_ID_PRV"  \
            --env COPR_GPG_KEY_PUB="$COPR_GPG_KEY_PUB" \
            --env COPR_GPG_KEY_PRV="$COPR_GPG_KEY_PRV" \
            --mount type=bind,source="$RPM_PACKAGE_DIR",target="/build_results" \
            -i rpm-publish-container-$distVersion \
            /bin/bash -c "/RpmPublish.sh $tag"
    else
        docker run --env COPR_SSH_ID_PUB="$COPR_SSH_ID_PUB" \
            --env COPR_SSH_ID_PRV="$COPR_SSH_ID_PRV"  \
            --env COPR_GPG_KEY_PUB="$COPR_GPG_KEY_PUB" \
            --env COPR_GPG_KEY_PRV="$COPR_GPG_KEY_PRV" \
            --mount type=bind,source="$RPM_PACKAGE_DIR",target="/build_results" \
            -i rpm-publish-container-$distVersion \
            /bin/bash -c "/RpmPublish.sh $tag dry"
    fi

    if [ "$?" != "0" ]; then 
        echo "Error during docker run"
        return 1
    fi
    return 0
}

# ################################################################################################################
# variable asigenment
# ################################################################################################################
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
BRANCH_ROOT="$SCRIPT_DIR/.."
WORK_DIR="$SCRIPT_DIR/rpm-publish"
RPM_PACKAGE_DIR="$BRANCH_ROOT/bin/rpm-packages"

# ################################################################################################################
# parameter and environment check
# ################################################################################################################

if [ "$1" == "-help" ]; then
    print_usage
    exit 0
fi  

if [ "$1" == "" ]; then
    echo "Error: No Tag given"
    print_usage
    exit 1
else 
    tag=$1
fi

if [ "$2" == "dry" ]; then
    dryRun="true"
    echo "It's a dry run! No changes will be uploaded or pushed to copr"
else
    dryRun="false"
fi

if [ "$COPR_SSH_ID_PUB" == "" ]; then
    echo "Error: Environment variables COPR_SSH_ID_PUB not set"
    print_usage
    exit 1
fi

if [ "$COPR_SSH_ID_PRV" == "" ]; then
    echo "Error: Environment variables COPR_SSH_ID_PRV not set"
    print_usage
    exit 1
fi


if [ "$COPR_GPG_KEY_PUB" == "" ]; then
    echo "Error: Environment variables COPR_GPG_KEY_PUB not set"
    print_usage
    exit 1
fi

if [ "$COPR_GPG_KEY_PRV" == "" ]; then
    echo "Error: Environment variables COPR_GPG_KEY_PRV not set"
    print_usage
    exit 1
fi


if [[ "$tag" =~ "-pre" ]]; then
    if [ "$dryRun" == "false" ]; then
        echo "Warinig: A pre release will be published!"
    else
        echo "Do a dry run with a pre release"
    fi
fi
# ################################################################################################################
# functional code
# ################################################################################################################

pushd "$BRANCH_ROOT"
echo "Get log messages for changelog update"
lastVersionChangeCommit=$(git log --pretty=format:"%H" -n 1 --follow "$BRANCH_ROOT/VersionMaster.txt")
echo "Get log messages since commit \"${lastVersionChangeCommit}\" "
mkdir -pv $BRANCH_ROOT/tmp/
git log "$lastVersionChangeCommit".. --pretty=format:"--::%an;;;;%ae;;;;%B" > "$BRANCH_ROOT/tmp/commit_logs"

popd

pushd "$WORK_DIR"
if [ -d "$RPM_PACKAGE_DIR" ]; then
    echo "Use existing $RPM_PACKAGE_DIR dir after cleanup"
    rm -rf $RPM_PACKAGE_DIR/*
else 
    echo "Create $RPM_PACKAGE_DIR dir"
    mkdir -p "$RPM_PACKAGE_DIR"
fi

if [ -d "$BRANCH_ROOT/logs" ]; then
    echo "Use existing $BRANCH_ROOT/logs dir"
    if ls $BRANCH_ROOT/logs/docker-build*.log 1> /dev/null 2>&1; then
        echo "Delete existing $BRANCH_ROOT/logs/docker-build*.log"
        rm -rf $BRANCH_ROOT/logs/docker-build*.log 
    fi

else 
    echo "Create $BRANCH_ROOT/logs dir"
    mkdir -p "$BRANCH_ROOT/logs"
fi


cp -v "$BRANCH_ROOT/tmp/commit_logs" "$RPM_PACKAGE_DIR"

dockerError="false"
echo "Publish tag $tag on corp within a docker cotainer for fedora 35"
echo "# ###################################################################"
buildAndRunDocker "fedora35"
if [ "$?" != "0" ]; then
    dockerError="true"
fi


popd

echo "# ###################################################################"
echo "Delete the container image when done"    
docker rmi -f $(docker images --filter=reference="rpm-publish*" -q) 

if [ "$dockerError" == "true" ];then 
    echo "Error detected"
    exit 1
fi

exit 0


