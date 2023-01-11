package models
import (
    "errors"
)
// Provides operations to manage the collection of dataPolicyOperation entities.
type DataPolicyOperationStatus int

const (
    NOTSTARTED_DATAPOLICYOPERATIONSTATUS DataPolicyOperationStatus = iota
    RUNNING_DATAPOLICYOPERATIONSTATUS
    COMPLETE_DATAPOLICYOPERATIONSTATUS
    FAILED_DATAPOLICYOPERATIONSTATUS
    UNKNOWNFUTUREVALUE_DATAPOLICYOPERATIONSTATUS
)

func (i DataPolicyOperationStatus) String() string {
    return []string{"notStarted", "running", "complete", "failed", "unknownFutureValue"}[i]
}
func ParseDataPolicyOperationStatus(v string) (interface{}, error) {
    result := NOTSTARTED_DATAPOLICYOPERATIONSTATUS
    switch v {
        case "notStarted":
            result = NOTSTARTED_DATAPOLICYOPERATIONSTATUS
        case "running":
            result = RUNNING_DATAPOLICYOPERATIONSTATUS
        case "complete":
            result = COMPLETE_DATAPOLICYOPERATIONSTATUS
        case "failed":
            result = FAILED_DATAPOLICYOPERATIONSTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_DATAPOLICYOPERATIONSTATUS
        default:
            return 0, errors.New("Unknown DataPolicyOperationStatus value: " + v)
    }
    return &result, nil
}
func SerializeDataPolicyOperationStatus(values []DataPolicyOperationStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
