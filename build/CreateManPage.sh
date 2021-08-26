#!/bin/bash
# ######################################################################################
# Copyright 2021 by tobi@backfrak.de. All
# rights reserved. Use of this source code is governed
# by a BSD-style license that can be found in the
# LICENSE file.
# ######################################################################################
# Script to convert markdown style man pages into actual man pages using ronn
# ######################################################################################

# ################################################################################################################
# variable asigenment
# ################################################################################################################
RONN=/usr/bin/ronn
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
BRANCH_ROOT="$SCRIPT_DIR/.."
SRC_DIR="$BRANCH_ROOT/src/man"
LOG_DIR="$BRANCH_ROOT/logs"
PACKAGE_NAME=$(cat $LOG_DIR/PackageName.txt)
DEFAULT_PACKAGE_ROOT="tmp"


# ################################################################################################################
# function definition
# ################################################################################################################
function print_usage()  {
    echo "Usage: $0 [package-root]"
    echo "-help             Print this help"
    echo "package-root     (optional) The path to the package. (default $DEFAULT_PACKAGE_ROOT)"
    echo ""
    echo "This script will convert markdown style man pages into actual man pages using ronn"
}

# ################################################################################################################
# functional code
# ################################################################################################################

if [ "$1" == "-help" ]; then
    print_usage
    exit 0
fi  

if [ "$1" != "" ]; then
    PACKAGE_ROOT="$BRANCH_ROOT/$1/$PACKAGE_NAME"
else
    PACKAGE_ROOT="$BRANCH_ROOT/$DEFAULT_PACKAGE_ROOT/$PACKAGE_NAME"
fi

if [ ! -f $RONN ];then 
    echo "ERROR: ronn package is not installed"
    exit 1
fi

pushd "$BRANCH_ROOT" >> /dev/null

# Clean old files
if [ -f "$SRC_DIR/samba_exporter.1" ];then
    rm -f "$SRC_DIR/samba_exporter.1"
fi 
if [ -f "$SRC_DIR/samba_exporter.1.html" ];then
    rm -f "$SRC_DIR/samba_exporter.1.html"
fi 
if [ -f "$SRC_DIR/samba_exporter.1.gz" ];then
    rm -f "$SRC_DIR/samba_exporter.1.gz"
fi 
if [ -f "$SRC_DIR/start_samba_statusd.1.html" ];then
    rm -f "$SRC_DIR/start_samba_statusd.1.html"
fi
if [ -f "$SRC_DIR/start_samba_statusd.1.gz" ];then
    rm -f "$SRC_DIR/start_samba_statusd.1.gz"
fi      
if [ -f "$SRC_DIR/start_samba_statusd.1" ];then
    rm -f "$SRC_DIR/start_samba_statusd.1"
fi  
if [ -f "$SRC_DIR/samba_statusd.1.html" ];then
    rm -f "$SRC_DIR/samba_statusd.1.html"
fi 
if [ -f "$SRC_DIR/samba_statusd.1.gz" ];then
    rm -f "$SRC_DIR/samba_statusd.1.gz"
fi    
if [ -f "$SRC_DIR/samba_statusd.1" ];then
    rm -f "$SRC_DIR/samba_statusd.1"
fi 


# Generate new files
$RONN "$SRC_DIR/samba_exporter.1.ronn"
gzip --keep "$SRC_DIR/samba_exporter.1"
$RONN "$SRC_DIR/samba_statusd.1.ronn"
gzip --keep "$SRC_DIR/samba_statusd.1"
$RONN "$SRC_DIR/start_samba_statusd.1.ronn"
gzip --keep "$SRC_DIR/start_samba_statusd.1"


# Install the man page into the package
echo "Install to tmp package $PACKAGE_ROOT"
mkdir -p "$PACKAGE_ROOT/usr/man/man1"
cp "$SRC_DIR/samba_statusd.1.gz" "$PACKAGE_ROOT/usr/man/man1"
cp "$SRC_DIR/start_samba_statusd.1.gz" "$PACKAGE_ROOT/usr/man/man1"
cp "$SRC_DIR/samba_exporter.1.gz" "$PACKAGE_ROOT/usr/man/man1"

popd
