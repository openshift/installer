package exec

import (
	"fmt"
	"strings"
	"sync"

	"github.com/hashicorp/terraform/terraform"
)

type countHookAction byte

const (
	countHookActionAdd countHookAction = iota
	countHookActionChange
	countHookActionRemove
)

// CountHook is a hook that counts the number of resources
// added, removed, changed during the course of an apply.
// Based on https://github.com/hashicorp/terraform/blob/2ad9f0513a2bd57ef49e88f9de0d2f2cc1592ff3/backend/local/hook_count.go
// Keeps track of total unique resources that are applied.
type CountHook struct {
	Added   int
	Changed int
	Removed int

	Total int

	pending   map[string]countHookAction
	resources map[string][]terraform.DiffChangeType

	sync.Mutex
	terraform.NilHook
}

// Reset resets the count internals.
func (h *CountHook) Reset() {
	h.Lock()
	defer h.Unlock()

	h.pending = nil
	h.resources = nil
	h.Added = 0
	h.Changed = 0
	h.Removed = 0
}

// PreApply is called before a single resource is applied.
func (h *CountHook) PreApply(
	n *terraform.InstanceInfo,
	s *terraform.InstanceState,
	d *terraform.InstanceDiff) (terraform.HookAction, error) {
	h.Lock()
	defer h.Unlock()

	if d.Empty() {
		return terraform.HookActionContinue, nil
	}

	if h.pending == nil {
		h.pending = make(map[string]countHookAction)
	}

	action := countHookActionChange
	if d.GetDestroy() {
		action = countHookActionRemove
	} else if s.ID == "" {
		action = countHookActionAdd
	}

	h.pending[n.HumanId()] = action

	return terraform.HookActionContinue, nil
}

// PostApply is called after a single resource is applied.
func (h *CountHook) PostApply(
	n *terraform.InstanceInfo,
	s *terraform.InstanceState,
	e error) (terraform.HookAction, error) {
	h.Lock()
	defer h.Unlock()

	if h.pending != nil {
		if a, ok := h.pending[n.HumanId()]; ok {
			delete(h.pending, n.HumanId())
			if e == nil {
				switch a {
				case countHookActionAdd:
					h.Added++
				case countHookActionChange:
					h.Changed++
				case countHookActionRemove:
					h.Removed++
				}
			}
		}
	}

	return terraform.HookActionContinue, nil
}

// PostDiff is called after a single resource is diffed.
func (h *CountHook) PostDiff(
	n *terraform.InstanceInfo, d *terraform.InstanceDiff) (
	terraform.HookAction, error) {
	h.Lock()
	defer h.Unlock()

	if h.resources == nil {
		h.resources = make(map[string][]terraform.DiffChangeType)
	}

	// We don't count anything for data sources
	if strings.HasPrefix(n.Id, "data.") {
		return terraform.HookActionContinue, nil
	}

	cts := h.resources[n.HumanId()]
	for _, c := range cts {
		if c == d.ChangeType() {
			return terraform.HookActionContinue, nil
		}
	}

	cts = append(cts, d.ChangeType())
	h.resources[n.HumanId()] = cts
	h.Total++

	return terraform.HookActionContinue, nil
}

// Progress tracks the progress of terraform actions.
type Progress struct {
	Added, Changed, Removed, Total int
}

func (s Progress) String() string {
	done := s.Added + s.Changed + s.Removed
	return fmt.Sprintf("done applying %d out of %d resources", done, s.Total)
}

func progressFromCountHook(ch *CountHook) Progress {
	ch.Lock()
	defer ch.Unlock()
	return Progress{Added: ch.Added, Changed: ch.Changed, Removed: ch.Removed, Total: ch.Total}
}
