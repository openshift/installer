package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// TeamsApp provides operations to manage the appCatalogs singleton.
type TeamsApp struct {
    Entity
    // The details for each version of the app.
    appDefinitions []TeamsAppDefinitionable
    // The name of the catalog app provided by the app developer in the Microsoft Teams zip app package.
    displayName *string
    // The method of distribution for the app. Read-only.
    distributionMethod *TeamsAppDistributionMethod
    // The ID of the catalog provided by the app developer in the Microsoft Teams zip app package.
    externalId *string
}
// NewTeamsApp instantiates a new teamsApp and sets the default values.
func NewTeamsApp()(*TeamsApp) {
    m := &TeamsApp{
        Entity: *NewEntity(),
    }
    return m
}
// CreateTeamsAppFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateTeamsAppFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTeamsApp(), nil
}
// GetAppDefinitions gets the appDefinitions property value. The details for each version of the app.
func (m *TeamsApp) GetAppDefinitions()([]TeamsAppDefinitionable) {
    return m.appDefinitions
}
// GetDisplayName gets the displayName property value. The name of the catalog app provided by the app developer in the Microsoft Teams zip app package.
func (m *TeamsApp) GetDisplayName()(*string) {
    return m.displayName
}
// GetDistributionMethod gets the distributionMethod property value. The method of distribution for the app. Read-only.
func (m *TeamsApp) GetDistributionMethod()(*TeamsAppDistributionMethod) {
    return m.distributionMethod
}
// GetExternalId gets the externalId property value. The ID of the catalog provided by the app developer in the Microsoft Teams zip app package.
func (m *TeamsApp) GetExternalId()(*string) {
    return m.externalId
}
// GetFieldDeserializers the deserialization information for the current model
func (m *TeamsApp) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["appDefinitions"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateTeamsAppDefinitionFromDiscriminatorValue , m.SetAppDefinitions)
    res["displayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDisplayName)
    res["distributionMethod"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseTeamsAppDistributionMethod , m.SetDistributionMethod)
    res["externalId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetExternalId)
    return res
}
// Serialize serializes information the current object
func (m *TeamsApp) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAppDefinitions() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetAppDefinitions())
        err = writer.WriteCollectionOfObjectValues("appDefinitions", cast)
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
    if m.GetDistributionMethod() != nil {
        cast := (*m.GetDistributionMethod()).String()
        err = writer.WriteStringValue("distributionMethod", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("externalId", m.GetExternalId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAppDefinitions sets the appDefinitions property value. The details for each version of the app.
func (m *TeamsApp) SetAppDefinitions(value []TeamsAppDefinitionable)() {
    m.appDefinitions = value
}
// SetDisplayName sets the displayName property value. The name of the catalog app provided by the app developer in the Microsoft Teams zip app package.
func (m *TeamsApp) SetDisplayName(value *string)() {
    m.displayName = value
}
// SetDistributionMethod sets the distributionMethod property value. The method of distribution for the app. Read-only.
func (m *TeamsApp) SetDistributionMethod(value *TeamsAppDistributionMethod)() {
    m.distributionMethod = value
}
// SetExternalId sets the externalId property value. The ID of the catalog provided by the app developer in the Microsoft Teams zip app package.
func (m *TeamsApp) SetExternalId(value *string)() {
    m.externalId = value
}
