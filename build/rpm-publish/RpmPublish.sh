#!/bin/bash
# ######################################################################################
# Copyright 2022 by tobi@backfrak.de. All
# rights reserved. Use of this source code is governed
# by a BSD-style license that can be found in the
# LICENSE file.
# ######################################################################################
# Script to do the following inside a container
# * import the given samba_exporter github sources to the rpm build area
# * do the needed conversation steps, so rpm package build can run
# * run rpm binary package  build
# * run rpm source package  build 
# ######################################################################################

# ################################################################################################################
# function definition
# ################################################################################################################
function print_usage()  {
    echo "Script to create the release RPM"
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

# ################################################################################################################
# variable asigenment
# ################################################################################################################
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
GITHUB_RELEASE_URL="https://github.com/imker25/samba_exporter/archive/refs/tags"

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
    preRelease="true"
    rpmVersion=${tag/-pre/}
else 
    preRelease="false"
    rpmVersion=${tag}
fi

distribution=$(lsb_release -is)
distVersionNumber=$(lsb_release -rs)

# ################################################################################################################
# functional code
# ################################################################################################################

if [ "$dryRun" == "false" ]; then
    echo "Release 'samba-exporter-${tag}' for $distribution $distVersionNumber as RPM Version $rpmVersion"
else
    echo "Dry run: Release 'samba-exporter-${tag}' for $distribution $distVersionNumber as RPM Version $rpmVersion"
fi 

echo "Prepare for operation"
mkdir -p ~/.ssh
echo "$COPR_SSH_ID_PUB" > ~/.ssh/id_rsa.pub
chmod 600 ~/.ssh/id_rsa.pub
echo "$COPR_SSH_ID_PRV" > ~/.ssh/id_rsa
chmod 600 ~/.ssh/id_rsa
mkdir -p ~/.gnupg
chmod 600 ~/.gnupg
echo "$COPR_GPG_KEY_PUB" > ~/.gnupg/imker-bienenkaefig.pub.asc
echo "$COPR_GPG_KEY_PRV" > ~/.gnupg/imker-bienenkaefig.asc

gpg --import --batch --no-tty ~/.gnupg/imker-bienenkaefig.asc
gpg --edit-key --batch --no-tty  CB6E90E9EC323850B16C1C14A38A1091C018AE68 trust quit
gpg --list-keys --batch --no-tty 

echo "Create rpm build folders"
mkdir -pv ~/rpmbuild/{BUILD,RPMS,SOURCES,SPECS,SRPMS}
cd ~/rpmbuild

echo "# ###################################################################"
echo "Download the source zip"
echo "wget -O ~/rpmbuild/SOURCES/${tag}.tar.gz \"$GITHUB_RELEASE_URL/$tag.tar.gz\""
wget -O ~/rpmbuild/SOURCES/${tag}.tar.gz "$GITHUB_RELEASE_URL/$tag.tar.gz"
if [ "$?" != "0" ]; then 
    echo "Error during sources download"
    return 1
fi

if [ ! -f ~/rpmbuild/SOURCES/${tag}.tar.gz ]; then
    echo "Error: Can not find '~/rpmbuild/SOURCES/${tag}.tar.gz'"
    return 1
fi 

echo "# ###################################################################"
echo "Prepare for build"
pushd ~/rpmbuild/SOURCES
tar -zxvf ~/rpmbuild/SOURCES/${tag}.tar.gz samba_exporter-${tag}/install/fedora/samba-exporter.from_source.spec
cp -v samba_exporter-${tag}/install/fedora/samba-exporter.from_source.spec ~/rpmbuild/SPECS/
popd

if [ ! -f ~/rpmbuild/SPECS/samba-exporter.from_source.spec ]; then
    echo "Error: Can not find '~/rpmbuild/SPECS/samba-exporter.from_source.spec'"
    return 1
fi 

echo "# ###################################################################"
echo "Patch the spec file"
sed -i "s/x.x.x-pre/${tag}/g" ~/rpmbuild/SPECS/samba-exporter.from_source.spec
sed -i "s/X.X.X-pre/${tag}/g" ~/rpmbuild/SPECS/samba-exporter.from_source.spec
sed -i "s/x.x.x/${rpmVersion}/g" ~/rpmbuild/SPECS/samba-exporter.from_source.spec

echo "# ###################################################################"
echo "~/rpmbuild/SPECS/samba-exporter.from_source.spec after modification"
echo "# ###################################################################"
cat ~/rpmbuild/SPECS/samba-exporter.from_source.spec
echo "# ###################################################################"


echo "Build the source package"
echo "rpmbuild -bs ~/rpmbuild/SPECS/samba-exporter.from_source.spec"
rpmbuild -bs ~/rpmbuild/SPECS/samba-exporter.from_source.spec
if [ "$?" != "0" ]; then 
    echo "Error during sources package build"
    return 1
fi

if [ ! -f ~/rpmbuild/SRPMS/samba-exporter-${rpmVersion}-1.src.rpm ]; then
    echo "Error: Can not find the source package '~/rpmbuild/SRPMS/samba-exporter-${rpmVersion}-1.src.rpm'"
    return 1
fi 

echo "# ###################################################################"
echo "Build the binary package"
echo "rpmbuild --rebuild ~/rpmbuild/SRPMS/samba-exporter-${rpmVersion}-1.src.rpm"
rpmbuild --rebuild ~/rpmbuild/SRPMS/samba-exporter-${rpmVersion}-1.src.rpm
if [ "$?" != "0" ]; then 
    echo "Error during binary package build"
    return 1
fi

if [ ! -f ~/samba-exporter-${rpmVersion}-1.x86_64.rpm ];then 
    echo "Error: Can not find the binary package '~/samba-exporter-${rpmVersion}-1.x86_64.rpm'"
fi 

cp -v ~/samba-exporter-${rpmVersion}-1.x86_64.rpm /build_results/
chmod -R 777 /build_results/*

exit 0