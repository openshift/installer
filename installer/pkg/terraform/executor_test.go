package terraform

import (
	"testing"

	"time"

	"io/ioutil"
	"path/filepath"

	"github.com/stretchr/testify/assert"
)

const tfTemplate = `
variable "version" {
  type = "string"
}

data "template_file" "foobar" {
  template = "foobar$${version}"

  vars {
    version = "${var.version}"
  }
}

resource "local_file" "foobar" {
  content = "${data.template_file.foobar.rendered}"
  filename = "/tmp/foobar.txt"
}

output "foobar" {
  value = "${data.template_file.foobar.rendered}"
}
`

// TestExecutorSimple executes TerraForm apply with a custom plugin, verifies it
// worked (State/Status), and then create a new executor at the path of the
// existing one and verify the state is shared.
func TestExecutorSimple(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "tectonic")
	if err != nil {
		t.Logf("Failed to create temporary directory: %s", err)
		t.FailNow()
	}

	// Create an executor.
	ex, err := NewExecutor(tmpDir)
	if err == ErrBinaryNotFound {
		t.Skip("TerraForm not found, skipping")
		return
	}
	defer ex.Cleanup()

	assert.Nil(t, err)
	assert.NotEmpty(t, ex.WorkingDirectory())
	assert.Nil(t, ex.State())

	// Add variables to it.
	err = ex.AddVariables([]byte("version = \"2000\""))
	assert.Nil(t, err)

	// Add a source file.
	mainTFPath := filepath.Join(ex.WorkingDirectory(), "main.tf")
	assert.Nil(t, ioutil.WriteFile(mainTFPath, []byte(tfTemplate), 0666))

	// Execute TerraForm init.
	id, done, err := ex.Execute("init")
	assert.Nil(t, err)
	assert.NotZero(t, id)

	// Wait for its termination.
	select {
	case <-done:
	case <-time.After(30 * time.Second):
		assert.FailNow(t, "TerraForm init timed out")
	}

	// Verify status, state and output.
	status, err := ex.Status(id)
	assert.Nil(t, err)
	assert.Equal(t, ExecutionStatusSuccess, status)

	// Execute TerraForm apply.
	id, done, err = ex.Execute("apply", "-auto-approve")
	assert.Nil(t, err)
	assert.NotZero(t, id)

	// Wait for its termination.
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		assert.FailNow(t, "TerraForm apply timed out")
	}

	// Verify status, state and output.
	status, err = ex.Status(id)
	assert.Nil(t, err)
	assert.Equal(t, ExecutionStatusSuccess, status)

	state := ex.State()
	assert.NotNil(t, state)

	output, err := ex.Output(id)
	assert.Nil(t, err)
	outputBytes, _ := ioutil.ReadAll(output)
	assert.NotZero(t, len(outputBytes))

	// Creates a new executor at the same existing one.
	ex2, err := NewExecutor(ex.WorkingDirectory())
	assert.Nil(t, err)
	assert.NotNil(t, ex2)

	state2 := ex2.State()
	if assert.NotNil(t, state2) {
		assert.Equal(t, state.Lineage, state2.Lineage)
	}
}

// TestExecutorMissingVar executes TerraForm apply with missing variables and
// ensures it failed.
func TestExecutorMissingVar(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "tectonic")
	if err != nil {
		t.Logf("Failed to create temporary directory: %s", err)
		t.FailNow()
	}

	// Create an executor.
	ex, err := NewExecutor(tmpDir)
	if err == ErrBinaryNotFound {
		t.Skip("TerraForm not found, skipping")
		return
	}
	defer ex.Cleanup()

	assert.Nil(t, err)
	assert.NotEmpty(t, ex.WorkingDirectory())
	assert.Nil(t, ex.State())

	// Add a source file.
	mainTFPath := filepath.Join(ex.WorkingDirectory(), "main.tf")
	assert.Nil(t, ioutil.WriteFile(mainTFPath, []byte(tfTemplate), 0666))

	// Execute TerraForm apply.
	id, done, err := ex.Execute("apply", "-input=false")
	assert.Nil(t, err)
	assert.NotZero(t, id)

	// Wait for its termination.
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		assert.FailNow(t, "TerraForm apply timed out")
	}

	// Verify status, state and output.
	status, err := ex.Status(id)
	assert.NotNil(t, err)
	assert.Equal(t, ExecutionStatusFailure, status)

	assert.Nil(t, ex.State())

	output, err := ex.Output(id)
	assert.Nil(t, err)
	outputBytes, _ := ioutil.ReadAll(output)
	assert.NotZero(t, len(outputBytes))
}
