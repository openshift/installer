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
	tags "google.golang.org/api/cloudresourcemanager/v3"
	"google.golang.org/api/option"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openshift/installer/pkg/types/gcp"
)

//go:generate mockgen -source=./usertags.go -destination=./mock/usertags_mock.go -package=mock

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

// tagManager handles resource tagging.
type tagManager struct {
	client API
}

// TagManager is the interface that wraps methods for resource tag operations.
type TagManager interface {
	GetProjectTags(ctx context.Context, projectID string) (sets.Set[string], error)
	GetNamespacedTagValue(ctx context.Context, tagNamespacedName string) (*tags.TagValue, error)
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

// NewTagManager creates a tagManager instance.
func NewTagManager(client API) TagManager {
	return &tagManager{client: client}
}

// GetUserTags returns the processed list of user provided tags if already available,
// else validates, persists in-memory and returns the processed tags.
func GetUserTags(ctx context.Context, mgr TagManager, projectID string, userTags []gcp.UserTag) (map[string]string, error) {
	if !processedTags.isProcessed() {
		if err := validateAndPersistUserTags(ctx, mgr, projectID, userTags); err != nil {
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
func validateAndPersistUserTags(ctx context.Context, mgr TagManager, project string, userTags []gcp.UserTag) error {
	if len(userTags) == 0 {
		return nil
	}

	start := time.Now()
	logrus.Debugf("user defined tags list: %v", userTags)

	if len(userTags) > maxUserTagLimit {
		return fmt.Errorf("more than %d user tags is not allowed, configured count: %d", maxUserTagLimit, len(userTags))
	}

	projectTags, err := mgr.GetProjectTags(ctx, project)
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
		tagValue, err := mgr.GetNamespacedTagValue(ctx, name)
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

// getCloudResourceServiceForTags returns the client required for querying resource manager resources.
func (m *tagManager) getCloudResourceServiceForTags(ctx context.Context) (*tags.Service, error) {
	svc, err := tags.NewService(ctx, option.WithCredentials(m.client.GetCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create cloud resource service: %w", err)
	}
	return svc, nil
}

// GetProjectTags returns the list of effective tags attached to the provided project resource.
func (m *tagManager) GetProjectTags(ctx context.Context, projectID string) (sets.Set[string], error) {
	const (
		// projectParentPathFmt is the format string for parent path of a project resource.
		projectParentPathFmt = "//cloudresourcemanager.googleapis.com/projects/%s"
	)

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	service, err := m.getCloudResourceServiceForTags(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create cloud resource service: %w", err)
	}

	effectiveTags := sets.New[string]()
	effectiveTagsService := tags.NewEffectiveTagsService(service)
	effectiveTagsRequest := effectiveTagsService.List().Context(ctx).Parent(fmt.Sprintf(projectParentPathFmt, projectID))
	if err := effectiveTagsRequest.Pages(ctx, func(page *tags.ListEffectiveTagsResponse) error {
		for _, effectiveTag := range page.EffectiveTags {
			effectiveTags.Insert(effectiveTag.NamespacedTagValue)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return effectiveTags, nil
}

// GetNamespacedTagValue returns the Tag Value metadata fetched using the tag's NamespacedName.
func (m *tagManager) GetNamespacedTagValue(ctx context.Context, tagNamespacedName string) (*tags.TagValue, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	service, err := m.getCloudResourceServiceForTags(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create cloud resource service: %w", err)
	}

	tagValuesService := tags.NewTagValuesService(service)

	return tagValuesService.GetNamespaced().
		Context(ctx).
		Name(tagNamespacedName).
		Do()
}
