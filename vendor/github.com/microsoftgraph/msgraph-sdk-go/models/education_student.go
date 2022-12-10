package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// EducationStudent 
type EducationStudent struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Birth date of the student.
    birthDate *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly
    // ID of the student in the source system.
    externalId *string
    // The possible values are: female, male, other, unknownFutureValue.
    gender *EducationGender
    // Current grade level of the student.
    grade *string
    // Year the student is graduating from the school.
    graduationYear *string
    // The OdataType property
    odataType *string
    // Student Number.
    studentNumber *string
}
// NewEducationStudent instantiates a new educationStudent and sets the default values.
func NewEducationStudent()(*EducationStudent) {
    m := &EducationStudent{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateEducationStudentFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateEducationStudentFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEducationStudent(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *EducationStudent) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetBirthDate gets the birthDate property value. Birth date of the student.
func (m *EducationStudent) GetBirthDate()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly) {
    return m.birthDate
}
// GetExternalId gets the externalId property value. ID of the student in the source system.
func (m *EducationStudent) GetExternalId()(*string) {
    return m.externalId
}
// GetFieldDeserializers the deserialization information for the current model
func (m *EducationStudent) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["birthDate"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetDateOnlyValue(m.SetBirthDate)
    res["externalId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetExternalId)
    res["gender"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseEducationGender , m.SetGender)
    res["grade"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetGrade)
    res["graduationYear"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetGraduationYear)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["studentNumber"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetStudentNumber)
    return res
}
// GetGender gets the gender property value. The possible values are: female, male, other, unknownFutureValue.
func (m *EducationStudent) GetGender()(*EducationGender) {
    return m.gender
}
// GetGrade gets the grade property value. Current grade level of the student.
func (m *EducationStudent) GetGrade()(*string) {
    return m.grade
}
// GetGraduationYear gets the graduationYear property value. Year the student is graduating from the school.
func (m *EducationStudent) GetGraduationYear()(*string) {
    return m.graduationYear
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *EducationStudent) GetOdataType()(*string) {
    return m.odataType
}
// GetStudentNumber gets the studentNumber property value. Student Number.
func (m *EducationStudent) GetStudentNumber()(*string) {
    return m.studentNumber
}
// Serialize serializes information the current object
func (m *EducationStudent) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteDateOnlyValue("birthDate", m.GetBirthDate())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("externalId", m.GetExternalId())
        if err != nil {
            return err
        }
    }
    if m.GetGender() != nil {
        cast := (*m.GetGender()).String()
        err := writer.WriteStringValue("gender", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("grade", m.GetGrade())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("graduationYear", m.GetGraduationYear())
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
        err := writer.WriteStringValue("studentNumber", m.GetStudentNumber())
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
func (m *EducationStudent) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetBirthDate sets the birthDate property value. Birth date of the student.
func (m *EducationStudent) SetBirthDate(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)() {
    m.birthDate = value
}
// SetExternalId sets the externalId property value. ID of the student in the source system.
func (m *EducationStudent) SetExternalId(value *string)() {
    m.externalId = value
}
// SetGender sets the gender property value. The possible values are: female, male, other, unknownFutureValue.
func (m *EducationStudent) SetGender(value *EducationGender)() {
    m.gender = value
}
// SetGrade sets the grade property value. Current grade level of the student.
func (m *EducationStudent) SetGrade(value *string)() {
    m.grade = value
}
// SetGraduationYear sets the graduationYear property value. Year the student is graduating from the school.
func (m *EducationStudent) SetGraduationYear(value *string)() {
    m.graduationYear = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *EducationStudent) SetOdataType(value *string)() {
    m.odataType = value
}
// SetStudentNumber sets the studentNumber property value. Student Number.
func (m *EducationStudent) SetStudentNumber(value *string)() {
    m.studentNumber = value
}
