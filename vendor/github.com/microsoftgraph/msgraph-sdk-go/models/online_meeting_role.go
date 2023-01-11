package models
import (
    "errors"
)
// Provides operations to manage the collection of agreement entities.
type OnlineMeetingRole int

const (
    ATTENDEE_ONLINEMEETINGROLE OnlineMeetingRole = iota
    PRESENTER_ONLINEMEETINGROLE
    UNKNOWNFUTUREVALUE_ONLINEMEETINGROLE
    PRODUCER_ONLINEMEETINGROLE
)

func (i OnlineMeetingRole) String() string {
    return []string{"attendee", "presenter", "unknownFutureValue", "producer"}[i]
}
func ParseOnlineMeetingRole(v string) (interface{}, error) {
    result := ATTENDEE_ONLINEMEETINGROLE
    switch v {
        case "attendee":
            result = ATTENDEE_ONLINEMEETINGROLE
        case "presenter":
            result = PRESENTER_ONLINEMEETINGROLE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_ONLINEMEETINGROLE
        case "producer":
            result = PRODUCER_ONLINEMEETINGROLE
        default:
            return 0, errors.New("Unknown OnlineMeetingRole value: " + v)
    }
    return &result, nil
}
func SerializeOnlineMeetingRole(values []OnlineMeetingRole) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
