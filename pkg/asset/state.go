package asset

// State is the state of an Asset.
type State struct {
	Contents []Content
}

type Content struct {
	Name string // the path on disk for this content.
	Data []byte
}
