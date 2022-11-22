package models
import (
    "errors"
)
// Provides operations to manage the collection of agreement entities.
type CountryLookupMethodType int

const (
    CLIENTIPADDRESS_COUNTRYLOOKUPMETHODTYPE CountryLookupMethodType = iota
    AUTHENTICATORAPPGPS_COUNTRYLOOKUPMETHODTYPE
    UNKNOWNFUTUREVALUE_COUNTRYLOOKUPMETHODTYPE
)

func (i CountryLookupMethodType) String() string {
    return []string{"clientIpAddress", "authenticatorAppGps", "unknownFutureValue"}[i]
}
func ParseCountryLookupMethodType(v string) (interface{}, error) {
    result := CLIENTIPADDRESS_COUNTRYLOOKUPMETHODTYPE
    switch v {
        case "clientIpAddress":
            result = CLIENTIPADDRESS_COUNTRYLOOKUPMETHODTYPE
        case "authenticatorAppGps":
            result = AUTHENTICATORAPPGPS_COUNTRYLOOKUPMETHODTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_COUNTRYLOOKUPMETHODTYPE
        default:
            return 0, errors.New("Unknown CountryLookupMethodType value: " + v)
    }
    return &result, nil
}
func SerializeCountryLookupMethodType(values []CountryLookupMethodType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
