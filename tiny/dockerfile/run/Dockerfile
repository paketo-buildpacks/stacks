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
    && chown 65532:65532 /tiny/home/nonroot \
    && export DEBIAN_FRONTEND=noninteractive \
    && apt -y update \
    && ./download-and-install-package.sh \
    && ./install-certs.sh

RUN find /tiny/usr/share/doc/*/* ! -name copyright | xargs rm -rf && \
  rm -rf \
    /tiny/etc/update-motd.d/* \
    /tiny/usr/share/man/* \
    /tiny/usr/share/lintian/*

ADD files/os-release /tmp/etc/os-release

RUN grep -v 'PRETTY_NAME=' "/tiny/etc/os-release" \
      | grep -v 'HOME_URL=' \
      | grep -v 'SUPPORT_URL=' \
      | grep -v 'BUG_REPORT_URL=' \
      | tee /tmp/etc/os-release-upstream \
    && rm /tiny/etc/os-release \
    && cat /tmp/etc/os-release-upstream /tmp/etc/os-release \
      | tee /tiny/etc/os-release

RUN echo "" > /tiny/var/lib/dpkg/status

FROM scratch

COPY --from=builder /tiny/ /
