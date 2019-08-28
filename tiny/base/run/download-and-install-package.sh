#!/bin/bash

set -eu -o pipefail

mkdir -p "/tiny/var/lib/dpkg/status.d"
apt download $(cat packagelist | tr '\n' ' ')

while read PACKAGE; do 
    echo "installing $PACKAGE..."
    ar p $PACKAGE*.deb data.tar.xz | unxz | tar x -C /tiny
    ar x $PACKAGE*.deb
    if [[ "$(ls control.*)" == "control.tar.xz" ]]; then
        unxz < control.tar.xz | tar x ./control && mv ./control "/tiny/var/lib/dpkg/status.d/$PACKAGE"
        rm control.tar.xz
    elif [[ "$(ls control.*)" == "control.tar.gz" ]]; then
        tar xzf control.tar.gz ./control && mv ./control "/tiny/var/lib/dpkg/status.d/$PACKAGE"
        rm control.tar.gz
    fi
    rm debian-binary data.tar.xz
done < packagelist
