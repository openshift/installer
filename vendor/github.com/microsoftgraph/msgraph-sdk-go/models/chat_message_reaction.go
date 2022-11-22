package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ChatMessageReaction 
type ChatMessageReaction struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
    createdDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // The OdataType property
    odataType *string
    // Supported values are like, angry, sad, laugh, heart, surprised.
    reactionType *string
    // The user property
    user ChatMessageReactionIdentitySetable
}
// NewChatMessageReaction instantiates a new chatMessageReaction and sets the default values.
func NewChatMessageReaction()(*ChatMessageReaction) {
    m := &ChatMessageReaction{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateChatMessageReactionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateChatMessageReactionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewChatMessageReaction(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *ChatMessageReaction) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetCreatedDateTime gets the createdDateTime property value. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
func (m *ChatMessageReaction) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.createdDateTime
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ChatMessageReaction) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["createdDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetCreatedDateTime)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["reactionType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetReactionType)
    res["user"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateChatMessageReactionIdentitySetFromDiscriminatorValue , m.SetUser)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *ChatMessageReaction) GetOdataType()(*string) {
    return m.odataType
}
// GetReactionType gets the reactionType property value. Supported values are like, angry, sad, laugh, heart, surprised.
func (m *ChatMessageReaction) GetReactionType()(*string) {
    return m.reactionType
}
// GetUser gets the user property value. The user property
func (m *ChatMessageReaction) GetUser()(ChatMessageReactionIdentitySetable) {
    return m.user
}
// Serialize serializes information the current object
func (m *ChatMessageReaction) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteTimeValue("createdDateTime", m.GetCreatedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("reactionType", m.GetReactionType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("user", m.GetUser())
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
func (m *ChatMessageReaction) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetCreatedDateTime sets the createdDateTime property value. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
func (m *ChatMessageReaction) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.createdDateTime = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *ChatMessageReaction) SetOdataType(value *string)() {
    m.odataType = value
}
// SetReactionType sets the reactionType property value. Supported values are like, angry, sad, laugh, heart, surprised.
func (m *ChatMessageReaction) SetReactionType(value *string)() {
    m.reactionType = value
}
// SetUser sets the user property value. The user property
func (m *ChatMessageReaction) SetUser(value ChatMessageReactionIdentitySetable)() {
    m.user = value
}
