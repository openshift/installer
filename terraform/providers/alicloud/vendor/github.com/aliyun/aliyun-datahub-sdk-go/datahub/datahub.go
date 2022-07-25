package datahub

func NewClientWithConfig(endpoint string, config *Config, account Account) DataHubApi {
    if config.UserAgent == "" {
        config.UserAgent = DefaultUserAgent()
    }
    if config.HttpClient == nil {
        config.HttpClient = DefaultHttpClient()
    }
    if !validateCompressorType(config.CompressorType) {
        config.CompressorType = NOCOMPRESS
    }

    dh := &DataHub{
        Client: NewRestClient(endpoint, config.UserAgent, config.HttpClient,
            account, config.CompressorType),
        cType: config.CompressorType,
    }

    if config.EnableSchemaRegistry {
        dh.schemaClient = NewSchemaClient(dh)

        // compress data in batch record, no need to compress http body
        if config.CompressorType != NOCOMPRESS {
            dh.Client.CompressorType = NOCOMPRESS
        }

        return &DataHubBatch{
            DataHub: *dh,
        }
    } else {
        if config.EnableBinary {
            return &DataHubPB{
                DataHub: *dh,
            }
        } else {
            return dh
        }
    }
}

func New(accessId, accessKey, endpoint string) DataHubApi {
    config := NewDefaultConfig()
    return NewClientWithConfig(endpoint, config, NewAliyunAccount(accessId, accessKey))
}

func NewBatchClient(accessId, accessKey, endpoint string) DataHubApi {
    config := NewDefaultConfig()
    config.EnableSchemaRegistry = true
    config.CompressorType = LZ4
    return NewClientWithConfig(endpoint, config, NewAliyunAccount(accessId, accessKey))
}

// Datahub provides restful apis for visiting examples service.
type DataHubApi interface {
    // List all projects the user owns.
    ListProject() (*ListProjectResult, error)

    // List all projects the user owns with filter.
    ListProjectWithFilter(filter string) (*ListProjectResult, error)

    // Create a examples project.
    CreateProject(projectName, comment string) (*CreateProjectResult, error)

    // Update project information. Only support comment
    UpdateProject(projectName, comment string) (*UpdateProjectResult, error)

    // Delete the specified project. If any topics exist in the project, the delete operation will fail.
    DeleteProject(projectName string) (*DeleteProjectResult, error)

    // Get the information of the specified project.
    GetProject(projectName string) (*GetProjectResult, error)

    // Update project vpc white list.
    UpdateProjectVpcWhitelist(projectName, vpcIds string) (*UpdateProjectVpcWhitelistResult, error)

    // Wait for all shards' status of this topic is ACTIVE. Default timeout is 60s.
    WaitAllShardsReady(projectName, topicName string) bool

    // Wait for all shards' status of this topic is ACTIVE.
    // The unit is seconds.
    // If timeout < 0, it will block util all shards ready
    WaitAllShardsReadyWithTime(projectName, topicName string, timeout int64) bool

    // List all topics in the project.
    ListTopic(projectName string) (*ListTopicResult, error)

    // List all topics in the project with filter.
    ListTopicWithFilter(projectName, filter string) (*ListTopicResult, error)

    // Create a examples topic with type: BLOB
    CreateBlobTopic(projectName, topicName, comment string, shardCount, lifeCycle int) (*CreateBlobTopicResult, error)

    // Create a examples topic with type: TUPLE
    CreateTupleTopic(projectName, topicName, comment string, shardCount, lifeCycle int, recordSchema *RecordSchema) (*CreateTupleTopicResult, error)

    // Create topic with specific parameter
    CreateTopicWithPara(projectName, topicName string, para *CreateTopicParameter) (*CreateTopicWithParaResult, error)

    // Update topic meta information.
    UpdateTopic(projectName, topicName, comment string) (*UpdateTopicResult, error)

    // Update topic meta information. Only support comment and lifeCycle now.
    UpdateTopicWithPara(projectName, topicName string, para *UpdateTopicParameter) (*UpdateTopicResult, error)

    // Delete a specified topic.
    DeleteTopic(projectName, topicName string) (*DeleteTopicResult, error)

    // Get the information of the specified topic.
    GetTopic(projectName, topicName string) (*GetTopicResult, error)

    // List shard information {ShardEntry} of a topic.
    ListShard(projectName, topicName string) (*ListShardResult, error)

    // Split a shard. In function, sdk will automatically compute the split key which is used to split shard.
    SplitShard(projectName, topicName, shardId string) (*SplitShardResult, error)

    // Split a shard by the specified splitKey.
    SplitShardBySplitKey(projectName, topicName, shardId, splitKey string) (*SplitShardResult, error)

    // Merge the specified shard and its adjacent shard. Only adjacent shards can be merged.
    MergeShard(projectName, topicName, shardId, adjacentShardId string) (*MergeShardResult, error)

    // Extend shard num.
    ExtendShard(projectName, topicName string, shardCount int) (*ExtendShardResult, error)

    // Get the data cursor of a shard. This function support OLDEST, LATEST, SYSTEM_TIME and SEQUENCE.
    // If choose OLDEST or LATEST, the last parameter will not be needed.
    // if choose SYSTEM_TIME or SEQUENCE. it needs to a parameter as sequence num or timestamp.
    GetCursor(projectName, topicName, shardId string, ctype CursorType, param ...int64) (*GetCursorResult, error)

    // Write data records into a DataHub topic.
    // The PutRecordsResult includes unsuccessfully processed records.
    // Datahub attempts to process all records in each record.
    // A single record failure does not stop the processing of subsequent records.
    PutRecords(projectName, topicName string, records []IRecord) (*PutRecordsResult, error)

    PutRecordsByShard(projectName, topicName, shardId string, records []IRecord) (*PutRecordsByShardResult, error)

    // Get the TUPLE records of a shard.
    GetTupleRecords(projectName, topicName, shardId, cursor string, limit int, recordSchema *RecordSchema) (*GetRecordsResult, error)

    // Get the BLOB records of a shard.
    GetBlobRecords(projectName, topicName, shardId, cursor string, limit int) (*GetRecordsResult, error)

    // Append a field to a TUPLE topic.
    // Field AllowNull should be true.
    AppendField(projectName, topicName string, field Field) (*AppendFieldResult, error)

    // Get metering info of the specified shard
    GetMeterInfo(projectName, topicName, shardId string) (*GetMeterInfoResult, error)

    // List name of connectors.
    ListConnector(projectName, topicName string) (*ListConnectorResult, error)

    // Create data connectors.
    CreateConnector(projectName, topicName string, cType ConnectorType, columnFields []string, config interface{}) (*CreateConnectorResult, error)

    // Create connector with start time(unit:ms)
    CreateConnectorWithStartTime(projectName, topicName string, cType ConnectorType,
        columnFields []string, sinkStartTime int64, config interface{}) (*CreateConnectorResult, error)

    // Create connector with parameter
    CreateConnectorWithPara(projectName, topicName string, para *CreateConnectorParameter) (*CreateConnectorResult, error)

    // Update connector config of the specified data connector.
    // Config should be SinkOdpsConfig, SinkOssConfig ...
    UpdateConnector(projectName, topicName, connectorId string, config interface{}) (*UpdateConnectorResult, error)

    // Update connector with parameter
    UpdateConnectorWithPara(projectName, topicName, connectorId string, para *UpdateConnectorParameter) (*UpdateConnectorResult, error)

    // Delete a data connector.
    DeleteConnector(projectName, topicName, connectorId string) (*DeleteConnectorResult, error)

    // Get information of the specified data connector.
    GetConnector(projectName, topicName, connectorId string) (*GetConnectorResult, error)

    // Get the done time of a data connector. This method mainly used to get MaxCompute synchronize point.
    GetConnectorDoneTime(projectName, topicName, connectorId string) (*GetConnectorDoneTimeResult, error)

    // Get the detail information of the shard task which belongs to the specified data connector.
    GetConnectorShardStatus(projectName, topicName, connectorId string) (*GetConnectorShardStatusResult, error)

    // Get the detail information of the shard task which belongs to the specified data connector.
    GetConnectorShardStatusByShard(projectName, topicName, connectorId, shardId string) (*GetConnectorShardStatusByShardResult, error)

    // Reload a data connector.
    ReloadConnector(projectName, topicName, connectorId string) (*ReloadConnectorResult, error)

    // Reload the specified shard of the data connector.
    ReloadConnectorByShard(projectName, topicName, connectorId, shardId string) (*ReloadConnectorByShardResult, error)

    // Update the state of the data connector
    UpdateConnectorState(projectName, topicName, connectorId string, state ConnectorState) (*UpdateConnectorStateResult, error)

    // Update connector sink offset. The operation must be operated after connector stopped.
    UpdateConnectorOffset(projectName, topicName, connectorId, shardId string, offset ConnectorOffset) (*UpdateConnectorOffsetResult, error)

    // Append data connector field.
    // Before run this method, you should ensure that this field is in both the topic and the connector.
    AppendConnectorField(projectName, topicName, connectorId, fieldName string) (*AppendConnectorFieldResult, error)

    // List subscriptions in the topic.
    ListSubscription(projectName, topicName string, pageIndex, pageSize int) (*ListSubscriptionResult, error)

    // Create a subscription, and then you should commit offsets with this subscription.
    CreateSubscription(projectName, topicName, comment string) (*CreateSubscriptionResult, error)

    // Update a subscription. Now only support update comment information.
    UpdateSubscription(projectName, topicName, subId, comment string) (*UpdateSubscriptionResult, error)

    // Delete a subscription.
    DeleteSubscription(projectName, topicName, subId string) (*DeleteSubscriptionResult, error)

    // Get the detail information of a subscription.
    GetSubscription(projectName, topicName, subId string) (*GetSubscriptionResult, error)

    // Update a subscription' state. You can change the state of a subscription to SUB_ONLINE or SUB_OFFLINE.
    // When offline, you can not commit offsets of the subscription.
    UpdateSubscriptionState(projectName, topicName, subId string, state SubscriptionState) (*UpdateSubscriptionStateResult, error)

    // Init and get a subscription session, and returns offset if any offset stored before.
    // Subscription should be initialized before use. This operation makes sure that only one client use this subscription.
    // If this function be called in elsewhere, the seesion will be invalid and can not commit offsets of the subscription.
    OpenSubscriptionSession(projectName, topicName, subId string, shardIds []string) (*OpenSubscriptionSessionResult, error)

    // Get offsets of a subscription.This method dost not return sessionId in SubscriptionOffset.
    // Only the SubscriptionOffset containing sessionId can commit offset.
    GetSubscriptionOffset(projectName, topicName, subId string, shardIds []string) (*GetSubscriptionOffsetResult, error)

    // Update offsets of shards to server. This operation allows you store offsets on the server side.
    CommitSubscriptionOffset(projectName, topicName, subId string, offsets map[string]SubscriptionOffset) (*CommitSubscriptionOffsetResult, error)

    // Reset offsets of shards to server. This operation allows you reset offsets on the server side.
    ResetSubscriptionOffset(projectName, topicName, subId string, offsets map[string]SubscriptionOffset) (*ResetSubscriptionOffsetResult, error)

    // Heartbeat request to let server know consumer status.
    Heartbeat(projectName, topicName, consumerGroup, consumerId string, versionId int64, holdShardList, readEndShardList []string) (*HeartbeatResult, error)

    // Join a consumer group.
    JoinGroup(projectName, topicName, consumerGroup string, sessionTimeout int64) (*JoinGroupResult, error)

    // Sync consumer group info.
    SyncGroup(projectName, topicName, consumerGroup, consumerId string, versionId int64, releaseShardList, readEndShardList []string) (*SyncGroupResult, error)

    // Leave consumer group info.
    LeaveGroup(projectName, topicName, consumerGroup, consumerId string, versionId int64) (*LeaveGroupResult, error)

    // List topic schema.
    ListTopicSchema(projectName, topicName string) (*ListTopicSchemaResult, error)

    // Get topic schema by versionId.
    GetTopicSchemaByVersion(projectName, topicName string, versionId int) (*GetTopicSchemaResult, error)

    // Get topic schema by schema string.
    GetTopicSchemaBySchema(projectName, topicName string, recordSchema *RecordSchema) (*GetTopicSchemaResult, error)

    // Register schema to a topic.
    RegisterTopicSchema(projectName, topicName string, recordSchema *RecordSchema) (*RegisterTopicSchemaResult, error)

    // Delete topic schema by versionId
    DeleteTopicSchema(projectName, topicName string, versionId int) (*DeleteTopicSchemaResult, error)
}
