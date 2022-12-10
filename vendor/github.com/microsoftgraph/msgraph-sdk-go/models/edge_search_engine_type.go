package models
import (
    "errors"
)
// Provides operations to manage the collection of agreement entities.
type EdgeSearchEngineType int

const (
    // Uses factory settings of Edge to assign the default search engine as per the user market
    DEFAULT_ESCAPED_EDGESEARCHENGINETYPE EdgeSearchEngineType = iota
    // Sets Bing as the default search engine
    BING_EDGESEARCHENGINETYPE
)

func (i EdgeSearchEngineType) String() string {
    return []string{"default", "bing"}[i]
}
func ParseEdgeSearchEngineType(v string) (interface{}, error) {
    result := DEFAULT_ESCAPED_EDGESEARCHENGINETYPE
    switch v {
        case "default":
            result = DEFAULT_ESCAPED_EDGESEARCHENGINETYPE
        case "bing":
            result = BING_EDGESEARCHENGINETYPE
        default:
            return 0, errors.New("Unknown EdgeSearchEngineType value: " + v)
    }
    return &result, nil
}
func SerializeEdgeSearchEngineType(values []EdgeSearchEngineType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
