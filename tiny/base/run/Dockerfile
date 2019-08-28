FROM cloudfoundry/cflinuxfs3:0.118.0 AS builder

ARG PACKAGE_LIST=packagelist

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
    && ./install-certs.sh \
    && ./download-and-install-package.sh

RUN rm /tiny/etc/os-release
ADD files/os-release /tiny/etc/os-release
RUN echo "" > /tiny/var/lib/dpkg/status

FROM scratch 

COPY --from=builder /tiny/ /
