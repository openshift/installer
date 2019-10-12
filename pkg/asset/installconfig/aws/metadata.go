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
	session *session.Session
	mutex   sync.Mutex
}

// Session holds an AWS session which can be used for AWS API calls
// during asset generation.
// GetMetadata loads metadata.
func (m *Metadata) Session(ctx context.Context) (*session.Session, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.session == nil {
		var err error
		m.session, err = GetSession()
		if err != nil {
			return nil, errors.Wrap(err, "creating AWS session")
		}
	}

	return m.session, nil
}
