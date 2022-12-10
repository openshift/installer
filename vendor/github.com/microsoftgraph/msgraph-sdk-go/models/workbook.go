package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Workbook 
type Workbook struct {
    Entity
    // The application property
    application WorkbookApplicationable
    // The comments property
    comments []WorkbookCommentable
    // The functions property
    functions WorkbookFunctionsable
    // Represents a collection of workbooks scoped named items (named ranges and constants). Read-only.
    names []WorkbookNamedItemable
    // The status of workbook operations. Getting an operation collection is not supported, but you can get the status of a long-running operation if the Location header is returned in the response. Read-only.
    operations []WorkbookOperationable
    // Represents a collection of tables associated with the workbook. Read-only.
    tables []WorkbookTableable
    // Represents a collection of worksheets associated with the workbook. Read-only.
    worksheets []WorkbookWorksheetable
}
// NewWorkbook instantiates a new workbook and sets the default values.
func NewWorkbook()(*Workbook) {
    m := &Workbook{
        Entity: *NewEntity(),
    }
    return m
}
// CreateWorkbookFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateWorkbookFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkbook(), nil
}
// GetApplication gets the application property value. The application property
func (m *Workbook) GetApplication()(WorkbookApplicationable) {
    return m.application
}
// GetComments gets the comments property value. The comments property
func (m *Workbook) GetComments()([]WorkbookCommentable) {
    return m.comments
}
// GetFieldDeserializers the deserialization information for the current model
func (m *Workbook) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["application"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateWorkbookApplicationFromDiscriminatorValue , m.SetApplication)
    res["comments"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateWorkbookCommentFromDiscriminatorValue , m.SetComments)
    res["functions"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateWorkbookFunctionsFromDiscriminatorValue , m.SetFunctions)
    res["names"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateWorkbookNamedItemFromDiscriminatorValue , m.SetNames)
    res["operations"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateWorkbookOperationFromDiscriminatorValue , m.SetOperations)
    res["tables"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateWorkbookTableFromDiscriminatorValue , m.SetTables)
    res["worksheets"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateWorkbookWorksheetFromDiscriminatorValue , m.SetWorksheets)
    return res
}
// GetFunctions gets the functions property value. The functions property
func (m *Workbook) GetFunctions()(WorkbookFunctionsable) {
    return m.functions
}
// GetNames gets the names property value. Represents a collection of workbooks scoped named items (named ranges and constants). Read-only.
func (m *Workbook) GetNames()([]WorkbookNamedItemable) {
    return m.names
}
// GetOperations gets the operations property value. The status of workbook operations. Getting an operation collection is not supported, but you can get the status of a long-running operation if the Location header is returned in the response. Read-only.
func (m *Workbook) GetOperations()([]WorkbookOperationable) {
    return m.operations
}
// GetTables gets the tables property value. Represents a collection of tables associated with the workbook. Read-only.
func (m *Workbook) GetTables()([]WorkbookTableable) {
    return m.tables
}
// GetWorksheets gets the worksheets property value. Represents a collection of worksheets associated with the workbook. Read-only.
func (m *Workbook) GetWorksheets()([]WorkbookWorksheetable) {
    return m.worksheets
}
// Serialize serializes information the current object
func (m *Workbook) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("application", m.GetApplication())
        if err != nil {
            return err
        }
    }
    if m.GetComments() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetComments())
        err = writer.WriteCollectionOfObjectValues("comments", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("functions", m.GetFunctions())
        if err != nil {
            return err
        }
    }
    if m.GetNames() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetNames())
        err = writer.WriteCollectionOfObjectValues("names", cast)
        if err != nil {
            return err
        }
    }
    if m.GetOperations() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetOperations())
        err = writer.WriteCollectionOfObjectValues("operations", cast)
        if err != nil {
            return err
        }
    }
    if m.GetTables() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetTables())
        err = writer.WriteCollectionOfObjectValues("tables", cast)
        if err != nil {
            return err
        }
    }
    if m.GetWorksheets() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetWorksheets())
        err = writer.WriteCollectionOfObjectValues("worksheets", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetApplication sets the application property value. The application property
func (m *Workbook) SetApplication(value WorkbookApplicationable)() {
    m.application = value
}
// SetComments sets the comments property value. The comments property
func (m *Workbook) SetComments(value []WorkbookCommentable)() {
    m.comments = value
}
// SetFunctions sets the functions property value. The functions property
func (m *Workbook) SetFunctions(value WorkbookFunctionsable)() {
    m.functions = value
}
// SetNames sets the names property value. Represents a collection of workbooks scoped named items (named ranges and constants). Read-only.
func (m *Workbook) SetNames(value []WorkbookNamedItemable)() {
    m.names = value
}
// SetOperations sets the operations property value. The status of workbook operations. Getting an operation collection is not supported, but you can get the status of a long-running operation if the Location header is returned in the response. Read-only.
func (m *Workbook) SetOperations(value []WorkbookOperationable)() {
    m.operations = value
}
// SetTables sets the tables property value. Represents a collection of tables associated with the workbook. Read-only.
func (m *Workbook) SetTables(value []WorkbookTableable)() {
    m.tables = value
}
// SetWorksheets sets the worksheets property value. Represents a collection of worksheets associated with the workbook. Read-only.
func (m *Workbook) SetWorksheets(value []WorkbookWorksheetable)() {
    m.worksheets = value
}
