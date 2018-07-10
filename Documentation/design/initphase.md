# Init Phase

## Goals

1. Generates `InstallConfig` by prompting user for required parameters, and populating the rest of configuration with default values.
2. Make useful decisions on default values for most options.
3. Init phase should never prompt the user for more than 10 options unless there is a very strong reason.

## Non Goals

1. Supporting all the configuration option through init is not a goal. Init is just a helper for humans.

## Overview

The output of the `init` command is a config file that describes a default cluster for the given provider, and is valid/immediately usable.

The `init` step is not required, this is merely a helper to generate a valid `InstallConfig`. A user could use examples from documentation, roll their own by hand, or re-use existing configs from source control.

## Detailed Design

### Prompts on CLI vs flags vs Environment variables

1. Flags don’t work well here because we want the fast-path to be single command.
2. Init is only valuable for humans and environment variables are similarly tedious.

Init phase should use prompts on CLI for top-level configurations that are required to generate `InstallConfig`.

### Remember user's answers

In case where users runs `installer init` and stops after platform prompt for some reason. Should the init command ask for platform again when it is rerun.

The number of prompts from users need to be very few and remembering user's answer is unnecessary optimization.

### Idempotency

Re-running the `installer init` prompts the user for options again and writes the `install-config.yaml` to disk, overwriting the old `install-config.yaml` file in the directory.

### Output

Init phase generates `install-config.yaml` file that contains `InstallConfig` object in the current working directory or directory specified by `--asset-dir` flag.

## TODOs

1. Add examples for user prompts for each platform and the corresponding `install-config.yaml`.
