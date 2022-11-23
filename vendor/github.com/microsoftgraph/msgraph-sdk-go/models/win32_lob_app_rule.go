package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Win32LobAppRule a base complex type to store the detection or requirement rule data for a Win32 LOB app.
type Win32LobAppRule struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The OdataType property
    odataType *string
    // Contains rule types for Win32 LOB apps.
    ruleType *Win32LobAppRuleType
}
// NewWin32LobAppRule instantiates a new win32LobAppRule and sets the default values.
func NewWin32LobAppRule()(*Win32LobAppRule) {
    m := &Win32LobAppRule{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateWin32LobAppRuleFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateWin32LobAppRuleFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
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
                    case "#microsoft.graph.win32LobAppFileSystemRule":
                        return NewWin32LobAppFileSystemRule(), nil
                    case "#microsoft.graph.win32LobAppPowerShellScriptRule":
                        return NewWin32LobAppPowerShellScriptRule(), nil
                    case "#microsoft.graph.win32LobAppProductCodeRule":
                        return NewWin32LobAppProductCodeRule(), nil
                    case "#microsoft.graph.win32LobAppRegistryRule":
                        return NewWin32LobAppRegistryRule(), nil
                }
            }
        }
    }
    return NewWin32LobAppRule(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *Win32LobAppRule) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
func (m *Win32LobAppRule) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["ruleType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseWin32LobAppRuleType , m.SetRuleType)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *Win32LobAppRule) GetOdataType()(*string) {
    return m.odataType
}
// GetRuleType gets the ruleType property value. Contains rule types for Win32 LOB apps.
func (m *Win32LobAppRule) GetRuleType()(*Win32LobAppRuleType) {
    return m.ruleType
}
// Serialize serializes information the current object
func (m *Win32LobAppRule) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    if m.GetRuleType() != nil {
        cast := (*m.GetRuleType()).String()
        err := writer.WriteStringValue("ruleType", &cast)
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
func (m *Win32LobAppRule) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *Win32LobAppRule) SetOdataType(value *string)() {
    m.odataType = value
}
// SetRuleType sets the ruleType property value. Contains rule types for Win32 LOB apps.
func (m *Win32LobAppRule) SetRuleType(value *Win32LobAppRuleType)() {
    m.ruleType = value
}
