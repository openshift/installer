package example

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/destroy/providers"
	"github.com/openshift/installer/pkg/types"
)

type ClusterUninstaller struct {
	ClusterID string
	InfraID   string
	Logger    logrus.FieldLogger
	counter   int
}

func New(logger logrus.FieldLogger, metadata *types.ClusterMetadata) (providers.Destroyer, error) {
	return &ClusterUninstaller{
		ClusterID: "ExampleClusterID",
		InfraID:   "ExampleInfraID",
		Logger:    logger,
	}, nil
}

func (o *ClusterUninstaller) GetStages() []providers.Stage {
	return []providers.Stage{
		{
			Name:  "Init clients",
			Funcs: []providers.StageFunc{o.initClients},
		},
		{
			Name:  "Shutdown VMs",
			Funcs: []providers.StageFunc{o.shutdownVMs},
		},
		{
			Name:  "Destroy VMs",
			Funcs: []providers.StageFunc{o.destroyVMs},
		},
		{
			Name:  "Cleanup clients",
			Funcs: []providers.StageFunc{o.cleanUp},
		},
	}
}

func (o *ClusterUninstaller) initClients(ctx context.Context) (bool, error) {
	fmt.Println("Initializing cloud provider clients...")
	return true, nil
}

func (o *ClusterUninstaller) setup(ctx context.Context) (bool, error) {
	fmt.Println("Setting up other stages...")
	// Fail this call every other time. By returning `false` we're saying this
	// is a recoverable error and the stage should be retried
	if o.counter%2 == 0 {
		o.counter += 1
		return false, nil
	}
	return true, nil
}

func (o *ClusterUninstaller) shutdownVMs(ctx context.Context) (bool, error) {
	fmt.Println("Shutting down VMs...")
	return true, nil
}

func (o *ClusterUninstaller) destroyVMs(ctx context.Context) (bool, error) {
	fmt.Println("Destroying VMs...")
	return true, nil
}

func (o *ClusterUninstaller) cleanUp(ctx context.Context) (bool, error) {
	fmt.Println("Cleaning up clients...")
	return true, nil
}
