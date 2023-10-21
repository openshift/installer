package datahub

import (
    "errors"
    "fmt"
    "github.com/aliyun/aliyun-datahub-sdk-go/datahub/util"
    "time"
)

type DataHub struct {
    Client *RestClient

    // for batch client
    cType        CompressorType
    schemaClient *schemaRegistryClient
}

// ListProjects list all projects
func (datahub *DataHub) ListProject() (*ListProjectResult, error) {
    path := projectsPath
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }

    responseBody, commonResp, err := datahub.Client.Get(path, reqPara)
    if err != nil {
        return nil, err
    }
    return NewListProjectResult(responseBody, commonResp)
}

// ListProjects list projects with filter
func (datahub *DataHub) ListProjectWithFilter(filter string) (*ListProjectResult, error) {
    path := projectsPath
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
        Query:  map[string]string{httpFilterQuery: filter},
    }

    responseBody, commonResp, err := datahub.Client.Get(path, reqPara)
    if err != nil {
        return nil, err
    }
    return NewListProjectResult(responseBody, commonResp)
}

// CreateProject create new project
func (datahub *DataHub) CreateProject(projectName, comment string) (*CreateProjectResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckComment(comment) {
        return nil, NewInvalidParameterErrorWithMessage(commentInvalid)
    }

    path := fmt.Sprintf(projectPath, projectName)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    requestBody := &CreateProjectRequest{
        Comment: comment,
    }

    _, commonResp, err := datahub.Client.Post(path, requestBody, reqPara)
    if err != nil {
        return nil, err
    }
    return NewCreateProjectResult(commonResp)
}

// UpdateProject update project
func (datahub *DataHub) UpdateProject(projectName, comment string) (*UpdateProjectResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckComment(comment) {
        return nil, NewInvalidParameterErrorWithMessage(commentInvalid)
    }

    path := fmt.Sprintf(projectPath, projectName)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    requestBody := &UpdateProjectRequest{
        Comment: comment,
    }

    _, commonResp, err := datahub.Client.Put(path, requestBody, reqPara)
    if err != nil {
        return nil, err
    }
    return NewUpdateProjectResult(commonResp)
}

// DeleteProject delete project
func (datahub *DataHub) DeleteProject(projectName string) (*DeleteProjectResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }

    path := fmt.Sprintf(projectPath, projectName)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }

    _, commonResp, err := datahub.Client.Delete(path, reqPara)
    if err != nil {
        return nil, err
    }
    return NewDeleteProjectResult(commonResp)
}

// GetProject get a project deatil named the given name
func (datahub *DataHub) GetProject(projectName string) (*GetProjectResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }

    path := fmt.Sprintf(projectPath, projectName)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }

    respBody, commonResp, err := datahub.Client.Get(path, reqPara)
    if err != nil {
        return nil, err
    }

    result, err := NewGetProjectResult(respBody, commonResp)
    if err != nil {
        return nil, err
    }

    result.ProjectName = projectName
    return result, nil
}

// Update project vpc white list.
func (datahub *DataHub) UpdateProjectVpcWhitelist(projectName, vpcIds string) (*UpdateProjectVpcWhitelistResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }

    path := fmt.Sprintf(projectPath, projectName)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    requestBody := &UpdateProjectVpcWhitelistRequest{
        VpcIds: vpcIds,
    }

    _, commonResp, err := datahub.Client.Put(path, requestBody, reqPara)
    if err != nil {
        return nil, err
    }

    return NewUpdateProjectVpcWhitelistResult(commonResp)
}

func (datahub *DataHub) WaitAllShardsReady(projectName, topicName string) bool {
    return datahub.WaitAllShardsReadyWithTime(projectName, topicName, minWaitingTimeInMs/1000)
}

func (datahub *DataHub) WaitAllShardsReadyWithTime(projectName, topicName string, timeout int64) bool {
    ready := make(chan bool)
    if timeout > 0 {
        go func(timeout int64) {
            time.Sleep(time.Duration(timeout) * time.Second)
            ready <- false
        }(timeout)
    }
    go func(datahub DataHubApi) {
        for {
            ls, err := datahub.ListShard(projectName, topicName)
            if err != nil {
                time.Sleep(1 * time.Microsecond)
                continue
            }
            ok := true
            for _, shard := range ls.Shards {
                switch shard.State {
                case ACTIVE, CLOSED:
                    continue
                default:
                    ok = false
                    break
                }
            }
            if ok {
                break
            }
        }
        ready <- true
    }(datahub)
    return <-ready
}

func (datahub *DataHub) ListTopic(projectName string) (*ListTopicResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }

    path := fmt.Sprintf(topicsPath, projectName)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    respBody, commonResp, err := datahub.Client.Get(path, reqPara)
    if err != nil {
        return nil, err
    }
    return NewListTopicResult(respBody, commonResp)
}

func (datahub *DataHub) ListTopicWithFilter(projectName, filter string) (*ListTopicResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }

    path := fmt.Sprintf(topicsPath, projectName)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
        Query:  map[string]string{httpFilterQuery: filter},
    }
    respBody, commonResp, err := datahub.Client.Get(path, reqPara)
    if err != nil {
        return nil, err
    }
    return NewListTopicResult(respBody, commonResp)
}

func (datahub *DataHub) CreateBlobTopic(projectName, topicName, comment string, shardCount, lifeCycle int) (*CreateBlobTopicResult, error) {
    para := &CreateTopicParameter{
        ShardCount:   shardCount,
        LifeCycle:    lifeCycle,
        Comment:      comment,
        RecordType:   BLOB,
        RecordSchema: nil,
        ExpandMode:   SPLIT_EXTEND,
    }

    ret, err := datahub.CreateTopicWithPara(projectName, topicName, para)
    if err != nil {
        return nil, err
    }
    return NewCreateBlobTopicResult(&ret.CommonResponseResult)
}

func (datahub *DataHub) CreateTupleTopic(projectName, topicName, comment string, shardCount, lifeCycle int, recordSchema *RecordSchema) (*CreateTupleTopicResult, error) {
    para := &CreateTopicParameter{
        ShardCount:   shardCount,
        LifeCycle:    lifeCycle,
        Comment:      comment,
        RecordType:   TUPLE,
        RecordSchema: recordSchema,
        ExpandMode:   SPLIT_EXTEND,
    }

    ret, err := datahub.CreateTopicWithPara(projectName, topicName, para)
    if err != nil {
        return nil, err
    }
    return NewCreateTupleTopicResult(&ret.CommonResponseResult)
}

func (datahub *DataHub) CreateTopicWithPara(projectName, topicName string, para *CreateTopicParameter) (*CreateTopicWithParaResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if para == nil {
        return nil, NewInvalidParameterErrorWithMessage(parameterNull)
    }
    if !util.CheckComment(para.Comment) {
        return nil, NewInvalidParameterErrorWithMessage(commentInvalid)
    }
    if para.RecordType != TUPLE && para.RecordType != BLOB {
        return nil, NewInvalidParameterErrorWithMessage(fmt.Sprintf("Invalid RecordType: %s", para.RecordType))
    }
    if para.RecordType == TUPLE && para.RecordSchema == nil {
        return nil, NewInvalidParameterErrorWithMessage("Tuple topic must set RecordSchema")
    }
    if para.LifeCycle <= 0 {
        return nil, NewInvalidParameterErrorWithMessage(lifecycleInvalid)
    }

    path := fmt.Sprintf(topicPath, projectName, topicName)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    ctr := &CreateTopicRequest{
        Action:       "create",
        ShardCount:   para.ShardCount,
        Lifecycle:    para.LifeCycle,
        RecordType:   para.RecordType,
        RecordSchema: para.RecordSchema,
        Comment:      para.Comment,
        ExpandMode:   para.ExpandMode,
    }

    _, commonResp, err := datahub.Client.Post(path, ctr, reqPara)
    if err != nil {
        return nil, err
    }
    return NewCreateTopicWithParaResult(commonResp)
}

func (datahub *DataHub) UpdateTopic(projectName, topicName, comment string) (*UpdateTopicResult, error) {
    para := &UpdateTopicParameter{
        Comment: comment,
    }

    return datahub.UpdateTopicWithPara(projectName, topicName, para)
}

// Update topic meta information. Only support comment and lifeCycle now.
func (datahub *DataHub) UpdateTopicWithPara(projectName, topicName string, para *UpdateTopicParameter) (*UpdateTopicResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if para == nil {
        return nil, NewInvalidParameterErrorWithMessage(parameterNull)
    }
    if !util.CheckComment(para.Comment) {
        return nil, NewInvalidParameterErrorWithMessage(commentInvalid)
    }

    path := fmt.Sprintf(topicPath, projectName, topicName)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    ut := &UpdateTopicRequest{
        Lifecycle: para.LifeCycle,
        Comment:   para.Comment,
    }

    _, commonResp, err := datahub.Client.Put(path, ut, reqPara)
    if err != nil {
        return nil, err
    }
    return NewUpdateTopicResult(commonResp)
}

func (datahub *DataHub) DeleteTopic(projectName, topicName string) (*DeleteTopicResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(topicPath, projectName, topicName)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    _, commonResp, err := datahub.Client.Delete(path, reqPara)
    if err != nil {
        return nil, err
    }
    return NewDeleteTopicResult(commonResp)
}

func (datahub *DataHub) GetTopic(projectName, topicName string) (*GetTopicResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(topicPath, projectName, topicName)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    respBody, commonResp, err := datahub.Client.Get(path, reqPara)
    if err != nil {
        return nil, err
    }
    result, err := NewGetTopicResult(respBody, commonResp)

    if err != nil {
        return nil, err
    }
    result.ProjectName = projectName
    result.TopicName = topicName
    return result, nil
}

func (datahub *DataHub) ListShard(projectName, topicName string) (*ListShardResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(shardsPath, projectName, topicName)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    respBody, commonResp, err := datahub.Client.Get(path, reqPara)
    if err != nil {
        return nil, err
    }
    return NewListShardResult(respBody, commonResp)
}

func (datahub *DataHub) SplitShard(projectName, topicName, shardId string) (*SplitShardResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    if !util.CheckShardId(shardId) {
        return nil, NewInvalidParameterErrorWithMessage(shardIdInvalid)
    }

    splitKey, err := generateSpliteKey(projectName, topicName, shardId, datahub)
    if err != nil {
        return nil, err
    }
    path := fmt.Sprintf(shardsPath, projectName, topicName)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    ssr := &SplitShardRequest{
        Action:   "split",
        ShardId:  shardId,
        SplitKey: splitKey,
    }

    respBody, commonResp, err := datahub.Client.Post(path, ssr, reqPara)
    if err != nil {
        return nil, err
    }
    return NewSplitShardResult(respBody, commonResp)

}

func (datahub *DataHub) SplitShardBySplitKey(projectName, topicName, shardId, splitKey string) (*SplitShardResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    if !util.CheckShardId(shardId) {
        return nil, NewInvalidParameterErrorWithMessage(shardIdInvalid)
    }

    path := fmt.Sprintf(shardsPath, projectName, topicName)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    ssr := &SplitShardRequest{
        Action:   "split",
        ShardId:  shardId,
        SplitKey: splitKey,
    }

    respBody, commonResp, err := datahub.Client.Post(path, ssr, reqPara)
    if err != nil {
        return nil, err
    }
    return NewSplitShardResult(respBody, commonResp)
}

func (datahub *DataHub) MergeShard(projectName, topicName, shardId, adjacentShardId string) (*MergeShardResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if !util.CheckShardId(shardId) || !util.CheckShardId(adjacentShardId) {
        return nil, NewInvalidParameterErrorWithMessage(shardIdInvalid)
    }

    path := fmt.Sprintf(shardsPath, projectName, topicName)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    mss := &MergeShardRequest{
        Action:          "merge",
        ShardId:         shardId,
        AdjacentShardId: adjacentShardId,
    }

    respBody, commonResp, err := datahub.Client.Post(path, mss, reqPara)
    if err != nil {
        return nil, err
    }
    return NewMergeShardResult(respBody, commonResp)
}

func (datahub *DataHub) ExtendShard(projectName, topicName string, shardCount int) (*ExtendShardResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if shardCount <= 0 {
        return nil, NewInvalidParameterErrorWithMessage("shardCount is invalid")
    }

    path := fmt.Sprintf(shardsPath, projectName, topicName)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    mss := &ExtendShardRequest{
        Action:     "extend",
        ExtendMode: "TO",
        ShardCount: shardCount,
    }

    _, commonResp, err := datahub.Client.Post(path, mss, reqPara)
    if err != nil {
        return nil, err
    }
    return NewExtendShardResult(commonResp)
}

func (datahub *DataHub) GetCursor(projectName, topicName, shardId string, ctype CursorType, param ...int64) (*GetCursorResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if !util.CheckShardId(shardId) {
        return nil, NewInvalidParameterErrorWithMessage(shardIdInvalid)
    }
    if len(param) > 1 {
        return nil, NewInvalidParameterErrorWithMessage(parameterNumInvalid)
    }

    path := fmt.Sprintf(shardPath, projectName, topicName, shardId)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    gcr := &GetCursorRequest{
        Action:     "cursor",
        CursorType: ctype,
    }

    switch ctype {
    case OLDEST, LATEST:
        if len(param) != 0 {
            return nil, NewInvalidParameterErrorWithMessage("Not need extra parameter when CursorType OLDEST or LATEST")
        }
    case SYSTEM_TIME:
        if len(param) != 1 {
            return nil, NewInvalidParameterErrorWithMessage("Timestamp must be set when CursorType is SYSTEM_TIME")
        }
        gcr.SystemTime = param[0]
    case SEQUENCE:
        if len(param) != 1 {
            return nil, NewInvalidParameterErrorWithMessage("Sequence must be set when CursorType is SEQUENCE")
        }
        gcr.Sequence = param[0]
    }

    respBody, commonResp, err := datahub.Client.Post(path, gcr, reqPara)
    if err != nil {
        return nil, err
    }
    return NewGetCursorResult(respBody, commonResp)
}
func (datahub *DataHub) PutRecords(projectName, topicName string, records []IRecord) (*PutRecordsResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if records == nil || len(records) == 0 {
        return nil, NewInvalidParameterErrorWithMessage(recordsInvalid)
    }

    path := fmt.Sprintf(shardsPath, projectName, topicName)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    prr := &PutRecordsRequest{
        Action:  "pub",
        Records: records,
    }
    respBody, commonResp, err := datahub.Client.Post(path, prr, reqPara)
    if err != nil {
        return nil, err
    }
    return NewPutRecordsResult(respBody, commonResp)
}

func (datahub *DataHub) PutRecordsByShard(projectName, topicName, shardId string, records []IRecord) (*PutRecordsByShardResult, error) {
    return nil, errors.New("not support this method")
}

func (datahub *DataHub) GetTupleRecords(projectName, topicName, shardId, cursor string, limit int, recordSchema *RecordSchema) (*GetRecordsResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if !util.CheckShardId(shardId) {
        return nil, NewInvalidParameterErrorWithMessage(shardIdInvalid)
    }
    if recordSchema == nil {
        return nil, NewInvalidParameterErrorWithMessage(missingRecordSchema)
    }

    path := fmt.Sprintf(shardPath, projectName, topicName, shardId)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    grr := &GetRecordRequest{
        Action: "sub",
        Cursor: cursor,
        Limit:  limit,
    }
    respBody, commonResp, err := datahub.Client.Post(path, grr, reqPara)
    if err != nil {
        return nil, err
    }

    ret, err := NewGetRecordsResult(respBody, recordSchema, commonResp)
    if err != nil {
        return nil, err
    }

    for _, record := range ret.Records {
        if _, ok := record.(*TupleRecord); !ok {
            return nil, NewInvalidParameterErrorWithMessage("Shouldn't call this method for BLOB topic")
        }
    }
    return ret, nil
}

func (datahub *DataHub) GetBlobRecords(projectName, topicName, shardId, cursor string, limit int) (*GetRecordsResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if !util.CheckShardId(shardId) {
        return nil, NewInvalidParameterErrorWithMessage(shardIdInvalid)
    }

    path := fmt.Sprintf(shardPath, projectName, topicName, shardId)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    grr := &GetRecordRequest{
        Action: "sub",
        Cursor: cursor,
        Limit:  limit,
    }
    respBody, commonResp, err := datahub.Client.Post(path, grr, reqPara)
    if err != nil {
        return nil, err
    }
    return NewGetRecordsResult(respBody, nil, commonResp)
}

func (datahub *DataHub) AppendField(projectName, topicName string, field Field) (*AppendFieldResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(topicPath, projectName, topicName)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    afr := &AppendFieldRequest{
        Action:    "AppendField",
        FieldName: field.Name,
        FieldType: field.Type,
    }

    _, commonResp, err := datahub.Client.Post(path, afr, reqPara)
    if err != nil {
        return nil, err
    }
    return NewAppendFieldResult(commonResp)
}

func (datahub *DataHub) GetMeterInfo(projectName, topicName, shardId string) (*GetMeterInfoResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if !util.CheckShardId(shardId) {
        return nil, NewInvalidParameterErrorWithMessage(shardIdInvalid)
    }

    path := fmt.Sprintf(shardPath, projectName, topicName, shardId)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    gmir := &GetMeterInfoRequest{
        Action: "meter",
    }
    respBody, commonResp, err := datahub.Client.Post(path, gmir, reqPara)
    if err != nil {
        return nil, err
    }
    return NewGetMeterInfoResult(respBody, commonResp)
}

func (datahub *DataHub) ListConnector(projectName, topicName string) (*ListConnectorResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(connectorsPath, projectName, topicName)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
        Query:  map[string]string{httpHeaderConnectorMode: "id"},
    }
    respBody, commonResp, err := datahub.Client.Get(path, reqPara)
    if err != nil {
        return nil, err
    }
    return NewListConnectorResult(respBody, commonResp)
}

func (datahub *DataHub) CreateConnector(projectName, topicName string, cType ConnectorType, columnFields []string, config interface{}) (*CreateConnectorResult, error) {
    return datahub.CreateConnectorWithStartTime(projectName, topicName, cType, columnFields, -1, config)
}

func (datahub *DataHub) CreateConnectorWithStartTime(projectName, topicName string, cType ConnectorType,
    columnFields []string, sinkStartTime int64, config interface{}) (*CreateConnectorResult, error) {
    para := &CreateConnectorParameter{
        SinkStartTime: sinkStartTime,
        ConnectorType: cType,
        ColumnFields:  columnFields,
        ColumnNameMap: nil,
        Config:        config,
    }

    return datahub.CreateConnectorWithPara(projectName, topicName, para)
}

func (datahub *DataHub) CreateConnectorWithPara(projectName, topicName string, para *CreateConnectorParameter) (*CreateConnectorResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if para == nil {
        return nil, NewInvalidParameterErrorWithMessage(parameterNull)
    }
    if !validateConnectorType(para.ConnectorType) {
        return nil, NewInvalidParameterErrorWithMessage(parameterTypeInvalid)
    }

    path := fmt.Sprintf(connectorPath, projectName, topicName, para.ConnectorType.String())
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    ccr := &CreateConnectorRequest{
        Action:        "create",
        Type:          para.ConnectorType,
        SinkStartTime: para.SinkStartTime,
        ColumnFields:  para.ColumnFields,
        ColumnNameMap: para.ColumnNameMap,
        Config:        para.Config,
    }
    respBody, commonResp, err := datahub.Client.Post(path, ccr, reqPara)
    if err != nil {
        return nil, err
    }
    return NewCreateConnectorResult(respBody, commonResp)
}

func (datahub *DataHub) GetConnector(projectName, topicName, connectorId string) (*GetConnectorResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(connectorPath, projectName, topicName, connectorId)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    respBody, commonResp, err := datahub.Client.Get(path, reqPara)
    if err != nil {
        return nil, err
    }
    return NewGetConnectorResult(respBody, commonResp)
}

func (datahub *DataHub) UpdateConnector(projectName, topicName, connectorId string, config interface{}) (*UpdateConnectorResult, error) {
    para := &UpdateConnectorParameter{
        ColumnFields:  nil,
        ColumnNameMap: nil,
        Config:        config,
    }

    return datahub.UpdateConnectorWithPara(projectName, topicName, connectorId, para)
}

func (datahub *DataHub) UpdateConnectorWithPara(projectName, topicName, connectorId string, para *UpdateConnectorParameter) (*UpdateConnectorResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if para == nil {
        return nil, NewInvalidParameterErrorWithMessage(parameterNull)
    }

    path := fmt.Sprintf(connectorPath, projectName, topicName, connectorId)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    ucr := &UpdateConnectorRequest{
        Action:        "updateconfig",
        ColumnFields:  para.ColumnFields,
        ColumnNameMap: para.ColumnNameMap,
        Config:        para.Config,
    }
    _, commonResp, err := datahub.Client.Post(path, ucr, reqPara)
    if err != nil {
        return nil, err
    }
    return NewUpdateConnectorResult(commonResp)
}

func (datahub *DataHub) DeleteConnector(projectName, topicName, connectorId string) (*DeleteConnectorResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(connectorPath, projectName, topicName, connectorId)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    _, commonResp, err := datahub.Client.Delete(path, reqPara)
    if err != nil {
        return nil, err
    }
    return NewDeleteConnectorResult(commonResp)
}

func (datahub *DataHub) GetConnectorDoneTime(projectName, topicName, connectorId string) (*GetConnectorDoneTimeResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(connectorPath, projectName, topicName, connectorId)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
        Query:  map[string]string{"donetime": ""},
    }

    respBody, commonResp, err := datahub.Client.Get(path, reqPara)
    if err != nil {
        return nil, err
    }
    return NewGetConnectorDoneTimeResult(respBody, commonResp)
}

func (datahub *DataHub) GetConnectorShardStatus(projectName, topicName, connectorId string) (*GetConnectorShardStatusResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(connectorPath, projectName, topicName, connectorId)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    gcss := &GetConnectorShardStatusRequest{
        Action: "Status",
    }
    respBody, commonResp, err := datahub.Client.Post(path, gcss, reqPara)
    if err != nil {
        return nil, err
    }
    return NewGetConnectorShardStatusResult(respBody, commonResp)
}

func (datahub *DataHub) GetConnectorShardStatusByShard(projectName, topicName, connectorId, shardId string) (*GetConnectorShardStatusByShardResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if !util.CheckShardId(shardId) {
        return nil, NewInvalidParameterErrorWithMessage(shardIdInvalid)
    }

    path := fmt.Sprintf(connectorPath, projectName, topicName, connectorId)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    gcss := &GetConnectorShardStatusRequest{
        Action:  "Status",
        ShardId: shardId,
    }
    respBody, commonResp, err := datahub.Client.Post(path, gcss, reqPara)
    if err != nil {
        return nil, err
    }
    return NewGetConnectorShardStatusByShardResult(respBody, commonResp)
}

func (datahub *DataHub) ReloadConnector(projectName, topicName, connectorId string) (*ReloadConnectorResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(connectorPath, projectName, topicName, connectorId)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    rcr := &ReloadConnectorRequest{
        Action: "Reload",
    }
    _, commonResp, err := datahub.Client.Post(path, rcr, reqPara)
    if err != nil {
        return nil, err
    }
    return NewReloadConnectorResult(commonResp)
}

func (datahub *DataHub) ReloadConnectorByShard(projectName, topicName, connectorId, shardId string) (*ReloadConnectorByShardResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if !util.CheckShardId(shardId) {
        return nil, NewInvalidParameterErrorWithMessage(shardIdInvalid)
    }

    path := fmt.Sprintf(connectorPath, projectName, topicName, connectorId)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    rcr := &ReloadConnectorRequest{
        Action:  "Reload",
        ShardId: shardId,
    }
    _, commonResp, err := datahub.Client.Post(path, rcr, reqPara)
    if err != nil {
        return nil, err
    }
    return NewReloadConnectorByShardResult(commonResp)
}

func (datahub *DataHub) UpdateConnectorState(projectName, topicName, connectorId string, state ConnectorState) (*UpdateConnectorStateResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if !validateConnectorState(state) {
        return nil, NewInvalidParameterErrorWithMessage(parameterTypeInvalid)
    }

    path := fmt.Sprintf(connectorPath, projectName, topicName, connectorId)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    ucsr := &UpdateConnectorStateRequest{
        Action: "updatestate",
        State:  state,
    }
    _, commonResp, err := datahub.Client.Post(path, ucsr, reqPara)
    if err != nil {
        return nil, err
    }
    return NewUpdateConnectorStateResult(commonResp)
}

func (datahub *DataHub) UpdateConnectorOffset(projectName, topicName, connectorId, shardId string, offset ConnectorOffset) (*UpdateConnectorOffsetResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if !util.CheckShardId(shardId) {
        return nil, NewInvalidParameterErrorWithMessage(shardIdInvalid)
    }

    path := fmt.Sprintf(connectorPath, projectName, topicName, connectorId)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    ucor := &UpdateConnectorOffsetRequest{
        Action:    "updateshardcontext",
        ShardId:   shardId,
        Timestamp: offset.Timestamp,
        Sequence:  offset.Sequence,
    }

    _, commonResp, err := datahub.Client.Post(path, ucor, reqPara)
    if err != nil {
        return nil, err
    }
    return NewUpdateConnectorOffsetResult(commonResp)
}

func (datahub *DataHub) AppendConnectorField(projectName, topicName, connectorId, fieldName string) (*AppendConnectorFieldResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(connectorPath, projectName, topicName, connectorId)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    acfr := &AppendConnectorFieldRequest{
        Action:    "appendfield",
        FieldName: fieldName,
    }
    _, commonResp, err := datahub.Client.Post(path, acfr, reqPara)
    if err != nil {
        return nil, err
    }
    return NewAppendConnectorFieldResult(commonResp)
}

func (datahub *DataHub) ListSubscription(projectName, topicName string, pageIndex, pageSize int) (*ListSubscriptionResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(subscriptionsPath, projectName, topicName)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    lsr := &ListSubscriptionRequest{
        Action:    "list",
        PageIndex: pageIndex,
        PageSize:  pageSize,
    }
    respBody, commonResp, err := datahub.Client.Post(path, lsr, reqPara)
    if err != nil {
        return nil, err
    }
    return NewListSubscriptionResult(respBody, commonResp)
}

func (datahub *DataHub) CreateSubscription(projectName, topicName, comment string) (*CreateSubscriptionResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if !util.CheckComment(comment) {
        return nil, NewInvalidParameterErrorWithMessage(commentInvalid)
    }

    path := fmt.Sprintf(subscriptionsPath, projectName, topicName)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    csr := &CreateSubscriptionRequest{
        Action:  "create",
        Comment: comment,
    }
    respBody, commonResp, err := datahub.Client.Post(path, csr, reqPara)
    if err != nil {
        return nil, err
    }
    return NewCreateSubscriptionResult(respBody, commonResp)
}

func (datahub *DataHub) UpdateSubscription(projectName, topicName, subId, comment string) (*UpdateSubscriptionResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if !util.CheckComment(comment) {
        return nil, NewInvalidParameterErrorWithMessage(commentInvalid)
    }

    path := fmt.Sprintf(subscriptionPath, projectName, topicName, subId)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    usr := &UpdateSubscriptionRequest{
        Comment: comment,
    }
    _, commonResp, err := datahub.Client.Put(path, usr, reqPara)
    if err != nil {
        return nil, err
    }
    return NewUpdateSubscriptionResult(commonResp)
}

func (datahub *DataHub) DeleteSubscription(projectName, topicName, subId string) (*DeleteSubscriptionResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(subscriptionPath, projectName, topicName, subId)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    _, commonResp, err := datahub.Client.Delete(path, reqPara)
    if err != nil {
        return nil, err
    }
    return NewDeleteSubscriptionResult(commonResp)
}

func (datahub *DataHub) GetSubscription(projectName, topicName, subId string) (*GetSubscriptionResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(subscriptionPath, projectName, topicName, subId)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    respBody, commonResp, err := datahub.Client.Get(path, reqPara)
    if err != nil {
        return nil, err
    }
    return NewGetSubscriptionResult(respBody, commonResp)
}

func (datahub *DataHub) UpdateSubscriptionState(projectName, topicName, subId string, state SubscriptionState) (*UpdateSubscriptionStateResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(subscriptionPath, projectName, topicName, subId)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    usr := &UpdateSubscriptionStateRequest{
        State: state,
    }
    _, commonResp, err := datahub.Client.Put(path, usr, reqPara)
    if err != nil {
        return nil, err
    }
    return NewUpdateSubscriptionStateResult(commonResp)
}

func (datahub *DataHub) OpenSubscriptionSession(projectName, topicName, subId string, shardIds []string) (*OpenSubscriptionSessionResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    for _, id := range shardIds {
        if !util.CheckShardId(id) {
            return nil, NewInvalidParameterErrorWithMessage(shardIdInvalid)
        }
    }

    path := fmt.Sprintf(offsetsPath, projectName, topicName, subId)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    ossr := &OpenSubscriptionSessionRequest{
        Action:   "open",
        ShardIds: shardIds,
    }
    respBody, commonResp, err := datahub.Client.Post(path, ossr, reqPara)
    if err != nil {
        return nil, err
    }
    return NewOpenSubscriptionSessionResult(respBody, commonResp)
}

func (datahub *DataHub) GetSubscriptionOffset(projectName, topicName, subId string, shardIds []string) (*GetSubscriptionOffsetResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    for _, id := range shardIds {
        if !util.CheckShardId(id) {
            return nil, NewInvalidParameterErrorWithMessage(shardIdInvalid)
        }
    }

    path := fmt.Sprintf(offsetsPath, projectName, topicName, subId)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    gsor := &GetSubscriptionOffsetRequest{
        Action:   "get",
        ShardIds: shardIds,
    }
    respBody, commonResp, err := datahub.Client.Post(path, gsor, reqPara)
    if err != nil {
        return nil, err
    }
    return NewGetSubscriptionOffsetResult(respBody, commonResp)
}

func (datahub *DataHub) CommitSubscriptionOffset(projectName, topicName, subId string, offsets map[string]SubscriptionOffset) (*CommitSubscriptionOffsetResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(offsetsPath, projectName, topicName, subId)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    req := &CommitSubscriptionOffsetRequest{
        Action:  "commit",
        Offsets: offsets,
    }

    _, commonResp, err := datahub.Client.Put(path, req, reqPara)
    if err != nil {
        return nil, err
    }
    return NewCommitSubscriptionOffsetResult(commonResp)
}

func (datahub *DataHub) ResetSubscriptionOffset(projectName, topicName, subId string, offsets map[string]SubscriptionOffset) (*ResetSubscriptionOffsetResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(offsetsPath, projectName, topicName, subId)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    req := &ResetSubscriptionOffsetRequest{
        Action:  "reset",
        Offsets: offsets,
    }
    _, commonResp, err := datahub.Client.Put(path, req, reqPara)
    if err != nil {
        return nil, err
    }
    return NewResetSubscriptionOffsetResult(commonResp)
}

func (datahub *DataHub) Heartbeat(projectName, topicName, consumerGroup, consumerId string, versionId int64, holdShardList, readEndShardList []string) (*HeartbeatResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    for _, id := range holdShardList {
        if !util.CheckShardId(id) {
            return nil, NewInvalidParameterErrorWithMessage(shardIdInvalid)
        }
    }
    for _, id := range readEndShardList {
        if !util.CheckShardId(id) {
            return nil, NewInvalidParameterErrorWithMessage(shardIdInvalid)
        }
    }

    path := fmt.Sprintf(consumerGroupPath, projectName, topicName, consumerGroup)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    hr := &HeartbeatRequest{
        Action:           "heartbeat",
        ConsumerId:       consumerId,
        VersionId:        versionId,
        HoldShardList:    holdShardList,
        ReadEndShardList: readEndShardList,
    }

    respBody, commonResp, err := datahub.Client.Post(path, hr, reqPara)
    if err != nil {
        return nil, err
    }
    return NewHeartbeatResult(respBody, commonResp)
}

func (datahub *DataHub) JoinGroup(projectName, topicName, consumerGroup string, sessionTimeout int64) (*JoinGroupResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(consumerGroupPath, projectName, topicName, consumerGroup)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    jgr := &JoinGroupRequest{
        Action:         "joinGroup",
        SessionTimeout: sessionTimeout,
    }
    respBody, commonResp, err := datahub.Client.Post(path, jgr, reqPara)
    if err != nil {
        return nil, err
    }
    return NewJoinGroupResult(respBody, commonResp)

}
func (datahub *DataHub) SyncGroup(projectName, topicName, consumerGroup, consumerId string, versionId int64, releaseShardList, readEndShardList []string) (*SyncGroupResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if len(releaseShardList) == 0 || len(readEndShardList) == 0 {
        return nil, NewInvalidParameterErrorWithMessage(shardListInvalid)
    }
    for _, id := range releaseShardList {
        if !util.CheckShardId(id) {
            return nil, NewInvalidParameterErrorWithMessage(shardIdInvalid)
        }
    }
    for _, id := range readEndShardList {
        if !util.CheckShardId(id) {
            return nil, NewInvalidParameterErrorWithMessage(shardIdInvalid)
        }
    }

    path := fmt.Sprintf(consumerGroupPath, projectName, topicName, consumerGroup)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    sgr := &SyncGroupRequest{
        Action:           "syncGroup",
        ConsumerId:       consumerId,
        VersionId:        versionId,
        ReleaseShardList: releaseShardList,
        ReadEndShardList: readEndShardList,
    }
    _, commonResp, err := datahub.Client.Post(path, sgr, reqPara)
    if err != nil {
        return nil, err
    }
    return NewSyncGroupResult(commonResp)
}

func (datahub *DataHub) LeaveGroup(projectName, topicName, consumerGroup, consumerId string, versionId int64) (*LeaveGroupResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(consumerGroupPath, projectName, topicName, consumerGroup)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    lgr := &LeaveGroupRequest{
        Action:     "leaveGroup",
        ConsumerId: consumerId,
        VersionId:  versionId,
    }
    _, commonResp, err := datahub.Client.Post(path, lgr, reqPara)
    if err != nil {
        return nil, err
    }
    return NewLeaveGroupResult(commonResp)
}

func (datahub *DataHub) ListTopicSchema(projectName, topicName string) (*ListTopicSchemaResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(topicPath, projectName, topicName)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    lts := &ListTopicSchemaRequest{
        Action: "ListSchema",
    }

    respBody, commonResp, err := datahub.Client.Post(path, lts, reqPara)
    if err != nil {
        return nil, err
    }
    return NewListTopicSchemaResult(respBody, commonResp)
}

func (datahub *DataHub) GetTopicSchemaByVersion(projectName, topicName string, versionId int) (*GetTopicSchemaResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(topicPath, projectName, topicName)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    lts := &GetTopicSchemaRequest{
        Action:       "GetSchema",
        VersionId:    versionId,
        RecordSchema: nil,
    }

    respBody, commonResp, err := datahub.Client.Post(path, lts, reqPara)
    if err != nil {
        return nil, err
    }
    return NewGetTopicSchemaResult(respBody, commonResp)
}

func (datahub *DataHub) GetTopicSchemaBySchema(projectName, topicName string, recordSchema *RecordSchema) (*GetTopicSchemaResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(topicPath, projectName, topicName)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    lts := &GetTopicSchemaRequest{
        Action:       "GetSchema",
        VersionId:    -1,
        RecordSchema: recordSchema,
    }

    respBody, commonResp, err := datahub.Client.Post(path, lts, reqPara)
    if err != nil {
        return nil, err
    }
    return NewGetTopicSchemaResult(respBody, commonResp)
}

func (datahub *DataHub) RegisterTopicSchema(projectName, topicName string, recordSchema *RecordSchema) (*RegisterTopicSchemaResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(topicPath, projectName, topicName)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    lts := &RegisterTopicSchemaRequest{
        Action:       "RegisterSchema",
        RecordSchema: recordSchema,
    }

    respBody, commonResp, err := datahub.Client.Post(path, lts, reqPara)
    if err != nil {
        return nil, err
    }
    return NewRegisterTopicSchemaResult(respBody, commonResp)
}

func (datahub *DataHub) DeleteTopicSchema(projectName, topicName string, versionId int) (*DeleteTopicSchemaResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(topicPath, projectName, topicName)
    reqPara := &RequestParameter{
        Header: map[string]string{httpHeaderContentType: httpJsonContent},
    }
    lts := &DeleteTopicSchemaRequest{
        Action:    "DeleteSchema",
        VersionId: versionId,
    }

    _, commonResp, err := datahub.Client.Post(path, lts, reqPara)
    if err != nil {
        return nil, err
    }
    return NewDeleteTopicSchemaResult(commonResp)
}

func (datahub *DataHub) getSchemaRegistry() *schemaRegistryClient {
    return datahub.schemaClient
}

type DataHubPB struct {
    DataHub
}

func (datahub *DataHubPB) PutRecords(projectName, topicName string, records []IRecord) (*PutRecordsResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(shardsPath, projectName, topicName)
    reqPara := &RequestParameter{
        Header: map[string]string{
            httpHeaderContentType:   httpProtoContent,
            httpHeaderRequestAction: httpPublistContent},
    }
    prr := &PutPBRecordsRequest{
        Records: records,
    }
    respBody, commonResp, err := datahub.Client.Post(path, prr, reqPara)
    if err != nil {
        return nil, err
    }
    return NewPutPBRecordsResult(respBody, commonResp)
}

func (datahub *DataHubPB) PutRecordsByShard(projectName, topicName, shardId string, records []IRecord) (*PutRecordsByShardResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if !util.CheckShardId(shardId) {
        return nil, NewInvalidParameterErrorWithMessage(shardIdInvalid)
    }

    path := fmt.Sprintf(shardPath, projectName, topicName, shardId)
    reqPara := &RequestParameter{
        Header: map[string]string{
            httpHeaderContentType:   httpProtoContent,
            httpHeaderRequestAction: httpPublistContent},
    }
    prr := &PutPBRecordsRequest{
        Records: records,
    }

    _, commonResp, err := datahub.Client.Post(path, prr, reqPara)
    if err != nil {
        return nil, err
    }
    return NewPutRecordsByShardResult(commonResp)
}

func (datahub *DataHubPB) GetTupleRecords(projectName, topicName, shardId, cursor string, limit int, recordSchema *RecordSchema) (*GetRecordsResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if !util.CheckShardId(shardId) {
        return nil, NewInvalidParameterErrorWithMessage(shardIdInvalid)
    }

    path := fmt.Sprintf(shardPath, projectName, topicName, shardId)
    reqPara := &RequestParameter{
        Header: map[string]string{
            httpHeaderContentType:   httpProtoContent,
            httpHeaderRequestAction: httpSubscribeContent},
    }
    grr := &GetPBRecordRequest{
        Cursor: cursor,
        Limit:  limit,
    }
    respBody, commonResp, err := datahub.Client.Post(path, grr, reqPara)
    if err != nil {
        return nil, err
    }
    return NewGetPBRecordsResult(respBody, recordSchema, commonResp)
}

func (datahub *DataHubPB) GetBlobRecords(projectName, topicName, shardId, cursor string, limit int) (*GetRecordsResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if !util.CheckShardId(shardId) {
        return nil, NewInvalidParameterErrorWithMessage(shardIdInvalid)
    }

    path := fmt.Sprintf(shardPath, projectName, topicName, shardId)
    reqPara := &RequestParameter{
        Header: map[string]string{
            httpHeaderContentType:   httpProtoContent,
            httpHeaderRequestAction: httpSubscribeContent},
    }
    grr := &GetPBRecordRequest{
        Cursor: cursor,
        Limit:  limit,
    }
    respBody, commonResp, err := datahub.Client.Post(path, grr, reqPara)
    if err != nil {
        return nil, err
    }
    return NewGetPBRecordsResult(respBody, nil, commonResp)
}

type DataHubBatch struct {
    DataHub
}

func (datahub *DataHubBatch) PutRecords(projectName, topicName string, records []IRecord) (*PutRecordsResult, error) {
    return nil, errors.New("not support this method")
}

func (datahub *DataHubBatch) PutRecordsByShard(projectName, topicName, shardId string, records []IRecord) (*PutRecordsByShardResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if !util.CheckShardId(shardId) {
        return nil, NewInvalidParameterErrorWithMessage(shardIdInvalid)
    }

    path := fmt.Sprintf(shardPath, projectName, topicName, shardId)
    reqPara := &RequestParameter{
        Header: map[string]string{
            httpHeaderContentType:   httpProtoBatchContent,
            httpHeaderRequestAction: httpPublistContent},
    }

    serializer := newBatchSerializer(projectName, topicName, datahub.cType, datahub.schemaClient)
    prr := &PutBatchRecordsRequest{
        serializer: serializer,
        Records:    records,
    }

    _, commonResp, err := datahub.Client.Post(path, prr, reqPara)
    if err != nil {
        return nil, err
    }
    return NewPutRecordsByShardResult(commonResp)
}

func (datahub *DataHubBatch) GetTupleRecords(projectName, topicName, shardId, cursor string, limit int, recordSchema *RecordSchema) (*GetRecordsResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if !util.CheckShardId(shardId) {
        return nil, NewInvalidParameterErrorWithMessage(shardIdInvalid)
    }

    path := fmt.Sprintf(shardPath, projectName, topicName, shardId)
    reqPara := &RequestParameter{
        Header: map[string]string{
            httpHeaderContentType:   httpProtoBatchContent,
            httpHeaderRequestAction: httpSubscribeContent},
    }
    gbr := &GetBatchRecordRequest{
        GetPBRecordRequest{
            Cursor: cursor,
            Limit:  limit,
        },
    }

    respBody, commonResp, err := datahub.Client.Post(path, gbr, reqPara)
    if err != nil {
        return nil, err
    }

    deserializer := newBatchDeserializer(projectName, topicName, shardId, recordSchema, datahub.schemaClient)
    return NewGetBatchRecordsResult(respBody, recordSchema, commonResp, deserializer)
}

func (datahub *DataHubBatch) GetBlobRecords(projectName, topicName, shardId, cursor string, limit int) (*GetRecordsResult, error) {
    return datahub.GetTupleRecords(projectName, topicName, shardId, cursor, limit, nil)
}
