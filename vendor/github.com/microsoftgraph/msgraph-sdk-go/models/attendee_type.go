package models
import (
    "errors"
)
// Provides operations to manage the collection of agreement entities.
type AttendeeType int

const (
    REQUIRED_ATTENDEETYPE AttendeeType = iota
    OPTIONAL_ATTENDEETYPE
    RESOURCE_ATTENDEETYPE
)

func (i AttendeeType) String() string {
    return []string{"required", "optional", "resource"}[i]
}
func ParseAttendeeType(v string) (interface{}, error) {
    result := REQUIRED_ATTENDEETYPE
    switch v {
        case "required":
            result = REQUIRED_ATTENDEETYPE
        case "optional":
            result = OPTIONAL_ATTENDEETYPE
        case "resource":
            result = RESOURCE_ATTENDEETYPE
        default:
            return 0, errors.New("Unknown AttendeeType value: " + v)
    }
    return &result, nil
}
func SerializeAttendeeType(values []AttendeeType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
