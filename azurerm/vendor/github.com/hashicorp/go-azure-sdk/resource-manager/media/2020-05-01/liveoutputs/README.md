
## `github.com/hashicorp/go-azure-sdk/resource-manager/media/2020-05-01/liveoutputs` Documentation

The `liveoutputs` SDK allows for interaction with the Azure Resource Manager Service `media` (API Version `2020-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/media/2020-05-01/liveoutputs"
```


### Client Initialization

```go
client := liveoutputs.NewLiveOutputsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `LiveOutputsClient.Create`

```go
ctx := context.TODO()
id := liveoutputs.NewLiveOutputID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "liveEventValue", "liveOutputValue")

payload := liveoutputs.LiveOutput{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `LiveOutputsClient.Delete`

```go
ctx := context.TODO()
id := liveoutputs.NewLiveOutputID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "liveEventValue", "liveOutputValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `LiveOutputsClient.Get`

```go
ctx := context.TODO()
id := liveoutputs.NewLiveOutputID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "liveEventValue", "liveOutputValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LiveOutputsClient.List`

```go
ctx := context.TODO()
id := liveoutputs.NewLiveEventID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "liveEventValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
