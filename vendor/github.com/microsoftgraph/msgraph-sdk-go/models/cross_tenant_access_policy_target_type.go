package models
import (
    "errors"
)
// Provides operations to manage the collection of application entities.
type CrossTenantAccessPolicyTargetType int

const (
    USER_CROSSTENANTACCESSPOLICYTARGETTYPE CrossTenantAccessPolicyTargetType = iota
    GROUP_CROSSTENANTACCESSPOLICYTARGETTYPE
    APPLICATION_CROSSTENANTACCESSPOLICYTARGETTYPE
    UNKNOWNFUTUREVALUE_CROSSTENANTACCESSPOLICYTARGETTYPE
)

func (i CrossTenantAccessPolicyTargetType) String() string {
    return []string{"user", "group", "application", "unknownFutureValue"}[i]
}
func ParseCrossTenantAccessPolicyTargetType(v string) (interface{}, error) {
    result := USER_CROSSTENANTACCESSPOLICYTARGETTYPE
    switch v {
        case "user":
            result = USER_CROSSTENANTACCESSPOLICYTARGETTYPE
        case "group":
            result = GROUP_CROSSTENANTACCESSPOLICYTARGETTYPE
        case "application":
            result = APPLICATION_CROSSTENANTACCESSPOLICYTARGETTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_CROSSTENANTACCESSPOLICYTARGETTYPE
        default:
            return 0, errors.New("Unknown CrossTenantAccessPolicyTargetType value: " + v)
    }
    return &result, nil
}
func SerializeCrossTenantAccessPolicyTargetType(values []CrossTenantAccessPolicyTargetType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
