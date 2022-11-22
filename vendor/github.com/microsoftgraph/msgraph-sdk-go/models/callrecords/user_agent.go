package callrecords

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// UserAgent 
type UserAgent struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Identifies the version of application software used by this endpoint.
    applicationVersion *string
    // User-agent header value reported by this endpoint.
    headerValue *string
    // The OdataType property
    odataType *string
}
// NewUserAgent instantiates a new userAgent and sets the default values.
func NewUserAgent()(*UserAgent) {
    m := &UserAgent{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateUserAgentFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateUserAgentFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    if parseNode != nil {
        mappingValueNode, err := parseNode.GetChildNode("@odata.type")
        if err != nil {
            return nil, err
        }
        if mappingValueNode != nil {
            mappingValue, err := mappingValueNode.GetStringValue()
            if err != nil {
                return nil, err
            }
            if mappingValue != nil {
                switch *mappingValue {
                    case "#microsoft.graph.callRecords.clientUserAgent":
                        return NewClientUserAgent(), nil
                    case "#microsoft.graph.callRecords.serviceUserAgent":
                        return NewServiceUserAgent(), nil
                }
            }
        }
    }
    return NewUserAgent(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *UserAgent) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetApplicationVersion gets the applicationVersion property value. Identifies the version of application software used by this endpoint.
func (m *UserAgent) GetApplicationVersion()(*string) {
    return m.applicationVersion
}
// GetFieldDeserializers the deserialization information for the current model
func (m *UserAgent) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["applicationVersion"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetApplicationVersion)
    res["headerValue"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetHeaderValue)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetHeaderValue gets the headerValue property value. User-agent header value reported by this endpoint.
func (m *UserAgent) GetHeaderValue()(*string) {
    return m.headerValue
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *UserAgent) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *UserAgent) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("applicationVersion", m.GetApplicationVersion())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("headerValue", m.GetHeaderValue())
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
func (m *UserAgent) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetApplicationVersion sets the applicationVersion property value. Identifies the version of application software used by this endpoint.
func (m *UserAgent) SetApplicationVersion(value *string)() {
    m.applicationVersion = value
}
// SetHeaderValue sets the headerValue property value. User-agent header value reported by this endpoint.
func (m *UserAgent) SetHeaderValue(value *string)() {
    m.headerValue = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *UserAgent) SetOdataType(value *string)() {
    m.odataType = value
}
