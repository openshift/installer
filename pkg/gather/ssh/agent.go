package ssh

import (
	"net"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/agent"

	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

// getAgent attempts to connect to the running SSH agent, returning a newly
// initialized static agent if that fails.
func getAgent(keys []string) (agent.Agent, error) {
	// Attempt to use the existing SSH agent if it's configured and no keys
	// were explicitly passed.
	if authSock := os.Getenv("SSH_AUTH_SOCK"); authSock != "" && len(keys) == 0 {
		logrus.Debugf("Using SSH_AUTH_SOCK %s to connect to an existing agent", authSock)
		if conn, err := net.Dial("unix", authSock); err == nil {
			return agent.NewClient(conn), nil
		}
	}

	return newAgent(keys)
}

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
	for name, key := range keys {
		if err := ag.Add(agent.AddedKey{PrivateKey: key}); err != nil {
			errs = append(errs, errors.Wrapf(err, "failed to add %s to agent", name))
		}
		logrus.Debugf("Added %s to installer's internal agent", name)
	}
	if agg := utilerrors.NewAggregate(errs); agg != nil {
		return nil, agg
	}
	return ag, nil
}

func loadKeys(paths []string) (map[string]interface{}, error) {
	keys := map[string]interface{}{}
	if len(paths) > 0 {
		pkeys, err := LoadPrivateSSHKeys(paths)
		if err != nil {
			return nil, err
		}
		for k, v := range pkeys {
			keys[k] = v
		}
	}
	dkeys, err := defaultPrivateSSHKeys()
	if err != nil && len(paths) == 0 {
		return nil, err
	}
	for k, v := range dkeys {
		keys[k] = v
	}
	return keys, nil
}
