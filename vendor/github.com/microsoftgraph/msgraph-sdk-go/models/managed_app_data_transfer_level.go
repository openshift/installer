package models
import (
    "errors"
)
// Provides operations to manage the collection of agreement entities.
type ManagedAppDataTransferLevel int

const (
    // All apps.
    ALLAPPS_MANAGEDAPPDATATRANSFERLEVEL ManagedAppDataTransferLevel = iota
    // Managed apps.
    MANAGEDAPPS_MANAGEDAPPDATATRANSFERLEVEL
    // No apps.
    NONE_MANAGEDAPPDATATRANSFERLEVEL
)

func (i ManagedAppDataTransferLevel) String() string {
    return []string{"allApps", "managedApps", "none"}[i]
}
func ParseManagedAppDataTransferLevel(v string) (interface{}, error) {
    result := ALLAPPS_MANAGEDAPPDATATRANSFERLEVEL
    switch v {
        case "allApps":
            result = ALLAPPS_MANAGEDAPPDATATRANSFERLEVEL
        case "managedApps":
            result = MANAGEDAPPS_MANAGEDAPPDATATRANSFERLEVEL
        case "none":
            result = NONE_MANAGEDAPPDATATRANSFERLEVEL
        default:
            return 0, errors.New("Unknown ManagedAppDataTransferLevel value: " + v)
    }
    return &result, nil
}
func SerializeManagedAppDataTransferLevel(values []ManagedAppDataTransferLevel) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
