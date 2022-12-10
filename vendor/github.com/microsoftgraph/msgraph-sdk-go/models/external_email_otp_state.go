package models
import (
    "errors"
)
// Provides operations to manage the collection of authenticationMethodConfiguration entities.
type ExternalEmailOtpState int

const (
    DEFAULT_ESCAPED_EXTERNALEMAILOTPSTATE ExternalEmailOtpState = iota
    ENABLED_EXTERNALEMAILOTPSTATE
    DISABLED_EXTERNALEMAILOTPSTATE
    UNKNOWNFUTUREVALUE_EXTERNALEMAILOTPSTATE
)

func (i ExternalEmailOtpState) String() string {
    return []string{"default", "enabled", "disabled", "unknownFutureValue"}[i]
}
func ParseExternalEmailOtpState(v string) (interface{}, error) {
    result := DEFAULT_ESCAPED_EXTERNALEMAILOTPSTATE
    switch v {
        case "default":
            result = DEFAULT_ESCAPED_EXTERNALEMAILOTPSTATE
        case "enabled":
            result = ENABLED_EXTERNALEMAILOTPSTATE
        case "disabled":
            result = DISABLED_EXTERNALEMAILOTPSTATE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_EXTERNALEMAILOTPSTATE
        default:
            return 0, errors.New("Unknown ExternalEmailOtpState value: " + v)
    }
    return &result, nil
}
func SerializeExternalEmailOtpState(values []ExternalEmailOtpState) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
