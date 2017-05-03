# Conventions

This document contains conventions for terraform code and module layout. All new platforms should adhere to these guidelines

## Directory Structure

### Root

```
├── config.tf
├── /modules
└── /platforms
```

- `config.tf`: Contains all common variable declarations which are required for all platforms.
- `/modules`: Contains all reusable modules. Platform specific modules go under their respective platform subdirectory.
- `/platforms`: All platform specific entrypoints

### Platforms

Each platform entrypoint should be laid out under the `platforms` top-level directory.
- `/platforms/<PLATFORM_NAME>`

Each platform should at minimum contain the following files:
- `config.tf`: A symlink to the main config.tf file in the repo root.
- `variables.tf`: Platform specific variable declararions. All should be namespaced as `tectonic_<PLATFORM_NAME>_<VAR_NAME>`
- `terraform.tfvars.example`: An example file contained in the root of the repo which demonstrates how users can configure the variables for provisioning. Users must customize this file.
- `README.md`: All platform related documentation to quickly get started and references to other docs.
- `main.tf`: The main entrypoint. This should configure all modules and pass along common variables. Users may customize this file as needed.
- `output.tf`: Contains all output variable declarations. These contain any generated output useful to users.

### Modules

Platform-agnostic modules are located in the root of the modules directory.

```
/modules/<COMMON_MODULE>
```

Platform-specific modules should go into their own subdirectory.

```
/modules/<COMMON_MODULE>/<PLATFORM_NAME>/<PLATFORM_MODULE>/
```
Each module should contain the following files:
- `output.tf`: Contains all output variable declarations. These contain any generated output useful to module consumers.
- `README.md`: All module related documentation.
- `variables.tf`: Module-specific variables. Should never contain defaults.

## Documentation

For more detailed documentation add docs to:

```
/Documentation/modules/<MODULE_NAME>
```
or
```
/Documentation/platforms/<PLATFORM_NAME>
```

