#!/bin/bash

set -eu -o pipefail

CERT_TMPDIR=$(mktemp -d)
apt download ca-certificates
ar -x ca-certificates*.deb data.tar.xz
tar -xf data.tar.xz -C "${CERT_TMPDIR}" ./usr/share/ca-certificates
tar -xf data.tar.xz -C /tiny/ ./usr/share/doc/ca-certificates/copyright

CERT_FILE=/tiny/etc/ssl/certs/ca-certificates.crt
mkdir -p $(dirname $CERT_FILE)

CERTS=$(find "${CERT_TMPDIR}/usr/share/ca-certificates" -type f | sort)
for cert in ${CERTS}; do
  cat "${cert}" >> "${CERT_FILE}"
done
