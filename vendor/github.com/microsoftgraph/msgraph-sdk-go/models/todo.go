package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Todo 
type Todo struct {
    Entity
    // The task lists in the users mailbox.
    lists []TodoTaskListable
}
// NewTodo instantiates a new todo and sets the default values.
func NewTodo()(*Todo) {
    m := &Todo{
        Entity: *NewEntity(),
    }
    return m
}
// CreateTodoFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateTodoFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTodo(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *Todo) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["lists"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateTodoTaskListFromDiscriminatorValue , m.SetLists)
    return res
}
// GetLists gets the lists property value. The task lists in the users mailbox.
func (m *Todo) GetLists()([]TodoTaskListable) {
    return m.lists
}
// Serialize serializes information the current object
func (m *Todo) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetLists() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetLists())
        err = writer.WriteCollectionOfObjectValues("lists", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetLists sets the lists property value. The task lists in the users mailbox.
func (m *Todo) SetLists(value []TodoTaskListable)() {
    m.lists = value
}
