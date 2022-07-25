package ovirtclient

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	ovirtclientlog "github.com/ovirt/go-ovirt-client-log/v3"
)

// retry is a function that will automatically retry calling the function specified in the what parameter until the
// timeouts specified in the howLong parameter are reached. It attempts to identify permanent failures and abort if
// they are encountered.
//
// - action is the action that is being performed in the "ing" form, for example "creating disk".
// - what is the function that should be called repeatedly.
// - logger is an optional logger that can be passed to log retry actions.
// - howLong is the retry configuration that should be used.
func retry(
	action string,
	logger ovirtclientlog.Logger,
	howLong []RetryStrategy,
	what func() error,
) error {
	retries := make([]RetryInstance, len(howLong))
	for i, factory := range howLong {
		retries[i] = factory.Get()
	}

	if logger == nil {
		logger = &noopLogger{}
	}
	logger.Infof("%s%s...", strings.ToUpper(action[:1]), action[1:])
	for {
		err := what()
		if err == nil {
			logger.Infof("Completed %s.", action)
			return nil
		}
		for _, r := range retries {
			if err := r.Continue(err, action); err != nil {
				logger.Infof("Giving up %s (%v)", action, err)
				return err
			}
		}

		if !recoverFailure(action, retries, err, logger) {
			logRetry(action, logger, err)
		}
		// Here we create a select statement with a dynamic number of cases. We use this because a) select{} only
		// supports fixed cases and b) the channel types are different. Context returns a <-chan struct{}, while
		// time.After() returns <-chan time.Time. Go doesn't support type assertions, so we have to result to
		// the reflection library to do this.
		var chans []reflect.SelectCase
		for _, r := range retries {
			c := r.Wait(err)
			if c != nil {
				chans = append(
					chans, reflect.SelectCase{
						Dir:  reflect.SelectRecv,
						Chan: reflect.ValueOf(c),
						Send: reflect.Value{},
					},
				)
			}
		}
		if len(chans) == 0 {
			logger.Errorf(
				"No retry strategies with waiting function specified for %s.",
				action,
			)
			return newError(EBug, "no retry strategies with waiting function specified for %s", action)
		}
		chosen, _, _ := reflect.Select(chans)
		if err := retries[chosen].OnWaitExpired(err, action); err != nil {
			logger.Infof("Giving up %s (%v)", action, err)
			return err
		}
	}
}

func recoverFailure(action string, retries []RetryInstance, err error, logger ovirtclientlog.Logger) bool {
	var e EngineError
	if !errors.As(err, &e) {
		return false
	}
	if !e.CanRecover() {
		return false
	}
	logger.Debugf("Failed %s, attempting automatic recovery... (%s)", err.Error())
	for _, r := range retries {
		recoveredErr := r.Recover(err)
		switch {
		case recoveredErr == nil:
			logger.Debugf("Automatic recovery successful, retrying %s...", action)
			return true
		// We purposefully ignore the wrapped errors here and only compare 1:1 to make sure we only use this case if
		// the errors don't match. If the errors match the retry strategy couldn't do anything with the error.
		// Do not change this to errors.Is!
		case recoveredErr != err: //nolint:errorlint
			logger.Errorf("Error encountered during automatic recovery of %s (%v).", action, recoveredErr)
			return false
		}
	}
	logger.Debugf("No appropriate recovery mechanism found for failure on %s, retrying without recovery (%v).")
	return false
}

func logRetry(action string, logger ovirtclientlog.Logger, err error) {
	var e EngineError
	isPending := false
	isConflict := false
	if errors.As(err, &e) {
		isPending = e.HasCode(EPending)
		isConflict = e.HasCode(EConflict) || e.HasCode(EDiskLocked) || e.HasCode(EVMLocked)
	}
	if isPending || isConflict {
		logger.Debugf("Still %s, retrying... (%s)", action, err.Error())
	} else {
		logger.Debugf("Failed %s, retrying... (%s)", action, err.Error())
	}
}

type noopLogger struct{}

func (n noopLogger) WithContext(_ context.Context) ovirtclientlog.Logger {
	return n
}

func (n noopLogger) Debugf(_ string, _ ...interface{}) {}

func (n noopLogger) Infof(_ string, _ ...interface{}) {}

func (n noopLogger) Warningf(_ string, _ ...interface{}) {}

func (n noopLogger) Errorf(_ string, _ ...interface{}) {}

// RetryStrategy is a function that creates a new copy of a RetryInstance. It is important because each
// RetryInstance may have an internal state, so reusing a RetryInstance won't work. RetryStrategy copies can be
// safely passed around between functions and reused multiple times.
type RetryStrategy interface {
	// Get returns an actual copy of the retry strategy. This can be used to initialize individual timers for
	// separate API calls within a larger call structure.
	Get() RetryInstance

	// CanClassifyErrors indicates if the strategy can determine if an error is retryable. At least one strategy with
	// this capability needs to be passed.
	CanClassifyErrors() bool
	// CanWait indicates if the retry strategy can wait in a loop. At least one strategy with this capability
	// needs to be passed.
	CanWait() bool
	// CanTimeout indicates that the retry strategy can properly abort a loop. At least one retry strategy with
	// this capability needs to be passed.
	CanTimeout() bool
	// CanRecover indicates that the retry strategy can recoverFailure from a failure. In this case the Recover method will be
	// called on errors.
	CanRecover() bool
}

type retryStrategyContainer struct {
	factory           func() RetryInstance
	canClassifyErrors bool
	canWait           bool
	canTimeout        bool
	canRecover        bool
}

func (r retryStrategyContainer) CanRecover() bool {
	return r.canRecover
}

func (r retryStrategyContainer) Get() RetryInstance {
	return r.factory()
}

func (r retryStrategyContainer) CanClassifyErrors() bool {
	return r.canClassifyErrors
}

func (r retryStrategyContainer) CanWait() bool {
	return r.canWait
}

func (r retryStrategyContainer) CanTimeout() bool {
	return r.canTimeout
}

// RetryInstance is an instance created by the RetryStrategy for a single use. It may have internal state
// and should not be reused.
type RetryInstance interface {
	// Continue returns an error if no more tries should be attempted. The error will be returned directly from the
	// retry function. The passed action parameters can be used to create a meaningful error message.
	Continue(err error, action string) error
	// Recover gives the strategy a chance to recoverFailure from a failure. This function is called if an execution errored
	// before the next cycle happens. The recovery method should return nil if the recovery was successful, and an
	// error otherwise. It may return the same error it received to indicate that it could not do anything with the
	// error.
	Recover(err error) error
	// Wait returns a channel that is closed when the wait time expires. The channel can have any content, so it is
	// provided as an interface{}. This function may return nil if it doesn't provide a wait time.
	Wait(err error) interface{}
	// OnWaitExpired is a hook that gives the strategy the option to return an error if its wait has expired. It will
	// only be called if it is the first to reach its wait. If no error is returned the loop is continued. The passed
	// action names can be incorporated into an error message.
	OnWaitExpired(err error, action string) error
}

// ContextStrategy provides a timeout based on a context in the ctx parameter. If the context is canceled the
// retry loop is aborted.
func ContextStrategy(ctx context.Context) RetryStrategy {
	return &retryStrategyContainer{
		func() RetryInstance {
			return &contextStrategy{
				ctx: ctx,
			}
		},
		false,
		false,
		true,
		false,
	}
}

type contextStrategy struct {
	ctx context.Context
}

func (c *contextStrategy) Recover(err error) error { return err }

func (c *contextStrategy) Name() string {
	return "context strategy"
}

func (c *contextStrategy) Continue(_ error, _ string) error {
	return nil
}

func (c *contextStrategy) Wait(_ error) interface{} {
	return c.ctx.Done()
}

func (c *contextStrategy) OnWaitExpired(err error, action string) error {
	return wrap(
		err,
		ETimeout,
		"timeout while %s",
		action,
	)
}

// ExponentialBackoff is a retry strategy that increases the wait time after each call by the specified factor.
func ExponentialBackoff(factor uint8) RetryStrategy {
	return &retryStrategyContainer{
		func() RetryInstance {
			waitTime := time.Second
			return &exponentialBackoff{
				waitTime: waitTime,
				factor:   factor,
			}
		},
		false,
		true,
		false,
		false,
	}
}

type exponentialBackoff struct {
	waitTime time.Duration
	factor   uint8
}

func (e *exponentialBackoff) Recover(err error) error { return err }

func (e *exponentialBackoff) Name() string {
	return fmt.Sprintf("exponential backoff strategy of %d seconds", e.waitTime/time.Second)
}

func (e *exponentialBackoff) Wait(_ error) interface{} {
	waitTime := e.waitTime
	e.waitTime *= time.Duration(e.factor)
	return time.After(waitTime)
}

func (e *exponentialBackoff) OnWaitExpired(_ error, _ string) error {
	return nil
}

func (e *exponentialBackoff) Continue(_ error, _ string) error {
	return nil
}

// AutoRetry retries an action only if it doesn't return a non-retryable error.
func AutoRetry() RetryStrategy {
	return &retryStrategyContainer{
		func() RetryInstance {
			return &autoRetryStrategy{}
		},
		true,
		false,
		false,
		false,
	}
}

type autoRetryStrategy struct{}

func (a *autoRetryStrategy) Recover(err error) error { return err }

func (a *autoRetryStrategy) Name() string {
	return "abort non-retryable errors strategy"
}

func (a *autoRetryStrategy) Continue(err error, action string) error {
	var engineErr EngineError
	if errors.As(err, &engineErr) {
		if !engineErr.CanAutoRetry() {
			return wrap(
				err,
				EUnidentified,
				"non-retryable error encountered while %s, giving up",
				action,
			)
		}
		return nil
	}
	identifiedError := realIdentify(err)
	if identifiedError == nil {
		return wrap(
			err,
			EUnidentified,
			"non-retryable error encountered while %s, giving up",
			action,
		)
	}
	if !identifiedError.CanAutoRetry() {
		return wrap(
			err,
			EUnidentified,
			"non-retryable error encountered while %s, giving up",
			action,
		)
	}
	return nil
}

func (a *autoRetryStrategy) Wait(_ error) interface{} {
	return nil
}

func (a *autoRetryStrategy) OnWaitExpired(_ error, _ string) error {
	return nil
}

// MaxTries is a strategy that will timeout individual API calls based on a maximum number of retries. The total number
// of API calls can be higher in case of a complex functions that involve multiple API calls.
func MaxTries(tries uint16) RetryStrategy {
	return &retryStrategyContainer{
		func() RetryInstance {
			return &maxTriesStrategy{
				maxTries: tries,
				tries:    0,
			}
		},
		false,
		false,
		true,
		false,
	}
}

type maxTriesStrategy struct {
	maxTries uint16
	tries    uint16
}

func (m *maxTriesStrategy) Recover(err error) error { return err }

func (m *maxTriesStrategy) Name() string {
	return fmt.Sprintf("maximum of %d retries strategy", m.maxTries)
}

func (m *maxTriesStrategy) Continue(err error, action string) error {
	m.tries++
	if m.tries > m.maxTries {
		return wrap(
			err,
			ETimeout,
			"maximum retries reached while %s, giving up",
			action,
		)
	}
	return nil
}

func (m *maxTriesStrategy) Wait(_ error) interface{} {
	return nil
}

func (m *maxTriesStrategy) OnWaitExpired(_ error, _ string) error {
	return nil
}

// Timeout is a strategy that will time out complex calls based on a timeout from the time the strategy factory was
// created. This is contrast to CallTimeout, which will evaluate timeouts for each individual API call.
func Timeout(timeout time.Duration) RetryStrategy {
	startTime := time.Now()
	return &retryStrategyContainer{
		func() RetryInstance {
			return &timeoutStrategy{
				duration:  timeout,
				startTime: startTime,
			}
		},
		false,
		false,
		true,
		false,
	}
}

// CallTimeout is a strategy that will timeout individual API call retries.
func CallTimeout(timeout time.Duration) RetryStrategy {
	return &retryStrategyContainer{
		func() RetryInstance {
			startTime := time.Now()
			return &timeoutStrategy{
				duration:  timeout,
				startTime: startTime,
			}
		},
		false,
		false,
		true,
		false,
	}
}

type timeoutStrategy struct {
	duration  time.Duration
	startTime time.Time
}

func (t *timeoutStrategy) Recover(err error) error { return err }

func (t *timeoutStrategy) Continue(err error, action string) error {
	if elapsedTime := time.Since(t.startTime); elapsedTime > t.duration {
		return wrap(
			err,
			ETimeout,
			"timeout of %d seconds while %s, giving up",
			t.duration/time.Second,
			action,
		)
	}
	return nil
}

func (t *timeoutStrategy) Wait(_ error) interface{} {
	return nil
}

func (t *timeoutStrategy) OnWaitExpired(_ error, _ string) error {
	return nil
}

// ReconnectStrategy triggers the client to reconnect if an EInvalidGrant error is encountered.
func ReconnectStrategy(client Client) RetryStrategy {
	return &retryStrategyContainer{
		func() RetryInstance {
			return &reconnectStrategy{
				client: client,
			}
		},
		false,
		false,
		false,
		true,
	}
}

type reconnectStrategy struct {
	client Client
}

func (r reconnectStrategy) Continue(_ error, _ string) error {
	return nil
}

func (r reconnectStrategy) Recover(err error) error {
	if HasErrorCode(err, EInvalidGrant) {
		return r.client.Reconnect()
	}
	return err
}

func (r reconnectStrategy) Wait(_ error) interface{} {
	return nil
}

func (r reconnectStrategy) OnWaitExpired(_ error, _ string) error {
	return nil
}

func defaultRetries(retries []RetryStrategy, timeout []RetryStrategy) []RetryStrategy {
	foundWait := false
	foundTimeout := false
	foundClassifier := false
	for _, r := range retries {
		if r.CanWait() {
			foundWait = true
		}
		if r.CanTimeout() {
			foundTimeout = true
		}
		if r.CanClassifyErrors() {
			foundClassifier = true
		}
	}
	if !foundWait {
		retries = append(retries, ExponentialBackoff(2))
	}
	if !foundTimeout {
		retries = append(retries, timeout...)
	}
	if !foundClassifier {
		retries = append(retries, AutoRetry())
	}
	return retries
}

// defaultReadTimeouts returns a list of retry strategies suitable for read calls. There are view retries and
// individual calls with retries shouldn't last longer than a minute, otherwise something went wrong.
func defaultReadTimeouts(client Client) []RetryStrategy {
	if ctx := client.GetContext(); ctx != nil {
		return []RetryStrategy{
			MaxTries(10),
			ContextStrategy(ctx),
			ReconnectStrategy(client),
		}
	}
	return []RetryStrategy{
		MaxTries(3),
		CallTimeout(time.Minute),
		Timeout(5 * time.Minute),
		ReconnectStrategy(client),
	}
}

// defaultWriteTimeouts has slightly higher tolerances for write API calls, as they may need longer waiting
// times.
func defaultWriteTimeouts(client Client) []RetryStrategy {
	if ctx := client.GetContext(); ctx != nil {
		return []RetryStrategy{
			MaxTries(10),
			ContextStrategy(ctx),
			ReconnectStrategy(client),
		}
	}
	return []RetryStrategy{
		MaxTries(10),
		CallTimeout(5 * time.Minute),
		Timeout(10 * time.Minute),
		ReconnectStrategy(client),
	}
}

// defaultLongTimeouts contains a strategy to wait for calls that typically take longer, for example waiting for a
// disk to become ready.
func defaultLongTimeouts(client Client) []RetryStrategy {
	if ctx := client.GetContext(); ctx != nil {
		return []RetryStrategy{
			MaxTries(10),
			ContextStrategy(ctx),
			ReconnectStrategy(client),
		}
	}
	return []RetryStrategy{
		MaxTries(30),
		CallTimeout(15 * time.Minute),
		Timeout(30 * time.Minute),
		ReconnectStrategy(client),
	}
}
