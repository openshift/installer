package security

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
)

// EdiscoveryCaseSettings 
type EdiscoveryCaseSettings struct {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entity
    // The OCR (Optical Character Recognition) settings for the case.
    ocr OcrSettingsable
    // The redundancy (near duplicate and email threading) detection settings for the case.
    redundancyDetection RedundancyDetectionSettingsable
    // The Topic Modeling (Themes) settings for the case.
    topicModeling TopicModelingSettingsable
}
// NewEdiscoveryCaseSettings instantiates a new ediscoveryCaseSettings and sets the default values.
func NewEdiscoveryCaseSettings()(*EdiscoveryCaseSettings) {
    m := &EdiscoveryCaseSettings{
        Entity: *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.NewEntity(),
    }
    return m
}
// CreateEdiscoveryCaseSettingsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateEdiscoveryCaseSettingsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEdiscoveryCaseSettings(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *EdiscoveryCaseSettings) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["ocr"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateOcrSettingsFromDiscriminatorValue , m.SetOcr)
    res["redundancyDetection"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateRedundancyDetectionSettingsFromDiscriminatorValue , m.SetRedundancyDetection)
    res["topicModeling"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateTopicModelingSettingsFromDiscriminatorValue , m.SetTopicModeling)
    return res
}
// GetOcr gets the ocr property value. The OCR (Optical Character Recognition) settings for the case.
func (m *EdiscoveryCaseSettings) GetOcr()(OcrSettingsable) {
    return m.ocr
}
// GetRedundancyDetection gets the redundancyDetection property value. The redundancy (near duplicate and email threading) detection settings for the case.
func (m *EdiscoveryCaseSettings) GetRedundancyDetection()(RedundancyDetectionSettingsable) {
    return m.redundancyDetection
}
// GetTopicModeling gets the topicModeling property value. The Topic Modeling (Themes) settings for the case.
func (m *EdiscoveryCaseSettings) GetTopicModeling()(TopicModelingSettingsable) {
    return m.topicModeling
}
// Serialize serializes information the current object
func (m *EdiscoveryCaseSettings) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("ocr", m.GetOcr())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("redundancyDetection", m.GetRedundancyDetection())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("topicModeling", m.GetTopicModeling())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetOcr sets the ocr property value. The OCR (Optical Character Recognition) settings for the case.
func (m *EdiscoveryCaseSettings) SetOcr(value OcrSettingsable)() {
    m.ocr = value
}
// SetRedundancyDetection sets the redundancyDetection property value. The redundancy (near duplicate and email threading) detection settings for the case.
func (m *EdiscoveryCaseSettings) SetRedundancyDetection(value RedundancyDetectionSettingsable)() {
    m.redundancyDetection = value
}
// SetTopicModeling sets the topicModeling property value. The Topic Modeling (Themes) settings for the case.
func (m *EdiscoveryCaseSettings) SetTopicModeling(value TopicModelingSettingsable)() {
    m.topicModeling = value
}
