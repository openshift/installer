/*
Copyright 2021 The Kubernetes Authors.

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

package converters

import (
	"encoding/base64"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/pkg/errors"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
)

// PollerToFuture converts an SDK poller to an infrav1.Future.
func PollerToFuture[T any](poller *runtime.Poller[T], futureType, service, resourceName, rgName string) (*infrav1.Future, error) {
	token, err := poller.ResumeToken()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get resume token")
	}
	return &infrav1.Future{
		Type:          futureType,
		ResourceGroup: rgName,
		ServiceName:   service,
		Name:          resourceName,
		Data:          base64.URLEncoding.EncodeToString([]byte(token)),
	}, nil
}

// FutureToResumeToken converts an infrav1.Future to an Azure SDK resume token.
func FutureToResumeToken(future infrav1.Future) (string, error) {
	if future.Data == "" {
		return "", errors.New("failed to unmarshal future data: data is empty")
	}
	token, err := base64.URLEncoding.DecodeString(future.Data)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode future data")
	}
	return string(token), nil
}
