package models
import (
    "errors"
)
// 
type CertificateStatus int

const (
    NOTPROVISIONED_CERTIFICATESTATUS CertificateStatus = iota
    PROVISIONED_CERTIFICATESTATUS
)

func (i CertificateStatus) String() string {
    return []string{"notProvisioned", "provisioned"}[i]
}
func ParseCertificateStatus(v string) (any, error) {
    result := NOTPROVISIONED_CERTIFICATESTATUS
    switch v {
        case "notProvisioned":
            result = NOTPROVISIONED_CERTIFICATESTATUS
        case "provisioned":
            result = PROVISIONED_CERTIFICATESTATUS
        default:
            return 0, errors.New("Unknown CertificateStatus value: " + v)
    }
    return &result, nil
}
func SerializeCertificateStatus(values []CertificateStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
