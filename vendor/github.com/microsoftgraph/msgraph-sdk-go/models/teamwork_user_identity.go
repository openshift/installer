package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// TeamworkUserIdentity 
type TeamworkUserIdentity struct {
    Identity
    // Type of user. Possible values are: aadUser, onPremiseAadUser, anonymousGuest, federatedUser, personalMicrosoftAccountUser, skypeUser, phoneUser, unknownFutureValue and emailUser.
    userIdentityType *TeamworkUserIdentityType
}
// NewTeamworkUserIdentity instantiates a new TeamworkUserIdentity and sets the default values.
func NewTeamworkUserIdentity()(*TeamworkUserIdentity) {
    m := &TeamworkUserIdentity{
        Identity: *NewIdentity(),
    }
    odataTypeValue := "#microsoft.graph.teamworkUserIdentity";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateTeamworkUserIdentityFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateTeamworkUserIdentityFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTeamworkUserIdentity(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *TeamworkUserIdentity) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Identity.GetFieldDeserializers()
    res["userIdentityType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseTeamworkUserIdentityType , m.SetUserIdentityType)
    return res
}
// GetUserIdentityType gets the userIdentityType property value. Type of user. Possible values are: aadUser, onPremiseAadUser, anonymousGuest, federatedUser, personalMicrosoftAccountUser, skypeUser, phoneUser, unknownFutureValue and emailUser.
func (m *TeamworkUserIdentity) GetUserIdentityType()(*TeamworkUserIdentityType) {
    return m.userIdentityType
}
// Serialize serializes information the current object
func (m *TeamworkUserIdentity) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Identity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetUserIdentityType() != nil {
        cast := (*m.GetUserIdentityType()).String()
        err = writer.WriteStringValue("userIdentityType", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetUserIdentityType sets the userIdentityType property value. Type of user. Possible values are: aadUser, onPremiseAadUser, anonymousGuest, federatedUser, personalMicrosoftAccountUser, skypeUser, phoneUser, unknownFutureValue and emailUser.
func (m *TeamworkUserIdentity) SetUserIdentityType(value *TeamworkUserIdentityType)() {
    m.userIdentityType = value
}
