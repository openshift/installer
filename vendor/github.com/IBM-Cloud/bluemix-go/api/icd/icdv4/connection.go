package icdv4

import (
	"fmt"
	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/utils"
)

type ConnectionReq struct {
	Password        string `json:"password,omitempty"`
	CertificateRoot string `json:"certificate_root,omitempty"`
}

type ConnectionRes struct {
	Connection Connection `json:"connection"`
}

type Connection struct {
	Rediss   Uri           `json:"rediss"`
	Grpc     Uri           `json:"grpc"`
	Postgres Uri           `json:"postgres"`
	Https    Uri           `json:"https"`
	Amqps    Uri           `json:"amqps"`
	Cli      CliConn       `json:"cli"`
	Mongo    Uri           `json:"mongodb"`
	Secure   CassandraUri  `json:"secure"`
}

type CassandraUri struct {
	Hosts []struct {
		HostName string `json:"hostname"`
		Port     int    `json:"port"`
	} `json:"hosts"`
	Authentication struct {
		Method   string `json:"method"`
		UserName string `json:"username"`
		Password string `json:"password"`
	}
	Bundle struct {
		Name         string `json:"name"`
		BundleBase64 string `json:"bundle_base64"`
	} `json:"bundle"`
}

type Uri struct {
	Type     string   `json:"type"`
	Composed []string `json:"composed"`
	Scheme   string   `json:"scheme"`
	Hosts    []struct {
		HostName string `json:"hostname"`
		Port     int    `json:"port"`
	} `json:"hosts"`
	Path           string      `json:"path"`
	QueryOptions   interface{} `json:"query_options"`
	Authentication struct {
		Method   string `json:"method"`
		UserName string `json:"username"`
		Password string `json:"password"`
	}
	Certificate struct {
		Name              string `json:"name"`
		CertificateBase64 string `json:"certificate_base64"`
	} `json:"certificate"`
	Database interface{} `json:"database"`
}

type CliConn struct {
	Type        string      `json:"type"`
	Composed    []string    `json:"composed"`
	Environment interface{} `json:"environment"`
	Bin         string      `json:"bin"`
	Arguments   [][]string  `json:"arguments"`
	Certificate struct {
		Name              string `json:"name"`
		CertificateBase64 string `json:"certificate_base64"`
	} `json:"certificate"`
}

type Connections interface {
	GetConnection(icdId string, userId string, endpoint ...string) (Connection, error)
	GetConnectionSubstitution(icdId string, userID string, connectionReq ConnectionReq) (Connection, error)
}

type connections struct {
	client *client.Client
}

func newConnectionAPI(c *client.Client) Connections {
	return &connections{
		client: c,
	}
}

func (r *connections) GetConnection(icdId string, userId string, endpoint ...string) (Connection, error) {
	connectionRes := ConnectionRes{}
	connectionEndpoint := "public"
	if len(endpoint) > 0 {
		connectionEndpoint = endpoint[0]
	}
	rawURL := fmt.Sprintf("/v4/ibm/deployments/%s/users/%s/connections/%s", utils.EscapeUrlParm(icdId), userId, connectionEndpoint)
	_, err := r.client.Get(rawURL, &connectionRes)
	if err != nil {
		return connectionRes.Connection, err
	}
	return connectionRes.Connection, nil
}

func (r *connections) GetConnectionSubstitution(icdId string, userID string, connectionReq ConnectionReq) (Connection, error) {
	connectionResSub := ConnectionRes{}
	rawURL := fmt.Sprintf("/v4/ibm/deployments/%s/connections", utils.EscapeUrlParm(icdId))
	_, err := r.client.Post(rawURL, &connectionReq, &connectionResSub)
	if err != nil {
		return connectionResSub.Connection, err
	}
	return connectionResSub.Connection, nil
}
