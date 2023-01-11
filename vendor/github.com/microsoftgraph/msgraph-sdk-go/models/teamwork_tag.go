package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// TeamworkTag provides operations to manage the collection of agreement entities.
type TeamworkTag struct {
    Entity
    // The description of the tag as it will appear to the user in Microsoft Teams.
    description *string
    // The name of the tag as it will appear to the user in Microsoft Teams.
    displayName *string
    // The number of users assigned to the tag.
    memberCount *int32
    // Users assigned to the tag.
    members []TeamworkTagMemberable
    // The type of the tag. Default is standard.
    tagType *TeamworkTagType
    // ID of the team in which the tag is defined.
    teamId *string
}
// NewTeamworkTag instantiates a new teamworkTag and sets the default values.
func NewTeamworkTag()(*TeamworkTag) {
    m := &TeamworkTag{
        Entity: *NewEntity(),
    }
    return m
}
// CreateTeamworkTagFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateTeamworkTagFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTeamworkTag(), nil
}
// GetDescription gets the description property value. The description of the tag as it will appear to the user in Microsoft Teams.
func (m *TeamworkTag) GetDescription()(*string) {
    return m.description
}
// GetDisplayName gets the displayName property value. The name of the tag as it will appear to the user in Microsoft Teams.
func (m *TeamworkTag) GetDisplayName()(*string) {
    return m.displayName
}
// GetFieldDeserializers the deserialization information for the current model
func (m *TeamworkTag) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["description"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDescription)
    res["displayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDisplayName)
    res["memberCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetMemberCount)
    res["members"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateTeamworkTagMemberFromDiscriminatorValue , m.SetMembers)
    res["tagType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseTeamworkTagType , m.SetTagType)
    res["teamId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetTeamId)
    return res
}
// GetMemberCount gets the memberCount property value. The number of users assigned to the tag.
func (m *TeamworkTag) GetMemberCount()(*int32) {
    return m.memberCount
}
// GetMembers gets the members property value. Users assigned to the tag.
func (m *TeamworkTag) GetMembers()([]TeamworkTagMemberable) {
    return m.members
}
// GetTagType gets the tagType property value. The type of the tag. Default is standard.
func (m *TeamworkTag) GetTagType()(*TeamworkTagType) {
    return m.tagType
}
// GetTeamId gets the teamId property value. ID of the team in which the tag is defined.
func (m *TeamworkTag) GetTeamId()(*string) {
    return m.teamId
}
// Serialize serializes information the current object
func (m *TeamworkTag) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("description", m.GetDescription())
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
        err = writer.WriteInt32Value("memberCount", m.GetMemberCount())
        if err != nil {
            return err
        }
    }
    if m.GetMembers() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetMembers())
        err = writer.WriteCollectionOfObjectValues("members", cast)
        if err != nil {
            return err
        }
    }
    if m.GetTagType() != nil {
        cast := (*m.GetTagType()).String()
        err = writer.WriteStringValue("tagType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("teamId", m.GetTeamId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDescription sets the description property value. The description of the tag as it will appear to the user in Microsoft Teams.
func (m *TeamworkTag) SetDescription(value *string)() {
    m.description = value
}
// SetDisplayName sets the displayName property value. The name of the tag as it will appear to the user in Microsoft Teams.
func (m *TeamworkTag) SetDisplayName(value *string)() {
    m.displayName = value
}
// SetMemberCount sets the memberCount property value. The number of users assigned to the tag.
func (m *TeamworkTag) SetMemberCount(value *int32)() {
    m.memberCount = value
}
// SetMembers sets the members property value. Users assigned to the tag.
func (m *TeamworkTag) SetMembers(value []TeamworkTagMemberable)() {
    m.members = value
}
// SetTagType sets the tagType property value. The type of the tag. Default is standard.
func (m *TeamworkTag) SetTagType(value *TeamworkTagType)() {
    m.tagType = value
}
// SetTeamId sets the teamId property value. ID of the team in which the tag is defined.
func (m *TeamworkTag) SetTeamId(value *string)() {
    m.teamId = value
}
