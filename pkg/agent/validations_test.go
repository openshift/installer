package agent

import (
	"regexp"
	"testing"

	"github.com/openshift/assisted-service/models"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestCheckHostsValidation(t *testing.T) {
	tests := []struct {
		name           string
		hosts          []*models.Host
		expectedResult bool
		expectedLogs   []string
	}{
		{
			name:           "no-validations",
			expectedResult: true,
		},
		{
			name: "no-failures",
			hosts: []*models.Host{
				{
					RequestedHostname: "master-0.ostest.test.metalkube.org",
					ValidationsInfo:   "{\"hardware\":[{\"id\":\"has-inventory\",\"status\":\"success\",\"message\":\"Valid inventory exists for the host\"}]}",
				},
			},
			expectedResult: true,
		},
		{
			name: "single-host-failure",
			hosts: []*models.Host{
				{
					RequestedHostname: "master-0.ostest.test.metalkube.org",
					ValidationsInfo:   `{"hardware":[{"id":"has-min-valid-disks","status":"failure","message":"No eligible disks were found, please check specific disks to see why they are not eligible"},{"id":"has-cpu-cores-for-role","status":"success","message":"Sufficient CPU cores for role master"},{"id":"has-memory-for-role","status":"success","message":"Sufficient RAM for role master"}]}`,
				},
			},
			expectedResult: false,
			expectedLogs: []string{
				`level=error msg="Validation failure found for master\-0.ostest.test.metalkube.org" category=hardware label="Minimum disks of required size" message="No eligible disks were found, please check specific disks to see why they are not eligible"`,
			},
		},
		{
			name: "multiple-hosts-failure",
			hosts: []*models.Host{
				{
					RequestedHostname: "master-0.ostest.test.metalkube.org",
					ValidationsInfo:   `{"hardware":[{"id":"has-min-valid-disks","status":"failure","message":"No eligible disks were found, please check specific disks to see why they are not eligible"},{"id":"has-cpu-cores-for-role","status":"success","message":"Sufficient CPU cores for role master"},{"id":"has-memory-for-role","status":"success","message":"Sufficient RAM for role master"}]}`,
				},
				{
					RequestedHostname: "master-1.ostest.test.metalkube.org",
					ValidationsInfo:   `{"hardware":[{"id":"has-min-valid-disks","status":"failure","message":"No eligible disks were found, please check specific disks to see why they are not eligible"},{"id":"has-cpu-cores-for-role","status":"success","message":"Sufficient CPU cores for role master"},{"id":"has-memory-for-role","status":"success","message":"Sufficient RAM for role master"}]}`,
				},
			},
			expectedResult: false,
			expectedLogs: []string{
				`level=error msg="Validation failure found for master\-0.ostest.test.metalkube.org" category=hardware label="Minimum disks of required size" message="No eligible disks were found, please check specific disks to see why they are not eligible"`,
				`level=error msg="Validation failure found for master\-1.ostest.test.metalkube.org" category=hardware label="Minimum disks of required size" message="No eligible disks were found, please check specific disks to see why they are not eligible"`,
			},
		},
		{
			name: "malformed-json",
			hosts: []*models.Host{
				{
					RequestedHostname: "master-0.ostest.test.metalkube.org",
					ValidationsInfo:   `not a valid info`,
				},
			},
			expectedResult: false,
			expectedLogs: []string{
				`Unable to verify cluster hosts validations`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cluster := &models.Cluster{
				Hosts: tt.hosts,
			}

			logger, hook := test.NewNullLogger()
			assert.Equal(t, tt.expectedResult, checkHostsValidations(cluster, logger))

			assert.Equal(t, len(tt.expectedLogs), len(hook.Entries))
			for _, expectedMsg := range tt.expectedLogs {

				matchFound := false
				for _, s := range hook.AllEntries() {
					logLine, err := s.String()
					assert.NoError(t, err)
					if regexp.MustCompile(expectedMsg).Match([]byte(logLine)) {
						matchFound = true
					}
				}
				assert.True(t, matchFound, "Unable to find log trace for `%s`", expectedMsg)
			}
		})
	}
}
