/*
Copyright (c) 2021 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ocm

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/openshift/rosa/pkg/aws"
	"github.com/openshift/rosa/pkg/logging"
)

const (
	clusterFlagName        = "cluster"
	clusterFlagShortHand   = "c"
	clusterFlagDescription = "Name or ID of the cluster."
)

var clusterKey string

func AddOptionalClusterFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(
		&clusterKey,
		clusterFlagName,
		clusterFlagShortHand,
		"",
		clusterFlagDescription,
	)
	cmd.RegisterFlagCompletionFunc(clusterFlagName, clusterCompletion)
}

func AddClusterFlag(cmd *cobra.Command) {
	AddOptionalClusterFlag(cmd)
	cmd.MarkFlagRequired(clusterFlagName)
}

func SetClusterKey(key string) {
	clusterKey = key
}

func GetClusterKey() (string, error) {
	// Check that the cluster key (name, identifier or external identifier) given by the user
	// is reasonably safe so that there is no risk of SQL injection:
	if !IsValidClusterKey(clusterKey) {
		return "", fmt.Errorf(
			"Cluster name, identifier or external identifier '%s' isn't valid: it "+
				"must contain only letters, digits, dashes and underscores",
			clusterKey,
		)
	}
	return clusterKey, nil
}

func clusterCompletion(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	logger := logging.NewLogger()

	ocmClient, err := NewClient().Logger(logger).Build()
	if err != nil {
		return []string{}, cobra.ShellCompDirectiveDefault
	}
	defer ocmClient.Close()

	awsClient, err := aws.NewClient().Logger(logger).Build()
	if err != nil {
		return []string{}, cobra.ShellCompDirectiveDefault
	}
	awsCreator, err := awsClient.GetCreator()
	if err != nil {
		return []string{}, cobra.ShellCompDirectiveDefault
	}

	clusters, err := ocmClient.GetClusters(awsCreator, 10)
	if err != nil {
		return []string{}, cobra.ShellCompDirectiveDefault
	}
	res := []string{}
	for _, cluster := range clusters {
		res = append(res, cluster.Name())
	}
	return res, cobra.ShellCompDirectiveDefault
}
