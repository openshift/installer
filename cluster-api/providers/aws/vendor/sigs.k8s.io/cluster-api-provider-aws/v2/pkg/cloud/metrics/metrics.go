/*
Copyright 2025 The Kubernetes Authors.

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

// Package metricsv2 provides a way to capture request metrics.
package metricsv2

import (
	"context"
	"fmt"
	"strconv"
	"time"

	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/metrics"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/record"
	"sigs.k8s.io/cluster-api-provider-aws/v2/version"
)

const (
	metricAWSSubsystem       = "aws"
	metricRequestCountKey    = "api_requests_total_v2"
	metricRequestDurationKey = "api_request_duration_seconds_v2"
	metricAPICallRetries     = "api_call_retries_v2"
	metricServiceLabel       = "service"
	metricRegionLabel        = "region"
	metricOperationLabel     = "operation"
	metricControllerLabel    = "controller"
	metricStatusCodeLabel    = "status_code"
	metricErrorCodeLabel     = "error_code"
)

var (
	awsRequestCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Subsystem: metricAWSSubsystem,
		Name:      metricRequestCountKey,
		Help:      "Total number of AWS requests",
	}, []string{metricControllerLabel, metricServiceLabel, metricRegionLabel, metricOperationLabel, metricStatusCodeLabel, metricErrorCodeLabel})

	awsRequestDurationSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem: metricAWSSubsystem,
		Name:      metricRequestDurationKey,
		Help:      "Latency of HTTP requests to AWS",
	}, []string{metricControllerLabel, metricServiceLabel, metricRegionLabel, metricOperationLabel})

	awsCallRetries = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem: metricAWSSubsystem,
		Name:      metricAPICallRetries,
		Help:      "Number of retries made against an AWS API",
		Buckets:   []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}, []string{metricControllerLabel, metricServiceLabel, metricRegionLabel, metricOperationLabel})
)

func init() {
	metrics.Registry.MustRegister(awsRequestCount)
	metrics.Registry.MustRegister(awsRequestDurationSeconds)
	metrics.Registry.MustRegister(awsCallRetries)
}

type requestContextKey struct{}

// RequestData holds information related to request metrics.
type RequestData struct {
	RequestStartTime time.Time
	RequestEndTime   time.Time
	StatusCode       int
	ErrorCode        string
	Service          string
	OperationName    string
	Region           string
	Controller       string
	Target           runtime.Object
	Attempts         int
}

// WithMiddlewares adds instrumentation middleware stacks to AWS GO SDK V2 service clients.
func WithMiddlewares(controller string, target runtime.Object) func(stack *middleware.Stack) error {
	return func(stack *middleware.Stack) error {
		if err := stack.Initialize.Add(getMetricCollectionMiddleware(controller, target), middleware.Before); err != nil {
			return fmt.Errorf("failed to add metric collection middleware: %w", err)
		}
		if err := stack.Finalize.Add(getRequestMetricContextMiddleware(), middleware.Before); err != nil {
			return fmt.Errorf("failed to add request metric context middleware: %w", err)
		}
		// Check if the Default Retry Middleware exists.
		if _, ok := stack.Finalize.Get((*retry.Attempt)(nil).ID()); ok {
			if err := stack.Finalize.Insert(getAttemptContextMiddleware(), (*retry.Attempt)(nil).ID(), middleware.After); err != nil {
				return fmt.Errorf("failed to add attempt context middleware: %w", err)
			}
		}
		return stack.Finalize.Add(getRecordAWSPermissionsIssueMiddleware(target), middleware.After)
	}
}

// WithCAPAUserAgentMiddleware returns User Agent middleware stack for AWS GO SDK V2 sessions.
func WithCAPAUserAgentMiddleware() func(*middleware.Stack) error {
	return awsmiddleware.AddUserAgentKeyValue("aws.cluster.x-k8s.io", version.Get().String())
}

// WithRequestMetricContextMiddleware returns Request Metric middleware stack for AWS GO SDK V2 sessions.
func WithRequestMetricContextMiddleware() func(*middleware.Stack) error {
	return func(stack *middleware.Stack) error {
		return stack.Finalize.Add(getRequestMetricContextMiddleware(), middleware.Before)
	}
}

func getMetricCollectionMiddleware(controller string, target runtime.Object) middleware.InitializeMiddleware {
	return middleware.InitializeMiddlewareFunc("capa/MetricCollectionMiddleware", func(ctx context.Context, input middleware.InitializeInput, handler middleware.InitializeHandler) (middleware.InitializeOutput, middleware.Metadata, error) {
		ctx = initRequestContext(ctx, controller, target)
		request := getContext(ctx)

		request.RequestStartTime = time.Now().UTC()
		out, metadata, err := handler.HandleInitialize(ctx, input)

		if responseAt, ok := awsmiddleware.GetResponseAt(metadata); ok {
			request.RequestEndTime = responseAt
		} else {
			request.RequestEndTime = time.Now().UTC()
		}
		request.CaptureRequestMetrics()
		return out, metadata, err
	})
}

func getRequestMetricContextMiddleware() middleware.FinalizeMiddleware {
	return middleware.FinalizeMiddlewareFunc("capa/RequestMetricContextMiddleware", func(ctx context.Context, input middleware.FinalizeInput, handler middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
		request := getContext(ctx)

		if request != nil {
			request.Service = awsmiddleware.GetServiceID(ctx)
			request.OperationName = awsmiddleware.GetOperationName(ctx)
			request.Region = awsmiddleware.GetRegion(ctx)
		}
		return handler.HandleFinalize(ctx, input)
	})
}

// getAttemptContextMiddleware will capture StatusCode and ErrorCode from API call attempt.
// This will result in the StatusCode and ErrorCode captured for last attempt when publishing to metrics.
func getAttemptContextMiddleware() middleware.FinalizeMiddleware {
	return middleware.FinalizeMiddlewareFunc("capa/AttemptMetricContextMiddleware", func(ctx context.Context, input middleware.FinalizeInput, handler middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
		request := getContext(ctx)
		if request != nil {
			request.Attempts++
		}

		out, metadata, err := handler.HandleFinalize(ctx, input)

		if request != nil {
			if rawResp := awsmiddleware.GetRawResponse(metadata); rawResp != nil {
				if httpResp, ok := rawResp.(*smithyhttp.Response); ok {
					request.StatusCode = httpResp.StatusCode
				}
			} else {
				request.StatusCode = -1
			}

			if err != nil {
				smithyErr := awserrors.ParseSmithyError(err)
				request.ErrorCode = smithyErr.ErrorCode()
				request.StatusCode = smithyErr.StatusCode()
			}
		}

		return out, metadata, err
	})
}

func getRecordAWSPermissionsIssueMiddleware(target runtime.Object) middleware.FinalizeMiddleware {
	return middleware.FinalizeMiddlewareFunc("capa/RecordAWSPermissionsIssueMiddleware", func(ctx context.Context, input middleware.FinalizeInput, handler middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
		output, metadata, err := handler.HandleFinalize(ctx, input)
		if err != nil {
			request := getContext(ctx)
			if request != nil {
				var errMessage string
				smithyErr := awserrors.ParseSmithyError(err)
				request.ErrorCode = smithyErr.ErrorCode()
				switch request.ErrorCode {
				case "AccessDenied", "AuthFailure", "UnauthorizedOperation", "NoCredentialProviders",
					"ExpiredToken", "InvalidClientTokenId", "SignatureDoesNotMatch", "ValidationError":
					record.Warnf(target, request.ErrorCode,
						"Operation %s failed with a credentials or permission issue: %s",
						request.OperationName,
						errMessage,
					)
				}
			}
		}
		return output, metadata, err
	})
}

func initRequestContext(ctx context.Context, controller string, target runtime.Object) context.Context {
	if middleware.GetStackValue(ctx, requestContextKey{}) == nil {
		ctx = middleware.WithStackValue(ctx, requestContextKey{}, &RequestData{
			Controller: controller,
			Target:     target,
			Attempts:   0,
		})
	}
	return ctx
}

func getContext(ctx context.Context) *RequestData {
	rctx := middleware.GetStackValue(ctx, requestContextKey{})
	if rctx == nil {
		return nil
	}
	return rctx.(*RequestData)
}

// CaptureRequestMetrics will monitor and capture request metrics.
func (r *RequestData) CaptureRequestMetrics() {
	if !r.IsIncomplete() {
		return
	}

	requestDuration := r.RequestEndTime.Sub(r.RequestStartTime)
	retryCount := max(r.Attempts-1, 0)
	statusCode := strconv.Itoa(r.StatusCode)
	errorCode := r.ErrorCode

	if errorCode == "" && r.StatusCode >= 400 {
		errorCode = fmt.Sprintf("HTTP%d", r.StatusCode)
	}

	awsRequestCount.WithLabelValues(
		r.Controller,
		r.Service,
		r.Region,
		r.OperationName,
		statusCode,
		errorCode,
	).Inc()

	awsRequestDurationSeconds.WithLabelValues(
		r.Controller,
		r.Service,
		r.Region,
		r.OperationName,
	).Observe(requestDuration.Seconds())

	awsCallRetries.WithLabelValues(
		r.Controller,
		r.Service,
		r.Region,
		r.OperationName,
	).Observe(float64(retryCount))
}

// IsIncomplete will return true if the RequestData was incomplete.
func (r *RequestData) IsIncomplete() bool {
	return r.Service == "" || r.Region == "" || r.OperationName == "" || r.Controller == ""
}
