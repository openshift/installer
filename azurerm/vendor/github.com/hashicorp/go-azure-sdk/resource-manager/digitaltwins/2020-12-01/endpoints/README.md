
## `github.com/hashicorp/go-azure-sdk/resource-manager/digitaltwins/2020-12-01/endpoints` Documentation

The `endpoints` SDK allows for interaction with the Azure Resource Manager Service `digitaltwins` (API Version `2020-12-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/digitaltwins/2020-12-01/endpoints"
```


### Client Initialization

```go
client := endpoints.NewEndpointsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `EndpointsClient.DigitalTwinsEndpointCreateOrUpdate`

```go
ctx := context.TODO()
id := endpoints.NewEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "digitalTwinsInstanceValue", "endpointValue")

payload := endpoints.DigitalTwinsEndpointResource{
	// ...
}


if err := client.DigitalTwinsEndpointCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `EndpointsClient.DigitalTwinsEndpointDelete`

```go
ctx := context.TODO()
id := endpoints.NewEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "digitalTwinsInstanceValue", "endpointValue")

if err := client.DigitalTwinsEndpointDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `EndpointsClient.DigitalTwinsEndpointGet`

```go
ctx := context.TODO()
id := endpoints.NewEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "digitalTwinsInstanceValue", "endpointValue")

read, err := client.DigitalTwinsEndpointGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EndpointsClient.DigitalTwinsEndpointList`

```go
ctx := context.TODO()
id := endpoints.NewDigitalTwinsInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "digitalTwinsInstanceValue")

// alternatively `client.DigitalTwinsEndpointList(ctx, id)` can be used to do batched pagination
items, err := client.DigitalTwinsEndpointListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
