// This build script is invoked as part of the container build to output the
// required platform permissions data into a file.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	iamv1 "sigs.k8s.io/cluster-api-provider-aws/v2/iam/api/v1beta1"

	icaws "github.com/openshift/installer/pkg/asset/installconfig/aws"
)

const awsPermissionsPolicyFilename = "bin/manifests/installer-aws-policy.json"

func run() error {
	permissions, err := icaws.PermissionsList(icaws.AllPermissionGroups())
	if err != nil {
		return fmt.Errorf("failed to generated permissions list: %w", err)
	}

	policy := &iamv1.PolicyDocument{
		Version: "2012-10-17",
		Statement: []iamv1.StatementEntry{
			{
				Effect:   "Allow",
				Resource: iamv1.Resources{"*"},
				Action:   permissions,
			},
		},
	}

	b, err := json.Marshal(policy)
	if err != nil {
		return fmt.Errorf("failed to marshal AWS permissions policy: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(awsPermissionsPolicyFilename), 0755); err != nil {
		return err
	}

	err = os.WriteFile(awsPermissionsPolicyFilename, b, 0o644) //#nosec G306 -- no sensitive data
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
