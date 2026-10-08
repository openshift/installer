/*
Copyright 2026 The Kubernetes Authors.

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

package v1alpha2

import (
	"context"

	vmoprv1alpha2 "github.com/vmware-tanzu/vm-operator/api/v1alpha2"
	vmoprv1alpha2common "github.com/vmware-tanzu/vm-operator/api/v1alpha2/common"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/conversion"
	vmoprvhub "sigs.k8s.io/cluster-api-provider-vsphere/pkg/conversion/api/vmoperator/hub"
)

func convert_v1alpha2_VirtualMachine_To_hub_VirtualMachine(_ context.Context, src *vmoprv1alpha2.VirtualMachine, dst *vmoprvhub.VirtualMachine) error {
	dst.ObjectMeta = src.ObjectMeta

	if src.Spec.Affinity != nil {
		dst.Spec.Affinity = &vmoprvhub.AffinitySpec{}
		if src.Spec.Affinity.VMAffinity != nil {
			dst.Spec.Affinity.VMAffinity = &vmoprvhub.VMAffinitySpec{}
			if src.Spec.Affinity.VMAffinity.RequiredDuringSchedulingPreferredDuringExecution != nil {
				terms := []vmoprvhub.VMAffinityTerm{}
				for _, term := range src.Spec.Affinity.VMAffinity.RequiredDuringSchedulingPreferredDuringExecution {
					t := vmoprvhub.VMAffinityTerm{
						TopologyKey: term.TopologyKey,
					}
					if term.LabelSelector != nil {
						t.LabelSelector = term.LabelSelector.DeepCopy()
					}
					terms = append(terms, t)
				}
				dst.Spec.Affinity.VMAffinity.RequiredDuringSchedulingPreferredDuringExecution = terms
			}
		}
		if src.Spec.Affinity.VMAntiAffinity != nil {
			dst.Spec.Affinity.VMAntiAffinity = &vmoprvhub.VMAntiAffinitySpec{}
			if src.Spec.Affinity.VMAntiAffinity.PreferredDuringSchedulingPreferredDuringExecution != nil {
				terms := []vmoprvhub.VMAffinityTerm{}
				for _, term := range src.Spec.Affinity.VMAntiAffinity.PreferredDuringSchedulingPreferredDuringExecution {
					t := vmoprvhub.VMAffinityTerm{
						TopologyKey: term.TopologyKey,
					}
					if term.LabelSelector != nil {
						t.LabelSelector = term.LabelSelector.DeepCopy()
					}
					terms = append(terms, t)
				}
				dst.Spec.Affinity.VMAntiAffinity.PreferredDuringSchedulingPreferredDuringExecution = terms
			}
		}
	}
	if src.Spec.Bootstrap != nil {
		dst.Spec.Bootstrap = &vmoprvhub.VirtualMachineBootstrapSpec{}
		if src.Spec.Bootstrap.CloudInit != nil {
			dst.Spec.Bootstrap.CloudInit = &vmoprvhub.VirtualMachineBootstrapCloudInitSpec{}
			if src.Spec.Bootstrap.CloudInit.RawCloudConfig != nil {
				dst.Spec.Bootstrap.CloudInit.RawCloudConfig = &vmoprvhub.SecretKeySelector{
					Name: src.Spec.Bootstrap.CloudInit.RawCloudConfig.Name,
					Key:  src.Spec.Bootstrap.CloudInit.RawCloudConfig.Key,
				}
			}
		}
	}
	dst.Spec.ClassName = src.Spec.ClassName
	dst.Spec.GroupName = src.Spec.GroupName
	dst.Spec.ImageName = src.Spec.ImageName
	dst.Spec.MinHardwareVersion = src.Spec.MinHardwareVersion
	if src.Spec.Network != nil {
		dst.Spec.Network = &vmoprvhub.VirtualMachineNetworkSpec{}
		if src.Spec.Network.Interfaces != nil {
			dst.Spec.Network.Interfaces = []vmoprvhub.VirtualMachineNetworkInterfaceSpec{}
			for _, iface := range src.Spec.Network.Interfaces {
				d := vmoprvhub.VirtualMachineNetworkInterfaceSpec{}
				d.Addresses = iface.Addresses
				// AdvancedProperties existing in hub but not in v1alpha2.VirtualMachineNetworkInterfaceSpec
				d.DHCP4 = iface.DHCP4
				d.DHCP6 = iface.DHCP6
				d.Gateway4 = iface.Gateway4
				d.Gateway6 = iface.Gateway6
				d.GuestDeviceName = iface.GuestDeviceName
				// IPAMModes existing in hub but not in v1alpha2.VirtualMachineNetworkInterfaceSpec
				d.MACAddr = iface.MACAddr
				if iface.MTU != nil {
					d.MTU = ptr.To(*iface.MTU)
				}
				d.Name = iface.Name
				d.Nameservers = iface.Nameservers
				if iface.Network != nil {
					d.Network = &vmoprvhub.PartialObjectRef{
						TypeMeta: metav1.TypeMeta{
							Kind:       iface.Network.Kind,
							APIVersion: iface.Network.APIVersion,
						},
						Name: iface.Network.Name,
					}
				}
				if iface.Routes != nil {
					d.Routes = []vmoprvhub.VirtualMachineNetworkRouteSpec{}
					for _, route := range iface.Routes {
						d.Routes = append(d.Routes, vmoprvhub.VirtualMachineNetworkRouteSpec{
							To:     route.To,
							Via:    route.Via,
							Metric: route.Metric,
						})
					}
				}
				d.SearchDomains = iface.SearchDomains
				// Type existing in hub but not in v1alpha2.VirtualMachineNetworkInterfaceSpec
				// VMXNet3 existing in hub but not in v1alpha2.VirtualMachineNetworkInterfaceSpec
				// VNUMANodeID existing in hub but not in v1alpha2.VirtualMachineNetworkInterfaceSpec
				dst.Spec.Network.Interfaces = append(dst.Spec.Network.Interfaces, d)
			}
		}
		// VLANs existing in hub but not in v1alpha2.VirtualMachineNetworkSpec
	}
	// Policies existing in hub but not in v1alpha2.VirtualMachineSpec
	dst.Spec.PowerOffMode = vmoprvhub.VirtualMachinePowerOpMode(src.Spec.PowerOffMode)
	dst.Spec.PowerState = vmoprvhub.VirtualMachinePowerState(src.Spec.PowerState)
	if src.Spec.ReadinessProbe != nil {
		dst.Spec.ReadinessProbe = &vmoprvhub.VirtualMachineReadinessProbeSpec{}
		if src.Spec.ReadinessProbe.TCPSocket != nil {
			dst.Spec.ReadinessProbe.TCPSocket = &vmoprvhub.TCPSocketAction{
				Port: src.Spec.ReadinessProbe.TCPSocket.Port,
				Host: src.Spec.ReadinessProbe.TCPSocket.Host,
			}
		}
	}
	if src.Spec.Reserved != nil {
		dst.Spec.Reserved = &vmoprvhub.VirtualMachineReservedSpec{
			ResourcePolicyName: src.Spec.Reserved.ResourcePolicyName,
		}
	}
	dst.Spec.StorageClass = src.Spec.StorageClass
	if src.Spec.Volumes != nil {
		dst.Spec.Volumes = []vmoprvhub.VirtualMachineVolume{}
		for _, volume := range src.Spec.Volumes {
			v := vmoprvhub.VirtualMachineVolume{}
			// ApplicationType existing in hub but not in v1alpha2.VirtualMachineVolume
			// ControllerBusNumber existing in hub but not in v1alpha2.VirtualMachineVolume
			// ControllerType existing in hub but not in v1alpha2.VirtualMachineVolume
			// DiskMode existing in hub but not in v1alpha2.VirtualMachineVolume
			v.Name = volume.Name
			if volume.PersistentVolumeClaim != nil {
				v.PersistentVolumeClaim = &vmoprvhub.PersistentVolumeClaimVolumeSource{
					PersistentVolumeClaimVolumeSource: corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: volume.PersistentVolumeClaim.ClaimName,
						ReadOnly:  volume.PersistentVolumeClaim.ReadOnly,
					},
				}
				if volume.PersistentVolumeClaim.InstanceVolumeClaim != nil {
					v.PersistentVolumeClaim.InstanceVolumeClaim = &vmoprvhub.InstanceVolumeClaimVolumeSource{
						StorageClass: volume.PersistentVolumeClaim.InstanceVolumeClaim.StorageClass,
						Size:         volume.PersistentVolumeClaim.InstanceVolumeClaim.Size,
					}
				}

			}
			// Removable existing in hub but not in v1alpha2.VirtualMachineVolume
			// SharingMode existing in hub but not in v1alpha2.VirtualMachineVolume
			// UnitNumber existing in hub but not in v1alpha2.VirtualMachineVolume
			dst.Spec.Volumes = append(dst.Spec.Volumes, v)
		}
	}

	dst.Status.BiosUUID = src.Status.BiosUUID
	if src.Status.Conditions != nil {
		dst.Status.Conditions = []metav1.Condition{}
		for _, condition := range src.Status.Conditions {
			dst.Status.Conditions = append(dst.Status.Conditions, condition)
		}
	}
	if src.Status.Network != nil {
		dst.Status.Network = &vmoprvhub.VirtualMachineNetworkStatus{}
		dst.Status.Network.HostName = src.Status.Network.HostName
		if src.Status.Network.Interfaces != nil {
			dst.Status.Network.Interfaces = []vmoprvhub.VirtualMachineNetworkInterfaceStatus{}
			for _, iface := range src.Status.Network.Interfaces {
				d := vmoprvhub.VirtualMachineNetworkInterfaceStatus{}
				d.DeviceKey = iface.DeviceKey
				if iface.DNS != nil {
					d.DNS = &vmoprvhub.VirtualMachineNetworkDNSStatus{
						DHCP:          iface.DNS.DHCP,
						HostName:      iface.DNS.HostName,
						DomainName:    iface.DNS.DomainName,
						Nameservers:   iface.DNS.Nameservers,
						SearchDomains: iface.DNS.SearchDomains,
					}
				}
				if iface.IP != nil {
					d.IP = &vmoprvhub.VirtualMachineNetworkInterfaceIPStatus{
						MACAddr: iface.IP.MACAddr,
					}
					if iface.IP.Addresses != nil {
						d.IP.Addresses = []vmoprvhub.VirtualMachineNetworkInterfaceIPAddrStatus{}
						for _, addr := range iface.IP.Addresses {
							d.IP.Addresses = append(d.IP.Addresses, vmoprvhub.VirtualMachineNetworkInterfaceIPAddrStatus{
								Address:  addr.Address,
								Lifetime: addr.Lifetime,
								Origin:   addr.Origin,
								State:    addr.State,
							})
						}
					}
					if iface.IP.AutoConfigurationEnabled != nil {
						d.IP.AutoConfigurationEnabled = ptr.To(*iface.IP.AutoConfigurationEnabled)
					}
					if iface.IP.DHCP != nil {
						d.IP.DHCP = &vmoprvhub.VirtualMachineNetworkDHCPStatus{
							IP4: vmoprvhub.VirtualMachineNetworkDHCPOptionsStatus{
								Enabled: iface.IP.DHCP.IP4.Enabled,
							},
							IP6: vmoprvhub.VirtualMachineNetworkDHCPOptionsStatus{
								Enabled: iface.IP.DHCP.IP6.Enabled,
							},
						}
						if iface.IP.DHCP.IP4.Config != nil {
							d.IP.DHCP.IP4.Config = []vmoprvhub.KeyValuePair{}
							for _, pair := range iface.IP.DHCP.IP4.Config {
								d.IP.DHCP.IP4.Config = append(d.IP.DHCP.IP4.Config, vmoprvhub.KeyValuePair{
									Key:   pair.Key,
									Value: pair.Value,
								})
							}
						}
						if iface.IP.DHCP.IP6.Config != nil {
							d.IP.DHCP.IP6.Config = []vmoprvhub.KeyValuePair{}
							for _, pair := range iface.IP.DHCP.IP6.Config {
								d.IP.DHCP.IP6.Config = append(d.IP.DHCP.IP6.Config, vmoprvhub.KeyValuePair{
									Key:   pair.Key,
									Value: pair.Value,
								})
							}
						}
					}
				}
				d.Name = iface.Name
				dst.Status.Network.Interfaces = append(dst.Status.Network.Interfaces, d)
			}
		}
		dst.Status.Network.PrimaryIP4 = src.Status.Network.PrimaryIP4
		dst.Status.Network.PrimaryIP6 = src.Status.Network.PrimaryIP6
	}
	dst.Status.NodeName = src.Status.Host // Field renamed in hub
	// Policies existing in hub but not in v1alpha2.VirtualMachineStatus
	dst.Status.PowerState = vmoprvhub.VirtualMachinePowerState(src.Status.PowerState)
	dst.Status.Zone = src.Status.Zone

	return nil
}

func convert_hub_VirtualMachine_To_v1alpha2_VirtualMachine(_ context.Context, src *vmoprvhub.VirtualMachine, dst *vmoprv1alpha2.VirtualMachine) error {
	dst.ObjectMeta = src.ObjectMeta

	if src.Spec.Affinity != nil {
		dst.Spec.Affinity = &vmoprv1alpha2.AffinitySpec{}
		if src.Spec.Affinity.VMAffinity != nil {
			dst.Spec.Affinity.VMAffinity = &vmoprv1alpha2.VMAffinitySpec{}
			if src.Spec.Affinity.VMAffinity.RequiredDuringSchedulingPreferredDuringExecution != nil {
				terms := []vmoprv1alpha2.VMAffinityTerm{}
				for _, term := range src.Spec.Affinity.VMAffinity.RequiredDuringSchedulingPreferredDuringExecution {
					t := vmoprv1alpha2.VMAffinityTerm{
						TopologyKey: term.TopologyKey,
					}
					if term.LabelSelector != nil {
						t.LabelSelector = term.LabelSelector.DeepCopy()
					}
					terms = append(terms, t)
				}
				dst.Spec.Affinity.VMAffinity.RequiredDuringSchedulingPreferredDuringExecution = terms
			}
		}
		if src.Spec.Affinity.VMAntiAffinity != nil {
			dst.Spec.Affinity.VMAntiAffinity = &vmoprv1alpha2.VMAntiAffinitySpec{}
			if src.Spec.Affinity.VMAntiAffinity.PreferredDuringSchedulingPreferredDuringExecution != nil {
				terms := []vmoprv1alpha2.VMAffinityTerm{}
				for _, term := range src.Spec.Affinity.VMAntiAffinity.PreferredDuringSchedulingPreferredDuringExecution {
					t := vmoprv1alpha2.VMAffinityTerm{
						TopologyKey: term.TopologyKey,
					}
					if term.LabelSelector != nil {
						t.LabelSelector = term.LabelSelector.DeepCopy()
					}
					terms = append(terms, t)
				}
				dst.Spec.Affinity.VMAntiAffinity.PreferredDuringSchedulingPreferredDuringExecution = terms
			}
		}
	}
	if src.Spec.Bootstrap != nil {
		dst.Spec.Bootstrap = &vmoprv1alpha2.VirtualMachineBootstrapSpec{}
		if src.Spec.Bootstrap.CloudInit != nil {
			dst.Spec.Bootstrap.CloudInit = &vmoprv1alpha2.VirtualMachineBootstrapCloudInitSpec{}
			if src.Spec.Bootstrap.CloudInit.RawCloudConfig != nil {
				dst.Spec.Bootstrap.CloudInit.RawCloudConfig = &vmoprv1alpha2common.SecretKeySelector{
					Name: src.Spec.Bootstrap.CloudInit.RawCloudConfig.Name,
					Key:  src.Spec.Bootstrap.CloudInit.RawCloudConfig.Key,
				}
			}
		}
	}
	dst.Spec.ClassName = src.Spec.ClassName
	dst.Spec.GroupName = src.Spec.GroupName
	dst.Spec.ImageName = src.Spec.ImageName
	dst.Spec.MinHardwareVersion = src.Spec.MinHardwareVersion
	if src.Spec.Network != nil {
		dst.Spec.Network = &vmoprv1alpha2.VirtualMachineNetworkSpec{}
		if src.Spec.Network.Interfaces != nil {
			dst.Spec.Network.Interfaces = []vmoprv1alpha2.VirtualMachineNetworkInterfaceSpec{}
			for _, iface := range src.Spec.Network.Interfaces {
				d := vmoprv1alpha2.VirtualMachineNetworkInterfaceSpec{}
				d.Addresses = iface.Addresses
				// AdvancedProperties existing in hub but not in v1alpha2.VirtualMachineNetworkInterfaceSpec
				d.DHCP4 = iface.DHCP4
				d.DHCP6 = iface.DHCP6
				d.Gateway4 = iface.Gateway4
				d.Gateway6 = iface.Gateway6
				d.GuestDeviceName = iface.GuestDeviceName
				// IPAMModes existing in hub but not in v1alpha2.VirtualMachineNetworkInterfaceSpec
				d.MACAddr = iface.MACAddr
				if iface.MTU != nil {
					d.MTU = ptr.To(*iface.MTU)
				}
				d.Name = iface.Name
				d.Nameservers = iface.Nameservers
				if iface.Network != nil {
					d.Network = &vmoprv1alpha2common.PartialObjectRef{
						TypeMeta: metav1.TypeMeta{
							Kind:       iface.Network.Kind,
							APIVersion: iface.Network.APIVersion,
						},
						Name: iface.Network.Name,
					}
				}
				if iface.Routes != nil {
					d.Routes = []vmoprv1alpha2.VirtualMachineNetworkRouteSpec{}
					for _, route := range iface.Routes {
						d.Routes = append(d.Routes, vmoprv1alpha2.VirtualMachineNetworkRouteSpec{
							To:     route.To,
							Via:    route.Via,
							Metric: route.Metric,
						})
					}
				}
				d.SearchDomains = iface.SearchDomains
				// Type existing in hub but not in v1alpha2.VirtualMachineNetworkInterfaceSpec
				// VMXNet3 existing in hub but not in v1alpha2.VirtualMachineNetworkInterfaceSpec
				// VNUMANodeID existing in hub but not in v1alpha2.VirtualMachineNetworkInterfaceSpec
				dst.Spec.Network.Interfaces = append(dst.Spec.Network.Interfaces, d)
			}
		}
		// VLANs existing in hub but not in v1alpha2.VirtualMachineNetworkSpec
	}
	// Policies existing in hub but not in v1alpha2.VirtualMachineSpec
	dst.Spec.PowerOffMode = vmoprv1alpha2.VirtualMachinePowerOpMode(src.Spec.PowerOffMode)
	dst.Spec.PowerState = vmoprv1alpha2.VirtualMachinePowerState(src.Spec.PowerState)
	if src.Spec.ReadinessProbe != nil {
		dst.Spec.ReadinessProbe = &vmoprv1alpha2.VirtualMachineReadinessProbeSpec{}
		if src.Spec.ReadinessProbe.TCPSocket != nil {
			dst.Spec.ReadinessProbe.TCPSocket = &vmoprv1alpha2.TCPSocketAction{
				Port: src.Spec.ReadinessProbe.TCPSocket.Port,
				Host: src.Spec.ReadinessProbe.TCPSocket.Host,
			}
		}
	}
	if src.Spec.Reserved != nil {
		dst.Spec.Reserved = &vmoprv1alpha2.VirtualMachineReservedSpec{
			ResourcePolicyName: src.Spec.Reserved.ResourcePolicyName,
		}
	}
	dst.Spec.StorageClass = src.Spec.StorageClass
	if src.Spec.Volumes != nil {
		dst.Spec.Volumes = []vmoprv1alpha2.VirtualMachineVolume{}
		for _, volume := range src.Spec.Volumes {
			v := vmoprv1alpha2.VirtualMachineVolume{}
			// ApplicationType existing in hub but not in v1alpha2.VirtualMachineVolume
			// ControllerBusNumber existing in hub but not in v1alpha2.VirtualMachineVolume
			// ControllerType existing in hub but not in v1alpha2.VirtualMachineVolume
			// DiskMode existing in hub but not in v1alpha2.VirtualMachineVolume
			v.Name = volume.Name
			if volume.PersistentVolumeClaim != nil {
				v.PersistentVolumeClaim = &vmoprv1alpha2.PersistentVolumeClaimVolumeSource{
					PersistentVolumeClaimVolumeSource: corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: volume.PersistentVolumeClaim.ClaimName,
						ReadOnly:  volume.PersistentVolumeClaim.ReadOnly,
					},
				}
				if volume.PersistentVolumeClaim.InstanceVolumeClaim != nil {
					v.PersistentVolumeClaim.InstanceVolumeClaim = &vmoprv1alpha2.InstanceVolumeClaimVolumeSource{
						StorageClass: volume.PersistentVolumeClaim.InstanceVolumeClaim.StorageClass,
						Size:         volume.PersistentVolumeClaim.InstanceVolumeClaim.Size,
					}
				}

			}
			// Removable existing in hub but not in v1alpha2.VirtualMachineVolume
			// SharingMode existing in hub but not in v1alpha2.VirtualMachineVolume
			// UnitNumber existing in hub but not in v1alpha2.VirtualMachineVolume
			dst.Spec.Volumes = append(dst.Spec.Volumes, v)
		}
	}

	dst.Status.BiosUUID = src.Status.BiosUUID
	if src.Status.Conditions != nil {
		dst.Status.Conditions = []metav1.Condition{}
		for _, condition := range src.Status.Conditions {
			dst.Status.Conditions = append(dst.Status.Conditions, condition)
		}
	}
	dst.Status.Host = src.Status.NodeName // Field renamed in hub
	if src.Status.Network != nil {
		dst.Status.Network = &vmoprv1alpha2.VirtualMachineNetworkStatus{}
		dst.Status.Network.HostName = src.Status.Network.HostName
		if src.Status.Network.Interfaces != nil {
			dst.Status.Network.Interfaces = []vmoprv1alpha2.VirtualMachineNetworkInterfaceStatus{}
			for _, iface := range src.Status.Network.Interfaces {
				d := vmoprv1alpha2.VirtualMachineNetworkInterfaceStatus{}
				d.DeviceKey = iface.DeviceKey
				if iface.DNS != nil {
					d.DNS = &vmoprv1alpha2.VirtualMachineNetworkDNSStatus{
						DHCP:          iface.DNS.DHCP,
						HostName:      iface.DNS.HostName,
						DomainName:    iface.DNS.DomainName,
						Nameservers:   iface.DNS.Nameservers,
						SearchDomains: iface.DNS.SearchDomains,
					}
				}
				if iface.IP != nil {
					d.IP = &vmoprv1alpha2.VirtualMachineNetworkInterfaceIPStatus{
						MACAddr: iface.IP.MACAddr,
					}
					if iface.IP.Addresses != nil {
						d.IP.Addresses = []vmoprv1alpha2.VirtualMachineNetworkInterfaceIPAddrStatus{}
						for _, addr := range iface.IP.Addresses {
							d.IP.Addresses = append(d.IP.Addresses, vmoprv1alpha2.VirtualMachineNetworkInterfaceIPAddrStatus{
								Address:  addr.Address,
								Lifetime: addr.Lifetime,
								Origin:   addr.Origin,
								State:    addr.State,
							})
						}
					}
					if iface.IP.AutoConfigurationEnabled != nil {
						d.IP.AutoConfigurationEnabled = ptr.To(*iface.IP.AutoConfigurationEnabled)
					}
					if iface.IP.DHCP != nil {
						d.IP.DHCP = &vmoprv1alpha2.VirtualMachineNetworkDHCPStatus{
							IP4: vmoprv1alpha2.VirtualMachineNetworkDHCPOptionsStatus{
								Enabled: iface.IP.DHCP.IP4.Enabled,
							},
							IP6: vmoprv1alpha2.VirtualMachineNetworkDHCPOptionsStatus{
								Enabled: iface.IP.DHCP.IP6.Enabled,
							},
						}
						if iface.IP.DHCP.IP4.Config != nil {
							d.IP.DHCP.IP4.Config = []vmoprv1alpha2common.KeyValuePair{}
							for _, pair := range iface.IP.DHCP.IP4.Config {
								d.IP.DHCP.IP4.Config = append(d.IP.DHCP.IP4.Config, vmoprv1alpha2common.KeyValuePair{
									Key:   pair.Key,
									Value: pair.Value,
								})
							}
						}
						if iface.IP.DHCP.IP6.Config != nil {
							d.IP.DHCP.IP6.Config = []vmoprv1alpha2common.KeyValuePair{}
							for _, pair := range iface.IP.DHCP.IP6.Config {
								d.IP.DHCP.IP6.Config = append(d.IP.DHCP.IP6.Config, vmoprv1alpha2common.KeyValuePair{
									Key:   pair.Key,
									Value: pair.Value,
								})
							}
						}
					}
				}
				d.Name = iface.Name
				dst.Status.Network.Interfaces = append(dst.Status.Network.Interfaces, d)
			}
		}
		dst.Status.Network.PrimaryIP4 = src.Status.Network.PrimaryIP4
		dst.Status.Network.PrimaryIP6 = src.Status.Network.PrimaryIP6
	}
	// Policies existing in hub but not in v1alpha2.VirtualMachineStatus
	dst.Status.PowerState = vmoprv1alpha2.VirtualMachinePowerState(src.Status.PowerState)
	dst.Status.Zone = src.Status.Zone

	return nil
}

func init() {
	converterBuilder.AddConversion(
		conversion.NewAddConversionBuilder(convert_hub_VirtualMachine_To_v1alpha2_VirtualMachine, convert_v1alpha2_VirtualMachine_To_hub_VirtualMachine),
	)
}
