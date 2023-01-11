package models
import (
    "errors"
)
// Provides operations to manage the collection of agreement entities.
type ApplicationType int

const (
    // The windows universal application
    UNIVERSAL_APPLICATIONTYPE ApplicationType = iota
    // The windows desktop application
    DESKTOP_APPLICATIONTYPE
)

func (i ApplicationType) String() string {
    return []string{"universal", "desktop"}[i]
}
func ParseApplicationType(v string) (interface{}, error) {
    result := UNIVERSAL_APPLICATIONTYPE
    switch v {
        case "universal":
            result = UNIVERSAL_APPLICATIONTYPE
        case "desktop":
            result = DESKTOP_APPLICATIONTYPE
        default:
            return 0, errors.New("Unknown ApplicationType value: " + v)
    }
    return &result, nil
}
func SerializeApplicationType(values []ApplicationType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
