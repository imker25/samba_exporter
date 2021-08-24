#!/bin/bash
# ######################################################################################
# Copyright 2021 by tobi@backfrak.de. All
# rights reserved. Use of this source code is governed
# by a BSD-style license that can be found in the
# LICENSE file.
# ######################################################################################
# Script to read infos from files into the github environment during a actions workflow
# ######################################################################################

# Read package info from install/debian/control file
controlFile="./install/debian/control"
while IFS= read -r line
do
  if [[ $line = Description:* ]]; then
    desc=${line/Description: /}
    echo "SAMBA_EXPORTER_PACKAGE_DESCRIPTION=$desc"
    echo "SAMBA_EXPORTER_PACKAGE_DESCRIPTION=$desc" >> $GITHUB_ENV 
  fi
  if [[ $line = Depends:* ]]; then
    dep=${line/Depends: /}
    echo "SAMBA_EXPORTER_PACKAGE_DEPENDS=$dep"
    echo "SAMBA_EXPORTER_PACKAGE_DEPENDS=$dep" >> $GITHUB_ENV 
  fi
  if [[ $line = Maintainer:* ]]; then
    dev=${line/Maintainer: /}
    echo "SAMBA_EXPORTER_PACKAGE_MAINTAINER=$dev"
    echo "SAMBA_EXPORTER_PACKAGE_MAINTAINER=$dev" >> $GITHUB_ENV 
  fi  
  if [[ $line = Package:* ]]; then
    pack=${line/Package: /}
    echo "SAMBA_EXPORTER_PACKAGE=$pack"
    echo "SAMBA_EXPORTER_PACKAGE=$pack" >> $GITHUB_ENV 
  fi   
done < "$controlFile"

# Read version infos from build log files
echo " GITHUB_REF= $GITHUB_REF"
packageName=$(cat logs/PackageName.txt)
echo "SAMBA_EXPORTER_PACKAGE_NAME=$packageName"
echo "SAMBA_EXPORTER_PACKAGE_NAME=$packageName" >> $GITHUB_ENV   
packageVersion=$(cat logs/Version.txt)  
echo "SAMBA_EXPORTER_VERSION=$packageVersion" 
echo "SAMBA_EXPORTER_VERSION=$packageVersion"  >> $GITHUB_ENV