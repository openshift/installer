# Fedora libvirt environment
This Libvirt container is designed to be a general purpose and disposable environment where users are freed from the process of setting up a Libvirt environment to install an `Openshift` cluster.

## Requirements
* `just` (https://github.com/casey/just)
* `jq` (https://stedolan.github.io/jq)

## Download images
`just dl-rhcos`

## Build & Run
* `just run [repo_owner] [branch]`        - Creates a background running container after the image is built. It must be stopped using `just stop`.
* `just run-attach [repo_owner] [branch]` - Creates a foreground running container. It will be terminated if the terminal is closed or the process stopped in any way.

`repo_owner` is an optional parameter which specifies the owner of the git repository to clone for the installer.
`branch` is an optional parameter which specifies the git branch to clone for the installer.

## Stop
`just stop`

## Enter the container from a different terminal
`just exec` [args]                        - Defaults to `/bin/bash`

`args` is an optional parameter which contains instructions to execute within the container.

## Execute a command as soon as the cluster is available
`just exec-ready` [max-wait] [args]        - Executes the specified command once the cluster is available. Defaults to the cluster's `bootstrap` journal follow.

`max-wait` is an optional parameter which specifies the maximum amount of seconds to wait for the cluster to be available. Defaults to 300 seconds.
`args` is an optional parameter which contains instructions to execute within the container. 

## Status
  1. The cluster completes but takes *a long time* on my machine, so the smoke tests time out first.
  2. Wait a long time (30 minutes or so) and execute the smoke tests manually.
