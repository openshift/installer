package models
import (
    "errors"
)
// Provides operations to manage the collection of agreement entities.
type ManagedAppFlaggedReason int

const (
    // No issue.
    NONE_MANAGEDAPPFLAGGEDREASON ManagedAppFlaggedReason = iota
    // The app registration is running on a rooted/unlocked device.
    ROOTEDDEVICE_MANAGEDAPPFLAGGEDREASON
)

func (i ManagedAppFlaggedReason) String() string {
    return []string{"none", "rootedDevice"}[i]
}
func ParseManagedAppFlaggedReason(v string) (interface{}, error) {
    result := NONE_MANAGEDAPPFLAGGEDREASON
    switch v {
        case "none":
            result = NONE_MANAGEDAPPFLAGGEDREASON
        case "rootedDevice":
            result = ROOTEDDEVICE_MANAGEDAPPFLAGGEDREASON
        default:
            return 0, errors.New("Unknown ManagedAppFlaggedReason value: " + v)
    }
    return &result, nil
}
func SerializeManagedAppFlaggedReason(values []ManagedAppFlaggedReason) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
