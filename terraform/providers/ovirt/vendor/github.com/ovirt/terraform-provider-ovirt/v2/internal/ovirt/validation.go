package ovirt

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ovirtclient "github.com/ovirt/go-ovirt-client/v3"
)

func validateDiskSize(i interface{}, path cty.Path) diag.Diagnostics {
	size, ok := i.(int)
	if !ok {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Disk size must be an integer.",
				Detail:        "The provided disk size is not an integer.",
				AttributePath: path,
			},
		}
	}
	if size <= 0 {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Disk size must be a positive integer.",
				Detail:   fmt.Sprintf("The provided disk size must be a positive integer, got %d.", size),
			},
		}
	}
	return nil
}

func validateLocalFile(i interface{}, p cty.Path) diag.Diagnostics {
	val, ok := i.(string)
	if !ok {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Local file path must be a string.",
				Detail:        "The local file path is not a string.",
				AttributePath: p,
			},
		}
	}
	cwd, err := os.Getwd()
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Failed to get determine current working directory",
				Detail:        err.Error(),
				AttributePath: p,
			},
		}
	}
	path, err := filepath.Abs(val)
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf("Failed to get absolute path for %s, cwd is %s", val, cwd),
				Detail:        err.Error(),
				AttributePath: p,
			},
		}
	}
	if _, err := os.Stat(path); err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf("Failed to stat %s", val),
				Detail:        err.Error(),
				AttributePath: p,
			},
		}
	}
	return nil
}

func validateFormat(i interface{}, path cty.Path) diag.Diagnostics {
	val, ok := i.(string)
	if !ok {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Disk format must be a string.",
				Detail:        "The provided disk format is not a string.",
				AttributePath: path,
			},
		}
	}
	format := ovirtclient.ImageFormat(val)
	if err := format.Validate(); err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Invalid disk image format.",
				Detail:        err.Error(),
				AttributePath: path,
			},
		}
	}
	return nil
}

func validateTLSSystem(value interface{}, path cty.Path) diag.Diagnostics {
	v, ok := value.(bool)
	if !ok {
		return diag.Diagnostics{
			{
				Severity:      diag.Error,
				Summary:       "Passed parameter is not a bool.",
				Detail:        "The passed parameter for the system cert pool is not a bool.",
				AttributePath: path,
			},
		}
	}

	if v && version.Must(
		version.NewVersion(
			strings.ReplaceAll(
				runtime.Version(),
				"go",
				"",
			),
		),
	).LessThan(version.Must(version.NewVersion("1.18.0"))) {
		if runtime.GOOS == "windows" {
			return diag.Diagnostics{
				{
					Severity:      diag.Error,
					Summary:       "The tls_ca_system option not available on Windows.",
					Detail:        "The tls_ca_system option is not available on Windows when the Terraform provider is compiled with Go versions before 1.18 due to Golang bug 16736.",
					AttributePath: path,
				},
			}
		}
		return diag.Diagnostics{
			{
				Severity:      diag.Warning,
				Summary:       "You are using the tls_ca_system option. Your Terraform code will not work on Windows.",
				Detail:        "The tls_ca_system option is not available on Windows when the Terraform provider is compiled with Go versions before 1.18 due to Golang bug 16736.",
				AttributePath: path,
			},
		}
	}

	return nil
}

func validateTLSInsecure(value interface{}, path cty.Path) diag.Diagnostics {
	v, ok := value.(bool)
	if !ok {
		return diag.Diagnostics{
			{
				Severity:      diag.Error,
				Summary:       "Passed parameter is not a bool.",
				Detail:        "The passed parameter for the insecure flag is not a bool.",
				AttributePath: path,
			},
		}
	}

	if v {
		return diag.Diagnostics{
			{
				Severity:      diag.Warning,
				Summary:       "Insecure connection mode enabled.",
				Detail:        "The insecure connection mode to oVirt is enabled. This is a very bad idea because Terraform will not validate the certificate of the oVirt engine.",
				AttributePath: path,
			},
		}
	}
	return nil
}

func validateNonEmpty(value interface{}, path cty.Path) diag.Diagnostics {
	v, ok := value.(string)
	if !ok {
		return diag.Diagnostics{
			{
				Severity:      diag.Error,
				Summary:       "Passed parameter is not a string.",
				Detail:        "The passed parameter is not a string.",
				AttributePath: path,
			},
		}
	}

	if v == "" {
		return diag.Diagnostics{
			{
				Severity:      diag.Error,
				Summary:       "This field must not be empty.",
				AttributePath: path,
			},
		}
	}
	return nil
}

func validateDiskInterface(i interface{}, path cty.Path) diag.Diagnostics {
	val, ok := i.(string)
	if !ok {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "The disk_interface should be a string.",
				Detail:        "The provided disk_interface value is not a string.",
				AttributePath: path,
			},
		}
	}
	interf := ovirtclient.DiskInterface(val)
	if err := interf.Validate(); err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Invalid disk_interface value.",
				Detail:        err.Error(),
				AttributePath: path,
			},
		}
	}
	return nil
}

var uuidRegexp = regexp.MustCompile(`^\b[0-9a-f]{8}\b-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-\b[0-9a-f]{12}\b$`)

func validateUUID(i interface{}, path cty.Path) diag.Diagnostics {
	val, ok := i.(string)
	if !ok {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Not a string",
				Detail:        "The specified value is not a string, but must be a string containing a UUID.",
				AttributePath: path,
			},
		}
	}

	if !uuidRegexp.MatchString(val) {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Not a UUID",
				Detail:        "The specified value is not a UUID.",
				AttributePath: path,
			},
		}
	}
	return nil
}

func validatePositiveInt(i interface{}, path cty.Path) diag.Diagnostics {
	val, ok := i.(int)
	if !ok {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Not an integer",
				Detail:        "The specified value is not an integer.",
				AttributePath: path,
			},
		}
	}
	if val < 1 {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Not a positive integer",
				Detail:        "The specified value is not a positive integer.",
				AttributePath: path,
			},
		}
	}
	return nil
}

func validateEnum(values []string) schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		val, ok := i.(string)
		if !ok {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity:      diag.Error,
					Summary:       "Not a string",
					Detail:        "The specified value is not a string.",
					AttributePath: path,
				},
			}
		}
		for _, value := range values {
			if value == val {
				return nil
			}
		}
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Not a valid value",
				Detail:        "The specified value is not one of the required values.",
				AttributePath: path,
			},
		}
	}
}

func validateHugePages(i interface{}, path cty.Path) diag.Diagnostics {
	hugePages := i.(int)
	if hugePages < 0 {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Not a valid huge page value.",
				Detail:        "must not be negative",
				AttributePath: path,
			},
		}
	}
	err := ovirtclient.VMHugePages(hugePages).Validate()
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Not a valid huge page value.",
				Detail:        err.Error(),
				AttributePath: path,
			},
		}
	}
	return nil
}
