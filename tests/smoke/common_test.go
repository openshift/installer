package smoke

import (
	"fmt"
	"testing"
	"time"
)

func testCommon(t *testing.T) {
	// verify that api server is up in 15 min
	t.Run("APIAvailable", testAPIAvailable)
}

func testAPIAvailable(t *testing.T) {
	wait := 15 * time.Minute
	err := retry(apiAvailable, t, 10*time.Second, wait)
	if err != nil {
		t.Fatalf("Failed to connect to API server in %v.", wait)
	}
	t.Log("API server is available.")
}

func apiAvailable(t *testing.T) error {
	client, _ := newClient(t)
	_, err := client.ServerVersion()
	if err != nil {
		return fmt.Errorf("failed to connect to API server: %v", err)
	}
	return nil
}
