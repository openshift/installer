package models
import (
    "errors"
)
// Provides operations to manage the collection of agreement entities.
type EdgeCookiePolicy int

const (
    // Allow the user to set.
    USERDEFINED_EDGECOOKIEPOLICY EdgeCookiePolicy = iota
    // Allow.
    ALLOW_EDGECOOKIEPOLICY
    // Block only third party cookies.
    BLOCKTHIRDPARTY_EDGECOOKIEPOLICY
    // Block all cookies.
    BLOCKALL_EDGECOOKIEPOLICY
)

func (i EdgeCookiePolicy) String() string {
    return []string{"userDefined", "allow", "blockThirdParty", "blockAll"}[i]
}
func ParseEdgeCookiePolicy(v string) (interface{}, error) {
    result := USERDEFINED_EDGECOOKIEPOLICY
    switch v {
        case "userDefined":
            result = USERDEFINED_EDGECOOKIEPOLICY
        case "allow":
            result = ALLOW_EDGECOOKIEPOLICY
        case "blockThirdParty":
            result = BLOCKTHIRDPARTY_EDGECOOKIEPOLICY
        case "blockAll":
            result = BLOCKALL_EDGECOOKIEPOLICY
        default:
            return 0, errors.New("Unknown EdgeCookiePolicy value: " + v)
    }
    return &result, nil
}
func SerializeEdgeCookiePolicy(values []EdgeCookiePolicy) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
