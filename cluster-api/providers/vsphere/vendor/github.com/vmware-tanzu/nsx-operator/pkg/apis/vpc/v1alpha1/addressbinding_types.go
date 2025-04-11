package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type AddressBindingSpec struct {
	// VMName contains the VM's name
	VMName string `json:"vmName"`
	// InterfaceName contains the interface name of the VM, if not set, the first interface of the VM will be used
	InterfaceName string `json:"interfaceName,omitempty"`
}

type AddressBindingStatus struct {
	IPAddress string `json:"ipAddress"`
}

// +genclient
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:storageversion

// AddressBinding is used to manage 1:1 NAT for a VM/NetworkInterface.
type AddressBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AddressBindingSpec   `json:"spec"`
	Status AddressBindingStatus `json:"status"`
}

//+kubebuilder:object:root=true

// AddressBindingList contains a list of AddressBinding.
type AddressBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AddressBinding `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AddressBinding{}, &AddressBindingList{})
}
