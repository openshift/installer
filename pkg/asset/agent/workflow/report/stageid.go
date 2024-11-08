package workflowreport

type stageID struct {
	Identifier string `json:"id"`
	Desc       string `json:"description,omitempty"`
}

// NewStageID creates a new StageID instance.
func NewStageID(id, desc string) StageID {
	return newStageID(id, desc)
}

// NewStageID creates a new stage id const.
func newStageID(id, desc string) stageID {
	return stageID{
		Identifier: id,
		Desc:       desc,
	}
}

// ID is the stage id.
func (s stageID) ID() string {
	return s.Identifier
}

// Description is the stage longer description.
func (s stageID) Description() string {
	return s.Desc
}

// String returns the stage id.
func (s stageID) String() string {
	return s.ID()
}

// Equals performs a comparison against another StageID instance.
func (s stageID) Equals(other StageID) bool {
	return other != nil && s.ID() == other.ID()
}
