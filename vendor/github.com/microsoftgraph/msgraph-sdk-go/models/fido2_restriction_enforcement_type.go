package models
import (
    "errors"
)
// Provides operations to manage the collection of authenticationMethodConfiguration entities.
type Fido2RestrictionEnforcementType int

const (
    ALLOW_FIDO2RESTRICTIONENFORCEMENTTYPE Fido2RestrictionEnforcementType = iota
    BLOCK_FIDO2RESTRICTIONENFORCEMENTTYPE
    UNKNOWNFUTUREVALUE_FIDO2RESTRICTIONENFORCEMENTTYPE
)

func (i Fido2RestrictionEnforcementType) String() string {
    return []string{"allow", "block", "unknownFutureValue"}[i]
}
func ParseFido2RestrictionEnforcementType(v string) (interface{}, error) {
    result := ALLOW_FIDO2RESTRICTIONENFORCEMENTTYPE
    switch v {
        case "allow":
            result = ALLOW_FIDO2RESTRICTIONENFORCEMENTTYPE
        case "block":
            result = BLOCK_FIDO2RESTRICTIONENFORCEMENTTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_FIDO2RESTRICTIONENFORCEMENTTYPE
        default:
            return 0, errors.New("Unknown Fido2RestrictionEnforcementType value: " + v)
    }
    return &result, nil
}
func SerializeFido2RestrictionEnforcementType(values []Fido2RestrictionEnforcementType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
