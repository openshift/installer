#!/bin/bash

set -eu
set -x

# Test domains
export DNS_DOMAIN_FORWARD="example.com."
export DNS_DOMAIN_REVERSE="1.168.192.in-addr.arpa."

DOCKER_CONTAINER_NAME=tf_acc_dns

cleanup_docker() {
	docker stop "$DOCKER_CONTAINER_NAME"
	docker rm "$DOCKER_CONTAINER_NAME"
}
failed() {
	cleanup_docker
	exit 1
}

# Run with no authentication

export DNS_UPDATE_SERVER=127.0.0.1
export DNS_UPDATE_PORT=55354
docker run -d -p "$DNS_UPDATE_PORT:53/udp" \
	-e BIND_DOMAIN_FORWARD=${DNS_DOMAIN_FORWARD} \
	-e BIND_DOMAIN_REVERSE=${DNS_DOMAIN_REVERSE} \
	-e BIND_INSECURE=true \
	--name "$DOCKER_CONTAINER_NAME" drebes/bind || failed
GO111MODULE=on GOFLAGS=-mod=vendor make testacc TEST=./dns || failed
cleanup_docker

# Run with authentication

export DNS_UPDATE_KEYNAME=${DNS_DOMAIN_FORWARD}
export DNS_UPDATE_KEYALGORITHM="hmac-md5"
export DNS_UPDATE_KEYSECRET="c3VwZXJzZWNyZXQ="
docker run -d -p "$DNS_UPDATE_PORT:53/udp" \
	-e BIND_DOMAIN_FORWARD=${DNS_DOMAIN_FORWARD} \
	-e BIND_DOMAIN_REVERSE=${DNS_DOMAIN_REVERSE} \
	-e BIND_KEY_NAME=${DNS_UPDATE_KEYNAME} \
	-e BIND_KEY_ALGORITHM=${DNS_UPDATE_KEYALGORITHM} \
	-e BIND_KEY_SECRET=${DNS_UPDATE_KEYSECRET} \
	--name "$DOCKER_CONTAINER_NAME" drebes/bind || failed
GO111MODULE=on GOFLAGS=-mod=vendor make testacc TEST=./dns || failed
cleanup_docker
