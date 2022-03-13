# IBM Cloud SDK for Go

[![Build Status](https://travis-ci.org/IBM-Cloud/bluemix-go.svg?branch=master)](https://travis-ci.org/IBM-Cloud/bluemix-go) [![GoDoc](https://godoc.org/github.com/IBM-Cloud/bluemix-go?status.svg)](https://godoc.org/github.com/IBM-Cloud/bluemix-go)

bluemix-go provides the Go implementation for operating the IBM Cloud platform, which is based on the [Cloud Foundry API][cloudfoundry_api].

## Installing

1. Install the SDK using the following command

```bash
go get github.com/IBM-Cloud/bluemix-go
```

2. Update the SDK to the latest version using the following command

```bash
go get -u github.com/IBM-Cloud/bluemix-go
```


## Using the SDK

You must have a working IBM Cloud account to use the APIs. [Sign up][ibmcloud_signup] if you don't have one.

The SDK has ```examples``` folder which cites few examples on how to use the SDK.
First you need to create a session.

```go
import "github.com/IBM-Cloud/bluemix-go/session"

func main(){

    s := session.New()
    .....
}
```

Creating session in this way creates a default configuration which reads the value from the environment variables.
You must export the following environment variables.
* IBMID - This is the IBM ID
* IBMID_PASSWORD - This is the password for the above ID

OR

* IC_API_KEY/IBMCLOUD_API_KEY - This is the Bluemix API Key. Login to [IBMCloud][ibmcloud_login] to create one if you don't already have one. See instructions below for creating an API Key.

The default region is _us_south_. You can override it in the [Config struct][ibmcloud_go_config]. You can also provide the value via environment variables; either via _IC_REGION_ or _IBMCLOUD_REGION_. Valid regions are -
* us-south
* us-east
* eu-gb
* eu-de
* au-syd
* jp-tok

The maximum retries is 3. You can override it in the [Config struct][ibmcloud_go_config]. You can also provide the value via environment variable; via MAX_RETRIES

## Creating an IBM Cloud API Key

First, navigate to the IBM Cloud console and use the Manage toolbar to access IAM.

![Access IAM from the Manage toolbar](.screenshots/screenshot_api_keys_iam.png)

On the left, click "IBM Cloud API Keys"

![Click IBM Cloud API Keys](.screenshots/screenshot_api_keys_iam_left.png)

Press "Create API Key"

![Press Create API Key](.screenshots/screenshot_api_keys_create_button.png)

Pick a name and description for your key

![Set name and description](.screenshots/screenshot_api_keys_create.png)

You have created a key! Press the eyeball to show the key. Copy or save it because keys can't be displayed or downloaded twice.

![Your key is now created](.screenshots/screenshot_api_keys_create_successful.png)


[ibmcloud_signup]: https://console.ng.bluemix.net/registration/?target=%2Fdashboard%2Fapps
[ibmcloud_login]: https://console.ng.bluemix.net
[ibmcloud_go_config]: https://godoc.org/github.com/IBM-Cloud/bluemix-go#Config
[cloudfoundry_api]: https://apidocs.cloudfoundry.org/264/
