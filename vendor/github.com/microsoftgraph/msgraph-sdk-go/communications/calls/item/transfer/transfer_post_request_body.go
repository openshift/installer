package transfer

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
)

// TransferPostRequestBody provides operations to call the transfer method.
type TransferPostRequestBody struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The transferee property
    transferee iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ParticipantInfoable
    // The transferTarget property
    transferTarget iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.InvitationParticipantInfoable
}
// NewTransferPostRequestBody instantiates a new transferPostRequestBody and sets the default values.
func NewTransferPostRequestBody()(*TransferPostRequestBody) {
    m := &TransferPostRequestBody{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateTransferPostRequestBodyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateTransferPostRequestBodyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTransferPostRequestBody(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *TransferPostRequestBody) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
func (m *TransferPostRequestBody) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["transferee"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateParticipantInfoFromDiscriminatorValue , m.SetTransferee)
    res["transferTarget"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateInvitationParticipantInfoFromDiscriminatorValue , m.SetTransferTarget)
    return res
}
// GetTransferee gets the transferee property value. The transferee property
func (m *TransferPostRequestBody) GetTransferee()(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ParticipantInfoable) {
    return m.transferee
}
// GetTransferTarget gets the transferTarget property value. The transferTarget property
func (m *TransferPostRequestBody) GetTransferTarget()(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.InvitationParticipantInfoable) {
    return m.transferTarget
}
// Serialize serializes information the current object
func (m *TransferPostRequestBody) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteObjectValue("transferee", m.GetTransferee())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("transferTarget", m.GetTransferTarget())
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
func (m *TransferPostRequestBody) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetTransferee sets the transferee property value. The transferee property
func (m *TransferPostRequestBody) SetTransferee(value iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ParticipantInfoable)() {
    m.transferee = value
}
// SetTransferTarget sets the transferTarget property value. The transferTarget property
func (m *TransferPostRequestBody) SetTransferTarget(value iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.InvitationParticipantInfoable)() {
    m.transferTarget = value
}
