package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// MessageRule provides operations to manage the collection of agreement entities.
type MessageRule struct {
    Entity
    // Actions to be taken on a message when the corresponding conditions are fulfilled.
    actions MessageRuleActionsable
    // Conditions that when fulfilled, will trigger the corresponding actions for that rule.
    conditions MessageRulePredicatesable
    // The display name of the rule.
    displayName *string
    // Exception conditions for the rule.
    exceptions MessageRulePredicatesable
    // Indicates whether the rule is in an error condition. Read-only.
    hasError *bool
    // Indicates whether the rule is enabled to be applied to messages.
    isEnabled *bool
    // Indicates if the rule is read-only and cannot be modified or deleted by the rules REST API.
    isReadOnly *bool
    // Indicates the order in which the rule is executed, among other rules.
    sequence *int32
}
// NewMessageRule instantiates a new messageRule and sets the default values.
func NewMessageRule()(*MessageRule) {
    m := &MessageRule{
        Entity: *NewEntity(),
    }
    return m
}
// CreateMessageRuleFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateMessageRuleFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMessageRule(), nil
}
// GetActions gets the actions property value. Actions to be taken on a message when the corresponding conditions are fulfilled.
func (m *MessageRule) GetActions()(MessageRuleActionsable) {
    return m.actions
}
// GetConditions gets the conditions property value. Conditions that when fulfilled, will trigger the corresponding actions for that rule.
func (m *MessageRule) GetConditions()(MessageRulePredicatesable) {
    return m.conditions
}
// GetDisplayName gets the displayName property value. The display name of the rule.
func (m *MessageRule) GetDisplayName()(*string) {
    return m.displayName
}
// GetExceptions gets the exceptions property value. Exception conditions for the rule.
func (m *MessageRule) GetExceptions()(MessageRulePredicatesable) {
    return m.exceptions
}
// GetFieldDeserializers the deserialization information for the current model
func (m *MessageRule) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["actions"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateMessageRuleActionsFromDiscriminatorValue , m.SetActions)
    res["conditions"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateMessageRulePredicatesFromDiscriminatorValue , m.SetConditions)
    res["displayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDisplayName)
    res["exceptions"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateMessageRulePredicatesFromDiscriminatorValue , m.SetExceptions)
    res["hasError"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetHasError)
    res["isEnabled"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsEnabled)
    res["isReadOnly"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsReadOnly)
    res["sequence"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetSequence)
    return res
}
// GetHasError gets the hasError property value. Indicates whether the rule is in an error condition. Read-only.
func (m *MessageRule) GetHasError()(*bool) {
    return m.hasError
}
// GetIsEnabled gets the isEnabled property value. Indicates whether the rule is enabled to be applied to messages.
func (m *MessageRule) GetIsEnabled()(*bool) {
    return m.isEnabled
}
// GetIsReadOnly gets the isReadOnly property value. Indicates if the rule is read-only and cannot be modified or deleted by the rules REST API.
func (m *MessageRule) GetIsReadOnly()(*bool) {
    return m.isReadOnly
}
// GetSequence gets the sequence property value. Indicates the order in which the rule is executed, among other rules.
func (m *MessageRule) GetSequence()(*int32) {
    return m.sequence
}
// Serialize serializes information the current object
func (m *MessageRule) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("actions", m.GetActions())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("conditions", m.GetConditions())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("exceptions", m.GetExceptions())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("hasError", m.GetHasError())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isEnabled", m.GetIsEnabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isReadOnly", m.GetIsReadOnly())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("sequence", m.GetSequence())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetActions sets the actions property value. Actions to be taken on a message when the corresponding conditions are fulfilled.
func (m *MessageRule) SetActions(value MessageRuleActionsable)() {
    m.actions = value
}
// SetConditions sets the conditions property value. Conditions that when fulfilled, will trigger the corresponding actions for that rule.
func (m *MessageRule) SetConditions(value MessageRulePredicatesable)() {
    m.conditions = value
}
// SetDisplayName sets the displayName property value. The display name of the rule.
func (m *MessageRule) SetDisplayName(value *string)() {
    m.displayName = value
}
// SetExceptions sets the exceptions property value. Exception conditions for the rule.
func (m *MessageRule) SetExceptions(value MessageRulePredicatesable)() {
    m.exceptions = value
}
// SetHasError sets the hasError property value. Indicates whether the rule is in an error condition. Read-only.
func (m *MessageRule) SetHasError(value *bool)() {
    m.hasError = value
}
// SetIsEnabled sets the isEnabled property value. Indicates whether the rule is enabled to be applied to messages.
func (m *MessageRule) SetIsEnabled(value *bool)() {
    m.isEnabled = value
}
// SetIsReadOnly sets the isReadOnly property value. Indicates if the rule is read-only and cannot be modified or deleted by the rules REST API.
func (m *MessageRule) SetIsReadOnly(value *bool)() {
    m.isReadOnly = value
}
// SetSequence sets the sequence property value. Indicates the order in which the rule is executed, among other rules.
func (m *MessageRule) SetSequence(value *int32)() {
    m.sequence = value
}
