package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

// AccessReviewScheduleSettings 
type AccessReviewScheduleSettings struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewAccessReviewScheduleSettings instantiates a new accessReviewScheduleSettings and sets the default values.
func NewAccessReviewScheduleSettings()(*AccessReviewScheduleSettings) {
    m := &AccessReviewScheduleSettings{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateAccessReviewScheduleSettingsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAccessReviewScheduleSettingsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAccessReviewScheduleSettings(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *AccessReviewScheduleSettings) GetAdditionalData()(map[string]any) {
    val , err :=  m.backingStore.Get("additionalData")
    if err != nil {
        panic(err)
    }
    if val == nil {
        var value = make(map[string]any);
        m.SetAdditionalData(value);
    }
    return val.(map[string]any)
}
// GetApplyActions gets the applyActions property value. Optional field. Describes the  actions to take once a review is complete. There are two types that are currently supported: removeAccessApplyAction (default) and disableAndDeleteUserApplyAction. Field only needs to be specified in the case of disableAndDeleteUserApplyAction.
func (m *AccessReviewScheduleSettings) GetApplyActions()([]AccessReviewApplyActionable) {
    val, err := m.GetBackingStore().Get("applyActions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AccessReviewApplyActionable)
    }
    return nil
}
// GetAutoApplyDecisionsEnabled gets the autoApplyDecisionsEnabled property value. Indicates whether decisions are automatically applied. When set to false, an admin must apply the decisions manually once the reviewer completes the access review. When set to true, decisions are applied automatically after the access review instance duration ends, whether or not the reviewers have responded. Default value is false.
func (m *AccessReviewScheduleSettings) GetAutoApplyDecisionsEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("autoApplyDecisionsEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetBackingStore gets the backingStore property value. Stores model information.
func (m *AccessReviewScheduleSettings) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetDecisionHistoriesForReviewersEnabled gets the decisionHistoriesForReviewersEnabled property value. Indicates whether decisions on previous access review stages are available for reviewers on an accessReviewInstance with multiple subsequent stages. If not provided, the default is disabled (false).
func (m *AccessReviewScheduleSettings) GetDecisionHistoriesForReviewersEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("decisionHistoriesForReviewersEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetDefaultDecision gets the defaultDecision property value. Decision chosen if defaultDecisionEnabled is enabled. Can be one of Approve, Deny, or Recommendation.
func (m *AccessReviewScheduleSettings) GetDefaultDecision()(*string) {
    val, err := m.GetBackingStore().Get("defaultDecision")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDefaultDecisionEnabled gets the defaultDecisionEnabled property value. Indicates whether the default decision is enabled or disabled when reviewers do not respond. Default value is false.
func (m *AccessReviewScheduleSettings) GetDefaultDecisionEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("defaultDecisionEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *AccessReviewScheduleSettings) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["applyActions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAccessReviewApplyActionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AccessReviewApplyActionable, len(val))
            for i, v := range val {
                res[i] = v.(AccessReviewApplyActionable)
            }
            m.SetApplyActions(res)
        }
        return nil
    }
    res["autoApplyDecisionsEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAutoApplyDecisionsEnabled(val)
        }
        return nil
    }
    res["decisionHistoriesForReviewersEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDecisionHistoriesForReviewersEnabled(val)
        }
        return nil
    }
    res["defaultDecision"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefaultDecision(val)
        }
        return nil
    }
    res["defaultDecisionEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefaultDecisionEnabled(val)
        }
        return nil
    }
    res["instanceDurationInDays"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInstanceDurationInDays(val)
        }
        return nil
    }
    res["justificationRequiredOnApproval"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetJustificationRequiredOnApproval(val)
        }
        return nil
    }
    res["mailNotificationsEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMailNotificationsEnabled(val)
        }
        return nil
    }
    res["@odata.type"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOdataType(val)
        }
        return nil
    }
    res["recommendationsEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRecommendationsEnabled(val)
        }
        return nil
    }
    res["recurrence"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePatternedRecurrenceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRecurrence(val.(PatternedRecurrenceable))
        }
        return nil
    }
    res["reminderNotificationsEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetReminderNotificationsEnabled(val)
        }
        return nil
    }
    return res
}
// GetInstanceDurationInDays gets the instanceDurationInDays property value. Duration of an access review instance in days. NOTE: If the stageSettings of the accessReviewScheduleDefinition object is defined, its durationInDays setting will be used instead of the value of this property.
func (m *AccessReviewScheduleSettings) GetInstanceDurationInDays()(*int32) {
    val, err := m.GetBackingStore().Get("instanceDurationInDays")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetJustificationRequiredOnApproval gets the justificationRequiredOnApproval property value. Indicates whether reviewers are required to provide justification with their decision. Default value is false.
func (m *AccessReviewScheduleSettings) GetJustificationRequiredOnApproval()(*bool) {
    val, err := m.GetBackingStore().Get("justificationRequiredOnApproval")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetMailNotificationsEnabled gets the mailNotificationsEnabled property value. Indicates whether emails are enabled or disabled. Default value is false.
func (m *AccessReviewScheduleSettings) GetMailNotificationsEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("mailNotificationsEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *AccessReviewScheduleSettings) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRecommendationsEnabled gets the recommendationsEnabled property value. Indicates whether decision recommendations are enabled or disabled. NOTE: If the stageSettings of the accessReviewScheduleDefinition object is defined, its recommendationsEnabled setting will be used instead of the value of this property.
func (m *AccessReviewScheduleSettings) GetRecommendationsEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("recommendationsEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetRecurrence gets the recurrence property value. Detailed settings for recurrence using the standard Outlook recurrence object. Note: Only dayOfMonth, interval, and type (weekly, absoluteMonthly) properties are supported. Use the property startDate on recurrenceRange to determine the day the review starts.
func (m *AccessReviewScheduleSettings) GetRecurrence()(PatternedRecurrenceable) {
    val, err := m.GetBackingStore().Get("recurrence")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PatternedRecurrenceable)
    }
    return nil
}
// GetReminderNotificationsEnabled gets the reminderNotificationsEnabled property value. Indicates whether reminders are enabled or disabled. Default value is false.
func (m *AccessReviewScheduleSettings) GetReminderNotificationsEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("reminderNotificationsEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AccessReviewScheduleSettings) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetApplyActions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetApplyActions()))
        for i, v := range m.GetApplyActions() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err := writer.WriteCollectionOfObjectValues("applyActions", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("autoApplyDecisionsEnabled", m.GetAutoApplyDecisionsEnabled())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("decisionHistoriesForReviewersEnabled", m.GetDecisionHistoriesForReviewersEnabled())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("defaultDecision", m.GetDefaultDecision())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("defaultDecisionEnabled", m.GetDefaultDecisionEnabled())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("instanceDurationInDays", m.GetInstanceDurationInDays())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("justificationRequiredOnApproval", m.GetJustificationRequiredOnApproval())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("mailNotificationsEnabled", m.GetMailNotificationsEnabled())
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
        err := writer.WriteBoolValue("recommendationsEnabled", m.GetRecommendationsEnabled())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("recurrence", m.GetRecurrence())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("reminderNotificationsEnabled", m.GetReminderNotificationsEnabled())
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
func (m *AccessReviewScheduleSettings) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetApplyActions sets the applyActions property value. Optional field. Describes the  actions to take once a review is complete. There are two types that are currently supported: removeAccessApplyAction (default) and disableAndDeleteUserApplyAction. Field only needs to be specified in the case of disableAndDeleteUserApplyAction.
func (m *AccessReviewScheduleSettings) SetApplyActions(value []AccessReviewApplyActionable)() {
    err := m.GetBackingStore().Set("applyActions", value)
    if err != nil {
        panic(err)
    }
}
// SetAutoApplyDecisionsEnabled sets the autoApplyDecisionsEnabled property value. Indicates whether decisions are automatically applied. When set to false, an admin must apply the decisions manually once the reviewer completes the access review. When set to true, decisions are applied automatically after the access review instance duration ends, whether or not the reviewers have responded. Default value is false.
func (m *AccessReviewScheduleSettings) SetAutoApplyDecisionsEnabled(value *bool)() {
    err := m.GetBackingStore().Set("autoApplyDecisionsEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the backingStore property value. Stores model information.
func (m *AccessReviewScheduleSettings) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetDecisionHistoriesForReviewersEnabled sets the decisionHistoriesForReviewersEnabled property value. Indicates whether decisions on previous access review stages are available for reviewers on an accessReviewInstance with multiple subsequent stages. If not provided, the default is disabled (false).
func (m *AccessReviewScheduleSettings) SetDecisionHistoriesForReviewersEnabled(value *bool)() {
    err := m.GetBackingStore().Set("decisionHistoriesForReviewersEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetDefaultDecision sets the defaultDecision property value. Decision chosen if defaultDecisionEnabled is enabled. Can be one of Approve, Deny, or Recommendation.
func (m *AccessReviewScheduleSettings) SetDefaultDecision(value *string)() {
    err := m.GetBackingStore().Set("defaultDecision", value)
    if err != nil {
        panic(err)
    }
}
// SetDefaultDecisionEnabled sets the defaultDecisionEnabled property value. Indicates whether the default decision is enabled or disabled when reviewers do not respond. Default value is false.
func (m *AccessReviewScheduleSettings) SetDefaultDecisionEnabled(value *bool)() {
    err := m.GetBackingStore().Set("defaultDecisionEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetInstanceDurationInDays sets the instanceDurationInDays property value. Duration of an access review instance in days. NOTE: If the stageSettings of the accessReviewScheduleDefinition object is defined, its durationInDays setting will be used instead of the value of this property.
func (m *AccessReviewScheduleSettings) SetInstanceDurationInDays(value *int32)() {
    err := m.GetBackingStore().Set("instanceDurationInDays", value)
    if err != nil {
        panic(err)
    }
}
// SetJustificationRequiredOnApproval sets the justificationRequiredOnApproval property value. Indicates whether reviewers are required to provide justification with their decision. Default value is false.
func (m *AccessReviewScheduleSettings) SetJustificationRequiredOnApproval(value *bool)() {
    err := m.GetBackingStore().Set("justificationRequiredOnApproval", value)
    if err != nil {
        panic(err)
    }
}
// SetMailNotificationsEnabled sets the mailNotificationsEnabled property value. Indicates whether emails are enabled or disabled. Default value is false.
func (m *AccessReviewScheduleSettings) SetMailNotificationsEnabled(value *bool)() {
    err := m.GetBackingStore().Set("mailNotificationsEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *AccessReviewScheduleSettings) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetRecommendationsEnabled sets the recommendationsEnabled property value. Indicates whether decision recommendations are enabled or disabled. NOTE: If the stageSettings of the accessReviewScheduleDefinition object is defined, its recommendationsEnabled setting will be used instead of the value of this property.
func (m *AccessReviewScheduleSettings) SetRecommendationsEnabled(value *bool)() {
    err := m.GetBackingStore().Set("recommendationsEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetRecurrence sets the recurrence property value. Detailed settings for recurrence using the standard Outlook recurrence object. Note: Only dayOfMonth, interval, and type (weekly, absoluteMonthly) properties are supported. Use the property startDate on recurrenceRange to determine the day the review starts.
func (m *AccessReviewScheduleSettings) SetRecurrence(value PatternedRecurrenceable)() {
    err := m.GetBackingStore().Set("recurrence", value)
    if err != nil {
        panic(err)
    }
}
// SetReminderNotificationsEnabled sets the reminderNotificationsEnabled property value. Indicates whether reminders are enabled or disabled. Default value is false.
func (m *AccessReviewScheduleSettings) SetReminderNotificationsEnabled(value *bool)() {
    err := m.GetBackingStore().Set("reminderNotificationsEnabled", value)
    if err != nil {
        panic(err)
    }
}
// AccessReviewScheduleSettingsable 
type AccessReviewScheduleSettingsable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetApplyActions()([]AccessReviewApplyActionable)
    GetAutoApplyDecisionsEnabled()(*bool)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetDecisionHistoriesForReviewersEnabled()(*bool)
    GetDefaultDecision()(*string)
    GetDefaultDecisionEnabled()(*bool)
    GetInstanceDurationInDays()(*int32)
    GetJustificationRequiredOnApproval()(*bool)
    GetMailNotificationsEnabled()(*bool)
    GetOdataType()(*string)
    GetRecommendationsEnabled()(*bool)
    GetRecurrence()(PatternedRecurrenceable)
    GetReminderNotificationsEnabled()(*bool)
    SetApplyActions(value []AccessReviewApplyActionable)()
    SetAutoApplyDecisionsEnabled(value *bool)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetDecisionHistoriesForReviewersEnabled(value *bool)()
    SetDefaultDecision(value *string)()
    SetDefaultDecisionEnabled(value *bool)()
    SetInstanceDurationInDays(value *int32)()
    SetJustificationRequiredOnApproval(value *bool)()
    SetMailNotificationsEnabled(value *bool)()
    SetOdataType(value *string)()
    SetRecommendationsEnabled(value *bool)()
    SetRecurrence(value PatternedRecurrenceable)()
    SetReminderNotificationsEnabled(value *bool)()
}
