#!/bin/bash

set -eu -o pipefail

mkdir -p "/tiny/var/lib/dpkg/status.d"
apt download $(cat packagelist | tr '\n' ' ')

while read PACKAGE; do 
    echo "installing $PACKAGE..."
    ar p $PACKAGE*.deb data.tar.xz | unxz | tar x -C /tiny
    ar x $PACKAGE*.deb
    source_package="$(dpkg-deb --showformat='${source:Package}' -W "$PACKAGE"*.deb)"
    source_version="$(dpkg-deb --showformat='${source:Version}' -W "$PACKAGE"*.deb)"
    source_upstream_version="$(dpkg-deb --showformat='${source:Upstream-Version}' -W "$PACKAGE"*.deb)"
    if [[ "$(ls control.*)" == "control.tar.xz" ]]; then
        unxz < control.tar.xz | tar x ./control && mv ./control "/tiny/var/lib/dpkg/status.d/$PACKAGE"
        echo "Source-Package: ${source_package}" >> "/tiny/var/lib/dpkg/status.d/$PACKAGE"
        echo "Source-Version: ${source_version}" >> "/tiny/var/lib/dpkg/status.d/$PACKAGE"
        echo "Source-Upstream-Version: ${source_upstream_version}" >> "/tiny/var/lib/dpkg/status.d/$PACKAGE"
        rm control.tar.xz
    elif [[ "$(ls control.*)" == "control.tar.gz" ]]; then
        tar xzf control.tar.gz ./control && mv ./control "/tiny/var/lib/dpkg/status.d/$PACKAGE"
        echo "Source-Package: ${source_package}" >> "/tiny/var/lib/dpkg/status.d/$PACKAGE"
        echo "Source-Version: ${source_version}" >> "/tiny/var/lib/dpkg/status.d/$PACKAGE"
        echo "Source-Upstream-Version: ${source_upstream_version}" >> "/tiny/var/lib/dpkg/status.d/$PACKAGE"
        rm control.tar.gz
    fi
    rm debian-binary data.tar.xz
done < packagelist
