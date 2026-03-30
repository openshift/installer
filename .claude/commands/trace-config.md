---
description: Trace how an install config field flows through validation, defaults, and manifests
argument-hint: "[field-name]"
---

## Name
trace-config

## Synopsis
```
/trace-config [field-name]
/trace-config pullSecret
/trace-config aws.region
/trace-config platform.aws.serviceEndpoints
```

## Description
The `trace-config` command traces an install config field through the entire installer pipeline to help you understand:
- The field's schema definition (from the CRD that powers `openshift-install explain`)
- Where the field is defined in Go types
- What validation logic applies to it
- Where default values are set
- How it's consumed by assets (manifests, ignition configs, Terraform/CAPI resources)
- Where it appears in the final cluster configuration

This is particularly useful when:
- Debugging validation errors
- Understanding field dependencies
- Modifying existing fields
- Adding new fields to the install config
- Troubleshooting why a field value isn't appearing in the cluster

The command performs comprehensive analysis:
1. **Schema Lookup** - Extracts field definition from `data/data/install.openshift.io_installconfigs.yaml` (the CRD schema)
2. **Type Definition Discovery** - Finds the struct field in `pkg/types/`
3. **Validation Tracing** - Locates validation functions in `pkg/types/*/validation/`
4. **Default Value Location** - Finds default assignment in `pkg/types/*/defaults/`
5. **Asset Usage Analysis** - Searches asset generation code for field references
6. **Manifest Generation** - Traces how the field flows into ignition configs, Terraform, or CAPI resources

## Implementation

The implementation follows these steps to provide comprehensive field analysis:

### Step 1: Parse the Field Name

Parse dot-notation field names (e.g., `platform.aws.region` or `pullSecret`):

```bash
field_name="$1"

if [ -z "$field_name" ]; then
  echo "Error: Field name is required"
  echo "Usage: /trace-config [field-name]"
  echo ""
  echo "Examples:"
  echo "  /trace-config pullSecret"
  echo "  /trace-config aws.region"
  echo "  /trace-config platform.aws.serviceEndpoints"
  exit 1
fi

# Parse field path (e.g., "platform.aws.region" -> platform="aws", field="region")
# For top-level fields like "pullSecret", platform will be empty
IFS='.' read -ra parts <<< "$field_name"
```

### Step 2: Look Up Field in CRD Schema

Extract the field definition from the install config CRD schema (what `openshift-install explain` uses):

```bash
echo "## Schema Definition"
echo ""
echo "From \`data/data/install.openshift.io_installconfigs.yaml\` (powers \`openshift-install explain\`):"
echo ""

# Convert field name to schema path (e.g., "aws.region" -> ".spec.versions[*].schema.openAPIV3Schema.properties.aws.properties.region")
# This is a simplified search - just look for the field in the YAML
schema_file="data/data/install.openshift.io_installconfigs.yaml"
if [ -f "$schema_file" ]; then
  # Search for the field in the schema and extract its definition
  # Look for the field name in the properties section
  field_to_find="${parts[-1]}"

  # Use yq or grep to extract the field definition
  # For now, use grep to find the field and show context
  grep -A 10 "^\s*${field_to_find}:" "$schema_file" | head -20
else
  echo "Schema file not found. Skipping schema lookup."
fi

echo ""
```

### Step 3: Find Type Definition

Search for the Go struct field definition in `pkg/types/`:

```bash
# Search for field definition in types
# Use grep with word boundaries to avoid partial matches
echo "## Type Definition"
echo ""

# Try multiple search patterns to find the field
field_pattern="${parts[-1]}"  # Last part is the field name

# Search in pkg/types for struct field definitions
grep -r "^\s*${field_pattern}\s\+" pkg/types/ --include="*.go" \
  -B 2 -A 1 | grep -v "vendor/" | head -20
```

### Step 4: Find Validation Logic

Locate validation functions that check this field:

```bash
echo ""
echo "## Validation"
echo ""

# Search for validation functions
# Look for the field name in validation code
grep -r "${field_pattern}" pkg/types/*/validation/ --include="*.go" \
  -B 3 -A 10 | grep -v "vendor/" | head -50

# Also check for validation in installconfig
grep -r "${field_pattern}" pkg/asset/installconfig/ --include="*.go" \
  -B 3 -A 10 | grep -v "vendor/" | head -50
```

### Step 5: Find Default Values

Search for default value assignment:

```bash
echo ""
echo "## Default Values"
echo ""

# Search in defaults packages
grep -r "${field_pattern}" pkg/types/*/defaults/ --include="*.go" \
  -B 3 -A 10 | grep -v "vendor/"

# Also check for defaults in platform-specific code
if [ ${#parts[@]} -gt 1 ]; then
  platform="${parts[1]}"
  grep -r "${field_pattern}" "pkg/types/${platform}/defaults/" --include="*.go" \
    -B 3 -A 10 2>/dev/null
fi
```

### Step 6: Trace Asset Usage

Find where the field is consumed by asset generation:

```bash
echo ""
echo "## Asset Usage"
echo ""

# Search across asset packages
grep -r "${field_pattern}" pkg/asset/ --include="*.go" \
  -B 2 -A 5 | grep -v "vendor/" | head -100

# Check specific asset areas
for dir in pkg/asset/manifests pkg/asset/ignition pkg/asset/machines pkg/tfvars; do
  if results=$(grep -r "${field_pattern}" "$dir" --include="*.go" -l 2>/dev/null); then
    echo "Used in $dir:"
    echo "$results" | sed 's/^/  /'
  fi
done
```

### Step 7: Check Template Usage

Search for field usage in data templates:

```bash
echo ""
echo "## Template Usage"
echo ""

# Search in data/data templates
grep -r "${field_pattern}" data/data/ -B 2 -A 2 | head -50

# Check for Go template references (e.g., .Config.FieldName)
field_camel=$(echo "$field_pattern" | sed 's/_\([a-z]\)/\U\1/g' | sed 's/^./\U&/')
grep -r "\\.${field_camel}" data/data/ -B 2 -A 2 2>/dev/null
```

### Step 8: Find Cluster API Usage

For fields used in CAPI providers:

```bash
echo ""
echo "## Cluster API Usage"
echo ""

# If platform is specified, check CAPI provider
if [ ${#parts[@]} -gt 1 ]; then
  platform="${parts[1]}"
  if [ -d "cluster-api/providers/${platform}" ]; then
    grep -r "${field_pattern}" "cluster-api/providers/${platform}/" \
      --include="*.go" -B 2 -A 5 | head -50
  fi
fi
```

### Step 9: Generate Summary

Provide a summary of findings:

```bash
echo ""
echo "## Summary"
echo ""
echo "Field: ${field_name}"
echo ""
echo "Flow:"
echo "  1. Schema: data/data/install.openshift.io_installconfigs.yaml (CRD definition)"
echo "  2. Type: pkg/types/[...] (Go struct definition)"
echo "  3. Validation: pkg/types/*/validation/ (validation rules)"
echo "  4. Defaults: pkg/types/*/defaults/ (default values)"
echo "  5. Usage: pkg/asset/ (asset generation)"
echo "  6. Output: Templates, Ignition configs, or CAPI resources"
echo ""
echo "💡 Tip: Use \`openshift-install explain installconfig.${field_name}\` to see user-facing docs"
echo ""
echo "Use the file paths above to explore the complete implementation."
```

## Examples

1. **Trace Top-Level Field**

```
/trace-config pullSecret
```

Output shows:
- Schema definition from the CRD (type, description, validation rules)
- Type definition in `pkg/types/installconfig.go`
- Validation in `pkg/asset/installconfig/installconfig.go`
- Usage in ignition config generation
- References in manifest templates

2. **Trace Platform-Specific Field**

```
/trace-config aws.region
```

Output shows:
- Schema definition with allowed values and description
- Type definition in `pkg/types/aws/platform.go`
- Validation in `pkg/types/aws/validation/platform.go`
- Default values in `pkg/types/aws/defaults/platform.go`
- Usage in AWS asset generation (`pkg/asset/machines/aws/`)
- References in Terraform or CAPI resources

3. **Trace Nested Field**

```
/trace-config platform.aws.serviceEndpoints
```

Output shows complete path through:
- Go struct path: `InstallConfig.Platform.AWS.ServiceEndpoints`
- Validation logic for endpoint format
- Usage in AWS client configuration
- How it affects cluster networking

## Arguments
- `field-name`: The install config field to trace using YAML notation. Can be simple (e.g., `pullSecret`) or dot-notation for nested fields (e.g., `aws.region`, `platform.aws.serviceEndpoints`). The command will map this to the Go struct path (e.g., `platform.aws.serviceEndpoints` → `InstallConfig.Platform.AWS.ServiceEndpoints`).

## Use Cases

**Debugging Validation Errors:**
When you see an error like "invalid value for `aws.region`", use this command to find the validation logic and understand what values are allowed.

**Understanding Field Flow:**
Before modifying a field, trace it to see all the places where changes might be needed.

**Adding New Fields:**
Use this to understand the pattern for similar existing fields when adding new configuration options.

**Troubleshooting Missing Values:**
If a field value isn't appearing in the cluster, trace it to see if it's being consumed correctly by assets.

## Tips

- For nested fields, use dot notation matching the YAML structure
- Field names are case-sensitive and should match Go field names
- Results are ordered by relevance: type → validation → defaults → usage
- Check both the output and the file paths for complete understanding

## Troubleshooting

**No results found:**
- Verify the field name spelling (case-sensitive)
- Check if the field might use a different name in code vs. YAML
- Some fields may be computed rather than directly specified

**Too many results:**
- Use a more specific field name
- Look for the most recent/relevant file paths
- Focus on files in `pkg/types/` first, then `pkg/asset/`
