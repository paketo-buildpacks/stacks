ARG cnb_uid=1000
ARG cnb_gid=1000
ARG stack_id="io.buildpacks.stacks.tiny"

FROM cloudfoundry/cnb-base:tiny AS tiny-source

FROM ubuntu:bionic AS builder
COPY --from=tiny-source / /tiny/

RUN echo "cnb:x:${cnb_uid}:${cnb_gid}:cnb:/home/cnb:/sbin/nologin" >> /tiny/etc/passwd \
  && echo "cnb:x:${cnb_gid}" >> /tiny/etc/group \
  && mkdir -p /tiny/home/cnb

FROM cloudfoundry/cnb-base:tiny
COPY --from=builder /tiny/ /

USER cnb
LABEL io.buildpacks.stack.id="${stack_id}"
