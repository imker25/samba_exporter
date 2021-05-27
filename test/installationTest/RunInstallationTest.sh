#!/bin/bash

# Run some installation tests for the samba_exporter package
echo "SAMBA_EXPORTER_PACKAGE_NAME=$SAMBA_EXPORTER_PACKAGE_NAME"

if [ -f "${SAMBA_EXPORTER_PACKAGE_NAME}_amd64.deb" ]; then
    echo "Install package found on : ${SAMBA_EXPORTER_PACKAGE_NAME}_amd64.deb"
else
    echo "Error: Installer \"${SAMBA_EXPORTER_PACKAGE_NAME}_amd64.deb\" not found"
fi
