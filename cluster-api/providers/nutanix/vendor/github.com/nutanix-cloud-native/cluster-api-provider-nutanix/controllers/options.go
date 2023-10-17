package controllers

import "errors"

// ControllerConfig is the configuration for cluster and machine controllers
type ControllerConfig struct {
	MaxConcurrentReconciles int
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
