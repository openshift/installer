# Contributing to go-ovirt-client

Thank you for wanting to contribute to go-ovirt-client! This document will describe the most important tooling, as well as the design concepts behind this library.

## Concept

This library is built to wrap [go-ovirt](https://github.com/ovirt/go-ovirt), which is an automatically generated library for oVirt. This library is mostly hand-written and aims to provide an easier to use interface and several convenience functions, as well as a proper abstraction to enable mocking for tests. Having mocks eliminates the need for an actual running oVirt engine to run tests.

## Tooling

The main tool to develop this library apart from Go itself is [golangci-lint](https://golangci-lint.run). You can use this tool to sanity-check your code before submitting a PR.

```
golangci-lint run
```

The provided [.golangci.yaml](.golangci.yml) describes the linting rules we are using. If you have a good reason to skip a rule, you can add a comment with the rule:

```go
// We are working through all template files here, so
// including these files is intentional and not a
// security issue.
fh, err := os.Open(templateFileName) // nolint:gosec
```

## Design principles

This library serves several projects, including the oVirt Terraform provider, the OpenShift CSI driver, etc. Since it is used in many places, the main consideration is **ease of use** and **providing mocks for tests**.

## Client interfaces

The main point of entry for most people will be the `New()` function located in [new.go](new.go). This function returns the `Client` interface. 

The `Client` interface in [client.go](client.go) aggregates all sub-interfaces, for example the `DiskClient`, `VMClient`, etc. This is done so that an application using this library can rely on a subset of the functionality. This, in turn, enables writing mocks for only a part of the functionality.

Each client function must be implemented in two copies: once for the live connection in the `client_` files, and once for the mock functionality in the `mock_` files.

## Auto-generated implementations

Some of the client implementations, such as Get and List functions are auto-generated using `go generate`. This is done using the generator located in [scripts/rest.go](scripts/rest.go).

## API objects

API calls, such as `GetVM`, will return API objects. These API objects are described by two interfaces instead of just returning structs. Using interfaces makes the library extensible in the future.

The two interfaces are `OBJECTNAME` and `OBJECTNAMEData`. For example, `VM` and `VMData`. The `*Data` interface describes the functions required for retrieving fields on this object, while the other interface incorporates the `*Data` object and adds the ability to call client functionality directly. For example, the `Disk` object will have a `Remove()` function on it. This is, again, done to enable easier mock-writing.

## Running tests

You can run the tests against the mock backend simply by running:

```
go generate
go test -v -client=mock ./...
```

If you don't specify `-client`, it will default to `all` and run each test against the mock backend as well as the live oVirt engine.  

Before you submit, we would like to ask you to run your tests against the live oVirt engine as we do not have one integrated with the CI at the moment. You can do so by running tests as follows:

```
OVIRT_URL=https://url-of-your-engine
OVIRT_USERNAME=admin@internal
OVIRT_PASSWORD=your-ovirt-password
# Use the system certificate store to verify the engine certificate.
OVIRT_SYSTEM=1
# Alternatively, use one of these options:
## Pass the certificate as a file:
# OVIRT_CA_FILE=/path/to/certificate.pem
## Read all files in a directory for certificates:
# OVIRT_CA_DIR=/path/to/certificates
## Pass the certificate directly:
# OVIRT_CA_CERT=cert-data-here
## Disable certificate verification:
# OVIRT_INSECURE=1

# Run the tests
go test -v -client=live ./...
```

If you want to connect to a live oVirt engine you need to define these environment variables.
To get a PR merged please run your tests against both the mock and the live backend.

In the test code you can then obtain the test helper using the `getHelper(t)` function:

```
helper := getHelper(t)
```

The client is then available using the `helper.GetClient()` function. 

## Submitting a PR

Once you are ready, please submit a pull request on GitHub. The pull request template will ask 3 questions, please answer them so we can merge your PR quickly:

- *Please describe the change you are making* - So we know what the motivation behind your change is.
- *Are you the owner of the code you are sending in, or do you have permission of the owner?* - This is important so we know that the license of your contribution is clean.
- *The code will be published under the BSD 3 clause license. Have you read and understood this license?* - So we know that you have read the [LICENSE file](LICENSE) and you are OK with your code being published under this license.

**Note:** If you are a member of the oVirt organization you may skip the last 2 questions.

Once your PR is submitted, GitHub actions will run on your PR and test your code against the mock backend. We currently do not have a live oVirt Engine for testing, so our maintainers will have to check your PR manually. This may take some time, but we should get back to you within a few days.

### Breaking changes

This library follows the Go module specification. Once v1 is released, any backwards-incompatible change (e.g. adding a field to an interface, etc) will require increasing the major version number. If your change involves a breaking change it may take a bit longer to get merged.

## Using your change

This library is versioned according to SemVer, as required by the Go module specification. That means, that your change will be in the next release once your PR is merged. Until then, you can grab the latest version for testing like this:

```
go get github.com/ovirt/go-ovirt-client@main
```

This will create a `go.mod` entry like this:

```
require(
  github.com/ovirt/go-ovirt-client v0.6.1-0.20210913092754-237ed78d23e7
)
```

However, we do not recommend using this in production. If you need to go into production soon, and you need your change tagged quickly, please open a separate issue so we know it's urgent.

## Recommended further reading

- [Go Patterns: Object-Oriented Programming](https://debugged.it/blog/go-patterns-oop/)
- [Go Patterns: Retries](https://debugged.it/blog/go-patterns-retries/)
- [Testing the oVirt Terraform Provider (video)](https://www.youtube.com/watch?v=1eciby2NiSM)
