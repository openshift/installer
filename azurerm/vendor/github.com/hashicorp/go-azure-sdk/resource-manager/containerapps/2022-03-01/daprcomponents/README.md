
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2022-03-01/daprcomponents` Documentation

The `daprcomponents` SDK allows for interaction with the Azure Resource Manager Service `containerapps` (API Version `2022-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2022-03-01/daprcomponents"
```


### Client Initialization

```go
client := daprcomponents.NewDaprComponentsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DaprComponentsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := daprcomponents.NewDaprComponentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentValue", "daprComponentValue")

payload := daprcomponents.DaprComponent{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DaprComponentsClient.Delete`

```go
ctx := context.TODO()
id := daprcomponents.NewDaprComponentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentValue", "daprComponentValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DaprComponentsClient.Get`

```go
ctx := context.TODO()
id := daprcomponents.NewDaprComponentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentValue", "daprComponentValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DaprComponentsClient.List`

```go
ctx := context.TODO()
id := daprcomponents.NewManagedEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DaprComponentsClient.ListSecrets`

```go
ctx := context.TODO()
id := daprcomponents.NewDaprComponentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentValue", "daprComponentValue")

read, err := client.ListSecrets(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
