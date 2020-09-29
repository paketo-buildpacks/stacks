FROM alpine

COPY entrypoint /entrypoint

RUN apk add \
    bash \
    git \
  && rm -rf /var/cache/apk/*

ENTRYPOINT ["/entrypoint"]
