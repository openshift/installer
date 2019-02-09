#!/bin/sh
#
# Given a release image, extract the referenced installer and use it
# to launch a cluster based on that same release image.
#
# Usage:
#
#   install-release-image.sh [ARG...]
#
# Requires:
#
# * oc, to extract the installer image
# * podman, to run the installer image
# * realpath, which is not in POSIX as a shell command [1], but is in
#   GNU coreutils [2].  The backing C function is in POSIX [3].
# * A pull secret in ~/.docker/config.json for both oc and podman.
#
# Optional:
#
# * ASSETS
#   The path is mounted into the installer container and used as the
#   installer's asset directory.  Defaults to the current working
#   directory.
# * RELEASE
#   The pull-spec for the release image.  Defaults to the development
#   OKD tip.
# * SSH_PUBKEY
#   If set, this is mounted into the installer container at
#   ~/.ssh/key.pub.  Otherwise all ~/.ssh/*.pub are mounted in.
# * AWS_PROFILE [4]
#   If set, this is also set in the installer container.
# * AWS_CONFIG_FILE [4]
#   This defaults to ~/.aws/config, and, if set to an existing file,
#   it is mounted into the installer container at ~/.aws/config
# * AWS_SHARED_CREDENTIALS_FILE [4]
#   If set, this is mounted into the installer container at
#   ~/.aws/credentials
#
# [1]: http://pubs.opengroup.org/onlinepubs/9699919799/idx/utilities.html
# [2]: http://man7.org/linux/man-pages/man1/realpath.1.html
# [3]: http://pubs.opengroup.org/onlinepubs/9699919799/functions/realpath.html
# [4]: https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html

die () {
	# shellcheck disable=SC2059
	printf "${@}" >&2
	exit 1
}

mkdir -p "${ASSETS:-.}" || die 'failed to create asset directory\n'
ASSETS="$(realpath -e "${ASSETS:-.}")" || die 'failed to resolve asset path\n'
RELEASE="${RELEASE:-registry.svc.ci.openshift.org/openshift/origin-release:v4.0}"
INSTALLER="$(oc adm release info --image-for=installer "${RELEASE}")" || die 'failed to resolve installer image\n'
set -- --env HOME=/home "${INSTALLER}" "${@}" || die 'failed to insert installer into podman arguments\n'
if test -n "${SSH_PUBKEY}"
then
	SSH_PUBKEY="$(realpath -e "${SSH_PUBKEY}")" || die 'failed to resolve SSH_PUBKEY\n'
	set -- --volume "${SSH_PUBKEY}:/home/.ssh/id_key.pub:ro,z" "${@}" || die 'failed to insert SSH_PUBKEY into podman arguments\n'
else
	for SSH_PUBKEY in ~/.ssh/*.pub
	do
		SSH_PUBKEY_FILENAME="$(basename "${SSH_PUBKEY}")" || die 'failed to split filename from %s\n' "${SSH_PUBKEY}"
		set -- --volume "${SSH_PUBKEY}:/home/.ssh/${SSH_PUBKEY_FILENAME}:ro,z" "${@}" || die 'failed to insert %s into podman arguments\n' "${SSH_PUBKEY}"
	done
fi

if test -n "${AWS_PROFILE}"
then
	set -- --env AWS_PROFILE="${AWS_PROFILE}" "${@}" || die 'failed to insert AWS_PROFILE into podman arguments\n'
fi

if test -n "${AWS_CONFIG_FILE}"
then
	AWS_CONFIG_FILE="$(realpath -e "${AWS_CONFIG_FILE}")" || die 'failed to resolve AWS_CONFIG_FILE\n'
else
	AWS_CONFIG_FILE=~/.aws/config
fi
if test -e "${AWS_CONFIG_FILE}"
then
	set -- --volume "${AWS_CONFIG_FILE}:/home/.aws/config:ro,z" "${@}" || die 'failed to insert AWS_CONFIG_FILE into podman arguments\n'
fi

if test -n "${AWS_SHARED_CREDENTIALS_FILE}"
then
	AWS_SHARED_CREDENTIALS_FILE="$(realpath -e "${AWS_SHARED_CREDENTIALS_FILE}")" || die 'failed to resolve AWS_SHARED_CREDENTIALS_FILE\n'
	set -- --volume "${AWS_SHARED_CREDENTIALS_FILE}:/home/.aws/credentials:z" "${@}" || die 'failed to insert AWS_SHARED_CREDENTIALS_FILE into podman arguments\n'
fi

exec podman run --rm -it \
	--user "$(id -u):$(id -g)" \
	--env OPENSHIFT_INSTALL_RELEASE_IMAGE_OVERRIDE="${RELEASE}" \
	--volume "${ASSETS}:/output:z" \
	"${@}"
