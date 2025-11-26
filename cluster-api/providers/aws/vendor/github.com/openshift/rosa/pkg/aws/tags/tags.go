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

// This file contains the tags that are used to store additional information in objects created
// in AWS.

package tags

// Prefix used by all the tag names:
const prefix = "rosa_"

// ClusterName is the name of the tag that will contain the name of the cluster.
const ClusterName = prefix + "cluster_name"

// ClusterID is the name of the tag that will contain the identifier of the cluster.
const ClusterID = prefix + "cluster_id"

// ClusterID is the name of the tag that will contain the identifier of the cluster.
const ClusterRegion = prefix + "region"

// RoleType is the name of the tag that will contain the purpose of the role (installer, support, etc.)
const RoleType = prefix + "role_type"

// RolePrefix is the name of the tag that will contain the user-set prefix of the role (installer, support, etc.)
const RolePrefix = prefix + "role_prefix"

// Environment is the name of the tag that will contain the environment of the role (integration/staging/production)
const Environment = prefix + "environment"

// AdminRole tags the role as admin (true/false)
const AdminRole = prefix + "admin_role"

// RedHatManaged tags the role as red_hat_managed
const RedHatManaged = "red-hat-managed"

// HcpSharedVpc tags are for resources related to HCP shared VPC
const HcpSharedVpc = "hcp-shared-vpc"

const HypershiftPolicies = prefix + "hcp_policies"

const OperatorNamespace = "operator_namespace"

const OperatorName = "operator_name"

const InUse = "in_use"

const True = "true"
