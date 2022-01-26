FROM fedora:35

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
                java-11-openjdk-headless \
                golang && \
    yum clean all

COPY BuildInDocker.sh /BuildInDocker.sh

RUN chmod 777 /BuildInDocker.sh