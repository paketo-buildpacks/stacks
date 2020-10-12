FROM ubuntu:bionic

ARG sources
ARG packages
ARG package_args='--allow-downgrades --allow-remove-essential --allow-change-held-packages --no-install-recommends'

RUN echo "$sources" > /etc/apt/sources.list

RUN echo "debconf debconf/frontend select noninteractive" | debconf-set-selections && \
  export DEBIAN_FRONTEND=noninteractive && \
  apt-get -y $package_args update && \
  apt-get -y $package_args install $packages && \
  rm -rf /var/lib/apt/lists/* /tmp/*

RUN curl -sL -o /usr/local/bin/yj https://github.com/sclevine/yj/releases/latest/download/yj-linux \
  && chmod +x /usr/local/bin/yj

