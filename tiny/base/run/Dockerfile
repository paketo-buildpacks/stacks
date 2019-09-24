FROM ubuntu:bionic AS builder

ARG PACKAGE_LIST=packagelist

RUN apt-get update && \
  apt-get install -y xz-utils binutils

ADD install-certs.sh .
ADD download-and-install-package.sh .
ADD $PACKAGE_LIST packagelist
ADD files/passwd /tiny/etc/passwd
ADD files/nsswitch.conf /tiny/etc/nsswitch.conf
ADD files/group /tiny/etc/group

RUN mkdir -p /tiny/tmp \
    /tiny/home/nonroot \
    && export DEBIAN_FRONTEND=noninteractive \
    && apt -y update \
    && ./download-and-install-package.sh \
    && ./install-certs.sh

RUN rm /tiny/etc/os-release
ADD files/os-release /tiny/etc/os-release
RUN echo "" > /tiny/var/lib/dpkg/status

FROM scratch 

COPY --from=builder /tiny/ /
