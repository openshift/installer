package controllers

import (
	"errors"

	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ControllerConfig is the configuration for cluster and machine controllers
type ControllerConfig struct {
	MaxConcurrentReconciles int
	RateLimiter             workqueue.TypedRateLimiter[reconcile.Request]
	SkipNameValidation      bool
}

// ControllerConfigOpts is a function that can be used to configure the controller config
type ControllerConfigOpts func(*ControllerConfig) error

// WithMaxConcurrentReconciles sets the maximum number of concurrent reconciles
func WithMaxConcurrentReconciles(max int) ControllerConfigOpts {
	return func(c *ControllerConfig) error {
		if max < 1 {
			return errors.New("max concurrent reconciles must be greater than 0")
		}
		c.MaxConcurrentReconciles = max
		return nil
	}
}

// WithRateLimiter sets the rate limiter for the controller
func WithRateLimiter(rateLimiter workqueue.TypedRateLimiter[reconcile.Request]) ControllerConfigOpts {
	return func(c *ControllerConfig) error {
		if rateLimiter == nil {
			return errors.New("rate limiter cannot be nil")
		}
		c.RateLimiter = rateLimiter
		return nil
	}
}

// WithSkipNameValidation sets whether to skip name validation in controller options
func WithSkipNameValidation(skip bool) ControllerConfigOpts {
	return func(c *ControllerConfig) error {
		c.SkipNameValidation = skip
		return nil
	}
}
