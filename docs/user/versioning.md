# Versioning

The installer uses [Semantic Versioning][semver] for its user-facing API.
Covered by the versioning are:

* `openshift-install [options] create install-config`, which will always create `install-config.yaml` in the asset directory, although the version of the generated install-config may change.
* `openshift-install [options] create node-config`, which will always create `master.ign` and `worker.ign` in the asset directory, although the content of the generated files may change.
* `openshift-install [options] create bootstrap-config`, which will always create `bootstrap.ign` in the asset directory, although the content of the generated files may change.
* `openshift-install [options] create pre-cluster`, which will always create `boostrap.ign`, `master.ign`, `worker.ign`, `metadata.json`, `auth/kubeconfig`, and `auth/kubeadmin-password` in the `pre-cluster` directory in the asset directory, although the contents of the generated files may change.
* `openshift-install [options] create cluster`, which will always launch a new cluster.
* `openshift-install [options] destroy bootstrap`, which will always destroy any bootstrap resources created for the cluster.
* `openshift-install [options] destroy cluster`, which will always destroy the cluster resources.
* `openshift-install [options] help`, which will always show help for the command, although available options and unstable commands may change.
* `openshift-install [options] version`, which will always show sufficient version information for maintainers to identify the installer, although the format and content of its output may change.
* The install-config format.  New versions of this format may be released, but within a minor version series, the `openshift-install` will continue to be able to read previous versions.

The following are explicitly not covered:

* `openshift-install [options] graph`
* `openshift-install [options] create manifest-templates`
* `openshift-install [options] create manifests`

That means that the only stable install-time configuration is [via the install-config](overview.md#multiple-invocations).
If you want a reliable way to alter, add, or remove Kubernetes objects, you should perform those actions as day-2 operations.

[semver]: https://semver.org/spec/v2.0.0.html
