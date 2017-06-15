# Running Tectonic Installer in a Docker container on Windows

Tectonic v1.6.4 does not include an Installer binary for Windows. Windows users can install Tectonic by running an Installer container with Docker Community Edition (CE). This document describes how to use Docker CE on Windows to run Tectonic Installer to create Tectonic clusters on supported cloud providers or physical hardware.

This document does not describe deploying Tectonic clusters on Windows hosts.

## Installing Tectonic from a Windows machine

### 1.  Install Docker Community Edition for Windows

Install the [Docker Community Edition for Windows][dce-win]. Docker for Windows transparently uses a Hyper-V Linux virtual machine to run the Docker engine. Windows 10 Professional is required.

### 2. Run the Tectonic Installer Container

Type `cmd` in the Windows menu search box and press the Enter key. A new Windows command prompt opens.

Issue the following command to fetch and run the Tectonic Installer container image from the Quay registry:

```sh
docker run --rm -p 4444:4444 -it quay.io/coreos/tectonic-installer:1.6.4-tectonic.1 /go/src/github.com/coreos/tectonic-installer/installer/bin/linux/installer -open-browser=false -address 0.0.0.0:4444
```

The status of the image download in progressively printed in the command window. Once Tectonic Installer is downloaded and running, `Starting Tectonic Installer on 0.0.0.0:4444` will be reported.

### 3. Install Tectonic

Direct a web browser on the Windows host to the URL `http://127.0.0.1:4444`. The Tectonic Installer GUI is shown. Proceed through the installation according to the instructions for the target provider, such as [AWS][aws-install] or [Bare Metal][bm-install].

### 4.  Download cluster assets

Once the cluster is installed, Installer will display a success message and a large green button labeled *Download assets*. Click this button to download the cluster's assets, which include configuration and authentication files, used for subsequent cluster infrastructure administration – for example, to destroy clusters to free their resources.

### 5. Use the new cluster

At the last Installer step, follow the link to Tectonic Console to use the new cluster. Check out the [Tectonic Tutorials to learn how to run, replicate, and scale a simple application][tut-firstapp] on the new cluster, and advance from there to more complex applications and concepts.

## Using Windows to delete a cluster

Deleting a cluster is referred to as "destroying" it, from the `destroy` subcommand. To destroy a cluster when using a Windows/Docker CE host, the Tectonic Installer container will again be used, but directly on the command line with Terraform commands, rather than a browser-based GUI. Most of the commands below are invoked within the executing container.

Ensure the cluster's [`assets.zip` file downloaded from Tectonic Installer][assets-anchor] is present in the Windows host's *Downloads* directory, from where the `docker run` command below expects to connect it to the Installer container.

### 1. Open Windows command prompt

Type `cmd` into the Windows menu search box and press the Enter key. A new command prompt window opens.

### 2. Run the Installer container

At the Windows command prompt, issue the command:

  ```sh
  docker run --rm -it -v %USERPROFILE%/Downloads:/Downloads quay.io/coreos/tectonic-installer:1.6.4-tectonic.1 bash
  ```

  The Installer container executes and presents a `bash` shell prompt within the container.

### 3. Export AWS credentials:

Replace the values in angle brackets `<...>` with the appropriate AWS credentials.

```
$ export AWS_ACCESS_KEY_ID=<ACCESSKEYID>
$ export AWS_SECRET_ACCESS_KEY=<SECRETACCESSKEY>`
```

### 4. Extract the assets bundle

The environment variable `PROJECT_DIR` is already set in the container shell for convenience.

```sh
$ unzip -o /Downloads/assets.zip -d $PROJECT_DIR/installer/bin/linux/clusters/
```

### 5. Destroy the cluster

```sh
$ cd $PROJECT_DIR/installer/bin/linux/clusters/*
$ TERRAFORM_CONFIG=$(pwd)/.terraformrc terraform get
$ TERRAFORM_CONFIG=$(pwd)/.terraformrc terraform destroy --force
```


[assets-anchor]: #download-cluster-assets
[aws-install]: aws/index.md#step-2-install
[bm-install]: bare-metal/index.md
[dce-win]: https://store.docker.com/editions/community/docker-ce-desktop-windows
[tut-firstapp]: ../tutorials/first-app.md
