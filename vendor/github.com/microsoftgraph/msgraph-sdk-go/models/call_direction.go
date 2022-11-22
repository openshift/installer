package models
import (
    "errors"
)
// Provides operations to manage the cloudCommunications singleton.
type CallDirection int

const (
    INCOMING_CALLDIRECTION CallDirection = iota
    OUTGOING_CALLDIRECTION
)

func (i CallDirection) String() string {
    return []string{"incoming", "outgoing"}[i]
}
func ParseCallDirection(v string) (interface{}, error) {
    result := INCOMING_CALLDIRECTION
    switch v {
        case "incoming":
            result = INCOMING_CALLDIRECTION
        case "outgoing":
            result = OUTGOING_CALLDIRECTION
        default:
            return 0, errors.New("Unknown CallDirection value: " + v)
    }
    return &result, nil
}
func SerializeCallDirection(values []CallDirection) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
