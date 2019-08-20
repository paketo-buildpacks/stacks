ARG base_image
FROM ${base_image}

ARG cnb_uid=1000
ARG cnb_gid=1000
ARG stack_id="io.buildpacks.stacks.bionic"

RUN groupadd cnb --gid ${cnb_gid} && \
  useradd --uid ${cnb_uid} --gid ${cnb_gid} -m -s /bin/bash cnb

USER cnb
LABEL io.buildpacks.stack.id="${stack_id}"
