# Container image that runs your code
FROM ubuntu:focal
ENV DEBIAN_FRONTEND="noninteractive"
ENV TZ="Europe/London"
RUN apt-get update && apt-get dist-upgrade -y && apt-get install -y wget \
                                        curl  \
                                        libc6  \
                                        golang-1.16 \
                                        gzip  \
                                        ronn\
                                        debhelper \
                                        golang-any \
                                        golang-github-go-kit-kit-dev \
                                        golang-github-prometheus-client-golang-dev \
                                        golang-github-prometheus-common-dev \
                                        golang-gopkg-alecthomas-kingpin.v2-dev \
                                        dh-golang \
                                        debhelper \ 
                                        dh-make \
                                        lintian \
                                        git-buildpackage \
                                        git \
                                        openssh-client \
                                        gpg

# Copies your code file from your action repository to the filesystem path `/` of the container
COPY PublishLaunchpad.sh /PublishLaunchpad.sh
RUN chmod 770 /PublishLaunchpad.sh
RUN git config --global user.name "Tobias Zellner"
RUN git config --global user.email imker@bienekaefig.de