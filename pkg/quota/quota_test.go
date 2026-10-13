package quota

import "testing"

func TestCheck(t *testing.T) {
	cases := []struct {
		name string
	}{{
		name: "missing quota",
	}, {
		name: "quota with low availability",
	}, {
		name: "quota with no availability",
	}, {
		name: "available quota",
	}, {
		name: "available quota, quota with low availability",
	}, {
		name: "available quota, quota with no availability",
	}, {
		name: "available quota, missing quota",
	}}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {

		})
	}
}
