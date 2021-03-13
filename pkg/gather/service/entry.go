package service

// Entry is an entry in a service entries file
type Entry struct {
	// Timestamp is the time at which the entry was recorded
	Timestamp string `json:"timestamp"`
	// Phase is the phase of the service
	Phase Phase `json:"phase"`
	// Result is the result of either the service, the stage, the pre-command, or the post-command. This is only
	// present when the phase is an ending phase.
	Result Result `json:"result,omitempty"`
	// Stage is the name of the stage being executed. This is only present when the phase is either StageStart or StageEnd.
	Stage string `json:"string,omitempty"`
	// PreCommand is the name of the pre-command being executed. This is only present when the phase is either
	// PreCommandStart or PreCommandEnd.
	PreCommand string `json:"preCommand,omitempty"`
	// PostCommand is the name of the post-command being executed. This is only present when the phase is either
	// PostCommandStart or PostCommandEnd.
	PostCommand string `json:"postCommand,omitempty"`
	// ErrorLine is the location where the error occurred that caused the failure. This is only present when the result
	// is Failure.
	ErrorLine string `json:"errorLine,omitempty"`
	// ErrorMessage is the last few output messages from the service prior to the error. This is only present when the
	// result is Failure.
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// Phase is the phase of the service.
type Phase string

const (
	// ServiceStart is the phase when the main command of the service starts.
	ServiceStart Phase = "service start"
	// ServiceEnd is the phase when the main command of the service ends.
	ServiceEnd Phase = "service end"
	// StageStart is the phase when a stage of the service starts.
	StageStart Phase = "stage start"
	// StageEnd is the phase when a stage of the service ends.
	StageEnd Phase = "stage end"
	// PreCommandStart is the phase when a pre-command of the service starts.
	PreCommandStart Phase = "pre-command start"
	// PreCommandEnd is the phase when a pre-command of the service ends.
	PreCommandEnd Phase = "pre-command end"
	// PostCommandStart is the phase when a post-command of the service starts.
	PostCommandStart Phase = "post-command start"
	// PostCommandEnd is the phase when a post-command of the service ends.
	PostCommandEnd Phase = "post-command end"
)

// Result is the result of either the service, the stage, the pre-command, or the post-command.
type Result string

const (
	// Success indicates that the service, stage, pre-command, or post-command was successful.
	Success Result = "success"
	// Failure indicates that the service, stage, pre-command, or post-command ended due to a failure.
	Failure Result = "failure"
)
