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

package mock

import (
	// Runtime dependency of mockgen, required when using vendoring so go mod knows
	// to pull it in.
	_ "github.com/golang/mock/mockgen/model"
)

//go:generate mockgen -package mock -destination=compute.go sigs.k8s.io/cluster-api-provider-openstack/pkg/clients ComputeClient
//go:generate /usr/bin/env bash -c "cat ../../../hack/boilerplate/boilerplate.generatego.txt compute.go > _compute.go && mv _compute.go compute.go"

//go:generate mockgen -package mock -destination=image.go sigs.k8s.io/cluster-api-provider-openstack/pkg/clients ImageClient
//go:generate /usr/bin/env bash -c "cat ../../../hack/boilerplate/boilerplate.generatego.txt image.go > _image.go && mv _image.go image.go"

//go:generate mockgen -package mock -destination=loadbalancer.go sigs.k8s.io/cluster-api-provider-openstack/pkg/clients LbClient
//go:generate /usr/bin/env bash -c "cat ../../../hack/boilerplate/boilerplate.generatego.txt loadbalancer.go > _loadbalancer.go && mv _loadbalancer.go loadbalancer.go"

//go:generate mockgen -package mock -destination=network.go sigs.k8s.io/cluster-api-provider-openstack/pkg/clients NetworkClient
//go:generate /usr/bin/env bash -c "cat ../../../hack/boilerplate/boilerplate.generatego.txt network.go > _network.go && mv _network.go network.go"

//go:generate mockgen -package mock -destination=volume.go sigs.k8s.io/cluster-api-provider-openstack/pkg/clients VolumeClient
//go:generate /usr/bin/env bash -c "cat ../../../hack/boilerplate/boilerplate.generatego.txt volume.go > _volume.go && mv _volume.go volume.go"
