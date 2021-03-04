ARG base_image
FROM ${base_image}

ARG cnb_uid=1000
ARG cnb_gid=1000
ARG distro_name="Ubuntu"
ARG distro_version="18.04"
ARG homepage="https://github.com/paketo-buildpacks/stacks"
ARG maintainer="Paketo Buildpacks"
ARG package_metadata
ARG stack_id="io.buildpacks.stacks.bionic"
ARG description
ARG metadata
ARG mixins
ARG released

RUN groupadd cnb --gid ${cnb_gid} && \
  useradd --uid ${cnb_uid} --gid ${cnb_gid} -m -s /bin/bash cnb

USER ${cnb_uid}:${cnb_gid}
LABEL io.buildpacks.stack.description=${description}
LABEL io.buildpacks.stack.distro.name=${distro_name}
LABEL io.buildpacks.stack.distro.version=${distro_version}
LABEL io.buildpacks.stack.homepage=${homepage}
LABEL io.buildpacks.stack.id=${stack_id}
LABEL io.buildpacks.stack.maintainer=${maintainer}
LABEL io.buildpacks.stack.metadata=${metadata}
LABEL io.buildpacks.stack.mixins=${mixins}
LABEL io.buildpacks.stack.released=${released}
LABEL io.paketo.stack.packages=${package_metadata}
