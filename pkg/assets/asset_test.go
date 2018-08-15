package assets

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssetWriteRead(t *testing.T) {
	ctx := context.Background()
	tempDir, err := ioutil.TempDir("", "openshift-install-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	asset := &Asset{
		Parents: []Reference{{
			Name: "a",
			Hash: []byte("\x00\x01\x02\x03"),
		}},
		Name: "b",
		Data: []byte("b-data"),
	}

	err = asset.Write(ctx, tempDir, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	recovered := &Asset{Name: asset.Name}
	err = recovered.Read(ctx, tempDir, t)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, asset, recovered)
}
