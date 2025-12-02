// Copyright 2023 Google LLC. All Rights Reserved.
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
package dataproc

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func DCLWorkflowTemplateSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Dataproc/WorkflowTemplate",
			Description: "The Dataproc WorkflowTemplate resource",
			StructName:  "WorkflowTemplate",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a WorkflowTemplate",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "workflowTemplate",
						Required:    true,
						Description: "A full instance of a WorkflowTemplate",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a WorkflowTemplate",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "workflowTemplate",
						Required:    true,
						Description: "A full instance of a WorkflowTemplate",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a WorkflowTemplate",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "workflowTemplate",
						Required:    true,
						Description: "A full instance of a WorkflowTemplate",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all WorkflowTemplate",
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
				Description: "The function used to list information about many WorkflowTemplate",
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
				"WorkflowTemplate": &dcl.Component{
					Title:           "WorkflowTemplate",
					ID:              "projects/{{project}}/locations/{{location}}/workflowTemplates/{{name}}",
					ParentContainer: "project",
					LabelsField:     "labels",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"name",
							"placement",
							"jobs",
							"project",
							"location",
						},
						Properties: map[string]*dcl.Property{
							"createTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "CreateTime",
								ReadOnly:    true,
								Description: "Output only. The time template was created.",
								Immutable:   true,
							},
							"dagTimeout": &dcl.Property{
								Type:        "string",
								GoName:      "DagTimeout",
								Description: "Optional. Timeout duration for the DAG of jobs, expressed in seconds (see [JSON representation of duration](https://developers.google.com/protocol-buffers/docs/proto3#json)). The timeout duration must be from 10 minutes (\"600s\") to 24 hours (\"86400s\"). The timer begins when the first job is submitted. If the workflow is running at the end of the timeout period, any remaining jobs are cancelled, the workflow is ended, and if the workflow was running on a [managed cluster](/dataproc/docs/concepts/workflows/using-workflows#configuring_or_selecting_a_cluster), the cluster is deleted.",
								Immutable:   true,
							},
							"jobs": &dcl.Property{
								Type:        "array",
								GoName:      "Jobs",
								Description: "Required. The Directed Acyclic Graph of Jobs to submit.",
								Immutable:   true,
								SendEmpty:   true,
								ListType:    "list",
								Items: &dcl.Property{
									Type:   "object",
									GoType: "WorkflowTemplateJobs",
									Required: []string{
										"stepId",
									},
									Properties: map[string]*dcl.Property{
										"hadoopJob": &dcl.Property{
											Type:        "object",
											GoName:      "HadoopJob",
											GoType:      "WorkflowTemplateJobsHadoopJob",
											Description: "Optional. Job is a Hadoop job.",
											Immutable:   true,
											Properties: map[string]*dcl.Property{
												"archiveUris": &dcl.Property{
													Type:        "array",
													GoName:      "ArchiveUris",
													Description: "Optional. HCFS URIs of archives to be extracted in the working directory of Hadoop drivers and tasks. Supported file types: .jar, .tar, .tar.gz, .tgz, or .zip.",
													Immutable:   true,
													SendEmpty:   true,
													ListType:    "list",
													Items: &dcl.Property{
														Type:   "string",
														GoType: "string",
													},
												},
												"args": &dcl.Property{
													Type:        "array",
													GoName:      "Args",
													Description: "Optional. The arguments to pass to the driver. Do not include arguments, such as `-libjars` or `-Dfoo=bar`, that can be set as job properties, since a collision may occur that causes an incorrect job submission.",
													Immutable:   true,
													SendEmpty:   true,
													ListType:    "list",
													Items: &dcl.Property{
														Type:   "string",
														GoType: "string",
													},
												},
												"fileUris": &dcl.Property{
													Type:        "array",
													GoName:      "FileUris",
													Description: "Optional. HCFS (Hadoop Compatible Filesystem) URIs of files to be copied to the working directory of Hadoop drivers and distributed tasks. Useful for naively parallel tasks.",
													Immutable:   true,
													SendEmpty:   true,
													ListType:    "list",
													Items: &dcl.Property{
														Type:   "string",
														GoType: "string",
													},
												},
												"jarFileUris": &dcl.Property{
													Type:        "array",
													GoName:      "JarFileUris",
													Description: "Optional. Jar file URIs to add to the CLASSPATHs of the Hadoop driver and tasks.",
													Immutable:   true,
													SendEmpty:   true,
													ListType:    "list",
													Items: &dcl.Property{
														Type:   "string",
														GoType: "string",
													},
												},
												"loggingConfig": &dcl.Property{
													Type:        "object",
													GoName:      "LoggingConfig",
													GoType:      "WorkflowTemplateJobsHadoopJobLoggingConfig",
													Description: "Optional. The runtime log config for job execution.",
													Immutable:   true,
													Properties: map[string]*dcl.Property{
														"driverLogLevels": &dcl.Property{
															Type: "object",
															AdditionalProperties: &dcl.Property{
																Type: "string",
															},
															GoName:      "DriverLogLevels",
															Description: "The per-package log levels for the driver. This may include \"root\" package name to configure rootLogger. Examples: 'com.google = FATAL', 'root = INFO', 'org.apache = DEBUG'",
															Immutable:   true,
														},
													},
												},
												"mainClass": &dcl.Property{
													Type:        "string",
													GoName:      "MainClass",
													Description: "The name of the driver's main class. The jar file containing the class must be in the default CLASSPATH or specified in `jar_file_uris`.",
													Immutable:   true,
												},
												"mainJarFileUri": &dcl.Property{
													Type:        "string",
													GoName:      "MainJarFileUri",
													Description: "The HCFS URI of the jar file containing the main class. Examples: 'gs://foo-bucket/analytics-binaries/extract-useful-metrics-mr.jar' 'hdfs:/tmp/test-samples/custom-wordcount.jar' 'file:///home/usr/lib/hadoop-mapreduce/hadoop-mapreduce-examples.jar'",
													Immutable:   true,
												},
												"properties": &dcl.Property{
													Type: "object",
													AdditionalProperties: &dcl.Property{
														Type: "string",
													},
													GoName:      "Properties",
													Description: "Optional. A mapping of property names to values, used to configure Hadoop. Properties that conflict with values set by the Dataproc API may be overwritten. Can include properties set in /etc/hadoop/conf/*-site and classes in user code.",
													Immutable:   true,
												},
											},
										},
										"hiveJob": &dcl.Property{
											Type:        "object",
											GoName:      "HiveJob",
											GoType:      "WorkflowTemplateJobsHiveJob",
											Description: "Optional. Job is a Hive job.",
											Immutable:   true,
											Properties: map[string]*dcl.Property{
												"continueOnFailure": &dcl.Property{
													Type:        "boolean",
													GoName:      "ContinueOnFailure",
													Description: "Optional. Whether to continue executing queries if a query fails. The default value is `false`. Setting to `true` can be useful when executing independent parallel queries.",
													Immutable:   true,
												},
												"jarFileUris": &dcl.Property{
													Type:        "array",
													GoName:      "JarFileUris",
													Description: "Optional. HCFS URIs of jar files to add to the CLASSPATH of the Hive server and Hadoop MapReduce (MR) tasks. Can contain Hive SerDes and UDFs.",
													Immutable:   true,
													SendEmpty:   true,
													ListType:    "list",
													Items: &dcl.Property{
														Type:   "string",
														GoType: "string",
													},
												},
												"properties": &dcl.Property{
													Type: "object",
													AdditionalProperties: &dcl.Property{
														Type: "string",
													},
													GoName:      "Properties",
													Description: "Optional. A mapping of property names and values, used to configure Hive. Properties that conflict with values set by the Dataproc API may be overwritten. Can include properties set in /etc/hadoop/conf/*-site.xml, /etc/hive/conf/hive-site.xml, and classes in user code.",
													Immutable:   true,
												},
												"queryFileUri": &dcl.Property{
													Type:        "string",
													GoName:      "QueryFileUri",
													Description: "The HCFS URI of the script that contains Hive queries.",
													Immutable:   true,
												},
												"queryList": &dcl.Property{
													Type:        "object",
													GoName:      "QueryList",
													GoType:      "WorkflowTemplateJobsHiveJobQueryList",
													Description: "A list of queries.",
													Immutable:   true,
													Required: []string{
														"queries",
													},
													Properties: map[string]*dcl.Property{
														"queries": &dcl.Property{
															Type:        "array",
															GoName:      "Queries",
															Description: "Required. The queries to execute. You do not need to end a query expression with a semicolon. Multiple queries can be specified in one string by separating each with a semicolon. Here is an example of a Dataproc API snippet that uses a QueryList to specify a HiveJob: \"hiveJob\": { \"queryList\": { \"queries\": [ \"query1\", \"query2\", \"query3;query4\", ] } }",
															Immutable:   true,
															SendEmpty:   true,
															ListType:    "list",
															Items: &dcl.Property{
																Type:   "string",
																GoType: "string",
															},
														},
													},
												},
												"scriptVariables": &dcl.Property{
													Type: "object",
													AdditionalProperties: &dcl.Property{
														Type: "string",
													},
													GoName:      "ScriptVariables",
													Description: "Optional. Mapping of query variable names to values (equivalent to the Hive command: `SET name=\"value\";`).",
													Immutable:   true,
												},
											},
										},
										"labels": &dcl.Property{
											Type: "object",
											AdditionalProperties: &dcl.Property{
												Type: "string",
											},
											GoName:      "Labels",
											Description: "Optional. The labels to associate with this job. Label keys must be between 1 and 63 characters long, and must conform to the following regular expression: p{Ll}p{Lo}{0,62} Label values must be between 1 and 63 characters long, and must conform to the following regular expression: [p{Ll}p{Lo}p{N}_-]{0,63} No more than 32 labels can be associated with a given job.",
											Immutable:   true,
										},
										"pigJob": &dcl.Property{
											Type:        "object",
											GoName:      "PigJob",
											GoType:      "WorkflowTemplateJobsPigJob",
											Description: "Optional. Job is a Pig job.",
											Immutable:   true,
											Properties: map[string]*dcl.Property{
												"continueOnFailure": &dcl.Property{
													Type:        "boolean",
													GoName:      "ContinueOnFailure",
													Description: "Optional. Whether to continue executing queries if a query fails. The default value is `false`. Setting to `true` can be useful when executing independent parallel queries.",
													Immutable:   true,
												},
												"jarFileUris": &dcl.Property{
													Type:        "array",
													GoName:      "JarFileUris",
													Description: "Optional. HCFS URIs of jar files to add to the CLASSPATH of the Pig Client and Hadoop MapReduce (MR) tasks. Can contain Pig UDFs.",
													Immutable:   true,
													SendEmpty:   true,
													ListType:    "list",
													Items: &dcl.Property{
														Type:   "string",
														GoType: "string",
													},
												},
												"loggingConfig": &dcl.Property{
													Type:        "object",
													GoName:      "LoggingConfig",
													GoType:      "WorkflowTemplateJobsPigJobLoggingConfig",
													Description: "Optional. The runtime log config for job execution.",
													Immutable:   true,
													Properties: map[string]*dcl.Property{
														"driverLogLevels": &dcl.Property{
															Type: "object",
															AdditionalProperties: &dcl.Property{
																Type: "string",
															},
															GoName:      "DriverLogLevels",
															Description: "The per-package log levels for the driver. This may include \"root\" package name to configure rootLogger. Examples: 'com.google = FATAL', 'root = INFO', 'org.apache = DEBUG'",
															Immutable:   true,
														},
													},
												},
												"properties": &dcl.Property{
													Type: "object",
													AdditionalProperties: &dcl.Property{
														Type: "string",
													},
													GoName:      "Properties",
													Description: "Optional. A mapping of property names to values, used to configure Pig. Properties that conflict with values set by the Dataproc API may be overwritten. Can include properties set in /etc/hadoop/conf/*-site.xml, /etc/pig/conf/pig.properties, and classes in user code.",
													Immutable:   true,
												},
												"queryFileUri": &dcl.Property{
													Type:        "string",
													GoName:      "QueryFileUri",
													Description: "The HCFS URI of the script that contains the Pig queries.",
													Immutable:   true,
												},
												"queryList": &dcl.Property{
													Type:        "object",
													GoName:      "QueryList",
													GoType:      "WorkflowTemplateJobsPigJobQueryList",
													Description: "A list of queries.",
													Immutable:   true,
													Required: []string{
														"queries",
													},
													Properties: map[string]*dcl.Property{
														"queries": &dcl.Property{
															Type:        "array",
															GoName:      "Queries",
															Description: "Required. The queries to execute. You do not need to end a query expression with a semicolon. Multiple queries can be specified in one string by separating each with a semicolon. Here is an example of a Dataproc API snippet that uses a QueryList to specify a HiveJob: \"hiveJob\": { \"queryList\": { \"queries\": [ \"query1\", \"query2\", \"query3;query4\", ] } }",
															Immutable:   true,
															SendEmpty:   true,
															ListType:    "list",
															Items: &dcl.Property{
																Type:   "string",
																GoType: "string",
															},
														},
													},
												},
												"scriptVariables": &dcl.Property{
													Type: "object",
													AdditionalProperties: &dcl.Property{
														Type: "string",
													},
													GoName:      "ScriptVariables",
													Description: "Optional. Mapping of query variable names to values (equivalent to the Pig command: `name=[value]`).",
													Immutable:   true,
												},
											},
										},
										"prerequisiteStepIds": &dcl.Property{
											Type:        "array",
											GoName:      "PrerequisiteStepIds",
											Description: "Optional. The optional list of prerequisite job step_ids. If not specified, the job will start at the beginning of workflow.",
											Immutable:   true,
											SendEmpty:   true,
											ListType:    "list",
											Items: &dcl.Property{
												Type:   "string",
												GoType: "string",
											},
										},
										"prestoJob": &dcl.Property{
											Type:        "object",
											GoName:      "PrestoJob",
											GoType:      "WorkflowTemplateJobsPrestoJob",
											Description: "Optional. Job is a Presto job.",
											Immutable:   true,
											Properties: map[string]*dcl.Property{
												"clientTags": &dcl.Property{
													Type:        "array",
													GoName:      "ClientTags",
													Description: "Optional. Presto client tags to attach to this query",
													Immutable:   true,
													SendEmpty:   true,
													ListType:    "list",
													Items: &dcl.Property{
														Type:   "string",
														GoType: "string",
													},
												},
												"continueOnFailure": &dcl.Property{
													Type:        "boolean",
													GoName:      "ContinueOnFailure",
													Description: "Optional. Whether to continue executing queries if a query fails. The default value is `false`. Setting to `true` can be useful when executing independent parallel queries.",
													Immutable:   true,
												},
												"loggingConfig": &dcl.Property{
													Type:        "object",
													GoName:      "LoggingConfig",
													GoType:      "WorkflowTemplateJobsPrestoJobLoggingConfig",
													Description: "Optional. The runtime log config for job execution.",
													Immutable:   true,
													Properties: map[string]*dcl.Property{
														"driverLogLevels": &dcl.Property{
															Type: "object",
															AdditionalProperties: &dcl.Property{
																Type: "string",
															},
															GoName:      "DriverLogLevels",
															Description: "The per-package log levels for the driver. This may include \"root\" package name to configure rootLogger. Examples: 'com.google = FATAL', 'root = INFO', 'org.apache = DEBUG'",
															Immutable:   true,
														},
													},
												},
												"outputFormat": &dcl.Property{
													Type:        "string",
													GoName:      "OutputFormat",
													Description: "Optional. The format in which query output will be displayed. See the Presto documentation for supported output formats",
													Immutable:   true,
												},
												"properties": &dcl.Property{
													Type: "object",
													AdditionalProperties: &dcl.Property{
														Type: "string",
													},
													GoName:      "Properties",
													Description: "Optional. A mapping of property names to values. Used to set Presto [session properties](https://prestodb.io/docs/current/sql/set-session.html) Equivalent to using the --session flag in the Presto CLI",
													Immutable:   true,
												},
												"queryFileUri": &dcl.Property{
													Type:        "string",
													GoName:      "QueryFileUri",
													Description: "The HCFS URI of the script that contains SQL queries.",
													Immutable:   true,
												},
												"queryList": &dcl.Property{
													Type:        "object",
													GoName:      "QueryList",
													GoType:      "WorkflowTemplateJobsPrestoJobQueryList",
													Description: "A list of queries.",
													Immutable:   true,
													Required: []string{
														"queries",
													},
													Properties: map[string]*dcl.Property{
														"queries": &dcl.Property{
															Type:        "array",
															GoName:      "Queries",
															Description: "Required. The queries to execute. You do not need to end a query expression with a semicolon. Multiple queries can be specified in one string by separating each with a semicolon. Here is an example of a Dataproc API snippet that uses a QueryList to specify a HiveJob: \"hiveJob\": { \"queryList\": { \"queries\": [ \"query1\", \"query2\", \"query3;query4\", ] } }",
															Immutable:   true,
															SendEmpty:   true,
															ListType:    "list",
															Items: &dcl.Property{
																Type:   "string",
																GoType: "string",
															},
														},
													},
												},
											},
										},
										"pysparkJob": &dcl.Property{
											Type:        "object",
											GoName:      "PysparkJob",
											GoType:      "WorkflowTemplateJobsPysparkJob",
											Description: "Optional. Job is a PySpark job.",
											Immutable:   true,
											Required: []string{
												"mainPythonFileUri",
											},
											Properties: map[string]*dcl.Property{
												"archiveUris": &dcl.Property{
													Type:        "array",
													GoName:      "ArchiveUris",
													Description: "Optional. HCFS URIs of archives to be extracted into the working directory of each executor. Supported file types: .jar, .tar, .tar.gz, .tgz, and .zip.",
													Immutable:   true,
													SendEmpty:   true,
													ListType:    "list",
													Items: &dcl.Property{
														Type:   "string",
														GoType: "string",
													},
												},
												"args": &dcl.Property{
													Type:        "array",
													GoName:      "Args",
													Description: "Optional. The arguments to pass to the driver. Do not include arguments, such as `--conf`, that can be set as job properties, since a collision may occur that causes an incorrect job submission.",
													Immutable:   true,
													SendEmpty:   true,
													ListType:    "list",
													Items: &dcl.Property{
														Type:   "string",
														GoType: "string",
													},
												},
												"fileUris": &dcl.Property{
													Type:        "array",
													GoName:      "FileUris",
													Description: "Optional. HCFS URIs of files to be placed in the working directory of each executor. Useful for naively parallel tasks.",
													Immutable:   true,
													SendEmpty:   true,
													ListType:    "list",
													Items: &dcl.Property{
														Type:   "string",
														GoType: "string",
													},
												},
												"jarFileUris": &dcl.Property{
													Type:        "array",
													GoName:      "JarFileUris",
													Description: "Optional. HCFS URIs of jar files to add to the CLASSPATHs of the Python driver and tasks.",
													Immutable:   true,
													SendEmpty:   true,
													ListType:    "list",
													Items: &dcl.Property{
														Type:   "string",
														GoType: "string",
													},
												},
												"loggingConfig": &dcl.Property{
													Type:        "object",
													GoName:      "LoggingConfig",
													GoType:      "WorkflowTemplateJobsPysparkJobLoggingConfig",
													Description: "Optional. The runtime log config for job execution.",
													Immutable:   true,
													Properties: map[string]*dcl.Property{
														"driverLogLevels": &dcl.Property{
															Type: "object",
															AdditionalProperties: &dcl.Property{
																Type: "string",
															},
															GoName:      "DriverLogLevels",
															Description: "The per-package log levels for the driver. This may include \"root\" package name to configure rootLogger. Examples: 'com.google = FATAL', 'root = INFO', 'org.apache = DEBUG'",
															Immutable:   true,
														},
													},
												},
												"mainPythonFileUri": &dcl.Property{
													Type:        "string",
													GoName:      "MainPythonFileUri",
													Description: "Required. The HCFS URI of the main Python file to use as the driver. Must be a .py file.",
													Immutable:   true,
												},
												"properties": &dcl.Property{
													Type: "object",
													AdditionalProperties: &dcl.Property{
														Type: "string",
													},
													GoName:      "Properties",
													Description: "Optional. A mapping of property names to values, used to configure PySpark. Properties that conflict with values set by the Dataproc API may be overwritten. Can include properties set in /etc/spark/conf/spark-defaults.conf and classes in user code.",
													Immutable:   true,
												},
												"pythonFileUris": &dcl.Property{
													Type:        "array",
													GoName:      "PythonFileUris",
													Description: "Optional. HCFS file URIs of Python files to pass to the PySpark framework. Supported file types: .py, .egg, and .zip.",
													Immutable:   true,
													SendEmpty:   true,
													ListType:    "list",
													Items: &dcl.Property{
														Type:   "string",
														GoType: "string",
													},
												},
											},
										},
										"scheduling": &dcl.Property{
											Type:        "object",
											GoName:      "Scheduling",
											GoType:      "WorkflowTemplateJobsScheduling",
											Description: "Optional. Job scheduling configuration.",
											Immutable:   true,
											Properties: map[string]*dcl.Property{
												"maxFailuresPerHour": &dcl.Property{
													Type:        "integer",
													Format:      "int64",
													GoName:      "MaxFailuresPerHour",
													Description: "Optional. Maximum number of times per hour a driver may be restarted as a result of driver exiting with non-zero code before job is reported failed. A job may be reported as thrashing if driver exits with non-zero code 4 times within 10 minute window. Maximum value is 10.",
													Immutable:   true,
												},
												"maxFailuresTotal": &dcl.Property{
													Type:        "integer",
													Format:      "int64",
													GoName:      "MaxFailuresTotal",
													Description: "Optional. Maximum number of times in total a driver may be restarted as a result of driver exiting with non-zero code before job is reported failed. Maximum value is 240.",
													Immutable:   true,
												},
											},
										},
										"sparkJob": &dcl.Property{
											Type:        "object",
											GoName:      "SparkJob",
											GoType:      "WorkflowTemplateJobsSparkJob",
											Description: "Optional. Job is a Spark job.",
											Immutable:   true,
											Properties: map[string]*dcl.Property{
												"archiveUris": &dcl.Property{
													Type:        "array",
													GoName:      "ArchiveUris",
													Description: "Optional. HCFS URIs of archives to be extracted into the working directory of each executor. Supported file types: .jar, .tar, .tar.gz, .tgz, and .zip.",
													Immutable:   true,
													SendEmpty:   true,
													ListType:    "list",
													Items: &dcl.Property{
														Type:   "string",
														GoType: "string",
													},
												},
												"args": &dcl.Property{
													Type:        "array",
													GoName:      "Args",
													Description: "Optional. The arguments to pass to the driver. Do not include arguments, such as `--conf`, that can be set as job properties, since a collision may occur that causes an incorrect job submission.",
													Immutable:   true,
													SendEmpty:   true,
													ListType:    "list",
													Items: &dcl.Property{
														Type:   "string",
														GoType: "string",
													},
												},
												"fileUris": &dcl.Property{
													Type:        "array",
													GoName:      "FileUris",
													Description: "Optional. HCFS URIs of files to be placed in the working directory of each executor. Useful for naively parallel tasks.",
													Immutable:   true,
													SendEmpty:   true,
													ListType:    "list",
													Items: &dcl.Property{
														Type:   "string",
														GoType: "string",
													},
												},
												"jarFileUris": &dcl.Property{
													Type:        "array",
													GoName:      "JarFileUris",
													Description: "Optional. HCFS URIs of jar files to add to the CLASSPATHs of the Spark driver and tasks.",
													Immutable:   true,
													SendEmpty:   true,
													ListType:    "list",
													Items: &dcl.Property{
														Type:   "string",
														GoType: "string",
													},
												},
												"loggingConfig": &dcl.Property{
													Type:        "object",
													GoName:      "LoggingConfig",
													GoType:      "WorkflowTemplateJobsSparkJobLoggingConfig",
													Description: "Optional. The runtime log config for job execution.",
													Immutable:   true,
													Properties: map[string]*dcl.Property{
														"driverLogLevels": &dcl.Property{
															Type: "object",
															AdditionalProperties: &dcl.Property{
																Type: "string",
															},
															GoName:      "DriverLogLevels",
															Description: "The per-package log levels for the driver. This may include \"root\" package name to configure rootLogger. Examples: 'com.google = FATAL', 'root = INFO', 'org.apache = DEBUG'",
															Immutable:   true,
														},
													},
												},
												"mainClass": &dcl.Property{
													Type:        "string",
													GoName:      "MainClass",
													Description: "The name of the driver's main class. The jar file that contains the class must be in the default CLASSPATH or specified in `jar_file_uris`.",
													Immutable:   true,
												},
												"mainJarFileUri": &dcl.Property{
													Type:        "string",
													GoName:      "MainJarFileUri",
													Description: "The HCFS URI of the jar file that contains the main class.",
													Immutable:   true,
												},
												"properties": &dcl.Property{
													Type: "object",
													AdditionalProperties: &dcl.Property{
														Type: "string",
													},
													GoName:      "Properties",
													Description: "Optional. A mapping of property names to values, used to configure Spark. Properties that conflict with values set by the Dataproc API may be overwritten. Can include properties set in /etc/spark/conf/spark-defaults.conf and classes in user code.",
													Immutable:   true,
												},
											},
										},
										"sparkRJob": &dcl.Property{
											Type:        "object",
											GoName:      "SparkRJob",
											GoType:      "WorkflowTemplateJobsSparkRJob",
											Description: "Optional. Job is a SparkR job.",
											Immutable:   true,
											Required: []string{
												"mainRFileUri",
											},
											Properties: map[string]*dcl.Property{
												"archiveUris": &dcl.Property{
													Type:        "array",
													GoName:      "ArchiveUris",
													Description: "Optional. HCFS URIs of archives to be extracted into the working directory of each executor. Supported file types: .jar, .tar, .tar.gz, .tgz, and .zip.",
													Immutable:   true,
													SendEmpty:   true,
													ListType:    "list",
													Items: &dcl.Property{
														Type:   "string",
														GoType: "string",
													},
												},
												"args": &dcl.Property{
													Type:        "array",
													GoName:      "Args",
													Description: "Optional. The arguments to pass to the driver. Do not include arguments, such as `--conf`, that can be set as job properties, since a collision may occur that causes an incorrect job submission.",
													Immutable:   true,
													SendEmpty:   true,
													ListType:    "list",
													Items: &dcl.Property{
														Type:   "string",
														GoType: "string",
													},
												},
												"fileUris": &dcl.Property{
													Type:        "array",
													GoName:      "FileUris",
													Description: "Optional. HCFS URIs of files to be placed in the working directory of each executor. Useful for naively parallel tasks.",
													Immutable:   true,
													SendEmpty:   true,
													ListType:    "list",
													Items: &dcl.Property{
														Type:   "string",
														GoType: "string",
													},
												},
												"loggingConfig": &dcl.Property{
													Type:        "object",
													GoName:      "LoggingConfig",
													GoType:      "WorkflowTemplateJobsSparkRJobLoggingConfig",
													Description: "Optional. The runtime log config for job execution.",
													Immutable:   true,
													Properties: map[string]*dcl.Property{
														"driverLogLevels": &dcl.Property{
															Type: "object",
															AdditionalProperties: &dcl.Property{
																Type: "string",
															},
															GoName:      "DriverLogLevels",
															Description: "The per-package log levels for the driver. This may include \"root\" package name to configure rootLogger. Examples: 'com.google = FATAL', 'root = INFO', 'org.apache = DEBUG'",
															Immutable:   true,
														},
													},
												},
												"mainRFileUri": &dcl.Property{
													Type:        "string",
													GoName:      "MainRFileUri",
													Description: "Required. The HCFS URI of the main R file to use as the driver. Must be a .R file.",
													Immutable:   true,
												},
												"properties": &dcl.Property{
													Type: "object",
													AdditionalProperties: &dcl.Property{
														Type: "string",
													},
													GoName:      "Properties",
													Description: "Optional. A mapping of property names to values, used to configure SparkR. Properties that conflict with values set by the Dataproc API may be overwritten. Can include properties set in /etc/spark/conf/spark-defaults.conf and classes in user code.",
													Immutable:   true,
												},
											},
										},
										"sparkSqlJob": &dcl.Property{
											Type:        "object",
											GoName:      "SparkSqlJob",
											GoType:      "WorkflowTemplateJobsSparkSqlJob",
											Description: "Optional. Job is a SparkSql job.",
											Immutable:   true,
											Properties: map[string]*dcl.Property{
												"jarFileUris": &dcl.Property{
													Type:        "array",
													GoName:      "JarFileUris",
													Description: "Optional. HCFS URIs of jar files to be added to the Spark CLASSPATH.",
													Immutable:   true,
													SendEmpty:   true,
													ListType:    "list",
													Items: &dcl.Property{
														Type:   "string",
														GoType: "string",
													},
												},
												"loggingConfig": &dcl.Property{
													Type:        "object",
													GoName:      "LoggingConfig",
													GoType:      "WorkflowTemplateJobsSparkSqlJobLoggingConfig",
													Description: "Optional. The runtime log config for job execution.",
													Immutable:   true,
													Properties: map[string]*dcl.Property{
														"driverLogLevels": &dcl.Property{
															Type: "object",
															AdditionalProperties: &dcl.Property{
																Type: "string",
															},
															GoName:      "DriverLogLevels",
															Description: "The per-package log levels for the driver. This may include \"root\" package name to configure rootLogger. Examples: 'com.google = FATAL', 'root = INFO', 'org.apache = DEBUG'",
															Immutable:   true,
														},
													},
												},
												"properties": &dcl.Property{
													Type: "object",
													AdditionalProperties: &dcl.Property{
														Type: "string",
													},
													GoName:      "Properties",
													Description: "Optional. A mapping of property names to values, used to configure Spark SQL's SparkConf. Properties that conflict with values set by the Dataproc API may be overwritten.",
													Immutable:   true,
												},
												"queryFileUri": &dcl.Property{
													Type:        "string",
													GoName:      "QueryFileUri",
													Description: "The HCFS URI of the script that contains SQL queries.",
													Immutable:   true,
												},
												"queryList": &dcl.Property{
													Type:        "object",
													GoName:      "QueryList",
													GoType:      "WorkflowTemplateJobsSparkSqlJobQueryList",
													Description: "A list of queries.",
													Immutable:   true,
													Required: []string{
														"queries",
													},
													Properties: map[string]*dcl.Property{
														"queries": &dcl.Property{
															Type:        "array",
															GoName:      "Queries",
															Description: "Required. The queries to execute. You do not need to end a query expression with a semicolon. Multiple queries can be specified in one string by separating each with a semicolon. Here is an example of a Dataproc API snippet that uses a QueryList to specify a HiveJob: \"hiveJob\": { \"queryList\": { \"queries\": [ \"query1\", \"query2\", \"query3;query4\", ] } }",
															Immutable:   true,
															SendEmpty:   true,
															ListType:    "list",
															Items: &dcl.Property{
																Type:   "string",
																GoType: "string",
															},
														},
													},
												},
												"scriptVariables": &dcl.Property{
													Type: "object",
													AdditionalProperties: &dcl.Property{
														Type: "string",
													},
													GoName:      "ScriptVariables",
													Description: "Optional. Mapping of query variable names to values (equivalent to the Spark SQL command: SET `name=\"value\";`).",
													Immutable:   true,
												},
											},
										},
										"stepId": &dcl.Property{
											Type:        "string",
											GoName:      "StepId",
											Description: "Required. The step id. The id must be unique among all jobs within the template. The step id is used as prefix for job id, as job `goog-dataproc-workflow-step-id` label, and in prerequisiteStepIds field from other steps. The id must contain only letters (a-z, A-Z), numbers (0-9), underscores (_), and hyphens (-). Cannot begin or end with underscore or hyphen. Must consist of between 3 and 50 characters.",
											Immutable:   true,
										},
									},
								},
							},
							"labels": &dcl.Property{
								Type: "object",
								AdditionalProperties: &dcl.Property{
									Type: "string",
								},
								GoName:      "Labels",
								Description: "Optional. The labels to associate with this template. These labels will be propagated to all jobs and clusters created by the workflow instance. Label **keys** must contain 1 to 63 characters, and must conform to [RFC 1035](https://www.ietf.org/rfc/rfc1035.txt). Label **values** may be empty, but, if present, must contain 1 to 63 characters, and must conform to [RFC 1035](https://www.ietf.org/rfc/rfc1035.txt). No more than 32 labels can be associated with a template.",
								Immutable:   true,
							},
							"location": &dcl.Property{
								Type:        "string",
								GoName:      "Location",
								Description: "The location for the resource",
								Immutable:   true,
							},
							"name": &dcl.Property{
								Type:        "string",
								GoName:      "Name",
								Description: "Output only. The resource name of the workflow template, as described in https://cloud.google.com/apis/design/resource_names. * For `projects.regions.workflowTemplates`, the resource name of the template has the following format: `projects/{project_id}/regions/{region}/workflowTemplates/{template_id}` * For `projects.locations.workflowTemplates`, the resource name of the template has the following format: `projects/{project_id}/locations/{location}/workflowTemplates/{template_id}`",
								Immutable:   true,
							},
							"parameters": &dcl.Property{
								Type:        "array",
								GoName:      "Parameters",
								Description: "Optional. Template parameters whose values are substituted into the template. Values for parameters must be provided when the template is instantiated.",
								Immutable:   true,
								SendEmpty:   true,
								ListType:    "list",
								Items: &dcl.Property{
									Type:   "object",
									GoType: "WorkflowTemplateParameters",
									Required: []string{
										"name",
										"fields",
									},
									Properties: map[string]*dcl.Property{
										"description": &dcl.Property{
											Type:        "string",
											GoName:      "Description",
											Description: "Optional. Brief description of the parameter. Must not exceed 1024 characters.",
											Immutable:   true,
										},
										"fields": &dcl.Property{
											Type:        "array",
											GoName:      "Fields",
											Description: "Required. Paths to all fields that the parameter replaces. A field is allowed to appear in at most one parameter's list of field paths. A field path is similar in syntax to a google.protobuf.FieldMask. For example, a field path that references the zone field of a workflow template's cluster selector would be specified as `placement.clusterSelector.zone`. Also, field paths can reference fields using the following syntax: * Values in maps can be referenced by key: * labels['key'] * placement.clusterSelector.clusterLabels['key'] * placement.managedCluster.labels['key'] * placement.clusterSelector.clusterLabels['key'] * jobs['step-id'].labels['key'] * Jobs in the jobs list can be referenced by step-id: * jobs['step-id'].hadoopJob.mainJarFileUri * jobs['step-id'].hiveJob.queryFileUri * jobs['step-id'].pySparkJob.mainPythonFileUri * jobs['step-id'].hadoopJob.jarFileUris[0] * jobs['step-id'].hadoopJob.archiveUris[0] * jobs['step-id'].hadoopJob.fileUris[0] * jobs['step-id'].pySparkJob.pythonFileUris[0] * Items in repeated fields can be referenced by a zero-based index: * jobs['step-id'].sparkJob.args[0] * Other examples: * jobs['step-id'].hadoopJob.properties['key'] * jobs['step-id'].hadoopJob.args[0] * jobs['step-id'].hiveJob.scriptVariables['key'] * jobs['step-id'].hadoopJob.mainJarFileUri * placement.clusterSelector.zone It may not be possible to parameterize maps and repeated fields in their entirety since only individual map values and individual items in repeated fields can be referenced. For example, the following field paths are invalid: - placement.clusterSelector.clusterLabels - jobs['step-id'].sparkJob.args",
											Immutable:   true,
											SendEmpty:   true,
											ListType:    "list",
											Items: &dcl.Property{
												Type:   "string",
												GoType: "string",
											},
										},
										"name": &dcl.Property{
											Type:        "string",
											GoName:      "Name",
											Description: "Required. Parameter name. The parameter name is used as the key, and paired with the parameter value, which are passed to the template when the template is instantiated. The name must contain only capital letters (A-Z), numbers (0-9), and underscores (_), and must not start with a number. The maximum length is 40 characters.",
											Immutable:   true,
										},
										"validation": &dcl.Property{
											Type:        "object",
											GoName:      "Validation",
											GoType:      "WorkflowTemplateParametersValidation",
											Description: "Optional. Validation rules to be applied to this parameter's value.",
											Immutable:   true,
											Properties: map[string]*dcl.Property{
												"regex": &dcl.Property{
													Type:        "object",
													GoName:      "Regex",
													GoType:      "WorkflowTemplateParametersValidationRegex",
													Description: "Validation based on regular expressions.",
													Immutable:   true,
													Required: []string{
														"regexes",
													},
													Properties: map[string]*dcl.Property{
														"regexes": &dcl.Property{
															Type:        "array",
															GoName:      "Regexes",
															Description: "Required. RE2 regular expressions used to validate the parameter's value. The value must match the regex in its entirety (substring matches are not sufficient).",
															Immutable:   true,
															SendEmpty:   true,
															ListType:    "list",
															Items: &dcl.Property{
																Type:   "string",
																GoType: "string",
															},
														},
													},
												},
												"values": &dcl.Property{
													Type:        "object",
													GoName:      "Values",
													GoType:      "WorkflowTemplateParametersValidationValues",
													Description: "Validation based on a list of allowed values.",
													Immutable:   true,
													Required: []string{
														"values",
													},
													Properties: map[string]*dcl.Property{
														"values": &dcl.Property{
															Type:        "array",
															GoName:      "Values",
															Description: "Required. List of allowed values for the parameter.",
															Immutable:   true,
															SendEmpty:   true,
															ListType:    "list",
															Items: &dcl.Property{
																Type:   "string",
																GoType: "string",
															},
														},
													},
												},
											},
										},
									},
								},
							},
							"placement": &dcl.Property{
								Type:        "object",
								GoName:      "Placement",
								GoType:      "WorkflowTemplatePlacement",
								Description: "Required. WorkflowTemplate scheduling information.",
								Immutable:   true,
								Properties: map[string]*dcl.Property{
									"clusterSelector": &dcl.Property{
										Type:        "object",
										GoName:      "ClusterSelector",
										GoType:      "WorkflowTemplatePlacementClusterSelector",
										Description: "Optional. A selector that chooses target cluster for jobs based on metadata. The selector is evaluated at the time each job is submitted.",
										Immutable:   true,
										Required: []string{
											"clusterLabels",
										},
										Properties: map[string]*dcl.Property{
											"clusterLabels": &dcl.Property{
												Type: "object",
												AdditionalProperties: &dcl.Property{
													Type: "string",
												},
												GoName:      "ClusterLabels",
												Description: "Required. The cluster labels. Cluster must have all labels to match.",
												Immutable:   true,
											},
											"zone": &dcl.Property{
												Type:        "string",
												GoName:      "Zone",
												Description: "Optional. The zone where workflow process executes. This parameter does not affect the selection of the cluster. If unspecified, the zone of the first cluster matching the selector is used.",
												Immutable:   true,
											},
										},
									},
									"managedCluster": &dcl.Property{
										Type:        "object",
										GoName:      "ManagedCluster",
										GoType:      "WorkflowTemplatePlacementManagedCluster",
										Description: "A cluster that is managed by the workflow.",
										Immutable:   true,
										Required: []string{
											"clusterName",
											"config",
										},
										Properties: map[string]*dcl.Property{
											"clusterName": &dcl.Property{
												Type:        "string",
												GoName:      "ClusterName",
												Description: "Required. The cluster name prefix. A unique cluster name will be formed by appending a random suffix. The name must contain only lower-case letters (a-z), numbers (0-9), and hyphens (-). Must begin with a letter. Cannot begin or end with hyphen. Must consist of between 2 and 35 characters.",
												Immutable:   true,
											},
											"config": &dcl.Property{
												Type:        "object",
												GoName:      "Config",
												GoType:      "WorkflowTemplatePlacementManagedClusterConfig",
												Description: "Required. The cluster configuration.",
												Immutable:   true,
												Properties: map[string]*dcl.Property{
													"autoscalingConfig": &dcl.Property{
														Type:        "object",
														GoName:      "AutoscalingConfig",
														GoType:      "WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig",
														Description: "Optional. Autoscaling config for the policy associated with the cluster. Cluster does not autoscale if this field is unset.",
														Immutable:   true,
														Properties: map[string]*dcl.Property{
															"policy": &dcl.Property{
																Type:        "string",
																GoName:      "Policy",
																Description: "Optional. The autoscaling policy used by the cluster. Only resource names including projectid and location (region) are valid. Examples: * `https://www.googleapis.com/compute/v1/projects/[project_id]/locations/[dataproc_region]/autoscalingPolicies/[policy_id]` * `projects/[project_id]/locations/[dataproc_region]/autoscalingPolicies/[policy_id]` Note that the policy must be in the same project and Dataproc region.",
																Immutable:   true,
																ResourceReferences: []*dcl.PropertyResourceReference{
																	&dcl.PropertyResourceReference{
																		Resource: "Dataproc/AutoscalingPolicy",
																		Field:    "name",
																	},
																},
															},
														},
													},
													"encryptionConfig": &dcl.Property{
														Type:        "object",
														GoName:      "EncryptionConfig",
														GoType:      "WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig",
														Description: "Optional. Encryption settings for the cluster.",
														Immutable:   true,
														Properties: map[string]*dcl.Property{
															"gcePdKmsKeyName": &dcl.Property{
																Type:        "string",
																GoName:      "GcePdKmsKeyName",
																Description: "Optional. The Cloud KMS key name to use for PD disk encryption for all instances in the cluster.",
																Immutable:   true,
																ResourceReferences: []*dcl.PropertyResourceReference{
																	&dcl.PropertyResourceReference{
																		Resource: "Cloudkms/CryptoKey",
																		Field:    "selfLink",
																	},
																},
															},
														},
													},
													"endpointConfig": &dcl.Property{
														Type:        "object",
														GoName:      "EndpointConfig",
														GoType:      "WorkflowTemplatePlacementManagedClusterConfigEndpointConfig",
														Description: "Optional. Port/endpoint configuration for this cluster",
														Immutable:   true,
														Properties: map[string]*dcl.Property{
															"enableHttpPortAccess": &dcl.Property{
																Type:        "boolean",
																GoName:      "EnableHttpPortAccess",
																Description: "Optional. If true, enable http access to specific ports on the cluster from external sources. Defaults to false.",
																Immutable:   true,
															},
															"httpPorts": &dcl.Property{
																Type: "object",
																AdditionalProperties: &dcl.Property{
																	Type: "string",
																},
																GoName:      "HttpPorts",
																ReadOnly:    true,
																Description: "Output only. The map of port descriptions to URLs. Will only be populated if enable_http_port_access is true.",
																Immutable:   true,
															},
														},
													},
													"gceClusterConfig": &dcl.Property{
														Type:        "object",
														GoName:      "GceClusterConfig",
														GoType:      "WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig",
														Description: "Optional. The shared Compute Engine config settings for all instances in a cluster.",
														Immutable:   true,
														Properties: map[string]*dcl.Property{
															"internalIPOnly": &dcl.Property{
																Type:          "boolean",
																GoName:        "InternalIPOnly",
																Description:   "Optional. If true, all instances in the cluster will only have internal IP addresses. By default, clusters are not restricted to internal IP addresses, and will have ephemeral external IP addresses assigned to each instance. This `internal_ip_only` restriction can only be enabled for subnetwork enabled networks, and all off-cluster dependencies must be configured to be accessible without external IP addresses.",
																Immutable:     true,
																ServerDefault: true,
															},
															"metadata": &dcl.Property{
																Type: "object",
																AdditionalProperties: &dcl.Property{
																	Type: "string",
																},
																GoName:      "Metadata",
																Description: "The Compute Engine metadata entries to add to all instances (see [Project and instance metadata](https://cloud.google.com/compute/docs/storing-retrieving-metadata#project_and_instance_metadata)).",
																Immutable:   true,
															},
															"network": &dcl.Property{
																Type:        "string",
																GoName:      "Network",
																Description: "Optional. The Compute Engine network to be used for machine communications. Cannot be specified with subnetwork_uri. If neither `network_uri` nor `subnetwork_uri` is specified, the \"default\" network of the project is used, if it exists. Cannot be a \"Custom Subnet Network\" (see [Using Subnetworks](https://cloud.google.com/compute/docs/subnetworks) for more information). A full URL, partial URI, or short name are valid. Examples: * `https://www.googleapis.com/compute/v1/projects/[project_id]/regions/global/default` * `projects/[project_id]/regions/global/default` * `default`",
																Immutable:   true,
																ResourceReferences: []*dcl.PropertyResourceReference{
																	&dcl.PropertyResourceReference{
																		Resource: "Compute/Network",
																		Field:    "selfLink",
																	},
																},
															},
															"nodeGroupAffinity": &dcl.Property{
																Type:        "object",
																GoName:      "NodeGroupAffinity",
																GoType:      "WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity",
																Description: "Optional. Node Group Affinity for sole-tenant clusters.",
																Immutable:   true,
																Required: []string{
																	"nodeGroup",
																},
																Properties: map[string]*dcl.Property{
																	"nodeGroup": &dcl.Property{
																		Type:        "string",
																		GoName:      "NodeGroup",
																		Description: "Required. The URI of a sole-tenant [node group resource](https://cloud.google.com/compute/docs/reference/rest/v1/nodeGroups) that the cluster will be created on. A full URL, partial URI, or node group name are valid. Examples: * `https://www.googleapis.com/compute/v1/projects/[project_id]/zones/us-central1-a/nodeGroups/node-group-1` * `projects/[project_id]/zones/us-central1-a/nodeGroups/node-group-1` * `node-group-1`",
																		Immutable:   true,
																		ResourceReferences: []*dcl.PropertyResourceReference{
																			&dcl.PropertyResourceReference{
																				Resource: "Compute/NodeGroup",
																				Field:    "selfLink",
																			},
																		},
																	},
																},
															},
															"privateIPv6GoogleAccess": &dcl.Property{
																Type:        "string",
																GoName:      "PrivateIPv6GoogleAccess",
																GoType:      "WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum",
																Description: "Optional. The type of IPv6 access for a cluster. Possible values: PRIVATE_IPV6_GOOGLE_ACCESS_UNSPECIFIED, INHERIT_FROM_SUBNETWORK, OUTBOUND, BIDIRECTIONAL",
																Immutable:   true,
																Enum: []string{
																	"PRIVATE_IPV6_GOOGLE_ACCESS_UNSPECIFIED",
																	"INHERIT_FROM_SUBNETWORK",
																	"OUTBOUND",
																	"BIDIRECTIONAL",
																},
															},
															"reservationAffinity": &dcl.Property{
																Type:        "object",
																GoName:      "ReservationAffinity",
																GoType:      "WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity",
																Description: "Optional. Reservation Affinity for consuming Zonal reservation.",
																Immutable:   true,
																Properties: map[string]*dcl.Property{
																	"consumeReservationType": &dcl.Property{
																		Type:        "string",
																		GoName:      "ConsumeReservationType",
																		GoType:      "WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum",
																		Description: "Optional. Type of reservation to consume Possible values: TYPE_UNSPECIFIED, NO_RESERVATION, ANY_RESERVATION, SPECIFIC_RESERVATION",
																		Immutable:   true,
																		Enum: []string{
																			"TYPE_UNSPECIFIED",
																			"NO_RESERVATION",
																			"ANY_RESERVATION",
																			"SPECIFIC_RESERVATION",
																		},
																	},
																	"key": &dcl.Property{
																		Type:        "string",
																		GoName:      "Key",
																		Description: "Optional. Corresponds to the label key of reservation resource.",
																		Immutable:   true,
																	},
																	"values": &dcl.Property{
																		Type:        "array",
																		GoName:      "Values",
																		Description: "Optional. Corresponds to the label values of reservation resource.",
																		Immutable:   true,
																		SendEmpty:   true,
																		ListType:    "list",
																		Items: &dcl.Property{
																			Type:   "string",
																			GoType: "string",
																		},
																	},
																},
															},
															"serviceAccount": &dcl.Property{
																Type:        "string",
																GoName:      "ServiceAccount",
																Description: "Optional. The [Dataproc service account](https://cloud.google.com/dataproc/docs/concepts/configuring-clusters/service-accounts#service_accounts_in_dataproc) (also see [VM Data Plane identity](https://cloud.google.com/dataproc/docs/concepts/iam/dataproc-principals#vm_service_account_data_plane_identity)) used by Dataproc cluster VM instances to access Google Cloud Platform services. If not specified, the [Compute Engine default service account](https://cloud.google.com/compute/docs/access/service-accounts#default_service_account) is used.",
																Immutable:   true,
																ResourceReferences: []*dcl.PropertyResourceReference{
																	&dcl.PropertyResourceReference{
																		Resource: "Iam/ServiceAccount",
																		Field:    "email",
																	},
																},
															},
															"serviceAccountScopes": &dcl.Property{
																Type:        "array",
																GoName:      "ServiceAccountScopes",
																Description: "Optional. The URIs of service account scopes to be included in Compute Engine instances. The following base set of scopes is always included: * https://www.googleapis.com/auth/cloud.useraccounts.readonly * https://www.googleapis.com/auth/devstorage.read_write * https://www.googleapis.com/auth/logging.write If no scopes are specified, the following defaults are also provided: * https://www.googleapis.com/auth/bigquery * https://www.googleapis.com/auth/bigtable.admin.table * https://www.googleapis.com/auth/bigtable.data * https://www.googleapis.com/auth/devstorage.full_control",
																Immutable:   true,
																SendEmpty:   true,
																ListType:    "list",
																Items: &dcl.Property{
																	Type:   "string",
																	GoType: "string",
																},
															},
															"shieldedInstanceConfig": &dcl.Property{
																Type:        "object",
																GoName:      "ShieldedInstanceConfig",
																GoType:      "WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig",
																Description: "Optional. Shielded Instance Config for clusters using Compute Engine Shielded VMs.",
																Immutable:   true,
																Properties: map[string]*dcl.Property{
																	"enableIntegrityMonitoring": &dcl.Property{
																		Type:        "boolean",
																		GoName:      "EnableIntegrityMonitoring",
																		Description: "Optional. Defines whether instances have integrity monitoring enabled. Integrity monitoring compares the most recent boot measurements to the integrity policy baseline and returns a pair of pass/fail results depending on whether they match or not.",
																		Immutable:   true,
																	},
																	"enableSecureBoot": &dcl.Property{
																		Type:        "boolean",
																		GoName:      "EnableSecureBoot",
																		Description: "Optional. Defines whether the instances have Secure Boot enabled. Secure Boot helps ensure that the system only runs authentic software by verifying the digital signature of all boot components, and halting the boot process if signature verification fails.",
																		Immutable:   true,
																	},
																	"enableVtpm": &dcl.Property{
																		Type:        "boolean",
																		GoName:      "EnableVtpm",
																		Description: "Optional. Defines whether the instance have the vTPM enabled. Virtual Trusted Platform Module protects objects like keys, certificates and enables Measured Boot by performing the measurements needed to create a known good boot baseline, called the integrity policy baseline.",
																		Immutable:   true,
																	},
																},
															},
															"subnetwork": &dcl.Property{
																Type:        "string",
																GoName:      "Subnetwork",
																Description: "Optional. The Compute Engine subnetwork to be used for machine communications. Cannot be specified with network_uri. A full URL, partial URI, or short name are valid. Examples: * `https://www.googleapis.com/compute/v1/projects/[project_id]/regions/us-east1/subnetworks/sub0` * `projects/[project_id]/regions/us-east1/subnetworks/sub0` * `sub0`",
																Immutable:   true,
																ResourceReferences: []*dcl.PropertyResourceReference{
																	&dcl.PropertyResourceReference{
																		Resource: "Compute/Subnetwork",
																		Field:    "selfLink",
																	},
																},
															},
															"tags": &dcl.Property{
																Type:        "array",
																GoName:      "Tags",
																Description: "The Compute Engine tags to add to all instances (see [Tagging instances](https://cloud.google.com/compute/docs/label-or-tag-resources#tags)).",
																Immutable:   true,
																SendEmpty:   true,
																ListType:    "set",
																Items: &dcl.Property{
																	Type:   "string",
																	GoType: "string",
																},
															},
															"zone": &dcl.Property{
																Type:        "string",
																GoName:      "Zone",
																Description: "Optional. The zone where the Compute Engine cluster will be located. On a create request, it is required in the \"global\" region. If omitted in a non-global Dataproc region, the service will pick a zone in the corresponding Compute Engine region. On a get request, zone will always be present. A full URL, partial URI, or short name are valid. Examples: * `https://www.googleapis.com/compute/v1/projects/[project_id]/zones/[zone]` * `projects/[project_id]/zones/[zone]` * `us-central1-f`",
																Immutable:   true,
															},
														},
													},
													"initializationActions": &dcl.Property{
														Type:        "array",
														GoName:      "InitializationActions",
														Description: "Optional. Commands to execute on each node after config is completed. By default, executables are run on master and all worker nodes. You can test a node's `role` metadata to run an executable on a master or worker node, as shown below using `curl` (you can also use `wget`): ROLE=$(curl -H Metadata-Flavor:Google http://metadata/computeMetadata/v1/instance/attributes/dataproc-role) if [[ \"${ROLE}\" == 'Master' ]]; then ... master specific actions ... else ... worker specific actions ... fi",
														Immutable:   true,
														SendEmpty:   true,
														ListType:    "list",
														Items: &dcl.Property{
															Type:   "object",
															GoType: "WorkflowTemplatePlacementManagedClusterConfigInitializationActions",
															Properties: map[string]*dcl.Property{
																"executableFile": &dcl.Property{
																	Type:        "string",
																	GoName:      "ExecutableFile",
																	Description: "Required. Cloud Storage URI of executable file.",
																	Immutable:   true,
																},
																"executionTimeout": &dcl.Property{
																	Type:        "string",
																	GoName:      "ExecutionTimeout",
																	Description: "Optional. Amount of time executable has to complete. Default is 10 minutes (see JSON representation of [Duration](https://developers.google.com/protocol-buffers/docs/proto3#json)). Cluster creation fails with an explanatory error message (the name of the executable that caused the error and the exceeded timeout period) if the executable is not completed at end of the timeout period.",
																	Immutable:   true,
																},
															},
														},
													},
													"lifecycleConfig": &dcl.Property{
														Type:        "object",
														GoName:      "LifecycleConfig",
														GoType:      "WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig",
														Description: "Optional. Lifecycle setting for the cluster.",
														Immutable:   true,
														Properties: map[string]*dcl.Property{
															"autoDeleteTime": &dcl.Property{
																Type:        "string",
																Format:      "date-time",
																GoName:      "AutoDeleteTime",
																Description: "Optional. The time when cluster will be auto-deleted (see JSON representation of [Timestamp](https://developers.google.com/protocol-buffers/docs/proto3#json)).",
																Immutable:   true,
															},
															"autoDeleteTtl": &dcl.Property{
																Type:        "string",
																GoName:      "AutoDeleteTtl",
																Description: "Optional. The lifetime duration of cluster. The cluster will be auto-deleted at the end of this period. Minimum value is 10 minutes; maximum value is 14 days (see JSON representation of [Duration](https://developers.google.com/protocol-buffers/docs/proto3#json)).",
																Immutable:   true,
															},
															"idleDeleteTtl": &dcl.Property{
																Type:        "string",
																GoName:      "IdleDeleteTtl",
																Description: "Optional. The duration to keep the cluster alive while idling (when no jobs are running). Passing this threshold will cause the cluster to be deleted. Minimum value is 5 minutes; maximum value is 14 days (see JSON representation of [Duration](https://developers.google.com/protocol-buffers/docs/proto3#json)).",
																Immutable:   true,
															},
															"idleStartTime": &dcl.Property{
																Type:        "string",
																Format:      "date-time",
																GoName:      "IdleStartTime",
																ReadOnly:    true,
																Description: "Output only. The time when cluster became idle (most recent job finished) and became eligible for deletion due to idleness (see JSON representation of [Timestamp](https://developers.google.com/protocol-buffers/docs/proto3#json)).",
																Immutable:   true,
															},
														},
													},
													"masterConfig": &dcl.Property{
														Type:          "object",
														GoName:        "MasterConfig",
														GoType:        "WorkflowTemplatePlacementManagedClusterConfigMasterConfig",
														Description:   "Optional. The Compute Engine config settings for the master instance in a cluster.",
														Immutable:     true,
														ServerDefault: true,
														Properties: map[string]*dcl.Property{
															"accelerators": &dcl.Property{
																Type:          "array",
																GoName:        "Accelerators",
																Description:   "Optional. The Compute Engine accelerator configuration for these instances.",
																Immutable:     true,
																ServerDefault: true,
																SendEmpty:     true,
																ListType:      "list",
																Items: &dcl.Property{
																	Type:   "object",
																	GoType: "WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators",
																	Properties: map[string]*dcl.Property{
																		"acceleratorCount": &dcl.Property{
																			Type:        "integer",
																			Format:      "int64",
																			GoName:      "AcceleratorCount",
																			Description: "The number of the accelerator cards of this type exposed to this instance.",
																			Immutable:   true,
																		},
																		"acceleratorType": &dcl.Property{
																			Type:        "string",
																			GoName:      "AcceleratorType",
																			Description: "Full URL, partial URI, or short name of the accelerator type resource to expose to this instance. See [Compute Engine AcceleratorTypes](https://cloud.google.com/compute/docs/reference/beta/acceleratorTypes). Examples: * `https://www.googleapis.com/compute/beta/projects/[project_id]/zones/us-east1-a/acceleratorTypes/nvidia-tesla-k80` * `projects/[project_id]/zones/us-east1-a/acceleratorTypes/nvidia-tesla-k80` * `nvidia-tesla-k80` **Auto Zone Exception**: If you are using the Dataproc [Auto Zone Placement](https://cloud.google.com/dataproc/docs/concepts/configuring-clusters/auto-zone#using_auto_zone_placement) feature, you must use the short name of the accelerator type resource, for example, `nvidia-tesla-k80`.",
																			Immutable:   true,
																		},
																	},
																},
															},
															"diskConfig": &dcl.Property{
																Type:          "object",
																GoName:        "DiskConfig",
																GoType:        "WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig",
																Description:   "Optional. Disk option config settings.",
																Immutable:     true,
																ServerDefault: true,
																Properties: map[string]*dcl.Property{
																	"bootDiskSizeGb": &dcl.Property{
																		Type:        "integer",
																		Format:      "int64",
																		GoName:      "BootDiskSizeGb",
																		Description: "Optional. Size in GB of the boot disk (default is 500GB).",
																		Immutable:   true,
																	},
																	"bootDiskType": &dcl.Property{
																		Type:        "string",
																		GoName:      "BootDiskType",
																		Description: "Optional. Type of the boot disk (default is \"pd-standard\"). Valid values: \"pd-balanced\" (Persistent Disk Balanced Solid State Drive), \"pd-ssd\" (Persistent Disk Solid State Drive), or \"pd-standard\" (Persistent Disk Hard Disk Drive). See [Disk types](https://cloud.google.com/compute/docs/disks#disk-types).",
																		Immutable:   true,
																	},
																	"numLocalSsds": &dcl.Property{
																		Type:          "integer",
																		Format:        "int64",
																		GoName:        "NumLocalSsds",
																		Description:   "Optional. Number of attached SSDs, from 0 to 4 (default is 0). If SSDs are not attached, the boot disk is used to store runtime logs and [HDFS](https://hadoop.apache.org/docs/r1.2.1/hdfs_user_guide.html) data. If one or more SSDs are attached, this runtime bulk data is spread across them, and the boot disk contains only basic config and installed binaries.",
																		Immutable:     true,
																		ServerDefault: true,
																	},
																},
															},
															"image": &dcl.Property{
																Type:        "string",
																GoName:      "Image",
																Description: "Optional. The Compute Engine image resource used for cluster instances. The URI can represent an image or image family. Image examples: * `https://www.googleapis.com/compute/beta/projects/[project_id]/global/images/[image-id]` * `projects/[project_id]/global/images/[image-id]` * `image-id` Image family examples. Dataproc will use the most recent image from the family: * `https://www.googleapis.com/compute/beta/projects/[project_id]/global/images/family/[custom-image-family-name]` * `projects/[project_id]/global/images/family/[custom-image-family-name]` If the URI is unspecified, it will be inferred from `SoftwareConfig.image_version` or the system default.",
																Immutable:   true,
																ResourceReferences: []*dcl.PropertyResourceReference{
																	&dcl.PropertyResourceReference{
																		Resource: "Compute/Image",
																		Field:    "selfLink",
																	},
																},
															},
															"instanceNames": &dcl.Property{
																Type:          "array",
																GoName:        "InstanceNames",
																ReadOnly:      true,
																Description:   "Output only. The list of instance names. Dataproc derives the names from `cluster_name`, `num_instances`, and the instance group.",
																Immutable:     true,
																ServerDefault: true,
																ListType:      "list",
																Items: &dcl.Property{
																	Type:   "string",
																	GoType: "string",
																	ResourceReferences: []*dcl.PropertyResourceReference{
																		&dcl.PropertyResourceReference{
																			Resource: "Compute/Instance",
																			Field:    "selfLink",
																		},
																	},
																},
															},
															"isPreemptible": &dcl.Property{
																Type:        "boolean",
																GoName:      "IsPreemptible",
																ReadOnly:    true,
																Description: "Output only. Specifies that this instance group contains preemptible instances.",
																Immutable:   true,
															},
															"machineType": &dcl.Property{
																Type:        "string",
																GoName:      "MachineType",
																Description: "Optional. The Compute Engine machine type used for cluster instances. A full URL, partial URI, or short name are valid. Examples: * `https://www.googleapis.com/compute/v1/projects/[project_id]/zones/us-east1-a/machineTypes/n1-standard-2` * `projects/[project_id]/zones/us-east1-a/machineTypes/n1-standard-2` * `n1-standard-2` **Auto Zone Exception**: If you are using the Dataproc [Auto Zone Placement](https://cloud.google.com/dataproc/docs/concepts/configuring-clusters/auto-zone#using_auto_zone_placement) feature, you must use the short name of the machine type resource, for example, `n1-standard-2`.",
																Immutable:   true,
															},
															"managedGroupConfig": &dcl.Property{
																Type:          "object",
																GoName:        "ManagedGroupConfig",
																GoType:        "WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig",
																ReadOnly:      true,
																Description:   "Output only. The config for Compute Engine Instance Group Manager that manages this group. This is only used for preemptible instance groups.",
																Immutable:     true,
																ServerDefault: true,
																Properties: map[string]*dcl.Property{
																	"instanceGroupManagerName": &dcl.Property{
																		Type:        "string",
																		GoName:      "InstanceGroupManagerName",
																		ReadOnly:    true,
																		Description: "Output only. The name of the Instance Group Manager for this group.",
																		Immutable:   true,
																	},
																	"instanceTemplateName": &dcl.Property{
																		Type:        "string",
																		GoName:      "InstanceTemplateName",
																		ReadOnly:    true,
																		Description: "Output only. The name of the Instance Template used for the Managed Instance Group.",
																		Immutable:   true,
																	},
																},
															},
															"minCpuPlatform": &dcl.Property{
																Type:          "string",
																GoName:        "MinCpuPlatform",
																Description:   "Optional. Specifies the minimum cpu platform for the Instance Group. See [Dataproc -> Minimum CPU Platform](https://cloud.google.com/dataproc/docs/concepts/compute/dataproc-min-cpu).",
																Immutable:     true,
																ServerDefault: true,
															},
															"numInstances": &dcl.Property{
																Type:        "integer",
																Format:      "int64",
																GoName:      "NumInstances",
																Description: "Optional. The number of VM instances in the instance group. For [HA cluster](/dataproc/docs/concepts/configuring-clusters/high-availability) [master_config](#FIELDS.master_config) groups, **must be set to 3**. For standard cluster [master_config](#FIELDS.master_config) groups, **must be set to 1**.",
																Immutable:   true,
															},
															"preemptibility": &dcl.Property{
																Type:        "string",
																GoName:      "Preemptibility",
																GoType:      "WorkflowTemplatePlacementManagedClusterConfigMasterConfigPreemptibilityEnum",
																Description: "Optional. Specifies the preemptibility of the instance group. The default value for master and worker groups is `NON_PREEMPTIBLE`. This default cannot be changed. The default value for secondary instances is `PREEMPTIBLE`. Possible values: PREEMPTIBILITY_UNSPECIFIED, NON_PREEMPTIBLE, PREEMPTIBLE",
																Immutable:   true,
																Enum: []string{
																	"PREEMPTIBILITY_UNSPECIFIED",
																	"NON_PREEMPTIBLE",
																	"PREEMPTIBLE",
																},
															},
														},
													},
													"secondaryWorkerConfig": &dcl.Property{
														Type:          "object",
														GoName:        "SecondaryWorkerConfig",
														GoType:        "WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig",
														Description:   "Optional. The Compute Engine config settings for additional worker instances in a cluster.",
														Immutable:     true,
														ServerDefault: true,
														Properties: map[string]*dcl.Property{
															"accelerators": &dcl.Property{
																Type:          "array",
																GoName:        "Accelerators",
																Description:   "Optional. The Compute Engine accelerator configuration for these instances.",
																Immutable:     true,
																ServerDefault: true,
																SendEmpty:     true,
																ListType:      "list",
																Items: &dcl.Property{
																	Type:   "object",
																	GoType: "WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators",
																	Properties: map[string]*dcl.Property{
																		"acceleratorCount": &dcl.Property{
																			Type:        "integer",
																			Format:      "int64",
																			GoName:      "AcceleratorCount",
																			Description: "The number of the accelerator cards of this type exposed to this instance.",
																			Immutable:   true,
																		},
																		"acceleratorType": &dcl.Property{
																			Type:        "string",
																			GoName:      "AcceleratorType",
																			Description: "Full URL, partial URI, or short name of the accelerator type resource to expose to this instance. See [Compute Engine AcceleratorTypes](https://cloud.google.com/compute/docs/reference/beta/acceleratorTypes). Examples: * `https://www.googleapis.com/compute/beta/projects/[project_id]/zones/us-east1-a/acceleratorTypes/nvidia-tesla-k80` * `projects/[project_id]/zones/us-east1-a/acceleratorTypes/nvidia-tesla-k80` * `nvidia-tesla-k80` **Auto Zone Exception**: If you are using the Dataproc [Auto Zone Placement](https://cloud.google.com/dataproc/docs/concepts/configuring-clusters/auto-zone#using_auto_zone_placement) feature, you must use the short name of the accelerator type resource, for example, `nvidia-tesla-k80`.",
																			Immutable:   true,
																		},
																	},
																},
															},
															"diskConfig": &dcl.Property{
																Type:          "object",
																GoName:        "DiskConfig",
																GoType:        "WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig",
																Description:   "Optional. Disk option config settings.",
																Immutable:     true,
																ServerDefault: true,
																Properties: map[string]*dcl.Property{
																	"bootDiskSizeGb": &dcl.Property{
																		Type:        "integer",
																		Format:      "int64",
																		GoName:      "BootDiskSizeGb",
																		Description: "Optional. Size in GB of the boot disk (default is 500GB).",
																		Immutable:   true,
																	},
																	"bootDiskType": &dcl.Property{
																		Type:        "string",
																		GoName:      "BootDiskType",
																		Description: "Optional. Type of the boot disk (default is \"pd-standard\"). Valid values: \"pd-balanced\" (Persistent Disk Balanced Solid State Drive), \"pd-ssd\" (Persistent Disk Solid State Drive), or \"pd-standard\" (Persistent Disk Hard Disk Drive). See [Disk types](https://cloud.google.com/compute/docs/disks#disk-types).",
																		Immutable:   true,
																	},
																	"numLocalSsds": &dcl.Property{
																		Type:          "integer",
																		Format:        "int64",
																		GoName:        "NumLocalSsds",
																		Description:   "Optional. Number of attached SSDs, from 0 to 4 (default is 0). If SSDs are not attached, the boot disk is used to store runtime logs and [HDFS](https://hadoop.apache.org/docs/r1.2.1/hdfs_user_guide.html) data. If one or more SSDs are attached, this runtime bulk data is spread across them, and the boot disk contains only basic config and installed binaries.",
																		Immutable:     true,
																		ServerDefault: true,
																	},
																},
															},
															"image": &dcl.Property{
																Type:        "string",
																GoName:      "Image",
																Description: "Optional. The Compute Engine image resource used for cluster instances. The URI can represent an image or image family. Image examples: * `https://www.googleapis.com/compute/beta/projects/[project_id]/global/images/[image-id]` * `projects/[project_id]/global/images/[image-id]` * `image-id` Image family examples. Dataproc will use the most recent image from the family: * `https://www.googleapis.com/compute/beta/projects/[project_id]/global/images/family/[custom-image-family-name]` * `projects/[project_id]/global/images/family/[custom-image-family-name]` If the URI is unspecified, it will be inferred from `SoftwareConfig.image_version` or the system default.",
																Immutable:   true,
																ResourceReferences: []*dcl.PropertyResourceReference{
																	&dcl.PropertyResourceReference{
																		Resource: "Compute/Image",
																		Field:    "selfLink",
																	},
																},
															},
															"instanceNames": &dcl.Property{
																Type:          "array",
																GoName:        "InstanceNames",
																ReadOnly:      true,
																Description:   "Output only. The list of instance names. Dataproc derives the names from `cluster_name`, `num_instances`, and the instance group.",
																Immutable:     true,
																ServerDefault: true,
																ListType:      "list",
																Items: &dcl.Property{
																	Type:   "string",
																	GoType: "string",
																	ResourceReferences: []*dcl.PropertyResourceReference{
																		&dcl.PropertyResourceReference{
																			Resource: "Compute/Instance",
																			Field:    "selfLink",
																		},
																	},
																},
															},
															"isPreemptible": &dcl.Property{
																Type:        "boolean",
																GoName:      "IsPreemptible",
																ReadOnly:    true,
																Description: "Output only. Specifies that this instance group contains preemptible instances.",
																Immutable:   true,
															},
															"machineType": &dcl.Property{
																Type:        "string",
																GoName:      "MachineType",
																Description: "Optional. The Compute Engine machine type used for cluster instances. A full URL, partial URI, or short name are valid. Examples: * `https://www.googleapis.com/compute/v1/projects/[project_id]/zones/us-east1-a/machineTypes/n1-standard-2` * `projects/[project_id]/zones/us-east1-a/machineTypes/n1-standard-2` * `n1-standard-2` **Auto Zone Exception**: If you are using the Dataproc [Auto Zone Placement](https://cloud.google.com/dataproc/docs/concepts/configuring-clusters/auto-zone#using_auto_zone_placement) feature, you must use the short name of the machine type resource, for example, `n1-standard-2`.",
																Immutable:   true,
															},
															"managedGroupConfig": &dcl.Property{
																Type:          "object",
																GoName:        "ManagedGroupConfig",
																GoType:        "WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig",
																ReadOnly:      true,
																Description:   "Output only. The config for Compute Engine Instance Group Manager that manages this group. This is only used for preemptible instance groups.",
																Immutable:     true,
																ServerDefault: true,
																Properties: map[string]*dcl.Property{
																	"instanceGroupManagerName": &dcl.Property{
																		Type:        "string",
																		GoName:      "InstanceGroupManagerName",
																		ReadOnly:    true,
																		Description: "Output only. The name of the Instance Group Manager for this group.",
																		Immutable:   true,
																	},
																	"instanceTemplateName": &dcl.Property{
																		Type:        "string",
																		GoName:      "InstanceTemplateName",
																		ReadOnly:    true,
																		Description: "Output only. The name of the Instance Template used for the Managed Instance Group.",
																		Immutable:   true,
																	},
																},
															},
															"minCpuPlatform": &dcl.Property{
																Type:          "string",
																GoName:        "MinCpuPlatform",
																Description:   "Optional. Specifies the minimum cpu platform for the Instance Group. See [Dataproc -> Minimum CPU Platform](https://cloud.google.com/dataproc/docs/concepts/compute/dataproc-min-cpu).",
																Immutable:     true,
																ServerDefault: true,
															},
															"numInstances": &dcl.Property{
																Type:        "integer",
																Format:      "int64",
																GoName:      "NumInstances",
																Description: "Optional. The number of VM instances in the instance group. For [HA cluster](/dataproc/docs/concepts/configuring-clusters/high-availability) [master_config](#FIELDS.master_config) groups, **must be set to 3**. For standard cluster [master_config](#FIELDS.master_config) groups, **must be set to 1**.",
																Immutable:   true,
															},
															"preemptibility": &dcl.Property{
																Type:        "string",
																GoName:      "Preemptibility",
																GoType:      "WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigPreemptibilityEnum",
																Description: "Optional. Specifies the preemptibility of the instance group. The default value for master and worker groups is `NON_PREEMPTIBLE`. This default cannot be changed. The default value for secondary instances is `PREEMPTIBLE`. Possible values: PREEMPTIBILITY_UNSPECIFIED, NON_PREEMPTIBLE, PREEMPTIBLE",
																Immutable:   true,
																Enum: []string{
																	"PREEMPTIBILITY_UNSPECIFIED",
																	"NON_PREEMPTIBLE",
																	"PREEMPTIBLE",
																},
															},
														},
													},
													"securityConfig": &dcl.Property{
														Type:        "object",
														GoName:      "SecurityConfig",
														GoType:      "WorkflowTemplatePlacementManagedClusterConfigSecurityConfig",
														Description: "Optional. Security settings for the cluster.",
														Immutable:   true,
														Properties: map[string]*dcl.Property{
															"kerberosConfig": &dcl.Property{
																Type:        "object",
																GoName:      "KerberosConfig",
																GoType:      "WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig",
																Description: "Optional. Kerberos related configuration.",
																Immutable:   true,
																Properties: map[string]*dcl.Property{
																	"crossRealmTrustAdminServer": &dcl.Property{
																		Type:        "string",
																		GoName:      "CrossRealmTrustAdminServer",
																		Description: "Optional. The admin server (IP or hostname) for the remote trusted realm in a cross realm trust relationship.",
																		Immutable:   true,
																	},
																	"crossRealmTrustKdc": &dcl.Property{
																		Type:        "string",
																		GoName:      "CrossRealmTrustKdc",
																		Description: "Optional. The KDC (IP or hostname) for the remote trusted realm in a cross realm trust relationship.",
																		Immutable:   true,
																	},
																	"crossRealmTrustRealm": &dcl.Property{
																		Type:        "string",
																		GoName:      "CrossRealmTrustRealm",
																		Description: "Optional. The remote realm the Dataproc on-cluster KDC will trust, should the user enable cross realm trust.",
																		Immutable:   true,
																	},
																	"crossRealmTrustSharedPassword": &dcl.Property{
																		Type:        "string",
																		GoName:      "CrossRealmTrustSharedPassword",
																		Description: "Optional. The Cloud Storage URI of a KMS encrypted file containing the shared password between the on-cluster Kerberos realm and the remote trusted realm, in a cross realm trust relationship.",
																		Immutable:   true,
																	},
																	"enableKerberos": &dcl.Property{
																		Type:        "boolean",
																		GoName:      "EnableKerberos",
																		Description: "Optional. Flag to indicate whether to Kerberize the cluster (default: false). Set this field to true to enable Kerberos on a cluster.",
																		Immutable:   true,
																	},
																	"kdcDbKey": &dcl.Property{
																		Type:        "string",
																		GoName:      "KdcDbKey",
																		Description: "Optional. The Cloud Storage URI of a KMS encrypted file containing the master key of the KDC database.",
																		Immutable:   true,
																	},
																	"keyPassword": &dcl.Property{
																		Type:        "string",
																		GoName:      "KeyPassword",
																		Description: "Optional. The Cloud Storage URI of a KMS encrypted file containing the password to the user provided key. For the self-signed certificate, this password is generated by Dataproc.",
																		Immutable:   true,
																	},
																	"keystore": &dcl.Property{
																		Type:        "string",
																		GoName:      "Keystore",
																		Description: "Optional. The Cloud Storage URI of the keystore file used for SSL encryption. If not provided, Dataproc will provide a self-signed certificate.",
																		Immutable:   true,
																	},
																	"keystorePassword": &dcl.Property{
																		Type:        "string",
																		GoName:      "KeystorePassword",
																		Description: "Optional. The Cloud Storage URI of a KMS encrypted file containing the password to the user provided keystore. For the self-signed certificate, this password is generated by Dataproc.",
																		Immutable:   true,
																	},
																	"kmsKey": &dcl.Property{
																		Type:        "string",
																		GoName:      "KmsKey",
																		Description: "Optional. The uri of the KMS key used to encrypt various sensitive files.",
																		Immutable:   true,
																		ResourceReferences: []*dcl.PropertyResourceReference{
																			&dcl.PropertyResourceReference{
																				Resource: "Cloudkms/CryptoKey",
																				Field:    "selfLink",
																			},
																		},
																	},
																	"realm": &dcl.Property{
																		Type:        "string",
																		GoName:      "Realm",
																		Description: "Optional. The name of the on-cluster Kerberos realm. If not specified, the uppercased domain of hostnames will be the realm.",
																		Immutable:   true,
																	},
																	"rootPrincipalPassword": &dcl.Property{
																		Type:        "string",
																		GoName:      "RootPrincipalPassword",
																		Description: "Optional. The Cloud Storage URI of a KMS encrypted file containing the root principal password.",
																		Immutable:   true,
																	},
																	"tgtLifetimeHours": &dcl.Property{
																		Type:        "integer",
																		Format:      "int64",
																		GoName:      "TgtLifetimeHours",
																		Description: "Optional. The lifetime of the ticket granting ticket, in hours. If not specified, or user specifies 0, then default value 10 will be used.",
																		Immutable:   true,
																	},
																	"truststore": &dcl.Property{
																		Type:        "string",
																		GoName:      "Truststore",
																		Description: "Optional. The Cloud Storage URI of the truststore file used for SSL encryption. If not provided, Dataproc will provide a self-signed certificate.",
																		Immutable:   true,
																	},
																	"truststorePassword": &dcl.Property{
																		Type:        "string",
																		GoName:      "TruststorePassword",
																		Description: "Optional. The Cloud Storage URI of a KMS encrypted file containing the password to the user provided truststore. For the self-signed certificate, this password is generated by Dataproc.",
																		Immutable:   true,
																	},
																},
															},
														},
													},
													"softwareConfig": &dcl.Property{
														Type:        "object",
														GoName:      "SoftwareConfig",
														GoType:      "WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig",
														Description: "Optional. The config settings for software inside the cluster.",
														Immutable:   true,
														Properties: map[string]*dcl.Property{
															"imageVersion": &dcl.Property{
																Type:        "string",
																GoName:      "ImageVersion",
																Description: "Optional. The version of software inside the cluster. It must be one of the supported [Dataproc Versions](https://cloud.google.com/dataproc/docs/concepts/versioning/dataproc-versions#supported_dataproc_versions), such as \"1.2\" (including a subminor version, such as \"1.2.29\"), or the [\"preview\" version](https://cloud.google.com/dataproc/docs/concepts/versioning/dataproc-versions#other_versions). If unspecified, it defaults to the latest Debian version.",
																Immutable:   true,
															},
															"optionalComponents": &dcl.Property{
																Type:        "array",
																GoName:      "OptionalComponents",
																Description: "Optional. The set of components to activate on the cluster.",
																Immutable:   true,
																SendEmpty:   true,
																ListType:    "list",
																Items: &dcl.Property{
																	Type:   "string",
																	GoType: "WorkflowTemplatePlacementManagedClusterConfigSoftwareConfigOptionalComponentsEnum",
																	Enum: []string{
																		"COMPONENT_UNSPECIFIED",
																		"ANACONDA",
																		"DOCKER",
																		"DRUID",
																		"FLINK",
																		"HBASE",
																		"HIVE_WEBHCAT",
																		"JUPYTER",
																		"KERBEROS",
																		"PRESTO",
																		"RANGER",
																		"SOLR",
																		"ZEPPELIN",
																		"ZOOKEEPER",
																	},
																},
															},
															"properties": &dcl.Property{
																Type: "object",
																AdditionalProperties: &dcl.Property{
																	Type: "string",
																},
																GoName:      "Properties",
																Description: "Optional. The properties to set on daemon config files. Property keys are specified in `prefix:property` format, for example `core:hadoop.tmp.dir`. The following are supported prefixes and their mappings: * capacity-scheduler: `capacity-scheduler.xml` * core: `core-site.xml` * distcp: `distcp-default.xml` * hdfs: `hdfs-site.xml` * hive: `hive-site.xml` * mapred: `mapred-site.xml` * pig: `pig.properties` * spark: `spark-defaults.conf` * yarn: `yarn-site.xml` For more information, see [Cluster properties](https://cloud.google.com/dataproc/docs/concepts/cluster-properties).",
																Immutable:   true,
															},
														},
													},
													"stagingBucket": &dcl.Property{
														Type:        "string",
														GoName:      "StagingBucket",
														Description: "Optional. A Cloud Storage bucket used to stage job dependencies, config files, and job driver console output. If you do not specify a staging bucket, Cloud Dataproc will determine a Cloud Storage location (US, ASIA, or EU) for your cluster's staging bucket according to the Compute Engine zone where your cluster is deployed, and then create and manage this project-level, per-location bucket (see [Dataproc staging bucket](https://cloud.google.com/dataproc/docs/concepts/configuring-clusters/staging-bucket)). **This field requires a Cloud Storage bucket name, not a URI to a Cloud Storage bucket.**",
														Immutable:   true,
														ResourceReferences: []*dcl.PropertyResourceReference{
															&dcl.PropertyResourceReference{
																Resource: "Storage/Bucket",
																Field:    "name",
															},
														},
													},
													"tempBucket": &dcl.Property{
														Type:        "string",
														GoName:      "TempBucket",
														Description: "Optional. A Cloud Storage bucket used to store ephemeral cluster and jobs data, such as Spark and MapReduce history files. If you do not specify a temp bucket, Dataproc will determine a Cloud Storage location (US, ASIA, or EU) for your cluster's temp bucket according to the Compute Engine zone where your cluster is deployed, and then create and manage this project-level, per-location bucket. The default bucket has a TTL of 90 days, but you can use any TTL (or none) if you specify a bucket. **This field requires a Cloud Storage bucket name, not a URI to a Cloud Storage bucket.**",
														Immutable:   true,
														ResourceReferences: []*dcl.PropertyResourceReference{
															&dcl.PropertyResourceReference{
																Resource: "Storage/Bucket",
																Field:    "name",
															},
														},
													},
													"workerConfig": &dcl.Property{
														Type:          "object",
														GoName:        "WorkerConfig",
														GoType:        "WorkflowTemplatePlacementManagedClusterConfigWorkerConfig",
														Description:   "Optional. The Compute Engine config settings for worker instances in a cluster.",
														Immutable:     true,
														ServerDefault: true,
														Properties: map[string]*dcl.Property{
															"accelerators": &dcl.Property{
																Type:          "array",
																GoName:        "Accelerators",
																Description:   "Optional. The Compute Engine accelerator configuration for these instances.",
																Immutable:     true,
																ServerDefault: true,
																SendEmpty:     true,
																ListType:      "list",
																Items: &dcl.Property{
																	Type:   "object",
																	GoType: "WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators",
																	Properties: map[string]*dcl.Property{
																		"acceleratorCount": &dcl.Property{
																			Type:        "integer",
																			Format:      "int64",
																			GoName:      "AcceleratorCount",
																			Description: "The number of the accelerator cards of this type exposed to this instance.",
																			Immutable:   true,
																		},
																		"acceleratorType": &dcl.Property{
																			Type:        "string",
																			GoName:      "AcceleratorType",
																			Description: "Full URL, partial URI, or short name of the accelerator type resource to expose to this instance. See [Compute Engine AcceleratorTypes](https://cloud.google.com/compute/docs/reference/beta/acceleratorTypes). Examples: * `https://www.googleapis.com/compute/beta/projects/[project_id]/zones/us-east1-a/acceleratorTypes/nvidia-tesla-k80` * `projects/[project_id]/zones/us-east1-a/acceleratorTypes/nvidia-tesla-k80` * `nvidia-tesla-k80` **Auto Zone Exception**: If you are using the Dataproc [Auto Zone Placement](https://cloud.google.com/dataproc/docs/concepts/configuring-clusters/auto-zone#using_auto_zone_placement) feature, you must use the short name of the accelerator type resource, for example, `nvidia-tesla-k80`.",
																			Immutable:   true,
																		},
																	},
																},
															},
															"diskConfig": &dcl.Property{
																Type:          "object",
																GoName:        "DiskConfig",
																GoType:        "WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig",
																Description:   "Optional. Disk option config settings.",
																Immutable:     true,
																ServerDefault: true,
																Properties: map[string]*dcl.Property{
																	"bootDiskSizeGb": &dcl.Property{
																		Type:        "integer",
																		Format:      "int64",
																		GoName:      "BootDiskSizeGb",
																		Description: "Optional. Size in GB of the boot disk (default is 500GB).",
																		Immutable:   true,
																	},
																	"bootDiskType": &dcl.Property{
																		Type:        "string",
																		GoName:      "BootDiskType",
																		Description: "Optional. Type of the boot disk (default is \"pd-standard\"). Valid values: \"pd-balanced\" (Persistent Disk Balanced Solid State Drive), \"pd-ssd\" (Persistent Disk Solid State Drive), or \"pd-standard\" (Persistent Disk Hard Disk Drive). See [Disk types](https://cloud.google.com/compute/docs/disks#disk-types).",
																		Immutable:   true,
																	},
																	"numLocalSsds": &dcl.Property{
																		Type:          "integer",
																		Format:        "int64",
																		GoName:        "NumLocalSsds",
																		Description:   "Optional. Number of attached SSDs, from 0 to 4 (default is 0). If SSDs are not attached, the boot disk is used to store runtime logs and [HDFS](https://hadoop.apache.org/docs/r1.2.1/hdfs_user_guide.html) data. If one or more SSDs are attached, this runtime bulk data is spread across them, and the boot disk contains only basic config and installed binaries.",
																		Immutable:     true,
																		ServerDefault: true,
																	},
																},
															},
															"image": &dcl.Property{
																Type:        "string",
																GoName:      "Image",
																Description: "Optional. The Compute Engine image resource used for cluster instances. The URI can represent an image or image family. Image examples: * `https://www.googleapis.com/compute/beta/projects/[project_id]/global/images/[image-id]` * `projects/[project_id]/global/images/[image-id]` * `image-id` Image family examples. Dataproc will use the most recent image from the family: * `https://www.googleapis.com/compute/beta/projects/[project_id]/global/images/family/[custom-image-family-name]` * `projects/[project_id]/global/images/family/[custom-image-family-name]` If the URI is unspecified, it will be inferred from `SoftwareConfig.image_version` or the system default.",
																Immutable:   true,
																ResourceReferences: []*dcl.PropertyResourceReference{
																	&dcl.PropertyResourceReference{
																		Resource: "Compute/Image",
																		Field:    "selfLink",
																	},
																},
															},
															"instanceNames": &dcl.Property{
																Type:          "array",
																GoName:        "InstanceNames",
																ReadOnly:      true,
																Description:   "Output only. The list of instance names. Dataproc derives the names from `cluster_name`, `num_instances`, and the instance group.",
																Immutable:     true,
																ServerDefault: true,
																ListType:      "list",
																Items: &dcl.Property{
																	Type:   "string",
																	GoType: "string",
																	ResourceReferences: []*dcl.PropertyResourceReference{
																		&dcl.PropertyResourceReference{
																			Resource: "Compute/Instance",
																			Field:    "selfLink",
																		},
																	},
																},
															},
															"isPreemptible": &dcl.Property{
																Type:        "boolean",
																GoName:      "IsPreemptible",
																ReadOnly:    true,
																Description: "Output only. Specifies that this instance group contains preemptible instances.",
																Immutable:   true,
															},
															"machineType": &dcl.Property{
																Type:        "string",
																GoName:      "MachineType",
																Description: "Optional. The Compute Engine machine type used for cluster instances. A full URL, partial URI, or short name are valid. Examples: * `https://www.googleapis.com/compute/v1/projects/[project_id]/zones/us-east1-a/machineTypes/n1-standard-2` * `projects/[project_id]/zones/us-east1-a/machineTypes/n1-standard-2` * `n1-standard-2` **Auto Zone Exception**: If you are using the Dataproc [Auto Zone Placement](https://cloud.google.com/dataproc/docs/concepts/configuring-clusters/auto-zone#using_auto_zone_placement) feature, you must use the short name of the machine type resource, for example, `n1-standard-2`.",
																Immutable:   true,
															},
															"managedGroupConfig": &dcl.Property{
																Type:          "object",
																GoName:        "ManagedGroupConfig",
																GoType:        "WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig",
																ReadOnly:      true,
																Description:   "Output only. The config for Compute Engine Instance Group Manager that manages this group. This is only used for preemptible instance groups.",
																Immutable:     true,
																ServerDefault: true,
																Properties: map[string]*dcl.Property{
																	"instanceGroupManagerName": &dcl.Property{
																		Type:        "string",
																		GoName:      "InstanceGroupManagerName",
																		ReadOnly:    true,
																		Description: "Output only. The name of the Instance Group Manager for this group.",
																		Immutable:   true,
																	},
																	"instanceTemplateName": &dcl.Property{
																		Type:        "string",
																		GoName:      "InstanceTemplateName",
																		ReadOnly:    true,
																		Description: "Output only. The name of the Instance Template used for the Managed Instance Group.",
																		Immutable:   true,
																	},
																},
															},
															"minCpuPlatform": &dcl.Property{
																Type:          "string",
																GoName:        "MinCpuPlatform",
																Description:   "Optional. Specifies the minimum cpu platform for the Instance Group. See [Dataproc -> Minimum CPU Platform](https://cloud.google.com/dataproc/docs/concepts/compute/dataproc-min-cpu).",
																Immutable:     true,
																ServerDefault: true,
															},
															"numInstances": &dcl.Property{
																Type:        "integer",
																Format:      "int64",
																GoName:      "NumInstances",
																Description: "Optional. The number of VM instances in the instance group. For [HA cluster](/dataproc/docs/concepts/configuring-clusters/high-availability) [master_config](#FIELDS.master_config) groups, **must be set to 3**. For standard cluster [master_config](#FIELDS.master_config) groups, **must be set to 1**.",
																Immutable:   true,
															},
															"preemptibility": &dcl.Property{
																Type:        "string",
																GoName:      "Preemptibility",
																GoType:      "WorkflowTemplatePlacementManagedClusterConfigWorkerConfigPreemptibilityEnum",
																Description: "Optional. Specifies the preemptibility of the instance group. The default value for master and worker groups is `NON_PREEMPTIBLE`. This default cannot be changed. The default value for secondary instances is `PREEMPTIBLE`. Possible values: PREEMPTIBILITY_UNSPECIFIED, NON_PREEMPTIBLE, PREEMPTIBLE",
																Immutable:   true,
																Enum: []string{
																	"PREEMPTIBILITY_UNSPECIFIED",
																	"NON_PREEMPTIBLE",
																	"PREEMPTIBLE",
																},
															},
														},
													},
												},
											},
											"labels": &dcl.Property{
												Type: "object",
												AdditionalProperties: &dcl.Property{
													Type: "string",
												},
												GoName:      "Labels",
												Description: "Optional. The labels to associate with this cluster. Label keys must be between 1 and 63 characters long, and must conform to the following PCRE regular expression: p{Ll}p{Lo}{0,62} Label values must be between 1 and 63 characters long, and must conform to the following PCRE regular expression: [p{Ll}p{Lo}p{N}_-]{0,63} No more than 32 labels can be associated with a given cluster.",
												Immutable:   true,
											},
										},
									},
								},
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
							},
							"updateTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "UpdateTime",
								ReadOnly:    true,
								Description: "Output only. The time template was last updated.",
								Immutable:   true,
							},
							"version": &dcl.Property{
								Type:        "integer",
								Format:      "int64",
								GoName:      "Version",
								ReadOnly:    true,
								Description: "Output only. The current version of this workflow template.",
								Immutable:   true,
							},
						},
					},
				},
			},
		},
	}
}
