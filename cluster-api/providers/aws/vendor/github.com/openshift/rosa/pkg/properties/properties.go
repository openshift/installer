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

// This file contains the OCM properties that are used to store additional information about crated
// clusters.

package properties

// Prefix used by all the property names:
const prefix = "rosa_"

const CLIVersion = prefix + "cli_version"

const FakeCluster = "fake_cluster"

// nolint:gosec // Linter thinks there are hardcoded credentials here...
const UseLocalCredentials = "use_local_credentials"

const ProvisionShardId = "provision_shard_id"

const KeyringEnvKey = "OCM_KEYRING"
