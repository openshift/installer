# Agent-Based Installer Authentication and Authorization - OpenShift 4.18

## Introduction
Since OpenShift 4.18, the Assisted Service (AS) APIs used by the Agent-Based Installer (ABI) have been secured through authentication and authorization mechanisms. This added security ensures that only authorized services can interact with AS. This document provides an overview of how the authentication and authorization mechanisms were implemented for ABI.

## Enabling Authentication and Authorization

To enable authentication and authorization, the `AuthType` setting must be configured in the `assisted-service.env.template` file. Set the following value:

```bash
AuthType = "agent-installer-local"
```

This setting ensures that the authentication mechanism required by the agent-based installer is activated.

## Token Workflow: Day 1
```mermaid
flowchart TD
  A[Create Agent ISO] --> B[Generate 3 tokens: agentAuth, userAuth, watcherAuth. Also generate public key]
  B --> C[Store tokens and base64 encoded public key in asset store hidden file '.openshift_state_install.json']
  B --> D[Embed tokens and base64 encoded public key in /usr/local/share/assisted-service/assisted-service.env file in ISO]
  D --> E[Boot hosts using ISO]
  E --> F[Assisted service running on rendezvous host validates the token from HTTP requests coming from each node]
  F --> G{Are tokens valid?}
  G -->|Yes| H[Allow successful cluster installation]
  G -->|No| I[Fail cluster installation]
```

## Tokens for Agent-Based Installer

When creating an Agent ISO (or PXE artifacts) for installing a new cluster, ABI generates three separate authentication tokens: `agentAuth`, `userAuth`, and `watcherAuth`. These tokens serve different purposes and have varying levels of access. The tokens are the part of the ignition that gets embeded in the CoreOS base ISO along with other ABI specific configurations for a cluster. The tokens are stored in `/usr/local/share/assisted-service/assisted-service.env` in the ISO. The base64 encoded public key required to validate these tokens is also saved in the same file. Additionally, this data also gets stored in a hidden file `.openshift_state_install.json` present in the directory location mentioned when creating ISO.

### Types of Tokens:
1. **agentAuth**:
   - **Access Level**: Read-Write (RW)
   - **Used By**: Agent services running in the live environment when host is booted.
   
2. **userAuth**:
   - **Access Level**: Read-Write (RW)
   - **Used By**: Systemd services (such as agent-register-infraenv.service, agent-register-cluster.service, etc) and the agent-installer-client in the Assisted Service.
   
3. **watcherAuth**:
   - **Access Level**: Read-Only (RO)
   - **Used By**: The monitoring commands (the `wait-for` command on Day 1 and `monitor-add-nodes` command on Day 2).

### Managing Access with Tokens and Security Schemes in swagger.yaml
Access control in swagger.yaml is managed through a combination of tokens and endpoint annotations. Each token corresponds to a specific security scheme and dictates the level of access based on the type of operation and annotation.

- agentAuth Token and Scheme:

A `GET` endpoint annotated with the `agentAuth` scheme grants `read-only (RO)` access.
A `POST` endpoint annotated with the same scheme enables `read-write (RW)` access.

- userAuth Token and Scheme:

Provides read-write (RW) access for both `GET` and `POST` endpoints.
- watcherAuth Token and Scheme:

Grants `read-only (RO)` access, applicable to `GET` endpoints only.

By combining the appropriate token and security annotations, `swagger.yaml` effectively controls access for various user personas.

### Token Expiry Behavior
- **Day 1** Tokens: The tokens generated on Day 1 never expire, meaning the agent ISO remains valid indefinitely. This is useful for use cases where customers create an agent ISO on one day and boot it on a later day (e.g., in a data center).
- **Day 2** Tokens: Tokens generated for nodes ISO on Day 2 expire after 48 hours. This ensures that the ISO is only valid for a short period.

### Day 2 token management
When a Day 2 node ISO is created, all three tokens are saved in the `openshift-config` namespace as a new secret named `agent-auth-token`. The following steps outline how token management works:
1. If the `openshift-config/agent-auth-token` secret does not exist in the cluster, the secret is created with the 3 tokens mentioned above along with the public key.
2. If the `openshift-config/agent-auth-token` secret already exists in the cluster, the tokens are retrieved from the secret (instead of being generated).
3. All tokens expire at the same time, so only one token needs to be checked for expiry.
4. If the token in the secret is older than 24 hours when creating a node ISO, the installer generates new tokens. These new tokens are saved in the asset store, and the secret in the cluster is updated with the new token values.
5. If the token in the secret is not older than 24 hours, the installer retrieves the existing token from the cluster, updates the asset store with the retrieved token, and uses the retrieved token for the new node ISO.
6. When a node is booted using the Day 2 `nodes.iso` to add worker nodes, the `agent-auth-token-status.service` systemd service will check if the authentication tokens have expired. If the tokens are expired, the following message will be displayed:

```
The authentication token has expired. Please generate a new ISO using the "oc adm node-image create" command, then reboot the node.
```

This ensures that the user is notified of expired tokens and can take appropriate action to generate a new ISO.

## Token Workflow: Day 2
```mermaid
flowchart TD
    A[Check if openshift-config/agent-auth-token exists] -->|No| B[Create secret with 3 tokens and public key]
    A -->|Yes| C[Retrieve tokens from existing secret]
    B --> D[Generate new tokens if secret does not exist]
    C --> D[Check if token is older than 24 hours]
    D -->|Yes| E[Generate new tokens, update asset store, and secret]
    D -->|No| F[Retrieve token, update asset store, use token for ISO]
    F --> G[Node boots and checks token expiry]
    G -->|Expired| H[Display message: 'Authentication token expired. Generate a new ISO.']
    G -->|Not expired| I[Proceed with installation using retrieved token]
```

## Authorization

To implement authorization for ABI, a new security definition is added in the `swagger.yaml` file in AS. The new security definition is as follows:

```yaml
watcherAuth:
  type: apiKey
  in: header
  name: Watcher-Authorization
```

The required endpoints are annotated with this security definition. For example:

```yaml
/v2/clusters:
  get:
    tags:
      - installer
    security:
      - userAuth: [admin, read-only-admin, user]
      - watcherAuth: []
    description: Retrieves the list of OpenShift clusters.
    operationId: v2ListClusters
```

### Authorization Handler

The Assisted Service API includes an authorization handler implemented in the `pkg/auth/agent_local_authz_handler.go` file. The `AgentLocalAuthzHandler`  is responsible for verifying that the `auth_scheme` in the token matches the required `authScheme` for the endpoint. The authorization logic works as follows:

1. The `auth_scheme` claim in the JWT token is extracted.
2. If the claim is malformed or missing, an error is returned.
3. The `auth_scheme` claim embedded in the token specifies the type of authorization (e.g., `watcherAuth`) and is compared with the `authScheme` expected by the endpoint. This ensures that the token matches the intended purpose of the request.  
4. If the `auth_scheme` in the token does not match the `authScheme` required by the endpoint, or if an incorrect token header is provided, the request is rejected with a `Forbidden (403)` error. 

Code implementation https://github.com/openshift/assisted-service/blame/973249d13ee9c83fd8a620a94956bcd16927b352/pkg/auth/agent_local_authz_handler.go#L46-L80

For example:  
- Suppose the endpoint `/api/assisted-install/v2/infra-envs` requires the `watcherAuth` scheme.  
- A request is made with the following command, but an **incorrect token header** is sent:  
  ```bash
  curl http://localhost:8090/api/assisted-install/v2/infra-envs \
       -H "Authorization: eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRoX3NjaGVtZSI6IndhdGNoZXJBdXRoIn0.hJ5EuBOttDQ1-ANWUJlc97H1k63d-clDCaFUZ84XSYdF03A938OfMU6RHWltv0AiWZEHNJ0gLrjLxR-1yUYqiA"
  ```  
- Although the token specifies the correct `auth_scheme` (`watcherAuth`), the **header sent is incorrect** for this request.

Expected response:  
```json
{
  "code": 403,
  "message": "authClaim watcherAuth is unauthorized to access"
}
```

**Outcome:** 
The system identifies that the provided header does not satisfy the endpoint’s requirements and correctly responds with an error. This confirms that tokens must not only align with the expected `authScheme` but also be passed in the correct format and header. The authorization logic ensures that only authorized requests can access protected endpoints in Assisted Service.

## Authentication API (Server Side)

The authentication server validates the token provided by the client using the `validateToken` function. The following steps explain the server-side validation process:

1. When an agent or nodes ISO is created, a set of public and private keys is generated. The private key is used to sign the JWT tokens, while the associated public key is made available to the API server.
2. The token is passed to the API server as part of the request. The API server validates the token by verifying its signature against the public key. This ensures the token was signed by the correct private key and has not been tampered with.
3. If the token is invalid or expired, an error is logged and the request is rejected with a `Unauthorized` error.
4. The claims in the token are extracted.
5. If the `exp` (expiration) claim is found, the server checks whether the token has expired.
   - In the install workflow, the `exp` claim will not be present.
   - In the add-nodes workflow, the token expiration is checked. If expired, a new token must be generated.
6. If the token is valid and not expired, the claims are returned and the request proceeds.
