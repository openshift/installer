package vsphere

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/event"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
)

const (
	eventTypeVMPoweredOffEvent = "VmPoweredOffEvent"
)

// virtualMachineCustomizationWaiter is an object that waits for customization
// of a VirtualMachine to complete, by watching for success or failure events.
//
// The waiter should be created with newWaiter **before** the start of the
// customization task to be 100% certain that completion events are not missed.
type virtualMachineCustomizationWaiter struct {
	// This channel will be closed upon completion, and should be blocked on.
	done chan struct{}

	// Any error received from the waiter - be it the customization failure
	// itself, timeouts waiting for the completion events, or other API-related
	// errors. This will always be nil until done is closed.
	err error
}

// Done returns the done channel. This channel will be closed upon completion,
// and should be blocked on.
func (w *virtualMachineCustomizationWaiter) Done() chan struct{} {
	return w.done
}

// Err returns any error received from the waiter. This will always be nil
// until the channel returned by Done is closed.
func (w *virtualMachineCustomizationWaiter) Err() error {
	return w.err
}

// newVirtualMachineCustomizationWaiter returns a new
// virtualMachineCustomizationWaiter to use to wait for customization on.
//
// This should be called **before** the start of the customization task to be
// 100% certain that completion events are not missed.
//
// The timeout value is in minutes - a value of less than 1 disables the waiter
// and returns immediately without error.
func newVirtualMachineCustomizationWaiter(client *govmomi.Client, vm *object.VirtualMachine, timeout int) *virtualMachineCustomizationWaiter {
	w := &virtualMachineCustomizationWaiter{
		done: make(chan struct{}),
	}
	go func() {
		w.err = w.wait(client, vm, timeout)
		close(w.done)
	}()
	return w
}

// wait waits for the customization of a supplied VirtualMachine to complete,
// either due to success or error. It does this by watching specifically for
// CustomizationSucceeded and CustomizationFailed events. If the customization
// failed due to some sort of error, the full formatted message is returned as
// an error.
func (w *virtualMachineCustomizationWaiter) wait(client *govmomi.Client, vm *object.VirtualMachine, timeout int) error {
	// A timeout of less than 1 minute (zero or negative value) skips the waiter,
	// so we return immediately.
	if timeout < 1 {
		return nil
	}

	// Our listener loop callback.
	cbErr := make(chan error, 1)
	cb := func(obj types.ManagedObjectReference, page []types.BaseEvent) error {
		for _, be := range page {
			switch e := be.(type) {
			case types.BaseCustomizationFailed:
				cbErr <- errors.New(e.GetCustomizationFailed().GetEvent().FullFormattedMessage)
			case *types.CustomizationSucceeded:
				close(cbErr)
			}
		}
		return nil
	}

	mgr := event.NewManager(client.Client)
	mgrErr := make(chan error, 1)
	// Make a proper background context so that we can gracefully cancel the
	// subscriber when we are done with it. This eventually gets passed down to
	// the property collector SOAP calls.
	pctx, pcancel := context.WithCancel(context.Background())
	defer pcancel()
	go func() {
		mgrErr <- mgr.Events(pctx, []types.ManagedObjectReference{vm.Reference()}, 10, true, false, cb)
	}()

	// Wait for any error condition (including nil from the closure of the
	// callback error channel on success). We also use a different context so
	// that we can give a better error message on timeout without interfering
	// with the subscriber's context.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Minute)
	defer cancel()
	var err error
	select {
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			err = fmt.Errorf("timeout waiting for customization to complete")
		}
	case err = <-mgrErr:
	case err = <-cbErr:
	}
	return err
}

// selectEventsForReference allows you to query events for a specific
// ManagedObjectReference.
//
// Event types can be supplied to this function via the eventTypes parameter.
// This is highly recommended when you expect the list of events to be large,
// as there is no limit on returned events.
func selectEventsForReference(client *govmomi.Client, ref types.ManagedObjectReference, eventTypes []string) ([]types.BaseEvent, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()
	filter := types.EventFilterSpec{
		Entity: &types.EventFilterSpecByEntity{
			Entity:    ref,
			Recursion: types.EventFilterSpecRecursionOptionAll,
		},
		EventTypeId: eventTypes,
	}
	mgr := event.NewManager(client.Client)
	return mgr.QueryEvents(ctx, filter)
}
