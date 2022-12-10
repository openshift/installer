package models
import (
    "errors"
)
// Provides operations to manage the collection of agreement entities.
type SharedPCAccountDeletionPolicyType int

const (
    // Delete immediately.
    IMMEDIATE_SHAREDPCACCOUNTDELETIONPOLICYTYPE SharedPCAccountDeletionPolicyType = iota
    // Delete at disk space threshold.
    DISKSPACETHRESHOLD_SHAREDPCACCOUNTDELETIONPOLICYTYPE
    // Delete at disk space threshold or inactive threshold.
    DISKSPACETHRESHOLDORINACTIVETHRESHOLD_SHAREDPCACCOUNTDELETIONPOLICYTYPE
)

func (i SharedPCAccountDeletionPolicyType) String() string {
    return []string{"immediate", "diskSpaceThreshold", "diskSpaceThresholdOrInactiveThreshold"}[i]
}
func ParseSharedPCAccountDeletionPolicyType(v string) (interface{}, error) {
    result := IMMEDIATE_SHAREDPCACCOUNTDELETIONPOLICYTYPE
    switch v {
        case "immediate":
            result = IMMEDIATE_SHAREDPCACCOUNTDELETIONPOLICYTYPE
        case "diskSpaceThreshold":
            result = DISKSPACETHRESHOLD_SHAREDPCACCOUNTDELETIONPOLICYTYPE
        case "diskSpaceThresholdOrInactiveThreshold":
            result = DISKSPACETHRESHOLDORINACTIVETHRESHOLD_SHAREDPCACCOUNTDELETIONPOLICYTYPE
        default:
            return 0, errors.New("Unknown SharedPCAccountDeletionPolicyType value: " + v)
    }
    return &result, nil
}
func SerializeSharedPCAccountDeletionPolicyType(values []SharedPCAccountDeletionPolicyType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
