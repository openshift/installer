package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Photo 
type Photo struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Camera manufacturer. Read-only.
    cameraMake *string
    // Camera model. Read-only.
    cameraModel *string
    // The denominator for the exposure time fraction from the camera. Read-only.
    exposureDenominator *float64
    // The numerator for the exposure time fraction from the camera. Read-only.
    exposureNumerator *float64
    // The F-stop value from the camera. Read-only.
    fNumber *float64
    // The focal length from the camera. Read-only.
    focalLength *float64
    // The ISO value from the camera. Read-only.
    iso *int32
    // The OdataType property
    odataType *string
    // The orientation value from the camera. Writable on OneDrive Personal.
    orientation *int32
    // Represents the date and time the photo was taken. Read-only.
    takenDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
}
// NewPhoto instantiates a new photo and sets the default values.
func NewPhoto()(*Photo) {
    m := &Photo{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreatePhotoFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreatePhotoFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPhoto(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *Photo) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetCameraMake gets the cameraMake property value. Camera manufacturer. Read-only.
func (m *Photo) GetCameraMake()(*string) {
    return m.cameraMake
}
// GetCameraModel gets the cameraModel property value. Camera model. Read-only.
func (m *Photo) GetCameraModel()(*string) {
    return m.cameraModel
}
// GetExposureDenominator gets the exposureDenominator property value. The denominator for the exposure time fraction from the camera. Read-only.
func (m *Photo) GetExposureDenominator()(*float64) {
    return m.exposureDenominator
}
// GetExposureNumerator gets the exposureNumerator property value. The numerator for the exposure time fraction from the camera. Read-only.
func (m *Photo) GetExposureNumerator()(*float64) {
    return m.exposureNumerator
}
// GetFieldDeserializers the deserialization information for the current model
func (m *Photo) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["cameraMake"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetCameraMake)
    res["cameraModel"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetCameraModel)
    res["exposureDenominator"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetFloat64Value(m.SetExposureDenominator)
    res["exposureNumerator"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetFloat64Value(m.SetExposureNumerator)
    res["fNumber"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetFloat64Value(m.SetFNumber)
    res["focalLength"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetFloat64Value(m.SetFocalLength)
    res["iso"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetIso)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["orientation"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetOrientation)
    res["takenDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetTakenDateTime)
    return res
}
// GetFNumber gets the fNumber property value. The F-stop value from the camera. Read-only.
func (m *Photo) GetFNumber()(*float64) {
    return m.fNumber
}
// GetFocalLength gets the focalLength property value. The focal length from the camera. Read-only.
func (m *Photo) GetFocalLength()(*float64) {
    return m.focalLength
}
// GetIso gets the iso property value. The ISO value from the camera. Read-only.
func (m *Photo) GetIso()(*int32) {
    return m.iso
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *Photo) GetOdataType()(*string) {
    return m.odataType
}
// GetOrientation gets the orientation property value. The orientation value from the camera. Writable on OneDrive Personal.
func (m *Photo) GetOrientation()(*int32) {
    return m.orientation
}
// GetTakenDateTime gets the takenDateTime property value. Represents the date and time the photo was taken. Read-only.
func (m *Photo) GetTakenDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.takenDateTime
}
// Serialize serializes information the current object
func (m *Photo) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("cameraMake", m.GetCameraMake())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("cameraModel", m.GetCameraModel())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat64Value("exposureDenominator", m.GetExposureDenominator())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat64Value("exposureNumerator", m.GetExposureNumerator())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat64Value("fNumber", m.GetFNumber())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat64Value("focalLength", m.GetFocalLength())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("iso", m.GetIso())
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
        err := writer.WriteInt32Value("orientation", m.GetOrientation())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteTimeValue("takenDateTime", m.GetTakenDateTime())
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
func (m *Photo) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetCameraMake sets the cameraMake property value. Camera manufacturer. Read-only.
func (m *Photo) SetCameraMake(value *string)() {
    m.cameraMake = value
}
// SetCameraModel sets the cameraModel property value. Camera model. Read-only.
func (m *Photo) SetCameraModel(value *string)() {
    m.cameraModel = value
}
// SetExposureDenominator sets the exposureDenominator property value. The denominator for the exposure time fraction from the camera. Read-only.
func (m *Photo) SetExposureDenominator(value *float64)() {
    m.exposureDenominator = value
}
// SetExposureNumerator sets the exposureNumerator property value. The numerator for the exposure time fraction from the camera. Read-only.
func (m *Photo) SetExposureNumerator(value *float64)() {
    m.exposureNumerator = value
}
// SetFNumber sets the fNumber property value. The F-stop value from the camera. Read-only.
func (m *Photo) SetFNumber(value *float64)() {
    m.fNumber = value
}
// SetFocalLength sets the focalLength property value. The focal length from the camera. Read-only.
func (m *Photo) SetFocalLength(value *float64)() {
    m.focalLength = value
}
// SetIso sets the iso property value. The ISO value from the camera. Read-only.
func (m *Photo) SetIso(value *int32)() {
    m.iso = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *Photo) SetOdataType(value *string)() {
    m.odataType = value
}
// SetOrientation sets the orientation property value. The orientation value from the camera. Writable on OneDrive Personal.
func (m *Photo) SetOrientation(value *int32)() {
    m.orientation = value
}
// SetTakenDateTime sets the takenDateTime property value. Represents the date and time the photo was taken. Read-only.
func (m *Photo) SetTakenDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.takenDateTime = value
}
