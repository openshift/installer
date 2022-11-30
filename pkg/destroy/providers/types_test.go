package providers

import "testing"

func TestStageSucceeds(t *testing.T) {
	succeeds := func(ctx context.Context) (bool, error) {
		return true, nil
	}
	stage := Stage {
		Name: "test-stage",
		Funcs: []StageFunc{succeeds, succeeds, succeeds },
	}
}
