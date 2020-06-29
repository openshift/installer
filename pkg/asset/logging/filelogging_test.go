package logging

import (
	"testing"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/machines"
	"github.com/stretchr/testify/assert"
)

func getAsset(filename string) *asset.File {
	return &asset.File{
		Filename: filename,
	}
}

func TestLogFilesChanged(t *testing.T) {
	cases := []struct {
		name                  string
		assets                []asset.WritableAsset
		cmdName               string
		directory             string
		expectedGenerationLog string
	}{
		{
			name:                  "test empty assets list",
			assets:                []asset.WritableAsset{},
			cmdName:               "test",
			directory:             "test",
			expectedGenerationLog: "",
		},
		{
			name: "test asset with one file",
			assets: []asset.WritableAsset{
				&installconfig.InstallConfig{
					File: &asset.File{
						Filename: "a.yaml",
					},
				},
			},
			cmdName:               "test install config",
			directory:             "test/",
			expectedGenerationLog: "Test Install Config created in: test",
		},
		{
			name: "test asset with two files, same directory",
			assets: []asset.WritableAsset{
				&machines.Master{
					MachineFiles: []*asset.File{
						getAsset("a.yaml"),
						getAsset("b.yaml"),
					},
				},
			},
			cmdName:               "machines",
			directory:             "test",
			expectedGenerationLog: "Machines created in: test",
		},
		{
			name: "test asset with two files, two directories",
			assets: []asset.WritableAsset{
				&machines.Master{
					MachineFiles: []*asset.File{
						getAsset("a.yaml"),
						getAsset("machines/b.yaml"),
					},
				},
			},
			cmdName:               "machines",
			directory:             "test",
			expectedGenerationLog: "Machines created in: test and test/machines",
		},
		{
			name: "test asset with two files, two directories, but three entries (same file twice",
			assets: []asset.WritableAsset{
				&machines.Master{
					MachineFiles: []*asset.File{
						getAsset("a.yaml"),
						getAsset("machines/b.yaml"),
						getAsset("machines/b.yaml"),
					},
				},
			},
			cmdName:               "machines",
			directory:             "test",
			expectedGenerationLog: "Machines created in: test and test/machines",
		},
		{
			name: "test asset with three files, two directories",
			assets: []asset.WritableAsset{
				&machines.Master{
					MachineFiles: []*asset.File{
						getAsset("a.yaml"),
						getAsset("machines/b.yaml"),
						getAsset("machines/c.yaml"),
					},
				},
			},
			cmdName:               "machine-config",
			directory:             "test/",
			expectedGenerationLog: "Machine-Config created in: test and test/machines",
		},
		{
			name: "test asset with three files, three directories",
			assets: []asset.WritableAsset{
				&machines.Master{
					MachineFiles: []*asset.File{
						getAsset("a.yaml"),
						getAsset("machines/b.yaml"),
						getAsset("control-plane/c.yaml"),
					},
				},
			},
			cmdName:               "machines",
			directory:             "test",
			expectedGenerationLog: "Machines created in: test, test/control-plane and test/machines",
		},
		{
			name: "test asset with five files, five directories",
			assets: []asset.WritableAsset{
				&machines.Master{
					MachineFiles: []*asset.File{
						getAsset("a.yaml"),
						getAsset("machines/b.yaml"),
						getAsset("control-plane/c.yaml"),
						getAsset("master/c.yaml"),
						getAsset("worker/c.yaml"),
					},
				},
			},
			cmdName:               "machines",
			directory:             "test",
			expectedGenerationLog: "Machines created in: test",
		},
		{
			name: "test asset with five files, five nested directories",
			assets: []asset.WritableAsset{
				&machines.Master{
					MachineFiles: []*asset.File{
						getAsset("machines/workers/b.yaml"),
						getAsset("control-plane/disks/c.yaml"),
						getAsset("control-plane/VM/d.yaml"),
						getAsset("machines/configurations/c.yaml"),
					},
				},
			},
			cmdName:               "machines",
			directory:             "test",
			expectedGenerationLog: "Machines created in: test/control-plane and test/machines",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			textOutput := LogCreatedFiles(tc.cmdName, tc.directory, tc.assets)
			assert.EqualValues(t, tc.expectedGenerationLog, textOutput)
		})

	}
}
