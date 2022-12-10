package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// CallOptions 
type CallOptions struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Indicates whether to hide the app after the call is escalated.
    hideBotAfterEscalation *bool
    // Indicates whether content sharing notifications should be enabled for the call.
    isContentSharingNotificationEnabled *bool
    // The OdataType property
    odataType *string
}
// NewCallOptions instantiates a new callOptions and sets the default values.
func NewCallOptions()(*CallOptions) {
    m := &CallOptions{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateCallOptionsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateCallOptionsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    if parseNode != nil {
        mappingValueNode, err := parseNode.GetChildNode("@odata.type")
        if err != nil {
            return nil, err
        }
        if mappingValueNode != nil {
            mappingValue, err := mappingValueNode.GetStringValue()
            if err != nil {
                return nil, err
            }
            if mappingValue != nil {
                switch *mappingValue {
                    case "#microsoft.graph.incomingCallOptions":
                        return NewIncomingCallOptions(), nil
                    case "#microsoft.graph.outgoingCallOptions":
                        return NewOutgoingCallOptions(), nil
                }
            }
        }
    }
    return NewCallOptions(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *CallOptions) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
func (m *CallOptions) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["hideBotAfterEscalation"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetHideBotAfterEscalation)
    res["isContentSharingNotificationEnabled"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsContentSharingNotificationEnabled)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetHideBotAfterEscalation gets the hideBotAfterEscalation property value. Indicates whether to hide the app after the call is escalated.
func (m *CallOptions) GetHideBotAfterEscalation()(*bool) {
    return m.hideBotAfterEscalation
}
// GetIsContentSharingNotificationEnabled gets the isContentSharingNotificationEnabled property value. Indicates whether content sharing notifications should be enabled for the call.
func (m *CallOptions) GetIsContentSharingNotificationEnabled()(*bool) {
    return m.isContentSharingNotificationEnabled
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *CallOptions) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *CallOptions) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteBoolValue("hideBotAfterEscalation", m.GetHideBotAfterEscalation())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isContentSharingNotificationEnabled", m.GetIsContentSharingNotificationEnabled())
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
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *CallOptions) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetHideBotAfterEscalation sets the hideBotAfterEscalation property value. Indicates whether to hide the app after the call is escalated.
func (m *CallOptions) SetHideBotAfterEscalation(value *bool)() {
    m.hideBotAfterEscalation = value
}
// SetIsContentSharingNotificationEnabled sets the isContentSharingNotificationEnabled property value. Indicates whether content sharing notifications should be enabled for the call.
func (m *CallOptions) SetIsContentSharingNotificationEnabled(value *bool)() {
    m.isContentSharingNotificationEnabled = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *CallOptions) SetOdataType(value *string)() {
    m.odataType = value
}
