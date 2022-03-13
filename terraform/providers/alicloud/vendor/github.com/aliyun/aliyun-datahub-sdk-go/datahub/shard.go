package datahub

import (
    "errors"
    "fmt"
    "math/big"
    "strings"
)

type ShardEntry struct {
    ShardId        string     `json:"ShardId"`
    State          ShardState `json:"State"`
    BeginHashKey   string     `json:"BeginHashKey"`
    EndHashKey     string     `json:"EndHashKey"`
    ClosedTime     int64      `json:"ClosedTime"`
    ParentShardIds []string   `json:"ParentShardIds"`
    LeftShardId    string     `json:"LeftShardId"`
    RightShardId   string     `json:"RightShardId"`
    Address        string     `json:"Address"`
}

func generateSpliteKey(projectName, topicName, shardId string, datahub DataHubApi) (string, error) {
    ls, err := datahub.ListShard(projectName, topicName)
    if err != nil {
        return "", err
    }
    shards := ls.Shards
    splitKey := ""
    for _, shard := range shards {
        if strings.EqualFold(shardId, shard.ShardId) {
            if shard.State != ACTIVE {
                return "", errors.New(fmt.Sprintf("Only active shard can be split,the shard %s state is %s", shard.ShardId, shard.State))
            }
            splitKey, err = getSplitKey(shard.BeginHashKey, shard.EndHashKey)
            splitKey = strings.ToUpper(splitKey)
            if err != nil {
                return "", err
            }
        }
    }
    if splitKey == "" {
        return "", errors.New(fmt.Sprintf("Shard not exist"))
    }
    return splitKey, nil
}

func getSplitKey(beginHashKey, endHashKey string) (string, error) {
    var begin, end, sum, quo big.Int
    base := 16

    if len(beginHashKey) != 32 || len(endHashKey) != 32 {
        return "", errors.New(fmt.Sprintf("Invalid Hash Key Range", ))
    }
    _, ok := begin.SetString(beginHashKey, base)
    if !ok {
        return "", errors.New(fmt.Sprintf("Invalid Hash Key Range"))
    }
    _, ok = end.SetString(endHashKey, base)
    if !ok {
        return "", errors.New(fmt.Sprintf("Invalid Hash Key Range"))
    }

    sum.Add(&begin, &end)
    quo.Quo(&sum, big.NewInt(2))
    return quo.Text(base), nil

}
