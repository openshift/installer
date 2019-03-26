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
	"fmt"
	"regexp"
	"strings"

	"github.com/coreos/ignition/config/shared/errors"
	"github.com/coreos/ignition/config/validate/report"
)

const (
	guidRegexStr = "^(|[[:xdigit:]]{8}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{12})$"
)

func (p Partition) Key() string {
	if p.Number != 0 {
		return fmt.Sprintf("number:%d", p.Number)
	} else {
		return fmt.Sprintf("label:%s", *p.Label)
	}
}

func (p Partition) Validate() (r report.Report) {
	if (p.Start != nil || p.Size != nil) && (p.StartMiB != nil || p.SizeMiB != nil) {
		r.AddOnError(errors.ErrPartitionsUnitsMismatch)
	}
	if p.ShouldExist != nil && !*p.ShouldExist &&
		(p.Label != nil || (p.TypeGUID != nil && *p.TypeGUID != "") || (p.GUID != nil && *p.GUID != "") || p.Start != nil || p.Size != nil) {
		r.AddOnError(errors.ErrShouldNotExistWithOthers)
	}
	if p.Number == 0 && p.Label == nil {
		r.AddOnError(errors.ErrNeedLabelOrNumber)
	}
	return
}

func (p Partition) ValidateSize() (r report.Report) {
	if p.Size != nil {
		r.AddOnDeprecated(errors.ErrSizeDeprecated)
	}
	return
}

func (p Partition) ValidateStart() (r report.Report) {
	if p.Start != nil {
		r.AddOnDeprecated(errors.ErrStartDeprecated)
	}
	return
}

func (p Partition) ValidateLabel() (r report.Report) {
	if p.Label == nil {
		return
	}
	// http://en.wikipedia.org/wiki/GUID_Partition_Table#Partition_entries:
	// 56 (0x38) 	72 bytes 	Partition name (36 UTF-16LE code units)

	// XXX(vc): note GPT calls it a name, we're using label for consistency
	// with udev naming /dev/disk/by-partlabel/*.
	if len(*p.Label) > 36 {
		r.AddOnError(errors.ErrLabelTooLong)
	}

	// sgdisk uses colons for delimitting compound arguments and does not allow escaping them.
	if strings.Contains(*p.Label, ":") {
		r.AddOnError(errors.ErrLabelContainsColon)
	}
	return
}

func (p Partition) ValidateTypeGUID() report.Report {
	return validateGUID(p.TypeGUID)
}

func (p Partition) ValidateGUID() report.Report {
	return validateGUID(p.GUID)
}

func validateGUID(guidPointer *string) (r report.Report) {
	if guidPointer == nil {
		return
	}
	guid := *guidPointer
	ok, err := regexp.MatchString(guidRegexStr, guid)
	if err != nil {
		r.AddOnError(fmt.Errorf("error matching guid regexp: %v", err))
	} else if !ok {
		r.AddOnError(errors.ErrDoesntMatchGUIDRegex)
	}
	return r
}
