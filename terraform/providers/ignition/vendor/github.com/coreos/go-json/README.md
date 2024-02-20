# go-json

This is a fork of Go's `encoding/json` library. It adds a third target for
unmarshalling, `json.Node`.

Unmarshalling to a `Node` behaves similarly to unmarshalling to an
`interface{}`, except that it also records the offsets for the start and end
of the value that is unmarshalled and, if the value is part of a JSON
object, the offsets of the start and end of the object's key. The `Value`
field of the `Node` is unmarshalled to the same type as if it were an
`interface{}`, except in the case of arrays and objects:

| JSON type | Go type, unmarshalled to `interface{}` | `Node.Value` type |
| --------- | -------------------------------------- | ----------------- |
| Array     | `[]interface{}`                        | `[]Node`          |
| Object    | `map[string]interface{}`               | `map[string]Node` |
| Other     | `interface{}`                          | `interface{}`     |
