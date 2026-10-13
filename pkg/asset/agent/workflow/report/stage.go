package workflowreport

import (
	"time"
)

type stageHeader struct {
	stageID
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`

	timeProvider timeProvider
}

type stage struct {
	stageHeader
	Result    string   `json:"result,omitempty"`
	SubStages []*stage `json:"sub_stages,omitempty"`
}

func newStage(id string, description string, tp timeProvider) *stage {
	return &stage{
		stageHeader: stageHeader{
			stageID:      newStageID(id, description),
			StartTime:    tp.Now(),
			timeProvider: tp,
		},
	}
}

func (s *stage) currSubStage() *stage {
	if len(s.SubStages) == 0 {
		return nil
	}
	return s.SubStages[len(s.SubStages)-1]
}

func (s *stage) getSubStageByID(id StageID) *stage {
	for _, s := range s.SubStages {
		if s.Equals(id) {
			return s
		}
	}
	return nil
}

func (s *stage) stop() {
	if sub := s.currSubStage(); sub != nil {
		sub.done()
	}
	s.done()
}

func (s *stage) done() {
	s.EndTime = s.timeProvider.Now()
}
