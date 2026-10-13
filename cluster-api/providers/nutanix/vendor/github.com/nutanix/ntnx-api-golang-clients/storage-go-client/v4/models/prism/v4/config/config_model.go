/*
 * Generated file models/prism/v4/config/config_model.go.
 *
 * Product version: 4.0.2-alpha-3
 *
 * Part of the Nutanix Storage Versioned APIs
 *
 * (c) 2023 Nutanix Inc.  All rights reserved
 *
 */

/*
  Module prism.v4.config of Nutanix Storage Versioned APIs
*/
package config

/**
Reference to a task tracking the async operation.
*/
type TaskReference struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/**
	  Globally unique identifier of a task.
	*/
	ExtId *string `json:"extId,omitempty"`
}

func NewTaskReference() *TaskReference {
	p := new(TaskReference)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "prism.v4.config.TaskReference"
	p.Reserved_ = map[string]interface{}{"$fqObjectType": "prism.v4.r0.a1.config.TaskReference"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}
