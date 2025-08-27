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

if [[ "$tag" =~ "-pre" ]]; then
    rpmVersion=${tag:0:-4}
else 
    rpmVersion="$tag"
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


# Only needed once, for every new debian version that is added to this repository
# echo "# ###################################################################"
# echo "Patch ./repos/debian/conf/distributions"
# # rm -v ./repos/debian/conf/distributions
# echo "" >> ./repos/debian/conf/distributions
# echo "Origin: imker25.github.io/samba_exporter/" >> ./repos/debian/conf/distributions
# echo "Label: imker25.github.io/samba_exporter/"  >> ./repos/debian/conf/distributions
# echo "Codename: trixie"  >> ./repos/debian/conf/distributions
# echo "Architectures: amd64"  >> ./repos/debian/conf/distributions
# echo "Components: main"  >> ./repos/debian/conf/distributions
# echo "Description: Personal repository for samba-exporter packages"  >> ./repos/debian/conf/distributions
# echo "SignWith: CB6E90E9EC323850B16C1C14A38A1091C018AE68"  >> ./repos/debian/conf/distributions
# echo "" >> ./repos/debian/conf/distributions
# reprepro clearvanished ./repos/debian/

echo "# ###################################################################"
echo "Update the debian repo with trixie package"
reprepro --basedir "./repos/debian/" includedeb trixie "$HOST_FOLDER/deb-packages/binary/samba-exporter_$tag~ppa1~debian13_amd64.deb"
if [ "$?" != "0" ]; then 
    echo "Error during reprepro for trixie"
    exit 1
fi
echo "# ###################################################################"
echo "Update the debian repo with bookworm package"
reprepro --basedir "./repos/debian/" includedeb bookworm "$HOST_FOLDER/deb-packages/binary/samba-exporter_$tag~ppa1~debian12_amd64.deb"
if [ "$?" != "0" ]; then 
    echo "Error during reprepro for bookworm"
    exit 1
fi

echo "# ###################################################################"
echo "Update the rpm repo with fc28 package"
if [ ! -f "$HOST_FOLDER/rpm-packages/Fedora-28/samba-exporter-${rpmVersion}-1.x86_64.rpm" ]; then
    echo "Error: Can not find the rpm package to publish"
    exit 1
fi 
mkdir -pv "./repos/rpm/fedora/releases/28/x86_64"
cp -v "$HOST_FOLDER/rpm-packages/Fedora-28/samba-exporter-${rpmVersion}-1.x86_64.rpm" "./repos/rpm/fedora/releases/28/x86_64/"
createrepo_c "./repos/rpm/fedora/releases/28/x86_64/"
if [ "$?" != "0" ]; then 
    echo "Error during createrepo for fc28"
    exit 1
fi

echo "# ###################################################################"
echo "Update the rpm repo with fc35 package"
if [ ! -f "$HOST_FOLDER/rpm-packages/Fedora-35/samba-exporter-${rpmVersion}-1.fc35.x86_64.rpm" ]; then
    echo "Error: Can not find the rpm package to publish"
    exit 1
fi 
mkdir -pv "./repos/rpm/fedora/releases/35/x86_64"
cp -v "$HOST_FOLDER/rpm-packages/Fedora-35/samba-exporter-${rpmVersion}-1.fc35.x86_64.rpm" "./repos/rpm/fedora/releases/35/x86_64/"
createrepo_c "./repos/rpm/fedora/releases/35/x86_64/"
if [ "$?" != "0" ]; then 
    echo "Error during createrepo for fc35"
    exit 1
fi


echo "# ###################################################################"
echo "git status"
git status

echo "# ###################################################################"
echo "Copy the updated repos out of the container"
mkdir -p "$HOST_FOLDER/pages/repos"
mv -v "./repos/debian" "$HOST_FOLDER/pages/repos"
mv -v "./repos/rpm" "$HOST_FOLDER/pages/repos"
chmod -R 777 "$HOST_FOLDER/pages"

exit 0
