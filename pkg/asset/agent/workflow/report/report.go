package workflowreport

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type report struct {
	stageHeader
	Stages []*stage      `json:"stages,omitempty"`
	Result *reportResult `json:"result"`
	writer io.Writer
}

type reportResult struct {
	ExitCode             int    `json:"exit_code"`
	ErrorMessage         string `json:"error_message,omitempty"`
	DetailedErrorMessage string `json:"detailed_error_message,omitempty"`
}

type timeProvider interface {
	Now() time.Time
}

type defaultTimeProvider struct{}

func (dtp *defaultTimeProvider) Now() time.Time {
	return time.Now().UTC()
}

type defaultWriter struct {
	reportFileName string
}

func (dw defaultWriter) Write(data []byte) (n int, err error) {
	return len(data), os.WriteFile(dw.reportFileName, data, 0600)
}

func newReport(id string, writer io.Writer, tp timeProvider) Report {
	return &report{
		stageHeader: stageHeader{
			stageID: stageID{
				Identifier: fmt.Sprintf("report-%s-%s", id, tp.Now().Format("200601021504")),
				Desc:       "",
			},

			timeProvider: tp,
		},
		writer: writer,
	}
}

func (r *report) update() error {
	_, err := r.writer.Write(r.json())
	return err
}

func (r *report) json() []byte {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return []byte(fmt.Sprintf(`{"errorMessage": "%s"}`, fmt.Errorf("error while generating the workflow report: %w", err)))
	}
	return data
}

func (r *report) stageKeys(stageID StageID) (string, string) {
	parts := strings.Split(stageID.ID(), ".")
	if len(parts) > 1 {
		return parts[0], parts[1]
	}
	return parts[0], ""
}

func (r *report) currStage() *stage {
	if len(r.Stages) == 0 {
		return nil
	}
	return r.Stages[len(r.Stages)-1]
}

func (r *report) getStageByID(id StageID) *stage {
	// id may be a substage id
	stageID, _ := r.stageKeys(id)

	for _, s := range r.Stages {
		if s.ID() == stageID {
			return s
		}
	}
	return nil
}

func (r *report) validateStageID(id StageID) error {
	sID, subID := r.stageKeys(id)
	if sID == "" {
		return fmt.Errorf("empty stage name")
	}
	if subID != "" {
		return fmt.Errorf("invalid stage name: %s", id)
	}
	return nil
}

// Stage adds a new stage to the report and makes it active.
func (r *report) Stage(id StageID) error {
	if err := r.validateStageID(id); err != nil {
		return fmt.Errorf("cannot add stage: %w", err)
	}

	if s := r.currStage(); s != nil {
		if s.Equals(id) {
			return fmt.Errorf("cannot add stage: %s is already active", id)
		}
		s.stop()
	}

	return r.addStage(&r.Stages, id)
}

func (r *report) addStage(stages *[]*stage, id StageID) error {
	s := newStage(id.ID(), id.Description(), r.timeProvider)
	if r.currStage() == nil {
		r.StartTime = s.StartTime
	}
	*stages = append(*stages, s)
	return r.update()
}

func (r *report) setResult(s *stage, value string) error {
	s.Result = value
	return r.update()
}

// StageResult adds a result for the given stage.
func (r *report) StageResult(id StageID, value string) error {
	if err := r.validateStageID(id); err != nil {
		return fmt.Errorf("cannot add result for stage %s: %w", id, err)
	}
	s := r.getStageByID(id)
	if s == nil {
		return fmt.Errorf("cannot add result for stage: %s not found", id)
	}
	return r.setResult(s, value)
}

func (r *report) validateSubStageID(id StageID) error {
	sID, subID := r.stageKeys(id)
	if sID == "" {
		return fmt.Errorf("empty stage name")
	}
	if subID == "" {
		return fmt.Errorf("empty substage name: %s", id)
	}
	return nil
}

func (r *report) getSubStageOwner(id StageID) (*stage, error) {
	if err := r.validateSubStageID(id); err != nil {
		return nil, err
	}
	owner := r.getStageByID(id)
	if owner == nil {
		return nil, fmt.Errorf("%s stage owner not found", id)
	}
	return owner, nil
}

// SubStage adds a new substage to the report and makes it active.
func (r *report) SubStage(id StageID) error {
	owner, err := r.getSubStageOwner(id)
	if err != nil {
		return fmt.Errorf("cannot add substage: %w", err)
	}
	if !owner.Equals(r.currStage()) {
		return fmt.Errorf("cannot add substage: %s stage owner is not active", id)
	}

	if s := owner.currSubStage(); s != nil {
		if s.Equals(id) {
			return fmt.Errorf("cannot add substage: substage %s already active", id)
		}
		s.stop()
	}

	return r.addStage(&owner.SubStages, id)
}

// SubStageResult adds a result for the given substage.
func (r *report) SubStageResult(id StageID, value string) error {
	owner, err := r.getSubStageOwner(id)
	if err != nil {
		return fmt.Errorf("cannot add substage result: %w", err)
	}

	sub := owner.getSubStageByID(id)
	if sub == nil {
		return fmt.Errorf("cannot add result for substage: %s substage not found", id)
	}
	return r.setResult(sub, value)
}

// Complete marks the report as done.
func (r *report) Complete(err error) error {
	r.Result = &reportResult{}

	if err != nil {
		r.Result.ExitCode = 1
		r.Result.DetailedErrorMessage = err.Error()

		for {
			unwrappedErr := errors.Unwrap(err)
			if unwrappedErr == nil {
				break
			}
			err = unwrappedErr
		}
		r.Result.ErrorMessage = err.Error()
	}

	if s := r.currStage(); s != nil {
		s.stop()
		r.EndTime = s.EndTime
	}
	return r.update()
}
