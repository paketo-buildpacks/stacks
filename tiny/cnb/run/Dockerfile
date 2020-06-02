ARG base_image
FROM ${base_image} AS tiny-source
FROM ubuntu:bionic AS builder
COPY --from=tiny-source / /tiny/

ARG cnb_uid=1000
ARG cnb_gid=1000

RUN echo "cnb:x:${cnb_uid}:${cnb_gid}:cnb:/home/cnb:/sbin/nologin" >> /tiny/etc/passwd \
  && echo "cnb:x:${cnb_gid}" >> /tiny/etc/group \
  && mkdir -p /tiny/home/cnb

FROM ${base_image}
COPY --from=builder /tiny/ /
ARG stack_id="io.paketo.stacks.tiny"
ARG cnb_uid=1000
ARG cnb_gid=1000

USER ${cnb_uid}:${cnb_gid}
LABEL io.buildpacks.stack.id="${stack_id}"

LABEL io.buildpacks.stack.mixins="[\
\"ca-certificates\", \
\"libssl1.1\", \
\"openssl\", \
\"run:ca-certificates\", \
\"run:libssl1.1\", \
\"run:openssl\"]"
