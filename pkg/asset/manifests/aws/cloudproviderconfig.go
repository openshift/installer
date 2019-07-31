package aws

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	ini "gopkg.in/ini.v1"

	awstypes "github.com/openshift/installer/pkg/types/aws"
)

// cloudConfig is the aws cloud provider config.
type cloudConfig struct {
	Global global
}

// global struct of cloudConfig which is currently not inialized.
type global struct {
}

// serviceOverride struct used for AWS service endpoint override.
type serviceOverride struct {
	// Service string describing the AWS service being overriden,
	// Currently support ec2, elb, iam, s3 and route53.
	Service string `ini:"Service"`
	// Region where the service is overridden. In current implementation,
	// its set to the cluster default region.
	Region string `ini:"Region"`
	// URL descibes the endpoint which will be used for service override.
	URL string `ini:"URL"`
}

// CloudProviderConfig builds the cloud provider config and reflects to an ini file.
func CloudProviderConfig(awsPlatform *awstypes.Platform) (string, error) {
	file := ini.Empty()
	config := &cloudConfig{
		Global: global{},
	}
	if err := file.ReflectFrom(config); err != nil {
		return "", errors.Wrap(err, "failed to reflect from config")
	}

	for index, endpoint := range awsPlatform.CustomRegionOverride {
		s, err := file.NewSection(fmt.Sprintf("ServiceOverride %q", strconv.Itoa(index)))
		if err != nil {
			return "", errors.Wrapf(err, "failed to create section for ServiceOverride")
		}
		if err := s.ReflectFrom(
			&serviceOverride{
				Service: endpoint.Service,
				Region:  awsPlatform.Region,
				URL:     endpoint.URL,
			}); err != nil {
			return "", errors.Wrapf(err, "failed to reflect from ServiceOverride")
		}
	}

	buf := &bytes.Buffer{}
	if _, err := file.WriteTo(buf); err != nil {
		return "", errors.Wrap(err, "failed to write out cloud provider config")
	}

	return buf.String(), nil
}
