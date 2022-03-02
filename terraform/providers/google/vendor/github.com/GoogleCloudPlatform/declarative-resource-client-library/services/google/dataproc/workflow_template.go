// Copyright 2021 Google LLC. All Rights Reserved.
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
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/api/googleapi"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

type WorkflowTemplate struct {
	Name       *string                      `json:"name"`
	Version    *int64                       `json:"version"`
	CreateTime *string                      `json:"createTime"`
	UpdateTime *string                      `json:"updateTime"`
	Labels     map[string]string            `json:"labels"`
	Placement  *WorkflowTemplatePlacement   `json:"placement"`
	Jobs       []WorkflowTemplateJobs       `json:"jobs"`
	Parameters []WorkflowTemplateParameters `json:"parameters"`
	DagTimeout *string                      `json:"dagTimeout"`
	Project    *string                      `json:"project"`
	Location   *string                      `json:"location"`
}

func (r *WorkflowTemplate) String() string {
	return dcl.SprintResource(r)
}

type WorkflowTemplatePlacement struct {
	empty           bool                                      `json:"-"`
	ManagedCluster  *WorkflowTemplatePlacementManagedCluster  `json:"managedCluster"`
	ClusterSelector *WorkflowTemplatePlacementClusterSelector `json:"clusterSelector"`
}

type jsonWorkflowTemplatePlacement WorkflowTemplatePlacement

func (r *WorkflowTemplatePlacement) UnmarshalJSON(data []byte) error {
	var res jsonWorkflowTemplatePlacement
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyWorkflowTemplatePlacement
	} else {

		r.ManagedCluster = res.ManagedCluster

		r.ClusterSelector = res.ClusterSelector

	}
	return nil
}

// This object is used to assert a desired state where this WorkflowTemplatePlacement is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyWorkflowTemplatePlacement *WorkflowTemplatePlacement = &WorkflowTemplatePlacement{empty: true}

func (r *WorkflowTemplatePlacement) Empty() bool {
	return r.empty
}

func (r *WorkflowTemplatePlacement) String() string {
	return dcl.SprintResource(r)
}

func (r *WorkflowTemplatePlacement) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type WorkflowTemplatePlacementManagedCluster struct {
	empty       bool                  `json:"-"`
	ClusterName *string               `json:"clusterName"`
	Config      *ClusterClusterConfig `json:"config"`
	Labels      map[string]string     `json:"labels"`
}

type jsonWorkflowTemplatePlacementManagedCluster WorkflowTemplatePlacementManagedCluster

func (r *WorkflowTemplatePlacementManagedCluster) UnmarshalJSON(data []byte) error {
	var res jsonWorkflowTemplatePlacementManagedCluster
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyWorkflowTemplatePlacementManagedCluster
	} else {

		r.ClusterName = res.ClusterName

		r.Config = res.Config

		r.Labels = res.Labels

	}
	return nil
}

// This object is used to assert a desired state where this WorkflowTemplatePlacementManagedCluster is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyWorkflowTemplatePlacementManagedCluster *WorkflowTemplatePlacementManagedCluster = &WorkflowTemplatePlacementManagedCluster{empty: true}

func (r *WorkflowTemplatePlacementManagedCluster) Empty() bool {
	return r.empty
}

func (r *WorkflowTemplatePlacementManagedCluster) String() string {
	return dcl.SprintResource(r)
}

func (r *WorkflowTemplatePlacementManagedCluster) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type WorkflowTemplatePlacementClusterSelector struct {
	empty         bool              `json:"-"`
	Zone          *string           `json:"zone"`
	ClusterLabels map[string]string `json:"clusterLabels"`
}

type jsonWorkflowTemplatePlacementClusterSelector WorkflowTemplatePlacementClusterSelector

func (r *WorkflowTemplatePlacementClusterSelector) UnmarshalJSON(data []byte) error {
	var res jsonWorkflowTemplatePlacementClusterSelector
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyWorkflowTemplatePlacementClusterSelector
	} else {

		r.Zone = res.Zone

		r.ClusterLabels = res.ClusterLabels

	}
	return nil
}

// This object is used to assert a desired state where this WorkflowTemplatePlacementClusterSelector is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyWorkflowTemplatePlacementClusterSelector *WorkflowTemplatePlacementClusterSelector = &WorkflowTemplatePlacementClusterSelector{empty: true}

func (r *WorkflowTemplatePlacementClusterSelector) Empty() bool {
	return r.empty
}

func (r *WorkflowTemplatePlacementClusterSelector) String() string {
	return dcl.SprintResource(r)
}

func (r *WorkflowTemplatePlacementClusterSelector) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type WorkflowTemplateJobs struct {
	empty               bool                             `json:"-"`
	StepId              *string                          `json:"stepId"`
	HadoopJob           *WorkflowTemplateJobsHadoopJob   `json:"hadoopJob"`
	SparkJob            *WorkflowTemplateJobsSparkJob    `json:"sparkJob"`
	PysparkJob          *WorkflowTemplateJobsPysparkJob  `json:"pysparkJob"`
	HiveJob             *WorkflowTemplateJobsHiveJob     `json:"hiveJob"`
	PigJob              *WorkflowTemplateJobsPigJob      `json:"pigJob"`
	SparkRJob           *WorkflowTemplateJobsSparkRJob   `json:"sparkRJob"`
	SparkSqlJob         *WorkflowTemplateJobsSparkSqlJob `json:"sparkSqlJob"`
	PrestoJob           *WorkflowTemplateJobsPrestoJob   `json:"prestoJob"`
	Labels              map[string]string                `json:"labels"`
	Scheduling          *WorkflowTemplateJobsScheduling  `json:"scheduling"`
	PrerequisiteStepIds []string                         `json:"prerequisiteStepIds"`
}

type jsonWorkflowTemplateJobs WorkflowTemplateJobs

func (r *WorkflowTemplateJobs) UnmarshalJSON(data []byte) error {
	var res jsonWorkflowTemplateJobs
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyWorkflowTemplateJobs
	} else {

		r.StepId = res.StepId

		r.HadoopJob = res.HadoopJob

		r.SparkJob = res.SparkJob

		r.PysparkJob = res.PysparkJob

		r.HiveJob = res.HiveJob

		r.PigJob = res.PigJob

		r.SparkRJob = res.SparkRJob

		r.SparkSqlJob = res.SparkSqlJob

		r.PrestoJob = res.PrestoJob

		r.Labels = res.Labels

		r.Scheduling = res.Scheduling

		r.PrerequisiteStepIds = res.PrerequisiteStepIds

	}
	return nil
}

// This object is used to assert a desired state where this WorkflowTemplateJobs is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyWorkflowTemplateJobs *WorkflowTemplateJobs = &WorkflowTemplateJobs{empty: true}

func (r *WorkflowTemplateJobs) Empty() bool {
	return r.empty
}

func (r *WorkflowTemplateJobs) String() string {
	return dcl.SprintResource(r)
}

func (r *WorkflowTemplateJobs) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type WorkflowTemplateJobsHadoopJob struct {
	empty          bool                                        `json:"-"`
	MainJarFileUri *string                                     `json:"mainJarFileUri"`
	MainClass      *string                                     `json:"mainClass"`
	Args           []string                                    `json:"args"`
	JarFileUris    []string                                    `json:"jarFileUris"`
	FileUris       []string                                    `json:"fileUris"`
	ArchiveUris    []string                                    `json:"archiveUris"`
	Properties     map[string]string                           `json:"properties"`
	LoggingConfig  *WorkflowTemplateJobsHadoopJobLoggingConfig `json:"loggingConfig"`
}

type jsonWorkflowTemplateJobsHadoopJob WorkflowTemplateJobsHadoopJob

func (r *WorkflowTemplateJobsHadoopJob) UnmarshalJSON(data []byte) error {
	var res jsonWorkflowTemplateJobsHadoopJob
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyWorkflowTemplateJobsHadoopJob
	} else {

		r.MainJarFileUri = res.MainJarFileUri

		r.MainClass = res.MainClass

		r.Args = res.Args

		r.JarFileUris = res.JarFileUris

		r.FileUris = res.FileUris

		r.ArchiveUris = res.ArchiveUris

		r.Properties = res.Properties

		r.LoggingConfig = res.LoggingConfig

	}
	return nil
}

// This object is used to assert a desired state where this WorkflowTemplateJobsHadoopJob is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyWorkflowTemplateJobsHadoopJob *WorkflowTemplateJobsHadoopJob = &WorkflowTemplateJobsHadoopJob{empty: true}

func (r *WorkflowTemplateJobsHadoopJob) Empty() bool {
	return r.empty
}

func (r *WorkflowTemplateJobsHadoopJob) String() string {
	return dcl.SprintResource(r)
}

func (r *WorkflowTemplateJobsHadoopJob) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type WorkflowTemplateJobsHadoopJobLoggingConfig struct {
	empty           bool              `json:"-"`
	DriverLogLevels map[string]string `json:"driverLogLevels"`
}

type jsonWorkflowTemplateJobsHadoopJobLoggingConfig WorkflowTemplateJobsHadoopJobLoggingConfig

func (r *WorkflowTemplateJobsHadoopJobLoggingConfig) UnmarshalJSON(data []byte) error {
	var res jsonWorkflowTemplateJobsHadoopJobLoggingConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyWorkflowTemplateJobsHadoopJobLoggingConfig
	} else {

		r.DriverLogLevels = res.DriverLogLevels

	}
	return nil
}

// This object is used to assert a desired state where this WorkflowTemplateJobsHadoopJobLoggingConfig is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyWorkflowTemplateJobsHadoopJobLoggingConfig *WorkflowTemplateJobsHadoopJobLoggingConfig = &WorkflowTemplateJobsHadoopJobLoggingConfig{empty: true}

func (r *WorkflowTemplateJobsHadoopJobLoggingConfig) Empty() bool {
	return r.empty
}

func (r *WorkflowTemplateJobsHadoopJobLoggingConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *WorkflowTemplateJobsHadoopJobLoggingConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type WorkflowTemplateJobsSparkJob struct {
	empty          bool                                       `json:"-"`
	MainJarFileUri *string                                    `json:"mainJarFileUri"`
	MainClass      *string                                    `json:"mainClass"`
	Args           []string                                   `json:"args"`
	JarFileUris    []string                                   `json:"jarFileUris"`
	FileUris       []string                                   `json:"fileUris"`
	ArchiveUris    []string                                   `json:"archiveUris"`
	Properties     map[string]string                          `json:"properties"`
	LoggingConfig  *WorkflowTemplateJobsSparkJobLoggingConfig `json:"loggingConfig"`
}

type jsonWorkflowTemplateJobsSparkJob WorkflowTemplateJobsSparkJob

func (r *WorkflowTemplateJobsSparkJob) UnmarshalJSON(data []byte) error {
	var res jsonWorkflowTemplateJobsSparkJob
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyWorkflowTemplateJobsSparkJob
	} else {

		r.MainJarFileUri = res.MainJarFileUri

		r.MainClass = res.MainClass

		r.Args = res.Args

		r.JarFileUris = res.JarFileUris

		r.FileUris = res.FileUris

		r.ArchiveUris = res.ArchiveUris

		r.Properties = res.Properties

		r.LoggingConfig = res.LoggingConfig

	}
	return nil
}

// This object is used to assert a desired state where this WorkflowTemplateJobsSparkJob is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyWorkflowTemplateJobsSparkJob *WorkflowTemplateJobsSparkJob = &WorkflowTemplateJobsSparkJob{empty: true}

func (r *WorkflowTemplateJobsSparkJob) Empty() bool {
	return r.empty
}

func (r *WorkflowTemplateJobsSparkJob) String() string {
	return dcl.SprintResource(r)
}

func (r *WorkflowTemplateJobsSparkJob) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type WorkflowTemplateJobsSparkJobLoggingConfig struct {
	empty           bool              `json:"-"`
	DriverLogLevels map[string]string `json:"driverLogLevels"`
}

type jsonWorkflowTemplateJobsSparkJobLoggingConfig WorkflowTemplateJobsSparkJobLoggingConfig

func (r *WorkflowTemplateJobsSparkJobLoggingConfig) UnmarshalJSON(data []byte) error {
	var res jsonWorkflowTemplateJobsSparkJobLoggingConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyWorkflowTemplateJobsSparkJobLoggingConfig
	} else {

		r.DriverLogLevels = res.DriverLogLevels

	}
	return nil
}

// This object is used to assert a desired state where this WorkflowTemplateJobsSparkJobLoggingConfig is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyWorkflowTemplateJobsSparkJobLoggingConfig *WorkflowTemplateJobsSparkJobLoggingConfig = &WorkflowTemplateJobsSparkJobLoggingConfig{empty: true}

func (r *WorkflowTemplateJobsSparkJobLoggingConfig) Empty() bool {
	return r.empty
}

func (r *WorkflowTemplateJobsSparkJobLoggingConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *WorkflowTemplateJobsSparkJobLoggingConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type WorkflowTemplateJobsPysparkJob struct {
	empty             bool                                         `json:"-"`
	MainPythonFileUri *string                                      `json:"mainPythonFileUri"`
	Args              []string                                     `json:"args"`
	PythonFileUris    []string                                     `json:"pythonFileUris"`
	JarFileUris       []string                                     `json:"jarFileUris"`
	FileUris          []string                                     `json:"fileUris"`
	ArchiveUris       []string                                     `json:"archiveUris"`
	Properties        map[string]string                            `json:"properties"`
	LoggingConfig     *WorkflowTemplateJobsPysparkJobLoggingConfig `json:"loggingConfig"`
}

type jsonWorkflowTemplateJobsPysparkJob WorkflowTemplateJobsPysparkJob

func (r *WorkflowTemplateJobsPysparkJob) UnmarshalJSON(data []byte) error {
	var res jsonWorkflowTemplateJobsPysparkJob
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyWorkflowTemplateJobsPysparkJob
	} else {

		r.MainPythonFileUri = res.MainPythonFileUri

		r.Args = res.Args

		r.PythonFileUris = res.PythonFileUris

		r.JarFileUris = res.JarFileUris

		r.FileUris = res.FileUris

		r.ArchiveUris = res.ArchiveUris

		r.Properties = res.Properties

		r.LoggingConfig = res.LoggingConfig

	}
	return nil
}

// This object is used to assert a desired state where this WorkflowTemplateJobsPysparkJob is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyWorkflowTemplateJobsPysparkJob *WorkflowTemplateJobsPysparkJob = &WorkflowTemplateJobsPysparkJob{empty: true}

func (r *WorkflowTemplateJobsPysparkJob) Empty() bool {
	return r.empty
}

func (r *WorkflowTemplateJobsPysparkJob) String() string {
	return dcl.SprintResource(r)
}

func (r *WorkflowTemplateJobsPysparkJob) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type WorkflowTemplateJobsPysparkJobLoggingConfig struct {
	empty           bool              `json:"-"`
	DriverLogLevels map[string]string `json:"driverLogLevels"`
}

type jsonWorkflowTemplateJobsPysparkJobLoggingConfig WorkflowTemplateJobsPysparkJobLoggingConfig

func (r *WorkflowTemplateJobsPysparkJobLoggingConfig) UnmarshalJSON(data []byte) error {
	var res jsonWorkflowTemplateJobsPysparkJobLoggingConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyWorkflowTemplateJobsPysparkJobLoggingConfig
	} else {

		r.DriverLogLevels = res.DriverLogLevels

	}
	return nil
}

// This object is used to assert a desired state where this WorkflowTemplateJobsPysparkJobLoggingConfig is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyWorkflowTemplateJobsPysparkJobLoggingConfig *WorkflowTemplateJobsPysparkJobLoggingConfig = &WorkflowTemplateJobsPysparkJobLoggingConfig{empty: true}

func (r *WorkflowTemplateJobsPysparkJobLoggingConfig) Empty() bool {
	return r.empty
}

func (r *WorkflowTemplateJobsPysparkJobLoggingConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *WorkflowTemplateJobsPysparkJobLoggingConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type WorkflowTemplateJobsHiveJob struct {
	empty             bool                                  `json:"-"`
	QueryFileUri      *string                               `json:"queryFileUri"`
	QueryList         *WorkflowTemplateJobsHiveJobQueryList `json:"queryList"`
	ContinueOnFailure *bool                                 `json:"continueOnFailure"`
	ScriptVariables   map[string]string                     `json:"scriptVariables"`
	Properties        map[string]string                     `json:"properties"`
	JarFileUris       []string                              `json:"jarFileUris"`
}

type jsonWorkflowTemplateJobsHiveJob WorkflowTemplateJobsHiveJob

func (r *WorkflowTemplateJobsHiveJob) UnmarshalJSON(data []byte) error {
	var res jsonWorkflowTemplateJobsHiveJob
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyWorkflowTemplateJobsHiveJob
	} else {

		r.QueryFileUri = res.QueryFileUri

		r.QueryList = res.QueryList

		r.ContinueOnFailure = res.ContinueOnFailure

		r.ScriptVariables = res.ScriptVariables

		r.Properties = res.Properties

		r.JarFileUris = res.JarFileUris

	}
	return nil
}

// This object is used to assert a desired state where this WorkflowTemplateJobsHiveJob is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyWorkflowTemplateJobsHiveJob *WorkflowTemplateJobsHiveJob = &WorkflowTemplateJobsHiveJob{empty: true}

func (r *WorkflowTemplateJobsHiveJob) Empty() bool {
	return r.empty
}

func (r *WorkflowTemplateJobsHiveJob) String() string {
	return dcl.SprintResource(r)
}

func (r *WorkflowTemplateJobsHiveJob) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type WorkflowTemplateJobsHiveJobQueryList struct {
	empty   bool     `json:"-"`
	Queries []string `json:"queries"`
}

type jsonWorkflowTemplateJobsHiveJobQueryList WorkflowTemplateJobsHiveJobQueryList

func (r *WorkflowTemplateJobsHiveJobQueryList) UnmarshalJSON(data []byte) error {
	var res jsonWorkflowTemplateJobsHiveJobQueryList
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyWorkflowTemplateJobsHiveJobQueryList
	} else {

		r.Queries = res.Queries

	}
	return nil
}

// This object is used to assert a desired state where this WorkflowTemplateJobsHiveJobQueryList is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyWorkflowTemplateJobsHiveJobQueryList *WorkflowTemplateJobsHiveJobQueryList = &WorkflowTemplateJobsHiveJobQueryList{empty: true}

func (r *WorkflowTemplateJobsHiveJobQueryList) Empty() bool {
	return r.empty
}

func (r *WorkflowTemplateJobsHiveJobQueryList) String() string {
	return dcl.SprintResource(r)
}

func (r *WorkflowTemplateJobsHiveJobQueryList) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type WorkflowTemplateJobsPigJob struct {
	empty             bool                                     `json:"-"`
	QueryFileUri      *string                                  `json:"queryFileUri"`
	QueryList         *WorkflowTemplateJobsPigJobQueryList     `json:"queryList"`
	ContinueOnFailure *bool                                    `json:"continueOnFailure"`
	ScriptVariables   map[string]string                        `json:"scriptVariables"`
	Properties        map[string]string                        `json:"properties"`
	JarFileUris       []string                                 `json:"jarFileUris"`
	LoggingConfig     *WorkflowTemplateJobsPigJobLoggingConfig `json:"loggingConfig"`
}

type jsonWorkflowTemplateJobsPigJob WorkflowTemplateJobsPigJob

func (r *WorkflowTemplateJobsPigJob) UnmarshalJSON(data []byte) error {
	var res jsonWorkflowTemplateJobsPigJob
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyWorkflowTemplateJobsPigJob
	} else {

		r.QueryFileUri = res.QueryFileUri

		r.QueryList = res.QueryList

		r.ContinueOnFailure = res.ContinueOnFailure

		r.ScriptVariables = res.ScriptVariables

		r.Properties = res.Properties

		r.JarFileUris = res.JarFileUris

		r.LoggingConfig = res.LoggingConfig

	}
	return nil
}

// This object is used to assert a desired state where this WorkflowTemplateJobsPigJob is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyWorkflowTemplateJobsPigJob *WorkflowTemplateJobsPigJob = &WorkflowTemplateJobsPigJob{empty: true}

func (r *WorkflowTemplateJobsPigJob) Empty() bool {
	return r.empty
}

func (r *WorkflowTemplateJobsPigJob) String() string {
	return dcl.SprintResource(r)
}

func (r *WorkflowTemplateJobsPigJob) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type WorkflowTemplateJobsPigJobQueryList struct {
	empty   bool     `json:"-"`
	Queries []string `json:"queries"`
}

type jsonWorkflowTemplateJobsPigJobQueryList WorkflowTemplateJobsPigJobQueryList

func (r *WorkflowTemplateJobsPigJobQueryList) UnmarshalJSON(data []byte) error {
	var res jsonWorkflowTemplateJobsPigJobQueryList
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyWorkflowTemplateJobsPigJobQueryList
	} else {

		r.Queries = res.Queries

	}
	return nil
}

// This object is used to assert a desired state where this WorkflowTemplateJobsPigJobQueryList is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyWorkflowTemplateJobsPigJobQueryList *WorkflowTemplateJobsPigJobQueryList = &WorkflowTemplateJobsPigJobQueryList{empty: true}

func (r *WorkflowTemplateJobsPigJobQueryList) Empty() bool {
	return r.empty
}

func (r *WorkflowTemplateJobsPigJobQueryList) String() string {
	return dcl.SprintResource(r)
}

func (r *WorkflowTemplateJobsPigJobQueryList) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type WorkflowTemplateJobsPigJobLoggingConfig struct {
	empty           bool              `json:"-"`
	DriverLogLevels map[string]string `json:"driverLogLevels"`
}

type jsonWorkflowTemplateJobsPigJobLoggingConfig WorkflowTemplateJobsPigJobLoggingConfig

func (r *WorkflowTemplateJobsPigJobLoggingConfig) UnmarshalJSON(data []byte) error {
	var res jsonWorkflowTemplateJobsPigJobLoggingConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyWorkflowTemplateJobsPigJobLoggingConfig
	} else {

		r.DriverLogLevels = res.DriverLogLevels

	}
	return nil
}

// This object is used to assert a desired state where this WorkflowTemplateJobsPigJobLoggingConfig is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyWorkflowTemplateJobsPigJobLoggingConfig *WorkflowTemplateJobsPigJobLoggingConfig = &WorkflowTemplateJobsPigJobLoggingConfig{empty: true}

func (r *WorkflowTemplateJobsPigJobLoggingConfig) Empty() bool {
	return r.empty
}

func (r *WorkflowTemplateJobsPigJobLoggingConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *WorkflowTemplateJobsPigJobLoggingConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type WorkflowTemplateJobsSparkRJob struct {
	empty         bool                                        `json:"-"`
	MainRFileUri  *string                                     `json:"mainRFileUri"`
	Args          []string                                    `json:"args"`
	FileUris      []string                                    `json:"fileUris"`
	ArchiveUris   []string                                    `json:"archiveUris"`
	Properties    map[string]string                           `json:"properties"`
	LoggingConfig *WorkflowTemplateJobsSparkRJobLoggingConfig `json:"loggingConfig"`
}

type jsonWorkflowTemplateJobsSparkRJob WorkflowTemplateJobsSparkRJob

func (r *WorkflowTemplateJobsSparkRJob) UnmarshalJSON(data []byte) error {
	var res jsonWorkflowTemplateJobsSparkRJob
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyWorkflowTemplateJobsSparkRJob
	} else {

		r.MainRFileUri = res.MainRFileUri

		r.Args = res.Args

		r.FileUris = res.FileUris

		r.ArchiveUris = res.ArchiveUris

		r.Properties = res.Properties

		r.LoggingConfig = res.LoggingConfig

	}
	return nil
}

// This object is used to assert a desired state where this WorkflowTemplateJobsSparkRJob is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyWorkflowTemplateJobsSparkRJob *WorkflowTemplateJobsSparkRJob = &WorkflowTemplateJobsSparkRJob{empty: true}

func (r *WorkflowTemplateJobsSparkRJob) Empty() bool {
	return r.empty
}

func (r *WorkflowTemplateJobsSparkRJob) String() string {
	return dcl.SprintResource(r)
}

func (r *WorkflowTemplateJobsSparkRJob) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type WorkflowTemplateJobsSparkRJobLoggingConfig struct {
	empty           bool              `json:"-"`
	DriverLogLevels map[string]string `json:"driverLogLevels"`
}

type jsonWorkflowTemplateJobsSparkRJobLoggingConfig WorkflowTemplateJobsSparkRJobLoggingConfig

func (r *WorkflowTemplateJobsSparkRJobLoggingConfig) UnmarshalJSON(data []byte) error {
	var res jsonWorkflowTemplateJobsSparkRJobLoggingConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyWorkflowTemplateJobsSparkRJobLoggingConfig
	} else {

		r.DriverLogLevels = res.DriverLogLevels

	}
	return nil
}

// This object is used to assert a desired state where this WorkflowTemplateJobsSparkRJobLoggingConfig is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyWorkflowTemplateJobsSparkRJobLoggingConfig *WorkflowTemplateJobsSparkRJobLoggingConfig = &WorkflowTemplateJobsSparkRJobLoggingConfig{empty: true}

func (r *WorkflowTemplateJobsSparkRJobLoggingConfig) Empty() bool {
	return r.empty
}

func (r *WorkflowTemplateJobsSparkRJobLoggingConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *WorkflowTemplateJobsSparkRJobLoggingConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type WorkflowTemplateJobsSparkSqlJob struct {
	empty           bool                                          `json:"-"`
	QueryFileUri    *string                                       `json:"queryFileUri"`
	QueryList       *WorkflowTemplateJobsSparkSqlJobQueryList     `json:"queryList"`
	ScriptVariables map[string]string                             `json:"scriptVariables"`
	Properties      map[string]string                             `json:"properties"`
	JarFileUris     []string                                      `json:"jarFileUris"`
	LoggingConfig   *WorkflowTemplateJobsSparkSqlJobLoggingConfig `json:"loggingConfig"`
}

type jsonWorkflowTemplateJobsSparkSqlJob WorkflowTemplateJobsSparkSqlJob

func (r *WorkflowTemplateJobsSparkSqlJob) UnmarshalJSON(data []byte) error {
	var res jsonWorkflowTemplateJobsSparkSqlJob
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyWorkflowTemplateJobsSparkSqlJob
	} else {

		r.QueryFileUri = res.QueryFileUri

		r.QueryList = res.QueryList

		r.ScriptVariables = res.ScriptVariables

		r.Properties = res.Properties

		r.JarFileUris = res.JarFileUris

		r.LoggingConfig = res.LoggingConfig

	}
	return nil
}

// This object is used to assert a desired state where this WorkflowTemplateJobsSparkSqlJob is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyWorkflowTemplateJobsSparkSqlJob *WorkflowTemplateJobsSparkSqlJob = &WorkflowTemplateJobsSparkSqlJob{empty: true}

func (r *WorkflowTemplateJobsSparkSqlJob) Empty() bool {
	return r.empty
}

func (r *WorkflowTemplateJobsSparkSqlJob) String() string {
	return dcl.SprintResource(r)
}

func (r *WorkflowTemplateJobsSparkSqlJob) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type WorkflowTemplateJobsSparkSqlJobQueryList struct {
	empty   bool     `json:"-"`
	Queries []string `json:"queries"`
}

type jsonWorkflowTemplateJobsSparkSqlJobQueryList WorkflowTemplateJobsSparkSqlJobQueryList

func (r *WorkflowTemplateJobsSparkSqlJobQueryList) UnmarshalJSON(data []byte) error {
	var res jsonWorkflowTemplateJobsSparkSqlJobQueryList
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyWorkflowTemplateJobsSparkSqlJobQueryList
	} else {

		r.Queries = res.Queries

	}
	return nil
}

// This object is used to assert a desired state where this WorkflowTemplateJobsSparkSqlJobQueryList is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyWorkflowTemplateJobsSparkSqlJobQueryList *WorkflowTemplateJobsSparkSqlJobQueryList = &WorkflowTemplateJobsSparkSqlJobQueryList{empty: true}

func (r *WorkflowTemplateJobsSparkSqlJobQueryList) Empty() bool {
	return r.empty
}

func (r *WorkflowTemplateJobsSparkSqlJobQueryList) String() string {
	return dcl.SprintResource(r)
}

func (r *WorkflowTemplateJobsSparkSqlJobQueryList) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type WorkflowTemplateJobsSparkSqlJobLoggingConfig struct {
	empty           bool              `json:"-"`
	DriverLogLevels map[string]string `json:"driverLogLevels"`
}

type jsonWorkflowTemplateJobsSparkSqlJobLoggingConfig WorkflowTemplateJobsSparkSqlJobLoggingConfig

func (r *WorkflowTemplateJobsSparkSqlJobLoggingConfig) UnmarshalJSON(data []byte) error {
	var res jsonWorkflowTemplateJobsSparkSqlJobLoggingConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyWorkflowTemplateJobsSparkSqlJobLoggingConfig
	} else {

		r.DriverLogLevels = res.DriverLogLevels

	}
	return nil
}

// This object is used to assert a desired state where this WorkflowTemplateJobsSparkSqlJobLoggingConfig is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyWorkflowTemplateJobsSparkSqlJobLoggingConfig *WorkflowTemplateJobsSparkSqlJobLoggingConfig = &WorkflowTemplateJobsSparkSqlJobLoggingConfig{empty: true}

func (r *WorkflowTemplateJobsSparkSqlJobLoggingConfig) Empty() bool {
	return r.empty
}

func (r *WorkflowTemplateJobsSparkSqlJobLoggingConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *WorkflowTemplateJobsSparkSqlJobLoggingConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type WorkflowTemplateJobsPrestoJob struct {
	empty             bool                                        `json:"-"`
	QueryFileUri      *string                                     `json:"queryFileUri"`
	QueryList         *WorkflowTemplateJobsPrestoJobQueryList     `json:"queryList"`
	ContinueOnFailure *bool                                       `json:"continueOnFailure"`
	OutputFormat      *string                                     `json:"outputFormat"`
	ClientTags        []string                                    `json:"clientTags"`
	Properties        map[string]string                           `json:"properties"`
	LoggingConfig     *WorkflowTemplateJobsPrestoJobLoggingConfig `json:"loggingConfig"`
}

type jsonWorkflowTemplateJobsPrestoJob WorkflowTemplateJobsPrestoJob

func (r *WorkflowTemplateJobsPrestoJob) UnmarshalJSON(data []byte) error {
	var res jsonWorkflowTemplateJobsPrestoJob
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyWorkflowTemplateJobsPrestoJob
	} else {

		r.QueryFileUri = res.QueryFileUri

		r.QueryList = res.QueryList

		r.ContinueOnFailure = res.ContinueOnFailure

		r.OutputFormat = res.OutputFormat

		r.ClientTags = res.ClientTags

		r.Properties = res.Properties

		r.LoggingConfig = res.LoggingConfig

	}
	return nil
}

// This object is used to assert a desired state where this WorkflowTemplateJobsPrestoJob is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyWorkflowTemplateJobsPrestoJob *WorkflowTemplateJobsPrestoJob = &WorkflowTemplateJobsPrestoJob{empty: true}

func (r *WorkflowTemplateJobsPrestoJob) Empty() bool {
	return r.empty
}

func (r *WorkflowTemplateJobsPrestoJob) String() string {
	return dcl.SprintResource(r)
}

func (r *WorkflowTemplateJobsPrestoJob) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type WorkflowTemplateJobsPrestoJobQueryList struct {
	empty   bool     `json:"-"`
	Queries []string `json:"queries"`
}

type jsonWorkflowTemplateJobsPrestoJobQueryList WorkflowTemplateJobsPrestoJobQueryList

func (r *WorkflowTemplateJobsPrestoJobQueryList) UnmarshalJSON(data []byte) error {
	var res jsonWorkflowTemplateJobsPrestoJobQueryList
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyWorkflowTemplateJobsPrestoJobQueryList
	} else {

		r.Queries = res.Queries

	}
	return nil
}

// This object is used to assert a desired state where this WorkflowTemplateJobsPrestoJobQueryList is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyWorkflowTemplateJobsPrestoJobQueryList *WorkflowTemplateJobsPrestoJobQueryList = &WorkflowTemplateJobsPrestoJobQueryList{empty: true}

func (r *WorkflowTemplateJobsPrestoJobQueryList) Empty() bool {
	return r.empty
}

func (r *WorkflowTemplateJobsPrestoJobQueryList) String() string {
	return dcl.SprintResource(r)
}

func (r *WorkflowTemplateJobsPrestoJobQueryList) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type WorkflowTemplateJobsPrestoJobLoggingConfig struct {
	empty           bool              `json:"-"`
	DriverLogLevels map[string]string `json:"driverLogLevels"`
}

type jsonWorkflowTemplateJobsPrestoJobLoggingConfig WorkflowTemplateJobsPrestoJobLoggingConfig

func (r *WorkflowTemplateJobsPrestoJobLoggingConfig) UnmarshalJSON(data []byte) error {
	var res jsonWorkflowTemplateJobsPrestoJobLoggingConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyWorkflowTemplateJobsPrestoJobLoggingConfig
	} else {

		r.DriverLogLevels = res.DriverLogLevels

	}
	return nil
}

// This object is used to assert a desired state where this WorkflowTemplateJobsPrestoJobLoggingConfig is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyWorkflowTemplateJobsPrestoJobLoggingConfig *WorkflowTemplateJobsPrestoJobLoggingConfig = &WorkflowTemplateJobsPrestoJobLoggingConfig{empty: true}

func (r *WorkflowTemplateJobsPrestoJobLoggingConfig) Empty() bool {
	return r.empty
}

func (r *WorkflowTemplateJobsPrestoJobLoggingConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *WorkflowTemplateJobsPrestoJobLoggingConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type WorkflowTemplateJobsScheduling struct {
	empty              bool   `json:"-"`
	MaxFailuresPerHour *int64 `json:"maxFailuresPerHour"`
	MaxFailuresTotal   *int64 `json:"maxFailuresTotal"`
}

type jsonWorkflowTemplateJobsScheduling WorkflowTemplateJobsScheduling

func (r *WorkflowTemplateJobsScheduling) UnmarshalJSON(data []byte) error {
	var res jsonWorkflowTemplateJobsScheduling
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyWorkflowTemplateJobsScheduling
	} else {

		r.MaxFailuresPerHour = res.MaxFailuresPerHour

		r.MaxFailuresTotal = res.MaxFailuresTotal

	}
	return nil
}

// This object is used to assert a desired state where this WorkflowTemplateJobsScheduling is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyWorkflowTemplateJobsScheduling *WorkflowTemplateJobsScheduling = &WorkflowTemplateJobsScheduling{empty: true}

func (r *WorkflowTemplateJobsScheduling) Empty() bool {
	return r.empty
}

func (r *WorkflowTemplateJobsScheduling) String() string {
	return dcl.SprintResource(r)
}

func (r *WorkflowTemplateJobsScheduling) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type WorkflowTemplateParameters struct {
	empty       bool                                  `json:"-"`
	Name        *string                               `json:"name"`
	Fields      []string                              `json:"fields"`
	Description *string                               `json:"description"`
	Validation  *WorkflowTemplateParametersValidation `json:"validation"`
}

type jsonWorkflowTemplateParameters WorkflowTemplateParameters

func (r *WorkflowTemplateParameters) UnmarshalJSON(data []byte) error {
	var res jsonWorkflowTemplateParameters
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyWorkflowTemplateParameters
	} else {

		r.Name = res.Name

		r.Fields = res.Fields

		r.Description = res.Description

		r.Validation = res.Validation

	}
	return nil
}

// This object is used to assert a desired state where this WorkflowTemplateParameters is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyWorkflowTemplateParameters *WorkflowTemplateParameters = &WorkflowTemplateParameters{empty: true}

func (r *WorkflowTemplateParameters) Empty() bool {
	return r.empty
}

func (r *WorkflowTemplateParameters) String() string {
	return dcl.SprintResource(r)
}

func (r *WorkflowTemplateParameters) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type WorkflowTemplateParametersValidation struct {
	empty  bool                                        `json:"-"`
	Regex  *WorkflowTemplateParametersValidationRegex  `json:"regex"`
	Values *WorkflowTemplateParametersValidationValues `json:"values"`
}

type jsonWorkflowTemplateParametersValidation WorkflowTemplateParametersValidation

func (r *WorkflowTemplateParametersValidation) UnmarshalJSON(data []byte) error {
	var res jsonWorkflowTemplateParametersValidation
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyWorkflowTemplateParametersValidation
	} else {

		r.Regex = res.Regex

		r.Values = res.Values

	}
	return nil
}

// This object is used to assert a desired state where this WorkflowTemplateParametersValidation is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyWorkflowTemplateParametersValidation *WorkflowTemplateParametersValidation = &WorkflowTemplateParametersValidation{empty: true}

func (r *WorkflowTemplateParametersValidation) Empty() bool {
	return r.empty
}

func (r *WorkflowTemplateParametersValidation) String() string {
	return dcl.SprintResource(r)
}

func (r *WorkflowTemplateParametersValidation) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type WorkflowTemplateParametersValidationRegex struct {
	empty   bool     `json:"-"`
	Regexes []string `json:"regexes"`
}

type jsonWorkflowTemplateParametersValidationRegex WorkflowTemplateParametersValidationRegex

func (r *WorkflowTemplateParametersValidationRegex) UnmarshalJSON(data []byte) error {
	var res jsonWorkflowTemplateParametersValidationRegex
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyWorkflowTemplateParametersValidationRegex
	} else {

		r.Regexes = res.Regexes

	}
	return nil
}

// This object is used to assert a desired state where this WorkflowTemplateParametersValidationRegex is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyWorkflowTemplateParametersValidationRegex *WorkflowTemplateParametersValidationRegex = &WorkflowTemplateParametersValidationRegex{empty: true}

func (r *WorkflowTemplateParametersValidationRegex) Empty() bool {
	return r.empty
}

func (r *WorkflowTemplateParametersValidationRegex) String() string {
	return dcl.SprintResource(r)
}

func (r *WorkflowTemplateParametersValidationRegex) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type WorkflowTemplateParametersValidationValues struct {
	empty  bool     `json:"-"`
	Values []string `json:"values"`
}

type jsonWorkflowTemplateParametersValidationValues WorkflowTemplateParametersValidationValues

func (r *WorkflowTemplateParametersValidationValues) UnmarshalJSON(data []byte) error {
	var res jsonWorkflowTemplateParametersValidationValues
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyWorkflowTemplateParametersValidationValues
	} else {

		r.Values = res.Values

	}
	return nil
}

// This object is used to assert a desired state where this WorkflowTemplateParametersValidationValues is
// empty.  Go lacks global const objects, but this object should be treated
// as one.  Modifying this object will have undesirable results.
var EmptyWorkflowTemplateParametersValidationValues *WorkflowTemplateParametersValidationValues = &WorkflowTemplateParametersValidationValues{empty: true}

func (r *WorkflowTemplateParametersValidationValues) Empty() bool {
	return r.empty
}

func (r *WorkflowTemplateParametersValidationValues) String() string {
	return dcl.SprintResource(r)
}

func (r *WorkflowTemplateParametersValidationValues) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *WorkflowTemplate) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "dataproc",
		Type:    "WorkflowTemplate",
		Version: "dataproc",
	}
}

func (r *WorkflowTemplate) ID() (string, error) {
	if err := extractWorkflowTemplateFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name":       dcl.ValueOrEmptyString(nr.Name),
		"version":    dcl.ValueOrEmptyString(nr.Version),
		"createTime": dcl.ValueOrEmptyString(nr.CreateTime),
		"updateTime": dcl.ValueOrEmptyString(nr.UpdateTime),
		"labels":     dcl.ValueOrEmptyString(nr.Labels),
		"placement":  dcl.ValueOrEmptyString(nr.Placement),
		"jobs":       dcl.ValueOrEmptyString(nr.Jobs),
		"parameters": dcl.ValueOrEmptyString(nr.Parameters),
		"dagTimeout": dcl.ValueOrEmptyString(nr.DagTimeout),
		"project":    dcl.ValueOrEmptyString(nr.Project),
		"location":   dcl.ValueOrEmptyString(nr.Location),
	}
	return dcl.Nprintf("projects/{{project}}/locations/{{location}}/workflowTemplates/{{name}}", params), nil
}

const WorkflowTemplateMaxPage = -1

type WorkflowTemplateList struct {
	Items []*WorkflowTemplate

	nextToken string

	pageSize int32

	resource *WorkflowTemplate
}

func (l *WorkflowTemplateList) HasNext() bool {
	return l.nextToken != ""
}

func (l *WorkflowTemplateList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listWorkflowTemplate(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListWorkflowTemplate(ctx context.Context, project, location string) (*WorkflowTemplateList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListWorkflowTemplateWithMaxResults(ctx, project, location, WorkflowTemplateMaxPage)

}

func (c *Client) ListWorkflowTemplateWithMaxResults(ctx context.Context, project, location string, pageSize int32) (*WorkflowTemplateList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &WorkflowTemplate{
		Project:  &project,
		Location: &location,
	}
	items, token, err := c.listWorkflowTemplate(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &WorkflowTemplateList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) GetWorkflowTemplate(ctx context.Context, r *WorkflowTemplate) (*WorkflowTemplate, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractWorkflowTemplateFields(r)

	b, err := c.getWorkflowTemplateRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalWorkflowTemplate(b, c)
	if err != nil {
		return nil, err
	}
	result.Project = r.Project
	result.Location = r.Location
	result.Name = r.Name

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizeWorkflowTemplateNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractWorkflowTemplateFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) DeleteWorkflowTemplate(ctx context.Context, r *WorkflowTemplate) error {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("WorkflowTemplate resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting WorkflowTemplate...")
	deleteOp := deleteWorkflowTemplateOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllWorkflowTemplate deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllWorkflowTemplate(ctx context.Context, project, location string, filter func(*WorkflowTemplate) bool) error {
	listObj, err := c.ListWorkflowTemplate(ctx, project, location)
	if err != nil {
		return err
	}

	err = c.deleteAllWorkflowTemplate(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllWorkflowTemplate(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplyWorkflowTemplate(ctx context.Context, rawDesired *WorkflowTemplate, opts ...dcl.ApplyOption) (*WorkflowTemplate, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	var resultNewState *WorkflowTemplate
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyWorkflowTemplateHelper(c, ctx, rawDesired, opts...)
		resultNewState = newState
		if err != nil {
			// If the error is 409, there is conflict in resource update.
			// Here we want to apply changes based on latest state.
			if dcl.IsConflictError(err) {
				return &dcl.RetryDetails{}, dcl.OperationNotDone{Err: err}
			}
			return nil, err
		}
		return nil, nil
	}, c.Config.RetryProvider)
	return resultNewState, err
}

func applyWorkflowTemplateHelper(c *Client, ctx context.Context, rawDesired *WorkflowTemplate, opts ...dcl.ApplyOption) (*WorkflowTemplate, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyWorkflowTemplate...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractWorkflowTemplateFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.workflowTemplateDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToWorkflowTemplateDiffs(c.Config, fieldDiffs, opts)
	if err != nil {
		return nil, err
	}

	// TODO(magic-modules-eng): 2.2 Feasibility check (all updates are feasible so far).

	// 2.3: Lifecycle Directive Check
	var create bool
	lp := dcl.FetchLifecycleParams(opts)
	if initial == nil {
		if dcl.HasLifecycleParam(lp, dcl.BlockCreation) {
			return nil, dcl.ApplyInfeasibleError{Message: fmt.Sprintf("Creation blocked by lifecycle params: %#v.", desired)}
		}
		create = true
	} else if dcl.HasLifecycleParam(lp, dcl.BlockAcquire) {
		return nil, dcl.ApplyInfeasibleError{
			Message: fmt.Sprintf("Resource already exists - apply blocked by lifecycle params: %#v.", initial),
		}
	} else {
		for _, d := range diffs {
			if d.RequiresRecreate {
				return nil, dcl.ApplyInfeasibleError{
					Message: fmt.Sprintf("infeasible update: (%v) would require recreation", d),
				}
			}
			if dcl.HasLifecycleParam(lp, dcl.BlockModification) {
				return nil, dcl.ApplyInfeasibleError{Message: fmt.Sprintf("Modification blocked, diff (%v) unresolvable.", d)}
			}
		}
	}

	// 2.4 Imperative Request Planning
	var ops []workflowTemplateApiOperation
	if create {
		ops = append(ops, &createWorkflowTemplateOperation{})
	} else {
		for _, d := range diffs {
			ops = append(ops, d.UpdateOp)
		}
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created plan: %#v", ops)

	// 2.5 Request Actuation
	for _, op := range ops {
		c.Config.Logger.InfoWithContextf(ctx, "Performing operation %T %+v", op, op)
		if err := op.do(ctx, desired, c); err != nil {
			c.Config.Logger.InfoWithContextf(ctx, "Failed operation %T %+v: %v", op, op, err)
			return nil, err
		}
		c.Config.Logger.InfoWithContextf(ctx, "Finished operation %T %+v", op, op)
	}
	return applyWorkflowTemplateDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyWorkflowTemplateDiff(c *Client, ctx context.Context, desired *WorkflowTemplate, rawDesired *WorkflowTemplate, ops []workflowTemplateApiOperation, opts ...dcl.ApplyOption) (*WorkflowTemplate, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetWorkflowTemplate(ctx, desired.urlNormalized())
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createWorkflowTemplateOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapWorkflowTemplate(r, c)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizeWorkflowTemplateNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizeWorkflowTemplateNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizeWorkflowTemplateDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractWorkflowTemplateFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractWorkflowTemplateFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffWorkflowTemplate(c, newDesired, newState)
	if err != nil {
		return newState, err
	}

	if len(newDiffs) == 0 {
		c.Config.Logger.InfoWithContext(ctx, "No diffs found. Apply was successful.")
	} else {
		c.Config.Logger.InfoWithContextf(ctx, "Found diffs: %v", newDiffs)
		diffMessages := make([]string, len(newDiffs))
		for i, d := range newDiffs {
			diffMessages[i] = fmt.Sprintf("%v", d)
		}
		return newState, dcl.DiffAfterApplyError{Diffs: diffMessages}
	}
	c.Config.Logger.InfoWithContext(ctx, "Done Apply.")
	return newState, nil
}
