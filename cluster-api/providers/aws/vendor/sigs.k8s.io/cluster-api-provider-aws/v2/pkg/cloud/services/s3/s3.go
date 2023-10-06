/*
Copyright 2021 The Kubernetes Authors.

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

package s3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/pkg/errors"

	iam "sigs.k8s.io/cluster-api-provider-aws/v2/iam/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/util/system"
)

// Service holds a collection of interfaces.
// The interfaces are broken down like this to group functions together.
// One alternative is to have a large list of functions from the ec2 client.
type Service struct {
	scope     scope.S3Scope
	S3Client  s3iface.S3API
	STSClient stsiface.STSAPI
}

// NewService returns a new service given the api clients.
func NewService(s3Scope scope.S3Scope) *Service {
	s3Client := scope.NewS3Client(s3Scope, s3Scope, s3Scope, s3Scope.InfraCluster())
	STSClient := scope.NewSTSClient(s3Scope, s3Scope, s3Scope, s3Scope.InfraCluster())

	return &Service{
		scope:     s3Scope,
		S3Client:  s3Client,
		STSClient: STSClient,
	}
}

func (s *Service) ReconcileBucket() error {
	if !s.bucketManagementEnabled() {
		return nil
	}

	bucketName := s.bucketName()

	if err := s.createBucketIfNotExist(bucketName); err != nil {
		return errors.Wrap(err, "ensuring bucket exists")
	}

	if err := s.ensureBucketPolicy(bucketName); err != nil {
		return errors.Wrap(err, "ensuring bucket policy")
	}

	return nil
}

func (s *Service) DeleteBucket() error {
	if !s.bucketManagementEnabled() {
		return nil
	}

	bucketName := s.bucketName()

	log := s.scope.WithValues("name", bucketName)

	log.Info("Deleting S3 Bucket")

	_, err := s.S3Client.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err == nil {
		return nil
	}

	aerr, ok := err.(awserr.Error)
	if !ok {
		return errors.Wrap(err, "deleting S3 bucket")
	}

	switch aerr.Code() {
	case s3.ErrCodeNoSuchBucket:
		log.Info("Bucket already removed")
	case "BucketNotEmpty":
		log.Info("Bucket not empty, skipping removal")
	default:
		return errors.Wrap(aerr, "deleting S3 bucket")
	}

	return nil
}

func (s *Service) Create(m *scope.MachineScope, data []byte) (string, error) {
	if !s.bucketManagementEnabled() {
		return "", errors.New("requested object creation but bucket management is not enabled")
	}

	if m == nil {
		return "", errors.New("machine scope can't be nil")
	}

	if len(data) == 0 {
		return "", errors.New("got empty data")
	}

	bucket := s.bucketName()
	key := s.bootstrapDataKey(m)

	s.scope.Info("Creating object", "bucket_name", bucket, "key", key)

	if _, err := s.S3Client.PutObject(&s3.PutObjectInput{
		Body:                 aws.ReadSeekCloser(bytes.NewReader(data)),
		Bucket:               aws.String(bucket),
		Key:                  aws.String(key),
		ServerSideEncryption: aws.String("aws:kms"),
	}); err != nil {
		return "", errors.Wrap(err, "putting object")
	}

	objectURL := &url.URL{
		Scheme: "s3",
		Host:   bucket,
		Path:   key,
	}

	return objectURL.String(), nil
}

func (s *Service) Delete(m *scope.MachineScope) error {
	if !s.bucketManagementEnabled() {
		return errors.New("requested object creation but bucket management is not enabled")
	}

	if m == nil {
		return errors.New("machine scope can't be nil")
	}

	bucket := s.bucketName()
	key := s.bootstrapDataKey(m)

	s.scope.Info("Deleting object", "bucket_name", bucket, "key", key)

	_, err := s.S3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err == nil {
		return nil
	}

	aerr, ok := err.(awserr.Error)
	if !ok {
		return errors.Wrap(err, "deleting S3 object")
	}

	switch aerr.Code() {
	case s3.ErrCodeNoSuchBucket:
	default:
		return errors.Wrap(aerr, "deleting S3 object")
	}

	return nil
}

func (s *Service) createBucketIfNotExist(bucketName string) error {
	input := &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	}

	_, err := s.S3Client.CreateBucket(input)
	if err == nil {
		s.scope.Info("Created bucket", "bucket_name", bucketName)

		return nil
	}

	aerr, ok := err.(awserr.Error)
	if !ok {
		return errors.Wrap(err, "creating S3 bucket")
	}

	switch aerr.Code() {
	// If bucket already exists, all good.
	//
	// TODO: This will fail if bucket is shared with other cluster.
	case s3.ErrCodeBucketAlreadyOwnedByYou:
		return nil
	default:
		return errors.Wrap(aerr, "creating S3 bucket")
	}
}

func (s *Service) ensureBucketPolicy(bucketName string) error {
	bucketPolicy, err := s.bucketPolicy(bucketName)
	if err != nil {
		return errors.Wrap(err, "generating Bucket policy")
	}

	input := &s3.PutBucketPolicyInput{
		Bucket: aws.String(bucketName),
		Policy: aws.String(bucketPolicy),
	}

	if _, err := s.S3Client.PutBucketPolicy(input); err != nil {
		return errors.Wrap(err, "creating S3 bucket policy")
	}

	s.scope.Trace("Updated bucket policy", "bucket_name", bucketName)

	return nil
}

func (s *Service) bucketPolicy(bucketName string) (string, error) {
	accountID, err := s.STSClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return "", errors.Wrap(err, "getting account ID")
	}

	bucket := s.scope.Bucket()
	partition := system.GetPartitionFromRegion(s.scope.Region())

	statements := []iam.StatementEntry{
		{
			Sid:    "control-plane",
			Effect: iam.EffectAllow,
			Principal: map[iam.PrincipalType]iam.PrincipalID{
				iam.PrincipalAWS: []string{fmt.Sprintf("arn:%s:iam::%s:role/%s", partition, *accountID.Account, bucket.ControlPlaneIAMInstanceProfile)},
			},
			Action:   []string{"s3:GetObject"},
			Resource: []string{fmt.Sprintf("arn:%s:s3:::%s/control-plane/*", partition, bucketName)},
		},
	}

	for _, iamInstanceProfile := range bucket.NodesIAMInstanceProfiles {
		statements = append(statements, iam.StatementEntry{
			Sid:    iamInstanceProfile,
			Effect: iam.EffectAllow,
			Principal: map[iam.PrincipalType]iam.PrincipalID{
				iam.PrincipalAWS: []string{fmt.Sprintf("arn:%s:iam::%s:role/%s", partition, *accountID.Account, iamInstanceProfile)},
			},
			Action:   []string{"s3:GetObject"},
			Resource: []string{fmt.Sprintf("arn:%s:s3:::%s/node/*", partition, bucketName)},
		})
	}

	policy := iam.PolicyDocument{
		Version:   "2012-10-17",
		Statement: statements,
	}

	policyRaw, err := json.Marshal(policy)
	if err != nil {
		return "", errors.Wrap(err, "building bucket policy")
	}

	return string(policyRaw), nil
}

func (s *Service) bucketManagementEnabled() bool {
	return s.scope.Bucket() != nil
}

func (s *Service) bucketName() string {
	return s.scope.Bucket().Name
}

func (s *Service) bootstrapDataKey(m *scope.MachineScope) string {
	// Use machine name as object key.
	return path.Join(m.Role(), m.Name())
}
