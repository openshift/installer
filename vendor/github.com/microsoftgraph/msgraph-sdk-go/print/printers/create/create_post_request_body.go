package create

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
)

// CreatePostRequestBody provides operations to call the create method.
type CreatePostRequestBody struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The certificateSigningRequest property
    certificateSigningRequest iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.PrintCertificateSigningRequestable
    // The connectorId property
    connectorId *string
    // The displayName property
    displayName *string
    // The hasPhysicalDevice property
    hasPhysicalDevice *bool
    // The manufacturer property
    manufacturer *string
    // The model property
    model *string
    // The physicalDeviceId property
    physicalDeviceId *string
}
// NewCreatePostRequestBody instantiates a new createPostRequestBody and sets the default values.
func NewCreatePostRequestBody()(*CreatePostRequestBody) {
    m := &CreatePostRequestBody{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateCreatePostRequestBodyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateCreatePostRequestBodyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCreatePostRequestBody(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *CreatePostRequestBody) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetCertificateSigningRequest gets the certificateSigningRequest property value. The certificateSigningRequest property
func (m *CreatePostRequestBody) GetCertificateSigningRequest()(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.PrintCertificateSigningRequestable) {
    return m.certificateSigningRequest
}
// GetConnectorId gets the connectorId property value. The connectorId property
func (m *CreatePostRequestBody) GetConnectorId()(*string) {
    return m.connectorId
}
// GetDisplayName gets the displayName property value. The displayName property
func (m *CreatePostRequestBody) GetDisplayName()(*string) {
    return m.displayName
}
// GetFieldDeserializers the deserialization information for the current model
func (m *CreatePostRequestBody) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["certificateSigningRequest"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreatePrintCertificateSigningRequestFromDiscriminatorValue , m.SetCertificateSigningRequest)
    res["connectorId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetConnectorId)
    res["displayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDisplayName)
    res["hasPhysicalDevice"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetHasPhysicalDevice)
    res["manufacturer"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetManufacturer)
    res["model"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetModel)
    res["physicalDeviceId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetPhysicalDeviceId)
    return res
}
// GetHasPhysicalDevice gets the hasPhysicalDevice property value. The hasPhysicalDevice property
func (m *CreatePostRequestBody) GetHasPhysicalDevice()(*bool) {
    return m.hasPhysicalDevice
}
// GetManufacturer gets the manufacturer property value. The manufacturer property
func (m *CreatePostRequestBody) GetManufacturer()(*string) {
    return m.manufacturer
}
// GetModel gets the model property value. The model property
func (m *CreatePostRequestBody) GetModel()(*string) {
    return m.model
}
// GetPhysicalDeviceId gets the physicalDeviceId property value. The physicalDeviceId property
func (m *CreatePostRequestBody) GetPhysicalDeviceId()(*string) {
    return m.physicalDeviceId
}
// Serialize serializes information the current object
func (m *CreatePostRequestBody) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteObjectValue("certificateSigningRequest", m.GetCertificateSigningRequest())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("connectorId", m.GetConnectorId())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("hasPhysicalDevice", m.GetHasPhysicalDevice())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("manufacturer", m.GetManufacturer())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("model", m.GetModel())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("physicalDeviceId", m.GetPhysicalDeviceId())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *CreatePostRequestBody) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetCertificateSigningRequest sets the certificateSigningRequest property value. The certificateSigningRequest property
func (m *CreatePostRequestBody) SetCertificateSigningRequest(value iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.PrintCertificateSigningRequestable)() {
    m.certificateSigningRequest = value
}
// SetConnectorId sets the connectorId property value. The connectorId property
func (m *CreatePostRequestBody) SetConnectorId(value *string)() {
    m.connectorId = value
}
// SetDisplayName sets the displayName property value. The displayName property
func (m *CreatePostRequestBody) SetDisplayName(value *string)() {
    m.displayName = value
}
// SetHasPhysicalDevice sets the hasPhysicalDevice property value. The hasPhysicalDevice property
func (m *CreatePostRequestBody) SetHasPhysicalDevice(value *bool)() {
    m.hasPhysicalDevice = value
}
// SetManufacturer sets the manufacturer property value. The manufacturer property
func (m *CreatePostRequestBody) SetManufacturer(value *string)() {
    m.manufacturer = value
}
// SetModel sets the model property value. The model property
func (m *CreatePostRequestBody) SetModel(value *string)() {
    m.model = value
}
// SetPhysicalDeviceId sets the physicalDeviceId property value. The physicalDeviceId property
func (m *CreatePostRequestBody) SetPhysicalDeviceId(value *string)() {
    m.physicalDeviceId = value
}
