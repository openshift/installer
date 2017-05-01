package server

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"
)

type dockercfg map[string]struct {
	Auth  string `json:"auth"`
	Email string `json:"email"`
}

type dockercfgjson struct {
	Auths dockercfg `json:"auths"`
}

func TestTectonicVersion(t *testing.T) {
	data, err := ioutil.ReadFile("../assets/self-hosted/kube-apiserver.yaml.tmpl")
	if err != nil {
		t.Fatal(err)
	}

	wantVersion := strings.Split(tectonicVersion, "-")[0]

	if !strings.Contains(string(data), wantVersion) {
		t.Errorf("mismatch of tectonic version (%s) and API server version", wantVersion)
	}
}

func TestGetDockercfgFormat(t *testing.T) {
	v1cfg := dockercfg{
		"quay.io": {
			Auth:  "ZmFrZXVzZXI6ZmFrZXBhc3MK",
			Email: "doesntmatter@example.com",
		},
	}
	v2cfg := dockercfgjson{
		Auths: v1cfg,
	}

	v1cfgJSON, err := json.Marshal(v1cfg)
	if err != nil {
		t.Fatalf("unable marshal dockercfg v1, err: %v", err)
	}
	v2cfgJSON, err := json.Marshal(v2cfg)
	if err != nil {
		t.Fatalf("unable marshal dockercfg v2, err: %v", err)
	}

	format, err := getDockercfgFormat(v1cfgJSON)
	if err != nil {
		t.Errorf("unexpected error when getting docker format for dockercfg v1: %v", err)
	}
	if format != dockerformatv1 {
		t.Errorf("expected format=%s, got format=%s", dockerformatv1, format)
	}

	format, err = getDockercfgFormat(v2cfgJSON)
	if err != nil {
		t.Errorf("unexpected error when getting docker format for dockercfg v2: %v", err)
	}
	if format != dockerformatv2 {
		t.Errorf("expected format=%s, got format=%s", dockerformatv2, format)
	}

	format, err = getDockercfgFormat([]byte(`{"invalid": "dockercfg"}`))
	if err != errInvalidDockercfg {
		t.Errorf("expected err=%v, got err=%v, format=%v", errInvalidDockercfg, err, format)
	}
}
