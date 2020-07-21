ARG base_image
FROM ${base_image}

ARG cnb_uid=2000
ARG cnb_gid=2000
ARG stack_id="org.cloudfoundry.stacks.cflinuxfs3"

USER ${cnb_uid}:${cnb_gid}
LABEL io.buildpacks.stack.id=${stack_id}
