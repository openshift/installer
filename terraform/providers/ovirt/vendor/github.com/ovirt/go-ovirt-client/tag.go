package ovirtclient

import (
	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

//go:generate go run scripts/rest/rest.go -i "Tag" -n "tag" -T TagID

// TagID is the UUID of a tag.
type TagID string

// TagClient describes the functions related to oVirt tags.
type TagClient interface {
	// GetTag returns a single tag based on its ID.
	GetTag(id TagID, retries ...RetryStrategy) (Tag, error)
	// ListTags returns all tags on the oVirt engine.
	ListTags(retries ...RetryStrategy) ([]Tag, error)
	// CreateTag creates a new tag with a name.
	CreateTag(name string, params CreateTagParams, retries ...RetryStrategy) (result Tag, err error)
	// RemoveTag removes the tag with the specified ID.
	RemoveTag(tagID TagID, retries ...RetryStrategy) error
}

// TagData is the core of Tag, providing only the data access functions, but not the client
// functions.
type TagData interface {
	// ID returns the auto-generated identifier for this tag.
	ID() TagID
	// Name returns the user-give name for this tag.
	Name() string
	// Description returns the user-give description for this tag. It may be nil if no decription is set.
	Description() *string
}

// Tag is the interface defining the fields for tag.
type Tag interface {
	TagData
	Remove(retries ...RetryStrategy) error
}

// CreateTagParams contains the optional parameters for tag creation.
type CreateTagParams interface {
	Description() *string
}

// BuildableTagParams is an buildable version of CreateTagParams.
type BuildableTagParams interface {
	CreateTagParams

	WithDescription(description string) (BuildableTagParams, error)
	MustWithDescription(description string) BuildableTagParams
}

// NewCreateTagParams creates a buildable set of CreateTagParams to pass to the CreateTag function.
func NewCreateTagParams() BuildableTagParams {
	return &createTagParams{}
}

type createTagParams struct {
	description *string
}

func (c *createTagParams) WithDescription(description string) (BuildableTagParams, error) {
	c.description = &description
	return c, nil
}

func (c *createTagParams) MustWithDescription(description string) BuildableTagParams {
	builder, err := c.WithDescription(description)
	if err != nil {
		panic(err)
	}
	return builder
}

func (c *createTagParams) Description() *string {
	return c.description
}

func convertSDKTag(sdkObject *ovirtsdk4.Tag, client *oVirtClient) (Tag, error) {
	id, ok := sdkObject.Id()
	if !ok {
		return nil, newFieldNotFound("tag", "id")
	}
	name, ok := sdkObject.Name()
	if !ok {
		return nil, newFieldNotFound("tag", "name")
	}
	desc, ok := sdkObject.Description()
	var description *string
	if ok {
		description = &desc
	}
	return &tag{
		client:      client,
		id:          TagID(id),
		name:        name,
		description: description,
	}, nil
}

type tag struct {
	client      Client
	id          TagID
	name        string
	description *string
}

func (n tag) ID() TagID {
	return n.id
}

func (n tag) Name() string {
	return n.name
}

func (n tag) Description() *string {
	return n.description
}

func (n *tag) Remove(retries ...RetryStrategy) error {
	return n.client.RemoveTag(n.id, retries...)
}
