package kms

import "fmt"

// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

func wrapError(err error, msg string) error {
	return fmt.Errorf("[ERROR] %s: %s", msg, err)
}
