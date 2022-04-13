package agent

import (
	"github.com/openshift-agent-team/fleeting/pkg/agent/imagebuilder"
	"github.com/openshift-agent-team/fleeting/pkg/agent/isosource"
)

func BuildImage() error {

	baseImage, err := isosource.EnsureIso()
	if err != nil {
		return err
	}

	err = imagebuilder.BuildImage(baseImage)
	if err != nil {
		return err
	}

	return nil
}
