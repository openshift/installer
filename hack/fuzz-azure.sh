#!/usr/bin/env bash
# Run the installer's Azure fuzz targets. Override FUZZ_DURATION for shorter
# local runs, for example: FUZZ_DURATION=5s ./hack/fuzz-azure.sh
set -uo pipefail

rc=0
FUZZ_WORKERS="${FUZZ_WORKERS:-2}"
FUZZ_DURATION="${FUZZ_DURATION:-2m}"
if [[ -z "${ARTIFACT_DIR:-}" ]]; then
  artifact_tmp_dir="${TMPDIR:-/tmp}"
  artifact_tmp_dir="${artifact_tmp_dir%/}"
  ARTIFACT_DIR=$(mktemp -d "$artifact_tmp_dir/fuzz-azure-artifacts.XXXXXX")
  echo "New fuzz corpus inputs will be saved under $ARTIFACT_DIR"
fi
FUZZ_DIRS=(
  pkg/destroy/azure
  pkg/types/azure
  pkg/types/azure/validation
  pkg/asset/installconfig/azure
  pkg/asset/machines/azure
  pkg/asset/manifests/azure
  pkg/rhcos/marketplace/azure
)

find_fuzz_corpus() {
  local dir
  for dir in "${FUZZ_DIRS[@]}"; do
    find "$dir" -path "*/testdata/fuzz/*/*" -type f -print
  done
}

corpus_snapshot=$(mktemp)
find_fuzz_corpus > "$corpus_snapshot"

copy_new_fuzz_corpus() {
  local f dest
  while IFS= read -r f; do
    if ! grep -Fqx -- "$f" "$corpus_snapshot"; then
      dest="$ARTIFACT_DIR/$(dirname "$f")"
      mkdir -p "$dest"
      cp "$f" "$dest/"
    fi
  done < <(find_fuzz_corpus)
  rm -f "$corpus_snapshot"
}
trap copy_new_fuzz_corpus EXIT

run_fuzz() {
  local pkg="$1" fn="$2" out
  echo "--- fuzz: $fn ($pkg) ---"
  out=$(mktemp)
  go test "./$pkg" -run='^$' -fuzz="^${fn}$" -parallel="$FUZZ_WORKERS" -fuzztime="$FUZZ_DURATION" 2>&1 | tee "$out" || rc=1
  # Go may exit 0 when a target or all test files are missing.
  if grep -Fq "no fuzz tests to fuzz" "$out" || grep -Fq "[no test files]" "$out"; then rc=1; fi
  rm -f "$out"
}

run_fuzz pkg/destroy/azure                          FuzzGroupFromID
run_fuzz pkg/types/azure                            FuzzAddressFamilySubnets
run_fuzz pkg/types/azure                            FuzzAzureResourceNames
run_fuzz pkg/types/azure                            FuzzSetBaseDomain
run_fuzz pkg/types/azure                            FuzzStorageAccountName
run_fuzz pkg/types/azure/validation                 FuzzDiskEncryptionSetID
run_fuzz pkg/types/azure/validation                 FuzzFindDuplicateTagKeys
run_fuzz pkg/types/azure/validation                 FuzzValidateDataDisk
run_fuzz pkg/types/azure/validation                 FuzzValidateDataDiskSetup
run_fuzz pkg/types/azure/validation                 FuzzValidateDiskEncryption
run_fuzz pkg/types/azure/validation                 FuzzValidateEncryptionAtHost
run_fuzz pkg/types/azure/validation                 FuzzValidateOSImage
run_fuzz pkg/types/azure/validation                 FuzzValidateSecurityProfile
run_fuzz pkg/types/azure/validation                 FuzzValidateTag
run_fuzz pkg/types/azure/validation                 FuzzValidateUserTags
run_fuzz pkg/asset/installconfig/azure              FuzzAcceleratedNetworkingCapability
run_fuzz pkg/asset/installconfig/azure              FuzzAzureStackClusterOSImage
run_fuzz pkg/asset/installconfig/azure              FuzzAzureStackClusterOSImageSet
run_fuzz pkg/asset/installconfig/azure              FuzzInstanceSecurityType
run_fuzz pkg/asset/installconfig/azure              FuzzPremiumDiskCapability
run_fuzz pkg/asset/installconfig/azure              FuzzWindowsVMFamily
run_fuzz pkg/asset/machines/azure                   FuzzBootDiagnosticsObject
run_fuzz pkg/asset/machines/azure                   FuzzBootstrapFailureDomain
run_fuzz pkg/asset/machines/azure                   FuzzCapzUserTags
run_fuzz pkg/asset/machines/azure                   FuzzMachineImage
run_fuzz pkg/asset/machines/azure                   FuzzTrimSubscriptionPrefix
run_fuzz pkg/asset/manifests/azure                  FuzzGetIPWithinCIDR
run_fuzz pkg/asset/manifests/azure                  FuzzSubnetAddressPrefixes
run_fuzz pkg/rhcos/marketplace/azure                FuzzCheckIfNewer
run_fuzz pkg/rhcos/marketplace/azure                FuzzConvertToSemver
run_fuzz pkg/rhcos/marketplace/azure                FuzzMatchesRHELMajor
run_fuzz pkg/rhcos/marketplace/azure                FuzzParseAROSKUs

exit "$rc"
