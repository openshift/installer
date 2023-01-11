package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AccessReviewHistoryScheduleSettings 
type AccessReviewHistoryScheduleSettings struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The OdataType property
    odataType *string
    // The recurrence property
    recurrence PatternedRecurrenceable
    // A duration string in ISO 8601 duration format specifying the lookback period of the generated review history data. For example, if a history definition is scheduled to run on the 1st of every month, the reportRange is P1M. In this case, on the first of every month, access review history data will be collected containing only the previous month's review data. Note: Only years, months, and days ISO 8601 properties are supported. Required.
    reportRange *string
}
// NewAccessReviewHistoryScheduleSettings instantiates a new accessReviewHistoryScheduleSettings and sets the default values.
func NewAccessReviewHistoryScheduleSettings()(*AccessReviewHistoryScheduleSettings) {
    m := &AccessReviewHistoryScheduleSettings{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateAccessReviewHistoryScheduleSettingsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAccessReviewHistoryScheduleSettingsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAccessReviewHistoryScheduleSettings(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *AccessReviewHistoryScheduleSettings) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
func (m *AccessReviewHistoryScheduleSettings) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["recurrence"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreatePatternedRecurrenceFromDiscriminatorValue , m.SetRecurrence)
    res["reportRange"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetReportRange)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *AccessReviewHistoryScheduleSettings) GetOdataType()(*string) {
    return m.odataType
}
// GetRecurrence gets the recurrence property value. The recurrence property
func (m *AccessReviewHistoryScheduleSettings) GetRecurrence()(PatternedRecurrenceable) {
    return m.recurrence
}
// GetReportRange gets the reportRange property value. A duration string in ISO 8601 duration format specifying the lookback period of the generated review history data. For example, if a history definition is scheduled to run on the 1st of every month, the reportRange is P1M. In this case, on the first of every month, access review history data will be collected containing only the previous month's review data. Note: Only years, months, and days ISO 8601 properties are supported. Required.
func (m *AccessReviewHistoryScheduleSettings) GetReportRange()(*string) {
    return m.reportRange
}
// Serialize serializes information the current object
func (m *AccessReviewHistoryScheduleSettings) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("recurrence", m.GetRecurrence())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("reportRange", m.GetReportRange())
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
func (m *AccessReviewHistoryScheduleSettings) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *AccessReviewHistoryScheduleSettings) SetOdataType(value *string)() {
    m.odataType = value
}
// SetRecurrence sets the recurrence property value. The recurrence property
func (m *AccessReviewHistoryScheduleSettings) SetRecurrence(value PatternedRecurrenceable)() {
    m.recurrence = value
}
// SetReportRange sets the reportRange property value. A duration string in ISO 8601 duration format specifying the lookback period of the generated review history data. For example, if a history definition is scheduled to run on the 1st of every month, the reportRange is P1M. In this case, on the first of every month, access review history data will be collected containing only the previous month's review data. Note: Only years, months, and days ISO 8601 properties are supported. Required.
func (m *AccessReviewHistoryScheduleSettings) SetReportRange(value *string)() {
    m.reportRange = value
}
