FROM ubuntu:bionic

RUN apt-get update && \
  apt-get install -y \
  build-essential \
  ca-certificates \
  curl \
  git \
  libssl1.1 \
  openssl \
  xz-utils \
  && rm -rf /var/lib/apt/lists/*

COPY ./bin/jq-linux64 /usr/local/bin/jq
COPY ./bin/yj-linux /usr/local/bin/yj
