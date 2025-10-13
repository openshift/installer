package powervc

const (
	// SGPolicyUnset represents the default empty string for the ServerGroupPolicy field of the MachinePool.
	SGPolicyUnset ServerGroupPolicy = ""

	// SGPolicyAffinity represents the "affinity" ServerGroupPolicy field of the MachinePool.
	SGPolicyAffinity ServerGroupPolicy = "affinity"

	// SGPolicySoftAffinity represents the "soft-affinity" ServerGroupPolicy field of the MachinePool.
	SGPolicySoftAffinity ServerGroupPolicy = "soft-affinity"

	// SGPolicyAntiAffinity represents the "anti-affinity" ServerGroupPolicy field of the MachinePool.
	SGPolicyAntiAffinity ServerGroupPolicy = "anti-affinity"

	// SGPolicySoftAntiAffinity represents the "soft-anti-affinity" ServerGroupPolicy field of the MachinePool.
	SGPolicySoftAntiAffinity ServerGroupPolicy = "soft-anti-affinity"
)

// ServerGroupPolicy is the policy to be applied to an OpenStack Server Group.
//
// +kubebuilder:validation:Enum="";affinity;soft-affinity;anti-affinity;soft-anti-affinity
// +optional
//
//nolint:godot
type ServerGroupPolicy string

// IsSet returns true when p is not the empty string.
func (p ServerGroupPolicy) IsSet() bool {
	return p != SGPolicyUnset
}
