#!/bin/bash
# ######################################################################################
# Copyright 2021 by tobi@backfrak.de. All
# rights reserved. Use of this source code is governed
# by a BSD-style license that can be found in the
# LICENSE file.
# ######################################################################################
# Script to build and run a docker container and do the following inside the container
# * clone the GitHub git repo
# * switch to the gh-pages banrch
# * update the content of repos/debian with new packages using reprepro
# * move the new conetent of repos/debian out of the container
# Once the container run is done the script does the following to patch the pages content
# * Copy the new repos/debian contetn to ./build/pages/site/repos/debian 
# # Copy the generated man pages from ./bin/deb-packages/man to ./build/pages/site/manpages/
# ######################################################################################


# ################################################################################################################
# function definition
# ################################################################################################################
function print_usage()  {
    echo "Script to transfer a github tag to launchpad and publish the package in a ppa"
    echo ""
    echo "Usage: $0 tag"
    echo "-help     Print this help"
    echo "tag       The tag on the github repo to import, e. g. 0.7.8"
    echo "The script expect the following environment variables to be set"
    echo "  LAUNCHPAD_GPG_KEY_PUB       Public GPG Key for the launchpad ppa"
    echo "  LAUNCHPAD_GPG_KEY_PRV       Private GPG Key for the launchpad ppa"    
}

# ################################################################################################################
# variable asigenment
# ################################################################################################################
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
BRANCH_ROOT="$SCRIPT_DIR/.."
WORK_DIR="$BRANCH_ROOT/build/additional-pages"
MOUNT_TO_DOCKER_DIR="$BRANCH_ROOT/bin"
MAN_PAGE_COPY_SOURCE="$BRANCH_ROOT/bin/deb-packages/man"
REPO_COPY_SOURCE="$MOUNT_TO_DOCKER_DIR/pages"
PAGES_COPY_TARGET="$BRANCH_ROOT/build/pages/site"

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
pushd "$WORK_DIR"
echo "Ensure target dir for pages exists"
mkdir -p "$PAGES_COPY_TARGET"
if [ ! -d "$PAGES_COPY_TARGET" ]; then
    echo "Error: Directory '$PAGES_COPY_TARGET' not found after creation" 
    popd
    exit 1
fi 

mkdir -p "$BRANCH_ROOT/logs"
if [ ! -d "$BRANCH_ROOT/logs" ]; then
    echo "Error: Directory '$BRANCH_ROOT/logs' not found after creation" 
    popd
    exit 1
fi

if [ -d "$REPO_COPY_SOURCE" ]; then
    echo "Clean debian repository tmp dir '$REPO_COPY_SOURCE'"
    rm -rf "$REPO_COPY_SOURCE"/*
fi

echo "# ###################################################################"
echo "Update repository with binariy packages for version '$tag'"
echo "# ###################################################################"
echo "Build the needed container from '$WORK_DIR/Dockerfile', logging to $BRANCH_ROOT/logs/docker-build-repoupdate.log"
docker build --file "$WORK_DIR/Dockerfile" --tag page-publish-container . > $BRANCH_ROOT/logs/docker-build-repoupdate.log 2>&1
if [ "$?" != "0" ]; then 
    echo "Error during docker build"
    popd
    exit 1
fi

echo "# ###################################################################"
echo "Run the docker container"
docker run --mount type=bind,source="$MOUNT_TO_DOCKER_DIR",target="/build_results" \
    --env LAUNCHPAD_GPG_KEY_PUB="$LAUNCHPAD_GPG_KEY_PUB" \
    --env LAUNCHPAD_GPG_KEY_PRV="$LAUNCHPAD_GPG_KEY_PRV" \
    -i page-publish-container \
    /bin/bash -c "/UpdatePagesRepo.sh $tag"
if [ "$?" != "0" ]; then 
    echo "Error during docker run"
    popd
    echo "# ###################################################################"
    echo "Delete containers"
    docker rmi -f $(docker images --filter=reference="page-publish*" -q) 
    exit 1
fi

echo "# ###################################################################"
echo "Delete containers"
docker rmi -f $(docker images --filter=reference="page-publish*" -q) 

echo "# ###################################################################"
echo "Copy files to '$PAGES_COPY_TARGET'"
if [ ! -d "$REPO_COPY_SOURCE/repos/debian" ]; then 
    echo "Error: The debian repository source dir '$REPO_COPY_SOURCE/repos/debian' was not found"
    popd
    exit 1
fi

if [ -d "$PAGES_COPY_TARGET/repos/debian" ]; then
    echo "Clean the debian repository dir target location '$PAGES_COPY_TARGET/repos/debian'"
    rm -rf "$PAGES_COPY_TARGET/repos/debian" 
fi

echo "Move the debian repository to pages dir"
mv -v "$REPO_COPY_SOURCE/repos" "$PAGES_COPY_TARGET"
if [ ! -d "$PAGES_COPY_TARGET/repos/debian" ]; then 
    echo "Error: The debian repository dir target location '$PAGES_COPY_TARGET/repos/debian' was not found "
    popd
    exit 1
fi

echo "Copy the html man pages"
mkdir -p "$PAGES_COPY_TARGET/manpages/"
cp -rv "$MAN_PAGE_COPY_SOURCE"/*.html "$PAGES_COPY_TARGET/manpages/"

echo "Copy the index.html"
cp -v "$BRANCH_ROOT/src/page/index.html" "$PAGES_COPY_TARGET"

popd

echo "# ###################################################################"
echo "done"