# Populating the RHCOS Stream Marketplace Extension

This tool is used to reach out to cloud sdks and populate the
data/data/coreos/marketplace-rhcos.json file. That file represents
the mareketplace extension to the rhcos stream and is merged into
the stream.

To run the program:

```shell
go run -mod=vendor ./hack/rhcos/populate-marketplace-imagestream.go
```

The program will find marketplace images based on the version of the
RHCOS stream. The version can be overriden with the
`STREAM_RELEASE_OVERRIDE` variable. This is useful if you are working
on the main branch, where up-to-date images are not available. For
example, looking up Azure images requires knowing the X.Y version
to populate the offer, so when working on the main branch (4.20),
it is necessary to run the following command to correctly populate
the stream:

```shell
STREAM_RELEASE_OVERRIDE=4.19 go run -mod=vendor ./hack/rhcos/populate-marketplace-imagestream.go
```