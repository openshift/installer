package asset

//go:generate mockgen -source=./filefetcher.go -destination=./mock/filefetcher_generated.go -package=mock

// FileFetcher fetches the asset files from disk.
type FileFetcher interface {
	// FetchByName returns the file with the given name.
	FetchByName(string) (*File, error)
	// FetchByPattern returns the files whose name match the given glob.
	FetchByPattern(pattern string) ([]*File, error)
}
