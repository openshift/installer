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

package secretsmanager

import (
	"fmt"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
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

	// we set the max secret size to well below the 10240 byte limit, because this is limit after base64 encoding,
	// but the aws sdk handles encoding for us, so we can't send a full 10240.
	maxSecretSizeBytes = 7000
)

var retryableErrors = []string{
	// Returned when the secret is scheduled for deletion
	secretsmanager.ErrCodeInvalidRequestException,
	// Returned during retries of deletes prior to recreation
	secretsmanager.ErrCodeResourceNotFoundException,
}

// Create stores data in AWS Secrets Manager for a given machine, chunking at 10kb per secret. The prefix of the secret
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
	// Split the data into chunks and create the secrets on demand.
	chunks := int32(0)
	var err error
	bytes.Split(data, false, maxSecretSizeBytes, func(chunk []byte) {
		name := fmt.Sprintf("%s-%d", prefix, chunks)
		retryFunc := func() (bool, error) { return s.retryableCreateSecret(name, chunk, tags) }
		// Default timeout is 5 mins, but if Secrets Manager has got to the state where the timeout is reached,
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
	_, err := s.SecretsManagerClient.CreateSecret(&secretsmanager.CreateSecretInput{
		Name:         aws.String(name),
		SecretBinary: chunk,
		Tags:         converters.MapToSecretsManagerTags(tags),
	})
	// If the secret already exists, delete it, return request to retry, as deletes are eventually consistent
	if awserrors.IsResourceExists(err) {
		return false, s.forceDeleteSecretEntry(name)
	}
	if err != nil {
		return false, err
	}
	return true, err
}

// forceDeleteSecretEntry deletes a single secret, ignoring if it is absent.
func (s *Service) forceDeleteSecretEntry(name string) error {
	_, err := s.SecretsManagerClient.DeleteSecret(&secretsmanager.DeleteSecretInput{
		SecretId:                   aws.String(name),
		ForceDeleteWithoutRecovery: aws.Bool(true),
	})
	if awserrors.IsNotFound(err) {
		return nil
	}
	return err
}

// Delete the secret belonging to a machine from AWS Secrets Manager.
func (s *Service) Delete(m *scope.MachineScope) error {
	var errs []error
	for i := range m.GetSecretCount() {
		if err := s.forceDeleteSecretEntry(fmt.Sprintf("%s-%d", m.GetSecretPrefix(), i)); err != nil {
			errs = append(errs, err)
		}
	}

	return kerrors.NewAggregate(errs)
}
