package asset

// Asset used to install OpenShift.
type Asset interface {
	// Dependencies returns the assets upon which this asset directly depends.
	Dependencies() []Asset

	// Generate generates this asset given the states of its dependent assets.
	Generate(map[Asset]*State) (*State, error)

	// Name returns the human-friendly name of the asset.
	Name() string
}
