package models
import (
    "errors"
)
// 
type ConditionalAccessExternalTenantsMembershipKind int

const (
    ALL_CONDITIONALACCESSEXTERNALTENANTSMEMBERSHIPKIND ConditionalAccessExternalTenantsMembershipKind = iota
    ENUMERATED_CONDITIONALACCESSEXTERNALTENANTSMEMBERSHIPKIND
    UNKNOWNFUTUREVALUE_CONDITIONALACCESSEXTERNALTENANTSMEMBERSHIPKIND
)

func (i ConditionalAccessExternalTenantsMembershipKind) String() string {
    return []string{"all", "enumerated", "unknownFutureValue"}[i]
}
func ParseConditionalAccessExternalTenantsMembershipKind(v string) (any, error) {
    result := ALL_CONDITIONALACCESSEXTERNALTENANTSMEMBERSHIPKIND
    switch v {
        case "all":
            result = ALL_CONDITIONALACCESSEXTERNALTENANTSMEMBERSHIPKIND
        case "enumerated":
            result = ENUMERATED_CONDITIONALACCESSEXTERNALTENANTSMEMBERSHIPKIND
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_CONDITIONALACCESSEXTERNALTENANTSMEMBERSHIPKIND
        default:
            return 0, errors.New("Unknown ConditionalAccessExternalTenantsMembershipKind value: " + v)
    }
    return &result, nil
}
func SerializeConditionalAccessExternalTenantsMembershipKind(values []ConditionalAccessExternalTenantsMembershipKind) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
