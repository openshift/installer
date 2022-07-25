package alicloud

type MongoDBShardingNodeType string

const (
	MongoDBShardingNodeMongos = MongoDBShardingNodeType("mongos")
	MongoDBShardingNodeShard  = MongoDBShardingNodeType("shard")
)
