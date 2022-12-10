package models
import (
    "errors"
)
// Provides operations to manage the collection of agreement entities.
type PrereleaseFeatures int

const (
    // User Defined, default value, no intent.
    USERDEFINED_PRERELEASEFEATURES PrereleaseFeatures = iota
    // Settings only pre-release features.
    SETTINGSONLY_PRERELEASEFEATURES
    // Settings and experimentations pre-release features.
    SETTINGSANDEXPERIMENTATIONS_PRERELEASEFEATURES
    // Pre-release features not allowed.
    NOTALLOWED_PRERELEASEFEATURES
)

func (i PrereleaseFeatures) String() string {
    return []string{"userDefined", "settingsOnly", "settingsAndExperimentations", "notAllowed"}[i]
}
func ParsePrereleaseFeatures(v string) (interface{}, error) {
    result := USERDEFINED_PRERELEASEFEATURES
    switch v {
        case "userDefined":
            result = USERDEFINED_PRERELEASEFEATURES
        case "settingsOnly":
            result = SETTINGSONLY_PRERELEASEFEATURES
        case "settingsAndExperimentations":
            result = SETTINGSANDEXPERIMENTATIONS_PRERELEASEFEATURES
        case "notAllowed":
            result = NOTALLOWED_PRERELEASEFEATURES
        default:
            return 0, errors.New("Unknown PrereleaseFeatures value: " + v)
    }
    return &result, nil
}
func SerializePrereleaseFeatures(values []PrereleaseFeatures) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
