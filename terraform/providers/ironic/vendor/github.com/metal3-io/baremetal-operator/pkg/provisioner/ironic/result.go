package ironic

import (
	"time"

	"github.com/metal3-io/baremetal-operator/pkg/provisioner"
)

func retryAfterDelay(delay time.Duration) (provisioner.Result, error) {
	// TODO(zaneb): this is currently indistinguishable from the result of
	// operationContinuing() from the caller's perspective. Changes are
	// required to the Result structure to enable this to be distinguished.
	return provisioner.Result{
		Dirty:        true,
		RequeueAfter: delay,
	}, nil
}

func operationContinuing(delay time.Duration) (provisioner.Result, error) {
	return provisioner.Result{
		Dirty:        true,
		RequeueAfter: delay,
	}, nil
}

func operationComplete() (provisioner.Result, error) {
	return provisioner.Result{}, nil
}

func operationFailed(message string) (provisioner.Result, error) {
	return provisioner.Result{ErrorMessage: message}, nil
}

func transientError(err error) (provisioner.Result, error) {
	return provisioner.Result{}, err
}
