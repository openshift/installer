package fencing

import (
	"testing"
)

func TestParseStonithEnabled(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect bool
	}{
		{
			name:   "enabled with colon",
			input:  " stonith-enabled: true\n",
			expect: true,
		},
		{
			name:   "enabled with equals",
			input:  "stonith-enabled=true\n",
			expect: true,
		},
		{
			name:   "disabled",
			input:  " stonith-enabled: false\n",
			expect: false,
		},
		{
			name:   "mixed case",
			input:  " Stonith-Enabled: True\n",
			expect: true,
		},
		{
			name:   "empty",
			input:  "",
			expect: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := parseStonithEnabled(tc.input)
			if got != tc.expect {
				t.Errorf("parseStonithEnabled(%q) = %v, want %v", tc.input, got, tc.expect)
			}
		})
	}
}

func TestParsePacemakerOnline(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect []string
	}{
		{
			name:   "two nodes with brackets",
			input:  "Online: [ node-a node-b ]\nStandby:\n",
			expect: []string{"node-a", "node-b"},
		},
		{
			name:   "two nodes without brackets",
			input:  "Online: node-a node-b\n",
			expect: []string{"node-a", "node-b"},
		},
		{
			name:   "one node",
			input:  "Online: [ node-a ]\n",
			expect: []string{"node-a"},
		},
		{
			name:   "indented",
			input:  "  Online: [ master-0.example.com master-1.example.com ]\n",
			expect: []string{"master-0.example.com", "master-1.example.com"},
		},
		{
			name:   "no online line",
			input:  "Standby: node-a\n",
			expect: nil,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := parsePacemakerOnline(tc.input)
			if len(got) != len(tc.expect) {
				t.Fatalf("parsePacemakerOnline() = %v, want %v", got, tc.expect)
			}
			for i := range got {
				if got[i] != tc.expect[i] {
					t.Errorf("parsePacemakerOnline()[%d] = %q, want %q", i, got[i], tc.expect[i])
				}
			}
		})
	}
}

func TestNodeInOnlineList(t *testing.T) {
	node := NodeInfo{Name: "master-0.example.com", PacemakerName: "master-0"}
	tests := []struct {
		name   string
		online []string
		expect bool
	}{
		{
			name:   "exact K8s name match",
			online: []string{"master-0.example.com", "master-1.example.com"},
			expect: true,
		},
		{
			name:   "exact pacemaker name match",
			online: []string{"master-0", "master-1"},
			expect: true,
		},
		{
			name:   "short hostname match",
			online: []string{"master-0.other.domain", "master-1.other.domain"},
			expect: true,
		},
		{
			name:   "no match",
			online: []string{"worker-0", "worker-1"},
			expect: false,
		},
		{
			name:   "empty list",
			online: nil,
			expect: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := nodeInOnlineList(node, tc.online)
			if got != tc.expect {
				t.Errorf("nodeInOnlineList() = %v, want %v", got, tc.expect)
			}
		})
	}
}

func TestParseDaemonStatus(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		missing []string
	}{
		{
			name: "all active",
			input: `Full List of Resources:
Daemon Status:
  corosync: active/enabled
  pacemaker: active/enabled
  pcsd: active/enabled
`,
			missing: nil,
		},
		{
			name: "pacemaker inactive",
			input: `Daemon Status:
  corosync: active/running
  pacemaker: inactive/disabled
  pcsd: active/running
`,
			missing: []string{"pacemaker"},
		},
		{
			name:    "no daemon section",
			input:   "some other output\n",
			missing: nil,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, _ := parseDaemonStatus(tc.input)
			if len(got) != len(tc.missing) {
				t.Fatalf("parseDaemonStatus() = %v, want %v", got, tc.missing)
			}
			for i := range got {
				if got[i] != tc.missing[i] {
					t.Errorf("parseDaemonStatus()[%d] = %q, want %q", i, got[i], tc.missing[i])
				}
			}
		})
	}
}

func TestParseEtcdHealth(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "both healthy",
			input:   `[{"health":true,"took":"10ms"},{"health":true,"took":"12ms"}]`,
			wantErr: false,
		},
		{
			name:    "one unhealthy",
			input:   `[{"health":true},{"health":false,"error":"context deadline exceeded"}]`,
			wantErr: true,
		},
		{
			name:    "empty array",
			input:   `[]`,
			wantErr: true,
		},
		{
			name:    "single endpoint",
			input:   `[{"health":true}]`,
			wantErr: true,
		},
		{
			name:    "invalid json",
			input:   `not json`,
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := parseEtcdHealth(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("parseEtcdHealth() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestParseEtcdMembers(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		ipA     string
		ipB     string
		wantErr bool
	}{
		{
			name: "two voters",
			input: `{"members":[
				{"isLearner":false,"clientURLs":["https://10.0.0.1:2379"]},
				{"isLearner":false,"clientURLs":["https://10.0.0.2:2379"]}
			]}`,
			ipA:     "10.0.0.1",
			ipB:     "10.0.0.2",
			wantErr: false,
		},
		{
			name: "one learner",
			input: `{"members":[
				{"isLearner":false,"clientURLs":["https://10.0.0.1:2379"]},
				{"isLearner":true,"clientURLs":["https://10.0.0.2:2379"]}
			]}`,
			ipA:     "10.0.0.1",
			ipB:     "10.0.0.2",
			wantErr: true,
		},
		{
			name: "missing node B",
			input: `{"members":[
				{"isLearner":false,"clientURLs":["https://10.0.0.1:2379"]}
			]}`,
			ipA:     "10.0.0.1",
			ipB:     "10.0.0.2",
			wantErr: true,
		},
		{
			name:    "invalid json",
			input:   `bad`,
			ipA:     "10.0.0.1",
			ipB:     "10.0.0.2",
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := parseEtcdMembers(tc.input, tc.ipA, tc.ipB)
			if (err != nil) != tc.wantErr {
				t.Errorf("parseEtcdMembers() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestShellQuote(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{
			name:   "simple string",
			input:  "pcs status nodes",
			expect: "'pcs status nodes'",
		},
		{
			name:   "string with single quotes",
			input:  "echo 'hello'",
			expect: "'echo '\"'\"'hello'\"'\"''",
		},
		{
			name:   "empty string",
			input:  "",
			expect: "''",
		},
		{
			name:   "subshell syntax",
			input:  "$(whoami)",
			expect: "'$(whoami)'",
		},
		{
			name:   "semicolon and pipe",
			input:  "cmd1; cmd2 | cmd3",
			expect: "'cmd1; cmd2 | cmd3'",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := shellQuote(tc.input)
			if got != tc.expect {
				t.Errorf("shellQuote(%q) = %q, want %q", tc.input, got, tc.expect)
			}
		})
	}
}

func TestFormatEtcdURL(t *testing.T) {
	tests := []struct {
		name   string
		ip     string
		expect string
	}{
		{
			name:   "IPv4 address",
			ip:     "10.0.0.1",
			expect: "https://10.0.0.1:2379",
		},
		{
			name:   "IPv6 address",
			ip:     "fd00::1",
			expect: "https://[fd00::1]:2379",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := formatEtcdURL(tc.ip)
			if got != tc.expect {
				t.Errorf("formatEtcdURL(%q) = %q, want %q", tc.ip, got, tc.expect)
			}
		})
	}
}
