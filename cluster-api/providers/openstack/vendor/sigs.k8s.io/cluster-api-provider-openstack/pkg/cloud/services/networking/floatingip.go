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

package networking

import (
	"fmt"
	"time"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/attributestags"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/utils/ptr"

	"sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha1"
	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/metrics"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/record"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/names"
)

func (s *Service) GetOrCreateFloatingIP(eventObject runtime.Object, openStackCluster *infrav1.OpenStackCluster, clusterResourceName string, ip *string) (*floatingips.FloatingIP, error) {
	var fp *floatingips.FloatingIP
	var err error
	var fpCreateOpts floatingips.CreateOpts

	if ptr.Deref(ip, "") != "" {
		fp, err = s.GetFloatingIP(*ip)
		if err != nil {
			return nil, err
		}
		if fp != nil {
			return fp, nil
		}
		// only admin can add ip address
		fpCreateOpts.FloatingIP = *ip
	}

	fpCreateOpts.FloatingNetworkID = openStackCluster.Status.ExternalNetwork.ID
	fpCreateOpts.Description = names.GetDescription(clusterResourceName)

	s.scope.Logger().Info("Creating floating IP", "ip", fpCreateOpts.FloatingIP, "floatingNetworkID", openStackCluster.Status.ExternalNetwork.ID, "description", fpCreateOpts.Description)

	fp, err = s.client.CreateFloatingIP(fpCreateOpts)
	if err != nil {
		record.Warnf(eventObject, "FailedCreateFloatingIP", "Failed to create floating IP %s: %v", fpCreateOpts.FloatingIP, err)
		return nil, err
	}

	if len(openStackCluster.Spec.Tags) > 0 {
		mc := metrics.NewMetricPrometheusContext("floating_ip", "update")
		_, err = s.client.ReplaceAllAttributesTags("floatingips", fp.ID, attributestags.ReplaceAllOpts{
			Tags: openStackCluster.Spec.Tags,
		})
		if mc.ObserveRequest(err) != nil {
			return nil, err
		}
	}

	record.Eventf(eventObject, "SuccessfulCreateFloatingIP", "Created floating IP %s with id %s", fp.FloatingIP, fp.ID)
	return fp, nil
}

func (s *Service) CreateFloatingIPForPool(pool *v1alpha1.OpenStackFloatingIPPool) (*floatingips.FloatingIP, error) {
	var fpCreateOpts floatingips.CreateOpts

	fpCreateOpts.FloatingNetworkID = pool.Status.FloatingIPNetwork.ID
	fpCreateOpts.Description = fmt.Sprintf("Created by cluster-api-provider-openstack OpenStackFloatingIPPool %s", pool.Name)

	fp, err := s.client.CreateFloatingIP(fpCreateOpts)
	if err != nil {
		record.Warnf(pool, "FailedCreateFloatingIP", "%s failed to create floating IP: %v", pool.Name, err)
		return nil, err
	}

	record.Eventf(pool, "SuccessfulCreateFloatingIP", "%s created floating IP %s with id %s", pool.Name, fp.FloatingIP, fp.ID)
	return fp, nil
}

func (s *Service) TagFloatingIP(ip string, tag string) error {
	fip, err := s.GetFloatingIP(ip)
	if err != nil {
		return err
	}
	if fip == nil {
		return nil
	}

	mc := metrics.NewMetricPrometheusContext("floating_ip", "update")
	_, err = s.client.ReplaceAllAttributesTags("floatingips", fip.ID, attributestags.ReplaceAllOpts{
		Tags: []string{tag},
	})
	if mc.ObserveRequest(err) != nil {
		return err
	}
	return nil
}

func (s *Service) GetFloatingIPsByTag(tag string) ([]floatingips.FloatingIP, error) {
	fipList, err := s.client.ListFloatingIP(floatingips.ListOpts{Tags: tag})
	if err != nil {
		return nil, err
	}
	return fipList, nil
}

func (s *Service) GetFloatingIP(ip string) (*floatingips.FloatingIP, error) {
	fpList, err := s.client.ListFloatingIP(floatingips.ListOpts{FloatingIP: ip})
	if err != nil {
		return nil, err
	}
	if len(fpList) == 0 {
		return nil, nil
	}
	return &fpList[0], nil
}

func (s *Service) GetFloatingIPByPortID(portID string) (*floatingips.FloatingIP, error) {
	fpList, err := s.client.ListFloatingIP(floatingips.ListOpts{PortID: portID})
	if err != nil {
		return nil, err
	}
	if len(fpList) == 0 {
		return nil, nil
	}
	return &fpList[0], nil
}

func (s *Service) DeleteFloatingIP(eventObject runtime.Object, ip string) error {
	fip, err := s.GetFloatingIP(ip)
	if err != nil {
		return err
	}
	if fip == nil {
		// nothing to do
		return nil
	}

	err = s.client.DeleteFloatingIP(fip.ID)
	if err != nil {
		record.Warnf(eventObject, "FailedDeleteFloatingIP", "Failed to delete floating IP %s: %v", ip, err)
		return err
	}

	record.Eventf(eventObject, "SuccessfulDeleteFloatingIP", "Deleted floating IP %s", ip)
	return nil
}

var backoff = wait.Backoff{
	Steps:    10,
	Duration: 1 * time.Second,
	Factor:   2.0,
	Jitter:   0.1,
	Cap:      30 * time.Second,
}

func (s *Service) AssociateFloatingIP(eventObject runtime.Object, fp *floatingips.FloatingIP, portID string) error {
	s.scope.Logger().Info("Associating floating IP", "ID", fp.ID, "IP", fp.FloatingIP)

	if fp.PortID == portID {
		s.scope.Logger().Info("Floating IP already associated:", "ID", fp.ID, "IP", fp.FloatingIP)
		return nil
	}

	fpUpdateOpts := &floatingips.UpdateOpts{
		PortID: &portID,
	}

	_, err := s.client.UpdateFloatingIP(fp.ID, fpUpdateOpts)
	if err != nil {
		record.Warnf(eventObject, "FailedAssociateFloatingIP", "Failed to associate floating IP %s with port %s: %v", fp.FloatingIP, portID, err)
		return err
	}

	if err = s.waitForFloatingIP(fp.ID, "ACTIVE"); err != nil {
		record.Warnf(eventObject, "FailedAssociateFloatingIP", "Failed to associate floating IP %s with port %s: wait for floating IP ACTIVE: %v", fp.FloatingIP, portID, err)
		return err
	}

	record.Eventf(eventObject, "SuccessfulAssociateFloatingIP", "Associated floating IP %s with port %s", fp.FloatingIP, portID)
	return nil
}

func (s *Service) DisassociateFloatingIP(eventObject runtime.Object, ip string) error {
	fip, err := s.GetFloatingIP(ip)
	if err != nil {
		return err
	}
	if fip == nil || fip.FloatingIP == "" {
		s.scope.Logger().Info("Floating IP not associated", "IP", ip)
		return nil
	}

	s.scope.Logger().Info("Disassociating floating IP", "ID", fip.ID, "IP", fip.FloatingIP)

	fpUpdateOpts := &floatingips.UpdateOpts{
		PortID: nil,
	}

	_, err = s.client.UpdateFloatingIP(fip.ID, fpUpdateOpts)
	if err != nil {
		record.Warnf(eventObject, "FailedDisassociateFloatingIP", "Failed to disassociate floating IP %s: %v", fip.FloatingIP, err)
		return err
	}

	if err = s.waitForFloatingIP(fip.ID, "DOWN"); err != nil {
		record.Warnf(eventObject, "FailedDisassociateFloatingIP", "Failed to disassociate floating IP: wait for floating IP DOWN: %v", fip.FloatingIP, err)
		return err
	}

	record.Eventf(eventObject, "SuccessfulDisassociateFloatingIP", "Disassociated floating IP %s", fip.FloatingIP)
	return nil
}

func (s *Service) waitForFloatingIP(id, target string) error {
	s.scope.Logger().Info("Waiting for floating IP", "ID", id, "status", target)
	return wait.ExponentialBackoff(backoff, func() (bool, error) {
		fip, err := s.client.GetFloatingIP(id)
		if err != nil {
			return false, err
		}
		return fip.Status == target, nil
	})
}
