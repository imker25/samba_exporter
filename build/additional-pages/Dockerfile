# ######################################################################################
# Copyright 2021 by tobi@backfrak.de. All
# rights reserved. Use of this source code is governed
# by a BSD-style license that can be found in the
# LICENSE file.
# ######################################################################################
# Container image used by ./../PatchPagesContent.sh to transfer the sources to the launchpad-ppa
# The actual transformation is done by ./PatchPagesContent.sh when this container get started
# ######################################################################################

FROM debian:trixie

# Setup the system and install needed packages
ENV DEBIAN_FRONTEND="noninteractive"
ENV TZ="Europe/London"
RUN apt-get update && apt-get dist-upgrade -y && apt-get install -y wget \
                                        curl  \
                                        git \
                                        reprepro \
                                        createrepo-c \
                                        gpg

# Copy the script
COPY UpdatePagesRepo.sh /UpdatePagesRepo.sh
RUN chmod 770 /UpdatePagesRepo.sh

# Configure git
RUN git config --global user.name "Tobias Zellner" && \
    git config --global user.email imker@bienekaefig.de