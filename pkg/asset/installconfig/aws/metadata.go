package aws

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"

	typesaws "github.com/openshift/installer/pkg/types/aws"
)

// Metadata holds additional metadata for InstallConfig resources that
// does not need to be user-supplied (e.g. because it can be retrieved
// from external APIs).
type Metadata struct {
	session           *session.Session
	config            *aws.Config
	availabilityZones []string
	privateSubnets    map[string]Subnet
	publicSubnets     map[string]Subnet
	vpc               string
	instanceTypes     map[string]InstanceType

	Region   string                     `json:"region,omitempty"`
	Subnets  []string                   `json:"subnets,omitempty"`
	Services []typesaws.ServiceEndpoint `json:"services,omitempty"`

	mutex sync.Mutex
}

// NewMetadata initializes a new Metadata object.
func NewMetadata(region string, subnets []string, services []typesaws.ServiceEndpoint) *Metadata {
	return &Metadata{Region: region, Subnets: subnets, Services: services}
}

// Session holds an AWS session which can be used for AWS API calls
// during asset generation.
func (m *Metadata) Session(ctx context.Context) (*session.Session, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.unlockedSession(ctx)
}

func (m *Metadata) unlockedSession(ctx context.Context) (*session.Session, error) {
	if m.session == nil {
		var err error
		m.session, err = GetSessionWithOptions(WithRegion(m.Region), WithServiceEndpoints(m.Region, m.Services))
		if err != nil {
			return nil, errors.Wrap(err, "creating AWS session")
		}
	}

	return m.session, nil
}

// Config holds an AWS config which can be used for AWS API calls
// during asset generation.
func (m *Metadata) Config(ctx context.Context) (*aws.Config, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.unlockedConfig(ctx)
}

func (m *Metadata) unlockedConfig(ctx context.Context) (*aws.Config, error) {
	if m.config == nil {
		cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(m.Region), WithCustomEndpointsResolver(m.Region, m.Services))
		if err != nil {
			return nil, errors.Wrap(err, "loading AWS default config")
		}
		m.config = &cfg
	}

	return m.config, nil
}

// AvailabilityZones retrieves a list of availability zones for the configured region.
func (m *Metadata) AvailabilityZones(ctx context.Context) ([]string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.availabilityZones) == 0 {
		config, err := m.unlockedConfig(ctx)
		if err != nil {
			return nil, err
		}

		m.availabilityZones, err = availabilityZones(ctx, *config, m.Region)
		if err != nil {
			return nil, err
		}
	}

	return m.availabilityZones, nil
}

// PrivateSubnets retrieves subnet metadata indexed by subnet ID, for
// subnets that the cloud-provider logic considers to be private
// (i.e. not public).
func (m *Metadata) PrivateSubnets(ctx context.Context) (map[string]Subnet, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	err := m.populateSubnets(ctx)
	if err != nil {
		return nil, err
	}

	return m.privateSubnets, nil
}

// PublicSubnets retrieves subnet metadata indexed by subnet ID, for
// subnets that the cloud-provider logic considers to be public
// (e.g. with suitable routing for hosting public load balancers).
func (m *Metadata) PublicSubnets(ctx context.Context) (map[string]Subnet, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	err := m.populateSubnets(ctx)
	if err != nil {
		return nil, err
	}

	return m.publicSubnets, nil
}

func (m *Metadata) populateSubnets(ctx context.Context) error {
	if len(m.publicSubnets) > 0 || len(m.privateSubnets) > 0 {
		return nil
	}

	if len(m.Subnets) == 0 {
		return errors.New("no subnets configured")
	}

	config, err := m.unlockedConfig(ctx)
	if err != nil {
		return err
	}

	m.vpc, m.privateSubnets, m.publicSubnets, err = subnets(ctx, *config, m.Subnets)
	return err
}

// VPC retrieves the VPC ID containing PublicSubnets and PrivateSubnets.
func (m *Metadata) VPC(ctx context.Context) (string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.vpc == "" {
		if len(m.Subnets) == 0 {
			return "", errors.New("cannot calculate VPC without configured subnets")
		}

		err := m.populateSubnets(ctx)
		if err != nil {
			return "", err
		}
	}

	return m.vpc, nil
}

// InstanceTypes retrieves instance type metadata indexed by InstanceType for the configured region.
func (m *Metadata) InstanceTypes(ctx context.Context) (map[string]InstanceType, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.instanceTypes) == 0 {
		config, err := m.unlockedConfig(ctx)
		if err != nil {
			return nil, err
		}

		m.instanceTypes, err = instanceTypes(ctx, *config)
		if err != nil {
			return nil, errors.Wrap(err, "listing instance types")
		}
	}

	return m.instanceTypes, nil
}
