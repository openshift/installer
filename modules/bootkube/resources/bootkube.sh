#!/bin/bash
set -e

echo "Rendering Kubernetes core manifests..."

# shellcheck disable=SC2154
/usr/bin/docker run \
    --volume "$(pwd)":/assets:z \
    --volume /etc/kubernetes:/etc/kubernetes:z \
    "${kube_core_renderer_image}" \
    --config=/assets/kco-config.yaml \
    --output=/assets

echo "Rendering TNC manifests..."

# shellcheck disable=SC2154
/usr/bin/docker run \
    --user 0 \
    --volume "$(pwd)":/assets:z \
    "${tnc_operator_image}" \
    --config=/assets/tnco-config.yaml \
    --render-bootstrap=true \
    --render-output=/assets/tnc-bootstrap

mkdir -p /etc/kubernetes/manifests/
cp "$(pwd)/tnc-bootstrap/tectonic-node-controller-pod.yaml" /etc/kubernetes/manifests/
cp "$(pwd)/tnc-bootstrap/tectonic-node-controller-config.yaml" /etc/kubernetes/tnc-config

# We originally wanted to run the etcd cert signer as
# a static pod, but kubelet could't remove static pod
# when API server is not up, so we have to run this as
# docker container.
# See https://github.com/kubernetes/kubernetes/issues/43292

echo "Starting etcd certificate signer..."

# shellcheck disable=SC2154,SC2034
signer_id=$(/usr/bin/docker run -d \
    --tmpfs /tmp \
    --volume /opt/tectonic/tls:/opt/tectonic/tls:ro,z \
    --network host \
    "${etcd_cert_signer_image}" \
    serve \
    --cacrt=/opt/tectonic/tls/etcd-client-ca.crt \
    --cakey=/opt/tectonic/tls/etcd-client-ca.key \
    --servcrt=/opt/tectonic/tls/apiserver.crt \
    --servkey=/opt/tectonic/tls/apiserver.key \
    --address=0.0.0.0:6443 \
    --csrdir=/tmp \
    --peercertdur=26280h \
    --servercertdur=26280h)

echo "Waiting for etcd cluster..."

# Wait for the etcd cluster to come up.
i=0
while true; do
    set +e
    # shellcheck disable=SC2154,SC2086
    /usr/bin/docker run \
        --rm \
        --name etcdctl \
        --env ETCDCTL_API=3 \
        --volume /opt/tectonic/tls:/opt/tectonic/tls:ro,z \
        "${etcdctl_image}" \
        /usr/local/bin/etcdctl \
            --dial-timeout=10m \
            --cacert=/opt/tectonic/tls/etcd-client-ca.crt \
            --cert=/opt/tectonic/tls/etcd-client.crt \
            --key=/opt/tectonic/tls/etcd-client.key \
            --endpoints=${etcd_cluster} \
            endpoint health
    status=$?
    set -e

    if [ "$status" -eq 0 ]; then
      break
    fi

    i=$((i+1))
    [ $i -eq 10 ] && echo "etcdctl failed too many times." && exit 1

    echo "etcdctl failed. Retrying in 5 seconds..."
    sleep 5
done

echo "etcd cluster up. Killing etcd certificate signer..."

# shellcheck disable=SC2154,SC1083
/usr/bin/docker kill $${signer_id}
rm /etc/kubernetes/manifests/tectonic-node-controller-pod.yaml

cp -r "$(pwd)/bootstrap-configs" /etc/kubernetes/bootstrap-configs

echo "Starting bootkube..."

# shellcheck disable=SC2154
/usr/bin/docker run \
    --volume "$(pwd)":/assets:z \
    --volume /etc/kubernetes:/etc/kubernetes:z \
    --network=host \
    --entrypoint=/bootkube \
    "${bootkube_image}" \
    start --asset-dir=/assets
