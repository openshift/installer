package server

import (
	"bytes"
	"strings"
	"testing"
	"unicode"

	"k8s.io/kubernetes/pkg/apis/abac"
)

func TestPolicyEncoding(t *testing.T) {
	dedent := func(s string) string {
		lines := strings.Split(s, "\n")
		var b bytes.Buffer
		for _, line := range lines {
			line := strings.TrimLeftFunc(line, unicode.IsSpace)
			if line == "" {
				continue
			}
			b.WriteString(line)
			b.WriteRune('\n')
		}
		return b.String()
	}

	tests := []struct {
		policies []abac.Policy
		want     string
	}{
		{
			policies: []abac.Policy{
				{
					Spec: abac.PolicySpec{
						User:            "*",
						Group:           "*",
						NonResourcePath: "*",
					},
				},
			},
			want: `
{"kind":"Policy","apiVersion":"abac.authorization.kubernetes.io/v1beta1","spec":{"user":"*","group":"*","nonResourcePath":"*"}}
			`,
		},
		{
			policies: []abac.Policy{
				{
					Spec: abac.PolicySpec{
						User:            "*",
						Group:           "*",
						NonResourcePath: "*",
					},
				},
				{
					Spec: abac.PolicySpec{
						User:            "foo",
						Group:           "bar",
						Resource:        "pods",
						NonResourcePath: "*",
					},
				},
			},
			want: `
{"kind":"Policy","apiVersion":"abac.authorization.kubernetes.io/v1beta1","spec":{"user":"*","group":"*","nonResourcePath":"*"}}
{"kind":"Policy","apiVersion":"abac.authorization.kubernetes.io/v1beta1","spec":{"user":"foo","group":"bar","resource":"pods","nonResourcePath":"*"}}
			`,
		},
	}

	for i, tc := range tests {
		jsonl, err := abacPolicyToJSONL(tc.policies)
		if err != nil {
			t.Errorf("case %d: %v", i, err)
			continue
		}
		got := string(jsonl)
		want := dedent(tc.want)
		if got != want {
			t.Errorf("case %d: got:\n`%s`\nwant:\n`%s`", i, got, want)
		}
	}
}
