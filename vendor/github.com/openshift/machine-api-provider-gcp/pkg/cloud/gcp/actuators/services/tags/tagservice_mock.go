package tagservice

import (
	"context"

	tags "google.golang.org/api/cloudresourcemanager/v3"
)

// MockTagService mocks TagService interface for tests.
type MockTagService struct {
	MockGetNamespacedName func(context.Context, string) (*tags.TagValue, error)
}

// NewMockTagService returns new mock of tagService.
func NewMockTagService() *MockTagService {
	return &MockTagService{}
}

// NewMockTagServiceBuilder returns new mock for creating GCP tag client.
func NewMockTagServiceBuilder(ctx context.Context, serviceAccountJSON string) (TagService, error) {
	return NewMockTagService(), nil
}

// GetNamespacedName returns mock metadata of the requested tag.
func (m *MockTagService) GetNamespacedName(ctx context.Context, name string) (*tags.TagValue, error) {
	if m.MockGetNamespacedName == nil {
		return nil, nil
	}
	return m.MockGetNamespacedName(ctx, name)
}
