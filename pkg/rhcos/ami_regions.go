package rhcos

import (
	"context"

	"github.com/coreos/stream-metadata-go/arch"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openshift/installer/pkg/types"
)

// AMIRegions returns the AWS regions in which an RHCOS AMI for the specified architecture is published.
func AMIRegions(architecture types.Architecture) sets.String {
	stream, err := FetchCoreOSBuild(context.Background())
	if err != nil {
		logrus.Errorf("could not fetch the rhcos stream data: %v", err)
		return nil
	}
	rpmArch := arch.RpmArch(string(architecture))
	awsImages := stream.Architectures[rpmArch].Images.Aws
	if awsImages == nil {
		return nil
	}
	regions := make([]string, 0, len(awsImages.Regions))
	for name, r := range awsImages.Regions {
		if r.Image == "" {
			continue
		}
		regions = append(regions, name)
	}
	return sets.NewString(regions...)
}
