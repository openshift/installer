package models
import (
    "errors"
)
// Provides operations to manage the collection of agreement entities.
type NotificationTemplateBrandingOptions int

const (
    // No Branding.
    NONE_NOTIFICATIONTEMPLATEBRANDINGOPTIONS NotificationTemplateBrandingOptions = iota
    // Include Company Logo.
    INCLUDECOMPANYLOGO_NOTIFICATIONTEMPLATEBRANDINGOPTIONS
    // Include Company Name.
    INCLUDECOMPANYNAME_NOTIFICATIONTEMPLATEBRANDINGOPTIONS
    // Include Contact Info.
    INCLUDECONTACTINFORMATION_NOTIFICATIONTEMPLATEBRANDINGOPTIONS
    // Include Device Details.
    INCLUDEDEVICEDETAILS_NOTIFICATIONTEMPLATEBRANDINGOPTIONS
)

func (i NotificationTemplateBrandingOptions) String() string {
    return []string{"none", "includeCompanyLogo", "includeCompanyName", "includeContactInformation", "includeDeviceDetails"}[i]
}
func ParseNotificationTemplateBrandingOptions(v string) (interface{}, error) {
    result := NONE_NOTIFICATIONTEMPLATEBRANDINGOPTIONS
    switch v {
        case "none":
            result = NONE_NOTIFICATIONTEMPLATEBRANDINGOPTIONS
        case "includeCompanyLogo":
            result = INCLUDECOMPANYLOGO_NOTIFICATIONTEMPLATEBRANDINGOPTIONS
        case "includeCompanyName":
            result = INCLUDECOMPANYNAME_NOTIFICATIONTEMPLATEBRANDINGOPTIONS
        case "includeContactInformation":
            result = INCLUDECONTACTINFORMATION_NOTIFICATIONTEMPLATEBRANDINGOPTIONS
        case "includeDeviceDetails":
            result = INCLUDEDEVICEDETAILS_NOTIFICATIONTEMPLATEBRANDINGOPTIONS
        default:
            return 0, errors.New("Unknown NotificationTemplateBrandingOptions value: " + v)
    }
    return &result, nil
}
func SerializeNotificationTemplateBrandingOptions(values []NotificationTemplateBrandingOptions) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
