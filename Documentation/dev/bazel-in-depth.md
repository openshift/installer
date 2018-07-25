# Bazel Under The Hood

The goal of this document is to detail the steps taken by Bazel when building the OpenShift Installer project so that users better understand the process. Ultimately, a user building the project could elect to build the project without Bazel, either by hand or otherwise. *Note*: building without Bazel is not recommended because it will lead to non-hermetic builds, which could lead to unpredictable results. We strongly recommend using the build instructions outlined in the README for consistent, reproducible builds.

This document covers the process of building the `tarball` Bazel target, which is the main target of the project.

## Build Layout
As noted in [build.md](build.md), the goal of the build process is to produce an archive with the following file structure:

```
openshift-installer
├── config.tf
├── examples
├── modules
├── steps
└── installer
    ├── darwin
    │   ├── openshift-install
    │   └── terraform
    └── linux
        ├── openshift-install
        └── terraform
```

## Steps
### Directories
Prepare the necessary output directories:

* `openshift-installer`
* `openshift-installer/examples`
* `openshift-installer/modules`
* `openshift-installer/steps`
* `openshift-installer/installer`
* `openshift-installer/installer/darwin`
* `openshift-installer/installer/linux`

### Go Build
Build the OpenShift Installer CLI Golang binary located in `openshift-installer/installer/cmd/openshift-install` using `go build …`
The binary should be built for both Darwin and Linux and placed in the corresponding output directory, i.e. `openshift-installer/installer/darwin`, or `openshift-installer/installer/linux`.

### Terraform Binaries
Download binaries for Terraform for both Darwin and Linux and place them in the corresponding output directories.

### Terraform Configuration
Copy all required Terraform configuration files from their source directories and place them in the correct output directory. Specifically, `config.tf`, `modules` and `steps` should be copied to the output directory at `openshift-installer/config.tf`, `openshift-installer/modules`, and `openshift-installer/steps`, respectively.

### Configuration Examples
Copy the OpenShift Installer configuration examples from `examples` to the output directory at `openshift-installer/examples`.

### Archive
Lastly, archive and gzip the output directory using the `tar` utility to produce the final asset.
