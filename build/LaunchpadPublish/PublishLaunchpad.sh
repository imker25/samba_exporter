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
BUILD_DIR="$HOME/build_area"
WORK_DIR="$BUILD_DIR/samba-exporter"
GITHUB_PROMETHEUS_VERSION="v1.11.0"
LAUNCHPAD_PROMETHEUS_VERSION="v1.9.0"
export DEBEMAIL="imker@bienekaefig.de"
export DEBFULLNAME="Tobias Zellner"


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

distribution=$(lsb_release -is)
distVersionNumber=$(lsb_release -rs)
distCodeName=$(lsb_release -cs)
packageVersion="${tag}~ppa1~${distribution,,}${distVersionNumber}"

# ################################################################################################################
# functional code
# ################################################################################################################

echo "Setup git"
echo "# ###################################################################"
git config --global user.name "Tobias Zellner"
git config --global user.email imker@bienekaefig.de

echo "Publish github release $tag to launchpad as version $packageVersion"
echo "# ###################################################################"


echo "Prepare for operation"
mkdir -p $HOME/.ssh
chmod 700 $HOME/.ssh
echo "$LAUNCHPAD_SSH_ID_PUB" > $HOME/.ssh/id_rsa.pub
chmod 600 $HOME/.ssh/id_rsa.pub
echo "$LAUNCHPAD_SSH_ID_PRV" > $HOME/.ssh/id_rsa
chmod 600 $HOME/.ssh/id_rsa
mkdir -p $HOME/.gpg 
echo "$LAUNCHPAD_GPG_KEY_PUB" > $HOME/.gpg/imker-bienenkaefig.pub.asc
echo "$LAUNCHPAD_GPG_KEY_PRV" > $HOME/.gpg/imker-bienenkaefig.asc

gpg --import --batch --no-tty $HOME/.gpg/imker-bienenkaefig.asc
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

    echo "Add new changelog entry"
    echo "# ###################################################################"
    dch --distribution "focal" -v $packageVersion "Work on Version $packageVersion" < /bin/true

    echo "Add git log to the changelog"
    echo "# ###################################################################"
    changes=$(cat /build_results/commit_logs)
    delimiter="--::"
    string=$changes$delimiter
    #Split the text changes on the delimiter
    changeEntries=()
    while [[ $string ]]; do
    changeEntries+=( "${string%%"$delimiter"*}" )
    string=${string#*"$delimiter"}
    done

    delimiter=";;;;"
    for entry in "${changeEntries[@]}"
    do
        
        string=$entry$delimiter
        entryFileds=()
        while [[ $string ]]; do
            entryFileds+=( "${string%%"$delimiter"*}" )
            string=${string#*"$delimiter"}
        done

        echo "Author: ${entryFileds[0]}"
        echo "Mail: ${entryFileds[1]}"
        echo "Message: ${entryFileds[2]}"
        if [ "${entryFileds[2]}" != "" ]; then
            dch -a "${entryFileds[2]} (by ${entryFileds[1]})" < /bin/true
        fi
    done
    echo "# ###################################################################"
    echo "git status"
    git status
    echo "# ###################################################################"
    echo "git commit"
    git commit -a -m "Add changelog for v${tag}"
    echo "# ###################################################################"
    if [ "$preRelease" == "true" ]; then
        echo "Tag with $gitTag"
        git tag upstream/$gitTag
    fi
else
    echo "Use already imported sources"
fi

echo "Create the branch '${distribution,,}-${distVersionNumber}/v${tag}' to work on"
git checkout -b "${distribution,,}-${distVersionNumber}/v${tag}"
git status


echo "# ###################################################################"
echo "# Patch the files"
cp -v "$WORK_DIR/debian/changelog" "$WORK_DIR/install/debian/changelog"
echo "$packageVersion" > "$WORK_DIR/VersionMaster.txt"

echo "Patch package dependencies acording the distribution and version"
if [ "$distVersionNumber" == "20.04" ] && [ "$distribution" == "Ubuntu" ]; then
    find . -name "go.mod" -exec sed -i "s/require github.com\\/prometheus\\/client_golang $GITHUB_PROMETHEUS_VERSION/require github.com\\/prometheus\\/client_golang $LAUNCHPAD_PROMETHEUS_VERSION/g" {} \;
    find . -name "*.go" -exec sed -i "s/github.com\\/shirou\\/gopsutil\\/v3\\/process/github.com\\/shirou\\/gopsutil\\/process/g" {} \;
    find . -name "*.go" -exec sed -i "s/|log.Lmsgprefix/ /g" {} \;
    mv -v "$WORK_DIR/src/tobi.backfrak.de/cmd/samba_statusd/go.mod" "$WORK_DIR/src/tobi.backfrak.de/cmd/samba_statusd/go.mod.gopsutil-v3"
    mv -v "$WORK_DIR/src/tobi.backfrak.de/cmd/samba_statusd/go.sum" "$WORK_DIR/src/tobi.backfrak.de/cmd/samba_statusd/go.sum.gopsutil-v3"
    mv -v "$WORK_DIR/src/tobi.backfrak.de/cmd/samba_statusd/go.mod.gopsutil-v2" "$WORK_DIR/src/tobi.backfrak.de/cmd/samba_statusd/go.mod" 
    mv -v "$WORK_DIR/src/tobi.backfrak.de/cmd/samba_statusd/go.sum.gopsutil-v2" "$WORK_DIR/src/tobi.backfrak.de/cmd/samba_statusd/go.sum"
    mv -v "$WORK_DIR/src/tobi.backfrak.de/internal/smbstatusdbl/go.mod" "$WORK_DIR/src/tobi.backfrak.de/internal/smbstatusdbl/go.mod.gopsutil-v3"
    mv -v "$WORK_DIR/src/tobi.backfrak.de/internal/smbstatusdbl/go.sum" "$WORK_DIR/src/tobi.backfrak.de/internal/smbstatusdbl/go.sum.gopsutil-v3"
    mv -v "$WORK_DIR/src/tobi.backfrak.de/internal/smbstatusdbl/go.mod.gopsutil-v2" "$WORK_DIR/src/tobi.backfrak.de/internal/smbstatusdbl/go.mod" 
    mv -v "$WORK_DIR/src/tobi.backfrak.de/internal/smbstatusdbl/go.sum.gopsutil-v2" "$WORK_DIR/src/tobi.backfrak.de/internal/smbstatusdbl/go.sum" 
else 
    echo "Not running on ubuntu 20.04 (focal)"
fi 


if [ "$distVersionNumber" == "22.04" ] && [ "$distribution" == "Ubuntu" ]; then
    sed -i "s/focal;/jammy;/g" $WORK_DIR/install/debian/changelog
    sed -i "s/ubuntu20.04/ubuntu22.04/g" $WORK_DIR/install/debian/changelog
    sed -i "s/golang-1.16,/golang-1.18,/g" $WORK_DIR/install/debian/control    
else 
    echo "Not running on ubuntu 22.04 (jammy)"
fi

if [ "$distVersionNumber" == "22.10" ] && [ "$distribution" == "Ubuntu" ]; then
    sed -i "s/focal;/kinetic;/g" $WORK_DIR/install/debian/changelog
    sed -i "s/ubuntu20.04/ubuntu22.10/g" $WORK_DIR/install/debian/changelog
    sed -i "s/golang-1.16,/golang-1.19,/g" $WORK_DIR/install/debian/control    
else 
    echo "Not running on ubuntu 22.10 (kinetic)"
fi

if [ "$distVersionNumber" == "23.04" ] && [ "$distribution" == "Ubuntu" ]; then
    sed -i "s/focal;/lunar;/g" $WORK_DIR/install/debian/changelog
    sed -i "s/ubuntu20.04/ubuntu23.04/g" $WORK_DIR/install/debian/changelog
    sed -i "s/golang-1.16,/golang-1.20,/g" $WORK_DIR/install/debian/control 
    find . -name "*.go" -exec sed -i "s/github.com\\/shirou\\/gopsutil\\/v3\\/process/github.com\\/shirou\\/gopsutil\\/process/g" {} \;
else 
    echo "Not running on lunar 23.04 (lunar)"
fi

if [ "$distVersionNumber" == "23.10" ] && [ "$distribution" == "Ubuntu" ]; then
    sed -i "s/focal;/mantic;/g" $WORK_DIR/install/debian/changelog
    sed -i "s/ubuntu20.04/ubuntu23.10/g" $WORK_DIR/install/debian/changelog
    sed -i "s/golang-1.16,/golang-1.21,/g" $WORK_DIR/install/debian/control 
    find . -name "*.go" -exec sed -i "s/github.com\\/shirou\\/gopsutil\\/v3\\/process/github.com\\/shirou\\/gopsutil\\/process/g" {} \;
else 
    echo "Not running on lunar 23.10 (mantic)"
fi

if [ "$distVersionNumber" == "24.04" ] && [ "$distribution" == "Ubuntu" ]; then
    sed -i "s/focal;/noble;/g" $WORK_DIR/install/debian/changelog
    sed -i "s/ubuntu20.04/ubuntu24.04/g" $WORK_DIR/install/debian/changelog
    sed -i "s/golang-1.16,/golang-1.22,/g" $WORK_DIR/install/debian/control 
    find . -name "*.go" -exec sed -i "s/github.com\\/shirou\\/gopsutil\\/v3\\/process/github.com\\/shirou\\/gopsutil\\/process/g" {} \;
else 
    echo "Not running on lunar 23.10 (mantic)"
fi

if [ "$distVersionNumber" == "12" ] && [ "$distribution" == "Debian" ]; then
    sed -i "s/focal;/bookworm;/g" $WORK_DIR/install/debian/changelog
    sed -i "s/ubuntu20.04/debian12/g" $WORK_DIR/install/debian/changelog
    sed -i "s/golang-1.16,/golang-1.19,/g" $WORK_DIR/install/debian/control 
    find . -name "*.go" -exec sed -i "s/github.com\\/shirou\\/gopsutil\\/v3\\/process/github.com\\/shirou\\/gopsutil\\/process/g" {} \;
else 
    echo "Not running on Debian 12 (Bookworm)"
fi

if [ "$distVersionNumber" == "11" ] && [ "$distribution" == "Debian" ]; then
    sed -i "s/focal;/bullseye;/g" $WORK_DIR/install/debian/changelog
    sed -i "s/ubuntu20.04/debian11/g" $WORK_DIR/install/debian/changelog
    sed -i "s/golang-1.16,/golang-1.15,/g" $WORK_DIR/install/debian/control
    find . -name "*.go" -exec sed -i "s/github.com\\/shirou\\/gopsutil\\/v3\\/process/github.com\\/shirou\\/gopsutil\\/process/g" {} \;
    mv -v "$WORK_DIR/src/tobi.backfrak.de/cmd/samba_statusd/go.mod" "$WORK_DIR/src/tobi.backfrak.de/cmd/samba_statusd/go.mod.gopsutil-v3"
    mv -v "$WORK_DIR/src/tobi.backfrak.de/cmd/samba_statusd/go.sum" "$WORK_DIR/src/tobi.backfrak.de/cmd/samba_statusd/go.sum.gopsutil-v3"
    mv -v "$WORK_DIR/src/tobi.backfrak.de/cmd/samba_statusd/go.mod.gopsutil-v2" "$WORK_DIR/src/tobi.backfrak.de/cmd/samba_statusd/go.mod" 
    mv -v "$WORK_DIR/src/tobi.backfrak.de/cmd/samba_statusd/go.sum.gopsutil-v2" "$WORK_DIR/src/tobi.backfrak.de/cmd/samba_statusd/go.sum"
    mv -v "$WORK_DIR/src/tobi.backfrak.de/internal/smbstatusdbl/go.mod" "$WORK_DIR/src/tobi.backfrak.de/internal/smbstatusdbl/go.mod.gopsutil-v3"
    mv -v "$WORK_DIR/src/tobi.backfrak.de/internal/smbstatusdbl/go.sum" "$WORK_DIR/src/tobi.backfrak.de/internal/smbstatusdbl/go.sum.gopsutil-v3"
    mv -v "$WORK_DIR/src/tobi.backfrak.de/internal/smbstatusdbl/go.mod.gopsutil-v2" "$WORK_DIR/src/tobi.backfrak.de/internal/smbstatusdbl/go.mod" 
    mv -v "$WORK_DIR/src/tobi.backfrak.de/internal/smbstatusdbl/go.sum.gopsutil-v2" "$WORK_DIR/src/tobi.backfrak.de/internal/smbstatusdbl/go.sum"     
else 
    echo "Not running on debian 11 (bullseye)"
fi

if [ "$distVersionNumber" == "10" ] && [ "$distribution" == "Debian" ]; then
    sed -i "s/focal;/buster;/g" $WORK_DIR/install/debian/changelog
    sed -i "s/ubuntu20.04/debian10/g" $WORK_DIR/install/debian/changelog
    sed -i "s/golang-1.16,/golang-1.15,/g" $WORK_DIR/install/debian/control
    sed -i "s/gzip (>=1.10)/gzip (>=1.9)/g" $WORK_DIR/install/debian/control
    sed -i "s/golang-any/man-db/g" $WORK_DIR/install/debian/control
    find . -name "*.go" -exec sed -i "s/github.com\\/shirou\\/gopsutil\\/v3\\/process/github.com\\/shirou\\/gopsutil\\/process/g" {} \;
    mv -v "$WORK_DIR/src/tobi.backfrak.de/cmd/samba_statusd/go.mod" "$WORK_DIR/src/tobi.backfrak.de/cmd/samba_statusd/go.mod.gopsutil-v3"
    mv -v "$WORK_DIR/src/tobi.backfrak.de/cmd/samba_statusd/go.sum" "$WORK_DIR/src/tobi.backfrak.de/cmd/samba_statusd/go.sum.gopsutil-v3"
    mv -v "$WORK_DIR/src/tobi.backfrak.de/cmd/samba_statusd/go.mod.gopsutil-v2" "$WORK_DIR/src/tobi.backfrak.de/cmd/samba_statusd/go.mod" 
    mv -v "$WORK_DIR/src/tobi.backfrak.de/cmd/samba_statusd/go.sum.gopsutil-v2" "$WORK_DIR/src/tobi.backfrak.de/cmd/samba_statusd/go.sum"
    mv -v "$WORK_DIR/src/tobi.backfrak.de/internal/smbstatusdbl/go.mod" "$WORK_DIR/src/tobi.backfrak.de/internal/smbstatusdbl/go.mod.gopsutil-v3"
    mv -v "$WORK_DIR/src/tobi.backfrak.de/internal/smbstatusdbl/go.sum" "$WORK_DIR/src/tobi.backfrak.de/internal/smbstatusdbl/go.sum.gopsutil-v3"
    mv -v "$WORK_DIR/src/tobi.backfrak.de/internal/smbstatusdbl/go.mod.gopsutil-v2" "$WORK_DIR/src/tobi.backfrak.de/internal/smbstatusdbl/go.mod" 
    mv -v "$WORK_DIR/src/tobi.backfrak.de/internal/smbstatusdbl/go.sum.gopsutil-v2" "$WORK_DIR/src/tobi.backfrak.de/internal/smbstatusdbl/go.sum"     
else 
    echo "Not running on debian 10 (buster)"
fi

rm -rf $WORK_DIR/debian/*
cp -rv -L $WORK_DIR/install/debian/* $WORK_DIR/debian

echo "# ###################################################################"
echo "Changelog content after mofifications"
cat $WORK_DIR/debian/changelog

# debug exit 
# exit 1

echo "# ###################################################################"
echo "# Build packages before git commit"
gbp buildpackage -kimker@bienenkaefig.de --git-ignore-new 
if [ "$?" != "0" ]; then 
    ls -l $BUILD_DIR/samba-exporter/bin/src/src
    echo "Error: Can not build the packages with default paramters"
    exit 1
fi

echo "# ###################################################################"
if [ -d "/build_results" ]; then 
    mkdir -p /build_results/binary
    echo "Copy binary debian package to the docker host using mount dir '/build_results/binary'"
    cp -v ../*.deb "/build_results/binary/"
    mkdir -p /build_results/man 
    echo "Copy the html man pages to /build_results/man"
    cp -v ./src/man/samba_exporter.1.html /build_results/man/
    cp -v ./src/man/samba_statusd.1.html /build_results/man/
    cp -v ./src/man/start_samba_statusd.1.html /build_results/man/
    chmod -R 777 /build_results/*
else 
    echo "Waring: /build_results does not exist, no debian packages copied to the docker host"
fi

if [ "$distribution" == "Debian" ]; then
    echo "# ###################################################################"
    echo "Running on debian, no source packages build"
    echo "# ###################################################################"
    echo "# git commit"
    git add .
    git status
    git commit -a -m "Deploy patches after GitHub V$tag import for $distribution $distVersionNumber"
    echo "# ###################################################################"
    echo "git status"
    git status
else 
    echo "Delete biniary packages before source package build"
    rm -rfv ../samba-exporter_$packageVersion*

    echo "# ###################################################################"
    echo "Prepeare for source package build"
    mkdir -p $WORK_DIR/debian/source
    echo "3.0 (native)" > $WORK_DIR/debian/source/format
    git add debian/source/*

    echo "# ###################################################################"
    echo "Source package test build"
    gbp buildpackage -kimker@bienenkaefig.de --git-builder="debuild -i -I -S " --git-ignore-new
    if [ "$?" != "0" ]; then 
        echo "Error: Can not build the source package"
        exit 1
    fi
    echo "# ###################################################################"
    echo "Delete souce package test build build"
    rm -rfv ../samba-exporter_$packageVersion*

    echo "# ###################################################################"
    echo "# git commit"
    git add .
    git status
    git commit -a -m "Deploy patches after GitHub V$tag import for $distribution $distVersionNumber"
    echo "# ###################################################################"
    echo "git status"
    git status

    echo "# ###################################################################"
    echo "# Build source package for upload"

    gbp buildpackage -kimker@bienenkaefig.de --git-builder="debuild -i -I -S" --git-tag --git-debian-branch="${distribution,,}-${distVersionNumber}/v${tag}"
    if [ "$?" != "0" ]; then 
        echo "Error: Can not build the source package for upload"
        exit 1
    fi

    if [ ! -f "$BUILD_DIR/samba-exporter_${packageVersion}_source.changes" ]; then
        echo "Can to find the source package at '$BUILD_DIR/samba-exporter_${packageVersion}_source.changes' as expected"
        echo "ls -l $BUILD_DIR/"
        ls -l $BUILD_DIR/
        exit 1
    fi
    echo "# ###################################################################"
    if [ "$dryRun" == "false" ]; then
        echo "Upload source package"
        echo "dput ppa:imker/samba-exporter-ppa \"$BUILD_DIR/samba-exporter_${packageVersion}_source.changes\" "
        dput ppa:imker/samba-exporter-ppa "$BUILD_DIR/samba-exporter_${packageVersion}_source.changes" 
        if [ "$?" != "0" ]; then 
            echo "Error: Can not upload the source package to the launchpad ppa"
            exit 1
        fi
    else
        echo "Upload skiped due to dry run"
    fi
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

