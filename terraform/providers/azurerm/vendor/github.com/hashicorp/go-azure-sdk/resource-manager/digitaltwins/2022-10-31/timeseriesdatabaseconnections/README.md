
## `github.com/hashicorp/go-azure-sdk/resource-manager/digitaltwins/2022-10-31/timeseriesdatabaseconnections` Documentation

The `timeseriesdatabaseconnections` SDK allows for interaction with the Azure Resource Manager Service `digitaltwins` (API Version `2022-10-31`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/digitaltwins/2022-10-31/timeseriesdatabaseconnections"
```


### Client Initialization

```go
client := timeseriesdatabaseconnections.NewTimeSeriesDatabaseConnectionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `TimeSeriesDatabaseConnectionsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := timeseriesdatabaseconnections.NewTimeSeriesDatabaseConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "digitalTwinsInstanceValue", "timeSeriesDatabaseConnectionValue")

payload := timeseriesdatabaseconnections.TimeSeriesDatabaseConnection{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `TimeSeriesDatabaseConnectionsClient.Delete`

```go
ctx := context.TODO()
id := timeseriesdatabaseconnections.NewTimeSeriesDatabaseConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "digitalTwinsInstanceValue", "timeSeriesDatabaseConnectionValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `TimeSeriesDatabaseConnectionsClient.Get`

```go
ctx := context.TODO()
id := timeseriesdatabaseconnections.NewTimeSeriesDatabaseConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "digitalTwinsInstanceValue", "timeSeriesDatabaseConnectionValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TimeSeriesDatabaseConnectionsClient.List`

```go
ctx := context.TODO()
id := timeseriesdatabaseconnections.NewDigitalTwinsInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "digitalTwinsInstanceValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
