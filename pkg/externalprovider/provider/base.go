package provider

// ExternalProvider is a provider that is not maintained by the core installer team.
type ExternalProvider interface {
	// Name returns the name of the external provider.
	Name() string

	// InstallConfigExternalProvider contains the installconfig-specific steps.
	InstallConfigExternalProvider
}
