/*
Package ovirtclient provides a human-friendly Go client for the oVirt Engine. It provides an abstraction layer
for the oVirt API, as well as a mocking facility for testing purposes.

Reading this documentation

This documentation contains two parts. This introduction explains setting up the client with the credentials. The API
doc contains the individual API calls.

When reading the API doc, start with the Client interface: it contains all components of the API. The individual
API's, their documentation and examples are located in subinterfaces, such as DiskClient.

Creating a client

There are several ways to create a client instance. The most basic way is to use the New() function as follows:

    // Create the client
    client, err := ovirtclient.New(
        // URL
        "https://localhost/ovirt-engine/api",
        // Username
        "admin@internal",
        // Password
        "super-secret",
        // Pull CA certificates from the operating system.
        // This won't work on Windows before Go 1.18. See below for an extended example.
        ovirtclient.TLS().CACertsFromSystem(),
        // Don't log.
        ovirtclientlog.NewNOOPLogger(),
        // No extra connection settings.
        nil,
    )
    if err != nil {
        panic(fmt.Errorf("failed to create oVirt client (%w)", err))
    }

Mock client

The mock client simulates the oVirt engine behavior in-memory without needing an actual running engine. This is a
good way to provide a testing facility.

It can be created using the NewMock method:

    client := ovirtclient.NewMock()

That's it! However, to make it really useful, you will need the test helper which can set up test fixtures.

Test helper

The test helper can work in two ways:

Either it sets up test fixtures in the mock client, or it sets up a live connection and identifies a usable storage
domain, cluster, etc. for testing purposes.

The ovirtclient.NewMockTestHelper() function can be used to create a test helper with a mock client in the backend:

    helper := ovirtclient.NewMockTestHelper(ovirtclientlog.NewNOOPLogger())

The easiest way to set up the test helper for a live connection is by using environment variables. To do that, you
can use the ovirtclient.NewLiveTestHelperFromEnv() function:

    helper := ovirtclient.NewLiveTestHelperFromEnv(ovirtclientlog.NewNOOPLogger())

This function will inspect environment variables to determine if a connection to a live oVirt engine can be
established. The following environment variables are supported:

  OVIRT_URL

URL of the oVirt engine API. Mandatory.

  OVIRT_USERNAME

The username for the oVirt engine. Mandatory.

  OVIRT_PASSWORD

The password for the oVirt engine. Mandatory.

  OVIRT_CAFILE

A file containing the CA certificate in PEM format.

  OVIRT_CA_BUNDLE

Provide the CA certificate in PEM format directly.

  OVIRT_INSECURE

Disable certificate verification if set. Not recommended.

  OVIRT_CLUSTER_ID

The cluster to use for testing. Will be automatically chosen if not provided.

  OVIRT_BLANK_TEMPLATE_ID

ID of the blank template. Will be automatically chosen if not provided.

  OVIRT_STORAGE_DOMAIN_ID

Storage domain to use for testing. Will be automatically chosen if not provided.

  OVIRT_VNIC_PROFILE_ID

VNIC profile to use for testing. Will be automatically chosen if not provided.

You can also create the test helper manually:

    import (
        "os"
        "testing"

        ovirtclient "github.com/ovirt/go-ovirt-client/v3"
        ovirtclientlog "github.com/ovirt/go-ovirt-client-log"
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

Logging

This library provides extensive logging. Each API interaction is logged on the debug level, and other messages are
added on other levels. In order to provide logging this library uses the go-ovirt-client-log
(https://github.com/oVirt/go-ovirt-client-log) interface definition.

As long as your logger implements this interface, you will be able to receive log messages. The logging
library also provides a few built-in loggers. For example, you can log via the default Go log interface:

    logger := ovirtclientlog.NewGoLogger()

Or, you can also log in tests:

    logger := ovirtclientlog.NewTestLogger(t)

You can also disable logging:

    logger := ovirtclientlog.NewNOOPLogger()

Finally, we also provide an adapter library for klog here: https://github.com/oVirt/go-ovirt-client-log-klog

TLS verification

Modern-day oVirt engines run secured with TLS. This means that the client needs a way to verify the certificate the
server is presenting. This is controlled by the tls parameter of the New() function. You can implement your own source
by implementing the TLSProvider interface, but the package also includes a ready-to-use provider.

Create the provider using the TLS() function:

    tls := ovirtclient.TLS()

This provider has several functions. The easiest to set up is using the system trust root for certificates. However,
this won't work own Windows:

    tls.CACertsFromSystem()

Now you need to add your oVirt engine certificate to your system trust root.

If you don't want to, or can't add the certificate to the system trust root, you can also directly provide it
to the client.

    // Add certificates from a certificate pool you have previously initialized.
    tls.CACertsFromCertPool(certpool)

    // Add certificates from an in-memory byte slice. Certificates must be in PEM format.
    tls.CACertsFromMemory([]byte("-----BEGIN CERTIFICATE-----\n..."))

    // Add certificates from a single file. Certificates must be in PEM format.
    tls.CACertsFromFile("/path/to/file.pem")

    // Add certificates from a directory. Optionally, regular expressions can be passed that must match the file
    // names.
    tls.CACertsFromDir("/path/to/certs", regexp.MustCompile(`\.pem`))

Finally, you can also disable certificate verification. Do we need to say that this is a very, very bad idea?

    tls.Insecure()

The configured tls variable can then be passed to the New() function to create an oVirt client.

Retries

This library attempts to retry API calls that can be retried if possible. Each function has a sensible retry policy.
However, you may want to customize the retries by passing one or more retry flags. The following retry flags are
supported:

    ovirtclient.ContextStrategy(ctx)

This strategy will stop retries when the context parameter is canceled.

    ovirtclient.ExponentialBackoff(factor)

This strategy adds a wait time after each time, which is increased by the given factor on each try. The default is a
backoff with a factor of 2.

    ovirtclient.AutoRetry()

This strategy will cancel retries if the error in question is a permanent error. This is enabled by default.

    ovirtclient.MaxTries(tries)

This strategy will abort retries if a maximum number of tries is reached. On complex calls the retries are counted per
underlying API call.

    ovirtclient.Timeout(duration)

This strategy will abort retries if a certain time has been elapsed for the higher level call.

    ovirtclient.CallTimeout(duration)

This strategy will abort retries if a certain underlying API call takes longer than the specified duration.

*/
package ovirtclient
