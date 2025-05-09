#!/bin/bash
# ######################################################################################
# Copyright 2022 by tobi@backfrak.de. All
# rights reserved. Use of this source code is governed
# by a BSD-style license that can be found in the
# LICENSE file.
# ######################################################################################
# Script to build and test the project in a fedora container with creating a binary rpm at the end

# ################################################################################################################
# variable asigenment
# ################################################################################################################
projcetRoot=/build_area

# ################################################################################################################
# functional code
# ################################################################################################################
echo ""
echo "# ###################################################################"
echo "Build in container - started"
echo "# ###################################################################"
if [ ! -d "$projcetRoot" ]; then
    echo "Error: Can not find the sources dir '$projcetRoot'"
    exit 1
fi

echo "Ensure ./build.sh can use the sources"
echo "git config --global --add safe.directory /build_area"
git config --global --add safe.directory /build_area

pushd "$projcetRoot"

echo ""
echo "# ###################################################################"
echo "Compile and unit test with ./build.sh"
echo "# ###################################################################"
./build.sh test preparePack
if [ "$?" != "0" ]; then 
    echo "Error: Compile and test run failed"
    popd
    exit 1
fi

echo ""
echo "# ###################################################################"
echo "Create the man pages"
echo "# ###################################################################"
./build/CreateManPage.sh 
if [ "$?" != "0" ]; then 
    echo "Error: Man page creation failed"
    popd
    exit 1
fi

echo ""
echo "# ###################################################################"
echo "Run the integration tests"
echo "# ###################################################################"
./test/integrationTest/scripts/RunIntegrationTests.sh
if [ "$?" != "0" ]; then 
    echo "Error: Integration tests failed"
    popd
    exit 1
fi

echo ""
echo "# ###################################################################"
echo "Prepare for rpm packaging"
echo "# ###################################################################"
rpmdev-setuptree
rpmVersion=$(cat ./logs/ShortVersion.txt)
fullVersion=$(cat ./logs/PackageName.txt)
if [ "$rpmVersion" == "" ]; then
    echo "Error: Can not read the package version from './logs/ShortVersion.txt'"
fi
if [ "$fullVersion" == "" ]; then
    echo "Error: Can not read the full version from './logs/PackageName.txt'"
fi
echo "RPM Version will be: '$rpmVersion'"
mkdir -pv "$HOME/rpmbuild/PREBINROOT/"
mv -v "./tmp/${fullVersion}/"* "$HOME/rpmbuild/PREBINROOT/"
pushd "$HOME/rpmbuild/"
sed -i "s/Version: x.x.x/Version: $rpmVersion/g" ./PREBINROOT/samba-exporter.spec
mv -v ./PREBINROOT/samba-exporter.spec ./SPECS/samba-exporter.spec

echo ""
echo "# ###################################################################"
echo "Run rpm packaging"
echo "# ###################################################################"
rpmbuild -bb ./SPECS/samba-exporter.spec
if [ "$?" != "0" ]; then 
    echo "Error: RPM creation failed"
    popd
    popd
    exit 1
fi
popd

if [ -f "$HOME/samba-exporter-${rpmVersion}-1.x86_64.rpm" ]; then
    mv -v "$HOME/samba-exporter-${rpmVersion}-1.x86_64.rpm" "./bin/"
else
    echo "Error: Can not find the package file '$HOME/samba-exporter-${rpmVersion}-1.x86_64.rpm'"
    popd
    exit 1
fi

popd

echo ""
echo "# ###################################################################"
echo "Build in container - done"
echo "# ###################################################################"
exit 0
