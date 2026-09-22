/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package entra

import (
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"k8s.io/apimachinery/pkg/types"
)

type Connection interface {
	Client() *msgraphsdk.GraphServiceClient
	CredentialFrom() types.NamespacedName
}
