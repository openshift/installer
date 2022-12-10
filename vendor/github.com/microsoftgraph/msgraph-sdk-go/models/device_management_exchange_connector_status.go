package models
import (
    "errors"
)
// Provides operations to manage the collection of agreement entities.
type DeviceManagementExchangeConnectorStatus int

const (
    // No Connector exists.
    NONE_DEVICEMANAGEMENTEXCHANGECONNECTORSTATUS DeviceManagementExchangeConnectorStatus = iota
    // Pending Connection to the Exchange Environment.
    CONNECTIONPENDING_DEVICEMANAGEMENTEXCHANGECONNECTORSTATUS
    // Connected to the Exchange Environment
    CONNECTED_DEVICEMANAGEMENTEXCHANGECONNECTORSTATUS
    // Disconnected from the Exchange Environment
    DISCONNECTED_DEVICEMANAGEMENTEXCHANGECONNECTORSTATUS
)

func (i DeviceManagementExchangeConnectorStatus) String() string {
    return []string{"none", "connectionPending", "connected", "disconnected"}[i]
}
func ParseDeviceManagementExchangeConnectorStatus(v string) (interface{}, error) {
    result := NONE_DEVICEMANAGEMENTEXCHANGECONNECTORSTATUS
    switch v {
        case "none":
            result = NONE_DEVICEMANAGEMENTEXCHANGECONNECTORSTATUS
        case "connectionPending":
            result = CONNECTIONPENDING_DEVICEMANAGEMENTEXCHANGECONNECTORSTATUS
        case "connected":
            result = CONNECTED_DEVICEMANAGEMENTEXCHANGECONNECTORSTATUS
        case "disconnected":
            result = DISCONNECTED_DEVICEMANAGEMENTEXCHANGECONNECTORSTATUS
        default:
            return 0, errors.New("Unknown DeviceManagementExchangeConnectorStatus value: " + v)
    }
    return &result, nil
}
func SerializeDeviceManagementExchangeConnectorStatus(values []DeviceManagementExchangeConnectorStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
