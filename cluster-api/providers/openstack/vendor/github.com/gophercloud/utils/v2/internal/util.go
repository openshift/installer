package internal

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"maps"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"strings"
)

// RemainingKeys will inspect a struct and compare it to a map. Any struct
// field that does not have a JSON tag that matches a key in the map or
// a matching lower-case field in the map will be returned as an extra.
//
// This is useful for determining the extra fields returned in response bodies
// for resources that can contain an arbitrary or dynamic number of fields.
func RemainingKeys(s any, m map[string]any) (extras map[string]any) {
	extras = make(map[string]any)
	maps.Copy(extras, m)

	valueOf := reflect.ValueOf(s)
	typeOf := reflect.TypeOf(s)
	for i := 0; i < valueOf.NumField(); i++ {
		field := typeOf.Field(i)

		lowerField := strings.ToLower(field.Name)
		delete(extras, lowerField)

		if tagValue := field.Tag.Get("json"); tagValue != "" && tagValue != "-" {
			delete(extras, tagValue)
		}
	}

	return
}

// PrepareTLSConfig generates TLS config based on the specifed parameters
func PrepareTLSConfig(caCertFile, clientCertFile, clientKeyFile string, insecure *bool) (*tls.Config, error) {
	config := &tls.Config{}
	if caCertFile != "" {
		caCert, err := pathOrContents(caCertFile)
		if err != nil {
			return nil, fmt.Errorf("error reading CA Cert: %s", err)
		}

		caCertPool := x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM(bytes.TrimSpace(caCert)); !ok {
			return nil, fmt.Errorf("error parsing CA Cert from %s", caCertFile)
		}
		config.RootCAs = caCertPool
	}

	if insecure == nil {
		config.InsecureSkipVerify = false
	} else {
		config.InsecureSkipVerify = *insecure
	}

	if clientCertFile != "" && clientKeyFile != "" {
		clientCert, err := pathOrContents(clientCertFile)
		if err != nil {
			return nil, fmt.Errorf("error reading Client Cert: %s", err)
		}
		clientKey, err := pathOrContents(clientKeyFile)
		if err != nil {
			return nil, fmt.Errorf("error reading Client Key: %s", err)
		}

		cert, err := tls.X509KeyPair(clientCert, clientKey)
		if err != nil {
			return nil, err
		}

		config.Certificates = []tls.Certificate{cert}
	}

	return config, nil
}

func pathOrContents(poc string) ([]byte, error) {
	if len(poc) == 0 {
		return nil, nil
	}

	path := poc
	if path[0] == '~' {
		usr, err := user.Current()
		if err != nil {
			return []byte(path), err
		}

		if len(path) == 1 {
			path = usr.HomeDir
		} else if strings.HasPrefix(path, "~/") {
			path = filepath.Join(usr.HomeDir, path[2:])
		}
	}

	if _, err := os.Stat(path); err == nil {
		contents, err := os.ReadFile(path)
		if err != nil {
			return contents, err
		}
		return contents, nil
	}

	return []byte(poc), nil
}
