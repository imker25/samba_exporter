#!/bin/bash
# ######################################################################################
# Copyright 2021 by tobi@backfrak.de. All
# rights reserved. Use of this source code is governed
# by a BSD-style license that can be found in the
# LICENSE file.
# ######################################################################################
# Script will run in a container and do the follwoing steps
# * clone the GitHub git repo
# * switch to the gh-pages banrch
# * update the content of repos/debian with new packages using reprepro
# * move the new conetent of repos/debian out of the container
# ######################################################################################


# ################################################################################################################
# variable asigenment
# ################################################################################################################
HOST_FOLDER="/build_results"
WORK_DIR="/tmp/"

# ################################################################################################################
# parameter and environment check
# ################################################################################################################
if [ "$1" == "" ]; then
    echo "Error: No Tag given"
    exit 1
else 
    tag=$1
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
# ################################################################################################################
# functional code
# ################################################################################################################
echo "# ###################################################################"
echo "Will update repofolder with packgages for version '$tag'"
echo "# ###################################################################"
if [ ! -d "$HOST_FOLDER" ]; then
    echo "Error: The folder mounted from the host '$HOST_FOLDER' was not found"
    exit 1
fi

mkdir -p "$WORK_DIR"
if [ ! -d "$WORK_DIR" ]; then
    echo "Error: The work dir '$WORK_DIR' was not found after creation"
    exit 1
fi

mkdir -p /root/.gpg 
echo "$LAUNCHPAD_GPG_KEY_PUB" > /root/.gpg/imker-bienenkaefig.pub.asc
echo "$LAUNCHPAD_GPG_KEY_PRV" > /root/.gpg/imker-bienenkaefig.asc

gpg --import --batch --no-tty /root/.gpg/imker-bienenkaefig.asc
gpg --edit-key --batch --no-tty  CB6E90E9EC323850B16C1C14A38A1091C018AE68 trust quit
gpg --list-keys --batch --no-tty 

cd "$WORK_DIR"
git clone https://github.com/imker25/samba_exporter.git
if [ "$?" != "0" ]; then 
    echo "Error during git clone"
    exit 1
fi
if [ ! -d "$WORK_DIR/samba_exporter" ]; then
    echo "Error: The repository dir '$WORK_DIR/samba_exporter' was not found after clone"
    exit 1
fi
cd "$WORK_DIR/samba_exporter"
git checkout --track origin/gh-pages
if [ "$?" != "0" ]; then 
    echo "Error during git checkout --track origin/gh-pages"
    exit 1
fi
git pull 
if [ "$?" != "0" ]; then 
    echo "Error during git pull"
    exit 1
fi
echo "# ###################################################################"
echo "git status"
git status

echo "# ###################################################################"
echo "Update the repo with bullseye package"
reprepro --basedir "./repos/debian/" includedeb bullseye "$HOST_FOLDER/deb-packages/binary/samba-exporter_$tag~ppa1~debian11_amd64.deb"
if [ "$?" != "0" ]; then 
    echo "Error during reprepro for bullseye"
    exit 1
fi
echo "# ###################################################################"
echo "Update the repo with buster package"
reprepro --basedir "./repos/debian/" includedeb buster "$HOST_FOLDER/deb-packages/binary/samba-exporter_$tag~ppa1~debian10_amd64.deb"
if [ "$?" != "0" ]; then 
    echo "Error during reprepro for buster"
    exit 1
fi

echo "# ###################################################################"
echo "git status"
git status

echo "# ###################################################################"
echo "Copy the update debian repo out of the container"
mkdir -p "$HOST_FOLDER/pages/repos"
mv "./repos/debian" "$HOST_FOLDER/pages/repos"
chmod -R 777 "$HOST_FOLDER/pages"

exit 0
