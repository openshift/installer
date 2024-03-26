/*
Copyright 2023 The Kubernetes Authors.

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

package transitgateway

import (
	"fmt"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	tgapiv1 "github.com/IBM/networking-go-sdk/transitgatewayapisv1"

	"k8s.io/utils/ptr"

	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/authenticator"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/utils"
)

var currentDate = fmt.Sprintf("%d-%02d-%02d", time.Now().Year(), time.Now().Month(), time.Now().Day())

// Service holds the IBM Cloud Resource Controller Service specific information.
type Service struct {
	tgClient *tgapiv1.TransitGatewayApisV1
}

// NewService returns a new service for the IBM Cloud Transit Gateway api client.
func NewService(options *tgapiv1.TransitGatewayApisV1Options) (TransitGateway, error) {
	if options == nil {
		options = &tgapiv1.TransitGatewayApisV1Options{}
	}
	if options.Authenticator == nil {
		auth, err := authenticator.GetAuthenticator()
		if err != nil {
			return nil, err
		}
		options.Authenticator = auth
	}
	options.Version = ptr.To(currentDate)
	tgClient, err := tgapiv1.NewTransitGatewayApisV1(options)
	if err != nil {
		return nil, err
	}

	return &Service{
		tgClient: tgClient,
	}, nil
}

// GetTransitGateway returns the specified transit gateway. If not found, returns error.
func (s *Service) GetTransitGateway(options *tgapiv1.GetTransitGatewayOptions) (*tgapiv1.TransitGateway, *core.DetailedResponse, error) {
	return s.tgClient.GetTransitGateway(options)
}

// GetTransitGatewayByName returns tranit gateway with given name. If not found, returns nil.
func (s *Service) GetTransitGatewayByName(name string) (*tgapiv1.TransitGateway, error) {
	var transitGateway tgapiv1.TransitGateway

	f := func(start string) (bool, string, error) {
		var listKeyOpt tgapiv1.ListTransitGatewaysOptions

		if start != "" {
			listKeyOpt.Start = &start
		}

		tgList, _, err := s.tgClient.ListTransitGateways(&listKeyOpt)
		if err != nil {
			return false, "", fmt.Errorf("failed to list transit gateway %w", err)
		}

		for _, tg := range tgList.TransitGateways {
			if tg.Name != nil && *tg.Name == name {
				transitGateway = tg
				return true, "", nil
			}
		}

		if tgList.Next != nil && *tgList.Next.Href != "" {
			return false, *tgList.Next.Href, nil
		}

		return true, "", nil
	}

	if err := utils.PagingHelper(f); err != nil {
		return nil, err
	}
	return &transitGateway, nil
}

// ListTransitGatewayConnections lists the transit gateway connections.
func (s *Service) ListTransitGatewayConnections(options *tgapiv1.ListTransitGatewayConnectionsOptions) (*tgapiv1.TransitGatewayConnectionCollection, *core.DetailedResponse, error) {
	return s.tgClient.ListTransitGatewayConnections(options)
}

// CreateTransitGateway creates a transit gateway.
func (s *Service) CreateTransitGateway(options *tgapiv1.CreateTransitGatewayOptions) (*tgapiv1.TransitGateway, *core.DetailedResponse, error) {
	return s.tgClient.CreateTransitGateway(options)
}

// CreateTransitGatewayConnection creates a transit gateway connection.
func (s *Service) CreateTransitGatewayConnection(options *tgapiv1.CreateTransitGatewayConnectionOptions) (*tgapiv1.TransitGatewayConnectionCust, *core.DetailedResponse, error) {
	return s.tgClient.CreateTransitGatewayConnection(options)
}

// GetTransitGatewayConnection returns a transit gateway connection.
func (s *Service) GetTransitGatewayConnection(options *tgapiv1.GetTransitGatewayConnectionOptions) (*tgapiv1.TransitGatewayConnectionCust, *core.DetailedResponse, error) {
	return s.tgClient.GetTransitGatewayConnection(options)
}

// DeleteTransitGateway deletes a transit gateway.
func (s *Service) DeleteTransitGateway(options *tgapiv1.DeleteTransitGatewayOptions) (*core.DetailedResponse, error) {
	return s.tgClient.DeleteTransitGateway(options)
}

// DeleteTransitGatewayConnection deletes a transit gateway connection.
func (s *Service) DeleteTransitGatewayConnection(options *tgapiv1.DeleteTransitGatewayConnectionOptions) (*core.DetailedResponse, error) {
	return s.tgClient.DeleteTransitGatewayConnection(options)
}
