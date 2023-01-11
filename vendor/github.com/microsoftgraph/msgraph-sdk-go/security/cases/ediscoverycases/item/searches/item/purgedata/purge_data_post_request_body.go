package purgedata

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae "github.com/microsoftgraph/msgraph-sdk-go/models/security"
)

// PurgeDataPostRequestBody provides operations to call the purgeData method.
type PurgeDataPostRequestBody struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The purgeAreas property
    purgeAreas *idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.PurgeAreas
    // The purgeType property
    purgeType *idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.PurgeType
}
// NewPurgeDataPostRequestBody instantiates a new purgeDataPostRequestBody and sets the default values.
func NewPurgeDataPostRequestBody()(*PurgeDataPostRequestBody) {
    m := &PurgeDataPostRequestBody{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreatePurgeDataPostRequestBodyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreatePurgeDataPostRequestBodyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPurgeDataPostRequestBody(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *PurgeDataPostRequestBody) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
func (m *PurgeDataPostRequestBody) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["purgeAreas"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.ParsePurgeAreas , m.SetPurgeAreas)
    res["purgeType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.ParsePurgeType , m.SetPurgeType)
    return res
}
// GetPurgeAreas gets the purgeAreas property value. The purgeAreas property
func (m *PurgeDataPostRequestBody) GetPurgeAreas()(*idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.PurgeAreas) {
    return m.purgeAreas
}
// GetPurgeType gets the purgeType property value. The purgeType property
func (m *PurgeDataPostRequestBody) GetPurgeType()(*idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.PurgeType) {
    return m.purgeType
}
// Serialize serializes information the current object
func (m *PurgeDataPostRequestBody) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetPurgeAreas() != nil {
        cast := (*m.GetPurgeAreas()).String()
        err := writer.WriteStringValue("purgeAreas", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetPurgeType() != nil {
        cast := (*m.GetPurgeType()).String()
        err := writer.WriteStringValue("purgeType", &cast)
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
func (m *PurgeDataPostRequestBody) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetPurgeAreas sets the purgeAreas property value. The purgeAreas property
func (m *PurgeDataPostRequestBody) SetPurgeAreas(value *idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.PurgeAreas)() {
    m.purgeAreas = value
}
// SetPurgeType sets the purgeType property value. The purgeType property
func (m *PurgeDataPostRequestBody) SetPurgeType(value *idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.PurgeType)() {
    m.purgeType = value
}
