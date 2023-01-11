package models
import (
    "errors"
)
// Provides operations to manage the cloudCommunications singleton.
type RoutingMode int

const (
    ONETOONE_ROUTINGMODE RoutingMode = iota
    MULTICAST_ROUTINGMODE
    UNKNOWNFUTUREVALUE_ROUTINGMODE
)

func (i RoutingMode) String() string {
    return []string{"oneToOne", "multicast", "unknownFutureValue"}[i]
}
func ParseRoutingMode(v string) (interface{}, error) {
    result := ONETOONE_ROUTINGMODE
    switch v {
        case "oneToOne":
            result = ONETOONE_ROUTINGMODE
        case "multicast":
            result = MULTICAST_ROUTINGMODE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_ROUTINGMODE
        default:
            return 0, errors.New("Unknown RoutingMode value: " + v)
    }
    return &result, nil
}
func SerializeRoutingMode(values []RoutingMode) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
