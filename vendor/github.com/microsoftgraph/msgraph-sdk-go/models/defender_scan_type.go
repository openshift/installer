package models
import (
    "errors"
)
// Provides operations to manage the collection of agreement entities.
type DefenderScanType int

const (
    // User Defined, default value, no intent.
    USERDEFINED_DEFENDERSCANTYPE DefenderScanType = iota
    // System scan disabled.
    DISABLED_DEFENDERSCANTYPE
    // Quick system scan.
    QUICK_DEFENDERSCANTYPE
    // Full system scan.
    FULL_DEFENDERSCANTYPE
)

func (i DefenderScanType) String() string {
    return []string{"userDefined", "disabled", "quick", "full"}[i]
}
func ParseDefenderScanType(v string) (interface{}, error) {
    result := USERDEFINED_DEFENDERSCANTYPE
    switch v {
        case "userDefined":
            result = USERDEFINED_DEFENDERSCANTYPE
        case "disabled":
            result = DISABLED_DEFENDERSCANTYPE
        case "quick":
            result = QUICK_DEFENDERSCANTYPE
        case "full":
            result = FULL_DEFENDERSCANTYPE
        default:
            return 0, errors.New("Unknown DefenderScanType value: " + v)
    }
    return &result, nil
}
func SerializeDefenderScanType(values []DefenderScanType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
