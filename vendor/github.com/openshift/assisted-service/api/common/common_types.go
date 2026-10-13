package common

// +genclient

// +kubebuilder:object:generate=true
type ValidationResult struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// +kubebuilder:object:root=false
// +kubebuilder:object:generate=true
// ValidationsStatus is the Schema for the ValidationsInfo field
type ValidationsStatus map[string]ValidationResults

// +kubebuilder:object:generate=true
type ValidationResults []ValidationResult
