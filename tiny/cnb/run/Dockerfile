ARG base_image
FROM ${base_image} AS tiny-source
FROM ubuntu:bionic AS builder
COPY --from=tiny-source / /tiny/

ARG cnb_uid=1000
ARG cnb_gid=1000

RUN echo "cnb:x:${cnb_uid}:${cnb_gid}:cnb:/home/cnb:/sbin/nologin" >> /tiny/etc/passwd \
  && echo "cnb:x:${cnb_gid}" >> /tiny/etc/group \
  && mkdir -p /tiny/home/cnb \
  && chown ${cnb_uid}:${cnb_gid} /tiny/home/cnb

FROM ${base_image}
COPY --from=builder /tiny/ /
ARG cnb_uid=1000
ARG cnb_gid=1000
ARG description="Distroless-like bionic + glibc + openssl + CA certs"
ARG distro_name="Ubuntu"
ARG distro_version="18.04"
ARG homepage="https://github.com/paketo-buildpacks/stacks"
ARG maintainer="Paketo Buildpacks"
ARG package_metadata
ARG stack_id="io.paketo.stacks.tiny"
ARG metadata
ARG released
ARG mixins

USER ${cnb_uid}:${cnb_gid}
LABEL io.buildpacks.stack.description=${description}
LABEL io.buildpacks.stack.distro.name=${distro_name}
LABEL io.buildpacks.stack.distro.version=${distro_version}
LABEL io.buildpacks.stack.homepage=${homepage}
LABEL io.buildpacks.stack.id=${stack_id}
LABEL io.buildpacks.stack.maintainer=${maintainer}
LABEL io.buildpacks.stack.metadata=${metadata}
LABEL io.buildpacks.stack.released=${released}
LABEL io.buildpacks.stack.mixins=${mixins}
LABEL io.paketo.stack.packages=${package_metadata}
