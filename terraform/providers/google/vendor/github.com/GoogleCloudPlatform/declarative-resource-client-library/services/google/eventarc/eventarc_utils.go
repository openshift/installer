// Copyright 2024 Google LLC. All Rights Reserved.
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
// Package eventarc contains DCL for Eventarc.
package eventarc

import (
	"context"
)

// Custom delete method for GoogleChannelConfig. Given that GoogleChannelConfig is a singleton resource, Eventarc's API only provides GET & PATCH requests for this resource
// However, the Terraform team requires a delete endpoint for resources, hence we are making the delete operation do nothing here and adding a custom method
// as a trait within the delete field in google_channel_config.textproto
func (op *deleteGoogleChannelConfigOperation) do(ctx context.Context, r *GoogleChannelConfig, c *Client) error {
	return nil
}
