# Workflow report

## Introduction
Every installer command is internally based on the [assets framework][../design/assetgeneration.md], which allows to define a directed acyclic graph of assets (a generic work item).
A workflow identifies a set of asset graph paths executed when running a specific command (note that the same asset may behave differently when triggered under different workflows).

A report allows to collect the main relevant workflow events - as well as any eventual error, and detailed results - in a human-readable format, providing an abstraction from the basic asset logging system.
A report also serializes immediately on the disk any update received, thus providing a simple mechanism to support streaming the progression of the current command execution (particulary
useful in case of remote execution).

## Stages
A report is a composed by a number of sequential _stage_. A stage represents a specific workflow phase (given by one or more assets) that could be relevant and/or informative from the user point of view.
A stage is composed by a _stage identifier_ (or briefly, a stage id) and a number of substages:

* The stage id is defined by a short internal identifier and a longer human-readable description, meant to be shown to the user.
* Substages are a sequential of zero or more stages, and they are meant to provide a further level of details for the owner stage, if required.

### Stages ids
A stage id can be defined via the `NewStageID` method, for example:
```
StageFetchBaseISO wr.StageID = wr.NewStageID("fetch-base-iso", "Retrieving the base ISO image")
```

A substage id can be defined by prefixing the internal id with the owner stage id, in the format `<owner stage>.<substage>`, for example:
```
StageFetchBaseISOExtract wr.StageID = wr.NewStageID("fetch-base-iso.extract-image", "Extracting base image from release payload")
StageFetchBaseISOVerify wr.StageID = wr.NewStageID("fetch-base-iso.verify-version", "Verifying base image version")
StageFetchBaseISODownload wr.StageID = wr.NewStageID("fetch-base-iso.download-image", "Downloading base ISO image")
```

### Substages
Substages are completely optional, and they may be useful for detailing a particularly big or lengthy stage. So, it's perfectly fine to define stages without any substage.
The current framework does not support more than two levels of stages, essentially for keeping both a simpler interface and for producing an easy to read output for the
final user.

### (Sub)Stage result
The workflow report framework allows to attach (optionally) an artifact for each stage (or substage), in order to capture the end result of the stage execution. The result field is
a free-text string, and its format depends on the specific stage (even though a JSON format is recommended).

## Populating a report
Every asset can access the current report from the `Generate()` context object (for the installer commands where the reporting have been enabled), using the `GetReport()` method.
Once retrieved, the report `Stage` method can be used to add a new stage to the current report:
```
workflowreport.GetReport(ctx).Stage(workflow.StageFetchBaseISO)
```

A similar approach could be follow to add a new substage to the report:
```
workflowreport.GetReport(ctx).SubStage(workflow.StageFetchBaseISOExtract)
```

_Note: a substage cannot be added before adding the related owner stage._

## Enabling the report for a command
If a report was not previously enabled for a specific command, the previously shown commands will have no effect. To active the reporting,
it is sufficient to use the `workflowreport.Context()` method which will create a dedicated context to be used in the assets generation: 

```
func NewAddNodesCommand(directory string, kubeConfig string) error {
	
	ctx := workflowreport.Context(string(workflow.AgentWorkflowTypeAddNodes), directory)

	fetcher := store.NewAssetsFetcher(directory)
	err = fetcher.FetchAndPersist(ctx, ...) //

	workflowreport.GetReport(ctx).Complete(err)
```

At the end, the `Complete(err)` method must be invoked to close the report (and eventually report any error).

## Appendix

### Sample of report.json file (without any result)
```
{
  "id": "report-addnodes-202410301546",
  "start_time": "2024-10-30T15:46:36.646915757Z",
  "end_time": "2024-10-30T15:47:23.689529009Z",
  "stages": [
    {
      "id": "add-nodes-cluster-inspection",
      "description": "Gathering additional information from the target cluster",
      "start_time": "2024-10-30T15:46:36.646915757Z",
      "end_time": "2024-10-30T15:46:36.814616705Z"
    },
    {
      "id": "create-manifest",
      "description": "Creating internal configuration manifests",
      "start_time": "2024-10-30T15:46:36.8146174Z",
      "end_time": "2024-10-30T15:46:37.890492356Z"
    },
    {
      "id": "ignition",
      "description": "Rendering ISO ignition",
      "start_time": "2024-10-30T15:46:37.890493244Z",
      "end_time": "2024-10-30T15:46:38.227771699Z"
    },
    {
      "id": "fetch-base-iso",
      "description": "Retrieving the base ISO image",
      "start_time": "2024-10-30T15:46:38.227772149Z",
      "end_time": "2024-10-30T15:47:00.63994573Z",
      "sub_stages": [
        {
          "id": "fetch-base-iso.extract-image",
          "description": "Extracting base image from release payload",
          "start_time": "2024-10-30T15:46:38.228032465Z",
          "end_time": "2024-10-30T15:46:51.428272041Z"
        },
        {
          "id": "fetch-base-iso.verify-version",
          "description": "Verifying base image version",
          "start_time": "2024-10-30T15:46:51.428273074Z",
          "end_time": "2024-10-30T15:47:00.639945356Z"
        }
      ]
    },
    {
      "id": "create-agent-artifacts",
      "description": "Creating agent artifacts for the final image",
      "start_time": "2024-10-30T15:47:00.639946343Z",
      "end_time": "2024-10-30T15:47:20.70873527Z",
      "sub_stages": [
        {
          "id": "create-agent-artifacts.agent-tui",
          "description": "Extracting required artifacts from release payload",
          "start_time": "2024-10-30T15:47:00.641937138Z",
          "end_time": "2024-10-30T15:47:10.419651275Z"
        },
        {
          "id": "create-agent-artifacts.prepare",
          "description": "Preparing artifacts",
          "start_time": "2024-10-30T15:47:10.419652528Z",
          "end_time": "2024-10-30T15:47:20.708735111Z"
        }
      ]
    },
    {
      "id": "generate-iso",
      "description": "Assembling ISO image",
      "start_time": "2024-10-30T15:47:20.708735596Z",
      "end_time": "2024-10-30T15:47:23.689529009Z"
    }
  ],
  "result": {
    "exit_code": 0
  }
}
```