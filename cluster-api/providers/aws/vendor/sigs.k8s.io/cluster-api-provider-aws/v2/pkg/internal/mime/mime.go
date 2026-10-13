/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package mime provides a function to generate a multipart MIME document.
package mime

import (
	"bytes"
	"fmt"
	"html/template"
	"mime/multipart"
	"net/textproto"
	"strings"
)

const (
	includePart = "file:///etc/secret-userdata.txt\n"
)

var (
	includeType = textproto.MIMEHeader{
		"content-type": {"text/x-include-url"},
	}

	boothookType = textproto.MIMEHeader{
		"content-type": {"text/cloud-boothook"},
	}

	multipartHeader = strings.Join([]string{
		"MIME-Version: 1.0",
		"Content-Type: multipart/mixed; boundary=\"%s\"",
		"\n",
	}, "\n")
)

type scriptVariables struct {
	SecretPrefix string
	Chunks       int32
	Region       string
	Endpoint     string
}

// GenerateInitDocument renders a given template, applies MIME properties
// and returns a series of byte chunks which put together represent a UserData
// script.
func GenerateInitDocument(secretPrefix string, chunks int32, region string, secretFetchScript string) ([]byte, error) {
	var secretFetchTemplate = template.Must(template.New("secret-fetch-script").Parse(secretFetchScript))

	var buf bytes.Buffer
	mpWriter := multipart.NewWriter(&buf)
	buf.WriteString(fmt.Sprintf(multipartHeader, mpWriter.Boundary()))
	scriptWriter, err := mpWriter.CreatePart(boothookType)
	if err != nil {
		return []byte{}, err
	}

	scriptVariables := scriptVariables{
		SecretPrefix: secretPrefix,
		Chunks:       chunks,
		Region:       region,
	}

	var scriptBuf bytes.Buffer
	if err := secretFetchTemplate.Execute(&scriptBuf, scriptVariables); err != nil {
		return []byte{}, err
	}
	_, err = scriptWriter.Write(scriptBuf.Bytes())
	if err != nil {
		return []byte{}, err
	}

	includeWriter, err := mpWriter.CreatePart(includeType)
	if err != nil {
		return []byte{}, err
	}

	_, err = includeWriter.Write([]byte(includePart))
	if err != nil {
		return []byte{}, err
	}

	if err := mpWriter.Close(); err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}
