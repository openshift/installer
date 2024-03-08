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

// Package throttle provides a way to limit the number of requests to AWS services.
package throttle

import (
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws/request"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/internal/rate"
)

// ServiceLimiters defines a mapping of service limiters.
type ServiceLimiters map[string]*ServiceLimiter

// ServiceLimiter defines a buffer of operation limiters.
type ServiceLimiter []*OperationLimiter

// NewMultiOperationMatch will create a multi operation matching.
func NewMultiOperationMatch(strs ...string) string {
	return "^" + strings.Join(strs, "|^")
}

// OperationLimiter defines the specs of an operation limiter.
type OperationLimiter struct {
	Operation  string
	RefillRate rate.Limit
	Burst      int
	regexp     *regexp.Regexp
	limiter    *rate.Limiter
}

// Wait will wait on a request.
func (o *OperationLimiter) Wait(r *request.Request) error {
	return o.getLimiter().Wait(r.Context())
}

// Match will match a request.
func (o *OperationLimiter) Match(r *request.Request) (bool, error) {
	if o.regexp == nil {
		var err error
		o.regexp, err = regexp.Compile("^" + o.Operation)
		if err != nil {
			return false, err
		}
	}
	return o.regexp.MatchString(r.Operation.Name), nil
}

// LimitRequest will limit a request.
func (s ServiceLimiter) LimitRequest(r *request.Request) {
	if ol, ok := s.matchRequest(r); ok {
		_ = ol.Wait(r)
	}
}

func (o *OperationLimiter) getLimiter() *rate.Limiter {
	if o.limiter == nil {
		o.limiter = rate.NewLimiter(o.RefillRate, o.Burst)
	}
	return o.limiter
}

// ReviewResponse will review the limits of a Request's response.
func (s ServiceLimiter) ReviewResponse(r *request.Request) {
	if r.Error != nil {
		if errorCode, ok := awserrors.Code(r.Error); ok {
			switch errorCode {
			case "Throttling", "RequestLimitExceeded":
				if ol, ok := s.matchRequest(r); ok {
					ol.limiter.ResetTokens()
				}
			}
		}
	}
}

func (s ServiceLimiter) matchRequest(r *request.Request) (*OperationLimiter, bool) {
	for _, ol := range s {
		match, err := ol.Match(r)
		if err != nil {
			return nil, false
		}
		if match {
			return ol, true
		}
	}
	return nil, false
}
