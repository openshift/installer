# go-ovirt-client: an easy-to-use overlay for the oVirt Go SDK

This library provides an easy-to-use overlay for the automatically generated [Go SDK for oVirt](https://github.com/ovirt/go-ovirt). It does *not* replace the Go SDK. It implements the functions of the SDK only partially and is primarily used by the [oVirt Terraform provider](https://github.com/ovirt/terraform-provider-ovirt/).

## Using this library

To use this library you will have to include it as a Go module dependency:

```
go get github.com/ovirt/go-ovirt-client github.com/ovirt/go-ovirt-client-log/v3
```

You can then create a client instance like this:

```go
package main

import (
	"crypto/x509"

	ovirtclient "github.com/ovirt/go-ovirt-client/v3"
	ovirtclientlog "github.com/ovirt/go-ovirt-client-log/v3"
)

func main() {
	// Create a logger that logs to the standard Go log here:
	logger := ovirtclientlog.NewGoLogger()

	// Create an ovirtclient.TLSProvider implementation. This allows for simple
	// TLS configuration.
	tls := ovirtclient.TLS()

	// Add certificates from an in-memory byte slice. Certificates must be in PEM format.
	tls.CACertsFromMemory(caCerts)

	// Add certificates from a single file. Certificates must be in PEM format.
	tls.CACertsFromFile("/path/to/file.pem")

	// Add certificates from a directory. Optionally, regular expressions can be passed that must match the file
	// names.
	tls.CACertsFromDir(
		"/path/to/certs",
		regexp.MustCompile(`\.pem`),
	)

	// Add system certificates. This doesn't work on Windows before Go 1.18.
	tls.CACertsFromSystem()

	// Add a custom cert pool as a source of certificates. This option is
	// incompatible with CACertsFromSystem.
	// tls.CACertsFromCertPool(x509.NewCertPool())

	// Disable certificate verification. This is a bad idea, please don't do this.
	tls.Insecure()

	// Create a new goVirt instance:
	client, err := ovirtclient.New(
		// URL to your oVirt engine API here:
		"https://your-ovirt-engine/ovirt-engine/api/",
		// Username here:
		"admin@internal",
		// Password here:
		"password-here",
		// Pass the TLS provider here:
		tls,
		// Pass the logger here:
		logger,
		// Pass in extra settings here. Must implement the ovirtclient.ExtraSettings interface.
		nil,
	)
	if err != nil {
		// Handle error, here in a really crude way:
		panic(err)
	}
	// Use client. Please use the code completion in your IDE to
	// discover the functions. Each is well documented.
	upload, err := client.StartUploadToNewDisk(
		//...
	)
	//....
}
```

## Test helper

The test helper can work in two ways:

Either it sets up test fixtures in the mock client, or it sets up a live connection and identifies a usable storage
domain, cluster, etc. for testing purposes.

The ovirtclient.NewMockTestHelper() function can be used to create a test helper with a mock client in the backend:

```go
helper := ovirtclient.NewMockTestHelper(ovirtclientlog.NewNOOPLogger())
```

The easiest way to set up the test helper for a live connection is by using environment variables. To do that, you 
can use the `ovirtclient.NewLiveTestHelperFromEnv()` function:

```go
helper := ovirtclient.NewLiveTestHelperFromEnv(ovirtclientlog.NewNOOPLogger())
```

This function will inspect environment variables to determine if a connection to a live oVirt engine can be estabilshed.
The following environment variables are supported:

- `OVIRT_URL`: URL of the oVirt engine API.
- `OVIRT_USERNAME`: The username for the oVirt engine.
- `OVIRT_PASSWORD`: The password for the oVirt engine
- `OVIRT_CAFILE`: A file containing the CA certificate in PEM format.
- `OVIRT_CA_BUNDLE`: Provide the CA certificate in PEM format directly.
- `OVIRT_INSECURE`: Disable certificate verification if set. Not recommended.
- `OVIRT_CLUSTER_ID`: The cluster to use for testing. Will be automatically chosen if not provided.
- `OVIRT_BLANK_TEMPLATE_ID`: ID of the blank template. Will be automatically chosen if not provided.
- `OVIRT_STORAGE_DOMAIN_ID`: Storage domain to use for testing. Will be automatically chosen if not provided.
- `OVIRT_VNIC_PROFILE_ID`: VNIC profile to use for testing. Will be automatically chosen if not provided.

You can also create the test helper manually:

```go
import (
	"os"
	"testing"

	ovirtclient "github.com/ovirt/go-ovirt-client/v3"
	ovirtclientlog "github.com/ovirt/go-ovirt-client-log/v3"
)

func TestSomething(t *testing.T) {
	// Create a logger that logs to the standard Go log here
	logger := ovirtclientlog.NewTestLogger(t)

	// Set to true to use in-memory mock, otherwise this will use a live connection
	isMock := true

	// The following parameters define which infrastructure parts to use for testing
	params := ovirtclient.TestHelperParams().
		WithClusterID(ovirtclient.ClusterID(os.Getenv("OVIRT_CLUSTER_ID"))).
		WithBlankTemplateID(ovirtclient.TemplateID(os.Getenv("OVIRT_BLANK_TEMPLATE_ID"))).
		WithStorageDomainID(ovirtclient.StorageDomainID(os.Getenv("OVIRT_STORAGE_DOMAIN_ID"))).
		WithSecondaryStorageDomainID(ovirtclient.StorageDomainID(os.Getenv("OVIRT_SECONDARY_STORAGE_DOMAIN_ID"))).
		WithVNICProfileID(ovirtclient.VNICProfileID(os.Getenv("OVIRT_VNIC_PROFILE_ID")))

	// Create the test helper
	helper, err := ovirtclient.NewTestHelper(
		"https://localhost/ovirt-engine/api",
		"admin@internal",
		"super-secret",
		// Leave these empty for auto-detection / fixture setup
		params,
		ovirtclient.TLS().CACertsFromSystem(),
		isMock,
		logger,
	)
	if err != nil {
		t.Fatal(err)
	}
	// Fetch the cluster ID for testing
	clusterID := helper.GetClusterID()
	//...
}
```

**Tip:** You can use any logger that satisfies the `Logger` interface described in [go-ovirt-client-log](https://github.com/oVirt/go-ovirt-client-log)

## Retries

This library attempts to retry API calls that can be retried if possible. Each function has a sensible retry policy. However, you may want to customize the retries by passing one or more retry flags. The following retry flags are supported:

- `ovirtclient.ContextStrategy(ctx)`: this strategy will stop retries when the context parameter is canceled.
- `ovirtclient.ExponentialBackoff(factor)`: this strategy adds a wait time after each time, which is increased by the given factor on each try. The default is a backoff with a factor of 2.
- `ovirtclient.AutoRetry()`: this strategy will cancel retries if the error in question is a permanent error. This is enabled by default.
- `ovirtclient.MaxTries(tries)`: this strategy will abort retries if a maximum number of tries is reached. On complex calls the retries are counted per underlying API call.
- `ovirtclient.Timeout(duration)`: this strategy will abort retries if a certain time has been elapsed for the higher level call.
- `ovirtclient.CallTimeout(duration)`: this strategy will abort retries if a certain underlying API call takes longer than the specified duration. 

## Mock client

This library also provides a mock oVirt client that doesn't need working oVirt engine to function. It stores all information in-memory and simulates a working oVirt system. You can instantiate the mock client like so:

```go
client := ovirtclient.NewMock()
```

We recommend using the `ovirtclient.Client` interface as a means to declare it as a dependency in your factory so you can pass both the mock and the real connection as a parameter:

```go
func NewMyoVirtUsingUtility(
    client ovirtclient.Client,
) *myOVirtUsingUtility {
    return &myOVirtUsingUtility{
        client: client,
    }
}
``` 

## FAQ

### Why doesn't the library return the underlying oVirt SDK objects?

It's a painful decision we made. We want to encourage anyone who needs a certain function to submit a PR instead of simply relying on the SDK objects. This will lead to some overhead when a new function needs to be added, but leads to cleaner code in the end and makes this library more comprehensive. It also makes it possible to create the mock client, which would not be possibly if we had to simulate all parts of the oVirt engine.

If you need to access the oVirt SDK client you can do so from the `ovirtclient.New()` function:

```go
client, err := ovirtclient.New(
    //...
)
if err != nil {
    //...
}
sdkClient := client.GetSDKClient()
```

You can also get a properly preconfigured HTTP client if you need it:

```go
httpClient := client.GetHTTPClient()
```

**🚧 Warning:** If your code relies on the SDK or HTTP clients you will not be able to use the mock functionality described above for testing.

## Contributing

You want to help out? Awesome! Please head over to our [contribution guide](CONTRIBUTING.md), which explains how this library is built in detail.
