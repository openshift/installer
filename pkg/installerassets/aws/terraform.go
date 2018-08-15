package aws

import (
	"context"
	"encoding/json"

	"github.com/ghodss/yaml"
	"github.com/openshift/installer/pkg/assets"
	"github.com/openshift/installer/pkg/installerassets"
	"github.com/pkg/errors"
)

type endpoints string

const (
	endpointsAll     endpoints = "all"
	endpointsPrivate endpoints = "private"
	endpointsPublic  endpoints = "public"
)

type terraformConfig struct {
	AMI                string            `json:"tectonic_aws_ec2_ami_override,omitempty"`
	Endpoints          endpoints         `json:"tectonic_aws_endpoints,omitempty"`
	ExternalVPCID      string            `json:"tectonic_aws_external_vpc_id,omitempty"`
	InstallerRole      string            `json:"tectonic_aws_installer_role,omitempty"`
	MasterInstanceType string            `json:"tectonic_aws_master_ec2_type,omitempty"`
	Region             string            `json:"tectonic_aws_region,omitempty"`
	UserTags           map[string]string `json:"tectonic_aws_extra_tags,omitempty"`
	VPCCIDR            string            `json:"tectonic_aws_vpc_cidr_block,omitempty"`
}

func terraformRebuilder(ctx context.Context, getByName assets.GetByString) (*assets.Asset, error) {
	asset := &assets.Asset{
		Name:          "terraform/aws-terraform.auto.tfvars",
		RebuildHelper: terraformRebuilder,
	}

	parents, err := asset.GetParents(
		ctx,
		getByName,
		"aws/ami",
		"aws/external-vpc-id",
		"aws/instance-type",
		"aws/region",
		"aws/user-tags",
		"network/node-cidr",
	)
	if err != nil {
		return nil, err
	}

	var userTags map[string]string
	err = yaml.Unmarshal(parents["aws/user-tags"].Data, &userTags)
	if err != nil {
		return nil, errors.Wrap(err, "parse user tags")
	}

	config := &terraformConfig{
		AMI:                string(parents["aws/ami"].Data),
		Endpoints:          endpointsAll,
		ExternalVPCID:      string(parents["aws/external-vpc-id"].Data),
		MasterInstanceType: string(parents["aws/instance-type"].Data),
		Region:             string(parents["aws/region"].Data),
		UserTags:           userTags,
		VPCCIDR:            string(parents["network/node-cidr"].Data),
	}

	asset.Data, err = json.Marshal(config)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

func init() {
	installerassets.Rebuilders["terraform/aws-terraform.auto.tfvars"] = terraformRebuilder
}
