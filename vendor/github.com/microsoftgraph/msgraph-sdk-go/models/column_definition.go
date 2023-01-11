package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ColumnDefinition provides operations to manage the collection of agreement entities.
type ColumnDefinition struct {
    Entity
    // This column stores boolean values.
    boolean BooleanColumnable
    // This column's data is calculated based on other columns.
    calculated CalculatedColumnable
    // This column stores data from a list of choices.
    choice ChoiceColumnable
    // For site columns, the name of the group this column belongs to. Helps organize related columns.
    columnGroup *string
    // This column stores content approval status.
    contentApprovalStatus ContentApprovalStatusColumnable
    // This column stores currency values.
    currency CurrencyColumnable
    // This column stores DateTime values.
    dateTime DateTimeColumnable
    // The default value for this column.
    defaultValue DefaultColumnValueable
    // The user-facing description of the column.
    description *string
    // The user-facing name of the column.
    displayName *string
    // If true, no two list items may have the same value for this column.
    enforceUniqueValues *bool
    // This column stores a geolocation.
    geolocation GeolocationColumnable
    // Specifies whether the column is displayed in the user interface.
    hidden *bool
    // This column stores hyperlink or picture values.
    hyperlinkOrPicture HyperlinkOrPictureColumnable
    // Specifies whether the column values can be used for sorting and searching.
    indexed *bool
    // Indicates whether this column can be deleted.
    isDeletable *bool
    // Indicates whether values in the column can be reordered. Read-only.
    isReorderable *bool
    // Specifies whether the column can be changed.
    isSealed *bool
    // This column's data is looked up from another source in the site.
    lookup LookupColumnable
    // The API-facing name of the column as it appears in the [fields][] on a [listItem][]. For the user-facing name, see displayName.
    name *string
    // This column stores number values.
    number NumberColumnable
    // This column stores Person or Group values.
    personOrGroup PersonOrGroupColumnable
    // If 'true', changes to this column will be propagated to lists that implement the column.
    propagateChanges *bool
    // Specifies whether the column values can be modified.
    readOnly *bool
    // Specifies whether the column value isn't optional.
    required *bool
    // The source column for the content type column.
    sourceColumn ColumnDefinitionable
    // ContentType from which this column is inherited from. Present only in contentTypes columns response. Read-only.
    sourceContentType ContentTypeInfoable
    // This column stores taxonomy terms.
    term TermColumnable
    // This column stores text values.
    text TextColumnable
    // This column stores thumbnail values.
    thumbnail ThumbnailColumnable
    // For site columns, the type of column. Read-only.
    type_escaped *ColumnTypes
    // This column stores validation formula and message for the column.
    validation ColumnValidationable
}
// NewColumnDefinition instantiates a new columnDefinition and sets the default values.
func NewColumnDefinition()(*ColumnDefinition) {
    m := &ColumnDefinition{
        Entity: *NewEntity(),
    }
    return m
}
// CreateColumnDefinitionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateColumnDefinitionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewColumnDefinition(), nil
}
// GetBoolean gets the boolean property value. This column stores boolean values.
func (m *ColumnDefinition) GetBoolean()(BooleanColumnable) {
    return m.boolean
}
// GetCalculated gets the calculated property value. This column's data is calculated based on other columns.
func (m *ColumnDefinition) GetCalculated()(CalculatedColumnable) {
    return m.calculated
}
// GetChoice gets the choice property value. This column stores data from a list of choices.
func (m *ColumnDefinition) GetChoice()(ChoiceColumnable) {
    return m.choice
}
// GetColumnGroup gets the columnGroup property value. For site columns, the name of the group this column belongs to. Helps organize related columns.
func (m *ColumnDefinition) GetColumnGroup()(*string) {
    return m.columnGroup
}
// GetContentApprovalStatus gets the contentApprovalStatus property value. This column stores content approval status.
func (m *ColumnDefinition) GetContentApprovalStatus()(ContentApprovalStatusColumnable) {
    return m.contentApprovalStatus
}
// GetCurrency gets the currency property value. This column stores currency values.
func (m *ColumnDefinition) GetCurrency()(CurrencyColumnable) {
    return m.currency
}
// GetDateTime gets the dateTime property value. This column stores DateTime values.
func (m *ColumnDefinition) GetDateTime()(DateTimeColumnable) {
    return m.dateTime
}
// GetDefaultValue gets the defaultValue property value. The default value for this column.
func (m *ColumnDefinition) GetDefaultValue()(DefaultColumnValueable) {
    return m.defaultValue
}
// GetDescription gets the description property value. The user-facing description of the column.
func (m *ColumnDefinition) GetDescription()(*string) {
    return m.description
}
// GetDisplayName gets the displayName property value. The user-facing name of the column.
func (m *ColumnDefinition) GetDisplayName()(*string) {
    return m.displayName
}
// GetEnforceUniqueValues gets the enforceUniqueValues property value. If true, no two list items may have the same value for this column.
func (m *ColumnDefinition) GetEnforceUniqueValues()(*bool) {
    return m.enforceUniqueValues
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ColumnDefinition) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["boolean"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateBooleanColumnFromDiscriminatorValue , m.SetBoolean)
    res["calculated"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateCalculatedColumnFromDiscriminatorValue , m.SetCalculated)
    res["choice"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateChoiceColumnFromDiscriminatorValue , m.SetChoice)
    res["columnGroup"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetColumnGroup)
    res["contentApprovalStatus"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateContentApprovalStatusColumnFromDiscriminatorValue , m.SetContentApprovalStatus)
    res["currency"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateCurrencyColumnFromDiscriminatorValue , m.SetCurrency)
    res["dateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateDateTimeColumnFromDiscriminatorValue , m.SetDateTime)
    res["defaultValue"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateDefaultColumnValueFromDiscriminatorValue , m.SetDefaultValue)
    res["description"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDescription)
    res["displayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDisplayName)
    res["enforceUniqueValues"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetEnforceUniqueValues)
    res["geolocation"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateGeolocationColumnFromDiscriminatorValue , m.SetGeolocation)
    res["hidden"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetHidden)
    res["hyperlinkOrPicture"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateHyperlinkOrPictureColumnFromDiscriminatorValue , m.SetHyperlinkOrPicture)
    res["indexed"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIndexed)
    res["isDeletable"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsDeletable)
    res["isReorderable"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsReorderable)
    res["isSealed"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsSealed)
    res["lookup"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateLookupColumnFromDiscriminatorValue , m.SetLookup)
    res["name"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetName)
    res["number"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateNumberColumnFromDiscriminatorValue , m.SetNumber)
    res["personOrGroup"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreatePersonOrGroupColumnFromDiscriminatorValue , m.SetPersonOrGroup)
    res["propagateChanges"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetPropagateChanges)
    res["readOnly"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetReadOnly)
    res["required"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetRequired)
    res["sourceColumn"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateColumnDefinitionFromDiscriminatorValue , m.SetSourceColumn)
    res["sourceContentType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateContentTypeInfoFromDiscriminatorValue , m.SetSourceContentType)
    res["term"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateTermColumnFromDiscriminatorValue , m.SetTerm)
    res["text"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateTextColumnFromDiscriminatorValue , m.SetText)
    res["thumbnail"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateThumbnailColumnFromDiscriminatorValue , m.SetThumbnail)
    res["type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseColumnTypes , m.SetType)
    res["validation"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateColumnValidationFromDiscriminatorValue , m.SetValidation)
    return res
}
// GetGeolocation gets the geolocation property value. This column stores a geolocation.
func (m *ColumnDefinition) GetGeolocation()(GeolocationColumnable) {
    return m.geolocation
}
// GetHidden gets the hidden property value. Specifies whether the column is displayed in the user interface.
func (m *ColumnDefinition) GetHidden()(*bool) {
    return m.hidden
}
// GetHyperlinkOrPicture gets the hyperlinkOrPicture property value. This column stores hyperlink or picture values.
func (m *ColumnDefinition) GetHyperlinkOrPicture()(HyperlinkOrPictureColumnable) {
    return m.hyperlinkOrPicture
}
// GetIndexed gets the indexed property value. Specifies whether the column values can be used for sorting and searching.
func (m *ColumnDefinition) GetIndexed()(*bool) {
    return m.indexed
}
// GetIsDeletable gets the isDeletable property value. Indicates whether this column can be deleted.
func (m *ColumnDefinition) GetIsDeletable()(*bool) {
    return m.isDeletable
}
// GetIsReorderable gets the isReorderable property value. Indicates whether values in the column can be reordered. Read-only.
func (m *ColumnDefinition) GetIsReorderable()(*bool) {
    return m.isReorderable
}
// GetIsSealed gets the isSealed property value. Specifies whether the column can be changed.
func (m *ColumnDefinition) GetIsSealed()(*bool) {
    return m.isSealed
}
// GetLookup gets the lookup property value. This column's data is looked up from another source in the site.
func (m *ColumnDefinition) GetLookup()(LookupColumnable) {
    return m.lookup
}
// GetName gets the name property value. The API-facing name of the column as it appears in the [fields][] on a [listItem][]. For the user-facing name, see displayName.
func (m *ColumnDefinition) GetName()(*string) {
    return m.name
}
// GetNumber gets the number property value. This column stores number values.
func (m *ColumnDefinition) GetNumber()(NumberColumnable) {
    return m.number
}
// GetPersonOrGroup gets the personOrGroup property value. This column stores Person or Group values.
func (m *ColumnDefinition) GetPersonOrGroup()(PersonOrGroupColumnable) {
    return m.personOrGroup
}
// GetPropagateChanges gets the propagateChanges property value. If 'true', changes to this column will be propagated to lists that implement the column.
func (m *ColumnDefinition) GetPropagateChanges()(*bool) {
    return m.propagateChanges
}
// GetReadOnly gets the readOnly property value. Specifies whether the column values can be modified.
func (m *ColumnDefinition) GetReadOnly()(*bool) {
    return m.readOnly
}
// GetRequired gets the required property value. Specifies whether the column value isn't optional.
func (m *ColumnDefinition) GetRequired()(*bool) {
    return m.required
}
// GetSourceColumn gets the sourceColumn property value. The source column for the content type column.
func (m *ColumnDefinition) GetSourceColumn()(ColumnDefinitionable) {
    return m.sourceColumn
}
// GetSourceContentType gets the sourceContentType property value. ContentType from which this column is inherited from. Present only in contentTypes columns response. Read-only.
func (m *ColumnDefinition) GetSourceContentType()(ContentTypeInfoable) {
    return m.sourceContentType
}
// GetTerm gets the term property value. This column stores taxonomy terms.
func (m *ColumnDefinition) GetTerm()(TermColumnable) {
    return m.term
}
// GetText gets the text property value. This column stores text values.
func (m *ColumnDefinition) GetText()(TextColumnable) {
    return m.text
}
// GetThumbnail gets the thumbnail property value. This column stores thumbnail values.
func (m *ColumnDefinition) GetThumbnail()(ThumbnailColumnable) {
    return m.thumbnail
}
// GetType gets the type property value. For site columns, the type of column. Read-only.
func (m *ColumnDefinition) GetType()(*ColumnTypes) {
    return m.type_escaped
}
// GetValidation gets the validation property value. This column stores validation formula and message for the column.
func (m *ColumnDefinition) GetValidation()(ColumnValidationable) {
    return m.validation
}
// Serialize serializes information the current object
func (m *ColumnDefinition) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("boolean", m.GetBoolean())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("calculated", m.GetCalculated())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("choice", m.GetChoice())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("columnGroup", m.GetColumnGroup())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("contentApprovalStatus", m.GetContentApprovalStatus())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("currency", m.GetCurrency())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("dateTime", m.GetDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("defaultValue", m.GetDefaultValue())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("description", m.GetDescription())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("enforceUniqueValues", m.GetEnforceUniqueValues())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("geolocation", m.GetGeolocation())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("hidden", m.GetHidden())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("hyperlinkOrPicture", m.GetHyperlinkOrPicture())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("indexed", m.GetIndexed())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isDeletable", m.GetIsDeletable())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isReorderable", m.GetIsReorderable())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isSealed", m.GetIsSealed())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("lookup", m.GetLookup())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("name", m.GetName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("number", m.GetNumber())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("personOrGroup", m.GetPersonOrGroup())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("propagateChanges", m.GetPropagateChanges())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("readOnly", m.GetReadOnly())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("required", m.GetRequired())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("sourceColumn", m.GetSourceColumn())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("sourceContentType", m.GetSourceContentType())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("term", m.GetTerm())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("text", m.GetText())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("thumbnail", m.GetThumbnail())
        if err != nil {
            return err
        }
    }
    if m.GetType() != nil {
        cast := (*m.GetType()).String()
        err = writer.WriteStringValue("type", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("validation", m.GetValidation())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetBoolean sets the boolean property value. This column stores boolean values.
func (m *ColumnDefinition) SetBoolean(value BooleanColumnable)() {
    m.boolean = value
}
// SetCalculated sets the calculated property value. This column's data is calculated based on other columns.
func (m *ColumnDefinition) SetCalculated(value CalculatedColumnable)() {
    m.calculated = value
}
// SetChoice sets the choice property value. This column stores data from a list of choices.
func (m *ColumnDefinition) SetChoice(value ChoiceColumnable)() {
    m.choice = value
}
// SetColumnGroup sets the columnGroup property value. For site columns, the name of the group this column belongs to. Helps organize related columns.
func (m *ColumnDefinition) SetColumnGroup(value *string)() {
    m.columnGroup = value
}
// SetContentApprovalStatus sets the contentApprovalStatus property value. This column stores content approval status.
func (m *ColumnDefinition) SetContentApprovalStatus(value ContentApprovalStatusColumnable)() {
    m.contentApprovalStatus = value
}
// SetCurrency sets the currency property value. This column stores currency values.
func (m *ColumnDefinition) SetCurrency(value CurrencyColumnable)() {
    m.currency = value
}
// SetDateTime sets the dateTime property value. This column stores DateTime values.
func (m *ColumnDefinition) SetDateTime(value DateTimeColumnable)() {
    m.dateTime = value
}
// SetDefaultValue sets the defaultValue property value. The default value for this column.
func (m *ColumnDefinition) SetDefaultValue(value DefaultColumnValueable)() {
    m.defaultValue = value
}
// SetDescription sets the description property value. The user-facing description of the column.
func (m *ColumnDefinition) SetDescription(value *string)() {
    m.description = value
}
// SetDisplayName sets the displayName property value. The user-facing name of the column.
func (m *ColumnDefinition) SetDisplayName(value *string)() {
    m.displayName = value
}
// SetEnforceUniqueValues sets the enforceUniqueValues property value. If true, no two list items may have the same value for this column.
func (m *ColumnDefinition) SetEnforceUniqueValues(value *bool)() {
    m.enforceUniqueValues = value
}
// SetGeolocation sets the geolocation property value. This column stores a geolocation.
func (m *ColumnDefinition) SetGeolocation(value GeolocationColumnable)() {
    m.geolocation = value
}
// SetHidden sets the hidden property value. Specifies whether the column is displayed in the user interface.
func (m *ColumnDefinition) SetHidden(value *bool)() {
    m.hidden = value
}
// SetHyperlinkOrPicture sets the hyperlinkOrPicture property value. This column stores hyperlink or picture values.
func (m *ColumnDefinition) SetHyperlinkOrPicture(value HyperlinkOrPictureColumnable)() {
    m.hyperlinkOrPicture = value
}
// SetIndexed sets the indexed property value. Specifies whether the column values can be used for sorting and searching.
func (m *ColumnDefinition) SetIndexed(value *bool)() {
    m.indexed = value
}
// SetIsDeletable sets the isDeletable property value. Indicates whether this column can be deleted.
func (m *ColumnDefinition) SetIsDeletable(value *bool)() {
    m.isDeletable = value
}
// SetIsReorderable sets the isReorderable property value. Indicates whether values in the column can be reordered. Read-only.
func (m *ColumnDefinition) SetIsReorderable(value *bool)() {
    m.isReorderable = value
}
// SetIsSealed sets the isSealed property value. Specifies whether the column can be changed.
func (m *ColumnDefinition) SetIsSealed(value *bool)() {
    m.isSealed = value
}
// SetLookup sets the lookup property value. This column's data is looked up from another source in the site.
func (m *ColumnDefinition) SetLookup(value LookupColumnable)() {
    m.lookup = value
}
// SetName sets the name property value. The API-facing name of the column as it appears in the [fields][] on a [listItem][]. For the user-facing name, see displayName.
func (m *ColumnDefinition) SetName(value *string)() {
    m.name = value
}
// SetNumber sets the number property value. This column stores number values.
func (m *ColumnDefinition) SetNumber(value NumberColumnable)() {
    m.number = value
}
// SetPersonOrGroup sets the personOrGroup property value. This column stores Person or Group values.
func (m *ColumnDefinition) SetPersonOrGroup(value PersonOrGroupColumnable)() {
    m.personOrGroup = value
}
// SetPropagateChanges sets the propagateChanges property value. If 'true', changes to this column will be propagated to lists that implement the column.
func (m *ColumnDefinition) SetPropagateChanges(value *bool)() {
    m.propagateChanges = value
}
// SetReadOnly sets the readOnly property value. Specifies whether the column values can be modified.
func (m *ColumnDefinition) SetReadOnly(value *bool)() {
    m.readOnly = value
}
// SetRequired sets the required property value. Specifies whether the column value isn't optional.
func (m *ColumnDefinition) SetRequired(value *bool)() {
    m.required = value
}
// SetSourceColumn sets the sourceColumn property value. The source column for the content type column.
func (m *ColumnDefinition) SetSourceColumn(value ColumnDefinitionable)() {
    m.sourceColumn = value
}
// SetSourceContentType sets the sourceContentType property value. ContentType from which this column is inherited from. Present only in contentTypes columns response. Read-only.
func (m *ColumnDefinition) SetSourceContentType(value ContentTypeInfoable)() {
    m.sourceContentType = value
}
// SetTerm sets the term property value. This column stores taxonomy terms.
func (m *ColumnDefinition) SetTerm(value TermColumnable)() {
    m.term = value
}
// SetText sets the text property value. This column stores text values.
func (m *ColumnDefinition) SetText(value TextColumnable)() {
    m.text = value
}
// SetThumbnail sets the thumbnail property value. This column stores thumbnail values.
func (m *ColumnDefinition) SetThumbnail(value ThumbnailColumnable)() {
    m.thumbnail = value
}
// SetType sets the type property value. For site columns, the type of column. Read-only.
func (m *ColumnDefinition) SetType(value *ColumnTypes)() {
    m.type_escaped = value
}
// SetValidation sets the validation property value. This column stores validation formula and message for the column.
func (m *ColumnDefinition) SetValidation(value ColumnValidationable)() {
    m.validation = value
}
