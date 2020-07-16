package gcp

import "testing"

func Test_newQuota(t *testing.T) {
	cases := []struct {
		name string
	}{{
		name: "missing usage",
	}, {
		name: "usage in single zone",
	}, {
		name: "usage in multiple zones",
	}, {
		name: "",
	}}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {

		})
	}
}
