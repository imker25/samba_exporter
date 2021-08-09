#!/bin/bash

echo " GITHUB_REF= $GITHUB_REF"
packageName=$(cat logs/PackageName.txt)
echo "SAMBA_EXPORTER_PACKAGE_NAME=$packageName"
echo "SAMBA_EXPORTER_PACKAGE_NAME=$packageName" >> $GITHUB_ENV   
packageVersion=$(cat logs/Version.txt)  
echo "SAMBA_EXPORTER_VERSION=$packageVersion" 
echo "SAMBA_EXPORTER_VERSION=$packageVersion"  >> $GITHUB_ENV