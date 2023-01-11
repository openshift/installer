package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ServiceHealth provides operations to manage the admin singleton.
type ServiceHealth struct {
    Entity
    // A collection of issues that happened on the service, with detailed information for each issue.
    issues []ServiceHealthIssueable
    // The service name. Use the list healthOverviews operation to get exact string names for services subscribed by the tenant.
    service *string
    // The status property
    status *ServiceHealthStatus
}
// NewServiceHealth instantiates a new serviceHealth and sets the default values.
func NewServiceHealth()(*ServiceHealth) {
    m := &ServiceHealth{
        Entity: *NewEntity(),
    }
    return m
}
// CreateServiceHealthFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateServiceHealthFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewServiceHealth(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ServiceHealth) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["issues"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateServiceHealthIssueFromDiscriminatorValue , m.SetIssues)
    res["service"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetService)
    res["status"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseServiceHealthStatus , m.SetStatus)
    return res
}
// GetIssues gets the issues property value. A collection of issues that happened on the service, with detailed information for each issue.
func (m *ServiceHealth) GetIssues()([]ServiceHealthIssueable) {
    return m.issues
}
// GetService gets the service property value. The service name. Use the list healthOverviews operation to get exact string names for services subscribed by the tenant.
func (m *ServiceHealth) GetService()(*string) {
    return m.service
}
// GetStatus gets the status property value. The status property
func (m *ServiceHealth) GetStatus()(*ServiceHealthStatus) {
    return m.status
}
// Serialize serializes information the current object
func (m *ServiceHealth) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetIssues() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetIssues())
        err = writer.WriteCollectionOfObjectValues("issues", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("service", m.GetService())
        if err != nil {
            return err
        }
    }
    if m.GetStatus() != nil {
        cast := (*m.GetStatus()).String()
        err = writer.WriteStringValue("status", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetIssues sets the issues property value. A collection of issues that happened on the service, with detailed information for each issue.
func (m *ServiceHealth) SetIssues(value []ServiceHealthIssueable)() {
    m.issues = value
}
// SetService sets the service property value. The service name. Use the list healthOverviews operation to get exact string names for services subscribed by the tenant.
func (m *ServiceHealth) SetService(value *string)() {
    m.service = value
}
// SetStatus sets the status property value. The status property
func (m *ServiceHealth) SetStatus(value *ServiceHealthStatus)() {
    m.status = value
}
