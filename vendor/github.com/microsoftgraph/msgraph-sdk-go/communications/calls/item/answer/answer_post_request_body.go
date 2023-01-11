package answer

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
)

// AnswerPostRequestBody provides operations to call the answer method.
type AnswerPostRequestBody struct {
    // The acceptedModalities property
    acceptedModalities []iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Modality
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The callbackUri property
    callbackUri *string
    // The callOptions property
    callOptions iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IncomingCallOptionsable
    // The mediaConfig property
    mediaConfig iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MediaConfigable
    // The participantCapacity property
    participantCapacity *int32
}
// NewAnswerPostRequestBody instantiates a new answerPostRequestBody and sets the default values.
func NewAnswerPostRequestBody()(*AnswerPostRequestBody) {
    m := &AnswerPostRequestBody{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateAnswerPostRequestBodyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAnswerPostRequestBodyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAnswerPostRequestBody(), nil
}
// GetAcceptedModalities gets the acceptedModalities property value. The acceptedModalities property
func (m *AnswerPostRequestBody) GetAcceptedModalities()([]iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Modality) {
    return m.acceptedModalities
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *AnswerPostRequestBody) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetCallbackUri gets the callbackUri property value. The callbackUri property
func (m *AnswerPostRequestBody) GetCallbackUri()(*string) {
    return m.callbackUri
}
// GetCallOptions gets the callOptions property value. The callOptions property
func (m *AnswerPostRequestBody) GetCallOptions()(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IncomingCallOptionsable) {
    return m.callOptions
}
// GetFieldDeserializers the deserialization information for the current model
func (m *AnswerPostRequestBody) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["acceptedModalities"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfEnumValues(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ParseModality , m.SetAcceptedModalities)
    res["callbackUri"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetCallbackUri)
    res["callOptions"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateIncomingCallOptionsFromDiscriminatorValue , m.SetCallOptions)
    res["mediaConfig"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateMediaConfigFromDiscriminatorValue , m.SetMediaConfig)
    res["participantCapacity"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetParticipantCapacity)
    return res
}
// GetMediaConfig gets the mediaConfig property value. The mediaConfig property
func (m *AnswerPostRequestBody) GetMediaConfig()(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MediaConfigable) {
    return m.mediaConfig
}
// GetParticipantCapacity gets the participantCapacity property value. The participantCapacity property
func (m *AnswerPostRequestBody) GetParticipantCapacity()(*int32) {
    return m.participantCapacity
}
// Serialize serializes information the current object
func (m *AnswerPostRequestBody) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetAcceptedModalities() != nil {
        err := writer.WriteCollectionOfStringValues("acceptedModalities", iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SerializeModality(m.GetAcceptedModalities()))
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("callbackUri", m.GetCallbackUri())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("callOptions", m.GetCallOptions())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("mediaConfig", m.GetMediaConfig())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("participantCapacity", m.GetParticipantCapacity())
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
// SetAcceptedModalities sets the acceptedModalities property value. The acceptedModalities property
func (m *AnswerPostRequestBody) SetAcceptedModalities(value []iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Modality)() {
    m.acceptedModalities = value
}
// SetAdditionalData sets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *AnswerPostRequestBody) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetCallbackUri sets the callbackUri property value. The callbackUri property
func (m *AnswerPostRequestBody) SetCallbackUri(value *string)() {
    m.callbackUri = value
}
// SetCallOptions sets the callOptions property value. The callOptions property
func (m *AnswerPostRequestBody) SetCallOptions(value iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IncomingCallOptionsable)() {
    m.callOptions = value
}
// SetMediaConfig sets the mediaConfig property value. The mediaConfig property
func (m *AnswerPostRequestBody) SetMediaConfig(value iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MediaConfigable)() {
    m.mediaConfig = value
}
// SetParticipantCapacity sets the participantCapacity property value. The participantCapacity property
func (m *AnswerPostRequestBody) SetParticipantCapacity(value *int32)() {
    m.participantCapacity = value
}
