# OCM SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/openshift-online/ocm-sdk-go.svg)](https://pkg.go.dev/github.com/openshift-online/ocm-sdk-go)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

This project contains a Go library that simplifies the use of the _OCM_
API, available in `api.openshift.com`.

## Usage

To use it import the `github.com/openshift-online/ocm-sdk-go` package, and then
use it to send requests to the API.

Note that the name of the directory is `ocm-sdk-go` but the name of the package
is just `sdk`, so to use it you will have to import it and then use `sdk` as
the package selector.

For example, if you need to create a cluster you can use the following code:

```go
package main

import (
        "fmt"
        "os"

        sdk "github.com/openshift-online/ocm-sdk-go"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

func main() {
	// Create a logger that has the debug level enabled:
	logger, err := sdk.NewGoLoggerBuilder().
		Debug(true).
		Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't build logger: %v\n", err)
		os.Exit(1)
	}

	// Create the connection, and remember to close it:
	token := os.Getenv("OCM_TOKEN")
	connection, err := sdk.NewConnectionBuilder().
		Logger(logger).
		Tokens(token).
		Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't build connection: %v\n", err)
		os.Exit(1)
	}
	defer connection.Close()

	// Get the client for the resource that manages the collection of clusters:
	collection := connection.ClustersMgmt().V1().Clusters()

	// Prepare the description of the cluster to create:
	cluster, err := cmv1.NewCluster().
		Name("mycluster").
		CloudProvider(
			cmv1.NewCloudProvider().
				ID("aws"),
		).
		Region(
			cmv1.NewCloudRegion().
				ID("us-east-1"),
		).
		Version(
			cmv1.NewVersion().
				ID("openshift-v4.0-beta4"),
		).
		Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't create cluster description: %v\n", err)
		os.Exit(1)
	}

	// Send a request to create the cluster:
	response, err := collection.Add().
		Body(cluster).
		Send()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't create cluster: %v\n", err)
		os.Exit(1)
	}

	// Print the result:
	cluster = response.Body()
	fmt.Printf("%s - %s\n", cluster.ID(), cluster.Name())
}
```

There are more examples in the [examples](examples) directory.

### Packages

The following are the packages that are most frequently needed in order to use
the SDK:

**main**

This is the top level package. The most important element is the `Connection`
type, as it is the mechanism to connect to the server and to get the reference
to the clients for the services that are part of the API.

**errors**

Contains the `Error` type that is used by the SDK to report errors.

**accesstransparency/v1** 

This package contains the types and clients for version 1 of the access
transparency service.

**accountsmgmt/v1**

This package contains the types and clients for version 1 of the accounts
management service.

**authorizations/v1**

This package contains the types and clients for version 1 of the
authorizations service.

**clustersmgmt/v1**

This package contains the types and clients for version 1 of the clusters
management service.

There are other packages, like `helpers` and `internal`.  Those contain
internal implementation details of the SDK. Refrain from using them, as they
may change in the future: backwards compatibility isn't guaranteed.

### Connecting to the server

To connect to the server import the `sdk` package. That contains the
`Connection` type, which is the entry point of the SDK, and gives you access to
the clients for the services that are part of the API:

```go
import (
	"github.com/openshift-online/ocm-sdk-go"
)

// Create the connection:
connection, err := sdk.NewConnectionBuilder().
	Tokens(token).
	Build()
if err != nil {
        fmt.Fprintf(os.Stderr, "Can't build connection: %v\n", err)
        os.Exit(1)
}
```

The connection holds expensive resources, including a pool of HTTP connections
to the server and an authentication token. It is important to release those
resources when they are no longer in use:

```go
// Close the connection:
connection.Close()
```

Consider using the _defer_ mechanism to ensure that the connection is always
closed when no longer needed.

### Using _types_

The Go types that correspond to the API data types live in the
`accountsmgmt/v1`, `authorizations/v1`, and `clustersmgmt/v1` packages. These types are pure data
containers, they don't have any logic or operation.  Instances can be created
at will.

Creation of objects of these types does *not* have any effect in the server
side, unless the object is explicitly passed to a call to one of the resource
methods described below. Changes in the server side are *not* automatically
reflected in the instances that already exist in memory.

Creation of objects of these types is done using the corresponding _builder_
type. For example, to create an object of the `Cluster` type create an object of
the `ClusterBuilder` type (using the `NewCluster` function) populate and then
build the object calling the `Build` method:

```go
// Create a new object of the `Cluster` type:
cluster, err := cmv1.NewCluster().
	Name("mycluster").
	CloudProvider(
		cmv1.NewCloudProvider().
			ID("aws"),
	).
	Region(
		cmv1.NewCloudRegion().
			ID("us-east-1"),
	).
	Version(
		cmv1.NewVersion().
			ID("openshift-v4.9.7"),
	).
	Build()
if err != nil {
	fmt.Fprintf(os.Stderr, "Can't create cluster object: %v\n", err)
	os.Exit(1)
}
```

Once created objects are immutable.

The fields containing the values of the attributes of these types are private.
To read them use the _access methods_. For example, to read the value of the
`name` attribute of a cluster:

```go
// Get the value of the `name` attribute:
name := cluster.Name()
fmt.Printf("Cluster name is '%s'\n", name)
```

The access methods return the value of the attribute, if it has a value, or the
zero value of the type (`""` for strings, `false` for booleans, `0` for
integers, etc) if the attribute doesn't have a value. That makes it impossible
to know if the attribute has a value or not. If you need that, use the `Get...`
variant of the accessor. For example, to get the value of the `name` attribute
and also check if the attribute has a value:

```go
// Get the value of the `name` attribute, and check if it has a value:
name, ok := cluster.GetName()
if !ok {
	fmt.Printf("Cluster has no name\n")
} else {
	fmt.Printf("Cluster name is '%s'\n", name)
}
```

Attributes that are defined as list of objects in the specification of the API
are implemented as objects of a `...List` type. For example, the value of the
`groups` attribute of the `Cluster` type is implemented as the `GroupList` type.
These list types provide methods to process the elements of the list. For
example, to print the names of a list of groups:

```go
// Get the list of groups:
groups := ...

// Print the name of each group:
groups.Each(func(group *cmv1.Group) bool {
	fmt.Printf("Group name is '%s'\n", group.Name())
	return true
})
```

The function passed to the `Each` method will be called once for each item of
the list. If it returns `true` the iteration will continue, otherwise will stop.
This is intended to mimic a `for` loop with an optional `break`.

If it is necessary to have access to the index of the item, then it is better to
use the `Range` method:

```go
// Get the list of groups:
groups := ...

// Print index and name of each group:
groups.Range(func(int i, group *cmv1.Group) bool {
	fmt.Printf("Group index is %d and is '%s'\n", i, group.Name())
	return true
})
```

It is also possible to convert the list to an slice, using the `Slice` method,
and the process it as usual:

```go
// Get the list of groups:
groups := ...

// Print the name of each group:
slice := groups.Slice()
for _, group := range slice {
	fmt.Printf("Group name is '%s'\n", group.Name())
}
```

It is in general better to use the `Each` or `Range` methods instead of the
`Slice` method, because `Slice` has the additional cost of allocating that slice
and copying the internal representation into it.

## CLI

See also the command-line tool https://github.com/openshift-online/ocm-cli built
on top of this SDK.
