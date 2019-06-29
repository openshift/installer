package providers

// Registry maps ClusterMetadata.Platform() to per-platform Destroyer creators.
var Registry = make(map[string]NewFunc)
