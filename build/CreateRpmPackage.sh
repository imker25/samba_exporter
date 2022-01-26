#!/bin/bash
# ######################################################################################
# Copyright 2022 by tobi@backfrak.de. All
# rights reserved. Use of this source code is governed
# by a BSD-style license that can be found in the
# LICENSE file.
# ######################################################################################
# Script to build and test the project using a fedora container with creating a binary rpm at the end

# ################################################################################################################
# variable asigenment
# ################################################################################################################
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
BRANCH_ROOT="$SCRIPT_DIR/.."

# ################################################################################################################
# functional code
# ################################################################################################################
echo "Create output folders"
echo "# ###################################################################"
mkdir -vp "$BRANCH_ROOT/bin"
mkdir -vp "$BRANCH_ROOT/logs"
mkdir -vp "$BRANCH_ROOT/tmp"

echo ""
echo "# ###################################################################"
echo "Build the container to run the package build"
echo "docker build --file \"$SCRIPT_DIR/rpm-build/Dockerfile\" --tag rpm-build \"$SCRIPT_DIR/rpm-build/\""
echo "Loggig to '$BRANCH_ROOT/logs/docker-build-fedora.log'"
docker build --file "$SCRIPT_DIR/rpm-build/Dockerfile" --tag rpm-build "$SCRIPT_DIR/rpm-build/" > $BRANCH_ROOT/logs/docker-build-fedora.log 2>&1
if [ "$?" != "0" ]; then 
    echo "Error: Docker container build failed"
    exit 1
fi
# Move the docker build log to other folder, so it will not be deleted during the "gradle clean" when the container runs
mv -v "$BRANCH_ROOT/logs/docker-build-fedora.log" "/tmp" 
echo ""
echo "# ###################################################################"
buildFailed="false"
echo "Run the container to create the rpm package"
docker run --mount type=bind,source="$BRANCH_ROOT",target="/build_area" \
            -i rpm-build \
            /bin/bash -c "/BuildInDocker.sh"

if [ "$?" != "0" ]; then    
    buildFailed="true"
fi
# Move the docker build log back to log folder
mv -v "/tmp/docker-build-fedora.log" "$BRANCH_ROOT/logs" 

echo "Delete the container image used for this build"
docker rmi -f $(docker images --filter=reference="rpm-build" -q) 

if [ "$buildFailed" == "true" ]; then
    echo "Error: RPM build in container failed"
    exit 1
fi

echo ""
echo "# ###################################################################"
echo "All done"
echo "# ###################################################################"
exit 0