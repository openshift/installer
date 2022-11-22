package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// FileSecurityState 
type FileSecurityState struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Complex type containing file hashes (cryptographic and location-sensitive).
    fileHash FileHashable
    // File name (without path).
    name *string
    // The OdataType property
    odataType *string
    // Full file path of the file/imageFile.
    path *string
    // Provider generated/calculated risk score of the alert file. Recommended value range of 0-1, which equates to a percentage.
    riskScore *string
}
// NewFileSecurityState instantiates a new fileSecurityState and sets the default values.
func NewFileSecurityState()(*FileSecurityState) {
    m := &FileSecurityState{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateFileSecurityStateFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateFileSecurityStateFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewFileSecurityState(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *FileSecurityState) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
func (m *FileSecurityState) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["fileHash"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateFileHashFromDiscriminatorValue , m.SetFileHash)
    res["name"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetName)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["path"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetPath)
    res["riskScore"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetRiskScore)
    return res
}
// GetFileHash gets the fileHash property value. Complex type containing file hashes (cryptographic and location-sensitive).
func (m *FileSecurityState) GetFileHash()(FileHashable) {
    return m.fileHash
}
// GetName gets the name property value. File name (without path).
func (m *FileSecurityState) GetName()(*string) {
    return m.name
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *FileSecurityState) GetOdataType()(*string) {
    return m.odataType
}
// GetPath gets the path property value. Full file path of the file/imageFile.
func (m *FileSecurityState) GetPath()(*string) {
    return m.path
}
// GetRiskScore gets the riskScore property value. Provider generated/calculated risk score of the alert file. Recommended value range of 0-1, which equates to a percentage.
func (m *FileSecurityState) GetRiskScore()(*string) {
    return m.riskScore
}
// Serialize serializes information the current object
func (m *FileSecurityState) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteObjectValue("fileHash", m.GetFileHash())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("name", m.GetName())
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
        err := writer.WriteStringValue("path", m.GetPath())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("riskScore", m.GetRiskScore())
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
func (m *FileSecurityState) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetFileHash sets the fileHash property value. Complex type containing file hashes (cryptographic and location-sensitive).
func (m *FileSecurityState) SetFileHash(value FileHashable)() {
    m.fileHash = value
}
// SetName sets the name property value. File name (without path).
func (m *FileSecurityState) SetName(value *string)() {
    m.name = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *FileSecurityState) SetOdataType(value *string)() {
    m.odataType = value
}
// SetPath sets the path property value. Full file path of the file/imageFile.
func (m *FileSecurityState) SetPath(value *string)() {
    m.path = value
}
// SetRiskScore sets the riskScore property value. Provider generated/calculated risk score of the alert file. Recommended value range of 0-1, which equates to a percentage.
func (m *FileSecurityState) SetRiskScore(value *string)() {
    m.riskScore = value
}
