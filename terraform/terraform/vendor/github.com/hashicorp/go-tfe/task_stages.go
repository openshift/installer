package tfe

import (
	"context"
	"fmt"
	"time"
)

// Compile-time proof of interface  implementation
var _ TaskStages = (*taskStages)(nil)

// TaskStages describes all the task stage related methods that the TFC/E API
// supports.
// **Note: This API is still in BETA and is subject to change.**
type TaskStages interface {
	// Read a task stage by ID
	Read(ctx context.Context, taskStageID string, options *TaskStageReadOptions) (*TaskStage, error)

	// List all task stages for a given rrun
	List(ctx context.Context, runID string, options *TaskStageListOptions) (*TaskStageList, error)
}

// taskStages implements TaskStages
type taskStages struct {
	client *Client
}

// Stage is an enum that represents the possible run stages for run tasks
type Stage string

const (
	PrePlan  Stage = "pre_plan"
	PostPlan Stage = "post_plan"
)

// TaskStage represents a TFC/E run's stage where run tasks can occur
type TaskStage struct {
	ID               string                    `jsonapi:"primary,task-stages"`
	Stage            Stage                     `jsonapi:"attr,stage"`
	StatusTimestamps TaskStageStatusTimestamps `jsonapi:"attr,status-timestamps"`
	CreatedAt        time.Time                 `jsonapi:"attr,created-at,iso8601"`
	UpdatedAt        time.Time                 `jsonapi:"attr,updated-at,iso8601"`

	Run         *Run          `jsonapi:"relation,run"`
	TaskResults []*TaskResult `jsonapi:"relation,task-results"`
}

// TaskStageList represents a list of task stages
type TaskStageList struct {
	*Pagination
	Items []*TaskStage
}

// TaskStageStatusTimestamps represents the set of timestamps recorded for a task stage
type TaskStageStatusTimestamps struct {
	ErroredAt  time.Time `jsonapi:"attr,errored-at,rfc3339"`
	RunningAt  time.Time `jsonapi:"attr,running-at,rfc3339"`
	CanceledAt time.Time `jsonapi:"attr,canceled-at,rfc3339"`
	FailedAt   time.Time `jsonapi:"attr,failed-at,rfc3339"`
	PassedAt   time.Time `jsonapi:"attr,passed-at,rfc3339"`
}

// TaskStageIncludeOpt represents the available options for include query params.
type TaskStageIncludeOpt string

const TaskStageTaskResults TaskStageIncludeOpt = "task_results"

// TaskStageReadOptions represents the set of options when reading a task stage
type TaskStageReadOptions struct {
	// Optional: A list of relations to include.
	Include []TaskStageIncludeOpt `url:"include,omitempty"`
}

// TaskStageListOptions represents the options for listing task stages for a run
type TaskStageListOptions struct {
	ListOptions
}

// Read a task stage by ID
func (s *taskStages) Read(ctx context.Context, taskStageID string, options *TaskStageReadOptions) (*TaskStage, error) {
	if !validStringID(&taskStageID) {
		return nil, ErrInvalidTaskStageID
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("task-stages/%s", taskStageID)
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	t := &TaskStage{}
	err = req.Do(ctx, t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// List task stages for a run
func (s *taskStages) List(ctx context.Context, runID string, options *TaskStageListOptions) (*TaskStageList, error) {
	if !validStringID(&runID) {
		return nil, ErrInvalidRunID
	}

	u := fmt.Sprintf("runs/%s/task-stages", runID)
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	tlist := &TaskStageList{}

	err = req.Do(ctx, tlist)
	if err != nil {
		return nil, err
	}

	return tlist, nil
}

func (o *TaskStageReadOptions) valid() error {
	if o == nil {
		return nil // nothing to validate
	}

	if err := validateTaskStageIncludeParams(o.Include); err != nil {
		return err
	}

	return nil
}

func validateTaskStageIncludeParams(params []TaskStageIncludeOpt) error {
	for _, p := range params {
		switch p {
		case TaskStageTaskResults:
			// do nothing
		default:
			return ErrInvalidIncludeValue
		}
	}

	return nil
}
