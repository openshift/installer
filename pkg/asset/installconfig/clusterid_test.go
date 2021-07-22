package installconfig

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_generateInfraID(t *testing.T) {
	tests := []struct {
		input      string
		envInfraID string
		expLen     int
		expNonRand string
	}{{
		input:      "qwertyuiop",
		expLen:     10 + randomLen + 1,
		expNonRand: "qwertyuiop",
	}, {
		input:      "qwertyuiopasdfghjklzxcvbnm",
		expLen:     27,
		expNonRand: "qwertyuiopasdfghjklzx",
	}, {
		input:      "qwertyuiopasdfghjklz-cvbnm",
		expLen:     26,
		expNonRand: "qwertyuiopasdfghjklz",
	}, {
		input:      "qwe.rty.@iop!",
		expLen:     11 + randomLen + 1,
		expNonRand: "qwe-rty-iop",
	}}
	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			for _, expRand := range []string{"", "xxxx", "xxxxx", "xxxxa", "xxXxx", "xxxxxx"} {
				envInfraID := ""
				if expRand != "" {
					envInfraID = fmt.Sprintf("%s", expRand)
				}
				os.Setenv("OPENSHIFT_INSTALL_INFRA_ID_OVERRIDE", envInfraID)
				t.Log("OPENSHIFT_INSTALL_INFRA_ID_OVERRIDE", envInfraID)

				got, err := generateInfraID(test.input, 27)
				t.Log("InfraID", got)

				randomIDPattern := fmt.Sprintf("^([b-df-hj-np-tv-z0-9]){%d}$", randomLen)
				randomIDRegexp := regexp.MustCompile(randomIDPattern)

				// Invalid "OPENSHIFT_INSTALL_INFRA_ID_OVERRIDE" configured
				if envInfraID != "" && !randomIDRegexp.MatchString(envInfraID) {
					assert.Error(t, err, "Invalid 'OPENSHIFT_INSTALL_INFRA_ID_OVERRIDE': '%s'", envInfraID)
				} else {
					infraIDPattern := fmt.Sprintf("^%s-([b-df-hj-np-tv-z0-9]){%d}$", test.expNonRand, randomLen)
					infraIDRegexp := regexp.MustCompile(infraIDPattern)
					// No "OPENSHIFT_INSTALL_INFRA_ID_OVERRIDE" env var configured (normal flow)
					if envInfraID == "" {
						assert.Regexp(t, infraIDRegexp, got)
					}
					// Valid "OPENSHIFT_INSTALL_INFRA_ID_OVERRIDE" configured
					if infraIDRegexp.MatchString(envInfraID) {
						assert.Equal(t, test.expNonRand+"-"+envInfraID, got)

					}
					assert.Equal(t, test.expLen, len(got))
					assert.Equal(t, test.expNonRand, got[:len(got)-randomLen-1])
				}
			}
		})
	}
}
