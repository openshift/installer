package datahub

import (
    "encoding/json"
    "errors"
    "fmt"
    "github.com/golang/protobuf/proto"
    "github.com/aliyun/aliyun-datahub-sdk-go/datahub/pbmodel"
    "github.com/aliyun/aliyun-datahub-sdk-go/datahub/util"
    "reflect"
)

// handel the http request
type RequestModel interface {
    // serialize the requestModel and maybe need add some message on http header
    requestBodyEncode() ([]byte, error)
}

// empty request
type EmptyRequest struct {
}

func (br *EmptyRequest) requestBodyEncode() ([]byte, error) {
    return nil, nil
}

type CreateProjectRequest struct {
    Comment string `json:"Comment"`
}

func (cpr *CreateProjectRequest) requestBodyEncode() ([]byte, error) {
    return json.Marshal(cpr)
}

type UpdateProjectRequest struct {
    Comment string `json:"Comment"`
}

func (upr *UpdateProjectRequest) requestBodyEncode() ([]byte, error) {
    return json.Marshal(upr)
}

type UpdateProjectVpcWhitelistRequest struct {
    VpcIds string `json:"VpcIds"`
}

func (upv *UpdateProjectVpcWhitelistRequest) requestBodyEncode() ([]byte, error) {
    return json.Marshal(upv)
}

type CreateTopicRequest struct {
    Action       string        `json:"Action"`
    ShardCount   int           `json:"ShardCount"`
    Lifecycle    int           `json:"Lifecycle"`
    RecordType   RecordType    `json:"RecordType"`
    RecordSchema *RecordSchema `json:"RecordSchema,omitempty"`
    Comment      string        `json:"Comment"`
    ExpandMode   ExpandMode    `json:"ExpandMode"`
}

func (ctr *CreateTopicRequest) MarshalJSON() ([]byte, error) {
    msg := &struct {
        Action       string     `json:"Action"`
        ShardCount   int        `json:"ShardCount"`
        Lifecycle    int        `json:"Lifecycle"`
        RecordType   RecordType `json:"RecordType"`
        RecordSchema string     `json:"RecordSchema,omitempty"`
        Comment      string     `json:"Comment"`
        ExpandMode   ExpandMode `json:"ExpandMode"`
    }{
        Action:     ctr.Action,
        ShardCount: ctr.ShardCount,
        Lifecycle:  ctr.Lifecycle,
        RecordType: ctr.RecordType,
        Comment:    ctr.Comment,
        ExpandMode: ctr.ExpandMode,
    }
    switch ctr.RecordType {
    case TUPLE:
        msg.RecordSchema = ctr.RecordSchema.String()
    default:
        msg.RecordSchema = ""

    }
    return json.Marshal(msg)
}

func (ctr *CreateTopicRequest) requestBodyEncode() ([]byte, error) {
    return json.Marshal(ctr)
}

type UpdateTopicRequest struct {
    Comment   string `json:"Comment,omitempty"`
    Lifecycle int    `json:"Lifecycle,omitempty"`
}

func (utr *UpdateTopicRequest) requestBodyEncode() ([]byte, error) {
    return json.Marshal(utr)
}

type SplitShardRequest struct {
    Action   string `json:"Action"`
    ShardId  string `json:"ShardId"`
    SplitKey string `json:"SplitKey,omitempty"`
}

func (ssr *SplitShardRequest) requestBodyEncode() ([]byte, error) {
    return json.Marshal(ssr)
}

type MergeShardRequest struct {
    Action          string `json:"Action"`
    ShardId         string `json:"ShardId"`
    AdjacentShardId string `json:"AdjacentShardId"`
}

func (msr *MergeShardRequest) requestBodyEncode() ([]byte, error) {
    return json.Marshal(msr)
}

type ExtendShardRequest struct {
    Action     string `json:"Action"`
    ExtendMode string `json:"ExtendMode"`
    ShardCount int    `json:"ShardNumber"`
}

func (esr *ExtendShardRequest) requestBodyEncode() ([]byte, error) {
    return json.Marshal(esr)
}

type GetCursorRequest struct {
    Action     string     `json:"Action"`
    CursorType CursorType `json:"Type"`
    SystemTime int64      `json:"SystemTime"`
    Sequence   int64      `json:"Sequence"`
}

func (gcr *GetCursorRequest) requestBodyEncode() ([]byte, error) {

    type ReqMsg struct {
        Action string     `json:"Action"`
        Type   CursorType `json:"Type"`
    }
    reqMsg := ReqMsg{
        Action: gcr.Action,
        Type:   gcr.CursorType,
    }
    switch gcr.CursorType {
    case OLDEST, LATEST:
        return json.Marshal(reqMsg)
    case SYSTEM_TIME:
        return json.Marshal(struct {
            ReqMsg
            SystemTime int64 `json:"SystemTime"`
        }{
            ReqMsg:     reqMsg,
            SystemTime: gcr.SystemTime,
        })
    case SEQUENCE:
        return json.Marshal(struct {
            ReqMsg
            Sequence int64 `json:"Sequence"`
        }{
            ReqMsg:   reqMsg,
            Sequence: gcr.Sequence,
        })
    default:
        return nil, errors.New(fmt.Sprintf("Cursor not support type %s", gcr.CursorType))
    }
}

type PutRecordsRequest struct {
    Action  string    `json:"Action"`
    Records []IRecord `json:"Records"`
}

func (prr *PutRecordsRequest) requestBodyEncode() ([]byte, error) {
    return json.Marshal(prr)
}

func (ptr *PutRecordsRequest) MarshalJSON() ([]byte, error) {
    msg := &struct {
        Action  string        `json:"Action"`
        Records []RecordEntry `json:"Records"`
    }{
        Action:  ptr.Action,
        Records: make([]RecordEntry, len(ptr.Records)),
    }
    for idx, val := range ptr.Records {
        msg.Records[idx].Data = val.GetData()
        msg.Records[idx].BaseRecord = val.GetBaseRecord()
    }
    return json.Marshal(msg)
}

type GetRecordRequest struct {
    Action string `json:"Action"`
    Cursor string `json:"Cursor"`
    Limit  int    `json:"Limit"`
}

func (grr *GetRecordRequest) requestBodyEncode() ([]byte, error) {
    return json.Marshal(grr)
}

type AppendFieldRequest struct {
    Action    string    `json:"Action"`
    FieldName string    `json:"FieldName"`
    FieldType FieldType `json:"FieldType"`
}

func (afr *AppendFieldRequest) requestBodyEncode() ([]byte, error) {
    return json.Marshal(afr)
}

type GetMeterInfoRequest struct {
    Action string `json:"Action"`
}

func (gmir *GetMeterInfoRequest) requestBodyEncode() ([]byte, error) {
    return json.Marshal(gmir)
}

type CreateConnectorRequest struct {
    Action        string            `json:"Action"`
    Type          ConnectorType     `json:"Type"`
    SinkStartTime int64             `json:"SinkStartTime"`
    ColumnFields  []string          `json:"ColumnFields"`
    ColumnNameMap map[string]string `json:"ColumnNameMap"`
    Config        interface{}       `json:"Config"`
}

func (ccr *CreateConnectorRequest) requestBodyEncode() ([]byte, error) {
    switch ccr.Type {
    case SinkOdps:
        return marshalCreateOdpsConnector(ccr)
    case SinkOss:
        return marshalCreateOssConnector(ccr)
    case SinkEs:
        return marshalCreateEsConnector(ccr)
    case SinkAds:
        return marshalCreateAdsConnector(ccr)
    case SinkMysql:
        return marshalCreateMysqlConnector(ccr)
    case SinkFc:
        return marshalCreateFcConnector(ccr)
    case SinkOts:
        return marshalCreateOtsConnector(ccr)
    case SinkDatahub:
        return marshalCreateDatahubConnector(ccr)
    case SinkHologres:
        return marshalCreateHologresConnector(ccr)
    default:
        return nil, errors.New(fmt.Sprintf("not support connector type config: %s", ccr.Type.String()))
    }
}

type UpdateConnectorRequest struct {
    Action        string            `json:"Action"`
    ColumnFields  []string          `json:"ColumnFields"`
    ColumnNameMap map[string]string `json:"ColumnNameMap"`
    Config        interface{}       `json:"Config"`
}

func (ucr *UpdateConnectorRequest) requestBodyEncode() ([]byte, error) {

    if ucr.Config == nil {
        return marshalUpdateConnector(ucr)
    }

    switch ucr.Config.(type) {
    case SinkOdpsConfig:
        return marshalUpdateOdpsConnector(ucr)
    case SinkOssConfig:
        return marshalUpdateOssConnector(ucr)
    case SinkEsConfig:
        return marshalUpdateEsConnector(ucr)
    case SinkAdsConfig:
        return marshalUpdateAdsConnector(ucr)
    case SinkMysqlConfig:
        return marshalUpdateMysqlConnector(ucr)
    case SinkFcConfig:
        return marshalUpdateFcConnector(ucr)
    case SinkOtsConfig:
        return marshalUpdateOtsConnector(ucr)
    case SinkDatahubConfig:
        return marshalUpdateDatahubConnector(ucr)
    case SinkHologresConfig:
        return marshalUpdateHologresConnector(ucr)
    default:
        return nil, errors.New(fmt.Sprintf("this connector type not support, %t", reflect.TypeOf(ucr.Config)))
    }
}

type ReloadConnectorRequest struct {
    Action  string `json:"Action"`
    ShardId string `json:"ShardId,omitempty"`
}

func (rcr *ReloadConnectorRequest) requestBodyEncode() ([]byte, error) {
    return json.Marshal(rcr)
}

type UpdateConnectorStateRequest struct {
    Action string         `json:"Action"`
    State  ConnectorState `json:"State"`
}

func (ucsr *UpdateConnectorStateRequest) requestBodyEncode() ([]byte, error) {
    return json.Marshal(ucsr)
}

type UpdateConnectorOffsetRequest struct {
    Action    string `json:"Action"`
    ShardId   string `json:"ShardId"`
    Timestamp int64  `json:"CurrentTime"`
    Sequence  int64  `json:"CurrentSequence"`
}

func (ucor *UpdateConnectorOffsetRequest) requestBodyEncode() ([]byte, error) {
    return json.Marshal(ucor)
}

type GetConnectorShardStatusRequest struct {
    Action  string `json:"Action"`
    ShardId string `json:"ShardId,omitempty"`
}

func (gcss *GetConnectorShardStatusRequest) requestBodyEncode() ([]byte, error) {
    return json.Marshal(gcss)
}

type AppendConnectorFieldRequest struct {
    Action    string `json:"Action"`
    FieldName string `json:"FieldName"`
}

func (acfr *AppendConnectorFieldRequest) requestBodyEncode() ([]byte, error) {
    return json.Marshal(acfr)
}

type CreateSubscriptionRequest struct {
    Action  string `json:"Action"`
    Comment string `json:"Comment"`
}

func (csr *CreateSubscriptionRequest) requestBodyEncode() ([]byte, error) {
    return json.Marshal(csr)
}

type ListSubscriptionRequest struct {
    Action    string `json:"Action"`
    PageIndex int    `json:"PageIndex"`
    PageSize  int    `json:"PageSize"`
}

func (lsr *ListSubscriptionRequest) requestBodyEncode() ([]byte, error) {
    return json.Marshal(lsr)
}

type UpdateSubscriptionRequest struct {
    //Action  string            `json:"Action"`
    Comment string `json:"Comment"`
}

func (usr *UpdateSubscriptionRequest) requestBodyEncode() ([]byte, error) {
    return json.Marshal(usr)
}

type UpdateSubscriptionStateRequest struct {
    State SubscriptionState `json:"State"`
}

func (ussr *UpdateSubscriptionStateRequest) requestBodyEncode() ([]byte, error) {
    return json.Marshal(ussr)
}

type OpenSubscriptionSessionRequest struct {
    Action   string   `json:"Action"`
    ShardIds []string `json:"ShardIds"`
}

func (ossr *OpenSubscriptionSessionRequest) requestBodyEncode() ([]byte, error) {
    return json.Marshal(ossr)
}

type GetSubscriptionOffsetRequest struct {
    Action   string   `json:"Action"`
    ShardIds []string `json:"ShardIds"`
}

func (gsor *GetSubscriptionOffsetRequest) requestBodyEncode() ([]byte, error) {
    return json.Marshal(gsor)
}

type CommitSubscriptionOffsetRequest struct {
    Action  string                        `json:"Action"`
    Offsets map[string]SubscriptionOffset `json:"Offsets"`
}

func (csor *CommitSubscriptionOffsetRequest) requestBodyEncode() ([]byte, error) {
    return json.Marshal(csor)
}

type ResetSubscriptionOffsetRequest struct {
    Action  string                        `json:"Action"`
    Offsets map[string]SubscriptionOffset `json:"Offsets"`
}

func (rsor *ResetSubscriptionOffsetRequest) requestBodyEncode() ([]byte, error) {
    return json.Marshal(rsor)
}

type HeartbeatRequest struct {
    Action           string    `json:"Action"`
    ConsumerId       string    `json:"ConsumerId"`
    VersionId        int64     `json:"VersionId"`
    HoldShardList    []string  `json:"HoldShardList,omitempty"`
    ReadEndShardList [] string `json:"ReadEndShardList,omitempty"`
}

func (hr *HeartbeatRequest) requestBodyEncode() ([]byte, error) {
    return json.Marshal(hr)
}

type JoinGroupRequest struct {
    Action         String `json:"Action"`
    SessionTimeout int64  `json:"SessionTimeout"`
}

func (jgr *JoinGroupRequest) requestBodyEncode() ([]byte, error) {
    return json.Marshal(jgr)
}

type SyncGroupRequest struct {
    Action           string   `json:"Action"`
    ConsumerId       string   `json:"ConsumerId"`
    VersionId        int64    `json:"VersionId"`
    ReleaseShardList []string `json:"ReleaseShardList,omitempty"`
    ReadEndShardList []string `json:"ReadEndShardList,omitempty"`
}

func (sgr *SyncGroupRequest) requestBodyEncode() ([]byte, error) {
    return json.Marshal(sgr)
}

type LeaveGroupRequest struct {
    Action     string `json:"Action"`
    ConsumerId string `json:"ConsumerId"`
    VersionId  int64  `json:"VersionId"`
}

func (lgr *LeaveGroupRequest) requestBodyEncode() ([]byte, error) {
    return json.Marshal(lgr)
}

type ListTopicSchemaRequest struct {
    Action string `json:"Action"`
}

func (lts *ListTopicSchemaRequest) requestBodyEncode() ([]byte, error) {
    return json.Marshal(lts)
}

type GetTopicSchemaRequest struct {
    Action       string        `json:"Action"`
    VersionId    int           `json:"VersionId"`
    RecordSchema *RecordSchema `json:"RecordSchema,omitempty"`
}

func (gts *GetTopicSchemaRequest) requestBodyEncode() ([]byte, error) {
    msg := &struct {
        Action       string `json:"Action"`
        VersionId    int    `json:"VersionId"`
        RecordSchema string `json:"RecordSchema,omitempty"`
    }{
        Action:    gts.Action,
        VersionId: gts.VersionId,
    }

    if gts.RecordSchema != nil {
        msg.RecordSchema = gts.RecordSchema.String()
    }
    return json.Marshal(msg)
}

type RegisterTopicSchemaRequest struct {
    Action       string        `json:"Action"`
    RecordSchema *RecordSchema `json:"RecordSchema"`
}

func (rts *RegisterTopicSchemaRequest) requestBodyEncode() ([]byte, error) {
    msg := &struct {
        Action       string `json:"Action"`
        RecordSchema string `json:"RecordSchema,omitempty"`
    }{
        Action: rts.Action,
    }

    if rts.RecordSchema != nil {
        msg.RecordSchema = rts.RecordSchema.String()
    }
    return json.Marshal(msg)
}

type DeleteTopicSchemaRequest struct {
    Action    string `json:"Action"`
    VersionId int    `json:"VersionId"`
}

func (lgr *DeleteTopicSchemaRequest) requestBodyEncode() ([]byte, error) {
    return json.Marshal(lgr)
}

type PutPBRecordsRequest struct {
    Records []IRecord `json:"Records"`
}

func (pr *PutPBRecordsRequest) requestBodyEncode() ([]byte, error) {
    res := make([]*pbmodel.RecordEntry, len(pr.Records))
    for idx, val := range pr.Records {
        bRecord := val.GetBaseRecord()
        data := val.GetData()

        fds := make([]*pbmodel.FieldData, 0)
        switch data.(type) {
        case []byte:
            fd := &pbmodel.FieldData{
                Value: data.([]byte),
            }
            fds = append(fds, fd)
        default:
            v, ok := data.([]interface{})
            if !ok {
                return nil, errors.New("data format is invalid")
            }
            for _, str := range v {
                fd := &pbmodel.FieldData{}
                if str == nil {
                    fd.Value = nil
                } else {
                    fd.Value = []byte(fmt.Sprintf("%s", str))
                }
                fds = append(fds, fd)
            }
        }
        rd := &pbmodel.RecordData{
            Data: fds,
        }

        recordEntry := &pbmodel.RecordEntry{
            ShardId: proto.String(bRecord.ShardId),
            Data:    rd,
        }

        if len(bRecord.Attributes) > 0 {
            sps := make([]*pbmodel.StringPair, len(bRecord.Attributes))
            index := 0
            for k, v := range bRecord.Attributes {
                strv := fmt.Sprintf("%v", v)
                sp := &pbmodel.StringPair{
                    Key:   proto.String(k),
                    Value: proto.String(strv),
                }
                sps[index] = sp
                index++
            }
            ra := &pbmodel.RecordAttributes{
                Attributes: sps,
            }
            recordEntry.Attributes = ra
        }
        res[idx] = recordEntry
    }

    prr := &pbmodel.PutRecordsRequest{
        Records: res,
    }
    buf, err := proto.Marshal(prr)
    if err != nil {
        return nil, err
    }
    x := util.WrapMessage(buf)
    return x, nil
}

type GetPBRecordRequest struct {
    Cursor string `json:"Cursor"`
    Limit  int    `json:"Limit"`
}

func (gpr *GetPBRecordRequest) requestBodyEncode() ([]byte, error) {
    limit := int32(gpr.Limit)
    grr := &pbmodel.GetRecordsRequest{
        Cursor: &gpr.Cursor,
        Limit:  &limit,
    }

    buf, err := proto.Marshal(grr)
    if err != nil {
        return nil, err
    }

    wBuf := util.WrapMessage(buf)
    return wBuf, nil
}

type PutBatchRecordsRequest struct {
    serializer     *batchSerializer
    Records        []IRecord
}

func (pbr *PutBatchRecordsRequest) requestBodyEncode() ([]byte, error) {
    batchBuf, err := pbr.serializer.serialize(pbr.Records)
    if err != nil {
        return nil, err
    }

    entry := &pbmodel.BinaryRecordEntry{
        Data: batchBuf,
    }
    protoReq := &pbmodel.PutBinaryRecordsRequest{
        Records: []*pbmodel.BinaryRecordEntry{entry},
    }

    buf, err := proto.Marshal(protoReq)
    if err != nil {
        return nil, err
    }
    return util.WrapMessage(buf), nil
}

type GetBatchRecordRequest struct {
    GetPBRecordRequest
}
