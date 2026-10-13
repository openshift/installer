package providers

// Registry maps ClusterMetadata.Platform() to per-platform Gather methods.
var Registry = make(map[string]NewFunc)
