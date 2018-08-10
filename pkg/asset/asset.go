package asset

type Asset interface {
	Dependencies() []Asset
	Generate(map[Asset]*State) (*State, error)
}
