#!/bin/bash

set -eu
set -x

cleanup_docker() {
	docker stop ns
	docker stop kdc || :
}
failed() {
	cleanup_docker
	exit 1
}

command -v docker >/dev/null 2>&1 || { echo >&2 "docker command not installed or in PATH"; exit 1; }
command -v go >/dev/null 2>&1 || { echo >&2 "go command not installed or in PATH"; exit 1; }
command -v kinit >/dev/null 2>&1 || { echo >&2 "kinit command not installed or in PATH"; exit 1; }
command -v make >/dev/null 2>&1 || { echo >&2 "make command not installed or in PATH"; exit 1; }
command -v terraform >/dev/null 2>&1 || test -n "${TF_ACC_TERRAFORM_PATH:-}" || { echo >&2 "terraform command not installed or in PATH, TF_ACC_TERRAFORM_PATH not set"; exit 1; }
grep -q "ns.example.com" /etc/hosts || echo >&2 "127.0.0.1 ns.example.com not found in /etc/hosts, ensure this mapping is handled in DNS resolution configuration"

docker buildx build --target kdc --tag kdc internal/provider/testdata/
docker buildx build --target ns --tag ns internal/provider/testdata/
docker buildx build --target keytab --output type=local,dest=internal/provider/testdata/ internal/provider/testdata/

export DNS_UPDATE_SERVER=127.0.0.1
export DNS_UPDATE_PORT=15353

# Run with no authentication

docker run --privileged --cgroupns=host -d --tmpfs /tmp --tmpfs /run \
	-v /sys/fs/cgroup:/sys/fs/cgroup:rw \
	-v /etc/localtime:/etc/localtime:ro \
	-v $PWD/internal/provider/testdata/named.conf.none:/etc/named.conf:ro \
	-p 127.0.0.1:15353:53 \
	-p 127.0.0.1:15353:53/udp \
	--rm --name ns --hostname ns.example.com ns || failed
GO111MODULE=on make testacc TEST=./internal/provider || failed
cleanup_docker

# Run with TSIG authentication (MD5)

docker run --privileged --cgroupns=host -d --tmpfs /tmp --tmpfs /run \
	-v /sys/fs/cgroup:/sys/fs/cgroup:rw \
	-v /etc/localtime:/etc/localtime:ro \
	-v $PWD/internal/provider/testdata/named.conf.md5:/etc/named.conf:ro \
	-p 127.0.0.1:15353:53 \
	-p 127.0.0.1:15353:53/udp \
	--rm --name ns --hostname ns.example.com ns || failed
DNS_UPDATE_KEYNAME="tsig.example.com." DNS_UPDATE_KEYALGORITHM="hmac-md5" DNS_UPDATE_KEYSECRET="mX9XKfw/RXBj5ZnZKMy4Nw==" GO111MODULE=on make testacc TEST=./internal/provider || failed
cleanup_docker

# Run with TSIG authentication (SHA256)

docker run --privileged --cgroupns=host -d --tmpfs /tmp --tmpfs /run \
	-v /sys/fs/cgroup:/sys/fs/cgroup:rw \
	-v /etc/localtime:/etc/localtime:ro \
	-v $PWD/internal/provider/testdata/named.conf.sha256:/etc/named.conf:ro \
	-p 127.0.0.1:15353:53 \
	-p 127.0.0.1:15353:53/udp \
	--rm --name ns --hostname ns.example.com ns || failed
DNS_UPDATE_KEYNAME="tsig.example.com." DNS_UPDATE_KEYALGORITHM="hmac-sha256" DNS_UPDATE_KEYSECRET="UHeh4Iv/DVmPhi6LqCPDs6PixnyjLH4fjGESBjYnOyE=" GO111MODULE=on make testacc TEST=./internal/provider || failed
cleanup_docker

export KRB5_CONFIG="${PWD}/internal/provider/testdata/krb5.conf"
export DNS_UPDATE_REALM="EXAMPLE.COM"
export DNS_UPDATE_SERVER="ns.example.com"

# Run with Kerberos authentication (password authentication)

docker run --privileged --cgroupns=host -d --tmpfs /tmp --tmpfs /run \
	-v /sys/fs/cgroup:/sys/fs/cgroup:rw \
	-v /etc/localtime:/etc/localtime:ro \
	-p 127.0.0.1:18888:88 \
	-p 127.0.0.1:18888:88/udp \
	-p 127.0.0.1:464:464 \
	-p 127.0.0.1:464:464/udp \
	--rm --name kdc kdc || failed
docker run --privileged --cgroupns=host -d --tmpfs /tmp --tmpfs /run \
	-v /sys/fs/cgroup:/sys/fs/cgroup:rw \
	-v /etc/localtime:/etc/localtime:ro \
	-v $PWD/internal/provider/testdata/named.conf.kerberos:/etc/named.conf:ro \
	-p 127.0.0.1:15353:53 \
	-p 127.0.0.1:15353:53/udp \
	--rm --name ns --hostname ns.example.com ns || failed
DNS_UPDATE_USERNAME="test" DNS_UPDATE_PASSWORD="password" GO111MODULE=on make testacc TEST=./internal/provider || failed
cleanup_docker

# Run with Kerberos authentication (keytab authentication)

docker run --privileged --cgroupns=host -d --tmpfs /tmp --tmpfs /run \
	-v /sys/fs/cgroup:/sys/fs/cgroup:rw \
	-v /etc/localtime:/etc/localtime:ro \
	-p 127.0.0.1:18888:88 \
	-p 127.0.0.1:18888:88/udp \
	-p 127.0.0.1:464:464 \
	-p 127.0.0.1:464:464/udp \
	--rm --name kdc kdc || failed
docker run --privileged --cgroupns=host -d --tmpfs /tmp --tmpfs /run \
	-v /sys/fs/cgroup:/sys/fs/cgroup:rw \
	-v /etc/localtime:/etc/localtime:ro \
	-v $PWD/internal/provider/testdata/named.conf.kerberos:/etc/named.conf:ro \
	-p 127.0.0.1:15353:53 \
	-p 127.0.0.1:15353:53/udp \
	--rm --name ns --hostname ns.example.com ns || failed
DNS_UPDATE_USERNAME="test" DNS_UPDATE_KEYTAB="${PWD}/internal/provider/testdata/test.keytab" GO111MODULE=on make testacc TEST=./internal/provider || failed
cleanup_docker

# Run with Kerberos authentication (session authentication)

docker run --privileged --cgroupns=host -d --tmpfs /tmp --tmpfs /run \
	-v /sys/fs/cgroup:/sys/fs/cgroup:rw \
	-v /etc/localtime:/etc/localtime:ro \
	-p 127.0.0.1:18888:88 \
	-p 127.0.0.1:18888:88/udp \
	-p 127.0.0.1:464:464 \
	-p 127.0.0.1:464:464/udp \
	--rm --name kdc kdc || failed
docker run --privileged --cgroupns=host -d --tmpfs /tmp --tmpfs /run \
	-v /sys/fs/cgroup:/sys/fs/cgroup:rw \
	-v /etc/localtime:/etc/localtime:ro \
	-v $PWD/internal/provider/testdata/named.conf.kerberos:/etc/named.conf:ro \
	-p 127.0.0.1:15353:53 \
	-p 127.0.0.1:15353:53/udp \
	--rm --name ns --hostname ns.example.com ns || failed
echo "password" | kinit --password-file=STDIN test@EXAMPLE.COM || echo "password" | kinit test@EXAMPLE.COM
GO111MODULE=on make testacc TEST=./internal/provider || failed
cleanup_docker
