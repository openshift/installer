# Building the binary for development

Just like in any other development build of the openshift-install, the process starts with running:

    $ ./hack/build.sh

This will result in the binary being written to *bin/openshift-install*

Since the agent based installer has version / release image checks, though, it is important to patch the binary with the release that you intend to use for your development deployments. If you have not built your own release payload, it is usually a good idea to pick the latest nightly. You can check what the latest nightly is by doing:

    $ RELEASE_VERSION=$(curl -s https://amd64.ocp.releases.ci.openshift.org/api/v1/releasestream/4.12.0-0.nightly/latest | jq -r '.name')
    $ RELEASE_IMAGE=$(curl -s https://amd64.ocp.releases.ci.openshift.org/api/v1/releasestream/4.12.0-0.nightly/latest | jq -r '.pullSpec')

Once we have these values, we can write them into the binary:

    $ RELEASE_VERSION_LOCATION=$(grep -oba ._RELEASE_VERSION_LOCATION_.XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX bin/openshift-install | cut -d':' -f1)
    $ RELEASE_IMAGE_LOCATION=$(grep -oba ._RELEASE_IMAGE_LOCATION_.XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX bin/openshift-install | cut -d':' -f1)
    $ printf "%s\0" "$RELEASE_VERSION" | dd of="bin/openshift-install" bs=1 seek="$RELEASE_VERSION_LOCATION" conv=notrunc
    $ printf "%s\0" "$RELEASE_IMAGE" | dd of="bin/openshift-install" bs=1 seek="$RELEASE_IMAGE_LOCATION" conv=notrunc

You can verify you got the right patching by running:

    $ ./bin/openshift-install version
    ./bin/openshift-install 4.12.0-0.nightly-2022-10-18-192348¹
    built from commit 6a5cc95fe77676c3d38d6aac5d6668c9f2d65ddb²
    release image registry.ci.openshift.org/ocp/release:4.12.0-0.nightly-2022-10-18-192348³
    release architecture amd64

Where
* <sup>1</sup> Should be the latest nightly version
* <sup>2</sup> Should be the commit id of the tip of your branch
* <sup>3</sup> Should be the latest nightly release image

Once all of this is done. You can use the *openshift-install* binary normally.
