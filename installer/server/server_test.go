package server

import (
	"bytes"
	"encoding/json"
	"io"
)

func validTectonicMetalCreate() (*CreateOperation, error) {
	cluster := validTectonicMetalCluster()
	b, err := json.Marshal(cluster)
	if err != nil {
		return nil, err
	}
	return &CreateOperation{
		ClusterKind: "tectonic-metal",
		ClusterData: json.RawMessage(b),
	}, nil
}

func validTectonicMetalCreateJSON() (io.Reader, error) {
	data, err := validTectonicMetalCreate()
	if err != nil {
		return nil, err
	}
	buf := &bytes.Buffer{}
	err = json.NewEncoder(buf).Encode(data)
	return buf, err
}
