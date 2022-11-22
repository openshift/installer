package security

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// UserSource 
type UserSource struct {
    DataSource
    // Email address of the user's mailbox.
    email *string
    // Specifies which sources are included in this group. Possible values are: mailbox, site.
    includedSources *SourceType
    // The URL of the user's OneDrive for Business site. Read-only.
    siteWebUrl *string
}
// NewUserSource instantiates a new UserSource and sets the default values.
func NewUserSource()(*UserSource) {
    m := &UserSource{
        DataSource: *NewDataSource(),
    }
    odataTypeValue := "#microsoft.graph.security.userSource";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateUserSourceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateUserSourceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUserSource(), nil
}
// GetEmail gets the email property value. Email address of the user's mailbox.
func (m *UserSource) GetEmail()(*string) {
    return m.email
}
// GetFieldDeserializers the deserialization information for the current model
func (m *UserSource) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DataSource.GetFieldDeserializers()
    res["email"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetEmail)
    res["includedSources"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseSourceType , m.SetIncludedSources)
    res["siteWebUrl"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetSiteWebUrl)
    return res
}
// GetIncludedSources gets the includedSources property value. Specifies which sources are included in this group. Possible values are: mailbox, site.
func (m *UserSource) GetIncludedSources()(*SourceType) {
    return m.includedSources
}
// GetSiteWebUrl gets the siteWebUrl property value. The URL of the user's OneDrive for Business site. Read-only.
func (m *UserSource) GetSiteWebUrl()(*string) {
    return m.siteWebUrl
}
// Serialize serializes information the current object
func (m *UserSource) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DataSource.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("email", m.GetEmail())
        if err != nil {
            return err
        }
    }
    if m.GetIncludedSources() != nil {
        cast := (*m.GetIncludedSources()).String()
        err = writer.WriteStringValue("includedSources", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("siteWebUrl", m.GetSiteWebUrl())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetEmail sets the email property value. Email address of the user's mailbox.
func (m *UserSource) SetEmail(value *string)() {
    m.email = value
}
// SetIncludedSources sets the includedSources property value. Specifies which sources are included in this group. Possible values are: mailbox, site.
func (m *UserSource) SetIncludedSources(value *SourceType)() {
    m.includedSources = value
}
// SetSiteWebUrl sets the siteWebUrl property value. The URL of the user's OneDrive for Business site. Read-only.
func (m *UserSource) SetSiteWebUrl(value *string)() {
    m.siteWebUrl = value
}
