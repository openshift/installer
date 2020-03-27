# RHCOS image customization

By default the installer creates an image in Glance called `<clusterID>-rhcos`, and uploads the binary data from a predetermined in the installer location. The Glance image exists throughout the life of the cluster and is removed along with it.

To change this behavior and upload binary data from a custom location the user may set `OPENSHIFT_INSTALL_OS_IMAGE_OVERRIDE` environment variable that points to that location, and then start the installation. In all other respects the process will be consistent with the default.

Example:

```sh
export OPENSHIFT_INSTALL_OS_IMAGE_OVERRIDE="https://example.com/my-rhcos.qcow2"
./openshift-install create cluster --dir ostest
```

**NOTE:** For this to work, the environment variable value must be a valid http(s) URL.

If the user wants to reuse an existing Glance image without any uploading of binary data, then it is possible to set `OPENSHIFT_INSTALL_OS_IMAGE_OVERRIDE` environment variable that specifies the Glance image name. In this case no new Glance images will be created, and the image will stay when the cluster is destroyed.

Example:

```sh
export OPENSHIFT_INSTALL_OS_IMAGE_OVERRIDE="my-rhcos"
./openshift-install create cluster --dir ostest
```

**NOTE:** The only difference in behavior with the previous example is that the value here is not an http(s) URL.
