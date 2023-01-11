package logteleconferencedevicequality

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
)

// LogTeleconferenceDeviceQualityPostRequestBody provides operations to call the logTeleconferenceDeviceQuality method.
type LogTeleconferenceDeviceQualityPostRequestBody struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The quality property
    quality iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.TeleconferenceDeviceQualityable
}
// NewLogTeleconferenceDeviceQualityPostRequestBody instantiates a new logTeleconferenceDeviceQualityPostRequestBody and sets the default values.
func NewLogTeleconferenceDeviceQualityPostRequestBody()(*LogTeleconferenceDeviceQualityPostRequestBody) {
    m := &LogTeleconferenceDeviceQualityPostRequestBody{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateLogTeleconferenceDeviceQualityPostRequestBodyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateLogTeleconferenceDeviceQualityPostRequestBodyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewLogTeleconferenceDeviceQualityPostRequestBody(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *LogTeleconferenceDeviceQualityPostRequestBody) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
func (m *LogTeleconferenceDeviceQualityPostRequestBody) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["quality"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateTeleconferenceDeviceQualityFromDiscriminatorValue , m.SetQuality)
    return res
}
// GetQuality gets the quality property value. The quality property
func (m *LogTeleconferenceDeviceQualityPostRequestBody) GetQuality()(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.TeleconferenceDeviceQualityable) {
    return m.quality
}
// Serialize serializes information the current object
func (m *LogTeleconferenceDeviceQualityPostRequestBody) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteObjectValue("quality", m.GetQuality())
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
func (m *LogTeleconferenceDeviceQualityPostRequestBody) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetQuality sets the quality property value. The quality property
func (m *LogTeleconferenceDeviceQualityPostRequestBody) SetQuality(value iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.TeleconferenceDeviceQualityable)() {
    m.quality = value
}
