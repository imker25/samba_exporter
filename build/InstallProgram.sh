#!/bin/bash
# ######################################################################################
# Copyright 2021 by tobi@backfrak.de. All
# rights reserved. Use of this source code is governed
# by a BSD-style license that can be found in the
# LICENSE file.
# ######################################################################################
# Script to install the program files. Used by debian/rules as well as ./build.sh
# ######################################################################################

# ################################################################################################################
# variable asigenment
# ################################################################################################################
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
BRANCH_ROOT="$SCRIPT_DIR/.."

# ################################################################################################################
# function definition
# ################################################################################################################
function print_usage()  {
    echo "Usage: $0 source-root binary-root package-root"
    echo "-help             Print this help"
    echo "source-root       The root folder from there all copy source pathes for sources are calculated"
    echo "binary-root       The root folder from there all copy source pathes for binaries are calculated"
    echo "package-root      The path to the package input root folder"
    echo ""
    echo "Script to install the program files. Used by debian/rules as well as ./build.sh build "
}

# ################################################################################################################
# functional code
# ################################################################################################################

if [ "$1" == "-help" ]; then
    print_usage
    exit 0
fi  

if [ "$1" != "" ]; then
    COPY_SOURCE="$1"
else
    echo "Error: No source-root parameter given"
    exit 1
fi

if [ "$2" != "" ]; then
    BIN_COPY_SOURCE="$2"
else
    echo "Error: No binary-root parameter given"
    exit 1
fi

if [ "$3" != "" ]; then
    PACKAGE_ROOT="$3"
else
    echo "Error: No package-root parameter given"
    exit 1
fi

echo "Install tree from ${COPY_SOURCE}"
install -d -m 775 "${PACKAGE_ROOT}/usr/bin"
install -s -m 775 "${BIN_COPY_SOURCE}/samba_exporter" "${PACKAGE_ROOT}/usr/bin/samba_exporter"
install -s -m 775 "${BIN_COPY_SOURCE}/samba_statusd" "${PACKAGE_ROOT}/usr/bin/samba_statusd"
install -m 775 "${COPY_SOURCE}/install/usr/bin/start_samba_statusd" "${PACKAGE_ROOT}/usr/bin/start_samba_statusd"
install -d -m 775 "${PACKAGE_ROOT}/lib/systemd/system"
install -m 664 "${COPY_SOURCE}/install/lib/systemd/system/samba_exporter.service" "${PACKAGE_ROOT}/lib/systemd/system/samba_exporter.service"
install -m 664 "${COPY_SOURCE}/install/lib/systemd/system/samba_statusd.service" "${PACKAGE_ROOT}/lib/systemd/system/samba_statusd.service"
install -d -m 775 "${PACKAGE_ROOT}/etc/default"
install -m 664 "${COPY_SOURCE}/install/etc/default/samba_exporter" "${PACKAGE_ROOT}/etc/default/samba_exporter"
install -m 664 "${COPY_SOURCE}/install/etc/default/samba_statusd" "${PACKAGE_ROOT}/etc/default/samba_statusd"
install -d -m 775 "${PACKAGE_ROOT}/usr/share/doc/samba-exporter/grafana"
install -m 664 "${COPY_SOURCE}/README.md" "${PACKAGE_ROOT}/usr/share/doc/samba-exporter/README.md"
install -m 664 "${COPY_SOURCE}/src/example/grafana/SambaService.json" "${PACKAGE_ROOT}/usr/share/doc/samba-exporter/grafana/SambaService.json"
install -d -m 775 "${PACKAGE_ROOT}/usr/share/doc/samba-exporter/docs/assets"
install -m 664 "${COPY_SOURCE}/docs/assets/Samba-Dashboard.png" "${PACKAGE_ROOT}/usr/share/doc/samba-exporter/docs/assets/Samba-Dashboard.png"
install -m 664 "${COPY_SOURCE}/docs/assets/samba-exporter.icon.png" "${PACKAGE_ROOT}/usr/share/doc/samba-exporter/docs/assets/samba-exporter.icon.png"
install -m 664 "${COPY_SOURCE}/docs/Index.md" "${PACKAGE_ROOT}/usr/share/doc/samba-exporter/docs/Index.md"
install -d -m 775 "${PACKAGE_ROOT}/usr/share/doc/samba-exporter/docs/DeveloperDocs"
install -m 664 "${COPY_SOURCE}/docs/DeveloperDocs/ActionsAndReleases.md" "${PACKAGE_ROOT}/usr/share/doc/samba-exporter/docs/DeveloperDocs/ActionsAndReleases.md"
install -m 664 "${COPY_SOURCE}/docs/DeveloperDocs/Compile.md" "${PACKAGE_ROOT}/usr/share/doc/samba-exporter/docs/DeveloperDocs/Compile.md"
install -m 664 "${COPY_SOURCE}/docs/DeveloperDocs/Hints.md" "${PACKAGE_ROOT}/usr/share/doc/samba-exporter/docs/DeveloperDocs/Hints.md"
install -d -m 775 "${PACKAGE_ROOT}/usr/share/doc/samba-exporter/docs/Installation"
install -m 664 "${COPY_SOURCE}/docs/Installation/InstallationGuide.md" "${PACKAGE_ROOT}/usr/share/doc/samba-exporter/docs/Installation/InstallationGuide.md"
install -m 664 "${COPY_SOURCE}/docs/Installation/SupportedVersions.md" "${PACKAGE_ROOT}/usr/share/doc/samba-exporter/docs/Installation/SupportedVersions.md"
install -d -m 775 "${PACKAGE_ROOT}/usr/share/doc/samba-exporter/docs/UserDocs"
install -m 664 "${COPY_SOURCE}/docs/UserDocs/Concept.md" "${PACKAGE_ROOT}/usr/share/doc/samba-exporter/docs/UserDocs/Concept.md"
install -m 664 "${COPY_SOURCE}/docs/UserDocs/ServiceIntegration.md" "${PACKAGE_ROOT}/usr/share/doc/samba-exporter/docs/UserDocs/ServiceIntegration.md"
install -m 664 "${COPY_SOURCE}/docs/UserDocs/UserGuide.md" "${PACKAGE_ROOT}/usr/share/doc/samba-exporter/docs/UserDocs/UserGuide.md"

