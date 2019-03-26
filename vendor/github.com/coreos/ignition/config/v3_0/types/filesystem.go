// Copyright 2016 CoreOS, Inc.
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

package types

import (
	"github.com/coreos/ignition/config/shared/errors"
	"github.com/coreos/ignition/config/validate/report"
)

func (f Filesystem) Key() string {
	return f.Device
}

func (f Filesystem) IgnoreDuplicates() map[string]struct{} {
	return map[string]struct{}{
		"Options": {},
	}
}

func (f Filesystem) ValidatePath() (r report.Report) {
	if f.Path == nil || *f.Path == "" {
		return
	}
	r.AddOnError(validatePath(*f.Path))
	return
}

func (f Filesystem) ValidateDevice() (r report.Report) {
	r.AddOnError(validatePath(f.Device))
	return
}

func (f Filesystem) ValidateFormat() (r report.Report) {
	if f.Format == nil || *f.Format == "" {
		if (f.Path == nil || *f.Path == "") &&
			(f.Label == nil || *f.Label == "") &&
			(f.UUID == nil || *f.UUID == "") &&
			len(f.Options) == 0 {
			return
		}
		r.AddOnError(errors.ErrFormatNilWithOthers)
		return
	}
	switch *f.Format {
	case "ext4", "btrfs", "xfs", "swap", "vfat":
	default:
		r.AddOnError(errors.ErrFilesystemInvalidFormat)
	}
	return
}

func (f Filesystem) ValidateLabel() (r report.Report) {
	if f.Label == nil || *f.Label == "" {
		return
	}
	if f.Format == nil || *f.Format == "" {
		r.AddOnError(errors.ErrLabelNeedsFormat)
		return
	}
	switch *f.Format {
	case "ext4":
		if len(*f.Label) > 16 {
			// source: man mkfs.ext4
			r.AddOnError(errors.ErrExt4LabelTooLong)
		}
	case "btrfs":
		if len(*f.Label) > 256 {
			// source: man mkfs.btrfs
			r.AddOnError(errors.ErrBtrfsLabelTooLong)
		}
	case "xfs":
		if len(*f.Label) > 12 {
			// source: man mkfs.xfs
			r.AddOnError(errors.ErrXfsLabelTooLong)
		}
	case "swap":
		// mkswap's man page does not state a limit on label size, but through
		// experimentation it appears that mkswap will truncate long labels to
		// 15 characters, so let's enforce that.
		if len(*f.Label) > 15 {
			r.AddOnError(errors.ErrSwapLabelTooLong)
		}
	case "vfat":
		if len(*f.Label) > 11 {
			// source: man mkfs.fat
			r.AddOnError(errors.ErrVfatLabelTooLong)
		}
	}
	return
}
