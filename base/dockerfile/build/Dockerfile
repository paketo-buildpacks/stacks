FROM ubuntu:bionic

RUN apt-get update && \
  apt-get install -y \
  build-essential \
  ca-certificates \
  curl \
  git \
  jq \
  libgmp-dev \
  libssl1.1 \
  libyaml-0-2 \
  openssl \
  xz-utils \
  zlib1g-dev \
  && rm -rf /var/lib/apt/lists/*

RUN curl -sL -o /usr/local/bin/yj https://github.com/sclevine/yj/releases/latest/download/yj-linux \
  && chmod +x /usr/local/bin/yj
