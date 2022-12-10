package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// FileHash 
type FileHash struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // File hash type. Possible values are: unknown, sha1, sha256, md5, authenticodeHash256, lsHash, ctph, peSha1, peSha256.
    hashType *FileHashType
    // Value of the file hash.
    hashValue *string
    // The OdataType property
    odataType *string
}
// NewFileHash instantiates a new fileHash and sets the default values.
func NewFileHash()(*FileHash) {
    m := &FileHash{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateFileHashFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateFileHashFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewFileHash(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *FileHash) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
func (m *FileHash) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["hashType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseFileHashType , m.SetHashType)
    res["hashValue"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetHashValue)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetHashType gets the hashType property value. File hash type. Possible values are: unknown, sha1, sha256, md5, authenticodeHash256, lsHash, ctph, peSha1, peSha256.
func (m *FileHash) GetHashType()(*FileHashType) {
    return m.hashType
}
// GetHashValue gets the hashValue property value. Value of the file hash.
func (m *FileHash) GetHashValue()(*string) {
    return m.hashValue
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *FileHash) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *FileHash) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetHashType() != nil {
        cast := (*m.GetHashType()).String()
        err := writer.WriteStringValue("hashType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("hashValue", m.GetHashValue())
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
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *FileHash) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetHashType sets the hashType property value. File hash type. Possible values are: unknown, sha1, sha256, md5, authenticodeHash256, lsHash, ctph, peSha1, peSha256.
func (m *FileHash) SetHashType(value *FileHashType)() {
    m.hashType = value
}
// SetHashValue sets the hashValue property value. Value of the file hash.
func (m *FileHash) SetHashValue(value *string)() {
    m.hashValue = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *FileHash) SetOdataType(value *string)() {
    m.odataType = value
}
