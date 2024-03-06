/*
Copyright 2022 The Kubernetes Authors.

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
	"github.com/IBM/go-sdk-core/v5/core"
	tgapiv1 "github.com/IBM/networking-go-sdk/transitgatewayapisv1"
)

// TransitGateway interface defines a method that a IBMCLOUD service object should implement in order to
// use the transitgateway package for listing resource instances.
type TransitGateway interface {
	GetTransitGateway(*tgapiv1.GetTransitGatewayOptions) (*tgapiv1.TransitGateway, *core.DetailedResponse, error)
	GetTransitGatewayByName(name string) (*tgapiv1.TransitGateway, error)
	ListTransitGatewayConnections(*tgapiv1.ListTransitGatewayConnectionsOptions) (*tgapiv1.TransitGatewayConnectionCollection, *core.DetailedResponse, error)
	CreateTransitGateway(*tgapiv1.CreateTransitGatewayOptions) (*tgapiv1.TransitGateway, *core.DetailedResponse, error)
	CreateTransitGatewayConnection(*tgapiv1.CreateTransitGatewayConnectionOptions) (*tgapiv1.TransitGatewayConnectionCust, *core.DetailedResponse, error)
	GetTransitGatewayConnection(*tgapiv1.GetTransitGatewayConnectionOptions) (*tgapiv1.TransitGatewayConnectionCust, *core.DetailedResponse, error)
	DeleteTransitGateway(deleteTransitGatewayOptions *tgapiv1.DeleteTransitGatewayOptions) (response *core.DetailedResponse, err error)
	DeleteTransitGatewayConnection(deleteTransitGatewayConnectionOptions *tgapiv1.DeleteTransitGatewayConnectionOptions) (response *core.DetailedResponse, err error)
}
