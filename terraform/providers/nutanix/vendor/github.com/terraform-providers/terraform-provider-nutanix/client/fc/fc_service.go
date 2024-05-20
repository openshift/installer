package foundationcentral

import (
	"context"
	"fmt"
	"net/http"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
)

// Operations implements Service interface
type Operations struct {
	client *client.Client
}

// Interface for foundation central apis
type Service interface {
	GetImagedNode(context.Context, string) (*ImagedNodeDetails, error)
	ListImagedNodes(context.Context, *ImagedNodesListInput) (*ImagedNodesListResponse, error)
	GetImagedCluster(context.Context, string) (*ImagedClusterDetails, error)
	ListImagedClusters(context.Context, *ImagedClustersListInput) (*ImagedClustersListResponse, error)
	CreateCluster(context.Context, *CreateClusterInput) (*CreateClusterResponse, error)
	UpdateCluster(context.Context, string, *UpdateClusterData) error
	DeleteCluster(context.Context, string) error
	CreateAPIKey(context.Context, *CreateAPIKeysInput) (*CreateAPIKeysResponse, error)
	GetAPIKey(context.Context, string) (*CreateAPIKeysResponse, error)
	ListAPIKeys(context.Context, *ListMetadataInput) (*ListAPIKeysResponse, error)
}

/*GetImagedNode Get the details of an imaged node
 * This operation fetches the node details of a node given it's node uuid
 *
 * @param nodeUUID The uuid of the node.
 * @return *ImagedNodeDetails
 */
func (op Operations) GetImagedNode(ctx context.Context, nodeUUID string) (*ImagedNodeDetails, error) {
	path := fmt.Sprintf("/imaged_nodes/%s", nodeUUID)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	imagedNodeDetails := new(ImagedNodeDetails)

	if err != nil {
		return nil, err
	}

	return imagedNodeDetails, op.client.Do(ctx, req, imagedNodeDetails)
}

/*ListImagedNodes Get all the imaged node within FC
 * This operation fetches all the imaged nodes within FC
 *
 * @param input
 * @return *ImagedNodesListResponse
 */
func (op Operations) ListImagedNodes(ctx context.Context, input *ImagedNodesListInput) (*ImagedNodesListResponse, error) {
	path := "/imaged_nodes/list"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, input)

	imagedNodesListResponse := new(ImagedNodesListResponse)

	if err != nil {
		return nil, err
	}

	return imagedNodesListResponse, op.client.Do(ctx, req, imagedNodesListResponse)
}

/*GetImagedCluster Get the details of an imaged cluster
 * This operation fetches the details of a cluster given it's uuid
 *
 * @param clusterUUID uuid of the imaged cluster
 * @return *ImagedClusterDetails
 */
func (op Operations) GetImagedCluster(ctx context.Context, clusterUUID string) (*ImagedClusterDetails, error) {
	path := fmt.Sprintf("/imaged_clusters/%s", clusterUUID)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	imagedClusterDetails := new(ImagedClusterDetails)

	if err != nil {
		return nil, err
	}

	return imagedClusterDetails, op.client.Do(ctx, req, imagedClusterDetails)
}

/*ListImagedNodes Get all the imaged clusters within FC
 * This operation fetches all the imaged clusters within FC
 *
 * @param input
 * @return *ImagedClustersListResponse
 */
func (op Operations) ListImagedClusters(ctx context.Context, input *ImagedClustersListInput) (*ImagedClustersListResponse, error) {
	path := "/imaged_clusters/list"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, input)

	imagedClustersListResponse := new(ImagedClustersListResponse)

	if err != nil {
		return nil, err
	}

	return imagedClustersListResponse, op.client.Do(ctx, req, imagedClustersListResponse)
}

/*CreateCluster Creates a Cluster
 * This operation submits a request to create a cluster or image nodes based on the input parameters.
 *
 * @param input
 * @return *CreateClusterResponse
 */
func (op Operations) CreateCluster(ctx context.Context, input *CreateClusterInput) (*CreateClusterResponse, error) {
	path := "/imaged_clusters"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, input)

	createClusterResponse := new(CreateClusterResponse)

	if err != nil {
		return nil, err
	}

	return createClusterResponse, op.client.Do(ctx, req, createClusterResponse)
}

/*UpdateCluster Updates a Cluster
 * This operation submits a request to archieve the cluster.
 *
 * @param clusterUUID
 * @return *UpdateClusterData
 */
func (op Operations) UpdateCluster(ctx context.Context, clusterUUID string, updateData *UpdateClusterData) error {
	path := fmt.Sprintf("/imaged_clusters/%s", clusterUUID)

	req, err := op.client.NewRequest(ctx, http.MethodPut, path, updateData)

	if err != nil {
		return err
	}

	return op.client.Do(ctx, req, nil)
}

/*DeleteCluster Deletes Cluster
 * This operation submits a request to delete the cluster on Foundation central
 *
 * @param clusterUUID
 * @return error
 */
func (op Operations) DeleteCluster(ctx context.Context, clusterUUID string) error {
	path := fmt.Sprintf("/imaged_clusters/%s", clusterUUID)

	req, err := op.client.NewRequest(ctx, http.MethodDelete, path, nil)

	if err != nil {
		return err
	}

	return op.client.Do(ctx, req, nil)
}

/*CreateAPIKey Creates a new API Key
 * This Operation creates a new api key which will be used by remote nodes to authenticate with Foundation Central
 *
 * @param input
 * @return CreateAPIKeysResponse
 */
func (op Operations) CreateAPIKey(ctx context.Context, input *CreateAPIKeysInput) (*CreateAPIKeysResponse, error) {
	path := "/api_keys"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, input)
	if err != nil {
		return nil, err
	}

	createAPIResponse := new(CreateAPIKeysResponse)
	return createAPIResponse, op.client.Do(ctx, req, createAPIResponse)
}

/*GetAPIKey Gets the details of an API key
 * Get an api key given its UUID.
 * @param uuid  The uuid of an API Key
 * @return CreateAPIKeysResponse
 */

func (op Operations) GetAPIKey(ctx context.Context, uuid string) (*CreateAPIKeysResponse, error) {
	path := fmt.Sprintf("/api_keys/%s", uuid)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, uuid)
	if err != nil {
		return nil, err
	}

	getAPIResponse := new(CreateAPIKeysResponse)
	return getAPIResponse, op.client.Do(ctx, req, getAPIResponse)
}

/*ListAPIKeys Gets all the API Keys within Foundation Central
 * List all the api keys.
 * @param body
 * @return ListAPIKeysResponse
 */

func (op Operations) ListAPIKeys(ctx context.Context, body *ListMetadataInput) (*ListAPIKeysResponse, error) {
	path := "/api_keys/list"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}

	listAPIKeysResponse := new(ListAPIKeysResponse)
	return listAPIKeysResponse, op.client.Do(ctx, req, listAPIKeysResponse)
}
