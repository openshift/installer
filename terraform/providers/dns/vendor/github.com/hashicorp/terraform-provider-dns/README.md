Terraform Provider
==================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
-	[Go](https://golang.org/doc/install) 1.19 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-dns`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
$ git clone git@github.com:terraform-providers/terraform-provider-dns
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-dns
$ make build
```

Using the provider
----------------------
## Fill in for each provider

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.19+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make bin
...
$ $GOPATH/bin/terraform-provider-dns
...
```

In order to run unit testing for the provider:

```sh
$ make test
```

In order to run acceptance tests, excluding ones requiring a `DNS_UPDATE_SERVER`:

```sh
$ make testacc
```

To run the full suite of acceptance tests:

```sh
$ ./internal/provider/acceptance.sh
```

Which has the following prerequisites:

- [Docker](https://www.docker.com/)
- [Go](https://golang.org/)
- [Kerberos Clients](https://web.mit.edu/kerberos/dist/) (e.g. `kinit`)
- [Make](https://www.gnu.org/software/make/)
- [Terraform CLI](https://terraform.io/)
- `/etc/hosts` entry (or equivalent): `127.0.0.1 ns.example.com`

### macOS Setup

- [Docker for Mac](https://docs.docker.com/docker-for-mac/install/)
- [Go](https://golang.org/dl/) or with Homebrew: `brew install go`
- [Terraform CLI](https://www.terraform.io/downloads.html) or with Homebrew: `brew install hashicorp/tap/terraform`

```shell
echo "127.0.0.1 ns.example.com" | sudo tee -a /etc/hosts
```

### Ubuntu Setup

- [Docker Engine](https://docs.docker.com/engine/install/ubuntu/)
- [Go](https://github.com/golang/go/wiki/Ubuntu)
- [Terraform CLI](https://www.terraform.io/docs/cli/install/apt.html)

```shell
echo "127.0.0.1 ns.example.com" | sudo tee -a /etc/hosts
sudo apt-get install krb5-user make
# If prompted for Kerberos configuration:
# Default Realm: EXAMPLE.COM
# Server: ns.example.com
# Administrative Server: ns.example.com
```
