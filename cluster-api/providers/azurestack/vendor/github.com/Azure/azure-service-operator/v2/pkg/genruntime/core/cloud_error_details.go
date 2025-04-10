/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package core

type CloudErrorDetails struct {
	Classification ErrorClassification
	Code           string
	Message        string
}
