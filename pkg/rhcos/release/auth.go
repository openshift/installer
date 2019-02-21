package release

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/containers/image/types"
	"github.com/docker/distribution/reference"
	"github.com/pkg/errors"
)

type auth struct {
	Auth string `json:"auth"`
}

type config struct {
	Auths map[string]auth `json:"auths"`
}

func addPullSecret(sys *types.SystemContext, pullSecret []byte, named reference.Named) error {
	var auths config
	err := json.Unmarshal(pullSecret, &auths)
	if err != nil {
		return err
	}

	authority := reference.Domain(named)
	auth, ok := auths.Auths[authority] // hack: skipping normalization
	if ok {
		decoded, err := base64.StdEncoding.DecodeString(auth.Auth)
		if err != nil {
			return err
		}

		parts := strings.SplitN(string(decoded), ":", 2)
		if len(parts) != 2 {
			return errors.Errorf("invalid pull-secret entry for %q", authority)
		}

		sys.DockerAuthConfig = &types.DockerAuthConfig{
			Username: parts[0],
			Password: parts[1],
		}
	}

	return nil
}
