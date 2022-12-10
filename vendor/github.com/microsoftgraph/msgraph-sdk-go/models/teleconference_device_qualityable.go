package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// TeleconferenceDeviceQualityable 
type TeleconferenceDeviceQualityable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCallChainId()(*string)
    GetCloudServiceDeploymentEnvironment()(*string)
    GetCloudServiceDeploymentId()(*string)
    GetCloudServiceInstanceName()(*string)
    GetCloudServiceName()(*string)
    GetDeviceDescription()(*string)
    GetDeviceName()(*string)
    GetMediaLegId()(*string)
    GetMediaQualityList()([]TeleconferenceDeviceMediaQualityable)
    GetOdataType()(*string)
    GetParticipantId()(*string)
    SetCallChainId(value *string)()
    SetCloudServiceDeploymentEnvironment(value *string)()
    SetCloudServiceDeploymentId(value *string)()
    SetCloudServiceInstanceName(value *string)()
    SetCloudServiceName(value *string)()
    SetDeviceDescription(value *string)()
    SetDeviceName(value *string)()
    SetMediaLegId(value *string)()
    SetMediaQualityList(value []TeleconferenceDeviceMediaQualityable)()
    SetOdataType(value *string)()
    SetParticipantId(value *string)()
}
