package ssh

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh/agent"

	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

// newAgent initializes an SSH Agent with the keys.
// If no keys are provided, it loads all the keys from the user's environment.
func newAgent(keyPaths []string) (agent.Agent, error) {
	keys, err := loadKeys(keyPaths)
	if err != nil {
		return nil, err
	}
	if len(keys) == 0 {
		return nil, errors.New("no keys found for SSH agent")
	}

	ag := agent.NewKeyring()
	var errs []error
	for idx := range keys {
		if err := ag.Add(agent.AddedKey{PrivateKey: keys[idx]}); err != nil {
			errs = append(errs, errors.Wrap(err, "failed to add key to agent"))
		}
	}
	if agg := utilerrors.NewAggregate(errs); agg != nil {
		return nil, agg
	}
	return ag, nil
}

func loadKeys(paths []string) ([]interface{}, error) {
	if len(paths) > 0 {
		return LoadPrivateSSHKeys(paths)
	}
	return defaultPrivateSSHKeys()
}
