Terraform Provider For Alibaba Cloud
==================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="400px"> 


<img src="https://www.datocms-assets.com/2885/1506527326-color.svg" width="400px">


Supported Versions
------------------

| Terraform version | minimum provider version |maxmimum provider version
| ---- | ---- | ----| 
| 0.11.x	| 1.0.0	| 1.122.0 |
| 0.12.x	| 1.0.0	| 1.122.0 |
| 0.13.x	| 1.0.0	| 1.122.0 |
| 0.14.x	| 1.0.0	| 1.122.0 |

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
-	[Go](https://golang.org/doc/install) 1.13 (to build the provider plugin)
-   [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports):
    ```
    go get golang.org/x/tools/cmd/goimports
    ```

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/aliyun/terraform-provider-alicloud`

```sh
$ mkdir -p $GOPATH/src/github.com/aliyun; cd $GOPATH/src/github.com/aliyun
$ git clone git@github.com:aliyun/terraform-provider-alicloud
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/aliyun/terraform-provider-alicloud
$ make build
```

Using the provider
----------------------
## Fill in for each provider

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.11+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-alicloud
...
```

Running `make dev` or `make devlinux` or `devwin` will only build the specified developing provider which matchs the local system.
And then, it will unarchive the provider binary and then replace the local provider plugin.

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

## Acceptance Testing
Before making a release, the resources and data sources are tested automatically with acceptance tests (the tests are located in the alicloud/*_test.go files).
You can run them by entering the following instructions in a terminal:
```
cd $GOPATH/src/github.com/aliyun/terraform-provider-alicloud
export ALICLOUD_ACCESS_KEY=xxx
export ALICLOUD_SECRET_KEY=xxx
export ALICLOUD_REGION=xxx
export ALICLOUD_ACCOUNT_ID=xxx
export ALICLOUD_RESOURCE_GROUP_ID=xxx
export outfile=gotest.out
TF_ACC=1 TF_LOG=INFO go test ./alicloud -v -run=TestAccAlicloud -timeout=1440m | tee $outfile
go2xunit -input $outfile -output $GOPATH/tests.xml
```

-> **Note:** The last line is optional, it allows to convert test results into a XML format compatible with xUnit.


-> **Note:** Most test cases will create PostPaid resources when running above test command. However, currently not all
 account site type support create PostPaid resources, so you need set your account site type before running the command:
```
# If your account belongs to domestic site
export ALICLOUD_ACCOUNT_SITE=Domestic

# If your account belongs to international site
export ALICLOUD_ACCOUNT_SITE=International
```
The setting of acount site type can skip some unsupported cases automatically.

-> **Note:** At present, there is missing CMS contact group resource and please create manually a contact group by web console and set it by environment variable `ALICLOUD_CMS_CONTACT_GROUP`, like:
 ```
 export ALICLOUD_CMS_CONTACT_GROUP=tf-testAccCms
 ```
 Otherwise, all of resource `alicloud_cms_alarm's` test cases will be skipped.

## Refer

Alibaba Cloud Provider [Official Docs](https://www.terraform.io/docs/providers/alicloud/index.html)
Alibaba Cloud Provider Modules [Official Modules](https://registry.terraform.io/browse?provider=alicloud)
