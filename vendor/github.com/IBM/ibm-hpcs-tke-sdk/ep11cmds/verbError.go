//
// Copyright 2021 IBM Inc. All rights reserved
// SPDX-License-Identifier: Apache2.0
//

package ep11cmds

import ()

/* Defines an error type with an associated return code and reason code */

type VerbError interface {
	error
	ReturnCode() int
	ReasonCode() int
}

type VerbErrorStruct struct {
	errorString string
	returnCode  int
	reasonCode  int
}

func (ves *VerbErrorStruct) Error() string {
	return ves.errorString
}

func (ves *VerbErrorStruct) ReturnCode() int {
	return ves.returnCode
}

func (ves *VerbErrorStruct) ReasonCode() int {
	return ves.reasonCode
}

func NewVerbError(text string, returnCode int, reasonCode int) VerbError {
	return &VerbErrorStruct{text, returnCode, reasonCode}
}
