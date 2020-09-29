FROM golang:alpine

COPY entrypoint /entrypoint
COPY release-notes-generator /tmp/release-notes-generator

RUN cd /tmp/release-notes-generator && go build -o /release-notes-generator .

RUN apk add \
    bash \
  && rm -rf /var/cache/apk/*

ENTRYPOINT ["/entrypoint"]
