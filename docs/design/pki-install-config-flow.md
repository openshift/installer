# PKI InstallConfig Field Flow

## Layer 1: Upstream API types (vendored)

The foundation is the `configv1alpha1` types from `github.com/openshift/api` (PR openshift/api#2645), vendored at `vendor/github.com/openshift/api/config/v1alpha1/types_pki.go`. This defines the full cluster-level `PKI` CR with its discriminated union types:

- **`KeyConfig`** — a `+union` with `Algorithm` as the discriminator, and `RSA`/`ECDSA` as members
- **`RSAKeyConfig`** — just `KeySize int32` (multiples of 1024, 2048-8192)
- **`ECDSAKeyConfig`** — just `Curve ECDSACurve` (P256, P384, P521)
- **`CertificateConfig`** — wraps a `KeyConfig` in a `Key` field
- **`PKIProfile`** — has `Defaults`, `SignerCertificates`, `ServingCertificates`, `ClientCertificates`

These are all **value types** (no pointers), so zero-value checks (`KeySize == 0`, `Curve == ""`) are used instead of nil checks.

## Layer 2: InstallConfig wrapper type

In `pkg/types/installconfig.go`, a **new local type** is defined rather than using the upstream `PKIProfile` directly:

```go
PKI *PKIConfig `json:"pki,omitempty"`
```

```go
type PKIConfig struct {
    SignerCertificates configv1alpha1.CertificateConfig `json:"signerCertificates"`
}
```

This is a deliberate narrowing — the installer only exposes `signerCertificates` (not defaults, serving, or client certs), because the installer only generates signer CAs. The upstream types are reused for the leaf structure (`CertificateConfig` -> `KeyConfig` -> `RSA`/`ECDSA`), keeping the YAML shape consistent with the cluster PKI CR.

The field is a pointer (`*PKIConfig`) so it's truly optional — `nil` means "not configured."

## Layer 3: CRD schema

The CRD at `data/data/install.openshift.io_installconfigs.yaml` gets the corresponding OpenAPI schema with all the kubebuilder validations from the upstream types (`enum`, `minimum`, `maximum`, `multipleOf`, `required`, plus CEL `x-kubernetes-validations` rules for the union discriminator). This was regenerated with `go generate ./pkg/types/installconfig.go`.

## Layer 4: Deep copy

`pkg/types/zz_generated.deepcopy.go` gets auto-generated additions:
- `InstallConfig.DeepCopyInto` handles the `*PKIConfig` pointer (allocates new, copies value)
- `PKIConfig.DeepCopyInto` copies `SignerCertificates` by value assignment

This works correctly because `CertificateConfig` and its nested types (`KeyConfig`, `RSAKeyConfig`, `ECDSAKeyConfig`) are all value types with no pointers — a plain struct copy is a full deep copy.

## Layer 5: Feature gate

In `pkg/types/validation/featuregates.go`, the field is registered as a gated feature:

```go
{
    FeatureGateName: features.FeatureGateConfigurablePKI,
    Condition:       c.PKI != nil,
    Field:           field.NewPath("pki"),
},
```

This means setting `pki:` in the install-config requires `ConfigurablePKI` to be enabled (via `TechPreviewNoUpgrade` or `CustomNoUpgrade` with the gate explicitly listed). If the gate is off and `pki` is set, `validateGatedFeatures()` produces an error.

## Layer 6: Validation

In `pkg/types/validation/installconfig.go`:

```go
if c.PKI != nil {
    allErrs = append(allErrs, pkivalidation.ValidatePKIConfig(c.PKI, field.NewPath("pki"), c.FIPS)...)
}
```

This delegates to `pkg/types/pki/validation.go` which validates:
1. `signerCertificates.key.algorithm` is required and must be `RSA` or `ECDSA`
2. For RSA: `rsa.keySize` required, must be multiple of 1024 in [2048, 8192]; `ecdsa` must not be set
3. For ECDSA: `ecdsa.curve` required, must be P256/P384/P521; `rsa` must not be set

Note: PKI validation runs **before** the feature gate check at `ValidateFeatureSet()`, so a user with an invalid PKI config and no feature gate will see both errors.

## Layer 7: Defaults / effective config

`pkg/types/pki/defaults.go` provides `EffectiveSignerPKIConfig(ic)` which is the central function all signer cert assets call:

| Scenario | Result |
|----------|--------|
| `ic.PKI != nil` | Returns user's config as-is |
| `ic.PKI == nil` + gate enabled | Returns synthetic RSA-4096 (from `DefaultPKIProfile().SignerCertificates`) |
| `ic.PKI == nil` + gate disabled | Returns `nil` -> legacy RSA-2048 path |

This means enabling the feature gate *without* specifying `pki:` still upgrades signer certs from RSA-2048 to RSA-4096, matching the upstream DefaultPKIProfile.

## Summary flow

```
install-config.yaml
  +-- pki.signerCertificates.key.algorithm: ECDSA
     pki.signerCertificates.key.ecdsa.curve: P384

  | parsed into

types.InstallConfig.PKI (*types.PKIConfig)
  +-- .SignerCertificates (configv1alpha1.CertificateConfig)
       +-- .Key.Algorithm = "ECDSA"
          .Key.ECDSA.Curve = "P384"

  | validated by

validation/featuregates.go  -> requires ConfigurablePKI gate
pki/validation.go           -> validates algorithm+params consistency

  | consumed by

pki/defaults.go:EffectiveSignerPKIConfig()
  | returns PKIConfig to

pkg/asset/tls/*.go  -> each signer asset calls Generate(ctx, cfg, name, pkiConfig)
  | which calls

tls.go:GenerateRSAPrivateKey() or GenerateECDSAPrivateKey()
  -> produces signer CA with configured algorithm
```