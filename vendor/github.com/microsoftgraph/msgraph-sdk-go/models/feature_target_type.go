package models
import (
    "errors"
)
// Provides operations to manage the collection of authenticationMethodConfiguration entities.
type FeatureTargetType int

const (
    GROUP_FEATURETARGETTYPE FeatureTargetType = iota
    ADMINISTRATIVEUNIT_FEATURETARGETTYPE
    ROLE_FEATURETARGETTYPE
    UNKNOWNFUTUREVALUE_FEATURETARGETTYPE
)

func (i FeatureTargetType) String() string {
    return []string{"group", "administrativeUnit", "role", "unknownFutureValue"}[i]
}
func ParseFeatureTargetType(v string) (interface{}, error) {
    result := GROUP_FEATURETARGETTYPE
    switch v {
        case "group":
            result = GROUP_FEATURETARGETTYPE
        case "administrativeUnit":
            result = ADMINISTRATIVEUNIT_FEATURETARGETTYPE
        case "role":
            result = ROLE_FEATURETARGETTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_FEATURETARGETTYPE
        default:
            return 0, errors.New("Unknown FeatureTargetType value: " + v)
    }
    return &result, nil
}
func SerializeFeatureTargetType(values []FeatureTargetType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
