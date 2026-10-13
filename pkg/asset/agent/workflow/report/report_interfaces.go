package workflowreport

// Report is used to generate a detailed activity of the workflow.
// A report is composed by a number of sequential stages.
// A stage identifies a specific workflow phase, and can provide a
// detailed result as an output of its activity.
// A stage can be composed by any number of sequential substages.
// Any update on the report is immediately serialized on the disk.
type Report interface {
	// Creates a new stage with the given id. The previously active
	// stage (and its active substage) is marked as completed.
	// The newly created stage becomes the active one.
	Stage(id StageID) error
	// Store the given result for the specified stage.
	StageResult(id StageID, value string) error
	// Creates a new substage with the given id.
	// The subStageID must include the stage id owner of the current substage,
	// in the form of `<stage id>.<sub id>`, and the specified
	// stage must be active.
	// The previously active substage is marked as completed.
	// The newly created substage becomes the active one.
	SubStage(subID StageID) error
	// Store the given result for the specified substage.
	SubStageResult(subID StageID, value string) error
	// Marks the report as done. In case of error, exit code and error messages
	// are store in the report.
	// The currently active stage (and its active substage) is marked as completed.
	Complete(err error) error
}

// StageID is an interface for the stage identifier.
type StageID interface {
	// The id of the stage.
	ID() string
	// A longer human-readable description of the stage.
	Description() string
}
