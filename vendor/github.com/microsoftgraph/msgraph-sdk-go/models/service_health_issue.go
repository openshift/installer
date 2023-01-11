package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ServiceHealthIssue 
type ServiceHealthIssue struct {
    ServiceAnnouncementBase
    // The classification property
    classification *ServiceHealthClassificationType
    // The feature name of the service issue.
    feature *string
    // The feature group name of the service issue.
    featureGroup *string
    // The description of the service issue impact.
    impactDescription *string
    // Indicates whether the issue is resolved.
    isResolved *bool
    // The origin property
    origin *ServiceHealthOrigin
    // Collection of historical posts for the service issue.
    posts []ServiceHealthIssuePostable
    // Indicates the service affected by the issue.
    service *string
    // The status property
    status *ServiceHealthStatus
}
// NewServiceHealthIssue instantiates a new ServiceHealthIssue and sets the default values.
func NewServiceHealthIssue()(*ServiceHealthIssue) {
    m := &ServiceHealthIssue{
        ServiceAnnouncementBase: *NewServiceAnnouncementBase(),
    }
    odataTypeValue := "#microsoft.graph.serviceHealthIssue";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateServiceHealthIssueFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateServiceHealthIssueFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewServiceHealthIssue(), nil
}
// GetClassification gets the classification property value. The classification property
func (m *ServiceHealthIssue) GetClassification()(*ServiceHealthClassificationType) {
    return m.classification
}
// GetFeature gets the feature property value. The feature name of the service issue.
func (m *ServiceHealthIssue) GetFeature()(*string) {
    return m.feature
}
// GetFeatureGroup gets the featureGroup property value. The feature group name of the service issue.
func (m *ServiceHealthIssue) GetFeatureGroup()(*string) {
    return m.featureGroup
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ServiceHealthIssue) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ServiceAnnouncementBase.GetFieldDeserializers()
    res["classification"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseServiceHealthClassificationType , m.SetClassification)
    res["feature"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetFeature)
    res["featureGroup"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetFeatureGroup)
    res["impactDescription"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetImpactDescription)
    res["isResolved"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsResolved)
    res["origin"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseServiceHealthOrigin , m.SetOrigin)
    res["posts"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateServiceHealthIssuePostFromDiscriminatorValue , m.SetPosts)
    res["service"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetService)
    res["status"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseServiceHealthStatus , m.SetStatus)
    return res
}
// GetImpactDescription gets the impactDescription property value. The description of the service issue impact.
func (m *ServiceHealthIssue) GetImpactDescription()(*string) {
    return m.impactDescription
}
// GetIsResolved gets the isResolved property value. Indicates whether the issue is resolved.
func (m *ServiceHealthIssue) GetIsResolved()(*bool) {
    return m.isResolved
}
// GetOrigin gets the origin property value. The origin property
func (m *ServiceHealthIssue) GetOrigin()(*ServiceHealthOrigin) {
    return m.origin
}
// GetPosts gets the posts property value. Collection of historical posts for the service issue.
func (m *ServiceHealthIssue) GetPosts()([]ServiceHealthIssuePostable) {
    return m.posts
}
// GetService gets the service property value. Indicates the service affected by the issue.
func (m *ServiceHealthIssue) GetService()(*string) {
    return m.service
}
// GetStatus gets the status property value. The status property
func (m *ServiceHealthIssue) GetStatus()(*ServiceHealthStatus) {
    return m.status
}
// Serialize serializes information the current object
func (m *ServiceHealthIssue) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ServiceAnnouncementBase.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetClassification() != nil {
        cast := (*m.GetClassification()).String()
        err = writer.WriteStringValue("classification", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("feature", m.GetFeature())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("featureGroup", m.GetFeatureGroup())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("impactDescription", m.GetImpactDescription())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isResolved", m.GetIsResolved())
        if err != nil {
            return err
        }
    }
    if m.GetOrigin() != nil {
        cast := (*m.GetOrigin()).String()
        err = writer.WriteStringValue("origin", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetPosts() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetPosts())
        err = writer.WriteCollectionOfObjectValues("posts", cast)
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
// SetClassification sets the classification property value. The classification property
func (m *ServiceHealthIssue) SetClassification(value *ServiceHealthClassificationType)() {
    m.classification = value
}
// SetFeature sets the feature property value. The feature name of the service issue.
func (m *ServiceHealthIssue) SetFeature(value *string)() {
    m.feature = value
}
// SetFeatureGroup sets the featureGroup property value. The feature group name of the service issue.
func (m *ServiceHealthIssue) SetFeatureGroup(value *string)() {
    m.featureGroup = value
}
// SetImpactDescription sets the impactDescription property value. The description of the service issue impact.
func (m *ServiceHealthIssue) SetImpactDescription(value *string)() {
    m.impactDescription = value
}
// SetIsResolved sets the isResolved property value. Indicates whether the issue is resolved.
func (m *ServiceHealthIssue) SetIsResolved(value *bool)() {
    m.isResolved = value
}
// SetOrigin sets the origin property value. The origin property
func (m *ServiceHealthIssue) SetOrigin(value *ServiceHealthOrigin)() {
    m.origin = value
}
// SetPosts sets the posts property value. Collection of historical posts for the service issue.
func (m *ServiceHealthIssue) SetPosts(value []ServiceHealthIssuePostable)() {
    m.posts = value
}
// SetService sets the service property value. Indicates the service affected by the issue.
func (m *ServiceHealthIssue) SetService(value *string)() {
    m.service = value
}
// SetStatus sets the status property value. The status property
func (m *ServiceHealthIssue) SetStatus(value *ServiceHealthStatus)() {
    m.status = value
}
