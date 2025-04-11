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

package ssm

import (
	"fmt"
	"path"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/uuid"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/converters"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/wait"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/internal/bytes"
)

const (
	entryPrefix = "aws.cluster.x-k8s.io"

	// max byte size for ssm is 4KB else we cross into the advanced-parameter tier.
	maxSecretSizeBytes = 4000
)

var (
	prefixRe        = regexp.MustCompile(`(?i)^[\/]?(aws|ssm)[.]?`)
	retryableErrors = []string{
		ssm.ErrCodeParameterLimitExceeded,
	}
)

// Create stores data in AWS SSM for a given machine, chunking at 4kb per secret. The prefix of the secret
// ARN and the number of chunks are returned.
func (s *Service) Create(m *scope.MachineScope, data []byte) (string, int32, error) {
	// Build the tags to apply to the secret.
	additionalTags := m.AdditionalTags()
	additionalTags[infrav1.ClusterAWSCloudProviderTagKey(s.scope.Name())] = string(infrav1.ResourceLifecycleOwned)
	tags := infrav1.Build(infrav1.BuildParams{
		ClusterName: s.scope.Name(),
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		Name:        aws.String(m.Name()),
		Role:        aws.String(m.Role()),
		Additional:  additionalTags,
	})

	// Build the prefix.
	prefix := m.GetSecretPrefix()
	if prefix == "" {
		prefix = path.Join(entryPrefix, string(uuid.NewUUID()))
	}
	// SSM Validation does not allow (/)aws|ssm in the beginning of the string
	prefix = prefixRe.ReplaceAllString(prefix, "")
	// Because the secret name has a slash in it, whole name must validate as a full path
	if prefix[0] != byte('/') {
		prefix = "/" + prefix
	}

	// Split the data into chunks and create the secrets on demand.
	chunks := int32(0)
	var err error
	bytes.Split(data, true, maxSecretSizeBytes, func(chunk []byte) {
		name := fmt.Sprintf("%s/%d", prefix, chunks)
		retryFunc := func() (bool, error) { return s.retryableCreateSecret(name, chunk, tags) }
		// Default timeout is 5 mins, but if SSM has got to the state where the timeout is reached,
		// makes sense to slow down machine creation until AWS weather improves.
		if err = wait.WaitForWithRetryable(wait.NewBackoff(), retryFunc, retryableErrors...); err != nil {
			return
		}
		chunks++
	})

	return prefix, chunks, err
}

// retryableCreateSecret is a function to be passed into a waiter. In a separate function for ease of reading.
func (s *Service) retryableCreateSecret(name string, chunk []byte, tags infrav1.Tags) (bool, error) {
	_, err := s.SSMClient.PutParameter(&ssm.PutParameterInput{
		Name:  aws.String(name),
		Value: aws.String(string(chunk)),
		Tags:  converters.MapToSSMTags(tags),
		Type:  aws.String("SecureString"),
	})
	if err != nil {
		return false, err
	}
	return true, err
}

// forceDeleteSecretEntry deletes a single secret, ignoring if it is absent.
func (s *Service) forceDeleteSecretEntry(name string) error {
	_, err := s.SSMClient.DeleteParameter(&ssm.DeleteParameterInput{
		Name: aws.String(name),
	})
	if awserrors.IsNotFound(err) {
		return nil
	}
	return err
}

// Delete the secret belonging to a machine from AWS SSM.
func (s *Service) Delete(m *scope.MachineScope) error {
	var errs []error
	for i := range m.GetSecretCount() {
		if err := s.forceDeleteSecretEntry(fmt.Sprintf("%s/%d", m.GetSecretPrefix(), i)); err != nil {
			errs = append(errs, err)
		}
	}

	return kerrors.NewAggregate(errs)
}
