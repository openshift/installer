FROM registry.fedoraproject.org/fedora

WORKDIR /opt/app-root/src

ENV GOPATH=/opt/app-root
ENV GOBIN=/opt/app-root/bin
ENV KUBECONFIG=/opt/app-root/src/github.com/openshift/installer/auth/kubeconfig

COPY ./ssh /root/.ssh
COPY ./scripts_build /scripts
COPY ./libvirt_config ./libvirt_config

RUN /scripts/build-stage0.sh
RUN /scripts/build-stage1.sh
RUN /scripts/build-stage2.sh

COPY ./scripts_command /scripts
COPY ./scripts_config /root

CMD [ "/scripts/cmd.sh" ]

