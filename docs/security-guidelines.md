# Security Guidelines

Rules and conventions for security-sensitive code in the OpenShift Installer.

## TLS Certificates

### Key size and algorithm
All RSA keys use 2048-bit via `rsa.GenerateKey(rand.Reader, 2048)` in `pkg/asset/tls/tls.go`. The constant `keySize = 2048` is the single source of truth. Do not introduce alternative key sizes without updating this constant.

### Certificate validity tiers
Use the existing validity helpers from `pkg/asset/tls/tls.go` -- never hardcode durations:

- `ValidityOneDay(installConfig)` -- short-lived certs (aggregator, some apiserver certs). Returns 24h normally, 2h when `ShortCertRotation` feature gate is enabled.
- `ValidityOneYear(installConfig)` -- medium-lived certs (kube control plane, apiserver CA). Returns 365 days normally, 4h when short rotation is enabled.
- `ValidityTenYears()` -- long-lived certs (root CA, admin kubeconfig signer, MCS cert, journal). Not affected by feature gates.

Choose the shortest tier that satisfies the cert's purpose. Leaf certs that will be auto-rotated by operators should use `ValidityOneDay`. CA certs and certs consumed only at install time may use `ValidityTenYears`.

### Certificate hierarchy
Every `SignedCertKey` must chain to a parent `CertKeyInterface`. Self-signed certs (`SelfSignedCertKey`) are only for CA certificates. The `CertCfg.Subject` must always set both `CommonName` and `OrganizationalUnit`; omitting either causes `SelfSignedCertificate` to return an error.

### Loading certs from disk
Certs are not loaded from disk by default (`Load` returns `false, nil`). The env var `OPENSHIFT_INSTALL_LOAD_CLUSTER_CERTS=true` enables loading pre-existing certs via `loadCertKey`. This is an intentional security boundary -- do not bypass it.

### Bound SA signing key
`BoundSASigningKey` (`pkg/asset/tls/boundsasigningkey.go`) is user-provided only -- its `Generate` is a no-op. The `Load` method validates that the key parses as a valid RSA private key before accepting it.

## Randomness

### Always use `crypto/rand`
All security-sensitive random values (passwords, keys, tokens, credentials) must use `crypto/rand.Reader` or `crypto/rand.Int`. The codebase follows this consistently in `pkg/asset/password/`, `pkg/asset/tls/`, `pkg/asset/ignition/bootstrap/baremetal/ironic_creds.go`, and `pkg/asset/agent/gencrypto/`.

Uses of `math/rand` exist in `pkg/types/baremetal/defaults/` (provisioning MAC generation) and `pkg/infrastructure/azure/` (non-security retry jitter). If you need randomness for any credential, token, or key material, use `crypto/rand` exclusively.

## Passwords and Credentials

### Kubeadmin password format
Generated in `pkg/asset/password/password.go`: 23 random alphanumeric characters formatted as `5char-5char-5char-5char` (with characters at indices 5, 11, and 17 replaced by dashes). Excludes ambiguous characters (0, 1, O, l). Hashed with `bcrypt.DefaultCost`. The plaintext is written to `auth/kubeadmin-password`; the hash is stored separately in `tls/kubeadmin-password.hash`.

### IRI registry credentials
Generated in `pkg/asset/tls/iriregistryauth.go` with 32 bytes (256-bit) of `crypto/rand` entropy, base64-encoded, and bcrypt-hashed into htpasswd format. This asset is in-memory only -- it must NOT write files to the `auth/` directory, because assisted-service deletes that directory and extra files break deployment.

### Ironic credentials
Generated in `pkg/asset/ignition/bootstrap/baremetal/ironic_creds.go` with 16-character random alphanumeric passwords via `crypto/rand`. Username is always `bootstrap-user`.

### Cloud credentials in manifests
Cloud credentials are base64-encoded and injected into Kubernetes secrets via Go templates in `pkg/asset/manifests/openshift.go`. Each platform has a typed struct in `pkg/asset/manifests/template.go` (e.g., `AwsCredsSecretData`, `AzureCredsSecretData`). The resulting secrets are scoped to `kube-system` namespace with RBAC limiting access via `role-cloud-creds-secret-reader.yaml.template`.

### Credential modes
`CredentialsMode` in `pkg/types/installconfig.go` controls how the cloud-credential-operator satisfies requests:
- `Manual` -- operator does not process credential requests (most restrictive)
- `Mint` -- operator creates scoped users per request
- `Passthrough` -- operator copies the install-time credential

When `CredentialsMode` is unset, the installer queries cloud permissions before proceeding. When it is set, permission checks are skipped -- the user takes responsibility.

## FIPS Compliance

### Host validation
`pkg/hostcrypt/` validates that the installer binary and host match the target cluster's FIPS setting. Two build variants exist:
- `fipscapable` build tag: requires `/proc/sys/crypto/fips_enabled` to be `1`
- Non-FIPS binary (default): always rejects FIPS clusters, directing users to the RHEL 9 FIPS binary

The env var `OPENSHIFT_INSTALL_SKIP_HOSTCRYPT_VALIDATION` bypasses this check (adds annotation `hostCryptBypassedAnnotation` to the config). This should only be used in testing.

### SSH key restrictions under FIPS
When `installConfig.FIPS == true`, SSH keys are validated against FIPS-allowed algorithms: only `ssh-rsa` and `ecdsa-sha2-nistp*` are accepted. Ed25519 and other types are rejected. This validation is in `validateFIPSconfig` in `pkg/types/validation/installconfig.go`.

## SSH Key Validation

SSH public keys are validated via `golang.org/x/crypto/ssh.ParseAuthorizedKey` in `pkg/validate/validate.go:SSHPublicKey`. The install-config validation supports newline-separated multiple keys. PowerVS platform requires `sshKey` to be non-empty.

## Pull Secret Validation

`pkg/validate/validate.go:ImagePullSecret` validates pull secrets by:
1. Parsing as JSON with `{"auths": {...}}` structure
2. Requiring at least one entry in `auths`
3. Requiring each registry entry to contain either an `auth` field or a `credsStore` field

Mirror registry hosts are cross-checked against pull secret entries; missing credentials produce warnings (not errors) via `validateMirrorCredentials`.

## File Permissions

Asset files are written with mode `0o640`. Credential files (Azure creds, PowerVS auth, oVirt config) use `0o600`. Do not use more permissive modes for any file containing secrets or keys.

## Agent Auth Tokens (JWT)

`pkg/asset/agent/gencrypto/authconfig.go` generates ECDSA P-256 key pairs and ES256-signed JWT tokens for agent-based installs. Three personas exist: `agentAuth`, `userAuth`, `watcherAuth`. Install workflow tokens have no expiry; add-nodes tokens expire after 48 hours and are refreshed when within 24 hours of expiry. The public key and tokens are stored in a `Secret` in `openshift-config` namespace.

## Gosec Annotations

Use `//nolint:gosec` or `// #nosec` only with a justification comment (e.g., `// not a hardcoded secret`, `// no sensitive info`). The codebase convention is to include the reason inline. False positives on constant names containing "secret", "password", or "token" are common and should be annotated.
