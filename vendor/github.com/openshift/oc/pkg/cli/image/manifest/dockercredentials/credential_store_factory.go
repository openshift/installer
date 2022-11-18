package dockercredentials

import (
	"github.com/containers/image/v5/docker/reference"
	"github.com/docker/distribution/registry/client/auth"
	"github.com/docker/docker/api/types"
	dockerregistry "github.com/docker/docker/registry"

	"github.com/openshift/library-go/pkg/image/registryclient"
)

// NewCredentialStoreFactory returns an entity capable of creating a CredentialStore
func NewCredentialStoreFactory(path string) (registryclient.CredentialStoreFactory, error) {
	authResolver, err := NewAuthResolver(path)
	if err != nil {
		return nil, err
	}
	return &credentialStoreFactory{authResolver}, nil
}

type credentialStoreFactory struct {
	authResolver *AuthResolver
}

func (c *credentialStoreFactory) CredentialStoreFor(image string) auth.CredentialStore {
	nocreds := registryclient.NoCredentials
	if c.authResolver == nil {
		return nocreds
	}

	ref, err := reference.ParseNormalizedNamed(image)
	if err != nil {
		return nocreds
	}

	authCfg, err := c.authResolver.findAuthentication(ref, reference.Domain(ref))
	if err != nil {
		return nocreds
	}

	return dockerregistry.NewStaticCredentialStore(&types.AuthConfig{
		Username:      authCfg.Username,
		Password:      authCfg.Password,
		IdentityToken: authCfg.IdentityToken,
	})
}
