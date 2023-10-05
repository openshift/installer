/*
Copyright 2020 The Kubernetes Authors.

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

package cloud

import (
	awsclient "github.com/aws/aws-sdk-go/aws/client"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/throttle"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
)

// Session represents an AWS session.
type Session interface {
	Session() awsclient.ConfigProvider
	ServiceLimiter(string) *throttle.ServiceLimiter
}

// ScopeUsage is used to indicate which controller is using a scope.
type ScopeUsage interface {
	// ControllerName returns the name of the controller that created the scope
	ControllerName() string
}

// ClusterObject represents a AWS cluster object.
type ClusterObject interface {
	conditions.Setter
}

// Logger represents the ability to log messages, both errors and not.
type Logger interface {
	// Enabled tests whether this Logger is enabled.  For example, commandline
	// flags might be used to set the logging verbosity and disable some info
	// logs.
	Enabled() bool

	// Info logs a non-error message with the given key/value pairs as context.
	//
	// The msg argument should be used to add some constant description to
	// the log line.  The key/value pairs can then be used to add additional
	// variable information.  The key/value pairs should alternate string
	// keys and arbitrary values.
	Info(msg string, keysAndValues ...interface{})

	// Error logs an error, with the given message and key/value pairs as context.
	// It functions similarly to calling Info with the "error" named value, but may
	// have unique behavior, and should be preferred for logging errors (see the
	// package documentations for more information).
	//
	// The msg field should be used to add context to any underlying error,
	// while the err field should be used to attach the actual error that
	// triggered this log line, if present.
	Error(err error, msg string, keysAndValues ...interface{})

	// V returns a Logger value for a specific verbosity level, relative to
	// this Logger.  In other words, V values are additive.  V higher verbosity
	// level means a log message is less important.  It's illegal to pass a log
	// level less than zero.
	V(level int) logr.Logger

	// WithValues adds some key-value pairs of context to a logger.
	// See Info for documentation on how key/value pairs work.
	WithValues(keysAndValues ...interface{}) logr.Logger

	// WithName adds a new element to the logger's name.
	// Successive calls with WithName continue to append
	// suffixes to the logger's name.  It's strongly recommended
	// that name segments contain only letters, digits, and hyphens
	// (see the package documentation for more information).
	WithName(name string) logr.Logger
}

// ClusterScoper is the interface for a cluster scope.
type ClusterScoper interface {
	Logger
	Session
	ScopeUsage

	// Name returns the CAPI cluster name.
	Name() string
	// Namespace returns the cluster namespace.
	Namespace() string
	// AWSClusterName returns the AWS cluster name.
	InfraClusterName() string
	// Region returns the cluster region.
	Region() string
	// KubernetesClusterName is the name of the Kubernetes cluster. For EKS this
	// will differ to the CAPI cluster name
	KubernetesClusterName() string

	// InfraCluster returns the AWS infrastructure cluster object.
	InfraCluster() ClusterObject

	// Cluster returns the cluster object.
	ClusterObj() ClusterObject

	// IdentityRef returns the AWS infrastructure cluster identityRef.
	IdentityRef() *infrav1.AWSIdentityReference

	// ListOptionsLabelSelector returns a ListOptions with a label selector for clusterName.
	ListOptionsLabelSelector() client.ListOption
	// APIServerPort returns the port to use when communicating with the API server.
	APIServerPort() int32
	// AdditionalTags returns any tags that you would like to attach to AWS resources. The returned value will never be nil.
	AdditionalTags() infrav1.Tags
	// SetFailureDomain sets the infrastructure provider failure domain key to the spec given as input.
	SetFailureDomain(id string, spec clusterv1.FailureDomainSpec)

	// PatchObject persists the cluster configuration and status.
	PatchObject() error
	// Close closes the current scope persisting the cluster configuration and status.
	Close() error
}
