package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// EducationRoot 
type EducationRoot struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The classes property
    classes []EducationClassable
    // The me property
    me EducationUserable
    // The OdataType property
    odataType *string
    // The schools property
    schools []EducationSchoolable
    // The users property
    users []EducationUserable
}
// NewEducationRoot instantiates a new EducationRoot and sets the default values.
func NewEducationRoot()(*EducationRoot) {
    m := &EducationRoot{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateEducationRootFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateEducationRootFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEducationRoot(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *EducationRoot) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetClasses gets the classes property value. The classes property
func (m *EducationRoot) GetClasses()([]EducationClassable) {
    return m.classes
}
// GetFieldDeserializers the deserialization information for the current model
func (m *EducationRoot) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["classes"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateEducationClassFromDiscriminatorValue , m.SetClasses)
    res["me"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateEducationUserFromDiscriminatorValue , m.SetMe)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["schools"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateEducationSchoolFromDiscriminatorValue , m.SetSchools)
    res["users"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateEducationUserFromDiscriminatorValue , m.SetUsers)
    return res
}
// GetMe gets the me property value. The me property
func (m *EducationRoot) GetMe()(EducationUserable) {
    return m.me
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *EducationRoot) GetOdataType()(*string) {
    return m.odataType
}
// GetSchools gets the schools property value. The schools property
func (m *EducationRoot) GetSchools()([]EducationSchoolable) {
    return m.schools
}
// GetUsers gets the users property value. The users property
func (m *EducationRoot) GetUsers()([]EducationUserable) {
    return m.users
}
// Serialize serializes information the current object
func (m *EducationRoot) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetClasses() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetClasses())
        err := writer.WriteCollectionOfObjectValues("classes", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("me", m.GetMe())
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
    if m.GetSchools() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetSchools())
        err := writer.WriteCollectionOfObjectValues("schools", cast)
        if err != nil {
            return err
        }
    }
    if m.GetUsers() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetUsers())
        err := writer.WriteCollectionOfObjectValues("users", cast)
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
func (m *EducationRoot) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetClasses sets the classes property value. The classes property
func (m *EducationRoot) SetClasses(value []EducationClassable)() {
    m.classes = value
}
// SetMe sets the me property value. The me property
func (m *EducationRoot) SetMe(value EducationUserable)() {
    m.me = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *EducationRoot) SetOdataType(value *string)() {
    m.odataType = value
}
// SetSchools sets the schools property value. The schools property
func (m *EducationRoot) SetSchools(value []EducationSchoolable)() {
    m.schools = value
}
// SetUsers sets the users property value. The users property
func (m *EducationRoot) SetUsers(value []EducationUserable)() {
    m.users = value
}
