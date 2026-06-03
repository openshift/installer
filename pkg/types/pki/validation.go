package pki

import (
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
)

// ValidatePKIConfig validates the PKI configuration.
// When pkiConfig is non-nil, signerCertificates must be fully specified.
func ValidatePKIConfig(pkiConfig *types.PKIConfig, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if pkiConfig == nil {
		return allErrs
	}

	allErrs = append(allErrs, ValidateKeyConfig(pkiConfig.SignerCertificates.Key,
		fldPath.Child("signerCertificates", "key"))...)

	return allErrs
}

// ValidateKeyConfig validates the KeyConfig structure.
func ValidateKeyConfig(config types.KeyConfig, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if config.Algorithm == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("algorithm"),
			"algorithm must be specified"))
		return allErrs
	}

	validAlgorithms := []types.KeyAlgorithm{types.KeyAlgorithmRSA, types.KeyAlgorithmECDSA}
	if !sets.New(validAlgorithms...).Has(config.Algorithm) {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("algorithm"),
			config.Algorithm, validAlgorithms))
		return allErrs
	}

	if config.Algorithm == types.KeyAlgorithmRSA {
		if config.RSA == nil {
			allErrs = append(allErrs, field.Required(fldPath.Child("rsa"),
				"rsa is required when algorithm is RSA"))
		} else {
			allErrs = append(allErrs, validateRSAKeyConfig(*config.RSA, fldPath.Child("rsa"))...)
		}

		if config.ECDSA != nil {
			allErrs = append(allErrs, field.Forbidden(fldPath.Child("ecdsa"),
				"ecdsa must not be set when algorithm is RSA"))
		}
	}

	if config.Algorithm == types.KeyAlgorithmECDSA {
		if config.ECDSA == nil {
			allErrs = append(allErrs, field.Required(fldPath.Child("ecdsa"),
				"ecdsa is required when algorithm is ECDSA"))
		} else {
			allErrs = append(allErrs, validateECDSAKeyConfig(*config.ECDSA, fldPath.Child("ecdsa"))...)
		}

		if config.RSA != nil {
			allErrs = append(allErrs, field.Forbidden(fldPath.Child("rsa"),
				"rsa must not be set when algorithm is ECDSA"))
		}
	}

	return allErrs
}

func validateRSAKeyConfig(config types.RSAKeyConfig, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if config.KeySize < 2048 || config.KeySize > 8192 || config.KeySize%1024 != 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("keySize"), config.KeySize,
			"must be a multiple of 1024 from 2048 to 8192"))
	}

	return allErrs
}

func validateECDSAKeyConfig(config types.ECDSAKeyConfig, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	validCurves := []types.ECDSACurve{types.ECDSACurveP256, types.ECDSACurveP384, types.ECDSACurveP521}
	if !sets.New(validCurves...).Has(config.Curve) {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("curve"),
			config.Curve, validCurves))
	}

	return allErrs
}
