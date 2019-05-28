package exec

import (
	"bytes"
	"os"

	"github.com/hashicorp/terraform/states/statefile"
	"github.com/pkg/errors"
)

// ReadState reads the terraform state from file and returns the contents in bytes
// It returns an error if reading the state was unsuccessful
// ReadState utilizes the terraform's internal wiring to upconvert versions of terraform state to return
// the state it currently recognizes.
func ReadState(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open %q", file)
	}
	defer f.Close()

	sf, err := statefile.Read(f)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read statefile from %q", file)
	}

	out := bytes.Buffer{}
	if err := statefile.Write(sf, &out); err != nil {
		return nil, errors.Wrapf(err, "failed to write statefile")
	}
	return out.Bytes(), nil
}
