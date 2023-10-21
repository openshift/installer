
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2022-03-01/containerapps` Documentation

The `containerapps` SDK allows for interaction with the Azure Resource Manager Service `containerapps` (API Version `2022-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2022-03-01/containerapps"
```


### Client Initialization

```go
client := containerapps.NewContainerAppsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ContainerAppsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := containerapps.NewContainerAppID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerAppValue")

payload := containerapps.ContainerApp{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ContainerAppsClient.Delete`

```go
ctx := context.TODO()
id := containerapps.NewContainerAppID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerAppValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ContainerAppsClient.Get`

```go
ctx := context.TODO()
id := containerapps.NewContainerAppID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerAppValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContainerAppsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := containerapps.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ContainerAppsClient.ListBySubscription`

```go
ctx := context.TODO()
id := containerapps.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ContainerAppsClient.ListCustomHostNameAnalysis`

```go
ctx := context.TODO()
id := containerapps.NewContainerAppID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerAppValue")

read, err := client.ListCustomHostNameAnalysis(ctx, id, containerapps.DefaultListCustomHostNameAnalysisOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContainerAppsClient.ListSecrets`

```go
ctx := context.TODO()
id := containerapps.NewContainerAppID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerAppValue")

read, err := client.ListSecrets(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContainerAppsClient.Update`

```go
ctx := context.TODO()
id := containerapps.NewContainerAppID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerAppValue")

payload := containerapps.ContainerApp{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
