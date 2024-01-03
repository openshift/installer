package clusterapi

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/installconfig/aws"
)

func createDNSRecords(cluster *clusterv1.Cluster, installConfig *installconfig.InstallConfig) error {
	logrus.Infof("Creating Route53 records for control plane load balancer")
	ssn, err := installConfig.AWS.Session(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	client := aws.NewClient(ssn)
	r53cfg := aws.GetR53ClientCfg(ssn, installConfig.Config.AWS.HostedZoneRole)
	err = client.CreateOrUpdateRecord(installConfig.Config, cluster.Spec.ControlPlaneEndpoint.Host, r53cfg)
	if err != nil {
		return fmt.Errorf("failed to create route53 records: %w", err)
	}
	logrus.Infof("Created Route53 records for control plane load balancer")
	return nil
}
