package agent

import "github.com/openshift-agent-team/fleeting/pkg/agent/isosource"

func BuildImage() error {
	_, err := isosource.EnsureIso()
	if err != nil {
		return err
	}

	return nil
}
