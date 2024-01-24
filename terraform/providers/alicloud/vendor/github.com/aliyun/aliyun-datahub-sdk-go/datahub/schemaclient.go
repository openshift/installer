package datahub

import (
    "sync"
)

type schemaRegistryClient struct {
    dh         DataHubApi
    schemaMap  sync.Map
    versionMap sync.Map
}

func NewSchemaClient(dh DataHubApi) *schemaRegistryClient {
    return &schemaRegistryClient{
        dh: dh,
    }
}

func (client *schemaRegistryClient) getSchemaByVersion(project, topic string, version int) (*RecordSchema, error) {
    key := client.genKey(project, topic)

    value, ok := client.versionMap.Load(key)
    if !ok {
        m := sync.Map{}
        client.versionMap.Store(key, m)
        value, _ = client.versionMap.Load(key)
    }

    topicCache := value.(sync.Map)
    schemaVal, ok := topicCache.Load(version)
    if !ok {
        newSchema, err := client.fetchSchemaById(project, topic, version)
        if err != nil {
            return nil, err
        }
        schemaVal = *newSchema
        topicCache.Store(version, schemaVal)
        client.versionMap.Store(key, topicCache)
    }
    schema := schemaVal.(RecordSchema)

    return &schema, nil
}

func (client *schemaRegistryClient) getVersionBySchema(project, topic string, schema *RecordSchema) (int, error) {
    key := client.genKey(project, topic)

    value, ok := client.schemaMap.Load(key)
    if !ok {
        m := sync.Map{}
        client.schemaMap.Store(key, m)
        value, _ = client.schemaMap.Load(key)
    }

    topicCache := value.(sync.Map)
    schemaKey := schema.String()
    versionVal, ok := topicCache.Load(schemaKey)
    if !ok {
        newVersion, err := client.fetchVersionBySchema(project, topic, schema)
        if err != nil {
            return -1, err
        }
        versionVal = newVersion
        topicCache.Store(schemaKey, versionVal)
        client.schemaMap.Store(key, topicCache)
    }

    version := versionVal.(int)
    return version, nil
}

func (client *schemaRegistryClient) fetchSchemaById(project, topic string, version int) (*RecordSchema, error) {
    ret, err := client.dh.GetTopicSchemaByVersion(project, topic, version)
    if err != nil {
        return nil, err
    }

    return &ret.RecordSchema, nil
}

func (client *schemaRegistryClient) fetchVersionBySchema(project, topic string, schema *RecordSchema) (int, error) {
    ret, err := client.dh.GetTopicSchemaBySchema(project, topic, schema)

    if err != nil {
        return -1, err
    }

    return ret.VersionId, nil
}

func (client *schemaRegistryClient) genKey(project, topic string) string {
    return project + "/" + topic
}
