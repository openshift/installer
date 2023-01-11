package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// EntitlementManagement 
type EntitlementManagement struct {
    Entity
    // Approval stages for decisions associated with access package assignment requests.
    accessPackageAssignmentApprovals []Approvalable
    // Access packages define the collection of resource roles and the policies for which subjects can request or be assigned access to those resources.
    accessPackages []AccessPackageable
    // Access package assignment policies govern which subjects can request or be assigned an access package via an access package assignment.
    assignmentPolicies []AccessPackageAssignmentPolicyable
    // Access package assignment requests created by or on behalf of a subject.
    assignmentRequests []AccessPackageAssignmentRequestable
    // The assignment of an access package to a subject for a period of time.
    assignments []AccessPackageAssignmentable
    // A container for access packages.
    catalogs []AccessPackageCatalogable
    // References to a directory or domain of another organization whose users can request access.
    connectedOrganizations []ConnectedOrganizationable
    // The settings that control the behavior of Azure AD entitlement management.
    settings EntitlementManagementSettingsable
}
// NewEntitlementManagement instantiates a new EntitlementManagement and sets the default values.
func NewEntitlementManagement()(*EntitlementManagement) {
    m := &EntitlementManagement{
        Entity: *NewEntity(),
    }
    return m
}
// CreateEntitlementManagementFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateEntitlementManagementFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEntitlementManagement(), nil
}
// GetAccessPackageAssignmentApprovals gets the accessPackageAssignmentApprovals property value. Approval stages for decisions associated with access package assignment requests.
func (m *EntitlementManagement) GetAccessPackageAssignmentApprovals()([]Approvalable) {
    return m.accessPackageAssignmentApprovals
}
// GetAccessPackages gets the accessPackages property value. Access packages define the collection of resource roles and the policies for which subjects can request or be assigned access to those resources.
func (m *EntitlementManagement) GetAccessPackages()([]AccessPackageable) {
    return m.accessPackages
}
// GetAssignmentPolicies gets the assignmentPolicies property value. Access package assignment policies govern which subjects can request or be assigned an access package via an access package assignment.
func (m *EntitlementManagement) GetAssignmentPolicies()([]AccessPackageAssignmentPolicyable) {
    return m.assignmentPolicies
}
// GetAssignmentRequests gets the assignmentRequests property value. Access package assignment requests created by or on behalf of a subject.
func (m *EntitlementManagement) GetAssignmentRequests()([]AccessPackageAssignmentRequestable) {
    return m.assignmentRequests
}
// GetAssignments gets the assignments property value. The assignment of an access package to a subject for a period of time.
func (m *EntitlementManagement) GetAssignments()([]AccessPackageAssignmentable) {
    return m.assignments
}
// GetCatalogs gets the catalogs property value. A container for access packages.
func (m *EntitlementManagement) GetCatalogs()([]AccessPackageCatalogable) {
    return m.catalogs
}
// GetConnectedOrganizations gets the connectedOrganizations property value. References to a directory or domain of another organization whose users can request access.
func (m *EntitlementManagement) GetConnectedOrganizations()([]ConnectedOrganizationable) {
    return m.connectedOrganizations
}
// GetFieldDeserializers the deserialization information for the current model
func (m *EntitlementManagement) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["accessPackageAssignmentApprovals"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateApprovalFromDiscriminatorValue , m.SetAccessPackageAssignmentApprovals)
    res["accessPackages"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateAccessPackageFromDiscriminatorValue , m.SetAccessPackages)
    res["assignmentPolicies"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateAccessPackageAssignmentPolicyFromDiscriminatorValue , m.SetAssignmentPolicies)
    res["assignmentRequests"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateAccessPackageAssignmentRequestFromDiscriminatorValue , m.SetAssignmentRequests)
    res["assignments"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateAccessPackageAssignmentFromDiscriminatorValue , m.SetAssignments)
    res["catalogs"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateAccessPackageCatalogFromDiscriminatorValue , m.SetCatalogs)
    res["connectedOrganizations"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateConnectedOrganizationFromDiscriminatorValue , m.SetConnectedOrganizations)
    res["settings"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateEntitlementManagementSettingsFromDiscriminatorValue , m.SetSettings)
    return res
}
// GetSettings gets the settings property value. The settings that control the behavior of Azure AD entitlement management.
func (m *EntitlementManagement) GetSettings()(EntitlementManagementSettingsable) {
    return m.settings
}
// Serialize serializes information the current object
func (m *EntitlementManagement) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAccessPackageAssignmentApprovals() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetAccessPackageAssignmentApprovals())
        err = writer.WriteCollectionOfObjectValues("accessPackageAssignmentApprovals", cast)
        if err != nil {
            return err
        }
    }
    if m.GetAccessPackages() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetAccessPackages())
        err = writer.WriteCollectionOfObjectValues("accessPackages", cast)
        if err != nil {
            return err
        }
    }
    if m.GetAssignmentPolicies() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetAssignmentPolicies())
        err = writer.WriteCollectionOfObjectValues("assignmentPolicies", cast)
        if err != nil {
            return err
        }
    }
    if m.GetAssignmentRequests() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetAssignmentRequests())
        err = writer.WriteCollectionOfObjectValues("assignmentRequests", cast)
        if err != nil {
            return err
        }
    }
    if m.GetAssignments() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetAssignments())
        err = writer.WriteCollectionOfObjectValues("assignments", cast)
        if err != nil {
            return err
        }
    }
    if m.GetCatalogs() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetCatalogs())
        err = writer.WriteCollectionOfObjectValues("catalogs", cast)
        if err != nil {
            return err
        }
    }
    if m.GetConnectedOrganizations() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetConnectedOrganizations())
        err = writer.WriteCollectionOfObjectValues("connectedOrganizations", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("settings", m.GetSettings())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAccessPackageAssignmentApprovals sets the accessPackageAssignmentApprovals property value. Approval stages for decisions associated with access package assignment requests.
func (m *EntitlementManagement) SetAccessPackageAssignmentApprovals(value []Approvalable)() {
    m.accessPackageAssignmentApprovals = value
}
// SetAccessPackages sets the accessPackages property value. Access packages define the collection of resource roles and the policies for which subjects can request or be assigned access to those resources.
func (m *EntitlementManagement) SetAccessPackages(value []AccessPackageable)() {
    m.accessPackages = value
}
// SetAssignmentPolicies sets the assignmentPolicies property value. Access package assignment policies govern which subjects can request or be assigned an access package via an access package assignment.
func (m *EntitlementManagement) SetAssignmentPolicies(value []AccessPackageAssignmentPolicyable)() {
    m.assignmentPolicies = value
}
// SetAssignmentRequests sets the assignmentRequests property value. Access package assignment requests created by or on behalf of a subject.
func (m *EntitlementManagement) SetAssignmentRequests(value []AccessPackageAssignmentRequestable)() {
    m.assignmentRequests = value
}
// SetAssignments sets the assignments property value. The assignment of an access package to a subject for a period of time.
func (m *EntitlementManagement) SetAssignments(value []AccessPackageAssignmentable)() {
    m.assignments = value
}
// SetCatalogs sets the catalogs property value. A container for access packages.
func (m *EntitlementManagement) SetCatalogs(value []AccessPackageCatalogable)() {
    m.catalogs = value
}
// SetConnectedOrganizations sets the connectedOrganizations property value. References to a directory or domain of another organization whose users can request access.
func (m *EntitlementManagement) SetConnectedOrganizations(value []ConnectedOrganizationable)() {
    m.connectedOrganizations = value
}
// SetSettings sets the settings property value. The settings that control the behavior of Azure AD entitlement management.
func (m *EntitlementManagement) SetSettings(value EntitlementManagementSettingsable)() {
    m.settings = value
}
