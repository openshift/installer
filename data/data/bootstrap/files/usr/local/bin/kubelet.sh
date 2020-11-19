#!/usr/bin/env bash

. /usr/local/bin/bootstrap-service-record.sh

/usr/bin/hyperkube \
  kubelet \
    --anonymous-auth=false \
    --container-runtime=remote \
    --container-runtime-endpoint=/var/run/crio/crio.sock \
    --runtime-request-timeout=${KUBELET_RUNTIME_REQUEST_TIMEOUT} \
    --pod-manifest-path=/etc/kubernetes/manifests \
    --minimum-container-ttl-duration=6m0s \
    --cluster-domain=cluster.local \
    --cgroup-driver=systemd \
    --serialize-image-pulls=false \
    --v=2 \
    --volume-plugin-dir=/etc/kubernetes/kubelet-plugins/volume/exec \
    --pod-infra-container-image=${MACHINE_CONFIG_INFRA_IMAGE}
