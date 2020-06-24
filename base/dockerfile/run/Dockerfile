FROM ubuntu:bionic

RUN apt-get update && \
  apt-get install -y \
  ca-certificates \
  libssl1.1 \
  openssl \
  zlib1g \
  && rm -rf /var/lib/apt/lists/*
