#!/usr/bin/env bash

WORKING_DIR=/home/core
REGISTRY_DIR=${WORKING_DIR}/registry
BASE_DOMAIN="ostest.test.metalkube.org"
LOCAL_REGISTRY_DNS_NAME=master-0.${BASE_DOMAIN}
LOCAL_REGISTRY_PORT=5000
REGISTRY_CRT=registry.2.crt
REGISTRY_KEY=registry.2.key
REGISTRY_USER=ocp-user
REGISTRY_PASS=ocp-pass
# DOCKER_REGISTRY_IMAGE=docker.io/library/registry:2
OC_MIRROR_FILE=${WORKING_DIR}/oc-mirror
REGISTRY_IMAGE=${WORKING_DIR}/registry_image.tar
HTPASSWD_FILE=${WORKING_DIR}/htpasswd
IMAGE_ARCHIVE=${WORKING_DIR}/mirror_seq*

function setup_podman_mirror_registry() {

    image=$1

    sudo mkdir -pv ${REGISTRY_DIR}/{auth,certs,data}
    sudo chown -R $USER:$GROUP ${REGISTRY_DIR}

    pushd $REGISTRY_DIR/certs

    openssl genrsa -out ${REGISTRY_DIR}/certs/${REGISTRY_KEY} 2048
    sudo sudo chmod 666 ${REGISTRY_DIR}/certs/${REGISTRY_KEY}

    # Format names as DNS:name1,DNS:name2
    SUBJECT_ALT_NAME="DNS:$(echo $LOCAL_REGISTRY_DNS_NAME | sed 's/ /,DNS:/g')"

    SSL_CONF=${REGISTRY_DIR}/certs/openssl.cnf
    cat > ${SSL_CONF} <<EOF
[req]
distinguished_name = req_distinguished_name
prompt = no

[req_distinguished_name]
C = US
ST = NC
L = Raleigh
O = Test Company
OU = Testing
CN = ${BASE_DOMAIN}

[SAN]
basicConstraints=CA:TRUE,pathlen:0
subjectAltName = ${SUBJECT_ALT_NAME}
EOF

    openssl req -x509 \
                -key ${REGISTRY_DIR}/certs/${REGISTRY_KEY} \
                -out  ${REGISTRY_DIR}/certs/${REGISTRY_CRT} \
                -days 365 \
                -config ${SSL_CONF} \
                -extensions SAN

    # Dump the certificate details to the log
    openssl x509 -in ${REGISTRY_DIR}/certs/${REGISTRY_CRT} -text
    popd

    # copy htpasswd, this will eventually be in the ISO
    cp ${HTPASSWD_FILE} ${REGISTRY_DIR}/auth/htpasswd

    sudo cp ${REGISTRY_DIR}/certs/${REGISTRY_CRT} /etc/pki/ca-trust/source/anchors/
    sudo update-ca-trust

    echo "Starting local registry for ${image}"
    sudo podman run -d --name registry --net=host --privileged \
        -v ${REGISTRY_DIR}/data:/var/lib/registry:z \
        -v ${REGISTRY_DIR}/auth:/auth:z \
        -e "REGISTRY_AUTH=htpasswd" \
        -e "REGISTRY_AUTH_HTPASSWD_REALM=Registry Realm" \
        -e REGISTRY_AUTH_HTPASSWD_PATH=/auth/htpasswd \
        -v ${REGISTRY_DIR}/certs:/certs:z \
        -e REGISTRY_HTTP_TLS_CERTIFICATE=/certs/${REGISTRY_CRT} \
        -e REGISTRY_HTTP_TLS_KEY=/certs/${REGISTRY_KEY} \
        ${image}

    # Add a /etc/hosts entry for $LOCAL_REGISTRY_DNS_NAME
    sudo sed -i "/${LOCAL_REGISTRY_DNS_NAME}/d" /etc/hosts
    node_ip=$(hostname -I | awk '{print $1}')
    echo "$node_ip ${LOCAL_REGISTRY_DNS_NAME}" | sudo tee -a /etc/hosts

}

pushd ${WORKING_DIR}

# Wait until binaries are available, these will eventually be in image
while [[ ! -e ${OC_MIRROR_FILE} ]] || [[ ! -e ${REGISTRY_IMAGE} ]]; do
    echo "Waiting for ${OC_MIRROR_FILE} and ${REGISTRY_IMAGE} to be available"
    sleep 5
done

# TODO - use docker-registry image from release image
# Copy image to container store
# sudo podman image trust set -f /etc/pki/rpm-gpg/RPM-GPG-KEY-redhat-release registry.access.redhat.com
# sudo podman image trust set -f /etc/pki/rpm-gpg/RPM-GPG-KEY-redhat-release registry.redhat.io

# tar xvzf docker-registry.tgz

# echo "Copying docker registry"
# docker_dir="$(ls -d ./docker-registry--**)"
# echo "sudo skopeo copy dir:./$docker_dir containers-storage:$DOCKER_REGISTRY_IMAGE"
# sudo skopeo copy dir:./$docker_dir containers-storage:$DOCKER_REGISTRY_IMAGE

# sudo podman image load -i registry_image.tar > image_load.out
# image_id=$(cat image_load.out | cut -d":" -f3)
# rm image_load.out

sudo podman image load -i ${REGISTRY_IMAGE}
image_id=$(sudo podman image list --format json | jq -c '.[] | select(.Names[] | contains("quay.io/libpod/registry")) | .Id' | tr -d '"')

# Start registry
setup_podman_mirror_registry ${image_id}

rm ${REGISTRY_IMAGE}

echo "Sleeping 10 seconds"
sleep 10

archive_file="$(ls ${IMAGE_ARCHIVE})"

# Mirror image tarfile using oc_mirror
chmod 777 ${OC_MIRROR_FILE}

${OC_MIRROR_FILE} --from $archive_file docker://${LOCAL_REGISTRY_DNS_NAME}:${LOCAL_REGISTRY_PORT} --dest-skip-tls --skip-metadata-check -v 9
echo "oc-mirror is complete"

rm ${archive_file}

popd
