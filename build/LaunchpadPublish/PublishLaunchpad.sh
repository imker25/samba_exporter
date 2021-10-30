#!/bin/bash
# ######################################################################################
# Copyright 2021 by tobi@backfrak.de. All
# rights reserved. Use of this source code is governed
# by a BSD-style license that can be found in the
# LICENSE file.
# ######################################################################################
# Script to do the following inside a container
# * clone the launchpad git repo
# * import the given samba_exporter github sources to the launchpad git repo
# * do the needed conversation steps, so debian package build can run
# * run debian binary package  build
# * run debian source package  build with tagging
# * commit the changes to the launchpad git repo
# * upload the debian source package to the launchpad ppa
# * push the launchpad git repo with tags
# ######################################################################################

# ################################################################################################################
# function definition
# ################################################################################################################
function print_usage()  {
    echo "Script to transfer a github tag to launchpad and publish the package in a ppa"
    echo ""
    echo "Usage: $0 tag <dry>"
    echo "-help     Print this help"
    echo "tag       The tag on the github repo to import, e. g. 0.7.5"
    echo "dry       Optional: Do not push the changes to launchpad git and not upload the sources to ppa"
    echo ""
    echo "The script expect the following environment variables to be set"
    echo "  LAUNCHPAD_SSH_ID_PUB        Public SSH key for the launchapd git repo"
    echo "  LAUNCHPAD_SSH_ID_PRV        Private SSH key for the launchapd git repo"
    echo "  LAUNCHPAD_GPG_KEY_PUB       Public GPG Key for the launchpad ppa"
    echo "  LAUNCHPAD_GPG_KEY_PRV       Private GPG Key for the launchpad ppa"
}

# ################################################################################################################
# variable asigenment
# ################################################################################################################
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
BRANCH_ROOT="$SCRIPT_DIR/../.."
LAUNCHPAD_GIT_REPO="ssh://imker@git.launchpad.net/samba-exporter"
GITHUB_RELEASE_URL="https://github.com/imker25/samba_exporter/archive/refs/tags"
BUILD_DIR="/build_area"
WORK_DIR="$BUILD_DIR/samba-exporter"
GITHUB_PROMETHEUS_VERSION="v1.11.0"
LAUNCHPAD_PROMETHEUS_VERSION="v1.9.0"


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
    echo "It's a dry run! No changes will be uploaded or pushed to launchpad"
else
    dryRun="false"
fi

if [ "$LAUNCHPAD_SSH_ID_PUB" == "" ]; then
    echo "Error: Environment variables LAUNCHPAD_SSH_ID_PUB not set"
    print_usage
    exit 1
fi

if [ "$LAUNCHPAD_SSH_ID_PRV" == "" ]; then
    echo "Error: Environment variables LAUNCHPAD_SSH_ID_PRV not set"
    print_usage
    exit 1
fi

if [ "$LAUNCHPAD_GPG_KEY_PUB" == "" ]; then
    echo "Error: Environment variables LAUNCHPAD_GPG_KEY_PUB not set"
    print_usage
    exit 1
fi

if [ "$LAUNCHPAD_GPG_KEY_PRV" == "" ]; then
    echo "Error: Environment variables LAUNCHPAD_GPG_KEY_PRV not set"
    print_usage
    exit 1
fi


if [[ "$tag" =~ "-pre" ]]; then
    if [ "$dryRun" == "false" ]; then
        echo "Warinig: A pre release will be imported to launchpad!"
    else
        echo "Do a dry run with a pre release"
    fi
    gitTag="${tag/-pre/}"
    preRelease="true"
else 
    gitTag=$tag
    preRelease="false"
fi

ubuntuVersion=$(lsb_release -rs)
packageVersion="${tag}~ppa1~ubuntu${ubuntuVersion}"

# ################################################################################################################
# functional code
# ################################################################################################################

echo "Publish github release $tag to launchpad as version $packageVersion"
echo "# ###################################################################"

echo "Prepare for operation"
mkdir -p /root/.ssh
echo "$LAUNCHPAD_SSH_ID_PUB" > /root/.ssh/id_rsa.pub
chmod 600 /root/.ssh/id_rsa.pub
echo "$LAUNCHPAD_SSH_ID_PRV" > /root/.ssh/id_rsa
chmod 600 /root/.ssh/id_rsa
mkdir -p /root/.gpg 
echo "$LAUNCHPAD_GPG_KEY_PUB" > /root/.gpg/imker-bienenkaefig.pub.asc
echo "$LAUNCHPAD_GPG_KEY_PRV" > /root/.gpg/imker-bienenkaefig.asc

gpg --import --batch --no-tty /root/.gpg/imker-bienenkaefig.asc
gpg --edit-key --batch --no-tty  CB6E90E9EC323850B16C1C14A38A1091C018AE68 trust quit
gpg --list-keys --batch --no-tty 

mkdir -p "$BUILD_DIR"
cd "$BUILD_DIR"
ssh-keyscan -t rsa git.launchpad.net >> ~/.ssh/known_hosts
git clone "$LAUNCHPAD_GIT_REPO"
if [ "$?" != "0" ]; then 
    echo "Error: Can not clone the launchpad repo $LAUNCHPAD_GIT_REPO"
    exit 1
fi

cd "$WORK_DIR"
git pull --all
git checkout --track origin/upstream
git checkout master
git branch
git tag | grep "${tag}"
if [ "$?" == "1" ]; then
    echo "Import the source package from github"
    gbp import-orig --merge-mode=replace --upstream-version=$tag $GITHUB_RELEASE_URL/$tag.tar.gz
    if [ "$?" != "0" ]; then 
        echo "Error: Can not import tag $tag from $GITHUB_RELEASE_URL"
        exit 1
    fi
    git checkout master
    if [ "$preRelease" == "true" ]; then
        echo "Tag with $gitTag"
        git tag upstream/$gitTag
    fi
else
    echo "Use already imported sources"
fi

echo "Create the branch 'ubuntu-${ubuntuVersion}/v${tag}' to work on"
git checkout -b "ubuntu-${ubuntuVersion}/v${tag}"
git status

echo "# ###################################################################"
echo "# Patch the files"
given_version=$(cat "$WORK_DIR/VersionMaster.txt")
echo "$packageVersion" > "$WORK_DIR/VersionMaster.txt"
echo "Version Prefix: $given_version"

sed -i "s/samba-exporter ($given_version)/samba-exporter ($packageVersion)/g" $WORK_DIR/changelog
# cat $WORK_DIR/changelog
rm -rf $WORK_DIR/debian/*
cp -rv -L $WORK_DIR/install/debian/* $WORK_DIR/debian

echo "Patch package dependencies acording the ubuntu version"
if [ "$ubuntuVersion" == "20.04"]; then
    find . -name "go.mod" -exec sed -i "s/require github.com\\/prometheus\\/client_golang $GITHUB_PROMETHEUS_VERSION/require github.com\\/prometheus\\/client_golang $LAUNCHPAD_PROMETHEUS_VERSION/g" {} \;
else 
    echo "Not running on ubuntu 20.04"
fi 

echo "# ###################################################################"
echo "# Build packages before git commit"
gbp buildpackage -kimker@bienenkaefig.de --git-ignore-new 
if [ "$?" != "0" ]; then 
    echo "Error: Can not build the packages with default paramters"
    exit 1
fi
rm -rfv ../samba-exporter_$packageVersion*

mkdir -p $WORK_DIR/debian/source
echo "3.0 (native)" > $WORK_DIR/debian/source/format
# echo "src/man/samba_exporter.1.gz" > debian/source/include-binaries
# echo "src/man/samba_statusd.1.gz" >> debian/source/include-binaries
# echo "src/man/start_samba_statusd.1.gz" >> debian/source/include-binaries
git add debian/source/*

gbp buildpackage -kimker@bienenkaefig.de --git-builder="debuild -i -I -S " --git-ignore-new
if [ "$?" != "0" ]; then 
    echo "Error: Can not build the source package"
    exit 1
fi
rm -rfv ../samba-exporter_$packageVersion*

echo "# ###################################################################"
echo "# git commit"
git add .
git status
git commit -a -m "Deploy patches after GitHub V$tag import"
echo "# ###################################################################"
echo "git status"
git status

echo "# ###################################################################"
echo "# Build source package for upload"

gbp buildpackage -kimker@bienenkaefig.de --git-builder="debuild -i -I -S" --git-tag --git-debian-branch="ubuntu-${ubuntuVersion}/v${tag}"
if [ "$?" != "0" ]; then 
    echo "Error: Can not build the source package for upload"
    exit 1
fi

echo "# ###################################################################"
if [ "$dryRun" == "false" ]; then
    echo "Upload source package"
    echo "dput ppa:imker/samba-exporter-ppa ../samba-exporter_${packageVersion}_source.changes "
    dput ppa:imker/samba-exporter-ppa ../samba-exporter_${packageVersion}_source.changes 
    if [ "$?" != "0" ]; then 
        echo "Error: Can not upload the source package to the launchpad ppa"
        exit 1
    fi
else
    echo "Upload skiped due to dry run"
fi

echo "# ###################################################################"
echo "# Push git to launchpad"
if [ "$dryRun" == "false" ]; then
    echo "git push --all origin"
    git push --all origin
    if [ "$?" != "0" ]; then 
        echo "Error: Can not push changes to lauchpad git"
        exit 1
    fi
    echo "git push --tag"
    git push --tag
    if [ "$?" != "0" ]; then 
        echo "Error: Can not push tags to launchpad git"
        exit 1
    fi
else 
    echo "Push skiped due to dry run"
fi

exit 0

