package getmailtips

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
)

// GetMailTipsPostRequestBody provides operations to call the getMailTips method.
type GetMailTipsPostRequestBody struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The EmailAddresses property
    emailAddresses []string
    // The MailTipsOptions property
    mailTipsOptions *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MailTipsType
}
// NewGetMailTipsPostRequestBody instantiates a new getMailTipsPostRequestBody and sets the default values.
func NewGetMailTipsPostRequestBody()(*GetMailTipsPostRequestBody) {
    m := &GetMailTipsPostRequestBody{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateGetMailTipsPostRequestBodyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateGetMailTipsPostRequestBodyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewGetMailTipsPostRequestBody(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *GetMailTipsPostRequestBody) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetEmailAddresses gets the emailAddresses property value. The EmailAddresses property
func (m *GetMailTipsPostRequestBody) GetEmailAddresses()([]string) {
    return m.emailAddresses
}
// GetFieldDeserializers the deserialization information for the current model
func (m *GetMailTipsPostRequestBody) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["emailAddresses"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("string" , m.SetEmailAddresses)
    res["mailTipsOptions"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ParseMailTipsType , m.SetMailTipsOptions)
    return res
}
// GetMailTipsOptions gets the mailTipsOptions property value. The MailTipsOptions property
func (m *GetMailTipsPostRequestBody) GetMailTipsOptions()(*iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MailTipsType) {
    return m.mailTipsOptions
}
// Serialize serializes information the current object
func (m *GetMailTipsPostRequestBody) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetEmailAddresses() != nil {
        err := writer.WriteCollectionOfStringValues("emailAddresses", m.GetEmailAddresses())
        if err != nil {
            return err
        }
    }
    if m.GetMailTipsOptions() != nil {
        cast := (*m.GetMailTipsOptions()).String()
        err := writer.WriteStringValue("mailTipsOptions", &cast)
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
func (m *GetMailTipsPostRequestBody) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetEmailAddresses sets the emailAddresses property value. The EmailAddresses property
func (m *GetMailTipsPostRequestBody) SetEmailAddresses(value []string)() {
    m.emailAddresses = value
}
// SetMailTipsOptions sets the mailTipsOptions property value. The MailTipsOptions property
func (m *GetMailTipsPostRequestBody) SetMailTipsOptions(value *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MailTipsType)() {
    m.mailTipsOptions = value
}
