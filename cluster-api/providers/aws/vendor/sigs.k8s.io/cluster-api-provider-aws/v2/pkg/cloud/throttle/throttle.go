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
	"context"
	"regexp"
	"strings"

	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/smithy-go/middleware"

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

// Wait will wait on a request for AWS SDK V2.
func (o *OperationLimiter) Wait(ctx context.Context) error {
	return o.getLimiter().Wait(ctx)
}

// Match will match a request for AWS SDK V2.
func (o *OperationLimiter) Match(ctx context.Context) (bool, error) {
	if o.regexp == nil {
		var err error
		o.regexp, err = regexp.Compile("^" + o.Operation)
		if err != nil {
			return false, err
		}
	}
	opName := awsmiddleware.GetOperationName(ctx)
	return o.regexp.MatchString(opName), nil
}

// LimitRequest will limit a request for AWS SDK V2.
func (s ServiceLimiter) LimitRequest(ctx context.Context) {
	if ol, ok := s.matchRequest(ctx); ok {
		_ = ol.Wait(ctx)
	}
}

func (o *OperationLimiter) getLimiter() *rate.Limiter {
	if o.limiter == nil {
		o.limiter = rate.NewLimiter(o.RefillRate, o.Burst)
	}
	return o.limiter
}

// ReviewResponse will review the limits of a Request's response for AWS SDK V2.
func (s ServiceLimiter) ReviewResponse(ctx context.Context, errorCode string) {
	switch errorCode {
	case "Throttling", "RequestLimitExceeded":
		if ol, ok := s.matchRequest(ctx); ok {
			ol.limiter.ResetTokens()
		}
	}
}

// matchRequest is used for matching request for AWS SDK V2.
func (s ServiceLimiter) matchRequest(ctx context.Context) (*OperationLimiter, bool) {
	for _, ol := range s {
		match, err := ol.Match(ctx)
		if err != nil {
			return nil, false
		}
		if match {
			return ol, true
		}
	}
	return nil, false
}

// WithServiceLimiterMiddleware returns ServiceLimiter middleware stack for specified service name.
func WithServiceLimiterMiddleware(limiter *ServiceLimiter) func(stack *middleware.Stack) error {
	return func(stack *middleware.Stack) error {
		// Inserts service Limiter middleware after RequestContext initialization.
		return stack.Finalize.Insert(getServiceLimiterMiddleware(limiter), "capa/RequestMetricContextMiddleware", middleware.After)
	}
}

// getServiceLimiterMiddleware implements serviceLimiter middleware.
func getServiceLimiterMiddleware(limiter *ServiceLimiter) middleware.FinalizeMiddleware {
	return middleware.FinalizeMiddlewareFunc("capa/ServiceLimiterMiddleware", func(ctx context.Context, input middleware.FinalizeInput, handler middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
		limiter.LimitRequest(ctx)

		out, metadata, err := handler.HandleFinalize(ctx, input)
		smithyErr := awserrors.ParseSmithyError(err)

		if smithyErr != nil {
			limiter.ReviewResponse(ctx, smithyErr.ErrorCode())
			return out, metadata, err
		}

		return out, metadata, err
	})
}
