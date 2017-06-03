#!/usr/bin/env bash

### Usage ###
# collect_project_licenses.sh aggregates license information from a packages'
# immediate dependencies for both Go and JS repositories listed in
# 'pkg_lists/go_pkg_inputs.txt' and 'pkg_lists/js_pkg_inputs.txt', respectively.
# Go dependency licenses are aggregated using 'license-bill-of-materials', a
# program that finds licenses using the 'go list' command (see
# https://github.com/coreos/license-bill-of-materials). JS dependency licenses
# are aggregated by parsing all node_modules/*/package.json files, which
# contain a 'license(s)' field. Each set of licenses has the form:
# [
#   {
#     "package": "package name",
#     "licenses": [{
#       "type": "license type",
#       "confidence": N
#     }, ...]
#   }
# ]
# Once all dependency licenses have been aggregated into JSON arrays in files,
# all arrays are concatenated into one large array in FINAL_LICENSE_FILE.

### Dependencies ###
# bash, jq, yarn, git, license-bill-of-materials

set -eu

# Set relative environment
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$DIR/../../.."
FINAL_LICENSE_FILE="$ROOT"/bill-of-materials.json
echo "[]" > "$FINAL_LICENSE_FILE"

# Temporary data locations
TMP_DIR=$(mktemp -d -p /tmp)
TMP_LICENSE_DIR="$TMP_DIR"/licenses
mkdir -p "$TMP_LICENSE_DIR"
tmp_shared_file=$(mktemp -p "$TMP_DIR")

# Package input files
go_pkg_inputs="$DIR"/pkg_lists/go_pkg_inputs.txt
js_pkg_inputs="$DIR"/pkg_lists/js_pkg_inputs.txt

# Temporarily create a new GOPATH
old_gopath="$GOPATH"
new_gopath="$TMP_DIR"/tmp_gopath
export GOPATH="$new_gopath"
mkdir -p "$GOPATH"/src

# Clean up transient files and unset vars
clean_up() {
  rm -rf "$TMP_DIR"
  export GOPATH="$old_gopath"
  exit
}
trap clean_up SIGHUP SIGINT SIGTERM EXIT

# Parse JS deps using jq and send to dst
parse_js_licenses() {
  local src="${1:-}"
  local dst="${2:-}"

  echo "[]" > "$dst"
  # Create JSON objects using string concatenation
  # TODO: remove duplicates (?)
  tmp_json_arr=$(mktemp -p "$TMP_DIR")
  # For every package.json found, parse out repository and license info.
  # jq runs through a package.json to find the 'repository' (or a 'url' child
  # if present) key, and check if a 'licenses' key exists:
  # if so, then use that array (remove 'url', append 'confidence' score); else
  # create a 'licenses' array by checking if there are multiple licenses
  # listed (will be separated by AND/OR); if so, split all licenses by separator
  # and create objects for each; else create a single object with license.
  # Finally, append each new object to the final array.
  while read -r dep; do
    obj=$(jq '{
      project: (.repository | (.url)?//.),
      licenses:
        (if (.licenses)? then
          .licenses | map((del(.url))?
          | .+{confidence: 1})
        else
          .license | if (. | test(" AND | OR ")) then
            . | gsub("\\(|\\)";"")
            | split(" AND | OR "; "g")
            | map({type: ., confidence: 1})
          else
            [{type: ., confidence: 1}]
          end
        end)
    }' < "$dep")
    jq ".+[${obj}]" > "$tmp_json_arr" < "$dst"
    cat "$tmp_json_arr" > "$dst"
  done < "$src"
}
# Clone repository and test if clone was unsuccessful
clone_package() {
  if ! git clone "$1" > /dev/null 2>&1; then
    echo "Failed to clone: $1"
    return 1
  fi
  echo "Cloned: $1"
  return 0
}
# Get all main.go file paths and write to a file
# NOTE: ignores paths that contain *test*, *vendor*, and *workspace*
find_trim_main_paths() {
  find . \
    -name main.go \
    -not -path '*test*' \
    -not -path '*vendor*' \
    -not -path '*workspace*' \
    -printf '%h\n' | uniq > "$1"
  sed -i 's@^\.\(\/\)\?@@g' "$1"
}

# Grab all Golang dependencies and install under TMP_DIR
pushd . > /dev/null
while read -r url; do
  # Continue loop if line is a comment
  if grep '^\s*#.*' <<< "$url" > /dev/null; then continue; fi

  # Make organization dir from SSH repo URL
  # Ex. git@github.com:org/repo.git -> $GOPATH/src/github.com/org
  # shellcheck disable=SC2001
  org_dir=$(echo "$url" | sed 's|.*@\(.*\):\(.*\)\/.*|'"$GOPATH"'\/src\/\1\/\2|g')
  # shellcheck disable=SC2001
  repo=$(echo "$url" | sed 's|.*\/\(.*\)\..*|\1|g')
  if [ ! -d "$org_dir" ]; then mkdir -p "$org_dir"; fi
  cd "$org_dir"

  # Exit loop if package cannot be downloaded
  if ! clone_package "$url"; then exit 1; fi
  # Find all main.go files
  cd "$repo"
  find_trim_main_paths "$tmp_shared_file"

  # Create an organization_repo_license.json file
  subd_path="${org_dir##*/}_${repo}"
  go_pkg_json_file="${TMP_LICENSE_DIR}/${subd_path}_license.json"
  # shellcheck disable=SC2001
  go_pkg_rel_path="$(echo "$org_dir" | sed 's|.*\/src\/\(.*\)|\1|g')/${repo}"
  # Construct full paths to each main.go parent dir, starting after src/
  dep_loc_set=$(sed 's|^\(.*\)$|'"${go_pkg_rel_path}"'\/\1|g' "$tmp_shared_file" | tr '\n' ' ')

  # If no main.go's can be found, 'go list' won't work, so skip
  if [ -z "$dep_loc_set" ]; then echo "No license info: $go_pkg_rel_path"; continue; fi

  # NOTE: license-bill-of-materials will return a non-zero exit code if one or
  # more dependencies' license cannot be found, even if other dependencies
  # in the same package/repo have licenses. Unset 'e' to prevent exit
  set +e
  # shellcheck disable=SC2086
  "${old_gopath}"/bin/license-bill-of-materials $dep_loc_set > "$go_pkg_json_file"
  set -e
done < "$go_pkg_inputs"
popd > /dev/null

# Grab all JS dependency license info from package.json files (usually in
# frontend dirs)
pushd . > /dev/null
while read -r loc; do
  if grep '^\s*#.*' <<< "$loc" > /dev/null; then continue; fi

  # Retrieve each dependency of the current package
  # shellcheck disable=SC2001
  subd_path=$(echo "${loc#*/}" | sed 's@\/@_@g')
  js_pkg_json_file="${TMP_LICENSE_DIR}/${subd_path}_license.json"
  # Install all non-dev deps. Exit if any errors occur while installing

  # Install non-dev node dependencies
  cd "${GOPATH}/src/${loc}"
  if ! yarn install --production --no-lockfiles; then
    echo "Error installing node modules."
    exit 1
  fi

  # Find all package.json files in node_modules
  find node_modules/ -type f -name package.json > "$tmp_shared_file"

  # Parse all package.json files for project and license fields
  parse_js_licenses "$tmp_shared_file" "$js_pkg_json_file"
done < "$js_pkg_inputs"
popd > /dev/null

# Append all licenses into one file
while read -r file; do
  # Concatenate all JSON arrays in a license file (might be 2-3 from
  # license-bill-of-materials output)
  json_license=$(jq '.[]' | jq -s 'map(if (.error)? then empty else . end)' < "$file")

  # Concatenate with final array
  jq ".+${json_license}" > "$tmp_shared_file" < "$FINAL_LICENSE_FILE"
  cat "$tmp_shared_file" > "$FINAL_LICENSE_FILE"
done < <(find "$TMP_LICENSE_DIR" -name '*_license.json')

# Remove duplicates and sort all objects by 'project' key, to avoid too many
# commit line changes
jq 'unique | sort_by(.project)' > "$tmp_shared_file" < "$FINAL_LICENSE_FILE"
cat "$tmp_shared_file" > "$FINAL_LICENSE_FILE"

# TODO: JSON output to Tectonic frontend (REST?)
echo "Aggregated licenses: $(readlink -f "$FINAL_LICENSE_FILE")"

clean_up
