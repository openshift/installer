package datahub

type CreateTopicParameter struct {
    ShardCount   int
    LifeCycle    int
    Comment      string
    RecordType   RecordType
    RecordSchema *RecordSchema
    ExpandMode   ExpandMode
}

type UpdateTopicParameter struct {
    LifeCycle int
    Comment   string
}

type CreateConnectorParameter struct {
    SinkStartTime int64
    ConnectorType ConnectorType
    ColumnFields  []string
    ColumnNameMap map[string]string
    Config        interface{}
}

type UpdateConnectorParameter struct {
    ColumnFields  []string
    ColumnNameMap map[string]string
    Config        interface{}
}
