# Bazel Under The Hood

The goal of this document is to detail the steps taken by Bazel when building the Tectonic Installer project so that users better understand the process. Ultimately, a user building the project could elect to build the project without Bazel, either by hand or otherwise. *Note*: building without Bazel is not recommended because it will lead to non-hermetic builds, which could lead to unpredictable results. We strongly recommend using the build instructions outlined in the README for consistent, reproducible builds.

This document covers the process of building the `tarball` Bazel target, which is the main target of the project.

## Build Layout
As noted in [build.md](build.md), the goal of the build process is to produce an archive with the following file structure:

```
tectonic
├── config.tf
├── examples
├── modules
├── steps
└── tectonic-installer
    ├── darwin
    │   ├── tectonic
    │   └── terraform
    └── linux
        ├── tectonic
        └── terraform
```

## Steps
### Directories
Prepare the necessary output directories:

* `tectonic`
* `tectonic/examples`
* `tectonic/modules`
* `tectonic/steps`
* `tectonic/tectonic-installer`
* `tectonic/tectonic-installer/darwin`
* `tectonic/tectonic-installer/linux`

### Go Build
Build the Tectonic CLI Golang binary located in `tectonic/installer/cmd/tectonic` using `go build …`
The binary should be built for both Darwin and Linux and placed in the corresponding output directory, i.e. `tectonic/tectonic-installer/darwin`, or `tectonic/tectonic-installer/linux`.

### Terraform Binaries
Download binaries for Terraform for both Darwin and Linux and place them in the corresponding output directories.

### Terraform Configuration
Copy all required Terraform configuration files from their source directories and place them in the correct output directory. Specifically, `config.tf`, `modules` and `steps` should be copied to the output directory at `tectonic/config.tf`, `tectonic/modules`, and `tectonic/steps`, respectively.

### Configuration Examples
Copy the Tectonic Installer configuration examples from `examples` to the output directory at `tectonic/examples`.

### Archive
Lastly, archive and gzip the output directory using the `tar` utility to produce the final asset.
