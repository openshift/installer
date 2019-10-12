package aws

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"
)

// Metadata holds additional metadata for InstallConfig resources that
// does not need to be user-supplied (e.g. because it can be retrieved
// from external APIs).
type Metadata struct {
	session           *session.Session
	availabilityZones []string
	region            string
	mutex             sync.Mutex
}

// NewMetadata initializes a new Metadata object.
func NewMetadata(region string) *Metadata {
	return &Metadata{region: region}
}

// Session holds an AWS session which can be used for AWS API calls
// during asset generation.
func (m *Metadata) Session(ctx context.Context) (*session.Session, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.unlockedSession(ctx)
}

func (m *Metadata) unlockedSession(ctx context.Context) (*session.Session, error) {
	if m.session == nil {
		var err error
		m.session, err = GetSession()
		if err != nil {
			return nil, errors.Wrap(err, "creating AWS session")
		}
	}

	return m.session, nil
}

// AvailabilityZones retrieves a list of availability zones for the configured region.
func (m *Metadata) AvailabilityZones(ctx context.Context) ([]string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.availabilityZones) == 0 {
		session, err := m.unlockedSession(ctx)
		if err != nil {
			return nil, err
		}

		m.availabilityZones, err = availabilityZones(ctx, session, m.region)
		if err != nil {
			return nil, errors.Wrap(err, "creating AWS session")
		}
	}

	return m.availabilityZones, nil
}
