package aws

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/quota"
)

func Test_aggregate(t *testing.T) {
	cases := []struct {
		input []quota.Constraint

		exp []quota.Constraint
	}{{
		input: []quota.Constraint{
			{Name: "q1", Region: "g", Count: 1},
			{Name: "q1", Region: "g", Count: 1},
		},
		exp: []quota.Constraint{
			{Name: "q1", Region: "g", Count: 2},
		},
	}, {
		input: []quota.Constraint{
			{Name: "q1", Region: "g", Count: 1},
			{Name: "q2", Region: "r1", Count: 1},
			{Name: "q3", Region: "r1", Count: 1},
			{Name: "q2", Region: "r1", Count: 1},
		},
		exp: []quota.Constraint{
			{Name: "q1", Region: "g", Count: 1},
			{Name: "q2", Region: "r1", Count: 2},
			{Name: "q3", Region: "r1", Count: 1},
		},
	}, {
		input: []quota.Constraint{
			{Name: "q1", Region: "g", Count: 1},
			{Name: "q2", Region: "r1", Count: 1},
			{Name: "q3", Region: "r1", Count: 1},
			{Name: "q4", Region: "r1", Count: 1},
		},
		exp: []quota.Constraint{
			{Name: "q1", Region: "g", Count: 1},
			{Name: "q2", Region: "r1", Count: 1},
			{Name: "q3", Region: "r1", Count: 1},
			{Name: "q4", Region: "r1", Count: 1},
		},
	}}

	for idx, test := range cases {
		t.Run(fmt.Sprintf("test %d", idx), func(t *testing.T) {
			got := aggregate(test.input)
			assert.EqualValues(t, test.exp, got)
		})
	}
}
