package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GpuResource struct {
	Count int64  `json:"count"`
	Sku   GpuSku `json:"sku"`
}
