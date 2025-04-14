# session package

This library provides a simple and consistent REST over HTTP library for accessing Akamai Endpoints

## Depedencies 

This library is dependent on the `github.com/akamai/AkamaiOPEN-edgegrid-golang/pkg/edgegrid` interface.

## Basic Example

```
func main() {
     edgerc := Must(New())

     s, err := session.New(
         session.WithConfig(edgerc),
     )
     if err != nil {
         panic(err)
     }

     var contracts struct {
		AccountID string         `json:"accountId"`
		Contracts ContractsItems `json:"contracts"`
        Items []struct {
            ContractID       string `json:"contractId"`
		    ContractTypeName string `json:"contractTypeName"`
        } `json:"items"`
     }

     req, _ := http.NewRequest(http.MethodGet, "/papi/v1/contracts", nil)

     _, err := s.Exec(r, &contracts)
     if err != nil {
         panic(err);
     }

     // do something with contracts
}
        
```

## Library Logging
The session package supports the structured logging interface from `github.com/apex`. These can be applied globally to the session or to the request context.

### Adding a logger to the session

```
    s, err := session.New(
         session.WithConfig(edgerc),
         session.WithLog(log.Log),
     )
     if err != nil {
         panic(err)
     }
```

### Request logging
The logger can be overidden for a specific request like this

```
    req, _ := http.NewRequest(http.MethodGet, "/papi/v1/contracts", nil)

    req = req.WithContext(
        session.ContextWithOptions(request.Context(),
            session.WithContextLog(otherlog),
        )
```

## Custom request headers
The context can also be updated to pass special http headers when necessary

```
    customHeader := make(http.Header)
    customHeader.Set("X-Custom-Header", "some custom value")

    req = req.WithContext(
        session.ContextWithOptions(request.Context(),
            session.WithContextHeaders(customHeader),
        )
```