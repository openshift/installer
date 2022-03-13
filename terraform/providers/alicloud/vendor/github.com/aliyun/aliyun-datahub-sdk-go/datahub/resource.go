package datahub

const (
    projectsPath = "/projects"
    projectPath  = "/projects/%s"
    topicsPath   = "/projects/%s/topics"
    topicPath    = "/projects/%s/topics/%s"
    shardsPath   = "/projects/%s/topics/%s/shards"
    shardPath    = "/projects/%s/topics/%s/shards/%s"

    connectorsPath    = "/projects/%s/topics/%s/connectors"
    connectorPath     = "/projects/%s/topics/%s/connectors/%s"
    consumerGroupPath = "/projects/%s/topics/%s/subscriptions/%s"

    subscriptionsPath = "/projects/%s/topics/%s/subscriptions"
    subscriptionPath  = "/projects/%s/topics/%s/subscriptions/%s"
    offsetsPath       = "/projects/%s/topics/%s/subscriptions/%s/offsets"
)
