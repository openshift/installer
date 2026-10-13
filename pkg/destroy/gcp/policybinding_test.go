package gcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_policyMemberToEmail(t *testing.T) {
	cases := []struct {
		member string
		email  string
	}{{
		member: "serviceAccount:operator@project",
		email:  "operator@project",
	}, {
		member: "deleted:serviceAccount:operator@project",
		email:  "operator@project",
	}, {
		member: "deleted:serviceAccount:operator@project?uid=1231243234",
		email:  "operator@project",
	}, {
		member: "user:user@project",
		email:  "user:user@project",
	}, {
		member: "deleted:user:user@project",
		email:  "user:user@project",
	}, {
		member: "deleted:user:user@project?uid=1232131243",
		email:  "user:user@project",
	}}
	for _, test := range cases {
		t.Run("", func(t *testing.T) {
			email := policyMemberToEmail(test.member)
			assert.Equal(t, email, test.email)
		})
	}
}
