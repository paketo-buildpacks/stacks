FROM ubuntu:bionic

RUN apt-get update && \
  apt-get install -y \
  ca-certificates \
  libssl1.1 \
  libyaml-0-2 \
  openssl \
  zlib1g \
  && rm -rf /var/lib/apt/lists/*
