#!/usr/bin/env bash
set -euo pipefail

# Generate Konnectivity certificates with a self-signed CA (1-day validity).
# These are needed for mTLS between the Konnectivity server and agents
# during the bootstrap phase.
#
# Usage: konnectivity-certs.sh <bootstrap-node-ip>

BOOTSTRAP_NODE_IP="${1:?Usage: konnectivity-certs.sh <bootstrap-node-ip>}"

KONNECTIVITY_CERT_DIR=/opt/openshift/tls/konnectivity
mkdir -p "${KONNECTIVITY_CERT_DIR}"

echo "Generating Konnectivity certificates in ${KONNECTIVITY_CERT_DIR}..."

# Generate self-signed Konnectivity CA
openssl req -x509 -newkey rsa:2048 -nodes \
    -keyout "${KONNECTIVITY_CERT_DIR}/ca.key" \
    -out "${KONNECTIVITY_CERT_DIR}/ca.crt" \
    -days 1 \
    -subj "/CN=konnectivity-signer/O=openshift"

# Server certificate for agent endpoint (needs bootstrap IP as SAN)
openssl req -new -newkey rsa:2048 -nodes \
    -keyout "${KONNECTIVITY_CERT_DIR}/server.key" \
    -out "${KONNECTIVITY_CERT_DIR}/server.csr" \
    -subj "/CN=konnectivity-server/O=openshift"

openssl x509 -req -in "${KONNECTIVITY_CERT_DIR}/server.csr" \
    -CA "${KONNECTIVITY_CERT_DIR}/ca.crt" \
    -CAkey "${KONNECTIVITY_CERT_DIR}/ca.key" \
    -CAcreateserial \
    -out "${KONNECTIVITY_CERT_DIR}/server.crt" \
    -days 1 \
    -extfile <(printf "extendedKeyUsage=serverAuth\nsubjectAltName=IP:%s" "${BOOTSTRAP_NODE_IP}")

# Agent client certificate (shared by all agents)
openssl req -new -newkey rsa:2048 -nodes \
    -keyout "${KONNECTIVITY_CERT_DIR}/agent.key" \
    -out "${KONNECTIVITY_CERT_DIR}/agent.csr" \
    -subj "/CN=konnectivity-agent/O=openshift"

openssl x509 -req -in "${KONNECTIVITY_CERT_DIR}/agent.csr" \
    -CA "${KONNECTIVITY_CERT_DIR}/ca.crt" \
    -CAkey "${KONNECTIVITY_CERT_DIR}/ca.key" \
    -CAcreateserial \
    -out "${KONNECTIVITY_CERT_DIR}/agent.crt" \
    -days 1 \
    -extfile <(printf "extendedKeyUsage=clientAuth")

# Clean up CSR files
rm -f "${KONNECTIVITY_CERT_DIR}"/*.csr

echo "Konnectivity certificates generated successfully."
