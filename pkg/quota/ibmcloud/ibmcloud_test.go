package ibmcloud

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/quota"
)

func Test_aggregate(t *testing.T) {
	cases := []struct {
		input []quota.Constraint
		exp   []quota.Constraint
	}{{
		input: []quota.Constraint{
			{Name: "is/instance", Region: "us-east", Count: 3},
			{Name: "is/instance", Region: "us-east", Count: 1},
		},
		exp: []quota.Constraint{
			{Name: "is/instance", Region: "us-east", Count: 4},
		},
	}, {
		input: []quota.Constraint{
			{Name: "is/floating-ip", Region: "us-east", Count: 5},
			{Name: "is/instance", Region: "us-east", Count: 3},
			{Name: "is/load-balancer", Region: "us-east", Count: 2},
			{Name: "is/instance", Region: "us-east", Count: 1},
		},
		exp: []quota.Constraint{
			{Name: "is/floating-ip", Region: "us-east", Count: 5},
			{Name: "is/instance", Region: "us-east", Count: 4},
			{Name: "is/load-balancer", Region: "us-east", Count: 2},
		},
	}, {
		input: []quota.Constraint{
			{Name: "is/floating-ip", Region: "us-east", Count: 5},
			{Name: "is/instance", Region: "us-east", Count: 3},
			{Name: "is/load-balancer", Region: "us-east", Count: 2},
			{Name: "is/security-group", Region: "us-east", Count: 6},
		},
		exp: []quota.Constraint{
			{Name: "is/floating-ip", Region: "us-east", Count: 5},
			{Name: "is/instance", Region: "us-east", Count: 3},
			{Name: "is/load-balancer", Region: "us-east", Count: 2},
			{Name: "is/security-group", Region: "us-east", Count: 6},
		},
	}, {
		input: nil,
		exp:   nil,
	}}

	for idx, test := range cases {
		t.Run(fmt.Sprintf("test %d", idx), func(t *testing.T) {
			got := aggregate(test.input)
			assert.EqualValues(t, test.exp, got)
		})
	}
}
