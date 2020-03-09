ARG base_image
FROM ${base_image}

ARG cnb_uid=2000
ARG cnb_gid=2000
ARG stack_id="org.cloudfoundry.stacks.cflinuxfs3"

ENV CNB_USER_ID=${cnb_uid}
ENV CNB_GROUP_ID=${cnb_gid}
ENV CNB_STACK_ID=${stack_id}

USER ${cnb_uid}:${cnb_gid}
LABEL io.buildpacks.stack.id=${stack_id}
