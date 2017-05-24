# Audiences and tools

The Tectonic installer has two main use-cases and audiences – end-users that want to deploy clusters, and developers that want to extend, improve or modify the codebase. The user experience is important but different for these audiences, and is explained in detail below.

## End-user experience

The primary audience of this installer is end-users that want to deploy one or more clusters on supported platforms. The ideal UX is to download a release of the installer that requires minimal configuration of the user's machine, including dependencies.

Freeing these users from installing many dependencies can isolate them from differences between platforms (macOS or Linux). This also reduces the documentation burden.

We should strive to _never require_ end-users to use or install `make`, `npm`, etc to install a cluster.

## Developer experience

The developer workflow is reflective of how often clusters will be created and destroyed. This project makes heavy use of `make` to make these repetitive actions easier.

It is expected that developers have a working knowledge of Terraform, including how to configure/use a `.terraformrc` and things of that nature.
