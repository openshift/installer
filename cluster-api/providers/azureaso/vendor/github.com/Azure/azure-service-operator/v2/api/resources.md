---
title: Supported Resources
---

These are the resources with Azure Service Operator support committed to our **main** branch, grouped by the originating ARM service. (Newly supported resources will appear in this list prior to inclusion in any ASO release.)

## Authorization

| Resource       | ARM Version        | CRD Version                | Supported From | Sample                                                                                                                                           |
|----------------|--------------------|----------------------------|----------------|--------------------------------------------------------------------------------------------------------------------------------------------------|
| RoleAssignment | 2020-08-01-preview | v1beta20200801preview      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/authorization/v1beta20200801preview_roleassignment.yaml)      |
| RoleAssignment | 2020-08-01-preview | v1alpha1api20200801preview | v2.0.0-alpha.2 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/authorization/v1alpha1api20200801preview_roleassignment.yaml) |

## Batch

| Resource     | ARM Version | CRD Version         | Supported From | Sample                                                                                                                          |
|--------------|-------------|---------------------|----------------|---------------------------------------------------------------------------------------------------------------------------------|
| BatchAccount | 2021-01-01  | v1beta20210101      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/batch/v1beta20210101_batchaccount.yaml)      |
| BatchAccount | 2021-01-01  | v1alpha1api20210101 | v2.0.0-alpha.1 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/batch/v1alpha1api20210101_batchaccount.yaml) |

## Cache

| Resource                | ARM Version | CRD Version         | Supported From | Sample                                                                                                                                     |
|-------------------------|-------------|---------------------|----------------|--------------------------------------------------------------------------------------------------------------------------------------------|
| Redis                   | 2020-12-01  | v1beta20201201      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/cache/v1beta20201201_redis.yaml)                        |
| Redis                   | 2020-12-01  | v1alpha1api20201201 | v2.0.0-alpha.4 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/cache/v1alpha1api20201201_redis.yaml)                   |
| RedisEnterprise         | 2021-03-01  | v1beta20210301      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/cache/v1beta20210301_redisenterprise.yaml)              |
| RedisEnterprise         | 2021-03-01  | v1alpha1api20210301 | v2.0.0-alpha.4 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/cache/v1alpha1api20210301_redisenterprise.yaml)         |
| RedisEnterpriseDatabase | 2021-03-01  | v1beta20210301      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/cache/v1beta20210301_redisenterprisedatabase.yaml)      |
| RedisEnterpriseDatabase | 2021-03-01  | v1alpha1api20210301 | v2.0.0-alpha.4 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/cache/v1alpha1api20210301_redisenterprisedatabase.yaml) |
| RedisFirewallRule       | 2020-12-01  | v1beta20201201      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/cache/v1beta20201201_redisfirewallrule.yaml)            |
| RedisFirewallRule       | 2020-12-01  | v1alpha1api20201201 | v2.0.0-alpha.4 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/cache/v1alpha1api20201201_redisfirewallrule.yaml)       |
| RedisLinkedServer       | 2020-12-01  | v1beta20201201      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/cache/v1beta20201201_redislinkedserver.yaml)            |
| RedisLinkedServer       | 2020-12-01  | v1alpha1api20201201 | v2.0.0-alpha.4 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/cache/v1alpha1api20201201_redislinkedserver.yaml)       |
| RedisPatchSchedule      | 2020-12-01  | v1beta20201201      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/cache/v1beta20201201_redispatchschedule.yaml)           |
| RedisPatchSchedule      | 2020-12-01  | v1alpha1api20201201 | v2.0.0-alpha.4 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/cache/v1alpha1api20201201_redispatchschedule.yaml)      |

## Compute

| Resource               | ARM Version | CRD Version         | Supported From | Sample                                                                                                                                      |
|------------------------|-------------|---------------------|----------------|---------------------------------------------------------------------------------------------------------------------------------------------|
| Disk                   | 2020-09-30  | v1beta20200930      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/compute/v1beta20200930_disk.yaml)                        |
| Disk                   | 2020-09-30  | v1alpha1api20200930 | v2.0.0-alpha.1 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/compute/v1alpha1api20200930_disk.yaml)                   |
| Image                  | 2021-07-01  | v1beta20210701      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/compute/v1beta20210701_image.yaml)                       |
| Image                  | 2021-07-01  | v1alpha1api20210701 | v2.0.0-alpha.6 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/compute/v1alpha1api20210701_image.yaml)                  |
| Snapshot               | 2020-09-30  | v1beta20200930      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/compute/v1beta20200930_snapshot.yaml)                    |
| Snapshot               | 2020-09-30  | v1alpha1api20200930 | v2.0.0-alpha.4 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/compute/v1alpha1api20200930_snapshot.yaml)               |
| VirtualMachine         | 2020-12-01  | v1beta20201201      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/compute/v1beta20201201_virtualmachine.yaml)              |
| VirtualMachine         | 2020-12-01  | v1alpha1api20201201 | v2.0.0-alpha.1 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/compute/v1alpha1api20201201_virtualmachine.yaml)         |
| VirtualMachineScaleSet | 2020-12-01  | v1beta20201201      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/compute/v1beta20201201_virtualmachinescaleset.yaml)      |
| VirtualMachineScaleSet | 2020-12-01  | v1alpha1api20201201 | v2.0.0-alpha.1 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/compute/v1alpha1api20201201_virtualmachinescaleset.yaml) |

## Containerregistry

| Resource | ARM Version | CRD Version         | Supported From | Sample                                                                                                                                  |
|----------|-------------|---------------------|----------------|-----------------------------------------------------------------------------------------------------------------------------------------|
| Registry | 2021-09-01  | v1beta20210901      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/containerregistry/v1beta20210901_registry.yaml)      |
| Registry | 2021-09-01  | v1alpha1api20210901 | v2.0.0-alpha.6 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/containerregistry/v1alpha1api20210901_registry.yaml) |

## Containerservice

| Resource                 | ARM Version | CRD Version         | Supported From | Sample                                                                                                                                                 |
|--------------------------|-------------|---------------------|----------------|--------------------------------------------------------------------------------------------------------------------------------------------------------|
| ManagedCluster           | 2021-05-01  | v1beta20210501      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/containerservice/v1beta20210501_managedcluster.yaml)                |
| ManagedCluster           | 2021-05-01  | v1alpha1api20210501 | v2.0.0-alpha.1 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/containerservice/v1alpha1api20210501_managedcluster.yaml)           |
| ManagedClustersAgentPool | 2021-05-01  | v1beta20210501      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/containerservice/v1beta20210501_managedclustersagentpool.yaml)      |
| ManagedClustersAgentPool | 2021-05-01  | v1alpha1api20210501 | v2.0.0-alpha.1 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/containerservice/v1alpha1api20210501_managedclustersagentpool.yaml) |

## Dbformariadb

| Resource      | ARM Version | CRD Version    | Supported From | Sample                                                                                                                             |
|---------------|-------------|----------------|----------------|------------------------------------------------------------------------------------------------------------------------------------|
| Configuration | 2018-06-01  | v1beta20180601 | v2.0.0-beta.1  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/dbformariadb/v1beta20180601_configuration.yaml) |
| Database      | 2018-06-01  | v1beta20180601 | v2.0.0-beta.1  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/dbformariadb/v1beta20180601_database.yaml)      |
| Server        | 2018-06-01  | v1beta20180601 | v2.0.0-beta.1  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/dbformariadb/v1beta20180601_server.yaml)        |

## Dbformysql

| Resource                    | ARM Version | CRD Version         | Supported From | Sample                                                                                                                                              |
|-----------------------------|-------------|---------------------|----------------|-----------------------------------------------------------------------------------------------------------------------------------------------------|
| FlexibleServer              | 2021-05-01  | v1beta20210501      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/dbformysql/v1beta20210501_flexibleserver.yaml)                   |
| FlexibleServer              | 2021-05-01  | v1alpha1api20210501 | v2.0.0-alpha.2 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/dbformysql/v1alpha1api20210501_flexibleserver.yaml)              |
| FlexibleServersDatabase     | 2021-05-01  | v1beta20210501      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/dbformysql/v1beta20210501_flexibleserversdatabase.yaml)          |
| FlexibleServersDatabase     | 2021-05-01  | v1alpha1api20210501 | v2.0.0-alpha.2 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/dbformysql/v1alpha1api20210501_flexibleserversdatabase.yaml)     |
| FlexibleServersFirewallRule | 2021-05-01  | v1beta20210501      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/dbformysql/v1beta20210501_flexibleserversfirewallrule.yaml)      |
| FlexibleServersFirewallRule | 2021-05-01  | v1alpha1api20210501 | v2.0.0-alpha.2 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/dbformysql/v1alpha1api20210501_flexibleserversfirewallrule.yaml) |

## Dbforpostgresql

| Resource                     | ARM Version | CRD Version         | Supported From | Sample                                                                                                                                                    |
|------------------------------|-------------|---------------------|----------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------|
| FlexibleServer               | 2021-06-01  | v1beta20210601      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/dbforpostgresql/v1beta20210601_flexibleserver.yaml)                    |
| FlexibleServer               | 2021-06-01  | v1alpha1api20210601 | v2.0.0-alpha.2 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/dbforpostgresql/v1alpha1api20210601_flexibleserver.yaml)               |
| FlexibleServersConfiguration | 2021-06-01  | v1beta20210601      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/dbforpostgresql/v1beta20210601_flexibleserversconfiguration.yaml)      |
| FlexibleServersConfiguration | 2021-06-01  | v1alpha1api20210601 | v2.0.0-alpha.4 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/dbforpostgresql/v1alpha1api20210601_flexibleserversconfiguration.yaml) |
| FlexibleServersDatabase      | 2021-06-01  | v1beta20210601      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/dbforpostgresql/v1beta20210601_flexibleserversdatabase.yaml)           |
| FlexibleServersDatabase      | 2021-06-01  | v1alpha1api20210601 | v2.0.0-alpha.2 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/dbforpostgresql/v1alpha1api20210601_flexibleserversdatabase.yaml)      |
| FlexibleServersFirewallRule  | 2021-06-01  | v1beta20210601      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/dbforpostgresql/v1beta20210601_flexibleserversfirewallrule.yaml)       |
| FlexibleServersFirewallRule  | 2021-06-01  | v1alpha1api20210601 | v2.0.0-alpha.2 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/dbforpostgresql/v1alpha1api20210601_flexibleserversfirewallrule.yaml)  |

## Documentdb

| Resource                                   | ARM Version | CRD Version         | Supported From | Sample                                                                                                                                                             |
|--------------------------------------------|-------------|---------------------|----------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| DatabaseAccount                            | 2021-05-15  | v1beta20210515      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/documentdb/v1beta20210515_databaseaccount.yaml)                                 |
| DatabaseAccount                            | 2021-05-15  | v1alpha1api20210515 | v2.0.0-alpha.1 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/documentdb/v1alpha1api20210515_databaseaccount.yaml)                            |
| MongodbDatabase                            | 2021-05-15  | v1beta20210515      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/documentdb/v1beta20210515_mongodbdatabase.yaml)                                 |
| MongodbDatabase                            | 2021-05-15  | v1alpha1api20210515 | v2.0.0-alpha.2 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/documentdb/v1alpha1api20210515_mongodbdatabase.yaml)                            |
| MongodbDatabaseCollection                  | 2021-05-15  | v1beta20210515      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/documentdb/v1beta20210515_mongodbdatabasecollection.yaml)                       |
| MongodbDatabaseCollection                  | 2021-05-15  | v1alpha1api20210515 | v2.0.0-alpha.2 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/documentdb/v1alpha1api20210515_mongodbdatabasecollection.yaml)                  |
| MongodbDatabaseCollectionThroughputSetting | 2021-05-15  | v1beta20210515      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/documentdb/v1beta20210515_mongodbdatabasecollectionthroughputsetting.yaml)      |
| MongodbDatabaseCollectionThroughputSetting | 2021-05-15  | v1alpha1api20210515 | v2.0.0-alpha.2 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/documentdb/v1alpha1api20210515_mongodbdatabasecollectionthroughputsetting.yaml) |
| MongodbDatabaseThroughputSetting           | 2021-05-15  | v1beta20210515      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/documentdb/v1beta20210515_mongodbdatabasethroughputsetting.yaml)                |
| MongodbDatabaseThroughputSetting           | 2021-05-15  | v1alpha1api20210515 | v2.0.0-alpha.2 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/documentdb/v1alpha1api20210515_mongodbdatabasethroughputsetting.yaml)           |
| SqlDatabase                                | 2021-05-15  | v1beta20210515      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/documentdb/v1beta20210515_sqldatabase.yaml)                                     |
| SqlDatabase                                | 2021-05-15  | v1alpha1api20210515 | v2.0.0-alpha.2 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/documentdb/v1alpha1api20210515_sqldatabase.yaml)                                |
| SqlDatabaseContainer                       | 2021-05-15  | v1beta20210515      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/documentdb/v1beta20210515_sqldatabasecontainer.yaml)                            |
| SqlDatabaseContainer                       | 2021-05-15  | v1alpha1api20210515 | v2.0.0-alpha.2 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/documentdb/v1alpha1api20210515_sqldatabasecontainer.yaml)                       |
| SqlDatabaseContainerStoredProcedure        | 2021-05-15  | v1beta20210515      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/documentdb/v1beta20210515_sqldatabasecontainerstoredprocedure.yaml)             |
| SqlDatabaseContainerStoredProcedure        | 2021-05-15  | v1alpha1api20210515 | v2.0.0-alpha.2 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/documentdb/v1alpha1api20210515_sqldatabasecontainerstoredprocedure.yaml)        |
| SqlDatabaseContainerThroughputSetting      | 2021-05-15  | v1beta20210515      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/documentdb/v1beta20210515_sqldatabasecontainerthroughputsetting.yaml)           |
| SqlDatabaseContainerThroughputSetting      | 2021-05-15  | v1alpha1api20210515 | v2.0.0-alpha.2 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/documentdb/v1alpha1api20210515_sqldatabasecontainerthroughputsetting.yaml)      |
| SqlDatabaseContainerTrigger                | 2021-05-15  | v1beta20210515      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/documentdb/v1beta20210515_sqldatabasecontainertrigger.yaml)                     |
| SqlDatabaseContainerTrigger                | 2021-05-15  | v1alpha1api20210515 | v2.0.0-alpha.2 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/documentdb/v1alpha1api20210515_sqldatabasecontainertrigger.yaml)                |
| SqlDatabaseContainerUserDefinedFunction    | 2021-05-15  | v1beta20210515      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/documentdb/v1beta20210515_sqldatabasecontaineruserdefinedfunction.yaml)         |
| SqlDatabaseContainerUserDefinedFunction    | 2021-05-15  | v1alpha1api20210515 | v2.0.0-alpha.2 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/documentdb/v1alpha1api20210515_sqldatabasecontaineruserdefinedfunction.yaml)    |
| SqlDatabaseThroughputSetting               | 2021-05-15  | v1beta20210515      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/documentdb/v1beta20210515_sqldatabasethroughputsetting.yaml)                    |
| SqlDatabaseThroughputSetting               | 2021-05-15  | v1alpha1api20210515 | v2.0.0-alpha.2 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/documentdb/v1alpha1api20210515_sqldatabasethroughputsetting.yaml)               |

## Eventgrid

| Resource          | ARM Version | CRD Version         | Supported From | Sample                                                                                                                                   |
|-------------------|-------------|---------------------|----------------|------------------------------------------------------------------------------------------------------------------------------------------|
| Domain            | 2020-06-01  | v1beta20200601      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/eventgrid/v1beta20200601_domain.yaml)                 |
| Domain            | 2020-06-01  | v1alpha1api20200601 | v2.0.0-alpha.5 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/eventgrid/v1alpha1api20200601_domain.yaml)            |
| DomainsTopic      | 2020-06-01  | v1beta20200601      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/eventgrid/v1beta20200601_domainstopic.yaml)           |
| DomainsTopic      | 2020-06-01  | v1alpha1api20200601 | v2.0.0-alpha.5 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/eventgrid/v1alpha1api20200601_domainstopic.yaml)      |
| EventSubscription | 2020-06-01  | v1beta20200601      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/eventgrid/v1beta20200601_eventsubscription.yaml)      |
| EventSubscription | 2020-06-01  | v1alpha1api20200601 | v2.0.0-alpha.5 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/eventgrid/v1alpha1api20200601_eventsubscription.yaml) |
| Topic             | 2020-06-01  | v1beta20200601      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/eventgrid/v1beta20200601_topic.yaml)                  |
| Topic             | 2020-06-01  | v1alpha1api20200601 | v2.0.0-alpha.3 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/eventgrid/v1alpha1api20200601_topic.yaml)             |

## Eventhub

| Resource                             | ARM Version | CRD Version         | Supported From | Sample                                                                                                                                                     |
|--------------------------------------|-------------|---------------------|----------------|------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Namespace                            | 2021-11-01  | v1beta20211101      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/eventhub/v1beta20211101_namespace.yaml)                                 |
| Namespace                            | 2021-11-01  | v1alpha1api20211101 | v2.0.0-alpha.1 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/eventhub/v1alpha1api20211101_namespace.yaml)                            |
| NamespacesAuthorizationRule          | 2021-11-01  | v1beta20211101      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/eventhub/v1beta20211101_namespacesauthorizationrule.yaml)               |
| NamespacesAuthorizationRule          | 2021-11-01  | v1alpha1api20211101 | v2.0.0-alpha.2 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/eventhub/v1alpha1api20211101_namespacesauthorizationrule.yaml)          |
| NamespacesEventhub                   | 2021-11-01  | v1beta20211101      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/eventhub/v1beta20211101_namespaceseventhub.yaml)                        |
| NamespacesEventhub                   | 2021-11-01  | v1alpha1api20211101 | v2.0.0-alpha.3 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/eventhub/v1alpha1api20211101_namespaceseventhub.yaml)                   |
| NamespacesEventhubsAuthorizationRule | 2021-11-01  | v1beta20211101      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/eventhub/v1beta20211101_namespaceseventhubsauthorizationrule.yaml)      |
| NamespacesEventhubsAuthorizationRule | 2021-11-01  | v1alpha1api20211101 | v2.0.0-alpha.3 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/eventhub/v1alpha1api20211101_namespaceseventhubsauthorizationrule.yaml) |
| NamespacesEventhubsConsumerGroup     | 2021-11-01  | v1beta20211101      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/eventhub/v1beta20211101_namespaceseventhubsconsumergroup.yaml)          |
| NamespacesEventhubsConsumerGroup     | 2021-11-01  | v1alpha1api20211101 | v2.0.0-alpha.1 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/eventhub/v1alpha1api20211101_namespaceseventhubsconsumergroup.yaml)     |

## Insights

| Resource  | ARM Version        | CRD Version                | Supported From | Sample                                                                                                                               |
|-----------|--------------------|----------------------------|----------------|--------------------------------------------------------------------------------------------------------------------------------------|
| Component | 2020-02-02         | v1beta20200202             | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/insights/v1beta20200202_component.yaml)           |
| Component | 2020-02-02         | v1alpha1api20200202        | v2.0.0-alpha.2 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/insights/v1alpha1api20200202_component.yaml)      |
| Webtest   | 2018-05-01-preview | v1beta20180501preview      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/insights/v1beta20180501preview_webtest.yaml)      |
| Webtest   | 2018-05-01-preview | v1alpha1api20180501preview | v2.0.0-alpha.4 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/insights/v1alpha1api20180501preview_webtest.yaml) |

## Keyvault

| Resource | ARM Version        | CRD Version           | Supported From | Sample                                                                                                                        |
|----------|--------------------|-----------------------|----------------|-------------------------------------------------------------------------------------------------------------------------------|
| Vault    | 2021-04-01-preview | v1beta20210401preview | v2.0.0-beta.1  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/keyvault/v1beta20210401preview_vault.yaml) |

## Managedidentity

| Resource             | ARM Version | CRD Version         | Supported From | Sample                                                                                                                                            |
|----------------------|-------------|---------------------|----------------|---------------------------------------------------------------------------------------------------------------------------------------------------|
| UserAssignedIdentity | 2018-11-30  | v1beta20181130      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/managedidentity/v1beta20181130_userassignedidentity.yaml)      |
| UserAssignedIdentity | 2018-11-30  | v1alpha1api20181130 | v2.0.0-alpha.1 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/managedidentity/v1alpha1api20181130_userassignedidentity.yaml) |

## Network

| Resource                             | ARM Version | CRD Version         | Supported From | Sample                                                                                                                                                    |
|--------------------------------------|-------------|---------------------|----------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------|
| LoadBalancer                         | 2020-11-01  | v1beta20201101      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/network/v1beta20201101_loadbalancer.yaml)                              |
| LoadBalancer                         | 2020-11-01  | v1alpha1api20201101 | v2.0.0-alpha.1 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/network/v1alpha1api20201101_loadbalancer.yaml)                         |
| NetworkInterface                     | 2020-11-01  | v1beta20201101      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/network/v1beta20201101_networkinterface.yaml)                          |
| NetworkInterface                     | 2020-11-01  | v1alpha1api20201101 | v2.0.0-alpha.3 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/network/v1alpha1api20201101_networkinterface.yaml)                     |
| NetworkSecurityGroup                 | 2020-11-01  | v1beta20201101      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/network/v1beta20201101_networksecuritygroup.yaml)                      |
| NetworkSecurityGroup                 | 2020-11-01  | v1alpha1api20201101 | v2.0.0-alpha.1 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/network/v1alpha1api20201101_networksecuritygroup.yaml)                 |
| NetworkSecurityGroupsSecurityRule    | 2020-11-01  | v1beta20201101      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/network/v1beta20201101_networksecuritygroupssecurityrule.yaml)         |
| NetworkSecurityGroupsSecurityRule    | 2020-11-01  | v1alpha1api20201101 | v2.0.0-alpha.1 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/network/v1alpha1api20201101_networksecuritygroupssecurityrule.yaml)    |
| PublicIPAddress                      | 2020-11-01  | v1beta20201101      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/network/v1beta20201101_publicipaddress.yaml)                           |
| PublicIPAddress                      | 2020-11-01  | v1alpha1api20201101 | v2.0.0-alpha.1 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/network/v1alpha1api20201101_publicipaddress.yaml)                      |
| RouteTable                           | 2020-11-01  | v1beta20201101      | v2.0.0-beta.1  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/network/v1beta20201101_routetable.yaml)                                |
| RouteTablesRoute                     | 2020-11-01  | v1beta20201101      | v2.0.0-beta.1  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/network/v1beta20201101_routetablesroute.yaml)                          |
| VirtualNetwork                       | 2020-11-01  | v1beta20201101      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/network/v1beta20201101_virtualnetwork.yaml)                            |
| VirtualNetwork                       | 2020-11-01  | v1alpha1api20201101 | v2.0.0-alpha.1 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/network/v1alpha1api20201101_virtualnetwork.yaml)                       |
| VirtualNetworkGateway                | 2020-11-01  | v1beta20201101      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/network/v1beta20201101_virtualnetworkgateway.yaml)                     |
| VirtualNetworkGateway                | 2020-11-01  | v1alpha1api20201101 | v2.0.0-alpha.1 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/network/v1alpha1api20201101_virtualnetworkgateway.yaml)                |
| VirtualNetworksSubnet                | 2020-11-01  | v1beta20201101      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/network/v1beta20201101_virtualnetworkssubnet.yaml)                     |
| VirtualNetworksSubnet                | 2020-11-01  | v1alpha1api20201101 | v2.0.0-alpha.1 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/network/v1alpha1api20201101_virtualnetworkssubnet.yaml)                |
| VirtualNetworksVirtualNetworkPeering | 2020-11-01  | v1beta20201101      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/network/v1beta20201101_virtualnetworksvirtualnetworkpeering.yaml)      |
| VirtualNetworksVirtualNetworkPeering | 2020-11-01  | v1alpha1api20201101 | v2.0.0-alpha.1 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/network/v1alpha1api20201101_virtualnetworksvirtualnetworkpeering.yaml) |

## Operationalinsights

| Resource  | ARM Version | CRD Version         | Supported From | Sample                                                                                                                                     |
|-----------|-------------|---------------------|----------------|--------------------------------------------------------------------------------------------------------------------------------------------|
| Workspace | 2021-06-01  | v1beta20210601      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/operationalinsights/v1beta20210601_workspace.yaml)      |
| Workspace | 2021-06-01  | v1alpha1api20210601 | v2.0.0-alpha.4 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/operationalinsights/v1alpha1api20210601_workspace.yaml) |

## Servicebus

| Resource        | ARM Version        | CRD Version                | Supported From | Sample                                                                                                                                         |
|-----------------|--------------------|----------------------------|----------------|------------------------------------------------------------------------------------------------------------------------------------------------|
| Namespace       | 2021-01-01-preview | v1beta20210101preview      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/servicebus/v1beta20210101preview_namespace.yaml)            |
| Namespace       | 2021-01-01-preview | v1alpha1api20210101preview | v2.0.0-alpha.1 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/servicebus/v1alpha1api20210101preview_namespace.yaml)       |
| NamespacesQueue | 2021-01-01-preview | v1beta20210101preview      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/servicebus/v1beta20210101preview_namespacesqueue.yaml)      |
| NamespacesQueue | 2021-01-01-preview | v1alpha1api20210101preview | v2.0.0-alpha.1 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/servicebus/v1alpha1api20210101preview_namespacesqueue.yaml) |
| NamespacesTopic | 2021-01-01-preview | v1beta20210101preview      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/servicebus/v1beta20210101preview_namespacestopic.yaml)      |
| NamespacesTopic | 2021-01-01-preview | v1alpha1api20210101preview | v2.0.0-alpha.1 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/servicebus/v1alpha1api20210101preview_namespacestopic.yaml) |

## Signalrservice

| Resource | ARM Version | CRD Version         | Supported From | Sample                                                                                                                              |
|----------|-------------|---------------------|----------------|-------------------------------------------------------------------------------------------------------------------------------------|
| SignalR  | 2021-10-01  | v1beta20211001      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/signalrservice/v1beta20211001_signalr.yaml)      |
| SignalR  | 2021-10-01  | v1alpha1api20211001 | v2.0.0-alpha.4 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/signalrservice/v1alpha1api20211001_signalr.yaml) |

## Storage

| Resource                             | ARM Version | CRD Version         | Supported From | Sample                                                                                                                                                    |
|--------------------------------------|-------------|---------------------|----------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------|
| StorageAccount                       | 2021-04-01  | v1beta20210401      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/storage/v1beta20210401_storageaccount.yaml)                            |
| StorageAccount                       | 2021-04-01  | v1alpha1api20210401 | v2.0.0-alpha.1 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/storage/v1alpha1api20210401_storageaccount.yaml)                       |
| StorageAccountsBlobService           | 2021-04-01  | v1beta20210401      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/storage/v1beta20210401_storageaccountsblobservice.yaml)                |
| StorageAccountsBlobService           | 2021-04-01  | v1alpha1api20210401 | v2.0.0-alpha.1 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/storage/v1alpha1api20210401_storageaccountsblobservice.yaml)           |
| StorageAccountsBlobServicesContainer | 2021-04-01  | v1beta20210401      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/storage/v1beta20210401_storageaccountsblobservicescontainer.yaml)      |
| StorageAccountsBlobServicesContainer | 2021-04-01  | v1alpha1api20210401 | v2.0.0-alpha.1 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/storage/v1alpha1api20210401_storageaccountsblobservicescontainer.yaml) |
| StorageAccountsManagementPolicy      | 2021-04-01  | v1beta20210401      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/storage/v1beta20210401_storageaccountsmanagementpolicy.yaml)           |
| StorageAccountsQueueService          | 2021-04-01  | v1beta20210401      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/storage/v1beta20210401_storageaccountsqueueservice.yaml)               |
| StorageAccountsQueueService          | 2021-04-01  | v1alpha1api20210401 | v2.0.0-alpha.5 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/storage/v1alpha1api20210401_storageaccountsqueueservice.yaml)          |
| StorageAccountsQueueServicesQueue    | 2021-04-01  | v1beta20210401      | v2.0.0-beta.0  | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/storage/v1beta20210401_storageaccountsqueueservicesqueue.yaml)         |
| StorageAccountsQueueServicesQueue    | 2021-04-01  | v1alpha1api20210401 | v2.0.0-alpha.5 | [View](https://github.com/Azure/azure-service-operator/tree/main/v2/config/samples/storage/v1alpha1api20210401_storageaccountsqueueservicesqueue.yaml)    |

