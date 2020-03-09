#!/bin/bash

set -eu -o pipefail

CERTS_DIR=/tiny/usr/share/ca-certificates
CERT_FILE=/tiny/etc/ssl/certs/ca-certificates.crt
mkdir -p $(dirname $CERT_FILE)

CERTS=$(find "${CERTS_DIR}" -type f | sort)
for cert in ${CERTS}; do
  cat "${cert}" >> "${CERT_FILE}"
done

rm -rf "${CERTS_DIR}"
