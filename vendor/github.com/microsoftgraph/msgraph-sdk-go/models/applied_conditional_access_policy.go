package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AppliedConditionalAccessPolicy 
type AppliedConditionalAccessPolicy struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Refers to the Name of the conditional access policy (example: 'Require MFA for Salesforce').
    displayName *string
    // Refers to the grant controls enforced by the conditional access policy (example: 'Require multi-factor authentication').
    enforcedGrantControls []string
    // Refers to the session controls enforced by the conditional access policy (example: 'Require app enforced controls').
    enforcedSessionControls []string
    // An identifier of the conditional access policy.
    id *string
    // The OdataType property
    odataType *string
    // Indicates the result of the CA policy that was triggered. Possible values are: success, failure, notApplied (Policy isn't applied because policy conditions were not met),notEnabled (This is due to the policy in disabled state), unknown, unknownFutureValue.
    result *AppliedConditionalAccessPolicyResult
}
// NewAppliedConditionalAccessPolicy instantiates a new appliedConditionalAccessPolicy and sets the default values.
func NewAppliedConditionalAccessPolicy()(*AppliedConditionalAccessPolicy) {
    m := &AppliedConditionalAccessPolicy{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateAppliedConditionalAccessPolicyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAppliedConditionalAccessPolicyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAppliedConditionalAccessPolicy(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *AppliedConditionalAccessPolicy) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetDisplayName gets the displayName property value. Refers to the Name of the conditional access policy (example: 'Require MFA for Salesforce').
func (m *AppliedConditionalAccessPolicy) GetDisplayName()(*string) {
    return m.displayName
}
// GetEnforcedGrantControls gets the enforcedGrantControls property value. Refers to the grant controls enforced by the conditional access policy (example: 'Require multi-factor authentication').
func (m *AppliedConditionalAccessPolicy) GetEnforcedGrantControls()([]string) {
    return m.enforcedGrantControls
}
// GetEnforcedSessionControls gets the enforcedSessionControls property value. Refers to the session controls enforced by the conditional access policy (example: 'Require app enforced controls').
func (m *AppliedConditionalAccessPolicy) GetEnforcedSessionControls()([]string) {
    return m.enforcedSessionControls
}
// GetFieldDeserializers the deserialization information for the current model
func (m *AppliedConditionalAccessPolicy) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["displayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDisplayName)
    res["enforcedGrantControls"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("string" , m.SetEnforcedGrantControls)
    res["enforcedSessionControls"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("string" , m.SetEnforcedSessionControls)
    res["id"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetId)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["result"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseAppliedConditionalAccessPolicyResult , m.SetResult)
    return res
}
// GetId gets the id property value. An identifier of the conditional access policy.
func (m *AppliedConditionalAccessPolicy) GetId()(*string) {
    return m.id
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *AppliedConditionalAccessPolicy) GetOdataType()(*string) {
    return m.odataType
}
// GetResult gets the result property value. Indicates the result of the CA policy that was triggered. Possible values are: success, failure, notApplied (Policy isn't applied because policy conditions were not met),notEnabled (This is due to the policy in disabled state), unknown, unknownFutureValue.
func (m *AppliedConditionalAccessPolicy) GetResult()(*AppliedConditionalAccessPolicyResult) {
    return m.result
}
// Serialize serializes information the current object
func (m *AppliedConditionalAccessPolicy) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    if m.GetEnforcedGrantControls() != nil {
        err := writer.WriteCollectionOfStringValues("enforcedGrantControls", m.GetEnforcedGrantControls())
        if err != nil {
            return err
        }
    }
    if m.GetEnforcedSessionControls() != nil {
        err := writer.WriteCollectionOfStringValues("enforcedSessionControls", m.GetEnforcedSessionControls())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("id", m.GetId())
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
    if m.GetResult() != nil {
        cast := (*m.GetResult()).String()
        err := writer.WriteStringValue("result", &cast)
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
func (m *AppliedConditionalAccessPolicy) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetDisplayName sets the displayName property value. Refers to the Name of the conditional access policy (example: 'Require MFA for Salesforce').
func (m *AppliedConditionalAccessPolicy) SetDisplayName(value *string)() {
    m.displayName = value
}
// SetEnforcedGrantControls sets the enforcedGrantControls property value. Refers to the grant controls enforced by the conditional access policy (example: 'Require multi-factor authentication').
func (m *AppliedConditionalAccessPolicy) SetEnforcedGrantControls(value []string)() {
    m.enforcedGrantControls = value
}
// SetEnforcedSessionControls sets the enforcedSessionControls property value. Refers to the session controls enforced by the conditional access policy (example: 'Require app enforced controls').
func (m *AppliedConditionalAccessPolicy) SetEnforcedSessionControls(value []string)() {
    m.enforcedSessionControls = value
}
// SetId sets the id property value. An identifier of the conditional access policy.
func (m *AppliedConditionalAccessPolicy) SetId(value *string)() {
    m.id = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *AppliedConditionalAccessPolicy) SetOdataType(value *string)() {
    m.odataType = value
}
// SetResult sets the result property value. Indicates the result of the CA policy that was triggered. Possible values are: success, failure, notApplied (Policy isn't applied because policy conditions were not met),notEnabled (This is due to the policy in disabled state), unknown, unknownFutureValue.
func (m *AppliedConditionalAccessPolicy) SetResult(value *AppliedConditionalAccessPolicyResult)() {
    m.result = value
}
