package server

type certificatesStrategy string

const (
	certificatesStrategyInstaller certificatesStrategy = "installerGeneratedCA"
	certificatesStrategyUser      certificatesStrategy = "userProvidedCA"
)

type metrics struct {
	// what was the platform submitted to the installer?
	installerPlatform string
	// what strategy is the cluster using for certificate infrastructure ?
	certificatesStrategy certificatesStrategy
	// is the Tectonic update operator enabled?
	tectonicUpdaterEnabled bool
}

// getCertificatesStrategy returns a certificatesStrategy that describes the
// certificates in use for the newly created cluster.
func getCertificatesStrategy(caCertificate string) certificatesStrategy {
	if caCertificate == "" {
		return certificatesStrategyInstaller
	}
	return certificatesStrategyUser
}
