package kms

import "github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"

// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

func wrapError(err error, msg string) error {
	return flex.FmtErrorf("[ERROR] %s: %s", msg, err)
}
