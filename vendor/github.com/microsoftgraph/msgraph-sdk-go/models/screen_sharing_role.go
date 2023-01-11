package models
import (
    "errors"
)
// Provides operations to call the changeScreenSharingRole method.
type ScreenSharingRole int

const (
    VIEWER_SCREENSHARINGROLE ScreenSharingRole = iota
    SHARER_SCREENSHARINGROLE
)

func (i ScreenSharingRole) String() string {
    return []string{"viewer", "sharer"}[i]
}
func ParseScreenSharingRole(v string) (interface{}, error) {
    result := VIEWER_SCREENSHARINGROLE
    switch v {
        case "viewer":
            result = VIEWER_SCREENSHARINGROLE
        case "sharer":
            result = SHARER_SCREENSHARINGROLE
        default:
            return 0, errors.New("Unknown ScreenSharingRole value: " + v)
    }
    return &result, nil
}
func SerializeScreenSharingRole(values []ScreenSharingRole) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
