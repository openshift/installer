/*
Copyright (c) 2025 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ocm

import (
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

// CreateImageMirror creates a new image mirror for the given cluster
func (c *Client) CreateImageMirror(clusterID, mirrorType, source string, mirrors []string) (*cmv1.ImageMirror, error) {
	imageMirror, err := cmv1.NewImageMirror().
		Type(mirrorType).
		Source(source).
		Mirrors(mirrors...).
		Build()
	if err != nil {
		return nil, err
	}

	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().
		Cluster(clusterID).
		ImageMirrors().
		Add().
		Body(imageMirror).
		Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}

	return response.Body(), nil
}

// ListImageMirrors lists image mirrors for the given cluster
func (c *Client) ListImageMirrors(clusterID string) ([]*cmv1.ImageMirror, error) {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().
		Cluster(clusterID).
		ImageMirrors().
		List().
		Page(1).
		Size(-1).
		Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}

	return response.Items().Slice(), nil
}

// GetImageMirror gets a specific image mirror by ID for the given cluster
func (c *Client) GetImageMirror(clusterID, id string) (*cmv1.ImageMirror, error) {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().
		Cluster(clusterID).
		ImageMirrors().
		ImageMirror(id).
		Get().
		Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}

	return response.Body(), nil
}

// UpdateImageMirror updates an existing image mirror
func (c *Client) UpdateImageMirror(clusterID, id string, mirrors []string, mirrorType *string) (*cmv1.ImageMirror, error) {
	imageMirrorBuilder := cmv1.NewImageMirror().
		Mirrors(mirrors...)

	if mirrorType != nil {
		imageMirrorBuilder = imageMirrorBuilder.Type(*mirrorType)
	}

	imageMirrorPatch, err := imageMirrorBuilder.Build()
	if err != nil {
		return nil, err
	}

	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().
		Cluster(clusterID).
		ImageMirrors().
		ImageMirror(id).
		Update().
		Body(imageMirrorPatch).
		Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}

	return response.Body(), nil
}

// DeleteImageMirror deletes an image mirror
func (c *Client) DeleteImageMirror(clusterID, id string) error {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().
		Cluster(clusterID).
		ImageMirrors().
		ImageMirror(id).
		Delete().
		Send()
	if err != nil {
		return handleErr(response.Error(), err)
	}

	return nil
}
