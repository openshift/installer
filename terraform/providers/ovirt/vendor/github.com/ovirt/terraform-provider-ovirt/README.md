Terraform oVirt Provider plugin
===============================

[![Build Status](https://travis-ci.org/oVirt/terraform-provider-ovirt.svg?branch=master)](https://travis-ci.org/oVirt/terraform-provider-ovirt)
[![Go Report Card](https://goreportcard.com/badge/github.com/ovirt/terraform-provider-ovirt)](https://goreportcard.com/report/github.com/ovirt/terraform-provider-ovirt)


This plugin allows Terraform to work with the oVirt Virtual Machine management platform.
It requires oVirt 4.x. 


Statements
-----------

Firstly, this project is inspired by [EMSL-MSC](http://github.com/EMSL-MSC/terraform-provider-ovirt), the author [@Maigard](https://github.com/EMSL-MSC/terraform-provider-ovirt/commits?author=Maigard) surely done a outstanding work and great thanks to him.

While in the last five months, the upstream project was not actively maintained and the pull request I committed is still not reviewed. Since this project is a heavy work in progress, for intuitive and convenient usage, I replaced the references of `EMSL-MSC` with `imjoey` in `main.go`, `README` and some other CI configuration files.

If possible, I would surely be happy to contribute back to the upstream again. ^_^ .


Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.11.x
-	[Go](https://golang.org/doc/install) 1.12 (to build the provider plugin)


Developing The Provider
-----------------------

If you wish to work on the provider, you'll first need [Go](https://golang.org) installed on your machine (version 1.12+ is *required*). You'll also need to correctly setup a [GOPATH](https://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

*Note:* This project uses [Go Modules](https://blog.golang.org/using-go-modules) making it safe to work with it outside of your existing [GOPATH](http://golang.org/doc/code.html#GOPATH). The instructions that follow assume a directory in your home directory outside of the standard GOPATH (i.e `$HOME/development/terraform-providers/`).

Clone repository to: `$HOME/development/terraform-providers/`

```sh
$ mkdir -p $HOME/development/terraform-providers/
$ cd $HOME/development/terraform-providers/
$ git clone git@github.com:ovirt/terraform-provider-ovirt
...
```

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-ovirt
...
```


Using the provider
------------------
If you're building the provider, follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) After placing it into your plugins directory,  run `terraform init` to initialize it.

Provider Usage
--------------

* Provider Configuration
```HCL
provider "ovirt" {
  username = "username@profile"
  url      = "https://ovirt/ovirt-engine/api"
  password = "Password"
  cafile   = "/path/to/ovirt/engine/ca.pem"
}
```
  * username - (Required) The username to access the oVirt api including the profile used
  * url - (Required) The url to the api endpoint (usually the ovirt server with a path of /ovirt-engine/api)
  * password - (Required) Password to access the server
  * cafile - (Optional) Path to the oVirt engine CA certificate for TLS verification
  * ca_bundle - (Optional) The oVirt engine CA certificate for TLS verification as a string input
  * insecure - (Optional) Disables TLS certificate verification
* Resources
  * ovirt_cluster
  * ovirt_datacenter
  * ovirt_disk
  * ovirt_disk_attachment
  * ovirt_host
  * ovirt_mac_pool
  * ovirt_network
  * ovirt_snapshot
  * ovirt_storage_domain
  * ovirt_tag
  * ovirt_user
  * ovirt_vm
  * ovirt_vnic
  * ovirt_vnic_profile
* Data Sources
  * ovirt_authzs
  * ovirt_clusters
  * ovirt_datacenters
  * ovirt_disks
  * ovirt_hosts
  * ovirt_mac_pools
  * ovirt_networks
  * ovirt_nics
  * ovirt_storagedomains
  * ovirt_template
  * ovirt_users
  * ovirt_vms
  * ovirt_vnic_profiles

Provider Documents
--------------
Currently the documents for this provider is not hosted by the official site [Terraform Providers](https://www.terraform.io/docs/providers/index.html). Please enter the provider directory and build the website locally.

```sh
$ make website
```

The commands above will start a docker-based web server powered by [Middleman](https://middlemanapp.com/), which hosts the documents in `website` directory. Simply open `http://localhost:4567/docs/providers/ovirt` and enjoy them.
