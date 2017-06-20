package smoke

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func testQA(t *testing.T) {
	t.Run("StatsEmitterLogs", testGetStatsEmitterLogs)
}

func getStatsEmitterLogs(t *testing.T) error {
	c := newClient(t)
	expected := "report successfully sent"
	namespace := "tectonic-system"
	podPrefix := "tectonic-stats-emitter"
	logs, err := validatePodLogging(c, namespace, podPrefix)
	if err != nil {
		return fmt.Errorf("failed to gather logs for %s/%s, %v", namespace, podPrefix, err)
	}
	if !bytes.Contains(logs, []byte(expected)) {
		return fmt.Errorf("expected logs to contain %q", expected)
	}
	return nil
}

func testGetStatsEmitterLogs(t *testing.T) {
	max := 3 * time.Minute
	err := retry(getStatsEmitterLogs, t, 3*time.Second, max)
	if err != nil {
		t.Fatalf("Failed to verify stats-emitter success in logs in %v.", max)
	}
	t.Log("Successfully verified stats-emitter success in logs.")
}
