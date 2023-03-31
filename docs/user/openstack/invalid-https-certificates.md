# OpenShift 4.10 refuses legacy HTTPS certificates

With OpenShift v4.10, HTTPS certificates not using the `Subject Alternative Names` fields will be rejected. Upgrades will be blocked if such certificates are detected in some areas; however OpenShift will not automatically check the underlying OpenStack infrastructure prior to upgrading or installing. This is what the following instructions will walk you through doing.

---

A script provided below automates the operation. However, it requires to have a set of tools available (including a relatively recent version of `python3-openstackclient`). To manually check your OpenStack infrastructure:

1. Collect the URL of the OpenStack public endpoints with `openstack catalog list` (HTTP (unsecured) endpoints do not need to be checked)
2. For each public HTTPS endpoint: collect the host (by removing the scheme, the port and the path) and the port
3. Run this openssl command to extract the SAN field of the certificate:

```plaintext
host=<the host part of the URL>
port=<the port part of the URL; 443 if unspecified>
openssl s_client -showcerts -servername "$host" -connect "$host:$port" </dev/null 2>/dev/null \
    | openssl x509 -noout -ext subjectAltName
```

If the output resembles this, the certificate is OK:
```plaintext
X509v3 Subject Alternative Name:
    DNS:yout.host.example.net
```

If instead there is no output, then the certificate is invalid and it needs to be re-issued.

---

This script automatically checks and reports on all HTTPS endpoints in an OpenStack catalog. Populate the environment with OpenStack credentials for the target cloud, then run the following Bash script.

Requirements:
* Bash v4+
* grep
* [Python OpenStack client][openstack-cli] v4+
* [jq][jq]
* [openssl 1.1.1l+][openssl]

```bash
#!/usr/bin/env bash

set -Eeuo pipefail

declare catalog san
catalog="$(mktemp)"
san="$(mktemp)"
readonly catalog san

declare invalid=0

openstack catalog list --format json --column Name --column Endpoints \
	| jq -r '.[] | .Name as $name | .Endpoints[] | select(.interface=="public") | [$name, .interface, .url] | join(" ")' \
	| sort \
	> "$catalog"

while read -r name interface url; do
	# Ignore HTTP
	if [[ ${url#"http://"} != "$url" ]]; then
		continue
	fi

	# Remove the schema from the URL
	noschema=${url#"https://"}
	
	# If the schema was not HTTPS, error
	if [[ "$noschema" == "$url" ]]; then
		echo "ERROR (unknown schema): $name $interface $url"
		exit 2
	fi

	# Remove the path and only keep host and port
	noschema="${noschema%%/*}"
	host="${noschema%%:*}"
	port="${noschema##*:}"

	# Add the port if was implicit
	if [[ "$port" == "$host" ]]; then
		port='443'
	fi

	# Get the SAN fields
	openssl s_client -showcerts -servername "$host" -connect "$host:$port" </dev/null 2>/dev/null \
		| openssl x509 -noout -ext subjectAltName \
		> "$san"

	# openssl returns the empty string if no SAN is found.
	# If a SAN is found, openssl is expected to return something like:
	#
	#    X509v3 Subject Alternative Name:
	#        DNS:standalone, DNS:osp1, IP Address:192.168.2.1, IP Address:10.254.1.2
	if [[ "$(grep -c "Subject Alternative Name" "$san" || true)" -gt 0 ]]; then
		echo "PASS: $name $interface $url"
	else
		invalid=$((invalid+1))
		echo "INVALID: $name $interface $url"
	fi
done < "$catalog"

# clean up temporary files
rm "$catalog" "$san"

if [[ $invalid -gt 0 ]]; then
	echo "${invalid} legacy certificates were detected. Update your certificates to include a SAN field."
	exit 1
else
	echo "All HTTPS certificates for this cloud are valid."
fi
```

In case the script reports INVALID certificates, they will have to be replaced with formally valid certificates (containing SAN fields for the server) before installing or upgrading to OpenShift v4.10.

## Legacy certificates

The [OpenShift v4.6 release notes][4.6-release-notes] warned about the deprecation in OpenShift of HTTPS certificates exposing the protected names in the CommonName field. With v4.10, this deprecation is enforced in the clusters. HTTPS certificates must now use the `Subject Alternative Names` fields or be rejected.

[openstack-cli]: https://docs.openstack.org/python-openstackclient/latest/ "OpenStackClient (aka OSC) is a command-line client for OpenStack that brings the command set for Compute, Identity, Image, Object Storage and Block Storage APIs together in a single shell with a uniform command structure"
[jq]: https://stedolan.github.io/jq/ "jq is a lightweight and flexible command-line JSON processor."
[openssl]: https://www.openssl.org/ "Cryptography and SSL/TLS Toolkit "
[4.6-release-notes]: https://docs.openshift.com/container-platform/4.6/release_notes/ocp-4-6-release-notes.html#ocp-4-6-tls-common-name "OpenShift Container Platform 4.6 release notes"
