package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// RubricQuality 
type RubricQuality struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The collection of criteria for this rubric quality.
    criteria []RubricCriterionable
    // The description of this rubric quality.
    description EducationItemBodyable
    // The name of this rubric quality.
    displayName *string
    // The OdataType property
    odataType *string
    // The ID of this resource.
    qualityId *string
    // If present, a numerical weight for this quality.  Weights must add up to 100.
    weight *float32
}
// NewRubricQuality instantiates a new rubricQuality and sets the default values.
func NewRubricQuality()(*RubricQuality) {
    m := &RubricQuality{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateRubricQualityFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateRubricQualityFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewRubricQuality(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *RubricQuality) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetCriteria gets the criteria property value. The collection of criteria for this rubric quality.
func (m *RubricQuality) GetCriteria()([]RubricCriterionable) {
    return m.criteria
}
// GetDescription gets the description property value. The description of this rubric quality.
func (m *RubricQuality) GetDescription()(EducationItemBodyable) {
    return m.description
}
// GetDisplayName gets the displayName property value. The name of this rubric quality.
func (m *RubricQuality) GetDisplayName()(*string) {
    return m.displayName
}
// GetFieldDeserializers the deserialization information for the current model
func (m *RubricQuality) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["criteria"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateRubricCriterionFromDiscriminatorValue , m.SetCriteria)
    res["description"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateEducationItemBodyFromDiscriminatorValue , m.SetDescription)
    res["displayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDisplayName)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["qualityId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetQualityId)
    res["weight"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetFloat32Value(m.SetWeight)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *RubricQuality) GetOdataType()(*string) {
    return m.odataType
}
// GetQualityId gets the qualityId property value. The ID of this resource.
func (m *RubricQuality) GetQualityId()(*string) {
    return m.qualityId
}
// GetWeight gets the weight property value. If present, a numerical weight for this quality.  Weights must add up to 100.
func (m *RubricQuality) GetWeight()(*float32) {
    return m.weight
}
// Serialize serializes information the current object
func (m *RubricQuality) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetCriteria() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetCriteria())
        err := writer.WriteCollectionOfObjectValues("criteria", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("description", m.GetDescription())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("displayName", m.GetDisplayName())
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
    {
        err := writer.WriteStringValue("qualityId", m.GetQualityId())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat32Value("weight", m.GetWeight())
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
func (m *RubricQuality) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetCriteria sets the criteria property value. The collection of criteria for this rubric quality.
func (m *RubricQuality) SetCriteria(value []RubricCriterionable)() {
    m.criteria = value
}
// SetDescription sets the description property value. The description of this rubric quality.
func (m *RubricQuality) SetDescription(value EducationItemBodyable)() {
    m.description = value
}
// SetDisplayName sets the displayName property value. The name of this rubric quality.
func (m *RubricQuality) SetDisplayName(value *string)() {
    m.displayName = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *RubricQuality) SetOdataType(value *string)() {
    m.odataType = value
}
// SetQualityId sets the qualityId property value. The ID of this resource.
func (m *RubricQuality) SetQualityId(value *string)() {
    m.qualityId = value
}
// SetWeight sets the weight property value. If present, a numerical weight for this quality.  Weights must add up to 100.
func (m *RubricQuality) SetWeight(value *float32)() {
    m.weight = value
}
