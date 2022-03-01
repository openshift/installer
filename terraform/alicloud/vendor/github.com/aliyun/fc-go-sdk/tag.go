package fc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	tagPath  = "/tag"
	tagsPath = "/tags"
)

// TagResourceInput defines input to tag a resource(service)
type TagResourceInput struct {
	ResourceArn *string           `json:"resourceArn"`
	Tags        map[string]string `json:"tags"`
}

// NewTagResourceInput ...
func NewTagResourceInput(resourceArn string, tags map[string]string) *TagResourceInput {
	return &TagResourceInput{ResourceArn: &resourceArn, Tags: tags}
}

func (s *TagResourceInput) WithResourceArn(resourceArn string) *TagResourceInput {
	s.ResourceArn = &resourceArn
	return s
}

func (s *TagResourceInput) WithTags(tags map[string]string) *TagResourceInput {
	s.Tags = tags
	return s
}

func (t *TagResourceInput) GetQueryParams() url.Values {
	out := url.Values{}
	return out
}

func (t *TagResourceInput) GetPath() string {
	return tagPath
}

func (t *TagResourceInput) GetHeaders() Header {
	return make(Header, 0)
}

func (t *TagResourceInput) GetPayload() interface{} {
	return t
}

func (t *TagResourceInput) Validate() error {
	if IsBlank(t.ResourceArn) {
		return fmt.Errorf("ResourceArn is required but not provided")
	}

	if t.Tags == nil || len(t.Tags) == 0 {
		return fmt.Errorf("At least 1 tag is required")
	}
	return nil
}

// TagResourceOut ...
type TagResourceOut struct {
	Header http.Header
}

func (o TagResourceOut) String() string {
	b, err := json.MarshalIndent(o, "", printIndent)
	if err != nil {
		return ""
	}
	return string(b)
}

func (o TagResourceOut) GetRequestID() string {
	return GetRequestID(o.Header)
}

// GetResourceTagsInput ...
type GetResourceTagsInput struct {
	ResourceArn *string `json:"resourceArn"`
}

// NewGetResourceTagsInput...
func NewGetResourceTagsInput(resourceArn string) *GetResourceTagsInput {
	return &GetResourceTagsInput{ResourceArn: &resourceArn}
}

func (s *GetResourceTagsInput) WithResourceArn(resourceArn string) *GetResourceTagsInput {
	s.ResourceArn = &resourceArn
	return s
}

func (t *GetResourceTagsInput) GetQueryParams() url.Values {
	out := url.Values{}
	if t.ResourceArn != nil {
		out.Set("resourceArn", *t.ResourceArn)
	}
	return out
}

func (t *GetResourceTagsInput) GetPath() string {
	return tagPath
}

func (t *GetResourceTagsInput) GetHeaders() Header {
	return make(Header, 0)
}

func (t *GetResourceTagsInput) GetPayload() interface{} {
	return nil
}

func (t *GetResourceTagsInput) Validate() error {
	if IsBlank(t.ResourceArn) {
		return fmt.Errorf("ResourceArn is required but not provided")
	}
	return nil
}

// GetResourceTagsOut ...
type GetResourceTagsOut struct {
	Header      http.Header
	ResourceArn *string           `json:"resourceArn"`
	Tags        map[string]string `json:"tags"`
}

func (o GetResourceTagsOut) String() string {
	b, err := json.MarshalIndent(o, "", printIndent)
	if err != nil {
		return ""
	}
	return string(b)
}

func (o GetResourceTagsOut) GetRequestID() string {
	return GetRequestID(o.Header)
}

// UnTagResourceInput ...
type UnTagResourceInput struct {
	ResourceArn *string  `json:"resourceArn"`
	TagKeys     []string `json:"tagKeys"`
	All         *bool    `json:"all"`
}

// NewUnTagResourceInput ...
func NewUnTagResourceInput(resourceArn string) *UnTagResourceInput {
	return &UnTagResourceInput{ResourceArn: &resourceArn}
}

func (s *UnTagResourceInput) WithResourceArn(resourceArn string) *UnTagResourceInput {
	s.ResourceArn = &resourceArn
	return s
}

func (s *UnTagResourceInput) WithTagKeys(tagKeys []string) *UnTagResourceInput {
	s.TagKeys = tagKeys
	return s
}

func (s *UnTagResourceInput) WithAll(all bool) *UnTagResourceInput {
	s.All = &all
	return s
}

func (t *UnTagResourceInput) GetQueryParams() url.Values {
	out := url.Values{}
	return out
}

func (t *UnTagResourceInput) GetPath() string {
	return tagPath
}

func (t *UnTagResourceInput) GetHeaders() Header {
	return make(Header, 0)
}

func (t *UnTagResourceInput) GetPayload() interface{} {
	return t
}

func (t *UnTagResourceInput) Validate() error {
	if IsBlank(t.ResourceArn) {
		return fmt.Errorf("ResourceArn is required but not provided")
	}

	all := true
	if t.All == nil {
		all = false
	} else {
		all = *t.All
	}

	if !all && len(t.TagKeys) == 0 {
		return fmt.Errorf("At least 1 tag key is required if all=false")
	}

	if len(t.TagKeys) > 20 {
		return fmt.Errorf("At most 20 tag is required")
	}
	return nil
}

// UnTagResourceOut ...
type UnTagResourceOut struct {
	Header http.Header
}

func (o UnTagResourceOut) String() string {
	b, err := json.MarshalIndent(o, "", printIndent)
	if err != nil {
		return ""
	}
	return string(b)
}

func (o UnTagResourceOut) GetRequestID() string {
	return GetRequestID(o.Header)
}
