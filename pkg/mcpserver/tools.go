package mcpserver

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/rhcos"
)

func GetCoreOS() (string, error) {
	logrus.Info("Getting CoreOS stream data")
	streamData, err := rhcos.FetchRawCoreOSStream(context.Background())
	if err != nil {
		return "", err
	}
	return string(streamData), nil
}
