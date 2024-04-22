package imagesource

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/openshift/library-go/pkg/image/reference"
)

// ErrAlreadyExists may be returned by the blob Create function to indicate that the blob already exists.
var ErrAlreadyExists = fmt.Errorf("blob already exists in the target location")

type DestinationType string

var (
	DestinationRegistry DestinationType = "docker"
	DestinationS3       DestinationType = "s3"
	DestinationFile     DestinationType = "file"
)

func (t DestinationType) Prefix() string {
	switch t {
	case DestinationFile:
		return "file://"
	case DestinationS3:
		return "s3://"
	default:
		return ""
	}
}

type TypedImageReference struct {
	Type DestinationType
	Ref  reference.DockerImageReference
}

func (t TypedImageReference) EqualRegistry(other TypedImageReference) bool {
	return t.Type == other.Type && t.Ref.Registry == other.Ref.Registry
}

func (t TypedImageReference) Equal(other TypedImageReference) bool {
	return t.Type == other.Type && t.Ref.Equal(other.Ref)
}

func (t TypedImageReference) String() string {
	switch t.Type {
	case DestinationFile:
		return fmt.Sprintf("file://%s", t.Ref.Exact())
	case DestinationS3:
		return fmt.Sprintf("s3://%s", t.Ref.Exact())
	default:
		return t.Ref.Exact()
	}
}

var rePossibleExpandableReference = regexp.MustCompile(`:([\w\d-\.*]+)$`)

func ParseSourceReference(ref string, expandFn func(ref TypedImageReference) ([]TypedImageReference, error)) ([]TypedImageReference, error) {
	if m := rePossibleExpandableReference.FindStringSubmatch(ref); len(m) == 2 && strings.Contains(m[1], "*") && expandFn != nil {
		subst := rePossibleExpandableReference.ReplaceAllString(ref, ":tag")
		src, err := ParseReference(subst)
		if err != nil {
			return nil, err
		}
		if src.Ref.Tag != "tag" {
			return nil, fmt.Errorf("source expansion is only possible in tags")
		}
		src.Ref.Tag = m[1]
		return expandFn(src)
	}
	src, err := ParseReference(ref)
	if err != nil {
		return nil, err
	}
	if len(src.Ref.Tag) == 0 && len(src.Ref.ID) == 0 {
		if expandFn == nil {
			return nil, fmt.Errorf("source references must have a tag or digest specified")
		}
		return expandFn(src)
	}
	return []TypedImageReference{src}, nil
}

func ParseDestinationReference(ref string) (TypedImageReference, error) {
	dst, err := ParseReference(ref)
	if err != nil {
		return dst, err
	}
	if len(dst.Ref.ID) != 0 {
		return dst, fmt.Errorf("you must specify a tag for DST or leave it blank to only push by digest")
	}
	return dst, err
}

func ParseReference(ref string) (TypedImageReference, error) {
	dstType := DestinationRegistry
	switch {
	case strings.HasPrefix(ref, "s3://"):
		dstType = DestinationS3
		ref = strings.TrimPrefix(ref, "s3://")
	case strings.HasPrefix(ref, "file://"):
		dstType = DestinationFile
		ref = strings.TrimPrefix(ref, "file://")
		if strings.HasPrefix(ref, "/") {
			ref = ref[1:]
		}
	}
	dst, err := reference.Parse(ref)
	if err != nil {
		return TypedImageReference{Ref: dst, Type: dstType}, fmt.Errorf("%q is not a valid image reference: %v", ref, err)
	}
	return TypedImageReference{Ref: dst, Type: dstType}, nil
}

// buildTagSearchRegexp creates a regexp from the provided tag value
// that can be used to filter tags. It supports standard '*' glob
// rules.
func buildTagSearchRegexp(tag string) (*regexp.Regexp, error) {
	search := tag
	if (len(search)) == 0 {
		search = "*"
	}
	var parts []string
	for _, part := range strings.Split(search, "*") {
		if len(part) == 0 {
			if len(parts) == 0 || parts[len(parts)-1] != ".*" {
				parts = append(parts, ".*")
			}
		} else {
			parts = append(parts, regexp.QuoteMeta(part))
		}
	}
	reText := fmt.Sprintf("^%s$", strings.Join(parts, ".*"))
	return regexp.Compile(reText)
}
