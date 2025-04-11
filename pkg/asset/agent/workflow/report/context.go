package workflowreport

import (
	"context"
	"path/filepath"
)

type contextReportKey string

const (
	reportKey = contextReportKey("report")
)

// Context creates a new context with a Report value.
func Context(id string, dir string) context.Context {
	w := &defaultWriter{
		reportFileName: filepath.Join(dir, "report.json"),
	}
	tp := &defaultTimeProvider{}

	return context.WithValue(
		context.Background(),
		reportKey,
		newReport(id, w, tp))
}

// GetReport retrieves the Report instance from the given context.
// If not present, a nullRepo will be returned.
func GetReport(ctx context.Context) Report {
	if report, ok := ctx.Value(reportKey).(Report); ok {
		return report
	}
	return nullRepo{}
}

// nullRepo allows to avoid additional checks when a Report is
// required, by providing a fully valid (but essentially empty)
// Report instance.
type nullRepo struct{}

// Stage returns an empty stage.
func (nullRepo) Stage(id StageID) error {
	return nil
}

// StageResult supports result setting.
func (nullRepo) StageResult(id StageID, value string) error {
	return nil
}

// SubStage returns an empty substage.
func (nullRepo) SubStage(subID StageID) error {
	return nil
}

// SubStageResult supports result setting.
func (nullRepo) SubStageResult(subID StageID, value string) error {
	return nil
}

// Complete marks the report as done.
func (nullRepo) Complete(err error) error {
	return nil
}
