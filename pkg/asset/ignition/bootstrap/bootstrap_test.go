package bootstrap

import (
	"fmt"
	"os"
	"syscall"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/mock"
	"github.com/stretchr/testify/assert"
)

func TestBootstrapLoad(t *testing.T) {
	tests := []struct {
		name       string
		config     []byte
		wantFound  bool
		fetchErr   error
		wantErrMsg string
	}{
		{
			name:      "valid",
			wantFound: true,
			config:    []byte(`{"ignition":{"config":{},"security":{"tls":{}},"timeouts":{},"version":"2.2.0"},"networkd":{},"passwd":{},"storage":{},"systemd":{}}`),
		},
		{
			name:       "no version",
			wantErrMsg: "error: invalid config version",
			config:     []byte(`{"ignition":{"config":{},"security":{"tls":{}},"timeouts":{},"version":""},"networkd":{},"passwd":{},"storage":{},"systemd":{}}`),
		},
		{
			name:       "empty",
			wantErrMsg: "failed to unmarshal bootstrap.ign",
			config:     []byte(""),
		},
		{
			name:       "custom fetch error",
			fetchErr:   fmt.Errorf("this test"),
			wantErrMsg: "this test",
			config:     []byte(""),
		},
		{
			name:     "file not found error",
			fetchErr: os.NewSyscallError("open", syscall.ENOENT),
			config:   []byte(""),
		},
		{
			name: "custom content error",
			config: []byte(`{"ignition":{"config":{},"security":{"tls":{}},"timeouts":{},"version":"2.2.0"},"networkd":{},"passwd":{},
			"storage":{"files":[
				{"filesystem":"root","path":"/etc/myfile","user":{"name":"root"},"append":true,
				"contents":{"source":"data:text/plain;charset=utf-8;base64,wrong","verification":{}},"mode":420}]},
			"systemd":{}}`),
			wantErrMsg: "illegal base64 data",
		},
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileFetcher := mock.NewMockFileFetcher(mockCtrl)
			if tt.fetchErr != nil {
				fileFetcher.EXPECT().FetchByName(bootstrapIgnFilename).Return(nil, tt.fetchErr)
			} else {
				fileFetcher.EXPECT().FetchByName(bootstrapIgnFilename).
					Return(
						&asset.File{
							Filename: bootstrapIgnFilename,
							Data:     tt.config},
						nil,
					)
			}
			a := &Bootstrap{}
			gotFound, err := a.Load(fileFetcher)

			if tt.wantErrMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.Contains(t, err.Error(), tt.wantErrMsg, "incorrect error message")
			}
			assert.Equal(t, gotFound, tt.wantFound, "Found is incorrect")
		})
	}
}
