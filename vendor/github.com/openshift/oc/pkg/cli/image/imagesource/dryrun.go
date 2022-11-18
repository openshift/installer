package imagesource

import (
	"errors"
	"net/http"

	"github.com/docker/distribution"
	"github.com/docker/distribution/reference"
	registryclient "github.com/docker/distribution/registry/client"
)

func NewDryRun(ref TypedImageReference) (distribution.Repository, error) {
	named, err := reference.WithName(ref.Ref.RepositoryName())
	if err != nil {
		return nil, err
	}
	return registryclient.NewRepository(named, ref.Ref.RegistryURL().String(), dryRunRoundTripper)
}

var dryRunRoundTripper = errorRoundTripper{errors.New("dry-run repository is not available")}

type errorRoundTripper struct {
	err error
}

func (rt errorRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, rt.err
}
