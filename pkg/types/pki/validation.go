package pki

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	configv1alpha1 "github.com/openshift/api/config/v1alpha1"
	"github.com/openshift/installer/pkg/types"
)

// ValidatePKIConfig validates the PKI configuration.
// When pkiConfig is non-nil, signerCertificates must be fully specified.
// NOTE: All fields are value types (not pointers). Use zero-value checks.
func ValidatePKIConfig(pkiConfig *types.PKIConfig, fldPath *field.Path, fips bool) field.ErrorList {
	allErrs := field.ErrorList{}

	if pkiConfig == nil {
		return allErrs
	}

	// signerCertificates.key must be fully specified when pki is present
	if pkiConfig.SignerCertificates.Key.Algorithm == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("signerCertificates", "key"),
			"signerCertificates.key is required when pki is specified"))
		return allErrs
	}

	allErrs = append(allErrs, ValidateKeyConfig(pkiConfig.SignerCertificates.Key,
		fldPath.Child("signerCertificates", "key"), fips)...)

	return allErrs
}

// ValidateKeyConfig validates the KeyConfig structure.
// KeyConfig fields are value types: RSA is RSAKeyConfig, ECDSA is ECDSAKeyConfig.
// Use zero-value checks (KeySize == 0, Curve == "") instead of nil checks.
func ValidateKeyConfig(config configv1alpha1.KeyConfig, fldPath *field.Path, fips bool) field.ErrorList {
	allErrs := field.ErrorList{}

	if config.Algorithm == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("algorithm"),
			"algorithm must be specified"))
		return allErrs
	}

	if config.Algorithm != configv1alpha1.KeyAlgorithmRSA && config.Algorithm != configv1alpha1.KeyAlgorithmECDSA {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("algorithm"),
			config.Algorithm, []string{string(configv1alpha1.KeyAlgorithmRSA), string(configv1alpha1.KeyAlgorithmECDSA)}))
		return allErrs
	}

	if config.Algorithm == configv1alpha1.KeyAlgorithmRSA {
		if config.RSA.KeySize == 0 {
			allErrs = append(allErrs, field.Required(fldPath.Child("rsa", "keySize"),
				"keySize must be specified when algorithm is RSA"))
		} else {
			allErrs = append(allErrs, validateRSAKeyConfig(config.RSA, fldPath.Child("rsa"), fips)...)
		}

		if config.ECDSA.Curve != "" {
			allErrs = append(allErrs, field.Forbidden(fldPath.Child("ecdsa"),
				"ecdsa must not be set when algorithm is RSA"))
		}
	}

	if config.Algorithm == configv1alpha1.KeyAlgorithmECDSA {
		if config.ECDSA.Curve == "" {
			allErrs = append(allErrs, field.Required(fldPath.Child("ecdsa", "curve"),
				"curve must be specified when algorithm is ECDSA"))
		} else {
			allErrs = append(allErrs, validateECDSAKeyConfig(config.ECDSA, fldPath.Child("ecdsa"), fips)...)
		}

		if config.RSA.KeySize != 0 {
			allErrs = append(allErrs, field.Forbidden(fldPath.Child("rsa"),
				"rsa must not be set when algorithm is ECDSA"))
		}
	}

	return allErrs
}

func validateRSAKeyConfig(config configv1alpha1.RSAKeyConfig, fldPath *field.Path, fips bool) field.ErrorList {
	allErrs := field.ErrorList{}

	// Validate key size — aligned with API kubebuilder validation:
	// multiples of 1024 from 2048 to 8192
	if config.KeySize < 2048 || config.KeySize > 8192 || config.KeySize%1024 != 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("keySize"), config.KeySize,
			"must be a multiple of 1024 from 2048 to 8192"))
	}

	return allErrs
}

func validateECDSAKeyConfig(config configv1alpha1.ECDSAKeyConfig, fldPath *field.Path, fips bool) field.ErrorList {
	allErrs := field.ErrorList{}

	validCurves := []configv1alpha1.ECDSACurve{
		configv1alpha1.ECDSACurveP256,
		configv1alpha1.ECDSACurveP384,
		configv1alpha1.ECDSACurveP521,
	}
	valid := false
	for _, curve := range validCurves {
		if config.Curve == curve {
			valid = true
			break
		}
	}

	if !valid {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("curve"), config.Curve,
			"must be P256, P384, or P521"))
	}

	return allErrs
}
