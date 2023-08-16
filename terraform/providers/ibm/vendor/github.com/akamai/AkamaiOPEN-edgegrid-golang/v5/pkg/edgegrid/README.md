# edgegrid authorization library

This library provides Akamai `.edgerc` configuration parsing and `http.Request` signing.

## EdgeGrid Configuration Files

The default location for the `.edgerc` file is `$HOME/.edgerc`. This file has a standard ini-file format. The default section is `default`. Multiple sections can be stored in the same ini-file for other configurations such as production, staging, or development.

```
[default]
client_secret = <default secret>
host = <default host>
access_token = <default access token>
client_token = <default client token>

[dev]
client_secret = <dev secret>
host = <dev host>
access_token = <dev access token>
client_token = <dev client token>
```

## Basic Example

```
func main() {
    edgerc := Must(New())
    
    client := http.Client{}

    req, _ := http.NewRequest(http.MethodGet, "/papi/v1/contracts", nil)

    edgerc.SignRequest(req)

    resp, err := client.Do(req)
    if err != nil {
        log.Fataln(err)
    }

    // do something with response
}
```

## Using a custom `.edgerc` file and section

```
    edgerc := Must(New(
        WithFile("/some/other/edgerc"),
        WithSection("production"),
    ))
}
```

## Loading from environment variables

By default, it uses `AKAMAI_HOST`, `AKAMAI_CLIENT_TOKEN`, `AKAMAI_CLIENT_SECRET`, `AKAMAI_ACCESS_TOKEN`, and `AKAMAI_MAX_BODY` variables.

You can define multiple configurations by prefixing with the section name specified, e.g. passing "ccu" will cause it to look for `AKAMAI_CCU_HOST`, etc.

If `AKAMAI_{SECTION}` does not exist, it will fall back to just `AKAMAI_`.

```
    // Load from AKAMA_CCU_
    edgerc := Must(New(
        WithEnv(true),
        WithSection("ccu"),
    ))
}
```
