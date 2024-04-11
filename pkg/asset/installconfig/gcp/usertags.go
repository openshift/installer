package gcp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/googleapis/gax-go/v2/apierror"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openshift/installer/pkg/types/gcp"
)

const (
	// maxUserTagLimit is the maximum userTags that can be configured as defined in openshift/api.
	// https://github.com/openshift/api/commit/ae73a19d05c35068af16c9aeff375d0b7c936a8a#diff-07b264a49084976b670fb699badaca1795027d6ea732a99226a5388104f6174fR604-R613
	maxUserTagLimit = 50
)

// processedTags is the global instance for storing the processed tags.
var processedTags = newProcessedUserTags()

// processedUserTags is for storing the validated and processed tags
// in-memory for later use. Validating and converting the tags to
// required format involves using same GCP APIs and to reduce the
// number of calls, making use of the in-memory storage approach.
type processedUserTags struct {
	processed bool
	tags      map[string]string

	sync.Mutex
}

// TagManager handles resource tagging.
type TagManager struct {
	client API
}

// newProcessedUserTags is for initializing an instance of processedUserTags.
func newProcessedUserTags() *processedUserTags {
	return &processedUserTags{}
}

// make is for initializing the internals of processedUserTags.
func (p *processedUserTags) make(size int) {
	p.Lock()
	defer p.Unlock()

	p.processed = false
	p.tags = make(map[string]string, size)
}

// addTag is for adding a tag-set to processedUserTags.
func (p *processedUserTags) addTag(k, v string) {
	p.Lock()
	defer p.Unlock()

	p.tags[k] = v
}

// markProcessed is for updating the tags processed status.
func (p *processedUserTags) markProcessed() {
	p.Lock()
	defer p.Unlock()

	p.processed = true
}

// isProcessed is for getting the tags processed status.
func (p *processedUserTags) isProcessed() bool {
	p.Lock()
	defer p.Unlock()

	return p.processed
}

// copy is for making copy of the tags in processedUserTags.
func (p *processedUserTags) copy() map[string]string {
	p.Lock()
	defer p.Unlock()

	t := make(map[string]string, len(p.tags))
	for k, v := range p.tags {
		t[k] = v
	}
	return t
}

// NewTagManager creates a TagManager instance.
func NewTagManager(client API) *TagManager {
	return &TagManager{client: client}
}

// GetUserTags returns the processed list of user provided tags if already available,
// else validates, persists in-memory and returns the processed tags.
func (t *TagManager) GetUserTags(ctx context.Context, projectID string, userTags []gcp.UserTag) (map[string]string, error) {
	if !processedTags.isProcessed() {
		if err := t.validateAndPersistUserTags(ctx, projectID, userTags); err != nil {
			return nil, err
		}
	}

	return processedTags.copy(), nil
}

// validateAndPersistUserTags validates user provided tags are accessible and exists.
// Converts and persists user-defined tags in NamespacedName
// `{parentID}/{tag_key_short_name}/{tag_value_short_name}` form to key:value pairs,
// with key of the form `tagKeys/{tag_key_id}` and value of the form
// `tagValues/{tag_value_id}`. Returns error when fetching a tag fails or when
// tag already exists on the project resource.
func (t *TagManager) validateAndPersistUserTags(ctx context.Context, project string, userTags []gcp.UserTag) error {
	if len(userTags) == 0 {
		return nil
	}

	start := time.Now()
	logrus.Debugf("user defined tags list: %v", userTags)

	if len(userTags) > maxUserTagLimit {
		return fmt.Errorf("more than %d user tags is not allowed, configured count: %d", maxUserTagLimit, len(userTags))
	}

	projectTags, err := t.client.GetProjectTags(ctx, project)
	if err != nil {
		return err
	}

	if dupTags := findDuplicateTags(userTags, projectTags); len(dupTags) != 0 {
		return fmt.Errorf("found duplicate tags, %v tags already exist on %s project resource", dupTags, project)
	}

	processedTags.make(len(userTags))
	nonexistentTags := make([]string, 0)
	for _, tag := range userTags {
		name := fmt.Sprintf("%s/%s/%s", tag.ParentID, tag.Key, tag.Value)
		tagValue, err := t.client.GetNamespacedTagValue(ctx, name)
		if err != nil {
			// check and return all non-existing tags at once
			// for user to fix in one go.
			var gErr *apierror.APIError
			// google API returns StatusForbidden or StatusNotFound when the tag
			// does not exist, since it could be because of permission issues
			// or genuinely tag does not exist.
			if errors.As(err, &gErr) && gErr.HTTPCode() == http.StatusForbidden {
				logrus.Debugf("does not have permission to access %s tag or does not exist: %d", name, gErr.HTTPCode())
				nonexistentTags = append(nonexistentTags, name)
				continue
			}
			return fmt.Errorf("failed to fetch user-defined tag %s(%d): %w", name, gErr.HTTPCode(), err)
		}
		processedTags.addTag(tagValue.Parent, tagValue.Name)
	}

	if len(nonexistentTags) != 0 {
		return fmt.Errorf("does not have permission to access %v tag(s) or does not exist", nonexistentTags)
	}

	processedTags.markProcessed()
	logrus.Debugf("user defined tags processed list: %v, took %s to finish", processedTags.tags, time.Since(start))
	return nil
}

// findDuplicateTags returns list of duplicate userTags which are already present
// in the parentTags list.
func findDuplicateTags(userTags []gcp.UserTag, parentTags sets.Set[string]) []string {
	dupTags := make([]string, 0)
	for _, tag := range userTags {
		tagNamespacedName := fmt.Sprintf("%s/%s/%s", tag.ParentID, tag.Key, tag.Value)
		if parentTags.Has(tagNamespacedName) {
			dupTags = append(dupTags, tagNamespacedName)
		}
	}
	return dupTags
}
