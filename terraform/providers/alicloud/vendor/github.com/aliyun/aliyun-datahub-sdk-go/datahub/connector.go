package datahub

import (
    "bytes"
    "encoding/json"
    "errors"
    "fmt"
    "reflect"
    "strconv"
)

type AuthMode string

const (
    AK  AuthMode = "ak"
    STS AuthMode = "sts"
)

type ConnectorType string

const (
    SinkOdps     ConnectorType = "sink_odps"
    SinkOss      ConnectorType = "sink_oss"
    SinkEs       ConnectorType = "sink_es"
    SinkAds      ConnectorType = "sink_ads"
    SinkMysql    ConnectorType = "sink_mysql"
    SinkFc       ConnectorType = "sink_fc"
    SinkOts      ConnectorType = "sink_ots"
    SinkDatahub  ConnectorType = "sink_datahub"
    SinkHologres ConnectorType = "sink_hologres"
)

func (ct *ConnectorType) String() string {
    return string(*ct)
}

func validateConnectorType(ct ConnectorType) bool {
    switch ct {
    case SinkOdps, SinkOss, SinkEs, SinkAds, SinkMysql, SinkFc, SinkOts, SinkDatahub, SinkHologres:
        return true
    default:
        return false
    }
}

type ConnectorState string

const (
    ConnectorStopped ConnectorState = "CONNECTOR_STOPPED"
    ConnectorRunning ConnectorState = "CONNECTOR_RUNNING"
)

func validateConnectorState(ct ConnectorState) bool {
    switch ct {
    case ConnectorStopped, ConnectorRunning:
        return true
    default:
        return false
    }
}

type ConnectorTimestampUnit string

const (
    ConnectorMicrosecond ConnectorTimestampUnit = "MICROSECOND"
    ConnectorMillisecond ConnectorTimestampUnit = "MILLISECOND"
    ConnectorSecond      ConnectorTimestampUnit = "SECOND"
)

type ConnectorConfig struct {
    TimestampUnit ConnectorTimestampUnit `json:"TimestampUnit"`
}

type PartitionMode string

const (
    UserDefineMode PartitionMode = "USER_DEFINE"
    SystemTimeMode PartitionMode = "SYSTEM_TIME"
    EventTimeMode  PartitionMode = "EVENT_TIME"
)

func (pm *PartitionMode) String() string {
    return string(*pm)
}

func NewPartitionConfig() *PartitionConfig {
    pc := &PartitionConfig{
        ConfigMap: make([]map[string]string, 0, 0),
    }
    return pc
}

type PartitionConfig struct {
    ConfigMap []map[string]string
}

func (pc *PartitionConfig) AddConfig(key, value string) {
    m := map[string]string{
        key: value,
    }
    pc.ConfigMap = append(pc.ConfigMap, m)
}

func (pc *PartitionConfig) MarshalJSON() ([]byte, error) {
    if pc == nil || len(pc.ConfigMap) == 0 {
        return nil, nil
    }
    buf := &bytes.Buffer{}
    buf.Write([]byte{'{'})

    length := len(pc.ConfigMap)
    for i, m := range pc.ConfigMap {
        for k, v := range m {
            if _, err := fmt.Fprintf(buf, "\"%s\":\"%s\"", k, v); err != nil {
                return nil, errors.New(fmt.Sprintf("partition config is invalid"))
            }
        }
        if i < length-1 {
            buf.WriteByte(',')
        }
    }
    buf.WriteByte('}')

    return buf.Bytes(), nil
}

func (pc *PartitionConfig) UnmarshalJSON(data []byte) error {
    //the data is "xxxxxx",should convert to xxxx, remove the ""
    var str *string = new(string)
    if err := json.Unmarshal(data, str); err != nil {
        return err
    }

    confParser := make([]map[string]string, 0)
    if err := json.Unmarshal([]byte(*str), &confParser); err != nil {
        return err
    }
    confMap := make([]map[string]string, len(confParser))

    //convert {"key":"ds","value":"%Y%m%d",...} to {"ds":"%Y%m%d",...}
    for i, m := range confParser {
        confMap[i] = map[string]string{
            m["key"]: m["value"],
        }
    }
    pc.ConfigMap = confMap
    return nil
}

/** ODPS CONFIG **/
type SinkOdpsConfig struct {
    ConnectorConfig
    Endpoint        string          `json:"OdpsEndpoint"`
    Project         string          `json:"Project"`
    Table           string          `json:"Table"`
    AccessId        string          `json:"AccessId"`
    AccessKey       string          `json:"AccessKey"`
    TimeRange       int             `json:"TimeRange"`
    TimeZone        string          `json:"TimeZone,omitempty"`
    PartitionMode   PartitionMode   `json:"PartitionMode"`
    PartitionConfig PartitionConfig `json:"PartitionConfig"`
    TunnelEndpoint  string          `json:"TunnelEndpoint,omitempty"`
    SplitKey        string          `json:"SplitKey,omitempty"`
    Base64Encode    bool            `json:"Base64Encode,omitempty"`
}

func marshalCreateOdpsConnector(ccr *CreateConnectorRequest) ([]byte, error) {
    soConf, ok := ccr.Config.(SinkOdpsConfig)
    if !ok {
        return nil, NewInvalidParameterErrorWithMessage(fmt.Sprintf("config type error,your input config type is %s,should be SinkOdpsConfig", reflect.TypeOf(ccr.Config)))
    }

    // set default value
    if soConf.TimestampUnit == "" {
        soConf.TimestampUnit = ConnectorMicrosecond
    }

    ct := &struct {
        Action        string            `json:"Action"`
        Type          string            `json:"Type"`
        SinkStartTime int64             `json:"SinkStartTime"`
        ColumnFields  []string          `json:"ColumnFields"`
        ColumnNameMap map[string]string `json:"ColumnNameMap,omitempty"`
        Config        SinkOdpsConfig    `json:"Config"`
    }{
        Action:        ccr.Action,
        Type:          ccr.Type.String(),
        SinkStartTime: ccr.SinkStartTime,
        ColumnFields:  ccr.ColumnFields,
        ColumnNameMap: ccr.ColumnNameMap,
        Config:        soConf,
    }
    return json.Marshal(ct)
}

func unmarshalGetOdpsConnector(commonResp *CommonResponseResult, data []byte) (*GetConnectorResult, error) {
    //the api return TimeRange is string, so need to convert to int64
    type SinkOdpsConfigHelper struct {
        SinkOdpsConfig
        TimeRange string `json:"TimeRange"`
    }
    ct := &struct {
        GetConnectorResult
        Config SinkOdpsConfigHelper `json:"Config"`
    }{}

    if err := json.Unmarshal(data, ct); err != nil {
        return nil, err
    }

    conf := ct.Config.SinkOdpsConfig
    t, err := strconv.Atoi(ct.Config.TimeRange)
    if err != nil {
        return nil, err
    }
    conf.TimeRange = t

    ret := &ct.GetConnectorResult
    ret.Config = conf
    ret.CommonResponseResult = *commonResp
    return ret, nil
}

// no config update
func marshalUpdateConnector(ucr *UpdateConnectorRequest) ([]byte, error) {
    ct := &struct {
        Action        string            `json:"Action"`
        ColumnFields  []string          `json:"ColumnFields,omitempty"`
        ColumnNameMap map[string]string `json:"ColumnNameMap,omitempty"`
    }{
        Action:        ucr.Action,
        ColumnFields:  ucr.ColumnFields,
        ColumnNameMap: ucr.ColumnNameMap,
    }
    return json.Marshal(ct)
}

func marshalUpdateOdpsConnector(ucr *UpdateConnectorRequest) ([]byte, error) {
    soConf, ok := ucr.Config.(SinkOdpsConfig)
    if !ok {
        return nil, NewInvalidParameterErrorWithMessage(fmt.Sprintf("config type error,your input config type is %s,should be SinkOdpsConfig", reflect.TypeOf(ucr.Config)))
    }

    // set default value
    if soConf.TimestampUnit == "" {
        soConf.TimestampUnit = ConnectorMicrosecond
    }

    ct := &struct {
        Action        string            `json:"Action"`
        ColumnFields  []string          `json:"ColumnFields,omitempty"`
        ColumnNameMap map[string]string `json:"ColumnNameMap,omitempty"`
        Config        SinkOdpsConfig    `json:"Config,omitempty"`
    }{
        Action:        ucr.Action,
        ColumnFields:  ucr.ColumnFields,
        ColumnNameMap: ucr.ColumnNameMap,
        Config:        soConf,
    }
    return json.Marshal(ct)
}

/*  Oss Config */
type SinkOssConfig struct {
    ConnectorConfig
    Endpoint    string   `json:"Endpoint"`
    Bucket      string   `json:"Bucket"`
    Prefix      string   `json:"Prefix"`
    TimeFormat  string   `json:"TimeFormat"`
    TimeRange   int      `json:"TimeRange"`
    AuthMode    AuthMode `json:"AuthMode"`
    AccessId    string   `json:"AccessId,omitempty"`
    AccessKey   string   `json:"AccessKey,omitempty"`
    MaxFileSize int64    `json:"MaxFileSize,omitempty"`
}

func marshalCreateOssConnector(ccr *CreateConnectorRequest) ([]byte, error) {
    soConf, ok := ccr.Config.(SinkOssConfig)
    if !ok {
        return nil, NewInvalidParameterErrorWithMessage(fmt.Sprintf("config type error,your input config type is %s,should be SinkOssConfig", reflect.TypeOf(ccr.Config)))
    }

    // set default value
    if soConf.TimestampUnit == "" {
        soConf.TimestampUnit = ConnectorMicrosecond
    }

    ct := &struct {
        Action        string            `json:"Action"`
        Type          ConnectorType     `json:"Type"`
        SinkStartTime int64             `json:"SinkStartTime"`
        ColumnFields  []string          `json:"ColumnFields"`
        ColumnNameMap map[string]string `json:"ColumnNameMap,omitempty"`
        Config        SinkOssConfig     `json:"Config"`
    }{
        Action:        "create",
        Type:          ccr.Type,
        SinkStartTime: ccr.SinkStartTime,
        ColumnFields:  ccr.ColumnFields,
        ColumnNameMap: ccr.ColumnNameMap,
        Config:        soConf,
    }
    return json.Marshal(ct)
}

func unmarshalGetOssConnector(commonResp *CommonResponseResult, data []byte) (*GetConnectorResult, error) {
    type SinkOssConfigHelper struct {
        SinkOssConfig
        TimeRange string `json:"TimeRange"`
    }
    ct := &struct {
        GetConnectorResult
        Config SinkOssConfigHelper `json:"Config"`
    }{}

    if err := json.Unmarshal(data, ct); err != nil {
        return nil, err
    }

    soConf := ct.Config.SinkOssConfig
    t, err := strconv.Atoi(ct.Config.TimeRange)
    if err != nil {
        return nil, err
    }
    soConf.TimeRange = t

    ret := &ct.GetConnectorResult
    ret.Config = soConf
    ret.CommonResponseResult = *commonResp
    return ret, nil
}

func marshalUpdateOssConnector(ucr *UpdateConnectorRequest) ([]byte, error) {
    soConf, ok := ucr.Config.(SinkOssConfig)
    if !ok {
        return nil, NewInvalidParameterErrorWithMessage(fmt.Sprintf("config type error,your input config type is %s,should be SinkOssConfig", reflect.TypeOf(ucr.Config)))
    }

    // set default value
    if soConf.TimestampUnit == "" {
        soConf.TimestampUnit = ConnectorMicrosecond
    }

    ct := &struct {
        Action        string            `json:"Action"`
        ColumnFields  []string          `json:"ColumnFields,omitempty"`
        ColumnNameMap map[string]string `json:"ColumnNameMap,omitempty"`
        Config        SinkOssConfig     `json:"Config,omitempty"`
    }{
        Action:        "create",
        ColumnFields:  ucr.ColumnFields,
        ColumnNameMap: ucr.ColumnNameMap,
        Config:        soConf,
    }
    return json.Marshal(ct)
}

/*  mysql Config */
type SinkMysqlConfig struct {
    ConnectorConfig
    Host     string     `json:"Host"`
    Port     string     `json:"Port"`
    Database string     `json:"Database"`
    Table    string     `json:"Table"`
    User     string     `json:"User"`
    Password string     `json:"Password"`
    Ignore   InsertMode `json:"Ignore"`
}

type InsertMode string

const (
    IGNORE    InsertMode = "true"
    OVERWRITE InsertMode = "false"
)

func marshalCreateMysqlConnector(ccr *CreateConnectorRequest) ([]byte, error) {
    soConf, ok := ccr.Config.(SinkMysqlConfig)
    if !ok {
        return nil, NewInvalidParameterErrorWithMessage(fmt.Sprintf("config type error,your input config type is %s,should be SinkMysqlConfig", reflect.TypeOf(ccr.Config)))
    }

    // set default value
    if soConf.TimestampUnit == "" {
        soConf.TimestampUnit = ConnectorMicrosecond
    }

    // set default value
    if soConf.TimestampUnit == "" {
        soConf.TimestampUnit = ConnectorMicrosecond
    }

    ct := &struct {
        Action        string            `json:"Action"`
        Type          ConnectorType     `json:"Type"`
        SinkStartTime int64             `json:"SinkStartTime"`
        ColumnFields  []string          `json:"ColumnFields"`
        ColumnNameMap map[string]string `json:"ColumnNameMap,omitempty"`
        Config        SinkMysqlConfig   `json:"Config"`
    }{
        Action:        "create",
        Type:          ccr.Type,
        SinkStartTime: ccr.SinkStartTime,
        ColumnFields:  ccr.ColumnFields,
        ColumnNameMap: ccr.ColumnNameMap,
        Config:        soConf,
    }
    return json.Marshal(ct)
}

func unmarshalGetMysqlConnector(commonResp *CommonResponseResult, data []byte) (*GetConnectorResult, error) {
    ct := &struct {
        GetConnectorResult
        Config SinkMysqlConfig `json:"Config"`
    }{}

    if err := json.Unmarshal(data, ct); err != nil {
        return nil, err
    }

    ret := &ct.GetConnectorResult
    ret.Config = ct.Config
    ret.CommonResponseResult = *commonResp
    return ret, nil
}

func marshalUpdateMysqlConnector(ucr *UpdateConnectorRequest) ([]byte, error) {
    soConf, ok := ucr.Config.(SinkMysqlConfig)
    if !ok {
        return nil, NewInvalidParameterErrorWithMessage(fmt.Sprintf("config type error,your input config type is %s,should be SinkMysqlConfig", reflect.TypeOf(ucr.Config)))
    }

    // set default value
    if soConf.TimestampUnit == "" {
        soConf.TimestampUnit = ConnectorMicrosecond
    }

    ct := &struct {
        Action        string            `json:"Action"`
        ColumnFields  []string          `json:"ColumnFields,omitempty"`
        ColumnNameMap map[string]string `json:"ColumnNameMap,omitempty"`
        Config        SinkMysqlConfig   `json:"Config,omitempty"`
    }{
        Action:        "create",
        ColumnFields:  ucr.ColumnFields,
        ColumnNameMap: ucr.ColumnNameMap,
        Config:        soConf,
    }
    return json.Marshal(ct)
}

/*  Ads Config */
type SinkAdsConfig struct {
    SinkMysqlConfig
}

func marshalCreateAdsConnector(ccr *CreateConnectorRequest) ([]byte, error) {
    soConf, ok := ccr.Config.(SinkAdsConfig)
    if !ok {
        return nil, NewInvalidParameterErrorWithMessage(fmt.Sprintf("config type error,your input config type is %s,should be SinkAdsConfig", reflect.TypeOf(ccr.Config)))
    }

    // set default value
    if soConf.TimestampUnit == "" {
        soConf.TimestampUnit = ConnectorMicrosecond
    }

    ct := &struct {
        Action        string            `json:"Action"`
        Type          ConnectorType     `json:"Type"`
        SinkStartTime int64             `json:"SinkStartTime"`
        ColumnFields  []string          `json:"ColumnFields"`
        ColumnNameMap map[string]string `json:"ColumnNameMap,omitempty"`
        Config        SinkAdsConfig     `json:"Config"`
    }{
        Action:        "create",
        Type:          ccr.Type,
        SinkStartTime: ccr.SinkStartTime,
        ColumnFields:  ccr.ColumnFields,
        ColumnNameMap: ccr.ColumnNameMap,
        Config:        soConf,
    }
    return json.Marshal(ct)
}

func unmarshalGetAdsConnector(commonResp *CommonResponseResult, data []byte) (*GetConnectorResult, error) {
    ct := &struct {
        GetConnectorResult
        Config SinkMysqlConfig `json:"Config"`
    }{}

    if err := json.Unmarshal(data, ct); err != nil {
        return nil, err
    }

    ret := &ct.GetConnectorResult
    ret.Config = ct.Config
    ret.CommonResponseResult = *commonResp
    return ret, nil
}

func marshalUpdateAdsConnector(ucr *UpdateConnectorRequest) ([]byte, error) {
    soConf, ok := ucr.Config.(SinkAdsConfig)
    if !ok {
        return nil, NewInvalidParameterErrorWithMessage(fmt.Sprintf("config type error,your input config type is %s,should be SinkAdsConfig", reflect.TypeOf(ucr.Config)))
    }

    // set default value
    if soConf.TimestampUnit == "" {
        soConf.TimestampUnit = ConnectorMicrosecond
    }

    ct := &struct {
        Action        string            `json:"Action"`
        ColumnFields  []string          `json:"ColumnFields,omitempty"`
        ColumnNameMap map[string]string `json:"ColumnNameMap,omitempty"`
        Config        SinkAdsConfig     `json:"Config,omitempty"`
    }{
        Action:        "create",
        ColumnFields:  ucr.ColumnFields,
        ColumnNameMap: ucr.ColumnNameMap,
        Config:        soConf,
    }
    return json.Marshal(ct)
}

/*  datahub Config */
type SinkDatahubConfig struct {
    ConnectorConfig
    Endpoint  string   `json:"Endpoint"`
    Project   string   `json:"Project"`
    Topic     string   `json:"Topic"`
    AuthMode  AuthMode `json:"AuthMode"`
    AccessId  string   `json:"AccessId,omitempty"`
    AccessKey string   `json:"AccessKey,omitempty"`
}

func marshalCreateDatahubConnector(ccr *CreateConnectorRequest) ([]byte, error) {
    soConf, ok := ccr.Config.(SinkDatahubConfig)
    if !ok {
        return nil, NewInvalidParameterErrorWithMessage(fmt.Sprintf("config type error,your input config type is %s,should be SinkDatahubConfig", reflect.TypeOf(ccr.Config)))
    }

    // set default value
    if soConf.TimestampUnit == "" {
        soConf.TimestampUnit = ConnectorMicrosecond
    }

    ct := &struct {
        Action        string            `json:"Action"`
        Type          ConnectorType     `json:"Type"`
        SinkStartTime int64             `json:"SinkStartTime"`
        ColumnFields  []string          `json:"ColumnFields"`
        ColumnNameMap map[string]string `json:"ColumnNameMap,omitempty"`
        Config        SinkDatahubConfig `json:"Config"`
    }{
        Action:        "create",
        Type:          ccr.Type,
        SinkStartTime: ccr.SinkStartTime,
        ColumnFields:  ccr.ColumnFields,
        ColumnNameMap: ccr.ColumnNameMap,
        Config:        soConf,
    }
    return json.Marshal(ct)
}

func unmarshalGetDatahubConnector(commonResp *CommonResponseResult, data []byte) (*GetConnectorResult, error) {
    ct := &struct {
        GetConnectorResult
        Config SinkDatahubConfig `json:"Config"`
    }{}

    if err := json.Unmarshal(data, ct); err != nil {
        return nil, err
    }

    ret := &ct.GetConnectorResult
    ret.Config = ct.Config
    ret.CommonResponseResult = *commonResp
    return ret, nil
}

func marshalUpdateDatahubConnector(ucr *UpdateConnectorRequest) ([]byte, error) {
    soConf, ok := ucr.Config.(SinkDatahubConfig)
    if !ok {
        return nil, NewInvalidParameterErrorWithMessage(fmt.Sprintf("config type error,your input config type is %s,should be SinkDatahubConfig", reflect.TypeOf(ucr.Config)))
    }

    // set default value
    if soConf.TimestampUnit == "" {
        soConf.TimestampUnit = ConnectorMicrosecond
    }

    ct := &struct {
        Action        string            `json:"Action"`
        ColumnFields  []string          `json:"ColumnFields,omitempty"`
        ColumnNameMap map[string]string `json:"ColumnNameMap,omitempty"`
        Config        SinkDatahubConfig `json:"Config,omitempty"`
    }{
        Action:        "create",
        ColumnFields:  ucr.ColumnFields,
        ColumnNameMap: ucr.ColumnNameMap,
        Config:        soConf,
    }
    return json.Marshal(ct)
}

/*  ES Config */
type SinkEsConfig struct {
    ConnectorConfig
    Index        string   `json:"Index"`
    Endpoint     string   `json:"Endpoint"`
    User         string   `json:"User"`
    Password     string   `json:"Password"`
    IDFields     []string `json:"IDFields"`
    TypeFields   []string `json:"TypeFields"`
    RouterFields []string `json:"RouterFields"`
    ProxyMode    bool     `json:"ProxyMode"`
}

func marshalCreateEsConnector(ccr *CreateConnectorRequest) ([]byte, error) {
    soConf, ok := ccr.Config.(SinkEsConfig)
    if !ok {
        return nil, NewInvalidParameterErrorWithMessage(fmt.Sprintf("config type error,your input config type is %s,should be SinkEsConfig", reflect.TypeOf(ccr.Config)))
    }

    // set default value
    if soConf.TimestampUnit == "" {
        soConf.TimestampUnit = ConnectorMicrosecond
    }

    // server need ProxyMode be string
    type SinkEsConfigHelper struct {
        SinkEsConfig
        ProxyMode string `json:"ProxyMode"`
    }
    confHelper := SinkEsConfigHelper{
        SinkEsConfig: soConf,
        ProxyMode:    strconv.FormatBool(soConf.ProxyMode),
    }

    ct := &struct {
        Action        string             `json:"Action"`
        Type          ConnectorType      `json:"Type"`
        SinkStartTime int64              `json:"SinkStartTime"`
        ColumnFields  []string           `json:"ColumnFields"`
        ColumnNameMap map[string]string  `json:"ColumnNameMap,omitempty"`
        Config        SinkEsConfigHelper `json:"Config"`
    }{
        Action:        "create",
        Type:          ccr.Type,
        SinkStartTime: ccr.SinkStartTime,
        ColumnFields:  ccr.ColumnFields,
        ColumnNameMap: ccr.ColumnNameMap,
        Config:        confHelper,
    }
    return json.Marshal(ct)
}

func unmarshalGetEsConnector(commonResp *CommonResponseResult, data []byte) (*GetConnectorResult, error) {
    type SinkEsConfigHelper struct {
        SinkEsConfig
        IDFields     string `json:"IDFields"`
        TypeFields   string `json:"TypeFields"`
        RouterFields string `json:"RouterFields"`
        ProxyMode    string `json:"ProxyMode"`
    }

    ct := &struct {
        GetConnectorResult
        Config SinkEsConfigHelper `json:"Config"`
    }{}

    if err := json.Unmarshal(data, ct); err != nil {
        return nil, err
    }

    conf := ct.Config.SinkEsConfig
    if ct.Config.ProxyMode != "" {
        proxy, err := strconv.ParseBool(ct.Config.ProxyMode)
        if err != nil {
            return nil, err
        }
        conf.ProxyMode = proxy
    }

    idFields := make([]string, 0)
    if ct.Config.IDFields != "" {
        if err := json.Unmarshal([]byte(ct.Config.IDFields), &idFields); err != nil {
            return nil, err
        }
    }
    conf.IDFields = idFields

    typeFields := make([]string, 0)
    if ct.Config.TypeFields != "" {
        if err := json.Unmarshal([]byte(ct.Config.TypeFields), &typeFields); err != nil {
            return nil, err
        }
        conf.TypeFields = typeFields
    }
    conf.TypeFields = typeFields

    routerFields := make([]string, 0)
    if ct.Config.RouterFields != "" {
        if err := json.Unmarshal([]byte(ct.Config.RouterFields), &routerFields); err != nil {
            return nil, err
        }
    }
    conf.RouterFields = routerFields

    ret := &ct.GetConnectorResult
    ret.CommonResponseResult = *commonResp
    ret.Config = conf
    return ret, nil
}

func marshalUpdateEsConnector(ucr *UpdateConnectorRequest) ([]byte, error) {
    soConf, ok := ucr.Config.(SinkEsConfig)
    if !ok {
        return nil, NewInvalidParameterErrorWithMessage(fmt.Sprintf("config type error,your input config type is %s,should be SinkEsConfig", reflect.TypeOf(ucr.Config)))
    }

    // set default value
    if soConf.TimestampUnit == "" {
        soConf.TimestampUnit = ConnectorMicrosecond
    }

    ct := &struct {
        Action        string            `json:"Action"`
        ColumnFields  []string          `json:"ColumnFields,omitempty"`
        ColumnNameMap map[string]string `json:"ColumnNameMap,omitempty"`
        Config        SinkEsConfig      `json:"Config,omitempty"`
    }{
        Action:        "create",
        ColumnFields:  ucr.ColumnFields,
        ColumnNameMap: ucr.ColumnNameMap,
        Config:        soConf,
    }
    return json.Marshal(ct)
}

type FcInvokeType string

const (
    FcSync  FcInvokeType = "sync"
    FcAsync FcInvokeType = "async"
)

/*  FC Config */
type SinkFcConfig struct {
    ConnectorConfig
    Endpoint   string       `json:"Endpoint"`
    Service    string       `json:"Service"`
    Function   string       `json:"Function"`
    AuthMode   AuthMode     `json:"AuthMode"`
    AccessId   string       `json:"AccessId,omitempty"`
    AccessKey  string       `json:"AccessKey,omitempty"`
    InvokeType FcInvokeType `json:"InvokeType"`
}

func marshalCreateFcConnector(ccr *CreateConnectorRequest) ([]byte, error) {
    soConf, ok := ccr.Config.(SinkFcConfig)
    if !ok {
        return nil, NewInvalidParameterErrorWithMessage(fmt.Sprintf("config type error,your input config type is %s,should be SinkFcConfig", reflect.TypeOf(ccr.Config)))
    }

    // set default value
    if soConf.TimestampUnit == "" {
        soConf.TimestampUnit = ConnectorMicrosecond
    }
    if soConf.InvokeType == "" {
        soConf.InvokeType = FcSync
    }

    ct := &struct {
        Action        string            `json:"Action"`
        Type          ConnectorType     `json:"Type"`
        SinkStartTime int64             `json:"SinkStartTime"`
        ColumnFields  []string          `json:"ColumnFields"`
        ColumnNameMap map[string]string `json:"ColumnNameMap,omitempty"`
        Config        SinkFcConfig      `json:"Config"`
    }{
        Action:        "create",
        Type:          ccr.Type,
        SinkStartTime: ccr.SinkStartTime,
        ColumnFields:  ccr.ColumnFields,
        ColumnNameMap: ccr.ColumnNameMap,
        Config:        soConf,
    }
    return json.Marshal(ct)
}

func unmarshalGetFcConnector(commonResp *CommonResponseResult, data []byte) (*GetConnectorResult, error) {
    ct := &struct {
        GetConnectorResult
        Config SinkFcConfig `json:"Config"`
    }{}

    if err := json.Unmarshal(data, ct); err != nil {
        return nil, err
    }

    ret := &ct.GetConnectorResult
    ret.Config = ct.Config
    ret.CommonResponseResult = *commonResp
    return ret, nil
}

func marshalUpdateFcConnector(ucr *UpdateConnectorRequest) ([]byte, error) {
    soConf, ok := ucr.Config.(SinkFcConfig)
    if !ok {
        return nil, NewInvalidParameterErrorWithMessage(fmt.Sprintf("config type error,your input config type is %s,should be SinkFcConfig", reflect.TypeOf(ucr.Config)))
    }

    // set default value
    if soConf.TimestampUnit == "" {
        soConf.TimestampUnit = ConnectorMicrosecond
    }

    ct := &struct {
        Action        string            `json:"Action"`
        ColumnFields  []string          `json:"ColumnFields,omitempty"`
        ColumnNameMap map[string]string `json:"ColumnNameMap,omitempty"`
        Config        SinkFcConfig      `json:"Config,omitempty"`
    }{
        Action:        "create",
        ColumnFields:  ucr.ColumnFields,
        ColumnNameMap: ucr.ColumnNameMap,
        Config:        soConf,
    }
    return json.Marshal(ct)
}

type OtsWriteMode string

const (
    OtsPut    OtsWriteMode = "PUT"
    OtsUpdate OtsWriteMode = "UPDATE"
)

/*  Ots Config */
type SinkOtsConfig struct {
    ConnectorConfig
    Endpoint     string       `json:"Endpoint"`
    InstanceName string       `json:"InstanceName"`
    TableName    string       `json:"TableName"`
    AuthMode     AuthMode     `json:"AuthMode"`
    AccessId     string       `json:"AccessId,omitempty"`
    AccessKey    string       `json:"AccessKey,omitempty"`
    WriteMode    OtsWriteMode `json:"WriteMode"`
}

func marshalCreateOtsConnector(ccr *CreateConnectorRequest) ([]byte, error) {
    soConf, ok := ccr.Config.(SinkOtsConfig)
    if !ok {
        return nil, NewInvalidParameterErrorWithMessage(fmt.Sprintf("config type error,your input config type is %s,should be SinkOtsConfig", reflect.TypeOf(ccr.Config)))
    }

    // set default value
    if soConf.TimestampUnit == "" {
        soConf.TimestampUnit = ConnectorMicrosecond
    }
    if soConf.WriteMode == "" {
        soConf.WriteMode = OtsPut
    }

    ct := &struct {
        Action        string            `json:"Action"`
        Type          ConnectorType     `json:"Type"`
        SinkStartTime int64             `json:"SinkStartTime"`
        ColumnFields  []string          `json:"ColumnFields"`
        ColumnNameMap map[string]string `json:"ColumnNameMap,omitempty"`
        Config        SinkOtsConfig     `json:"Config"`
    }{
        Action:        "create",
        Type:          ccr.Type,
        SinkStartTime: ccr.SinkStartTime,
        ColumnFields:  ccr.ColumnFields,
        ColumnNameMap: ccr.ColumnNameMap,
        Config:        soConf,
    }
    return json.Marshal(ct)
}

func unmarshalGetOtsConnector(commonResp *CommonResponseResult, data []byte) (*GetConnectorResult, error) {
    ct := &struct {
        GetConnectorResult
        Config SinkOtsConfig `json:"Config"`
    }{}

    if err := json.Unmarshal(data, ct); err != nil {
        return nil, err
    }

    ret := &ct.GetConnectorResult
    ret.Config = ct.Config
    ret.CommonResponseResult = *commonResp
    return ret, nil
}

func marshalUpdateOtsConnector(ucr *UpdateConnectorRequest) ([]byte, error) {
    soConf, ok := ucr.Config.(SinkOtsConfig)
    if !ok {
        return nil, NewInvalidParameterErrorWithMessage(fmt.Sprintf("config type error,your input config type is %s,should be SinkMysqlConfig", reflect.TypeOf(ucr.Config)))
    }

    // set default value
    if soConf.TimestampUnit == "" {
        soConf.TimestampUnit = ConnectorMicrosecond
    }

    ct := &struct {
        Action        string            `json:"Action"`
        ColumnFields  []string          `json:"ColumnFields,omitempty"`
        ColumnNameMap map[string]string `json:"ColumnNameMap,omitempty"`
        Config        SinkOtsConfig     `json:"Config,omitempty"`
    }{
        Action:        "create",
        ColumnFields:  ucr.ColumnFields,
        ColumnNameMap: ucr.ColumnNameMap,
        Config:        soConf,
    }
    return json.Marshal(ct)
}

/*  datahub Config */
type SinkHologresConfig struct {
    SinkDatahubConfig
    InstanceId string `json:"InstanceId,omitempty"`
}

func marshalCreateHologresConnector(ccr *CreateConnectorRequest) ([]byte, error) {
    soConf, ok := ccr.Config.(SinkHologresConfig)
    if !ok {
        return nil, NewInvalidParameterErrorWithMessage(fmt.Sprintf("config type error,your input config type is %s,should be SinkHologresConfig", reflect.TypeOf(ccr.Config)))
    }

    // set default value
    if soConf.TimestampUnit == "" {
        soConf.TimestampUnit = ConnectorMicrosecond
    }

    ct := &struct {
        Action        string             `json:"Action"`
        Type          ConnectorType      `json:"Type"`
        SinkStartTime int64              `json:"SinkStartTime"`
        ColumnFields  []string           `json:"ColumnFields"`
        ColumnNameMap map[string]string  `json:"ColumnNameMap,omitempty"`
        Config        SinkHologresConfig `json:"Config"`
    }{
        Action:        "create",
        Type:          ccr.Type,
        SinkStartTime: ccr.SinkStartTime,
        ColumnFields:  ccr.ColumnFields,
        ColumnNameMap: ccr.ColumnNameMap,
        Config:        soConf,
    }
    return json.Marshal(ct)
}

func unmarshalGetHologresConnector(commonResp *CommonResponseResult, data []byte) (*GetConnectorResult, error) {
    ct := &struct {
        GetConnectorResult
        Config SinkHologresConfig `json:"Config"`
    }{}

    if err := json.Unmarshal(data, ct); err != nil {
        return nil, err
    }

    ret := &ct.GetConnectorResult
    ret.Config = ct.Config
    ret.CommonResponseResult = *commonResp
    return ret, nil
}

func marshalUpdateHologresConnector(ucr *UpdateConnectorRequest) ([]byte, error) {
    soConf, ok := ucr.Config.(SinkHologresConfig)
    if !ok {
        return nil, NewInvalidParameterErrorWithMessage(fmt.Sprintf("config type error,your input config type is %s,should be SinkHologresConfig", reflect.TypeOf(ucr.Config)))
    }

    // set default value
    if soConf.TimestampUnit == "" {
        soConf.TimestampUnit = ConnectorMicrosecond
    }

    ct := &struct {
        Action        string             `json:"Action"`
        ColumnFields  []string           `json:"ColumnFields,omitempty"`
        ColumnNameMap map[string]string  `json:"ColumnNameMap,omitempty"`
        Config        SinkHologresConfig `json:"Config,omitempty"`
    }{
        Action:        "create",
        ColumnFields:  ucr.ColumnFields,
        ColumnNameMap: ucr.ColumnNameMap,
        Config:        soConf,
    }
    return json.Marshal(ct)
}

type ConnectorOffset struct {
    Timestamp int64 `json:"Timestamp"`
    Sequence  int64 `json:"Sequence"`
}

type ConnectorShardState string

// Deprecated, will be removed in a future version
const (
    Created   ConnectorShardState = "CONTEXT_PLANNED"
    Eexcuting ConnectorShardState = "CONTEXT_EXECUTING"
    Stopped   ConnectorShardState = "CONTEXT_PAUSED"
    Finished  ConnectorShardState = "CONTEXT_FINISHED"
)

const (
    ConnectorShardHang      ConnectorShardState = "CONTEXT_HANG"
    ConnectorShardPlanned   ConnectorShardState = "CONTEXT_PLANNED"
    ConnectorShardExecuting ConnectorShardState = "CONTEXT_EXECUTING"
    ConnectorShardStopped   ConnectorShardState = "CONTEXT_STOPPED"
    ConnectorShardFinished  ConnectorShardState = "CONTEXT_FINISHED"
)

type ConnectorShardStatusEntry struct {
    StartSequence    int64               `json:"StartSequence"`
    EndSequence      int64               `json:"EndSequence"`
    CurrentSequence  int64               `json:"CurrentSequence"`
    CurrentTimestamp int64               `json:"CurrentTimestamp"`
    UpdateTime       int64               `json:"UpdateTime"`
    State            ConnectorShardState `json:"State"`
    LastErrorMessage string              `json:"LastErrorMessage"`
    DiscardCount     int64               `json:"DiscardCount"`
    DoneTime         int64               `json:"DoneTime"`
    WorkerAddress    string              `json:"WorkerAddress"`
}
