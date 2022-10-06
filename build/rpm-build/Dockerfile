FROM fedora:36

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

RUN chmod 777 /BuildInDocker.sh