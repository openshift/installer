// Copyright 2018 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package errors includes errors that are used in multiple config versions
package errors

import (
	"errors"
	"fmt"
)

var (
	// Parsing / general errors
	ErrInvalid = errors.New("config is not valid")
	ErrEmpty   = errors.New("not a config (empty)")

	// Ignition section errors
	ErrInvalidVersion = errors.New("invalid config version (couldn't parse)")
	ErrUnknownVersion = errors.New("unsupported config version")

	ErrDeprecated         = errors.New("config format deprecated")
	ErrCompressionInvalid = errors.New("invalid compression method")

	// Storage section errors
	ErrFilePermissionsUnset      = errors.New("permissions unset, defaulting to 0644")
	ErrDirectoryPermissionsUnset = errors.New("permissions unset, defaulting to 0755")
	ErrDiskDeviceRequired        = errors.New("disk device is required")
	ErrPartitionNumbersCollide   = errors.New("partition numbers collide")
	ErrPartitionsOverlap         = errors.New("partitions overlap")
	ErrPartitionsMisaligned      = errors.New("partitions misaligned")
	ErrAppendAndOverwrite        = errors.New("cannot set both append and overwrite to true")
	ErrFilesystemInvalidFormat   = errors.New("invalid filesystem format")
	ErrLabelNeedsFormat          = errors.New("filesystem must specify format if label is specified")
	ErrFormatNilWithOthers       = errors.New("format cannot be empty when path, label, uuid, or options are specified")
	ErrExt4LabelTooLong          = errors.New("filesystem labels cannot be longer than 16 characters when using ext4")
	ErrBtrfsLabelTooLong         = errors.New("filesystem labels cannot be longer than 256 characters when using btrfs")
	ErrXfsLabelTooLong           = errors.New("filesystem labels cannot be longer than 12 characters when using xfs")
	ErrSwapLabelTooLong          = errors.New("filesystem labels cannot be longer than 15 characters when using swap")
	ErrVfatLabelTooLong          = errors.New("filesystem labels cannot be longer than 11 characters when using vfat")
	ErrFileIllegalMode           = errors.New("illegal file mode")
	ErrNoFilesystem              = errors.New("no filesystem specified")
	ErrBothIDAndNameSet          = errors.New("cannot set both id and name")
	ErrLabelTooLong              = errors.New("partition labels may not exceed 36 characters")
	ErrDoesntMatchGUIDRegex      = errors.New("doesn't match the form \"01234567-89AB-CDEF-EDCB-A98765432101\"")
	ErrLabelContainsColon        = errors.New("partition label will be truncated to text before the colon")
	ErrNoPath                    = errors.New("path not specified")
	ErrPathRelative              = errors.New("path not absolute")
	ErrDirtyPath                 = errors.New("path is not fully simplified")
	ErrSparesUnsupportedForLevel = errors.New("spares unsupported for arrays with a level greater than 0")
	ErrUnrecognizedRaidLevel     = errors.New("unrecognized raid level")
	ErrShouldNotExistWithOthers  = errors.New("shouldExist specified false with other options also specified")
	ErrZeroesWithShouldNotExist  = errors.New("shouldExist is false for a partition and other partition(s) has start or size 0")
	ErrPartitionsUnitsMismatch   = errors.New("cannot mix MBs and sectors within a disk")
	ErrSizeDeprecated            = errors.New("size is deprecated; use sizeMB instead")
	ErrStartDeprecated           = errors.New("start is deprecated; use startMB instead")
	ErrNeedLabelOrNumber         = errors.New("a partition number >= 1 or a label must be specified")
	ErrDuplicateLabels           = errors.New("cannot use the same partition label twice")

	// Passwd section errors
	ErrPasswdCreateAndGecos        = errors.New("cannot use both the create object and the user-level gecos field")
	ErrPasswdCreateAndGroups       = errors.New("cannot use both the create object and the user-level groups field")
	ErrPasswdCreateAndHomeDir      = errors.New("cannot use both the create object and the user-level homeDir field")
	ErrPasswdCreateAndNoCreateHome = errors.New("cannot use both the create object and the user-level noCreateHome field")
	ErrPasswdCreateAndNoLogInit    = errors.New("cannot use both the create object and the user-level noLogInit field")
	ErrPasswdCreateAndNoUserGroup  = errors.New("cannot use both the create object and the user-level noUserGroup field")
	ErrPasswdCreateAndPrimaryGroup = errors.New("cannot use both the create object and the user-level primaryGroup field")
	ErrPasswdCreateAndShell        = errors.New("cannot use both the create object and the user-level shell field")
	ErrPasswdCreateAndSystem       = errors.New("cannot use both the create object and the user-level system field")
	ErrPasswdCreateAndUID          = errors.New("cannot use both the create object and the user-level uid field")

	// Systemd section errors
	ErrInvalidSystemdExt       = errors.New("invalid systemd unit extension")
	ErrInvalidSystemdDropinExt = errors.New("invalid systemd drop-in extension")

	// Misc errors
	ErrInvalidScheme       = errors.New("invalid url scheme")
	ErrInvalidUrl          = errors.New("unable to parse url")
	ErrHashMalformed       = errors.New("malformed hash specifier")
	ErrHashWrongSize       = errors.New("incorrect size for hash sum")
	ErrHashUnrecognized    = errors.New("unrecognized hash function")
	ErrEngineConfiguration = errors.New("engine incorrectly configured")

	// AWS S3 specific errors
	ErrInvalidS3ObjectVersionId = errors.New("invalid S3 object VersionId")
)

// NewNoInstallSectionError produces an error indicating the given unit, named
// name, is missing an Install section.
func NewNoInstallSectionError(name string) error {
	return fmt.Errorf("unit %q is enabled, but has no install section so enable does nothing", name)
}
