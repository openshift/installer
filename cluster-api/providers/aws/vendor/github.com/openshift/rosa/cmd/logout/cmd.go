/*
Copyright (c) 2020 Red Hat, Inc.

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

package logout

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/openshift/rosa/pkg/config"
	rprtr "github.com/openshift/rosa/pkg/reporter"
)

var Cmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out",
	Long:  "Log out, removing the configuration file.",
	Run:   run,
	Args:  cobra.NoArgs,
}

func run(_ *cobra.Command, _ []string) {
	reporter := rprtr.CreateReporter()
	// Remove the configuration file:
	err := config.Remove()
	if err != nil {
		reporter.Errorf("Failed to remove config file: %v", err)
		os.Exit(1)
	}
}
