package asset

// State is the state of an Asset.
type State struct {
	Contents []Content
}

// Content is a generated portion of an Asset.
type Content struct {
	Name string // the path on disk for this content.
	Data []byte
}
