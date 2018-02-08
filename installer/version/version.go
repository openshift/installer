package version

// TectonicVersion - defaults to 'dev-build' if TECTONIC_VERSION was not set
var TectonicVersion = "dev-build"

// BuildTime - set to $(date -u '+%d_%b_%Y_%I:%M:%S%p')
var BuildTime string
