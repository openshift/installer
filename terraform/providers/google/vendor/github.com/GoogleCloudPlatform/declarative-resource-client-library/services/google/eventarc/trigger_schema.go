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
package eventarc

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func DCLTriggerSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Eventarc/Trigger",
			Description: "The Eventarc Trigger resource",
			StructName:  "Trigger",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a Trigger",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "trigger",
						Required:    true,
						Description: "A full instance of a Trigger",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a Trigger",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "trigger",
						Required:    true,
						Description: "A full instance of a Trigger",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a Trigger",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "trigger",
						Required:    true,
						Description: "A full instance of a Trigger",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all Trigger",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:     "project",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
					dcl.PathParameters{
						Name:     "location",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
				},
			},
			List: &dcl.Path{
				Description: "The function used to list information about many Trigger",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:     "project",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
					dcl.PathParameters{
						Name:     "location",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
				},
			},
		},
		Components: &dcl.Components{
			Schemas: map[string]*dcl.Component{
				"Trigger": &dcl.Component{
					Title:           "Trigger",
					ID:              "projects/{{project}}/locations/{{location}}/triggers/{{name}}",
					ParentContainer: "project",
					LabelsField:     "labels",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"name",
							"matchingCriteria",
							"destination",
							"project",
							"location",
						},
						Properties: map[string]*dcl.Property{
							"channel": &dcl.Property{
								Type:        "string",
								GoName:      "Channel",
								Description: "Optional. The name of the channel associated with the trigger in `projects/{project}/locations/{location}/channels/{channel}` format. You must provide a channel to receive events from Eventarc SaaS partners.",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Eventarc/Channel",
										Field:    "name",
									},
								},
							},
							"conditions": &dcl.Property{
								Type: "object",
								AdditionalProperties: &dcl.Property{
									Type: "string",
								},
								GoName:      "Conditions",
								ReadOnly:    true,
								Description: "Output only. The reason(s) why a trigger is in FAILED state.",
								Immutable:   true,
							},
							"createTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "CreateTime",
								ReadOnly:    true,
								Description: "Output only. The creation time.",
								Immutable:   true,
							},
							"destination": &dcl.Property{
								Type:        "object",
								GoName:      "Destination",
								GoType:      "TriggerDestination",
								Description: "Required. Destination specifies where the events should be sent to.",
								Properties: map[string]*dcl.Property{
									"cloudFunction": &dcl.Property{
										Type:        "string",
										GoName:      "CloudFunction",
										Description: "[WARNING] Configuring a Cloud Function in Trigger is not supported as of today. The Cloud Function resource name. Format: projects/{project}/locations/{location}/functions/{function}",
										Conflicts: []string{
											"cloudRunService",
											"gke",
											"workflow",
										},
										ResourceReferences: []*dcl.PropertyResourceReference{
											&dcl.PropertyResourceReference{
												Resource: "Cloudfunctions/Function",
												Field:    "name",
											},
										},
										HasLongForm: true,
									},
									"cloudRunService": &dcl.Property{
										Type:        "object",
										GoName:      "CloudRunService",
										GoType:      "TriggerDestinationCloudRunService",
										Description: "Cloud Run fully-managed service that receives the events. The service should be running in the same project of the trigger.",
										Conflicts: []string{
											"cloudFunction",
											"gke",
											"workflow",
										},
										Required: []string{
											"service",
											"region",
										},
										Properties: map[string]*dcl.Property{
											"path": &dcl.Property{
												Type:        "string",
												GoName:      "Path",
												Description: "Optional. The relative path on the Cloud Run service the events should be sent to. The value must conform to the definition of URI path segment (section 3.3 of RFC2396). Examples: \"/route\", \"route\", \"route/subroute\".",
											},
											"region": &dcl.Property{
												Type:        "string",
												GoName:      "Region",
												Description: "Required. The region the Cloud Run service is deployed in.",
											},
											"service": &dcl.Property{
												Type:        "string",
												GoName:      "Service",
												Description: "Required. The name of the Cloud Run service being addressed. See https://cloud.google.com/run/docs/reference/rest/v1/namespaces.services. Only services located in the same project of the trigger object can be addressed.",
												ResourceReferences: []*dcl.PropertyResourceReference{
													&dcl.PropertyResourceReference{
														Resource: "Run/Service",
														Field:    "name",
													},
												},
											},
										},
									},
									"gke": &dcl.Property{
										Type:        "object",
										GoName:      "Gke",
										GoType:      "TriggerDestinationGke",
										Description: "A GKE service capable of receiving events. The service should be running in the same project as the trigger.",
										Conflicts: []string{
											"cloudRunService",
											"cloudFunction",
											"workflow",
										},
										Required: []string{
											"cluster",
											"location",
											"namespace",
											"service",
										},
										Properties: map[string]*dcl.Property{
											"cluster": &dcl.Property{
												Type:        "string",
												GoName:      "Cluster",
												Description: "Required. The name of the cluster the GKE service is running in. The cluster must be running in the same project as the trigger being created.",
												ResourceReferences: []*dcl.PropertyResourceReference{
													&dcl.PropertyResourceReference{
														Resource: "Container/Cluster",
														Field:    "selfLink",
													},
												},
											},
											"location": &dcl.Property{
												Type:        "string",
												GoName:      "Location",
												Description: "Required. The name of the Google Compute Engine in which the cluster resides, which can either be compute zone (for example, us-central1-a) for the zonal clusters or region (for example, us-central1) for regional clusters.",
											},
											"namespace": &dcl.Property{
												Type:        "string",
												GoName:      "Namespace",
												Description: "Required. The namespace the GKE service is running in.",
											},
											"path": &dcl.Property{
												Type:        "string",
												GoName:      "Path",
												Description: "Optional. The relative path on the GKE service the events should be sent to. The value must conform to the definition of a URI path segment (section 3.3 of RFC2396). Examples: \"/route\", \"route\", \"route/subroute\".",
											},
											"service": &dcl.Property{
												Type:        "string",
												GoName:      "Service",
												Description: "Required. Name of the GKE service.",
											},
										},
									},
									"httpEndpoint": &dcl.Property{
										Type:        "object",
										GoName:      "HttpEndpoint",
										GoType:      "TriggerDestinationHttpEndpoint",
										Description: "An HTTP endpoint destination described by an URI.",
										Required: []string{
											"uri",
										},
										Properties: map[string]*dcl.Property{
											"uri": &dcl.Property{
												Type:        "string",
												GoName:      "Uri",
												Description: "Required. The URI of the HTTP enpdoint. The value must be a RFC2396 URI string. Examples: `http://10.10.10.8:80/route`, `http://svc.us-central1.p.local:8080/`. Only HTTP and HTTPS protocols are supported. The host can be either a static IP addressable from the VPC specified by the network config, or an internal DNS hostname of the service resolvable via Cloud DNS.",
											},
										},
									},
									"networkConfig": &dcl.Property{
										Type:        "object",
										GoName:      "NetworkConfig",
										GoType:      "TriggerDestinationNetworkConfig",
										Description: "Optional. Network config is used to configure how Eventarc resolves and connect to a destination. This should only be used with HttpEndpoint destination type.",
										Required: []string{
											"networkAttachment",
										},
										Properties: map[string]*dcl.Property{
											"networkAttachment": &dcl.Property{
												Type:        "string",
												GoName:      "NetworkAttachment",
												Description: "Required. Name of the NetworkAttachment that allows access to the destination VPC. Format: `projects/{PROJECT_ID}/regions/{REGION}/networkAttachments/{NETWORK_ATTACHMENT_NAME}`",
												ResourceReferences: []*dcl.PropertyResourceReference{
													&dcl.PropertyResourceReference{
														Resource: "Compute/NetworkAttachment",
														Field:    "selfLink",
													},
												},
											},
										},
									},
									"workflow": &dcl.Property{
										Type:        "string",
										GoName:      "Workflow",
										Description: "The resource name of the Workflow whose Executions are triggered by the events. The Workflow resource should be deployed in the same project as the trigger. Format: `projects/{project}/locations/{location}/workflows/{workflow}`",
										Conflicts: []string{
											"cloudRunService",
											"cloudFunction",
											"gke",
										},
										ResourceReferences: []*dcl.PropertyResourceReference{
											&dcl.PropertyResourceReference{
												Resource: "Workflows/Workflow",
												Field:    "name",
											},
										},
										HasLongForm: true,
									},
								},
							},
							"etag": &dcl.Property{
								Type:        "string",
								GoName:      "Etag",
								ReadOnly:    true,
								Description: "Output only. This checksum is computed by the server based on the value of other fields, and may be sent only on create requests to ensure the client has an up-to-date value before proceeding.",
								Immutable:   true,
							},
							"eventDataContentType": &dcl.Property{
								Type:          "string",
								GoName:        "EventDataContentType",
								Description:   "Optional. EventDataContentType specifies the type of payload in MIME format that is expected from the CloudEvent data field. This is set to `application/json` if the value is not defined.",
								ServerDefault: true,
							},
							"labels": &dcl.Property{
								Type: "object",
								AdditionalProperties: &dcl.Property{
									Type: "string",
								},
								GoName:      "Labels",
								Description: "Optional. User labels attached to the triggers that can be used to group resources.",
							},
							"location": &dcl.Property{
								Type:        "string",
								GoName:      "Location",
								Description: "The location for the resource",
								Immutable:   true,
								Parameter:   true,
							},
							"matchingCriteria": &dcl.Property{
								Type:        "array",
								GoName:      "MatchingCriteria",
								Description: "Required. null The list of filters that applies to event attributes. Only events that match all the provided filters will be sent to the destination.",
								SendEmpty:   true,
								ListType:    "set",
								Items: &dcl.Property{
									Type:   "object",
									GoType: "TriggerMatchingCriteria",
									Required: []string{
										"attribute",
										"value",
									},
									Properties: map[string]*dcl.Property{
										"attribute": &dcl.Property{
											Type:        "string",
											GoName:      "Attribute",
											Description: "Required. The name of a CloudEvents attribute. Currently, only a subset of attributes are supported for filtering. All triggers MUST provide a filter for the 'type' attribute.",
										},
										"operator": &dcl.Property{
											Type:        "string",
											GoName:      "Operator",
											Description: "Optional. The operator used for matching the events with the value of the filter. If not specified, only events that have an exact key-value pair specified in the filter are matched. The only allowed value is `match-path-pattern`.",
										},
										"value": &dcl.Property{
											Type:        "string",
											GoName:      "Value",
											Description: "Required. The value for the attribute. See https://cloud.google.com/eventarc/docs/creating-triggers#trigger-gcloud for available values.",
										},
									},
								},
							},
							"name": &dcl.Property{
								Type:        "string",
								GoName:      "Name",
								Description: "Required. The resource name of the trigger. Must be unique within the location on the project.",
								Immutable:   true,
								HasLongForm: true,
							},
							"project": &dcl.Property{
								Type:        "string",
								GoName:      "Project",
								Description: "The project for the resource",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Cloudresourcemanager/Project",
										Field:    "name",
										Parent:   true,
									},
								},
								Parameter: true,
							},
							"serviceAccount": &dcl.Property{
								Type:        "string",
								GoName:      "ServiceAccount",
								Description: "Optional. The IAM service account email associated with the trigger. The service account represents the identity of the trigger. The principal who calls this API must have `iam.serviceAccounts.actAs` permission in the service account. See https://cloud.google.com/iam/docs/understanding-service-accounts#sa_common for more information. For Cloud Run destinations, this service account is used to generate identity tokens when invoking the service. See https://cloud.google.com/run/docs/triggering/pubsub-push#create-service-account for information on how to invoke authenticated Cloud Run services. In order to create Audit Log triggers, the service account should also have `roles/eventarc.eventReceiver` IAM role.",
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Iam/ServiceAccount",
										Field:    "email",
									},
								},
							},
							"transport": &dcl.Property{
								Type:          "object",
								GoName:        "Transport",
								GoType:        "TriggerTransport",
								Description:   "Optional. In order to deliver messages, Eventarc may use other GCP products as transport intermediary. This field contains a reference to that transport intermediary. This information can be used for debugging purposes.",
								Immutable:     true,
								ServerDefault: true,
								Properties: map[string]*dcl.Property{
									"pubsub": &dcl.Property{
										Type:        "object",
										GoName:      "Pubsub",
										GoType:      "TriggerTransportPubsub",
										Description: "The Pub/Sub topic and subscription used by Eventarc as delivery intermediary.",
										Immutable:   true,
										Properties: map[string]*dcl.Property{
											"subscription": &dcl.Property{
												Type:        "string",
												GoName:      "Subscription",
												ReadOnly:    true,
												Description: "Output only. The name of the Pub/Sub subscription created and managed by Eventarc system as a transport for the event delivery. Format: `projects/{PROJECT_ID}/subscriptions/{SUBSCRIPTION_NAME}`.",
												Immutable:   true,
												ResourceReferences: []*dcl.PropertyResourceReference{
													&dcl.PropertyResourceReference{
														Resource: "Pubsub/Subscription",
														Field:    "name",
													},
												},
												HasLongForm: true,
											},
											"topic": &dcl.Property{
												Type:        "string",
												GoName:      "Topic",
												Description: "Optional. The name of the Pub/Sub topic created and managed by Eventarc system as a transport for the event delivery. Format: `projects/{PROJECT_ID}/topics/{TOPIC_NAME}. You may set an existing topic for triggers of the type google.cloud.pubsub.topic.v1.messagePublished` only. The topic you provide here will not be deleted by Eventarc at trigger deletion.",
												Immutable:   true,
												ResourceReferences: []*dcl.PropertyResourceReference{
													&dcl.PropertyResourceReference{
														Resource: "Pubsub/Topic",
														Field:    "name",
													},
												},
												HasLongForm: true,
											},
										},
									},
								},
							},
							"uid": &dcl.Property{
								Type:        "string",
								GoName:      "Uid",
								ReadOnly:    true,
								Description: "Output only. Server assigned unique identifier for the trigger. The value is a UUID4 string and guaranteed to remain unchanged until the resource is deleted.",
								Immutable:   true,
							},
							"updateTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "UpdateTime",
								ReadOnly:    true,
								Description: "Output only. The last-modified time.",
								Immutable:   true,
							},
						},
					},
				},
			},
		},
	}
}
