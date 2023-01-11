package termstore

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
)

// Store 
type Store struct {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entity
    // Default language of the term store.
    defaultLanguageTag *string
    // Collection of all groups available in the term store.
    groups []Groupable
    // List of languages for the term store.
    languageTags []string
    // Collection of all sets available in the term store. This relationship can only be used to load a specific term set.
    sets []Setable
}
// NewStore instantiates a new store and sets the default values.
func NewStore()(*Store) {
    m := &Store{
        Entity: *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.NewEntity(),
    }
    return m
}
// CreateStoreFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateStoreFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewStore(), nil
}
// GetDefaultLanguageTag gets the defaultLanguageTag property value. Default language of the term store.
func (m *Store) GetDefaultLanguageTag()(*string) {
    return m.defaultLanguageTag
}
// GetFieldDeserializers the deserialization information for the current model
func (m *Store) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["defaultLanguageTag"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDefaultLanguageTag)
    res["groups"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateGroupFromDiscriminatorValue , m.SetGroups)
    res["languageTags"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("string" , m.SetLanguageTags)
    res["sets"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateSetFromDiscriminatorValue , m.SetSets)
    return res
}
// GetGroups gets the groups property value. Collection of all groups available in the term store.
func (m *Store) GetGroups()([]Groupable) {
    return m.groups
}
// GetLanguageTags gets the languageTags property value. List of languages for the term store.
func (m *Store) GetLanguageTags()([]string) {
    return m.languageTags
}
// GetSets gets the sets property value. Collection of all sets available in the term store. This relationship can only be used to load a specific term set.
func (m *Store) GetSets()([]Setable) {
    return m.sets
}
// Serialize serializes information the current object
func (m *Store) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("defaultLanguageTag", m.GetDefaultLanguageTag())
        if err != nil {
            return err
        }
    }
    if m.GetGroups() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetGroups())
        err = writer.WriteCollectionOfObjectValues("groups", cast)
        if err != nil {
            return err
        }
    }
    if m.GetLanguageTags() != nil {
        err = writer.WriteCollectionOfStringValues("languageTags", m.GetLanguageTags())
        if err != nil {
            return err
        }
    }
    if m.GetSets() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetSets())
        err = writer.WriteCollectionOfObjectValues("sets", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDefaultLanguageTag sets the defaultLanguageTag property value. Default language of the term store.
func (m *Store) SetDefaultLanguageTag(value *string)() {
    m.defaultLanguageTag = value
}
// SetGroups sets the groups property value. Collection of all groups available in the term store.
func (m *Store) SetGroups(value []Groupable)() {
    m.groups = value
}
// SetLanguageTags sets the languageTags property value. List of languages for the term store.
func (m *Store) SetLanguageTags(value []string)() {
    m.languageTags = value
}
// SetSets sets the sets property value. Collection of all sets available in the term store. This relationship can only be used to load a specific term set.
func (m *Store) SetSets(value []Setable)() {
    m.sets = value
}
