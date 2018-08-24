package asset

// Store is a store for the states of assets.
type Store interface {
	// Fetch retrieves the state of the given asset, generating it and its
	// dependencies if necessary.
	Fetch(Asset) (*State, error)
}

// StoreImpl is the implementation of Store.
type StoreImpl struct {
	assets map[Asset]*State
}

// Fetch retrieves the state of the given asset, generating it and its
// dependencies if necessary.
func (s *StoreImpl) Fetch(asset Asset) (*State, error) {
	state, ok := s.assets[asset]
	if ok {
		return state, nil
	}
	dependies := asset.Dependencies()
	dependenciesStates := make(map[Asset]*State, len(dependies))
	for _, d := range dependies {
		ds, err := s.Fetch(d)
		if err != nil {
			return nil, err
		}
		dependenciesStates[d] = ds
	}
	state, err := asset.Generate(dependenciesStates)
	if err != nil {
		return nil, err
	}
	if s.assets == nil {
		s.assets = make(map[Asset]*State)
	}
	s.assets[asset] = state
	return state, nil
}
