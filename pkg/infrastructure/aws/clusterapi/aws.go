package clusterapi

import (
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	awstypes "github.com/openshift/installer/pkg/types/aws"
)

var _ clusterapi.Provider = (*Provider)(nil)

type Provider struct{}

func (_ *Provider) Name() string { return awstypes.Name }
