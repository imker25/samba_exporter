FROM fedora:40

RUN yum -y update && \
    yum clean all
RUN yum -y install redhat-lsb \
                rubygem-ronn-ng \
                rpmdevtools \
                rpmlint \
                rpm-build \
                gzip \
                curl \
                man-db \
                binutils \
                wget \
                procps-ng \
                golang && \
    yum clean all

COPY BuildInDocker.sh /BuildInDocker.sh

ARG USER
ARG UID
ARG GID

RUN chmod 777 /BuildInDocker.sh

RUN mkdir -p /home/${USER} &&\
    chmod 770 /home/${USER} && \
    chown ${UID}:${GID} /home/${USER} && \
    chmod 777 /var/lib/rpm/