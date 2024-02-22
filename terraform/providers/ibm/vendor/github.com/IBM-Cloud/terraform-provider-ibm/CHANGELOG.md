# 1.61.0 (Jan 05, 2024)
Features
* Support for MQ on Cloud
    - **Datasources**
       - ibm_mqcloud_queue_manager
       - ibm_mqcloud_queue_manager_status
       - ibm_mqcloud_application 
       - ibm_mqcloud_user
       - ibm_mqcloud_keystore_certificate
       - ibm_mqcloud_truststore_certificate
    - **Resources**
       - ibm_mqcloud_queue_manager
       - ibm_mqcloud_application 
       - ibm_mqcloud_user
       - ibm_mqcloud_keystore_certificate
       - ibm_mqcloud_truststore_certificate
* Support for Secret Manager
     - **Datasources**
        - ibm_sm_service_credentials_secret_metadata
        - ibm_sm_service_credentials_secret
     - **Resources**
        - ibm_sm_service_credentials_secret
* Support for VPC
     - **Datasources**
        - ibm_is_snapshot_consistency_group
        - ibm_is_snapshot_consistency_groups
     - **Resources**
        - ibm_is_snapshot_consistency_group

Enhancements
* feat(Cloud Databases): Database user password complexity validation ([4931](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4931))
* Update pi_user_data to accept string input ([4974](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4974))
* support host_link_agent_endpoint for Satellite host ([4970](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4970))
* Add mtu and accessConfig flags to subnet create commands for terraform ([4690](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4690))
* feat(Cloud Databases): Redis Database User RBAC support ([4982](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4982))
* fix(Cloud Databases): fix Unwrap return value for go 1.18 compat ([4991](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4991))
* update issue fixed ibm_is_subnet_reserved_ip ([4988](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4988))
* Adding Flexible IOPS ([4992](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4992))
* Removing Support For Power VPN Create ([4993](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4993))
* Feature(share-crr): Share cross region replication ([4995](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4995))
* Enhancement: Added operating system attributes to is images datasources ([4998](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4998))
* added enhancement to one step delegate resolver in is_vpc ([5000](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/5000))
* resolved delete issue for the floated nics on bm server ([5001](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/5001))
* Regenerate projects provider based off the latest go sdk ([5003](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/5003))
* Support route advertising in vpc ([5005](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/5005))
* Add a nil check for boottarget of bms ([5014](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/5014))
* Delete wait logic changes ([5017](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/5017))

BugFixes
* Fix IBM pi documentation bug ([4969](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4969))
* Incorrect key_algorithm handling forces delete & replace of ibm_sm_private_certificate on every apply ([4978](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4978))
* ibm_sm_private_certificate_configuration_template arguments ttl and max_ttl are not documented ([4977](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4977))
* ibm_sm_private_certificate unsupported argument: rotation.rotate_keys ([4976](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4976))
* data ibm_schematics_workspace bug ([4990](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4990))
* Secret Manager docs bug fix ([5018](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/5018))

# 1.61.0-beta0 (Dec 20, 2023)
Features
* Support for MQ on Cloud
    - **Datasources**
       - ibm_mqcloud_queue_manager
       - ibm_mqcloud_queue_manager_status
       - ibm_mqcloud_application 
       - ibm_mqcloud_user
       - ibm_mqcloud_keystore_certificate
       - ibm_mqcloud_truststore_certificate
    - **Resources**
       - ibm_mqcloud_queue_manager
       - ibm_mqcloud_application 
       - ibm_mqcloud_user
       - ibm_mqcloud_keystore_certificate
       - ibm_mqcloud_truststore_certificate
* Support for Secret Manager
     - **Datasources**
        - ibm_sm_service_credentials_secret_metadata
        - ibm_sm_service_credentials_secret
     - **Resources**
        - ibm_sm_service_credentials_secret
* Support for VPC
     - **Datasources**
        - ibm_is_snapshot_consistency_group
        - ibm_is_snapshot_consistency_groups
     - **Resources**
        - ibm_is_snapshot_consistency_group

Enhancements
* feat(Cloud Databases): Database user password complexity validation ([4931](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4931))
* Update pi_user_data to accept string input ([4974](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4974))
* support host_link_agent_endpoint for Satellite host ([4970](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4970))
* Add mtu and accessConfig flags to subnet create commands for terraform ([4690](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4690))
* feat(Cloud Databases): Redis Database User RBAC support ([4982](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4982))
* fix(Cloud Databases): fix Unwrap return value for go 1.18 compat ([4991](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4991))
* update issue fixed ibm_is_subnet_reserved_ip ([4988](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4988))
* Adding Flexible IOPS ([4992](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4992))
* Removing Support For Power VPN Create ([4993](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4993))
* Feature(share-crr): Share cross region replication ([4995](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4995))
* Enhancement: Added operating system attributes to is images datasources ([4998](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4998))
* added enhancement to one step delegate resolver in is_vpc ([5000](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/5000))
* resolved delete issue for the floated nics on bm server ([5001](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/5001))

BugFixes
* Fix IBM pi documentation bug ([4969](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4969))

# 1.60.1 (Nov 09, 2023)

BugFixes
* Regenerate Projects TF to fix generated doc and samples
 ([4961](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4961))

# 1.60.0 (Nov 29, 2023)
Features
* Support for Projects
   - **Datasources**
       - ibm_project
       - ibm_project_config 
       - ibm_project_environment
   - **Resources**
        - ibm_project
        - ibm_project_config 
        - ibm_project_environment

* Support for Code Engine
   - **Datasources**
       - ibm_code_engine_domain_mapping
   - **Resources**
        - ibm_code_engine_domain_mapping

* Support for Power Instance
   - **Resources**
        - ibm_pi_workspace

Enhancements
* support offline restore for MongoDB EE PITR ([4601](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4601))
* bump ContinuousDelivery Go SDK version ([4918](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4918))
* Added nest conditions to rule.conditions IAM Policies ([4896](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4896))
* Updates to SCC tool ([4920](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4920))
* Add entitlement option to Satellite cluster/workerpool create ([4894](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4894))
* VPC ID Filter is added when Subnet Name is specified ([4892](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4892))
* add optional account id to kms config ([4944](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4944))

BugFixes
* resolved nil pointer issue on vpn gateway resource ([4903](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4903))
* Private/direct COS endpoint settings conflicts with IBM Cloud docs, VPE options, and COS config endpoint ([4919](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4919))
* add missing required argument name to the doc ([4909](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4909))
* Fix wrong sintax in doc for Ingress Secret Opaque ([4917](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4917))
* schematics agent doc fixes ([4933](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4933))
* Fix some job parameters for code engine ([4923](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4923))
* CIS - remove deafult value for min_tls_version ([4947](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4947))
* Update Power Workspace/s and Datacenter/s DataSource Documentation ([4904](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4904))
* validation on encryption with catalog images fixed ([4940](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4940))
* issue-13603-fix attachment terraform ([4952](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4952))
* Fixed catalog service extensions and values metadata params of workspace ds ([4957](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4957))

# 1.60.0-beta1 (Nov 23, 2023)
Features
* Support for Projects
   - **Datasources**
       - ibm_project
       - ibm_project_config 
       - ibm_project_environment
   - **Resources**
        - ibm_project
        - ibm_project_config 
        - ibm_project_environment

* Support for Code Engine
   - **Datasources**
       - ibm_code_engine_domain_mapping
   - **Resources**
        - ibm_code_engine_domain_mapping

* Support for Power Instance
   - **Resources**
        - ibm_pi_workspace

Enhancements
* support offline restore for MongoDB EE PITR ([4601](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4601))
* bump ContinuousDelivery Go SDK version ([4918](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4918))
* Added nest conditions to rule.conditions IAM Policies ([4896](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4896))
* Updates to SCC tool ([4920](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4920))
* Add entitlement option to Satellite cluster/workerpool create ([4894](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4894))
* VPC ID Filter is added when Subnet Name is specified ([4892](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4892))
* add optional account id to kms config ([4944](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4944))

BugFixes
* resolved nil pointer issue on vpn gateway resource ([4903](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4903))
* Private/direct COS endpoint settings conflicts with IBM Cloud docs, VPE options, and COS config endpoint ([4919](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4919))
* add missing required argument name to the doc ([4909](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4909))
* Fix wrong sintax in doc for Ingress Secret Opaque ([4917](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4917))
* schematics agent doc fixes ([4933](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4933))
* Fix some job parameters for code engine ([4923](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4923))
* CIS - remove deafult value for min_tls_version ([4947](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4947))
* Update Power Workspace/s and Datacenter/s DataSource Documentation ([4904](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4904))
* validation on encryption with catalog images fixed ([4940](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4940))

# 1.60.0-beta0 (Nov 15, 2023)
Features
* Support for Projects
   - **Datasources**
       - ibm_project
       - ibm_project_config 
   - **Resources**
        - ibm_project
        - ibm_project_config 
Enhancements
* support offline restore for MongoDB EE PITR ([4601](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4601))
* bump ContinuousDelivery Go SDK version ([4918](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4918))
* Added nest conditions to rule.conditions IAM Policies ([4896](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4896))
* Updates to SCC tool ([4920](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4920))

BugFixes
* resolved nil pointer issue on vpn gateway resource ([4903](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4903))
* Private/direct COS endpoint settings conflicts with IBM Cloud docs, VPE options, and COS config endpoint ([4919](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4919))
* add missing required argument name to the doc ([4909](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4909))

# 1.59.1 (Nov 20, 2023)

Bug Fixes
* ibm_schematics_workspace: provider crash during terraform plan ([4907](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4907))


# 1.59.0 (Oct 30, 2023)

Features
* Support Usage Reports
    - **Datasources**
        - ibm_billing_snapshot_list
    - **Resources**
        - ibm_billing_report_snapshot

* Support Power Instance
    - **Datasources**
        - ibm_pi_workspace
        - ibm_pi_workspaces
        - ibm_pi_datacenter
        - ibm_pi_datacenters

* Support Schematics Agents
    - **Datasources**
        - ibm_schematics_policies
        - ibm_schematics_policy
        - ibm_schematics_agents
        - ibm_schematics_agent
        - ibm_schematics_agent_prs
        - ibm_schematics_agent_deploy
        - ibm_schematics_agent_health
    - **Resources**
        - ibm_schematics_policy
        - ibm_schematics_agent
        - ibm_schematics_agent_prs
        - ibm_schematics_agent_deploy
        - ibm_schematics_agent_health

* Support Event Notification
    - **Datasources**
        - ibm_en_destination_custom_email
        - ibm_en_subscription_custom_email
    - **Resources**
        - ibm_en_destination_custom_email
        - ibm_en_subscription_custom_email


Enhancements
* Get secret by name ([4825](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4825))
* VPN for VPC: Customer should be able to recover their gateway or server for unhealthy status ([4858](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4858))
* support for tf 1.5 in schematics workspace ([4853](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4853))
* Deprecated match_resource_types and Intoduced match_resource_type ([4863](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4863))
* Enterprise BaaS feature ([4845](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4845))
* sarama golang library update ([4810](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4810))
* Adding NUMA and Profile Status to instance, instance profile and dedicated hosts ([4871](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4871))
* update terraform as per latest eventstreams go sdk release 1.4.0 ([4862](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4862))
* Add path to invoke update without determining changed CRN data, add validator for name of secrets ([4859](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4859))
* Feature: ReplicationEnabledField for Storage Pool ([4875](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4875))
* Remove deprecated scaling attributes ([4481](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4481))
* update CD Go SDK version ([4887](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4887))
* Adding updates to the scc resources/datasources ([4865](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4865))
* Add support for Elasticsearch Platinum ([4712](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4712))
* Add support for security groups for network load balancers ([4888](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4888))


Bug Fixes
* Fix handling of bundle_certs in Secrets Manager public cert ([4854](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4854))
* add description in docs for Key Protect ([4846](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4846))
* Update iam_service_policy.html.markdown ([4836](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4836))
* v1.58.0 ibm_container_cluster_config: new endpoint_type returning self-signed private endpoint ([4861](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4861))
* Doc correction: Share mount target doc corrections ([4889](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4889))
* Fix indentation for subcategory ([4891](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4891))

# 1.59.0-beta0 (Oct 25, 2023)

Features
* Support Usage Reports
    - **Datasources**
        - ibm_billing_snapshot_list
    - **Resources**
        - ibm_billing_report_snapshot

* Support Power Instance
    - **Datasources**
        - ibm_pi_workspace
        - ibm_pi_workspaces
        - ibm_pi_datacenter
        - ibm_pi_datacenters

* Support Schematics Agents
    - **Datasources**
        - ibm_schematics_policies
        - ibm_schematics_policy
        - ibm_schematics_agents
        - ibm_schematics_agent
        - ibm_schematics_agent_prs
        - ibm_schematics_agent_deploy
        - ibm_schematics_agent_health
    - **Resources**
        - ibm_schematics_policy
        - ibm_schematics_agent
        - ibm_schematics_agent_prs
        - ibm_schematics_agent_deploy
        - ibm_schematics_agent_health

* Support Event Notification
    - **Datasources**
        - ibm_en_destination_custom_email
        - ibm_en_subscription_custom_email
    - **Resources**
        - ibm_en_destination_custom_email
        - ibm_en_subscription_custom_email


Enhancements
* Get secret by name ([4825](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4825))
* VPN for VPC: Customer should be able to recover their gateway or server for unhealthy status ([4858](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4858))
* support for tf 1.5 in schematics workspace ([4853](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4853))
* Deprecated match_resource_types and Intoduced match_resource_type ([4863](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4863))
* Enterprise BaaS feature ([4845](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4845))
* sarama golang library update ([4810](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4810))
* Adding NUMA and Profile Status to instance, instance profile and dedicated hosts ([4871](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4871))
* update terraform as per latest eventstreams go sdk release 1.4.0 ([4862](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4862))
* Add path to invoke update without determining changed CRN data, add validator for name of secrets ([4859](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4859))
* Feature: ReplicationEnabledField for Storage Pool ([4875](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4875))

Bug Fixes
* Fix handling of bundle_certs in Secrets Manager public cert ([4854](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4854))
* add description in docs for Key Protect ([4846](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4846))
* Update iam_service_policy.html.markdown ([4836](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4836))
* v1.58.0 ibm_container_cluster_config: new endpoint_type returning self-signed private endpoint ([4861](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4861))

# 1.58.1 (Oct 04, 2023)

Bug Fixes
* Timing issue while destroying Key Protect resources ([4837](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4837))
* ibm_cos_bucket data lookup is throwing NoSuchWebsiteConfiguration in new version 1.58.0 ([4838](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4838))
* Metrics router and atracker: Updated platform-services-go-sdk to fetch Madrid endpoint ([4830](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4830))

# 1.58.0 (Sep 29, 2023)

Features
* Support Satellite Cluster
    - **Resources**
        - ibm_satellite_storage_configuration
        - ibm_satellite_storage_assignment
    - **Datasources**
        - ibm_satellite_storage_configuration
        - ibm_satellite_storage_assignment
* Support Security and Compliance
    - **Resources**
        - ibm_scc_rule
        - ibm_scc_control_library
        - ibm_scc_profile
        - ibm_scc_profile_attachment
        - ibm_scc_provider_type_instance
    - **Datasources**
        - ibm_scc_instance_settings
        - ibm_scc_control_library
        - ibm_scc_profile
        - ibm_scc_profile_attachment
        - ibm_scc_provider_type
        - ibm_scc_provider_type_collection
        - ibm_scc_provider_type_instance
        - ibm_scc_latest_reports
        - ibm_scc_report
        - ibm_scc_report_controls
        - ibm_scc_report_evaluations
        - ibm_scc_report_resources
        - ibm_scc_report_rule
        - ibm_scc_report_summary
        - ibm_scc_report_tags
        - ibm_scc_report_violation_drift
        - ibm_scc_rule

* Support CD Toolchain
     - **Datasources**
        - ibm_cd_toolchains
* Support Virtual Private Cloud
    - **Resources**
        - ibm_is_vpc_dns_resolution_binding
    - **Datasources**
        - ibm_is_vpc_dns_resolution_binding
        - ibm_is_vpc_dns_resolution_bindings


Enhancements
* Added retry mechanism and new SDK generator 3.78 ([4776](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4776))
* Add default cluster versions to cluster versions data source ([4799](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4799))
* Add description for keys and force_delete for deleteKeyRings for IBM ([4767](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4767))
* Retry cloud connection create/update when vpc is unavailable ([4766](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4766))
* Adding support for COS Static Web hosting ([4766](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4766))
* add support for endpoint parameter in cluster_config ([4793](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4793))
* fix(IAM Policy Management): allow sourceServiceName to be optional for authorizational policies ([4804](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4804))

BugFixes
* ops_manager User Creation ([4755](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4755))
* fix(share-iops): Share Iops range fix for dp2 ([4807](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4807))
* fix(VSI-Profile-patch): Remove validation for VSI profile patching ([4824](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4824))


# 1.58.0-beta1 (Sep 25, 2023)
Features
* Support CD Toolchain
     - **Datasources**
        - ibm_cd_toolchains
* Support Virtual Private Cloud
    - **Resources**
        - ibm_is_vpc_dns_resolution_binding
    - **Datasources**
        - ibm_is_vpc_dns_resolution_binding
        - ibm_is_vpc_dns_resolution_bindings

# 1.58.0-beta0 (Sep 10, 2023)

Features
* Support Satellite Cluster
    - **Resources**
        - ibm_satellite_storage_configuration
        - ibm_satellite_storage_assignment
    - **Datasources**
        - ibm_satellite_storage_configuration
        - ibm_satellite_storage_assignment
* Support Security and Compliance
    - **Resources**
        - ibm_scc_rule
        - ibm_scc_control_library
        - ibm_scc_profile
        - ibm_scc_profile_attachment
        - ibm_scc_provider_type_instance
    - **Datasources**
        - ibm_scc_instance_settings
        - ibm_scc_control_library
        - ibm_scc_profile
        - ibm_scc_profile_attachment
        - ibm_scc_provider_type
        - ibm_scc_provider_type_collection
        - ibm_scc_provider_type_instance
        - ibm_scc_latest_reports
        - ibm_scc_report
        - ibm_scc_report_controls
        - ibm_scc_report_evaluations
        - ibm_scc_report_resources
        - ibm_scc_report_rule
        - ibm_scc_report_summary
        - ibm_scc_report_tags
        - ibm_scc_report_violation_drift
        - ibm_scc_rule

Enhancements
* Added retry mechanism and new SDK generator 3.78 ([4776](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4776))
* Add default cluster versions to cluster versions data source ([4799](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4799))
* Add description for keys and force_delete for deleteKeyRings for IBM ([4767](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4767))
* Retry cloud connection create/update when vpc is unavailable ([4766](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4766))
* Adding support for COS Static Web hosting ([4766](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4766))
* add support for endpoint parameter in cluster_config ([4793](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4793))

BugFixes
* ops_manager User Creation ([4755](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4755))
* fix(share-iops): Share Iops range fix for dp2 ([4807](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4807))


# 1.57.0 (Sep 13, 2023)

Features
* Support IAM Trusted Profile
    - **Resources**
        - ibm_iam_trusted_profile_identity
    - **Datasources**
        - ibm_iam_trusted_profile_identity
        - ibm_iam_trusted_profile_identities

* Support IAM Identity Enterprise Templates 
     - **Resources**
        - ibm_iam_account_settings_template
        - ibm_iam_trusted_profile_template
        - ibm_iam_account_settings_template_assignment
        - ibm_iam_account_settings_template_assignment
    - **Datasources**
        - ibm_iam_account_settings_template
        - ibm_iam_trusted_profile_template
        - ibm_iam_account_settings_template_assignment
        - ibm_iam_trusted_profile_template_assignment

* Support IAM Access Group Templates 
    - **Resources**
        - ibm_iam_access_group_template
        - ibm_iam_access_group_template_version
        - ibm_iam_access_group_template_assignment
    - **Datasources**
        - ibm_iam_access_group_template_versions
        - ibm_iam_access_group_template_assignment
    
* Support IAM Policy Templates 
    - **Resources**
        - ibm_iam_policy_template
        - ibm_iam_policy_template_version
    - **Datasources**
        - ibm_iam_policy_template
        - ibm_iam_policy_template_version
        - ibm_iam_policy_assignments
        - ibm_iam_policy_assignment

     
Enhancements
* Support `instance_crn` argument for ibm_cd_toolchain_tool_securitycompliance ([4746](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4746))
* Remove `whitelist` argument for ibm_database ([4714](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4714))
* Remove deprecated share target resource and data sources ([4739](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4739))
* fix force_new to resource fields for ibm_iam_trusted_profile_identity ([4762](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4762))
* CD Toolchain SCC tool param description update ([4753](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4753))
* enhancement(File-share): GA Preview, beta SDK upgraded ([4770](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4770))
* feat: support new target_account_contexts field in Catalog Management ([4773](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4773))
* ODF Example Doc Update ([4757](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4757))
* Feature(File-share-GA): Promoting File share from Beta to GA ([4759](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4759))
* Share profile data source added with capacity and iops ([4789](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4789))
* remove update from IAM Access Group Templates ([4796](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4796))

BugFixes
* Fix documentation of Secrets Manager private certificate resource ([4760](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4760))
* Update iam_access_group_members.html.markdown ([4760](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4760))
* fix(tekton): Trigger updates ([4731](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4731))
* feat: Fixed TGW Route report update issue dependency: None ([4777](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4777))
* DeleteVPNServerWithContext failure occurring when attempting to destroy ([4758](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4758))
* fix(tekton): update trigger sample ([4787](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4787))

# 1.57.0-beta0 (Sep 05, 2023)

Features
* Support IAM Trusted Profile
    - **Resources**
        - ibm_iam_trusted_profile_identity
    - **Datasources**
        - ibm_iam_trusted_profile_identity
        - ibm_iam_trusted_profile_identities

* Support IAM Identity Enterprise Templates 
     - **Resources**
        - ibm_iam_account_settings_template
        - ibm_iam_trusted_profile_template
        - ibm_iam_account_settings_template_assignment
        - ibm_iam_account_settings_template_assignment
    - **Datasources**
        - ibm_iam_account_settings_template
        - ibm_iam_trusted_profile_template
        - ibm_iam_account_settings_template_assignment
        - ibm_iam_trusted_profile_template_assignment

* Support IAM Access Group Templates 
    - **Resources**
        - ibm_iam_access_group_template
        - ibm_iam_access_group_template_version
        - ibm_iam_access_group_template_assignment
    - **Datasources**
        - ibm_iam_access_group_template_versions
        - ibm_iam_access_group_template_assignment
    
* Support IAM Policy Templates 
    - **Resources**
        - ibm_iam_policy_template
        - ibm_iam_policy_template_version
    - **Datasources**
        - ibm_iam_policy_template
        - ibm_iam_policy_template_version
        - ibm_iam_policy_assignments
        - ibm_iam_policy_assignment

     
Enhancements
* Support `instance_crn` argument for ibm_cd_toolchain_tool_securitycompliance ([4746](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4746))
* Remove `whitelist` argument for ibm_database ([4714](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4714))
* Remove deprecated share target resource and data sources ([4739](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4739))
* fix force_new to resource fields for ibm_iam_trusted_profile_identity ([4762](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4762))
* CD Toolchain SCC tool param description update ([4753](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4753))
* enhancement(File-share): GA Preview, beta SDK upgraded ([4770](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4770))
* feat: support new target_account_contexts field in Catalog Management ([4773](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4773))
* ODF Example Doc Update ([4757](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4757))
* Feature(File-share-GA): Promoting File share from Beta to GA ([4759](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4759))

BugFixes
* Fix documentation of Secrets Manager private certificate resource ([4760](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4760))
* Update iam_access_group_members.html.markdown ([4760](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4760))
* fix(tekton): Trigger updates ([4731](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4731))
* feat: Fixed TGW Route report update issue dependency: None ([4777](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4777))

# 1.56.2 (Aug 29, 2023)

BugFixes
* upgrades to new beta SDK version which has updated version date for vpc api ([4770](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4770))

# 1.56.1 (Aug 21, 2023)

BugFixes
* Fix adding of schematics tags if they are service/access tags ([4755](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4755))
* allow workspace type 1.4 for schematic ([4756](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4756))

# 1.56.0 (Aug 08, 2023)

Removal
* Remove SCC V1 version resources/datasources ([4689](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4689))
* Remove atracker v1 changes  ([4684](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4684))

Features
* Support Secret Manager
    - **Resources**
        - ibm_sm_public_certificate_action_validate_manual_dns

* Support IAM
    - **Resources**
        - ibm_iam_user_mfa_enrollments

* Support VPC
    - **Datasources**
        - ibm_is_virtual_network_interface
        - ibm_is_virtual_network_interfaces
* Support EVent Notification
    - **Resources**
        - ibm_en_destination_huawei
        - ibm_en_subscription_huawei
        - ibm_en_ibmsource
    - **Datasources**
        - ibm_en_destination_huawei
        - ibm_en_subscription_huawei
        - ibm_en_ibmsource

     
Enhancements
* Modified TGW connection resource to support powerVS dependency ([4657](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4657))
* DL Route Report changes ([4570](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4570))
* added trait field in enterprise child account ([4696](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4696))
* Tekton simplify CRON validation ([4668](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4668))
* Remove Key Protect ID and Transition to v5 API ([4420](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4420))
* Git broker parameter updates ([4675](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4675))
* fix(4630): added computed to lb pool member target id and address ([4642](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4642))
* add(image): support for image lifecycle management ([4698](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4698))
* SSH Key: Support additional encryption algorithms ([4677](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4677))
* backup: added suuport for cross region copy for vpc backup ([4692](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4692))
* enhancement(snapshot): added support for cross region snapshot copy ([4678](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4678))
* Update container-registry SDK  ([4706](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4706))
* Remove deprecated connection string and cert path from database datasource ([4406](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4406))
* Akamai authentication method upgraded ([4700](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4700))
* support `no_sg_acl_rules` argument in VPC ([4702](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4702))
* enhancement(bm): added support for console type and nic count for bm server ([4716](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4716))


BugFixes
* update doc files for Projects ([4681](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4681))
* handle nil network id in dhcp detail ([4683](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4683))
* fix: CIS - documentation fix ([4688](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4688))
* Update initialization block ([4675](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4675))
* fix import example in code engine document ([4646](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4646))
* pDNS documention correction for CR forwarding rule ([4695](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4695))
* always read and set datacentervalue ([4703](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4703))
* Doc Fix: Updated IBM Cloud Metrics Routing ([4705](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4705))
* Failed to update resource ibm_is_instance_group_manager_policy update ([4647](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4647))
* Crash while creating COS bucket ([4699](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4699))
* fix(Cloud Databases): marshal blank string for latest PITR time ([4718](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4718))
* doc update for container addons ([4730](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4730))
* Doc fix: Updated subcategory of ibm_metrics_router_targets ([4737](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4737))
* fixed(backup_policy_plan): delete_over_count null issue fixed ([4735](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4735))
* move trait field from acc import to acc create ([4729](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4729))
* Fix ibm_pi_dhcps example ([4727](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4727))
* Updated the documentation to cleanup l1bm ([4741](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4741))
* fix(ssh): resolved password decryption for pkcs8 ([4707](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4707))
* Update ibm_pi_instance Doc ([4728](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4728))
* Fix ibm_database plan documentation error ([4742](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4742))

# 1.56.0-beta0 (Jul 25, 2023)

Removal
* Remove SCC V1 version resources/datasources ([4689](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4689))
* Remove atracker v1 changes  ([4684](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4684))

Features
* Support Secret Manager
    - **Resources**
        - ibm_sm_public_certificate_action_validate_manual_dns

* Support IAM
    - **Resources**
        - ibm_iam_user_mfa_enrollments

* Support VPC
    - **Datasources**
        - ibm_is_virtual_network_interface
        - ibm_is_virtual_network_interfaces
* Support EVent Notification
    - **Resources**
        - ibm_en_destination_huawei
        - ibm_en_subscription_huawei
        - ibm_en_ibmsource
    - **Datasources**
        - ibm_en_destination_huawei
        - ibm_en_subscription_huawei
        - ibm_en_ibmsource

     
Enhancements
* Modified TGW connection resource to support powerVS dependency ([4657](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4657))
* DL Route Report changes ([4570](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4570))
* added trait field in enterprise child account ([4696](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4696))
* Tekton simplify CRON validation ([4668](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4668))
* Remove Key Protect ID and Transition to v5 API ([4420](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4420))
* Git broker parameter updates ([4675](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4675))
* fix(4630): added computed to lb pool member target id and address ([4642](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4642))
* add(image): support for image lifecycle management ([4698](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4698))
* SSH Key: Support additional encryption algorithms ([4677](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4677))
* backup: added suuport for cross region copy for vpc backup ([4692](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4692))
* enhancement(snapshot): added support for cross region snapshot copy ([4678](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4678))
* Update container-registry SDK  ([4706](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4706))
* Remove deprecated connection string and cert path from database datasource ([4406](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4406))

BugFixes
* update doc files for Projects ([4681](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4681))
* handle nil network id in dhcp detail ([4683](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4683))
* fix: CIS - documentation fix ([4688](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4688))
* Update initialization block ([4675](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4675))
* fix import example in code engine document ([4646](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4646))
* pDNS documention correction for CR forwarding rule ([4695](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4695))
* always read and set datacentervalue ([4703](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4703))
* Doc Fix: Updated IBM Cloud Metrics Routing ([4705](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4705))
* Failed to update resource ibm_is_instance_group_manager_policy update ([4647](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4647))

# 1.55.0 (Jul 10, 2023)

Features
* Support VPC
    - **Resources**
        - ibm_is_share_mount_target
    - **Datasources**
        - ibm_is_share_mount_target
        - ibm_is_share_mount_targets

* Support DNS
    - **Resources**
        - ibm_dns_linked_zone
       
* Support IKS
    - **Datasources**
        - ibm_container_ingress_secret_tls
        - ibm_container_ingress_secret_opaque
    - **Resources**
        - ibm_container_ingress_secret_tls
        - ibm_container_ingress_secret_opaque

* Support Metrics Router
    - **Datasources**
        - ibm_metrics_router_targets
        - ibm_metrics_router_routes
    - **Resources**
        - ibm_metrics_router_route
        - ibm_metrics_router_target
        - ibm_metrics_router_settings

* Support Code Engine
    - **Datasources**
        - ibm_code_engine_binding
    - **Resources**
        - ibm_code_engine_binding

* Support CIS
    - **Datasources**
        - ibm_cis_bot_managements
        - ibm_cis_bot_analytics
    - **Resources**
        - ibm_cis_bot_management
     
Enhancements
* Secrets Manager - Additional tests ([4613](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4613))
* Tekton: Improved property type handling ([4595](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4595))
* Tagging for CD Toolchains ([4607](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4607))
* Support for ODF Worker Replace ([4600](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4600))
* Documentation for ODF Add-on and Worker Replace ([4627](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4627))
* Document private net DNS limit on PER workspaces ([4629](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4629))
* added support for vcpu manufacturer(vsi, dh) ([4637](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4637))
* Adding changes for replacement of key_protect parameter ([4618](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4618))
* Adding more examples and test cases for Atracker ([4645](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4645))

BugFixes
* fix(iam-service-api-key): added nil check on apikey *string ([4617](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4617))
* update project_instance.html.markdown ([4621](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4621))
* add import example of code engine ([4623](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4623))
* Update pi_cloud_connection documentation ([4625](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4625))
* Documenation correction instance group membership ([4632](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4632))
* CIS - Documentation Update for plan  ([4638](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4638))
* Fix indentation  ([4633](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4633))
* managed_addons fix for container_addons ([4606](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4606))
* cos: importing a bucket does not import the key_protect attribute ([3394](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3394))
* Private cert attribute fixes ([4641](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4641))
* auto_rotate fix ([4649](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4649))
* Fix the private endpoint for global search API ([4666](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4666))
* Fix dnssvcs module broken after SDK release ([4667](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4667))
* DNS ut fix: strfmt to string conversion to fix broken UT ([4671](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4671))
* fix: catalog management version resource patch fix to correct operation type ([4673](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4673))
* Fix unclosed code block in r/iam_access_group_policy docs ([4659](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4659))


# 1.55.0-beta0(Jun 21, 2023)

Features
* Support VPC
    - **Resources**
        - ibm_is_share_mount_target
    - **Datasources**
        - ibm_is_share_mount_target
        - ibm_is_share_mount_targets

* Support DNS
    - **Resources**
        - ibm_dns_linked_zone
       
* Support IKS
    - **Datasources**
        - ibm_container_ingress_secret_tls
        - ibm_container_ingress_secret_opaque
    - **Resources**
        - ibm_container_ingress_secret_tls
        - ibm_container_ingress_secret_opaque

* Support Metrics Router
    - **Datasources**
        - ibm_metrics_router_targets
        - ibm_metrics_router_routes
    - **Resources**
        - ibm_metrics_router_route
        - ibm_metrics_router_target
        - ibm_metrics_router_settings

* Support Code Engine
    - **Datasources**
        - ibm_code_engine_binding
    - **Resources**
        - ibm_code_engine_binding

* Support CIS
    - **Datasources**
        - ibm_cis_bot_managements
    - **Resources**
        - ibm_cis_bot_management
     
Enhancements
* Secrets Manager - Additional tests ([4613](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4613))
* Tekton: Improved property type handling ([4595](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4595))
* Tagging for CD Toolchains ([4607](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4607))
* Support for ODF Worker Replace ([4600](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4600))
* Documentation for ODF Add-on and Worker Replace ([4627](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4627))
* Document private net DNS limit on PER workspaces ([4629](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4629))
* added support for vcpu manufacturer(vsi, dh) ([4637](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4637))
* Adding changes for replacement of key_protect parameter ([4618](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4618))

BugFixes
* fix(iam-service-api-key): added nil check on apikey *string ([4617](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4617))
* update project_instance.html.markdown ([4621](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4621))
* add import example of code engine ([4623](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4623))
* Update pi_cloud_connection documentation ([4625](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4625))
* Documenation correction instance group membership ([4632](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4632))
* CIS - Documentation Update for plan  ([4638](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4638))
* Fix indentation  ([4633](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4633))
* managed_addons fix for container_addons ([4606](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4606))
* cos: importing a bucket does not import the key_protect attribute ([3394](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3394))


# 1.54.0(Jun 07, 2023)

Features
* Support Project
    - **Resources**
        - ibm_project_instance
    - **Datasources**
        - ibm_project_event_notification

* Support Event Notification
    - **Resources**
        - ibm_en_destination_cos
        - ibm_en_subscription_cos
    - **Datasources**
        - ibm_en_destination_cos
        - ibm_en_subscription_cos
* Support Code Engine
    - **Datasources**
        - ibm_code_engine_app
        - ibm_code_engine_build
        - ibm_code_engine_config_map
        - ibm_code_engine_job
        - ibm_code_engine_job

Enhancements
* suppress resize for autoscaled workerpools for autoscale enabled cluster ([4533](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4533))
* Feature(Image Export Jobs): Support for image export jobs ([4566](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4566))
* feat(Enterprise):added support for destroy method ([4565](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4565))
* ODF Terraform Documentation Updated ([4549](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4549))
* CD Security Compliance param deprecation ([4581](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4581))
* Remove weak ciphers support from ike and ipsec policies ([4584](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4584))
* feat(Cloud Databases): Deprecate auto_scaling cpu attribute ([4599](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4599))


BugFixes
* disable creating new resource when endpoint type changes ([4554](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4554))
* retry added for read authorization policy ([4554](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4413))
* Update example for ibm_is_vpc_routing_table ([4564](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4564))
* IBM Cloud Object destination resources and data sources and docs update ([4534](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4534))
* updated docs to not include force new resource ([4577](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4577))
* fix(doc): doc fix for vpc_routing_table_route next_hop ([4586](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4586))
* addon doc corrected  ([4596](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4596)) 
* fix(Cloud Databases): use correct redis maxmemory configuration ([4588](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4588))
* corrected the argument repetation. ([4604](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4604))
* fix(bm) : pci nic deletion change ([4532](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4532))

# 1.54.0-beta0(May 13, 2023)

Features
* Support Project
    - **Resources**
        - ibm_project_instance
    - **Datasources**
        - ibm_project_event_notification

* Support Event Notification
    - **Resources**
        - ibm_en_destination_cos
        - ibm_en_subscription_cos
    - **Datasources**
        - ibm_en_destination_cos
        - ibm_en_subscription_cos
* Support Code Engine
    - **Datasources**
        - ibm_code_engine_app
        - ibm_code_engine_build
        - ibm_code_engine_config_map
        - ibm_code_engine_job
        - ibm_code_engine_job

Enhancements
* suppress resize for autoscaled workerpools for autoscale enabled cluster ([4533](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4533))
* Feature(Image Export Jobs): Support for image export jobs ([4566](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4566))
* feat(Enterprise):added support for destroy method ([4565](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4565))
* ODF Terraform Documentation Updated ([4549](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4549))

BugFixes
* disable creating new resource when endpoint type changes ([4554](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4554))
* retry added for read authorization policy ([4554](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4413))
* Update example for ibm_is_vpc_routing_table ([4564](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4564))
* IBM Cloud Object destination resources and data sources and docs update ([4534](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4534)) 


# 1.53.0(May 04, 2023)
Deprecation
* Added V1 deprecation message for SCC Posture resources and datasources ([4459](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4459))
* Added V1 deprecation message for secret manager datasources ([4523](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4523))

Removed
* The IBM Cloud Certificate Manager service ([4449](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4449))
* Remove resource - is_vpc_route ([4496](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4496))

Features
* Support Direct Link Gateway
    - **Resources**
        - ibm_dl_gateway_action
    - **datasources**
        - ibm_dl_export_route_filters
		- ibm_dl_export_route_filter
		- ibm_dl_import_route_filters
		- ibm_dl_import_route_filter
    
* Support Kubernetes
    - **Resources**
        - ibm_container_ingress_instance
    - **datasources**
        - ibm_container_ingress_instance

* Support Event Notification
    - **Resources**
        - ibm_en_destination_sn
		- ibm_en_subscription_sn
		- ibm_en_destination_ce
		- ibm_en_subscription_ce
    - **datasources**
        - ibm_en_destination_sn
		- ibm_en_subscription_sn
		- ibm_en_destination_ce
		- ibm_en_subscription_ce
* Support Continous Delivery
    - **Resources**
        - ibm_cd_toolchain_tool_eventnotifications
    - **datasources**
        - ibm_cd_toolchain_tool_eventnotifications
* Support VPC
    - **datasources**
        - ibm_is_lb_profile
* Support Code Engine
    - **Resources**
        - ibm_code_engine_secret
    - **datasources**
        - ibm_code_engine_secret

Enhancements
* Update sm_iam_credentials_secret.html.markdown doc ([4468](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4459))
* Adding new resources to support BGP Route Filters in direct link gateway ([4294](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4294))
* Support optional service_group_id param support to policies and roles API ([4455](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4455))
* Add sample Elasticsearch enterprise instance for Docs ([4497](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4497))
* Adding some changes to the object lock documentation ([4486](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4486))
* Support new attributes for secrets manager configuration datasource ([4449](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4449))
* Test cases for IBM PI Instance Update Flow ([4477](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4477))
* Support priority argument for VPC route resource and datasource ([4435](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4435))
* Private DNS support for load balancers ([4463](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4463))
* VSI from existing boot volume ([4433](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4433))
* Support service_endpoints for VPC Virtual Endpoint Gateway resource and datasource ([4514](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4514))
* Support optional attributes and attributes for CD toolchain security compliance resource and datasource ([4525](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4525))
* Support encryption_key attribute in VPC snapshot resource and datasource ([4519](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4519))
* Support network_interface_count attribute for VPC Instance profile datasources ([4527](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4527))
* Support for AddOn parameters ([4408](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4408))

Bugfixes
* VPC load balancer listener idle connection timeout fix ([4482](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4482))
* KMS Instance Policies not setting key_create_import_access values correctly ([4340](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4340))
* Moving an attached public gateway to another resource group fails during apply ([4503](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4503))
* don't throw exception for list TP datsource when there is no TP in the account ([4473](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4473))
* fix(image) : added length check for volume attachments ([4512](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4512))
* support for target crn attribute for VPC subnet reserved resource and datasource ([4513](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4513))
* Support name argument to filter VPC load balancer resource ([4516](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4516))
* fix changes shown on creator attribute for VPC Routing table resource ([4518](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4518))
* ibm_is_bare_metal_server: change to security groups does not trigger change ([4504](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4504))
* [ibm_sm_arbitrary_secret] Allow updating of existing certificate ([4465](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4465))
* Secrets manager documentation fixes ([4530](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4530))
* Unable to deploy the SG rule for ICMP type "8" and code "Any" ([4051](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4051))
* VPC VPN Server Document Correction ([4531](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4531))
* fix: lifecycle of ibm_container_vpc_worker_pool when cluster is deleted ([4539](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4539))
* bug: ibm_is_volume are not attached in sequence to ibm_is_instance ([2390](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2390))
* Added support for s1022 sys type in document update ([4535](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4535))
* Secrets manager documentation fixes ([4547](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4547))
* Documentation correction and ux update ([4529](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4529))

# 1.53.0-beta0(Apr 20, 2023)
Deprecation
* Added V1 deprecation message for SCC Posture resources and datasources ([4459](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4459))
* Added V1 deprecation message for secret manager datasources ([4523](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4523))

Removed
* The IBM Cloud Certificate Manager service ([4449](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4449))
* Remove resource - is_vpc_route ([4496](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4496))

Features
* Support Direct Link Gateway
    - **Resources**
        - ibm_dl_gateway_action
    - **datasources**
        - ibm_dl_export_route_filters
		- ibm_dl_export_route_filter
		- ibm_dl_import_route_filters
		- ibm_dl_import_route_filter
    
* Support Kubernetes
    - **Resources**
        - ibm_container_ingress_instance
    - **datasources**
        - ibm_container_ingress_instance

* Support Event Notification
    - **Resources**
        - ibm_en_destination_sn
		- ibm_en_subscription_sn
		- ibm_en_destination_ce
		- ibm_en_subscription_ce
    - **datasources**
        - ibm_en_destination_sn
		- ibm_en_subscription_sn
		- ibm_en_destination_ce
		- ibm_en_subscription_ce
* Support Continous Delivery
    - **Resources**
        - ibm_cd_toolchain_tool_eventnotifications
    - **datasources**
        - ibm_cd_toolchain_tool_eventnotifications
* Support VPC
    - **datasources**
        - ibm_is_lb_profile
* Support Code Engine
    - **Resources**
        - ibm_code_engine_secret
    - **datasources**
        - ibm_code_engine_secret

Enhancements
* Update sm_iam_credentials_secret.html.markdown doc ([4468](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4459))
* Adding new resources to support BGP Route Filters in direct link gateway ([4294](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4294))
* Support optional service_group_id param support to policies and roles API ([4455](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4455))
* Add sample Elasticsearch enterprise instance for Docs ([4497](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4497))
* Adding some changes to the object lock documentation ([4486](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4486))
* Support new attributes for secrets manager configuration datasource ([4449](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4449))
* Test cases for IBM PI Instance Update Flow ([4477](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4477))
* Support priority argument for VPC route resource and datasource ([4435](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4435))
* Private DNS support for load balancers ([4463](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4463))
* VSI from existing boot volume ([4433](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4433))
* Support service_endpoints for VPC Virtual Endpoint Gateway resource and datasource ([4514](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4514))
* Support optional attributes and attributes for CD toolchain security compliance resource and datasource ([4525](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4525))
* Support encryption_key attribute in VPC snapshot resource and datasource ([4519](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4519))
* Support network_interface_count attribute for VPC Instance profile datasources ([4527](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4527))

Bugfixes
* VPC load balancer listener idle connection timeout fix ([4482](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4482))
* KMS Instance Policies not setting key_create_import_access values correctly ([4340](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4340))
* Moving an attached public gateway to another resource group fails during apply ([4503](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4503))
* don't throw exception for list TP datsource when there is no TP in the account ([4473](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4473))
* fix(image) : added length check for volume attachments ([4512](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4512))
* support for target crn attribute for VPC subnet reserved resource and datasource ([4513](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4513))
* Support name argument to filter VPC load balancer resource ([4516](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4516))
* fix changes shown on creator attribute for VPC Routing table resource ([4518](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4518))
* ibm_is_bare_metal_server: change to security groups does not trigger change ([4504](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4504))
* [ibm_sm_arbitrary_secret] Allow updating of existing certificate ([4465](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4465))
# 1.52.0(Apr 05, 2023)
Features
* Support for Transist Gateway
    - **Resources**
        - ibm_tg_connection_action
* Support for Code Engine
    - **Resources**
        - ibm_code_engine_project
        - ibm_code_engine_app
        - ibm_code_engine_build
        - ibm_code_engine_config_map
        - ibm_code_engine_job
    - **DataSources**
        - ibm_code_engine_project
* Beta support for VPC File Share
    - **Resources**
        - ibm_is_share
        - ibm_is_share_replica_operations
        - ibm_is_share_target
    - **DataSources**
        - ibm_is_share
        - ibm_is_shares
        - ibm_is_source_share
        - ibm_is_share_profile
        - ibm_is_share_profiles
        - ibm_is_share_target
        - ibm_is_share_targets

Enhancements
* Support Idle Connection Timeout for VPC instance and loadbalancer listener ([4399](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4399))
* Support Object Lock feature for COS ([4418](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4418))
* Catalog Management enhancements ([4415](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4415))
* removed the resource: ibm_is_security_group_network_interface_attachment ([4416](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4416))
* enhancement(tpm): support for bare metal secure boot and tpm ([4343](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4343))
* feat(IAM Policy Management): add support for v2/policies API ([4381](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4381))
* Update UKO plugin to use the UKO 4.7 version of the sdk ([4409](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4409))


Bugfixes
* document the CBR custom service endpoint ([4344](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4344))
* Fix to allow memory and processor updates to server instances in either shutdown or active state ([4383](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4383))
* Error trying to read ibm_satellite_endpoint resource ([4057](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4057))
* CD Artifactory parameter description update ([4405](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4405))
* Secrets Manager fixes ([4402](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4402))
* PowerVS start and possibly stop actions failing when resource already in the desired state ([4400](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4400))
* Update pi_network documentation to provide information about default DNS configurations ([4404](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4404))
* Use resource timeout instead of hardcoded value in instance action ([4417](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4417))
* ibm_cos_bucket object versioning block type should be object ([4069](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4069))
* Link from documentation pointing at 404 ([3615](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3615))
* set instance_count in instance group read ([4421](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4421))
* fix(tekton): worker ID fix and nextBuildNumber ([4425](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4425))
* refactor how we set and update the taints on worker pools ([4419](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4419))
* fix: check state for ibm_tg_connection_action correctly ([4446](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4446))
* Code engine docs changes and example update ([4448](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4448))
* Fix regex used for validating a secret group name,Update resource_ibm_sm_secret_group.go ([4451](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4451))
* refactor(cloudant): handle capacity change duration more gracefully ([4414](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4414))
* Added code, tests, and documentation for certificate_root ([4457](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4457))
* Fix #4452: Add 1.3 version of terraform to schematics workspace resource ([4453](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4453))
* fix(tekton): sorting issue for events ([4462](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4462))
* Update platform-services-go-sdk to v0.34.0 and CM fixes ([4474](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4474))

# 1.52.0-beta0(Mar 27, 2023)
Features
* Support for Transist Gateway
    - **Resources**
        - ibm_tg_connection_action
* Support for Code Engine
    - **Resources**
        - ibm_code_engine_project
        - ibm_code_engine_app
        - ibm_code_engine_build
        - ibm_code_engine_config_map
        - ibm_code_engine_job
    - **DataSources**
        - ibm_code_engine_project
* Beta support for VPC File Share
    - **Resources**
        - ibm_is_share
        - ibm_is_share_replica_operations
        - ibm_is_share_target
    - **DataSources**
        - ibm_is_share
        - ibm_is_shares
        - ibm_is_source_share
        - ibm_is_share_profile
        - ibm_is_share_profiles
        - ibm_is_share_target
        - ibm_is_share_targets

Enhancements
* Support Idle Connection Timeout for VPC instance and loadbalancer listener ([4399](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4399))
* Support Object Lock feature for COS ([4418](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4418))
* Catalog Management enhancements ([4415](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4415))
* removed the resource: ibm_is_security_group_network_interface_attachment ([4416](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4416))
* enhancement(tpm): support for bare metal secure boot and tpm ([4343](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4343))
* feat(IAM Policy Management): add support for v2/policies API ([4381](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4381))
* Update UKO plugin to use the UKO 4.7 version of the sdk ([4409](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4409))


Bugfixes
* document the CBR custom service endpoint ([4344](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4344))
* Fix to allow memory and processor updates to server instances in either shutdown or active state ([4383](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4383))
* Error trying to read ibm_satellite_endpoint resource ([4057](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4057))
* CD Artifactory parameter description update ([4405](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4405))
* Secrets Manager fixes ([4402](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4402))
* PowerVS start and possibly stop actions failing when resource already in the desired state ([4400](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4400))
* Update pi_network documentation to provide information about default DNS configurations ([4404](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4404))
* Use resource timeout instead of hardcoded value in instance action ([4417](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4417))
* ibm_cos_bucket object versioning block type should be object ([4069](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4069))
* Link from documentation pointing at 404 ([3615](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3615))
* set instance_count in instance group read ([4421](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4421))
* fix(tekton): worker ID fix and nextBuildNumber ([4425](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4425))
* refactor how we set and update the taints on worker pools ([4419](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4419))

# 1.51.0(Mar 03, 2023)
Features
* Support for Virtual Private Cloud
    - **DataSources**
        - ibm_is_snapshot_clone
        - ibm_is_snapshot_clones
* Beta support for Secrets Manager
    - **DataSources**
        - ibm_sm_secret_group
        - ibm_sm_secret_groups
        - ibm_sm_private_certificate_configuration_intermediate_ca
        - ibm_sm_private_certificate_configuration_root_ca
        - ibm_sm_private_certificate_configuration_template
        - ibm_sm_public_certificate_configuration_ca_lets_encrypt
        - ibm_sm_public_certificate_configuration_dns_cis
        - ibm_sm_public_certificate_configuration_dns_classic_infrastructure
        - ibm_sm_configuration_iam_credentials
        - ibm_sm_secrets
        - ibm_sm_arbitrary_secret_metadata
        - ibm_sm_imported_certificate_metadata
        - ibm_sm_public_certificate_metadata
        - ibm_sm_private_certificate_metadata
        - ibm_sm_iam_credentials_secret_metadata
        - ibm_sm_kv_secret_metadata
        - ibm_sm_username_password_secret_metadata
        - ibm_sm_arbitrary_secret
        - ibm_sm_imported_certificate
        - ibm_sm_public_certificate
        - ibm_sm_private_certificate
        - ibm_sm_iam_credentials_secret
        - ibm_sm_kv_secret
        - ibm_sm_username_password_secret
        - ibm_sm_en_registration
    - **Resources**
        - ibm_sm_secret_group
        - ibm_sm_arbitrary_secret
        - ibm_sm_imported_certificate
        - ibm_sm_public_certificate
        - ibm_sm_private_certificate
        - ibm_sm_iam_credentials_secret
        - ibm_sm_username_password_secret
        - ibm_sm_kv_secret
        - ibm_sm_public_certificate_configuration_ca_lets_encrypt
        - ibm_sm_public_certificate_configuration_dns_cis
        - ibm_sm_public_certificate_configuration_dns_classic_infrastructure
        - ibm_sm_private_certificate_configuration_root_ca
        - ibm_sm_private_certificate_configuration_intermediate_ca
        - ibm_sm_private_certificate_configuration_template
        - ibm_sm_configuration_iam_credentials
        - ibm_sm_en_registration
* Support for Event Notification
    - **Data Sources**
        -ibm_en_destination_msteams
        - ibm_en_subscription_msteams
        - ibm_en_destination_cf
        - ibm_en_subscription_cf
        - ibm_en_destination_pagerduty
        - ibm_en_subscription_pagerduty
        - ibm_en_integration
         - ibm_en_integrations
    - **Resources**
        -ibm_en_destination_msteams
        - ibm_en_subscription_msteams
        - ibm_en_destination_cf
        - ibm_en_subscription_cf
        - ibm_en_destination_pagerduty
        - ibm_en_subscription_pagerduty
        - ibm_en_integration
        - ibm_en_integrations

Enhancements
* Deprecate ibm_private_dns_custom_resolver_location resource ([4296](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4296))
* Add support for single node infrastructure in ibm_satellite_cluster ([4269](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4269))
* Allow data to be set/updated when already unquoted for ibm_cm_object ([4320](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4320))
* Refactor(Cloud Databases): configuration uses cloud-databases-go-sdk ([4234](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4234))
* Enable Context-based Restriction service Private Endpoint support ([4268](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4268))
* Feature: metadata post ga changes in instance and templates ([4352](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4352))
* pass ibm cloud session token along to iaas ([4254](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4254))
* New options for secrets manager tool for ibm_cd_toolchain_tool_secretsmanager ([4350](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4350))
* enhancement(clone): support for backup policy plan clone ([4341](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4341))
* enhancement(clone): support for snapshot clone ([4342](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4342))
* enhancemnet(bare metal): support for z systems in VPC ([4387](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4387))

Bugfixes
* fix: add fields to version metadata for onboarding fix ([4348](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4348))
* Object versioning and bucket replication fix ([4355](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4355))
* fix adding taints for classic kubernetes worker pool ([4377](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4377))
* fix(clone): update on backup policy plan clone policy ([4371](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4371))
* docs: small doc fix for validation resource of cm_object ([4361](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4361))
* fix(doc): vpn_server_route resource doc fix vpn_server_id to vpn_server ([4388](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4388))

# 1.51.0-beta0(Feb 22, 2023)
Features
* Support for Virtual Private Cloud
    - **DataSources**
        - ibm_is_snapshot_clone
        - ibm_is_snapshot_clones
* Beta support for Secrets Manager
    - **DataSources**
        - ibm_sm_secret_group
        - ibm_sm_secret_groups
        - ibm_sm_configuration_private_certificate_intermediate_ca
        - ibm_sm_configuration_private_certificate_root_ca
        - ibm_sm_configuration_private_certificate_template
        - ibm_sm_configuration_public_certificate_ca_lets_encrypt
        - ibm_sm_configuration_public_certificate_dns_cis
        - ibm_sm_configuration_public_certificate_dns_classic_infrastructure
        - ibm_sm_configuration_iam_credentials
        - ibm_sm_secrets
        - ibm_sm_arbitrary_secret_metadata
        - ibm_sm_imported_certificate_metadata
        - ibm_sm_public_certificate_metadata
        - ibm_sm_private_certificate_metadata
        - ibm_sm_iam_credentials_secret_metadata
        - ibm_sm_kv_secret_metadata
        - ibm_sm_username_password_secret_metadata
        - ibm_sm_arbitrary_secret
        - ibm_sm_imported_certificate
        - ibm_sm_public_certificate
        - ibm_sm_private_certificate
        - ibm_sm_iam_credentials_secret
        - ibm_sm_kv_secret
        - ibm_sm_username_password_secret
        - ibm_sm_en_registration
    - **Resources**
        - ibm_sm_secret_group
        - ibm_sm_arbitrary_secret
        - ibm_sm_imported_certificate
        - ibm_sm_public_certificate
        - ibm_sm_private_certificate
        - ibm_sm_iam_credentials_secret
        - ibm_sm_username_password_secret
        - ibm_sm_kv_secret
        - ibm_sm_configuration_public_certificate_ca_lets_encrypt
        - ibm_sm_configuration_public_certificate_dns_cis
        - ibm_sm_configuration_public_certificate_dns_classic_infrastructure
        - ibm_sm_configuration_private_certificate_root_ca
        - ibm_sm_configuration_private_certificate_intermediate_ca
        - ibm_sm_configuration_private_certificate_template
        - ibm_sm_configuration_iam_credentials
        - ibm_sm_en_registration
* Support for Event Notification
    - **Data Sources**
        -ibm_en_destination_msteams
        - ibm_en_subscription_msteams
        - ibm_en_destination_cf
        - ibm_en_subscription_cf
        - ibm_en_destination_pagerduty
        - ibm_en_subscription_pagerduty
        - ibm_en_integration
         - ibm_en_integrations
    - **Resources**
        -ibm_en_destination_msteams
        - ibm_en_subscription_msteams
        - ibm_en_destination_cf
        - ibm_en_subscription_cf
        - ibm_en_destination_pagerduty
        - ibm_en_subscription_pagerduty
        - ibm_en_integration
        - ibm_en_integrations

Enhancements
* Deprecate ibm_private_dns_custom_resolver_location resource ([4296](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4296))
* Add support for single node infrastructure in ibm_satellite_cluster ([4269](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4269))
* Allow data to be set/updated when already unquoted for ibm_cm_object ([4320](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4320))
* Refactor(Cloud Databases): configuration uses cloud-databases-go-sdk ([4234](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4234))
* Enable Context-based Restriction service Private Endpoint support ([4268](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4268))
* Feature: metadata post ga changes in instance and templates ([4352](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4352))
* pass ibm cloud session token along to iaas ([4254](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4254))
* New options for secrets manager tool for ibm_cd_toolchain_tool_secretsmanager ([4350](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4350))
* enhancement(clone): support for backup policy plan clone ([4341](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4341))
* enhancement(clone): support for snapshot clone ([4342](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4342))

Bugfixes
* fix: add fields to version metadata for onboarding fix ([4348](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4348))

# 1.50.0-beta0(Jan 23, 2023)
Features
* Support for Virtual Private Cloud
    - **DataSources**
        - ibm_is_backup_policy_job
        - ibm_is_backup_policy_jobs

Enhancements
* Updates to tekton-pipeline support resources and datasource ([4235](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4235))
* Git PAT support for continous delivery ([4276](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4276))
* made name optional to allow create_before_destroy for VPC network acl ([4289](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4289))
* Remove certificate manager service support ([4290](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4290))

Bugfixes
* add ForceNew tag to OS fields ([4279](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4279))
* warning state in pending list, WARN logs about warning and critical states of cluster ([4283](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4283))
* vpc address prefix and subnet doc fix ([4262](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4262))
* Schematics managing schematics fails after creation ([4132](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4132))
* Fix ibm_appid_cloud_directory_user missing userName ([4284](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4284))

# 1.50.0-beta0(Jan 23, 2023)
Features
* Support for Virtual Private Cloud
    - **DataSources**
        - ibm_is_backup_policy_job
        - ibm_is_backup_policy_jobs

Enhancements
* Updates to tekton-pipeline support resources and datasource ([4235](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4235))
* Git PAT support for continous delivery ([4276](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4276))
* made name optional to allow create_before_destroy for VPC network acl ([4289](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4289))
* Remove certificate manager service support ([4290](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4290))

Bugfixes
* add ForceNew tag to OS fields ([4279](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4279))
* warning state in pending list, WARN logs about warning and critical states of cluster ([4283](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4283))
* vpc address prefix and subnet doc fix ([4262](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4262))
* Schematics managing schematics fails after creation ([4132](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4132))
* Fix ibm_appid_cloud_directory_user missing userName ([4284](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4284))
# 1.49.0(Jan 04, 2023)
Features
* Support for Catalog Management
    - **DataSources**
        - ibm_cm_preset
        - ibm_cm_object
    - **Resources**
        - ibm_cm_object
* Support for Virtual Private Cloud
    - **Resources**
        - ibm_is_instance_network_interface_floating_ip

Enhancements
* Create and Delete Logical Replication Slots for databases-for-postgresql ([4116](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4116))
* Remove bluemix-go dependency for cloud-database allowlist ([4222](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4222))
* Removed usage of direct tags API for retrieving tags using resource CRN ([4209](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4209))
* Support resource_group_id as optional argument for catalogs ([4224](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4224))
* Fix the last operation as per new SDK for resource controller ([4228](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4228))
* Support access_tags for VPC, DedicatedHost, Image, Instance, SSHKey, Network, Instance Group, VPN Server and Bare Metal Server
* Support for transit gateway unbound gre tunnel connections ([4213](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4213))
* Support user_mfa , system_access_token_expiration_in_seconds, system_refresh_token_expiration_in_seconds arguments for account settings ([4221](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4221))
* Added removal notification in docs for ibm_is_security_group_network_interface_attachment ([4232](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4232))
* wait for virtual endpoint gateway to be available after creation ([4206](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4206))
* Support filters in vpc service data sources ([4119](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4119))
* Removed direct call to tagging API to retrieve single resource tags ([4210](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4210))
* Support auto delete for vpc vsi boot volume ([4191](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4191))
* Support VPC Volume creation from Snapshot ([4245](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4245))
* Support for catalog images for VPC instance template ([4249](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4249))

Bugfixes
* Dont set python for IBM Satellite host attachment ([4226](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4226))
* [ibm_container_vpc_cluster ] wait_till = "Normal" does not work ([4214](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4214))
* fix doc for ibmcloud cli command for Schematics ([4243](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4243))
* change from fixed to to append new line to slice ([4239](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4239))
* Added security group target list for supported resources ([4247](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4247))
* Fix volume not found error while destroying ibm_pi_volume resource ([4252](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4252))
* PowerVS VSI update in-place with 'ibm_pi_instance' times out when VSI status is SHUTDOWN ([4258](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4258))
* Support Terraform 1.2 verison in schematics resources ([4238](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4238))
# 1.49.0(Dec 19, 2022)
Features
* Support for Catalog Management
    - **DataSources**
        - ibm_cm_preset
        - ibm_cm_object
    - **Resources**
        - ibm_cm_object
* Support for Virtual Private Cloud
    - **Resources**
        - ibm_is_instance_network_interface_floating_ip

Enhancements
* Create and Delete Logical Replication Slots for databases-for-postgresql ([4116](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4116))
* Remove bluemix-go dependency for cloud-database allowlist ([4222](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4222))
* Removed usage of direct tags API for retrieving tags using resource CRN ([4209](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4209))
* Support resource_group_id as optional argument for catalogs ([4224](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4224))
* Fix the last operation as per new SDK for resource controller ([4228](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4228))
* Support access_tags for VPC, DedicatedHost, Image, Instance, SSHKey, Network, Instance Group, VPN Server and Bare Metal Server
* Support for transit gateway unbound gre tunnel connections ([4213](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4213))
* Support user_mfa , system_access_token_expiration_in_seconds, system_refresh_token_expiration_in_seconds arguments for account settings ([4221](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4221))
* Added removal notification in docs for ibm_is_security_group_network_interface_attachment ([4232](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4232))
* wait for virtual endpoint gateway to be available after creation ([4206](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4206))
* Support filters in vpc service data sources ([4119](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4119))
* Removed direct call to tagging API to retrieve single resource tags ([4210](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4210))
* Support auto delete for vpc vsi boot volume ([4191](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4191))
* Support VPC Volume creation from Snapshot ([4245](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4245))
* Support for catalog images for VPC instance template ([4249](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4249))

Bugfixes
* Dont set python for IBM Satellite host attachment ([4226](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4226))
* [ibm_container_vpc_cluster ] wait_till = "Normal" does not work ([4214](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4214))
* fix doc for ibmcloud cli command for Schematics ([4243](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4243))
* change from fixed to to append new line to slice ([4239](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4239))
* Added security group target list for supported resources ([4247](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4247))

# 1.48.0(Dec 01, 2022)
Features
* Support for Powervs
    - **DataSources**
        - ibm_pi_disaster_recovery_location
        - ibm_pi_disaster_recovery_locations
        - ibm_pi_volume_group
        - ibm_pi_volume_groups
        - ibm_pi_volume_group_details
        - ibm_pi_volume_groups_details
        - ibm_pi_volume_group_storage_details
        - ibm_pi_volume_group_remote_copy_relationships
        - ibm_pi_volume_flash_copy_mappings
        - ibm_pi_volume_remote_copy_relationship
        - ibm_pi_volume_onboardings
        - ibm_pi_volume_onboarding
    - **Resources**
        - ibm_pi_volume_onboarding
        - ibm_pi_volume_group
        - ibm_pi_volume_group_action
* Support for KMS
    - **DataSources**
        - ibm_kms_instance_policies
    - **Resources**
        - ibm_kms_instance_policies
        - ibm_kms_key_with_policy_overrides

Enhancements
* Enable cross account KMS boot volume encryption for VPC Cluster ([4128](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4128))
* Deprecate Whitelist for IBM-cloud-databases and Introduce allowlisting ([3852](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3852)) 
* Support wait_till for cluster provisioning with normal state ([4139](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4139)) 
* Disable/Enable Rotation Policy support for kms key ([4110](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4110)) 
* Support operating_system argument for cluster workerpool creation ([4133](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4133)) 
* Update Catalog Management resources and datasoruces with latest API changes ([4126](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4126))
* Added support for Event streams targets ([4161](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4161))
* Support Public Ingress Routing for VPC routing table routes ([4157](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4157))
* Support VPC network acl before rule patch ([4136](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4136))
* Improvements to the Continuous Delivery resources and datasources for the ga release ([4145](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4145))
* feat(Cloud Databases): Allow users to edit configuration on database creation ([4186](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4186))
* Support AS prepend `specific_prefixes` for Direct Link ([4179](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4179))
* LoadBalancerPool HealthMonitor Port Nullable ([4129](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4129))
* Access Tags for VPC Volumes ([4127](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4127))
* Add tagging and access tags support for is_snapshot ([4134](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4134))
* Remove credential passphrase and group support from scc ([4140](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4140))

Bugfixes
* Fix endpoint URL via the env variable for KMS resources/datasources ([4120](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4120))
* Fix docs typo for cloud databases ([4143](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4143))
* Firewall rules Paused is not reflected only at the first execution ([4142](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4142))
* set subscription manager release and disable eus for Satellite host script ([4175](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4175)) 
* Wait for scaling task to complete ([4188](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4188)) 
* Update worker_pool host_labels usage and fix blank entry in array ([4189](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4189)) 
* ibm_resource_instance data source should not need resource group id ([4137](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4137)) 
* private endpoint in data resource ibm_secrets_manager_secret is incorrect ([4187](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4187)) 
* Fix: fixed pi_volume_group_action, pi_volume_group & data pi_volume_group_storage_details resources ([4193](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4193))
* Fix the provision failure of lbass ([4197](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4197))
# 1.48.0-beta0(Nov 16, 2022)
Features
* Support for Powervs
    - **DataSources**
        - ibm_pi_disaster_recovery_location
        - ibm_pi_disaster_recovery_locations
        - ibm_pi_volume_group
        - ibm_pi_volume_groups
        - ibm_pi_volume_group_details
        - ibm_pi_volume_groups_details
        - ibm_pi_volume_group_storage_details
        - ibm_pi_volume_group_remote_copy_relationships
        - ibm_pi_volume_flash_copy_mappings
        - ibm_pi_volume_remote_copy_relationship
        - ibm_pi_volume_onboardings
        - ibm_pi_volume_onboarding
    - **Resources**
        - ibm_pi_volume_onboarding
        - ibm_pi_volume_group
        - ibm_pi_volume_group_action
* Support for KMS
    - **DataSources**
        - ibm_kms_instance_policies
    - **Resources**
        - ibm_kms_instance_policies
        - ibm_kms_key_with_policy_overrides

Enhancements
* Enable cross account KMS boot volume encryption for VPC Cluster ([4128](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4128))
* Deprecate Whitelist for IBM-cloud-databases and Introduce allowlisting ([3852](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3852)) 
* Support wait_till for cluster provisioning with normal state ([4139](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4139)) 
* Disable/Enable Rotation Policy support for kms key ([4110](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4110)) 
* Support operating_system argument for cluster workerpool creation ([4133](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4133)) 
* Update Catalog Management resources and datasoruces with latest API changes ([4126](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4126))
* Added support for Event streams targets ([4161](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4161))
* Support Public Ingress Routing for VPC routing table routes ([4157](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4157))
* Support VPC network acl before rule patch ([4136](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4136))
* Improvements to the Continuous Delivery resources and datasources for the ga release ([4145](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4145))

Bugfixes
* Fix endpoint URL via the env variable for KMS resources/datasources ([4120](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4120))
* Fix docs typo for cloud databases ([4143](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4143))
* Firewall rules Paused is not reflected only at the first execution ([4142](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4142))
* set subscription manager release and disable eus for Satellite host script ([4175](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4175)) 



# 1.47.1(Nov 08, 2022)
BUG FIXES
* Support to add retry fetch of tags using global search API's when we reach ratelimit ([4125](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4125))
* Fix len out of range for IBM Satellite Cluster host labels ([4149](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4149))
# 1.47.0(Nov 02, 2022)
Features
* Support for Kubernetes Service
    - **Resources**
        - ibm_container_vpc_worker
* Support for APP Configuration
    - **DataSources**
        - ibm_app_config_snapshot
        - ibm_app_config_snapshots
    - **Resources**
        - ibm_app_config_snapshot
* Support for DirectLink
    - **DataSources**
        - ibm_dl_route_reports
        - ibm_dl_route_report
    - **Resources**
        - ibm_dl_route_report
Enhancements
* Update schema validation to latest schema methods ([4068](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4068))
* Support Private end-point for APP Configuration ([4048](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4048))
* supported CoreOS-enabled clusters (default worker pool) and workerPools in satellite clusters ([3985](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3985))
* support cos one rate plan ([4092](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4092))
* Fix GRE tunnel in cloud connection ([4093](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4093))
* support floating_bare_metal_server attribute for resource ibm_is_bare_metal_server_network_interface_allow_float ([4115](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4115))
* added support for catalog images for enterprises ([3994](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3994))
* Support rhel8 in attach script ([4033](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4033))

BUG FIXES
* data source: ibm_is_backup_policies, expected 'the empty list', but returned 'no BackupPolicies found' ([4079](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4079))
* Add fix for storage capacity & system pool output ([4074](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4074))
* fixes the endpoint; updates go sdk for Atracker ([4073](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4073))
* fix(vpn_gateway_connection): adding support to null patch ike_policy and ipsec_policy on vpn_gateway_connection ([4058](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4058))
* fix(bare_metal_server_network_interface_allow_float) : reordered the wait logic in bare metal nics ([4101](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4101))
* Update VNF Scalability NIC in docs ([4047](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4047))
* fix(iam-access-groups): added retry logic for access groups read ([4098](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4098))
* ibm_satellite_attach_host_script shows incorrect syntax for labels ([4099](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4099))
* fix(bm-nic): added check for 0.0.0.0 reserved ip on nic availability ([4122](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4122))

# 1.47.0-beta3 (Oct 26, 2022)
BUG FIXES
* Fix attach script loop for satellite host attach ([4117](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4117))
* fix(bm-nic): added check for 0.0.0.0 reserved ip on nic availability ([4122](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4122))
# 1.47.0-beta2 (Oct 18, 2022)
Features
* Support for DirectLink
    - **DataSources**
        - ibm_dl_route_reports
        - ibm_dl_route_report
    - **Resource**
        - ibm_dl_route_report
        
Enhancements
* supported CoreOS-enabled clusters (default worker pool) and workerPools in satellite clusters ([3985](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3985))
* support cos one rate plan ([4092](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4092))
* Fix GRE tunnel in cloud connection ([4093](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4093))
* support floating_bare_metal_server attribute for resource ibm_is_bare_metal_server_network_interface_allow_float ([4115](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4115))
* added support for catalog images for enterprises ([3994](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3994))

BUG FIXES
* data source: ibm_is_backup_policies, expected 'the empty list', but returned 'no BackupPolicies found' ([4079](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4079))
* Add fix for storage capacity & system pool output ([4074](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4074))
* fixes the endpoint; updates go sdk for Atracker ([4073](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4073))
* fix(vpn_gateway_connection): adding support to null patch ike_policy and ipsec_policy on vpn_gateway_connection ([4058](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4058))
* fix(bare_metal_server_network_interface_allow_float) : reordered the wait logic in bare metal nics ([4101](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4101))
* Update VNF Scalability NIC in docs ([4047](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4047))
* fix(iam-access-groups): added retry logic for access groups read ([4098](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4098))
* ibm_satellite_attach_host_script shows incorrect syntax for labels ([4099](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4099))

# 1.47.0-beta0 (Oct 10, 2022)
Features
* Support for Kubernetes Service
    - **Resource**
        - ibm_container_vpc_worker
* Support for APP Configuration
    - **DataSources**
        - ibm_app_config_snapshot
        - ibm_app_config_snapshots
    - **Resources**
        - ibm_app_config_snapshot
Enhancements
* Update schema validation to latest schema methods ([4068](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4068))
* Support Private end-point for APP Configuration ([4048](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4048))

# 1.46.0 (Oct 03, 2022)
Features
* Support for Power Instance
    - **DataSources**
        - ibm_pi_instance_action
        - ibm_pi_shared_processor_pool
        - ibm_pi_shared_processor_pools
        - ibm_pi_spp_placement_group
        - ibm_pi_spp_placement_groups
     - **Resources**
        - ibm_pi_shared_processor_pool
        - ibm_pi_spp_placement_group
* Support Security and Compilance 
    - **DataSources**
        - ibm_scc_posture_profile_import
        - ibm_scc_posture_scan_initiate_validation
    - **Resources**
        - ibm_scc_posture_scan_initiate_validation
* Support App Configuration
    - **DataSources**
        - ibm_app_config_collection
        - ibm_app_config_collections
        - ibm_app_config_property
        - ibm_app_config_properties
    - **Resources**
        - ibm_app_config_collection
        - ibm_app_config_collection
* Support Virtual Private Cloud
    - **DataSources**
        - ibm_is_instance_groups
        - ibm_is_bare_metal_server_network_interface_reserved_ip
        - ibm_is_bare_metal_server_network_interface_reserved_ips
        
Enhancements
* Update CD Toolchain resources and datasources with latest SDK ([3933](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3933)) 
* Support uset tags in volumes of instance, instance template and volume attachement ([3993](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3993)) 
* Add replication enabled attribute in pi volume resource ([4007](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4007)) 
* Update Subnets on LoadBalancer ([4026](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4026)) 
* Support for lifecycle status and reasons in VPC instance ([4017](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4017)) 
* Support VPC reference in VPN gateway ([4039](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4039)) 
* Examples and docs updated for HPCS COS support ([4034](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4034)) 
* Deprecate connection strings for IBM-cloud-databases ([4050](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4050)) 
* Revert volume replication attribute in pi volume resource ([4059](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4059))
* added support for new ciphers ([4018](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4018))

BUG FIXES
* terraform ibm_is_instance_template volume_attachments removal from configuration not changing resource ([3972](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3972)) 
* docs in is_instance_template do not cover volume_attachments completely ([3967](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3967)) 
* Note section is not in proper format ([3970](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3970)) 
* Wrong tagging for TOC ([4019](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4019)) 
* ibm_kms_key unhelpful deprecation box ([3923](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3923)) 
* "data iam_authorization_policy" is "data iam_authorization_policies" ([4015](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4015)) 
* ibm_scc_rule vs ibm_scc_configuration_rule ([3908](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3908)) 
* ibm_cos_bucket.abort_incomplete_multipart_upload_days is undocumented ([3799](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3799)) 
* Updated Description for Deployment type EPIC ([4037](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4037)) 
* Doc corrections for resource group to intimate id and subnet floating ip changes ([4025](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4025)) 
* fix(ibm_is_vpn_gateway): reordered the setting of id in ibm_is_vpn_gateway resource to taint the resource properly ([4055](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4055))
* fix(docs): added details about user data and update link for ibmcloud docs ([4066](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4066))
# 1.46.0-beta0 (Sep 19, 2022)
Features
* Support for Power Instance
    - **DataSources**
        - ibm_pi_instance_action
        - ibm_pi_shared_processor_pool
        - ibm_pi_shared_processor_pools
        - ibm_pi_spp_placement_group
        - ibm_pi_spp_placement_groups
     - **Resources**
        - ibm_pi_shared_processor_pool
        - ibm_pi_shared_processor_pool
* Support Security and Compilance 
    - **DataSources**
        - ibm_scc_posture_profile_import
        - ibm_scc_posture_scan_initiate_validation
    - **Resources**
        - ibm_scc_posture_scan_initiate_validation
* Support App Configuration
    - **DataSources**
        - ibm_app_config_collection
        - ibm_app_config_collections
        - ibm_app_config_property
        - ibm_app_config_properties
    - **Resources**
        - ibm_app_config_collection
        - ibm_app_config_collection
* Support Virtual Private Cloud
    - **DataSources**
        - ibm_is_instance_groups
        - ibm_is_bare_metal_server_network_interface_reserved_ip
        - ibm_is_bare_metal_server_network_interface_reserved_ips
        
Enhancements
* Update CD Toolchain resources and datasources with latest SDK ([3933](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3933)) 
* Support uset tags in volumes of instance, instance template and volume attachement ([3993](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3993)) 
* Add replication enabled attribute in pi volume resource ([4007](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4007)) 
* Update Subnets on LoadBalancer ([4026](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4026)) 
* Support for lifecycle status and reasons in VPC instance ([4017](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4017)) 
* Support VPC reference in VPN gateway ([4039](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4039)) 
* Examples and docs updated for HPCS COS support ([4034](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4034)) 

BUG FIXES
* terraform ibm_is_instance_template volume_attachments removal from configuration not changing resource ([3972](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3972)) 
* docs in is_instance_template do not cover volume_attachments completely ([3967](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3967)) 
* Note section is not in proper format ([3970](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3970)) 
* Wrong tagging for TOC ([4019](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4019)) 
* ibm_kms_key unhelpful deprecation box ([3923](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3923)) 
* "data iam_authorization_policy" is "data iam_authorization_policies" ([4015](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/4015)) 
* ibm_scc_rule vs ibm_scc_configuration_rule ([3908](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3908)) 
* ibm_cos_bucket.abort_incomplete_multipart_upload_days is undocumented ([3799](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3799)) 
* Updated Description for Deployment type EPIC ([4037](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4037)) 
* Doc corrections for resource group to intimate id and subnet floating ip changes ([4025](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4025)) 

# 1.45.1 (Sep 14, 2022)
Enhancements
* Support location argument to target the respective Schematics region service ([4030](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4030))

# 1.45.0 (Sep 05, 2022)
Features
* Support App Configuration
    - **DataSources**
        - ibm_app_config_segment
        - ibm_app_config_segment
    - **Resources**
        - ibm_app_config_segment
* Support IAM Access Group
    - **Resources**
        - ibm_iam_access_group_account_settings

Enhancements
* CIS CNAME Setup : Create Partial Zone ([3937](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3937))
* deprecating Security Insights ([3755](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3755))
* Add SNAP enabled bool for DHCP create function ([3932](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3932))
* Support for EPIC offering Create flow ([3949](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3949))
* Add minimum role notes ([3962](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3962))
* enhancement(is_vpc): added support for identifier in vpc data source ([3959](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3959))
* Atracker v2 metadata backup ([3887](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3887))
* coreos host attach support ([3968](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3968))
* add CBR rule API type support ([3971](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3971))
* Mark key_protect_key_id as deprecate in database datasource ([3939](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3939))
* doc changes for VPC Load Balancer ([3872](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3872))
* Add source units required parameter ([3991](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3991))
* Add default_network_acl computed attribute for VPC ([3997](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3997))


BUG FIXES
* added log entry for redacted credentials ([3942](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3942))
* stack trace with ibm_pi_dhcps data resource ([3951](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3951)) 
* Fix crash on schematics_action resource ([3969](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3969)) 
* Fix duplicate VPC entries in cloud connection create
operation ([3958](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3958)) 
* fix: account_id issue when service id apikey is used ([3950](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3950))
* Delete Over Count Bug Fix ([3946](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3946))
* Fix the setting of monitoring if it fails with error during update ([3974](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3974))
* CIS - Delete filter on deletion of Firewall Rules ([3963](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3963))
* Fix VPN Server clientAuthentication ([3947](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3947))
* fix(security_group): added wait logic to wait for target removal to avoid 409 ([3957](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3957))
* The IBM database fails with Error: Unprocessable Entity ([3964](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3964))
* Prevent runtime error on 0 member allocations ([3992](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3992))
* key_protect_key_id doc typo fix for database ([3980](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3980))
* Fix links in database docs ([3977](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3977))
* fix(ibm_is_security_group_target): missing set statements in sg target resource ([4002](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4002))
* target_http_status_code correction for VPC LoadBalancer listener ([4005](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4005))
* fix the diff on cis tags if provisioned from Schematics ([4008](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/4008))
* CIS Firewall Rules : added priority key ([3998](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3998))
# 1.44.3 (Aug 29, 2022)
BUG FIXES
* The IBM database fails with Error: Unprocessable Entity ([3964](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3964))
* Prevent runtime error on 0 member allocations ([3992](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3992))

# 1.45.0-beta0 (Aug 18, 2022)
Features
* Support App Configuration
    - **DataSources**
        - ibm_app_config_segment
        - ibm_app_config_segment
    - **Resources**
        - ibm_app_config_segment
* Support IAM Access Group
    - **Resources**
        - ibm_iam_access_group_account_settings

Enhancements
* CIS CNAME Setup : Create Partial Zone ([3937](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3937))
* deprecating Security Insights ([3755](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3755))
* Add SNAP enabled bool for DHCP create function ([3932](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3932))
* Support for EPIC offering Create flow ([3949](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3949))
* Add minimum role notes ([3962](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3962))
* enhancement(is_vpc): added support for identifier in vpc data source ([3959](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3959))
* Atracker v2 metadata backup ([3887](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3887))
* coreos host attach support ([3968](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3968))
* add CBR rule API type support ([3971](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3971))
* Mark key_protect_key_id as deprecate in database datasource ([3939](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3939))
* doc changes for VPC Load Balancer ([3872](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3872))

BUG FIXES
* added log entry for redacted credentials ([3942](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3942))
* stack trace with ibm_pi_dhcps data resource ([3951](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3951)) 
* Fix crash on schematics_action resource ([3969](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3969)) 
* Fix duplicate VPC entries in cloud connection create
operation ([3958](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3958)) 
* fix: account_id issue when service id apikey is used ([3950](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3950))
* Delete Over Count Bug Fix ([3946](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3946))
* Fix the setting of monitoring if it fails with error during update ([3974](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3974))
* CIS - Delete filter on deletion of Firewall Rules ([3963](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3963))
* Fix VPN Server clientAuthentication ([3947](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3947))
* fix(security_group): added wait logic to wait for target removal to avoid 409 ([3957](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3957))

# 1.44.2 (Aug 8, 2022)
BUG FIXES
* fix: schematics template metadata issue ([3953](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3953))

# 1.44.1 (Aug 3, 2022)
BUG FIXES
* Fix the breaking change of cos bucket with replication configuration ([3945](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3945))
* fix the firewall allowed_ip issue ([3894](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3894))
* Fix the nil pointer for coreos enabled ([3948](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3948))

# 1.44.0 (Aug 2, 2022)
Features
* Support Internt Services
    - **DataSources**
        - ibm_cis_mtlss
        - ibm_cis_mtls_apps
        - ibm_cis_origin_auths
    - **Resources**
        - ibm_cis_mtls
        - ibm_cis_mtls_app
        - ibm_cis_origin_auth
* Support Internt Services
    - **DataSources**
        - ibm_database_task
        - ibm_database_tasks
* Support Security and Compilance 
    - **DataSources**
        - ibm_scc_posture_credential
        - ibm_scc_posture_collector
        - ibm_scc_posture_scope
        - ibm_scc_posture_credentials
        - ibm_scc_posture_collectors
* Support Virtual Private Cloud
    - **DataSources**
        - ibm_is_backup_policy
        - ibm_is_backup_policies
        - ibm_is_backup_policy_plan
        - ibm_is_backup_policy_plans
        - ibm_is_vpn_server
        - ibm_is_vpn_servers
        - ibm_is_vpn_server_client
        - ibm_is_vpn_server_client_configuration
        - ibm_is_vpn_server_clients
        - ibm_is_vpn_server_route
        - ibm_is_vpn_server_routes
    - **Resources**
        - ibm_is_backup_policy
        - ibm_is_backup_policy_plan
        - ibm_is_vpn_server
        - ibm_is_vpn_server_route
* Support IBM Cloud Storage 
    - **DataSources**
        - ibm_cos_bucket_replication_rule
* Support Private DNS
    - **DataSources**
        - ibm_dns_custom_resolver_secondary_zones
    - **Resources**
        - ibm_dns_custom_resolver_secondary_zones

Enhancements
* Routing Table in Subnet Datasource ([3909](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3909))
* Added cidr param to dhcp ([3916](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3916))
* Add CoreOS-enabled option for Satellite location create ([3914](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3914))

BUG FIXES
* fix(documentation): error statement fix on security group and is_volumes ([3993](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3893))
* CIS WAF Group Documentation Fix ([3997](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3897))
* CD doc fix from generator utility ([3905](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3905))
* added note for default cipher values ([3915](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3915))
* Deprecation of vpc_route resource & correction in vpc_routing_table_route ([3919](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3919))
* Roll back the strogae types for existing buckets ([3928](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3928))
* Fix the subcategory for cbr docs ([3928](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3928))
* terraform versions 1.1 added to validator check ([3930](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3920))
* fix(is_bare_metal_server): removed nics with allow float from server nics ([3938](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3928))
# 1.44.0-beta0 (Jul 20, 2022)
Features
* Support Internt Services
    - **DataSources**
        - ibm_cis_mtlss
        - ibm_cis_mtls_apps
        - ibm_cis_origin_auths
    - **Resources**
        - ibm_cis_mtls
        - ibm_cis_mtls_app
        - ibm_cis_origin_auth
* Support Internt Services
    - **DataSources**
        - ibm_database_task
        - ibm_database_tasks
* Support Security and Compilance 
    - **DataSources**
        - ibm_scc_posture_credential
        - ibm_scc_posture_collector
        - ibm_scc_posture_scope
        - ibm_scc_posture_credentials
        - ibm_scc_posture_collectors
* Support Virtual Private Cloud
    - **DataSources**
        - ibm_is_backup_policy
        - ibm_is_backup_policies
        - ibm_is_backup_policy_plan
        - ibm_is_backup_policy_plans
        - ibm_is_vpn_server
        - ibm_is_vpn_servers
        - ibm_is_vpn_server_client
        - ibm_is_vpn_server_client_configuration
        - ibm_is_vpn_server_clients
        - ibm_is_vpn_server_route
        - ibm_is_vpn_server_routes
    - **Resources**
        - ibm_is_backup_policy
        - ibm_is_backup_policy_plan
        - ibm_is_vpn_server
        - ibm_is_vpn_server_route
* Support IBM Cloud Storage 
    - **DataSources**
        - ibm_cos_bucket_replication_rule
* Support Private DNS
    - **DataSources**
        - ibm_dns_custom_resolver_secondary_zones
    - **Resources**
        - ibm_dns_custom_resolver_secondary_zones

Enhancements
* Routing Table in Subnet Datasource ([3909](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3909))
* Added cidr param to dhcp ([3916](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3916))
* Add CoreOS-enabled option for Satellite location create ([3914](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3914))

BUG FIXES
* fix(documentation): error statement fix on security group and is_volumes ([3993](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3893))
* CIS WAF Group Documentation Fix ([3997](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3897))
* CD doc fix from generator utility ([3905](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3905))
* added note for default cipher values ([3915](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3915))
* Deprecation of vpc_route resource & correction in vpc_routing_table_route ([3919](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3919))
* Roll back the strogae types for existing buckets ([3928](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3928))
* Fix the subcategory for cbr docs ([3928](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3928))

# 1.43.0 (Jul 01, 2022)
Features
* Support Kubernetes
    - **DataSources**
        - ibm_container_dedicated_host_pool
        - ibm_container_dedicated_host_flavor
        - ibm_container_dedicated_host_flavors
        - ibm_container_dedicated_host
    - **Resources**
        - ibm_container_dedicated_host_pool
        - ibm_container_dedicated_host

* Support EventNotification
    - **DataSources**
        - ibm_en_source
        - ibm_en_destination_slack
        - ibm_en_subscription_slack
        - ibm_en_subscription_safari
        - ibm_en_destination_safari
    - **Resources**
        - ibm_en_source
        - ibm_en_destination_slack
        - ibm_en_subscription_slack
        - ibm_en_subscription_safari
        - ibm_en_destination_safari

* Support HPCS
    - **DataSources**
        - ibm_hpcs_managed_key
        - ibm_hpcs_key_template
        - ibm_hpcs_keystore
        - ibm_hpcs_vault
    - **Resources**
        - ibm_hpcs_managed_key
        - ibm_hpcs_key_template
        - ibm_hpcs_keystore
        - ibm_hpcs_vault

* Support CD Toolchain
    - **DataSources**
        - ibm_cd_toolchain
		- ibm_cd_toolchain_tool_keyprotect         
		- ibm_cd_toolchain_tool_secretsmanager    
		- ibm_cd_toolchain_tool_bitbucketgit       
		- ibm_cd_toolchain_tool_githubintegrated   
		- ibm_cd_toolchain_tool_githubconsolidated 
		- ibm_cd_toolchain_tool_gitlab 
		- ibm_cd_toolchain_tool_hostedgit          
		- ibm_cd_toolchain_tool_artifactory        
		- ibm_cd_toolchain_tool_custom 
		- ibm_cd_toolchain_tool_pipeline          
		- ibm_cd_toolchain_tool_devopsinsights     
		- ibm_cd_toolchain_tool_slack  
		- ibm_cd_toolchain_tool_sonarqube         
		- ibm_cd_toolchain_tool_hashicorpvault     
		- ibm_cd_toolchain_tool_securitycompliance 
		- ibm_cd_toolchain_tool_privateworker     
		- ibm_cd_toolchain_tool_appconfig          
		- ibm_cd_toolchain_tool_jenkins
		- ibm_cd_toolchain_tool_nexus  
		- ibm_cd_toolchain_tool_pagerduty          
		- ibm_cd_toolchain_tool_saucelabs          
    - **Resources**
        - ibm_cd_toolchain
		- ibm_cd_toolchain_tool_keyprotect         
		- ibm_cd_toolchain_tool_secretsmanager    
		- ibm_cd_toolchain_tool_bitbucketgit       
		- ibm_cd_toolchain_tool_githubintegrated   
		- ibm_cd_toolchain_tool_githubconsolidated 
		- ibm_cd_toolchain_tool_gitlab 
		- ibm_cd_toolchain_tool_hostedgit          
		- ibm_cd_toolchain_tool_artifactory        
		- ibm_cd_toolchain_tool_custom 
		- ibm_cd_toolchain_tool_pipeline          
		- ibm_cd_toolchain_tool_devopsinsights     
		- ibm_cd_toolchain_tool_slack  
		- ibm_cd_toolchain_tool_sonarqube         
		- ibm_cd_toolchain_tool_hashicorpvault     
		- ibm_cd_toolchain_tool_securitycompliance 
		- ibm_cd_toolchain_tool_privateworker     
		- ibm_cd_toolchain_tool_appconfig          
		- ibm_cd_toolchain_tool_jenkins
		- ibm_cd_toolchain_tool_nexus  
		- ibm_cd_toolchain_tool_pagerduty          
		- ibm_cd_toolchain_tool_saucelabs

* Support CD Tekton Pipeline
    - **Datasources**
        - ibm_cd_tekton_pipeline_definition
        - ibm_cd_tekton_pipeline_trigger_property
        - ibm_cd_tekton_pipeline_property
        - ibm_cd_tekton_pipeline_trigger
        - ibm_cd_tekton_pipeline
    - **Resources**
        - ibm_cd_tekton_pipeline_definition
        - ibm_cd_tekton_pipeline_trigger_property
        - ibm_cd_tekton_pipeline_property
        - ibm_cd_tekton_pipeline_trigger
        - ibm_cd_tekton_pipeline
        
Enhancements
* support for creation of bucket on Satellite location ([3727](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3727))
* Allow status filtering in ibm_is_images ([3841](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3841))
* IBM Cloud CD terraform resources examples - empty toolchain template ([3858](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3858))
* updated platform-services-go-sdk version to 0.26.1 ([3881](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3881))
* placement target patch support ([3809](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3809))
* Added multithreading for 49 requests for Private DNS Resource Record ([3886](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3886))
* add vpc cluster and workerpool to the dhost example ([3878](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3878))
* CBR: context-based-restriction update for Enforcement mode support ([3853](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3853))
* profile patch enhancement for VPC instance ([3860](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3860))

BUG Fixes
* ibm_dns_glb docs to not specify the enabled argument ([3818](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3818))
* floating ip data source nil fix ([3843](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3843))
* support TLS 1.3 supported ciphers for ibm_cis_domain_settings ([3736](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3736))
* Update provider version in docs ([3750](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3750))
* Update cloud_shell_account_settings.html.markdown ([3812](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3812))
* fix: placement targets empty list ([3787](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3787))
* fix: instance network interface mutex sync ([3791](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3791))
* ibm_scc_posture_scope does not allow for blank-space in description ([3733](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3733))
* SCC provider should output User friendly error messages clearly indicating the problem ([3844](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3844))
* ibm_is_instance should not not suppress change and force new on boot_volume.0.snapshot change ([3819](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3819))
* fix(bare_metal_server) : vlan id in multi nic fix ([3767](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3767))
* fix ibm_container_storage_attachment import ([3862](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3862))
* Remove rogue CD website document ([3865](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3865))
* IAM User inivte always shows a diff on "invited_users" ([3863](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3863))
* name corrections for VPC instance group manager docs ([3884](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3884))
* UKO Terraform Doc and Bugfixes ([3889](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3889))
* Make CD Pagerduty prop computed and general tidy up ([3891](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3891))
* cis waf group settings fix ([3882](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3882))
* Fix: validator identifier names ([3870](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3870))
* Fix validateSchema from one element to zero ([3870](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3870))
# 1.43.0-beta0 (Jun 22, 2022)
Features
* Support Kubernetes
    - **DataSources**
        - ibm_container_dedicated_host_pool
        - ibm_container_dedicated_host_flavor
        - ibm_container_dedicated_host_flavors
        - ibm_container_dedicated_host
    - **Resources**
        - ibm_container_dedicated_host_pool
        - ibm_container_dedicated_host

* Support EventNotification
    - **DataSources**
        - ibm_en_source
        - ibm_en_destination_slack
        - ibm_en_subscription_slack
        - ibm_en_subscription_safari
        - ibm_en_destination_safari
    - **Resources**
        - ibm_en_source
        - ibm_en_destination_slack
        - ibm_en_subscription_slack
        - ibm_en_subscription_safari
        - ibm_en_destination_safari

* Support HPCS
    - **DataSources**
        - ibm_hpcs_managed_key
        - ibm_hpcs_key_template
        - ibm_hpcs_keystore
        - ibm_hpcs_vault
    - **Resources**
        - ibm_hpcs_managed_key
        - ibm_hpcs_key_template
        - ibm_hpcs_keystore
        - ibm_hpcs_vault

* Support CD Toolchain
    - **DataSources**
        - ibm_cd_toolchain
		- ibm_cd_toolchain_tool_keyprotect         
		- ibm_cd_toolchain_tool_secretsmanager    
		- ibm_cd_toolchain_tool_bitbucketgit       
		- ibm_cd_toolchain_tool_githubintegrated   
		- ibm_cd_toolchain_tool_githubconsolidated 
		- ibm_cd_toolchain_tool_gitlab 
		- ibm_cd_toolchain_tool_hostedgit          
		- ibm_cd_toolchain_tool_artifactory        
		- ibm_cd_toolchain_tool_custom 
		- ibm_cd_toolchain_tool_pipeline          
		- ibm_cd_toolchain_tool_devopsinsights     
		- ibm_cd_toolchain_tool_slack  
		- ibm_cd_toolchain_tool_sonarqube         
		- ibm_cd_toolchain_tool_hashicorpvault     
		- ibm_cd_toolchain_tool_securitycompliance 
		- ibm_cd_toolchain_tool_privateworker     
		- ibm_cd_toolchain_tool_appconfig          
		- ibm_cd_toolchain_tool_jenkins
		- ibm_cd_toolchain_tool_nexus  
		- ibm_cd_toolchain_tool_pagerduty          
		- ibm_cd_toolchain_tool_saucelabs          
    - **Resources**
        - ibm_cd_toolchain
		- ibm_cd_toolchain_tool_keyprotect         
		- ibm_cd_toolchain_tool_secretsmanager    
		- ibm_cd_toolchain_tool_bitbucketgit       
		- ibm_cd_toolchain_tool_githubintegrated   
		- ibm_cd_toolchain_tool_githubconsolidated 
		- ibm_cd_toolchain_tool_gitlab 
		- ibm_cd_toolchain_tool_hostedgit          
		- ibm_cd_toolchain_tool_artifactory        
		- ibm_cd_toolchain_tool_custom 
		- ibm_cd_toolchain_tool_pipeline          
		- ibm_cd_toolchain_tool_devopsinsights     
		- ibm_cd_toolchain_tool_slack  
		- ibm_cd_toolchain_tool_sonarqube         
		- ibm_cd_toolchain_tool_hashicorpvault     
		- ibm_cd_toolchain_tool_securitycompliance 
		- ibm_cd_toolchain_tool_privateworker     
		- ibm_cd_toolchain_tool_appconfig          
		- ibm_cd_toolchain_tool_jenkins
		- ibm_cd_toolchain_tool_nexus  
		- ibm_cd_toolchain_tool_pagerduty          
		- ibm_cd_toolchain_tool_saucelabs

* Support CD Tekton Pipeline
    - **Datasources**
        - ibm_cd_tekton_pipeline_definition
        - ibm_cd_tekton_pipeline_trigger_property
        - ibm_cd_tekton_pipeline_property
        - ibm_cd_tekton_pipeline_trigger
        - ibm_cd_tekton_pipeline
    - **Resources**
        - ibm_cd_tekton_pipeline_definition
        - ibm_cd_tekton_pipeline_trigger_property
        - ibm_cd_tekton_pipeline_property
        - ibm_cd_tekton_pipeline_trigger
        - ibm_cd_tekton_pipeline
        
Enhancements
* support for creation of bucket on Satellite location ([3727](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3727))

BUG Fixes
* ibm_dns_glb docs to not specify the enabled argument ([3818](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3818))
* floating ip data source nil fix ([3843](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3843))
* support TLS 1.3 supported ciphers for ibm_cis_domain_settings ([3736](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3736))
* Update provider version in docs ([3750](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3750))
* Update cloud_shell_account_settings.html.markdown ([3812](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3812))
* fix: placement targets empty list ([3787](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3787))
* fix: instance network interface mutex sync ([3791](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3791))
* ibm_scc_posture_scope does not allow for blank-space in description ([3733](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3733))
* SCC provider should output User friendly error messages clearly indicating the problem ([3844](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3844))
* ibm_is_instance should not not suppress change and force new on boot_volume.0.snapshot change ([3819](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3819))
* fix(bare_metal_server) : vlan id in multi nic fix ([3767](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3767))
* fix ibm_container_storage_attachment import ([3862](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3862))

# 1.42.0 (Jun 07, 2022)
Breaking Changes
* Redesign Dns Custom resolver resource and deprecate Custom resolver location ([3820](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3820))

Features
* Support Databases
    - **DataSources**
        - ibm_database_remotes
        - ibm_database_point_in_time_recovery
        - ibm_database_backup
        - ibm_database_backups
* Support PowerSystem
    - **DataSources**
        - ibm_pi_system_pools
* Support EventNotification
    - **DataSources**
        - ibm_en_destination_chrome
        - ibm_en_destination_firefox
    - **Resources**
        - ibm_en_destination_chrome
        - ibm_en_destination_firefox
* Support VPC
    - **DataSources**
        - ibm_is_ssh_keys
        - ibm_is_volumes
* Support Atracker
    - **DataSources**
        - ibm_atracker_settings
    - **Resources**
        - ibm_atracker_settings
* Support Cloudant
    - **Datasources**
        - ibm_cloudant_database
    - **Resources**
        - ibm_cloudant_database
Enhancements
* enhancement(Cloud Databases): add support for User types, roles ([3475](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3475))
* added in preliminary changes for v2 atracker ([3724](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3724))
*  add support to retrieve all secret types ([3793](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3793))
* Support of Transaction-Id for IAM Policy Management ([3518](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3518))
* Deprecate flat list of scaling attributes ([3782](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3782))
* added docs & enhanced network-port-attach resource ([3756](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3756))
* Add support for transit enabled cloud connections ([3758](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3758))
* Support boot volume encryption in cluster and workerpool ([3776](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3776))
* Alias support for key policies ([3768](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3768))
* Enable VPC cluster and worker pool to accept a dedicated host pool ID ([3781](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3781))
* Support SAP deployment type ([3822](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3822))
*

BUGFIXES
* Name Validation Error Fix for VPC ([3675](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3675))
* fix for route by name VPC Routing Table ([3754](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3754))
* Documentation Update VPC Lb Listener and SG Rule ([3673](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3673))
* Added fix for optional weight attribute's value ([3684](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3684))
* Error on failed bucket parsing ([3757](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3757))
* Fix to list the policies even policy contains service specific attributes ([3811](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3811))
* Resource: ibm_resource_key is always recreated even when no changes are available for service Role ([3803](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3803))
* browser cachec issue fixed for zero valid input ([3823](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3823))
* Add new Healthcheck regions ([3827](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3827))


# 1.42.0-beta0 (May 23, 2022)
Features
* Support Databases
    - **DataSources**
        - ibm_database_remotes
        - ibm_database_point_in_time_recovery
        - ibm_database_backup
        - ibm_database_backups
* Support PowerSystem
    - **DataSources**
        - ibm_pi_system_pools
* Support EventNotification
    - **DataSources**
        - ibm_en_destination_chrome
        - ibm_en_destination_firefox
    - **Resources**
        - ibm_en_destination_chrome
        - ibm_en_destination_firefox
* Support VPC
    - **DataSources**
        - ibm_is_ssh_keys
        - ibm_is_volumes
* Support Atracker
    - **DataSources**
        - ibm_atracker_settings
    - **Resources**
        - ibm_atracker_settings
* Support Cloudant
    - **Datasources**
        - ibm_cloudant_database
    - **Resources**
        - ibm_cloudant_database
Enhancements
* enhancement(Cloud Databases): add support for User types, roles ([3475](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3475))
* added in preliminary changes for v2 atracker ([3724](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3724))
*  add support to retrieve all secret types ([3793](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3793))
* Support of Transaction-Id for IAM Policy Management ([3518](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3518))
* Deprecate flat list of scaling attributes ([3782](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3782))

BUGFIXES
* Name Validation Error Fix for VPC ([3675](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3675))
* fix for route by name VPC Routing Table ([3754](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3754))
* Documentation Update VPC Lb Listener and SG Rule ([3673](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3673))
* Added fix for optional weight attribute's value ([3684](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3684))
* Error on failed bucket parsing ([3757](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3757))

# 1.41.1 (May 17, 2022)
BUGFIXES
* fix(floating_ip): fixed nil check on floating ip target ([3783](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3783))

# 1.41.0 (May 04, 2022)
Features
* Support Databases
    - **DataSources**
        - ibm_database_connection
* Support Power Instances
    - **Resources**
        - ibm_pi_cloud_connection_network_attach

* Support Event Notifications
    - **Resources**
        - ibm_en_destination_webhook
        - ibm_en_destination_android
        - ibm_en_destination_ios
        - ibm_en_subscription_sms
        - ibm_en_subscription_email
        - ibm_en_subscription_webhook
        - ibm_en_subscription_android
        - ibm_en_subscription_ios
    - **DataSources**
        - ibm_en_destination_webhook
        - ibm_en_destination_android
        - ibm_en_destination_ios
        - ibm_en_subscription_sms
        - ibm_en_subscription_email
        - ibm_en_subscription_webhook
        - ibm_en_subscription_android
        - ibm_en_subscription_ios
* Support SCC
    - **Resources**
        - ibm_scc_rule
        - ibm_scc_rule_attachment
        - ibm_scc_template
        - ibm_scc_template_attachment
* Support Transist Gateway
    - **Resources**
        - ibm_tg_connection_prefix_filter
     - **DataSources**
        - ibm_tg_connection_prefix_filter
        - ibm_tg_connection_prefix_filters


Enhancements
* Key Policies Deprecate ([3670](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3670))
* Support dynamic network type address ([3659](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3659))
* Update connections schema with mongodbee analytics and bi_connector ([3705](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3705))
* Support port speed attribute in primary nic of instance ([3367](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3367))
* Support for reserved ip changes for VPC ([3712](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3712))
* Support for udp protocol in load balancers ([3711](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3711))
* Fix(status_reasons): fixed status reasons in volume, instance and baremetalserver service ([3725](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3725))
* Support group scaling for IBM Cloud Databases ([3699](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3699))
* Support for port range for public network load balancers ([3660](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3660))
* Add gateway & ip_address_range attribute for network resource ([3585](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3585))

BUGFIXES
* read cloud connection for details ([3650](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3650))
* Fix for optional weight attribute's value ([3685](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3685))
* add DiffSuppressFunc for public key for ibm_is_ssh_key ([3701](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3701))
* Changed list logic in datasource by name ([3415](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3415))
* SCC-PostureManagement Create scope issue ([3708](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3708))
* Limitation on the number of characters in the resource ibm_resource_tag ([3703](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3703))
* Correct a minor typo in docs ([3735](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3735))
* fix(Cloud Databases): fix group scaling crash during provisioning ([3737](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3737))
* fix(Cloud Databases): remove nodeCount update during group resource validation ([3739](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3739))
*  show pi instance provisioning time error ([3738](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3738))
* fix(reserved ip) : bare metal server reserved ip multi nic changes ([3747](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3747))
* Fix: authorisation policy docs ([3753](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3753))
* Broken link for Host Auto Assigment in Cluster ([3752](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3752))
* iam_access_group data source can't see more than 50 access groups ([3728](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3728))

# 1.41.0-beta0 (Apr18, 2022)
Features
* Support Databases
    - **DataSources**
        - ibm_database_connection
* Support Power Instances
    - **Resources**
        - ibm_pi_cloud_connection_network_attach

* Support Event Notifications
    - **Resources**
        - ibm_en_destination_webhook
        - ibm_en_destination_android
        - ibm_en_destination_ios
        - ibm_en_subscription_sms
        - ibm_en_subscription_email
        - ibm_en_subscription_webhook
        - ibm_en_subscription_android
        - ibm_en_subscription_ios
    - **DataSources**
        - ibm_en_destination_webhook
        - ibm_en_destination_android
        - ibm_en_destination_ios
        - ibm_en_subscription_sms
        - ibm_en_subscription_email
        - ibm_en_subscription_webhook
        - ibm_en_subscription_android
        - ibm_en_subscription_ios
* Support SCC
    - **Resources**
        - ibm_scc_rule
        - ibm_scc_rule_attachment
        - ibm_scc_template
        - ibm_scc_template_attachment
* Support Transist Gateway
    - **Resources**
        - ibm_tg_connection_prefix_filter
     - **DataSources**
        - ibm_tg_connection_prefix_filter
        - ibm_tg_connection_prefix_filters


Enhancements
* Key Policies Deprecate ([3670](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3670))
* Support dynamic network type address ([3659](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3659))
* Update connections schema with mongodbee analytics and bi_connector ([3705](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3705))
* Support port speed attribute in primary nic of instance ([3367](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3367))
* Support for reserved ip changes for VPC ([3712](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3712))
* Support for udp protocol in load balancers ([3711](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3711))
* Fix(status_reasons): fixed status reasons in volume, instance and baremetalserver service ([3725](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3725))
* Support group scaling for IBM Cloud Databases ([3699](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3699))
* Support for port range for public network load balancers ([3660](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3660))

BUGFIXES
* read cloud connection for details ([3650](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3650))
* Fix for optional weight attribute's value ([3685](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3685))
* add DiffSuppressFunc for public key for ibm_is_ssh_key ([3701](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3701))
* Changed list logic in datasource by name ([3415](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3415))
* SCC-PostureManagement Create scope issue ([3708](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3708))

# 1.40.1 (Apr08, 2022)
BUGFIXES:
 Allow 0 allocation for optional scaling group types ([3714](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3714))
 

# 1.40.0 (Mar31, 2022)
Features
* Support VPC Infrastructure
    - **DataSources**
        - ibm_is_lb_listener
        - ibm_is_lb_listeners
        - ibm_is_lb_listener_policies
        - ibm_is_lb_listener_policy
        - ibm_is_lb_listener_policy_rule
        - ibm_is_lb_listener_policy_rules
        - ibm_is_security_groups
        - ibm_is_security_group_rule
        - ibm_is_security_group_rules
        - ibm_is_ipsec_policy
        - ibm_is_ipsec_policies
        - ibm_is_ike_policies
        - ibm_is_ike_policy
    - **Resources**
        - ibm_is_subnet_public_gateway_attachment
        - ibm_is_subnet_routing_table_attachment
* Support Transist Gateway
    - **DataSources**
        - ibm_tg_route_report
        - ibm_tg_route_reports
    - **Resources**
        - ibm_tg_route_report
* Support CIS
    - **DataSources**
        - ibm_cis_alerts
    - **Resources**
        - ibm_cis_alerts

* Support Power Instance
    - **Datasources**
        - ibm_pi_storage_pool_capacity
        - ibm_pi_storage_pools_capacity
        - ibm_pi_storage_type_capacity
        - ibm_pi_storage_types_capacity

ENHANCEMENTS:
* Support resize boot volume for VPC instance ([2205](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2205))
* Support to allow SAP create with a placement group ([3633](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3633))
* Support VM Host Failure - Available Policy Development ([3604](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3604))
* Support default forwarding rule in custom resolver ([3588](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3588))
* Support trusted-profiles for access group members ([3651](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3651))
* Authorization policy to support any service specific attributes ([3482](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3482))
* Move `name` as optional argument in ibm_is_region datasource ([3431](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3431))
* Add note for disabled VTL creation ([3693](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3693))

BUGFIXES:
* resource_ibm_is_virtual_endpoint_gateway does not allow for full name character limit ([3606](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3606))
* Added the doc support for endpoint gateway identifier ([3649](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3649))
* Fixing resource_attributes to support the service specific roles  ([3648](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3648))
* endpoint gateways target datasource fix, interchanging the order and including other targets ([3492](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3492))

# 1.40.0-beta0 (Mar17, 2022)
Features
* Support VPC Infrastructure
    - **DataSources**
        - ibm_is_lb_listener
        - ibm_is_lb_listeners
        - ibm_is_lb_listener_policies
        - ibm_is_lb_listener_policy
        - ibm_is_lb_listener_policy_rule
        - ibm_is_lb_listener_policy_rules
        - ibm_is_security_groups
        - ibm_is_security_group_rule
        - ibm_is_security_group_rules
        - ibm_is_ipsec_policy
        - ibm_is_ipsec_policies
        - ibm_is_ike_policies
        - ibm_is_ike_policy
    - **Resources**
        - ibm_is_subnet_public_gateway_attachment
        - ibm_is_subnet_routing_table_attachment
* Support Transist Gateway
    - **DataSources**
        - ibm_tg_route_report
        - ibm_tg_route_reports
    - **Resources**
        - ibm_tg_route_report
* Support CIS
    - **DataSources**
        - ibm_cis_alerts
    - **Resources**
        - ibm_cis_alerts

* Support Power Instance
    - **Datasources**
        - ibm_pi_storage_pool_capacity
        - ibm_pi_storage_pools_capacity
        - ibm_pi_storage_type_capacity
        - ibm_pi_storage_types_capacity

ENHANCEMENTS:
* Support resize boot volume for VPC instance ([2205](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2205))
* Support to allow SAP create with a placement group ([3633](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3633))
* Support VM Host Failure - Available Policy Development ([3604](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3604))
* Support default forwarding rule in custom resolver ([3588](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3588))
* Support trusted-profiles for access group members ([3651](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3651))
* Authorization policy to support any service specific attributes ([3482](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3482))
* Move `name` as optional argument in ibm_is_region datasource ([3431](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3431))

BUGFIXES:
* resource_ibm_is_virtual_endpoint_gateway does not allow for full name character limit ([3606](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3606))
* Added the doc support for endpoint gateway identifier ([3649](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3649))
* Fixing resource_attributes to support the service specific roles  ([3648](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3648))
* endpoint gateways target datasource fix, interchanging the order and including other targets ([3492](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3492))

# 1.39.2 (Mar11, 2022)
BUGFIXES:
* fix proxy on kms client ([3634](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3634))
* Fix crash on resource_instance ([3643](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3643))
* SCC data.ibm_scc_posture_latest_scans failing ([3594](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3594))
* Error: Collector Creation in Security Posture Management API service: us region not supported ([3630](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3630))
* Fix ibm_cloud_shell_account_settings IAM issue ([3654](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3654))

# 1.39.1 (Mar02, 2022)
BUGFIXES:
* Fix the breaking change of 1.39.0 release for ibm_is_instance resource for force_new on VSI([3619](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3619))

# 1.39.0 (Feb28, 2022)
Features
* Support VPC Infrastructure
    - **DataSources**
        - ibm_is_lb_pool
        - ibm_is_lb_pools
        - ibm_is_lb_pool_member
        - ibm_is_lb_pool_members
        - ibm_is_vpc_address_prefix
        - ibm_is_vpc_routing_table
        - ibm_is_vpc_routing_table_route
        - ibm_is_vpn_gateway
        - ibm_is_vpn_gateway_connection
        - ibm_is_bare_metal_server_disk
        - ibm_is_bare_metal_server_disks
        - ibm_is_bare_metal_server_initialization
        - ibm_is_bare_metal_server_network_interface_floating_ip
        - ibm_is_bare_metal_server_network_interface_floating_ips
        - ibm_is_bare_metal_server_network_interface
        - ibm_is_bare_metal_server_network_interfaces
        - ibm_is_bare_metal_server_profile
        - ibm_is_bare_metal_server_profiles
        - ibm_is_bare_metal_server
        - ibm_is_bare_metal_servers
    - **Resources**
        - ibm_is_bare_metal_server_action
        - ibm_is_bare_metal_server_disk
        - ibm_is_bare_metal_server_network_interface_allow_float
        - ibm_is_bare_metal_server_network_interface_floating_ip
        - ibm_is_bare_metal_server_network_interface
        - ibm_is_bare_metal_server
* Support for CIS
    - **Datasources**
        - ibm_cis_webhooks
    - **Resources**
        - ibm_cis_webhook

       
ENHANCEMENTS:
* Remove hmac-md5-96 authentication from IPSec VPN policy ([3515](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3515))
* added status reason in case of failure in is_instance ([3505](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3505))
* Updated PI IKEPolicy and IPSecPolicy ([3530](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3530))
* IBM Cloud VPC Documentation Update ([3479](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3479))
* IKS image security enforcement ([3344](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3344))
* add check for empty slice ([3528](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3528))
* push notifications deprecation label ([3550](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3550))
* Allow passing cloud connection while creating DHCP server ([3559](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3559))
* Add dns attribute to power network data source ([3551](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3551))
* Support arm , arm64 for terraform provider ([2626](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2626))
* Add support to use the new IAMAuthenticatior for token ([3542](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3542))
* add network update support ([3565](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3565))
* P10 enablement for SAP ([3534](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3534))
* implement provider gateway vlan assignment ([3558](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3558))
* Support storage pool and affinity while importing image ([3543](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3543))
* Add custom_script argument to satellite_host_script datasource to avoid provider dependancy ([3579](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3579))
* Add wait_till to satellite_host resource to wait until location normal ([3579](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3579))
* Snapshot captured at deletable ([3555](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3555))
* Support access tags in ibm_iam_access_group_policy ([3290](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3290))
* Bumpup: terraform-plugin-sdk to 2.10.1 ([3613](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3613))
* support instance metadata and trusted profile ([3616](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3616))

BUGFIXES:
* Not able to list all the subnets for a given VPC ([3506](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3506))
* panic while creating cloud connections ([3498](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3498))
* schema and test fix for instance group manager action ([3495](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3495))
* API endpoints for SCC are wrong and cannot be changed ([3524](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3524))
* log improvement: subnet delete retry logs ([3485](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3485))
* added a check to check if remote sg id exists in is_security_group_rule ([3493](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3493))
*  document fix for routing table ([3489](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3489))
* Fix the sub category ([3560](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3560))
* Added more info for href ([3486](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3486))
* Subnet pagination fix in is_vpcs datasource ([3537](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3537))
* Destroy resource group failed with 204 ([3529](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3529))
* ibm_is_lb resource is silently updating subnet details in state with no actual affect on lb ([3538](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3538))
* Fixed satellite location immutable error ([3537](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3537))
* ibm_cloud_shell_account_settings issue ([3289](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3289))
* typo in "network_interaces" ([3548](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3548))
* Add additional example to ibm_iam_trusted_profile_policy ([3532](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3532))
* remove default from resource_group doc as datasource doesnot support ([3564](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3564)) 
* Fix: visibility private for vpc ([3578](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3578))
* Satellite Host Attach Configuration Adds Incorrect Repos ([3576](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3576))
* Missing Resource Group Id argument in Data block docs ([3586](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3586))
* container_addons doesnt update version as expected ([3584](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3584))
* added description for managed_from parameter in satellite location ([3591](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3591)) 
* Fix for ibm_pi_volume response ([3598](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3598))
* Added a check for "Deleted" status of the VPC worker pool ([3595](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3595))
* Fix: Ignore 500 on getclusterkey in container cluster datasource ([3613](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3613))
* Power client update to refresh token for every API call ([3614](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3614))


# 1.39.0-beta0 (Feb16, 2022)
Features
* Support VPC Infrastructure
    - **DataSources**
        - ibm_is_lb_pool
        - ibm_is_lb_pools
        - ibm_is_lb_pool_member
        - ibm_is_lb_pool_members
        - ibm_is_vpc_address_prefix
        - ibm_is_vpc_routing_table
        - ibm_is_vpc_routing_table_route
        - ibm_is_vpn_gateway
        - ibm_is_vpn_gateway_connection
        - ibm_is_bare_metal_server_disk
        - ibm_is_bare_metal_server_disks
        - ibm_is_bare_metal_server_initialization
        - ibm_is_bare_metal_server_network_interface_floating_ip
        - ibm_is_bare_metal_server_network_interface_floating_ips
        - ibm_is_bare_metal_server_network_interface
        - ibm_is_bare_metal_server_network_interfaces
        - ibm_is_bare_metal_server_profile
        - ibm_is_bare_metal_server_profiles
        - ibm_is_bare_metal_server
        - ibm_is_bare_metal_servers
    - **Resources**
        - ibm_is_bare_metal_server_action
        - ibm_is_bare_metal_server_disk
        - ibm_is_bare_metal_server_network_interface_allow_float
        - ibm_is_bare_metal_server_network_interface_floating_ip
        - ibm_is_bare_metal_server_network_interface
        - ibm_is_bare_metal_server
* Support for CIS
    - **Datasources**
        - ibm_cis_webhooks
    - **Resources**
        - ibm_cis_webhook

       
ENHANCEMENTS:
* Remove hmac-md5-96 authentication from IPSec VPN policy ([3515](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3515))
* added status reason in case of failure in is_instance ([3505](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3505))
* Updated PI IKEPolicy and IPSecPolicy ([3530](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3530))
* IBM Cloud VPC Documentation Update ([3479](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3479))
* IKS image security enforcement ([3344](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3344))
* add check for empty slice ([3528](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3528))
* push notifications deprecation label ([3550](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3550))
* Allow passing cloud connection while creating DHCP server ([3559](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3559))
* Add dns attribute to power network data source ([3551](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3551))
* Support arm , arm64 for terraform provider ([2626](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2626))
* Add support to use the new IAMAuthenticatior for token ([3542](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3542))
* add network update support ([3565](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3565))
* P10 enablement for SAP ([3534](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3534))
* implement provider gateway vlan assignment ([3558](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3558))
* Support storage pool and affinity while importing image ([3543](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3543))
* Add custom_script argument to satellite_host_script datasource to avoid provider dependancy ([3579](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3579))
* Add wait_till to satellite_host resource to wait until location normal ([3579](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3579))
* Snapshot captured at deletable ([3555](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3555))
* Support access tags in ibm_iam_access_group_policy ([3290](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3290))

BUGFIXES:
* Not able to list all the subnets for a given VPC ([3506](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3506))
* panic while creating cloud connections ([3498](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3498))
* schema and test fix for instance group manager action ([3495](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3495))
* API endpoints for SCC are wrong and cannot be changed ([3524](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3524))
* log improvement: subnet delete retry logs ([3485](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3485))
* added a check to check if remote sg id exists in is_security_group_rule ([3493](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3493))
*  document fix for routing table ([3489](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3489))
* Fix the sub category ([3560](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3560))
* Added more info for href ([3486](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3486))
* Subnet pagination fix in is_vpcs datasource ([3537](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3537))
* Destroy resource group failed with 204 ([3529](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3529))
* ibm_is_lb resource is silently updating subnet details in state with no actual affect on lb ([3538](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3538))
* Fixed satellite location immutable error ([3537](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3537))
* ibm_cloud_shell_account_settings issue ([3289](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3289))
* typo in "network_interaces" ([3548](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3548))
* Add additional example to ibm_iam_trusted_profile_policy ([3532](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3532))
* remove default from resource_group doc as datasource doesnot support ([3564](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3564)) 
* Fix: visibility private for vpc ([3578](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3578))
* Satellite Host Attach Configuration Adds Incorrect Repos ([3576](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3576))
* Missing Resource Group Id argument in Data block docs ([3586](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3586))
* container_addons doesnt update version as expected ([3584](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3584))
* added description for managed_from parameter in satellite location ([3591](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3591)) 

# 1.38.2 (Feb09, 2022)
BugFIXES:
* Updating members_cpu_allocation_count to 0 fails with ibm_database ([3567](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3567))


# 1.38.1 (Jan27, 2022)
ENHANCEMENTS:
* Support intergation of security groups with virtual private endpoint gateway ([3488](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3488))
* Update the deprecation notice of policy management from KMS resource ([3520](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3520)

# 1.38.0 (Jan12, 2022)
Features
* Support Power Instance
    - **DataSources**
        - ibm_pi_cloud_connections
        - ibm_pi_keys
        - ibm_pi_console_languages
        - ibm_pi_placement_group
        - ibm_pi_placement_groups
    - **Resources**
        - ibm_pi_console_language
        - ibm_pi_volume_attach
        - ibm_pi_image_export
        - ibm_pi_placement_group
        - ibm_pi_capture
* Support Security and Compliance Center
    - **DataSources**
        - ibm_scc_posture_profile
        - ibm_scc_posture_group_profile
        - ibm_scc_posture_scope_correlation
    - **Resources**
        - ibm_scc_posture_collector
        - ibm_scc_posture_scope
        - ibm_scc_posture_credential
* Support IAM Authorization Policies
    - **Datasources**
        - ibm_iam_authorization_policies
* Support Satellite Cluster
    - **Datasources**
        - ibm_satellite_cluster_worker_pool_zone_attachment
    - **Resources**
        - ibm_satellite_cluster_worker_pool_zone_attachment
       
ENHANCEMENTS:
* Add issue labeler workflow and update issue template ([3430](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3430))
* Enhance log statements ([3442](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3442))
* Allow retrieval of IAM access tag (data source) ([3287](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3287))
* Refactor multiple pi resources with context awareness ([3429](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3429)) 
* Support: mysql in ibm_database ([3454](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3454))
* ibm_resource_instance cannot except json in "parameters" ([3458](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3458))
* Add release process and maintainers ([3472](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3472))
* Allow mixed storage for pi_instance ([3484](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3484))

BUGFIXES:
* Invite Users module it's not working as expected ([3226](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3226))
* Docs: data iam_trusted_profile_claim_rules actually is documented as iam_trusted_profiles_claim_rules ([3421](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3421))
* ibm_iam_custom_role does not honor description changes ([3353](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3353))
* ibm_cos_bucket documented storge_class flex does not work ([3349](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3349))
* DocFix: satellite_host resource ([3437](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3437))
* ibm_cos_bucket missing docs and missing s3_endpoint_direct ([3436](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3436))
* Docs for data ibm_iam_trusted_profile_policy: Wrong title for examples ([3425](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3425))
* Not documented: how to access ibm_resource_key.objectstorage.credentials.cos_hmac_keys.access_key_id ([2180](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2180))
* Cloud Object Storage access key and secret ([1860](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1860))
* ibm_resource_key for COS does not nest HMAC object ([1741](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1741))
* IBM Cloud Databases CLI Behavior Update ([1387](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1387))
* Feature Request: Make ibm_resource_key resources linked to ibm_database resources more importable ([1232](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1232))
* Bug Fix Policy ([3413](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3413))
* Fix: iam url issue in kms client ([3417](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3417))
* Update doc of container service bind ([3450](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3450))
* removed required tag from name in instance template ([3432](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3432))
* DocFix: dns services ([3462](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3462))
* Fix the entitlement for VPC worker pool ([3464](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3464))
* Add Watson query example ([3451](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3451))
* Issue with Access group creation ([3476](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3476))

# 1.37.1 (Dec09, 2021)
BUGFIXES:
* Regression: Breaking change on policy resources ([3410](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3410))
* ibm_pi_instance is not able to complete with pi_health_status "WARNING" ([3401](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3401))
* trusted profiles- doc fixes ([3407](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3407))
* DocFix: cis_page_rule resource and flowlog and cbr_rule datasource ([3402](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3402))

# 1.37.0 (Dec07, 2021)
Features
* Support VPC Infrastructure
    - **DataSources**
        - ibm_is_instance_network_interface
        - ibm_is_instance_network_interfaces
    - **Resources**
        - ibm_is_instance_network_interface
* Support Power Instance
    - **Resources**
        - ibm_pi_ike_policy
        - ibm_pi_ipsec_policy
        - ibm_pi_vpn_connection
    - **DataSources**
        - ibm_pi_sap_profiles
        - ibm_pi_sap_profile
* Support Security and Compliance Center
    - **DataSources**
        - ibm_scc_account_location
        - ibm_scc_account_locations
        - ibm_scc_account_settings
    - **Resources**
        - ibm_scc_account_settings
* Support Container, Satellite nlb and ALB
    - **Resources**
        - ibm_container_nlb_dns
        - ibm_satellite_location_nlb_dns
        - ibm_container_alb_create
        - ibm_container_vpc_alb_create
* Support IAM Trusted Profiles
     - **DataSources**
        - ibm_iam_trusted_profile_claim_rules
        - ibm_iam_trusted_profile_links
        - ibm_iam_trusted_profiles
* Support Context Based Restriction
    - **Resources**
        - ibm_cbr_zone
        - ibm_cbr_rule
     - **DataSources**
        - ibm_cbr_zone
        - ibm_cbr_rule
* Support IAM Access Group
     - **DataSources**
        - ibm_iam_access_group_policy
       


ENHANCEMENTS:
* Support abort_incomplete_multipart_upload_days, expired object delete markers and noncurrent_version_expiration feature for Cloud Object Storage ([3359](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3359))
* Added new resource attribute service_type for access policies ([3347](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3347))
* vpc-go-sdk migration to 0.14.0 ([3376](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3376))
* Add support for VTL in Power Instance ([3328](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3328))
* Added filters to VPC Volume Snapshot collection datasource ([3238](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3238))


BUGFIXES:
* Documentation fixes for Security and Compilance ([3342](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3342))
* Update container-registry SDK and fix default region ([3356](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3356))
* Fix: private endpoint for secrets manager ([3378](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3378))
* Fix ibm_appid_token_config has source "roles", but missing in docs ([3370](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3370))
* Bug in documentation for ibm_access_group_policy ([3365](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3365))
* Bug in documentation for ibm_iam_api_key datasource ([3363](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3363))
* Rename App ID provider ([3355](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3355))
* DocFix: remove API from Activity Tracker subcategory ([3379](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3379))
* doc fix for instance and subnet ([3372](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3372))
* Inconsistent examples ibm_kms_key_rings example ([3279](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3279))
* added wait logic for security group target ([3373](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3373))
* do not return error when topic exists in creation ([3223](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3223))
* private endpoints doesn't work for iam_access_group resources ([3340](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3340))
* ibm_kms_key and ibm_kp_key produce inconsistent plan/apply ([3314](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3314))
* Actions fail to import for ibm_cis_page_rule resources ([2765](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2765))
* Upgrade of MongoDB from standard to enterprise should not work ([3327](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3327))

# 1.36.0 (Nov16, 2021)
Features
* Support VPC Infrastructure
    - **DataSources**
        - ibm_is_network_acls
        - ibm_is_network_acl
        - ibm_is_flow_log
        - ibm_is_floating_ips
    - **Resources**
        - ibm_is_instance_action
* Support Power Instance
    - **DataSources**
        - ibm_pi_dhcp
        - ibm_pi_dhcps
        - ibm_pi_cloud_connection
    - **Resources**
        - ibm_pi_dhcp
        - ibm_pi_cloud_connection
* Support SCC Security Insights
    - **DataSources**
        - ibm_scc_si_occurrence
        - ibm_scc_si_occurrences
    - **Resources**
        - ibm_scc_si_occurrence
* Support Schematics
    - **DataSources**
        - ibm_schematics_inventory
        - ibm_schematics_resource_query
    - **Resources**
        - ibm_schematics_inventory
        - ibm_schematics_resource_query


ENHANCEMENTS:
* Support storage pool and affinity for instance and volume ([3270](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3270))

* Import image from public and private COS bucket ([3265](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3265))

* Support br-sao region for Container Registry ([3258](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3258))

* Support gpu for instance profile datasource ([3158](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3158))

* add resource group to cm_catalog resource and datasource ([3291](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3291))

* Support: default headers for service client ([3257](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3257))

* Support VPC instance bandwidth ([3156](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3156))

* Postgres configuration through terraform ([3278](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3278))

* Configure Redis database ([1428](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1428))

* Add force_create to create classic infrastructure reserved capacity with exisitng name ([3306](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3306))

* vtl and sap options for Power Instance stock images ([3310](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3310))

* Allow Power Instance volume update when in-use ([3323](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3323))

* added support for enabling sdk debug logging ([3268](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3268))


BUGFIXES:
* updated docs and examples for vpn gateway and gateway connections ([3283](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3283))

* Load Balancer cannot be updated because its status is 'UPDATE_PENDING' ([3006](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3006))

* VPC Address Prefix can't delete even though Subnet Deletion process complete. ([2759](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2759))

* RC api returning incorrect response when instance already exists ([3187](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3187))

* Failures when creating ibm_database because of bad values for allocations should be more clear ([3294](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3294))

* Fix: nil pointer on pi_key ([3133](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3133))

* Regression: "AuthorizationDelegator" no longer works in 1.30.0 ([3013](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3013))

* adding custom retry to fix the enabling of logging and moniotirg ([3319](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3319))

* Update ibm-hpcs-tke-sdk ([3313](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3313))

* Update schematics terraform resources and datasources based on latest API's ([2901](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2901))

* ibm_schematics_workspace adding template_inputs causes panic ([3295](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3295))

* ibm_schematics_workspace add template_git_url fails ([3296](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3296))

* Cannot provision Schematics workspaces with recent versions of Terraform ([3048](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3048))


# 1.35.0 (Oct28, 2021)

Features
* Support Event Notifications
    - **Resources**
        - ibm_en_destination
        - ibm_en_topic
        - ibm_en_subscription
    - **DataSources**
        - ibm_en_destination
        - ibm_en_destinations
        - ibm_en_topic
        - ibm_en_topics
        - ibm_en_subscription
        - ibm_en_subscriptions

* Support Container Storage Attachment
    - **Resources**
        - ibm_container_storage_attachment
    - **DataSources**
        - ibm_container_storage_attachment

ENHANCEMENTS:

*  Implemented feature to edit BGP IPs and ASN values for non provider and provider flow  gateways ([3186](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3186))

* Support VPC load balancer https redirect ([3115](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3115))

* Support route_mode for VPC NLB vnf ([3208](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3208))

* Support direct endpoints for cos_bucket ([3252](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3252))

* Add support for pi_network jumbo option ([3255](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3255))

* Support enable/update the BFD config for the DirectLink gateways ([3194](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3194))

* Added port range support for VPC NLBs with route mode enabled ([3207](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3207))

* Support high availablity for custom resolver ([3190](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3190))

* Add resource_group_id ibm_container_vpc_alb resource and datasource ([2768](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2768))


BUGFIXES:

* ibm_is_lb_pool_member member weight is not working as expected when weight is 0 ([3124](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3124))

* Documentation for data "ibm_iam_trusted_profile_policy" is wrong ([3201](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3201))

* Wrong documentation for iam_trusted_profile_claim_rule ([3216](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3216))

* added placement_group documentation in instance_template ([3210](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3210))

* VPE documentation issue ([3225](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3225))

* Do not ignore label to enable node-local-dns-enabled for kubernetes service ([3232](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3232))

* ibm_is_lb how to create a network load balancer, what is profile, resource group name ([3108](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3108))

* Fix VPC lb listener access protocol ([3240](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3240))

* Fix documentation for atracker resources ([3246](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3246))

* Conrefs in https://registry.terraform.io/providers/IBM-Cloud/ibm/latest/docs/resources/database documentation ([3233](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3233))

* Fix ibm_iam_trusted_profile_link, argument "namespace" is not required if cr_type is VSI ([3219](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3219))

* Add links to App ID documentation ([3220](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3220))

* What are the valid values for the attributes? (App ID) ([3221](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3221))

* ibm_iam_access_group_members incorrect state when members list >50 users ([3189](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3189))

* ibm_iam_account_settings modify issue ([3249](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3249))

* ibm_cloud_shell_account_settings modify issue ([3247](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3247))

* ibm_cloud_shell_account_settings delete issue ([3242](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3242))

* Error downloading the cluster config - config.zip: no such file or directory ([2806](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2806))

* IBM Cloud Shell data resource ([3275](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3275))

# 1.34.0 (Oct12, 2021)

**Note**
This release replace github.com/dgrijalva/jwt-go dependency with github.com/golang-jwt/jwt to fix ([CVE-2020-26160] (https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3191))
FEATURES:

* Support IBM IAM Trusted Profile
    - **Resources**
        - ibm_iam_trusted_profile
        - ibm_iam_trusted_profile_claim_rule
        - ibm_iam_trusted_profile_link
        - ibm_iam_trusted_profile_policy
    - **DataSources**
        - ibm_iam_trusted_profile
        - ibm_iam_trusted_profile_claim_rule
        - ibm_iam_trusted_profile_link
        - ibm_iam_trusted_profile_policy

* Support Classic Infrastructure Reserved Capacity
    - **Resources**
        - ibm_compute_reserved_capacity
    - **DataSources**
        - ibm_compute_reserved_capacity
ENHANCEMENTS:

* Support for reading the endpoints for supported IBM Cloud Services via file ([3071](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3071))

* Support for provisioning Monthly based servers on reserved capacity ([3185](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3185))

## 1.33.1 (Oct1, 2021)
ENHANCEMENTS

* Ability to provide IP address for provisioning Power Systems using ibm_pi_instance resource ([3102](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3102))


BUGFIXES

* Regression: ibm_is_instance 1.33.0 no longer creates user_data ([3163](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3163))
* Regression: "AuthorizationDelegator" no longer works in 1.30.0 ([3013](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3013))
* Fix: ibm_database datasource returns nil for connectionstrings ([3166](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3166))
* Fix: appid token config destination_claim should be optional ([3143](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3143))
* Dereferencing of rg id, name and target update in floating ip ([3164](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3164))
* DocFix: Security and compliance doc updates ([3155](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3155))

# 1.33.0 (Sep28, 2021)
FEATURES:

* Support IBM Cloud Shell
    - **Resources**
        - ibm_cloud_shell_account_settings
    - **DataSources**
        - ibm_cloud_shell_account_settings

* Support Security and Compliance Center
    - **Resources**
        - ibm_scc_si_note
    - **DataSources**
        - ibm_scc_si_note
        - ibm_scc_si_notes
        - ibm_scc_si_providers
        - ibm_scc_posture_scopes
        - ibm_scc_posture_latest_scans
        - ibm_scc_posture_profiles
        - ibm_scc_posture_scan_summary
        -ibm_scc_posture_scan_summaries

* Support Event Streams Schema
    - **Resources**
        - ibm_event_streams_schema
    - **DataSources**
        - ibm_event_streams_schema

* Support AppID
    - **Resources**
        - ibm_appid_idp_google
        - ibm_appid_mfa_channel
    - **DataSources**
        - ibm_appid_idp_google
        - ibm_appid_mfa_channel

* Support Cloudant database
    - **Resources**
        - ibm_cloudant
    - **DataSources**
        - ibm_cloudant

* Support CIS Firewall Rules
    - **Resources**
        - ibm_cis_firewall_rules
    - **DataSources**
        - ibm_cis_firewall_rule

**DeprecationMessage**: Resource ibm_is_security_group_network_interface_attachment is deprecated. Use ibm_is_security_group_target to attach a network interface to a security group

ENHANCEMENTS

* Feature: add support for Transit Gateway DLaaS ([3105](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3105))
* Added changes for adjustable iops, capacity and volume profile ([3068](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3068))


BUGFIXES

* suppressing the change in wait_before_delete on import of is_instance which showed update in place ([3075](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3075))

* ibm_is_instance_volume_attachment - Multiple volume creation and attachment failure ([3077](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3077))

* data "ibm_container_cluster_config" results in intermittent authentication as 'system:anonymous ([2811](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2811))

* Added missing crn to is resources ([3130](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3130))

* Added a document update regarding allow_ip_spoofing on network interfaces for VSIs ([3145](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3145))

* Changing an ibm_is_instance causes the associated ibm_is_floating_ip to fail ([3110](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3110))

* Updated security group target APIs ([2896](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2896))

* Failure modifying volume_name in ibm_is_instance_volume_attachment ([3089](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3089))

* DocFix: Satellite Link and Endpoint resources ([3152](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3152))

* Added zones parameter to immutable list and added hosts parameter to satellite location data source ([3137](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3137))


## 1.32.1 (Sep20, 2021)
ENHANCEMENTS

* Add support for provisioning Enterprise DB, Enterprise Mongo, Cassandra Databases ([3097](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3097))

* upgraded the vpc-go-sdk to use the latest 0.10.0 version ([3067](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3067))

* Add missing crn to vpc resources and datasources ([3111](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3111))

* Support expandabale volumes for VPC ([2668](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2668))


BUGFIXES

* IBM Cloud Terraform Provider is unable to upgrade the "vpc-block-csi-driver" addon on VPC Gen 2 Clusters ([2988](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2988))


## 1.32.0 (Sep14, 2021)
FEATURES:

* Support APP ID 
    - **Resources**
        - ibm_appid_idp_facebook
		- ibm_appid_cloud_directory_user
        - ibm_appid_mfa
    - **DataSources**
        - ibm_appid_idp_facebook
		- ibm_appid_cloud_directory_user
		- ibm_appid_mfa 

* Support VPC Placement group
    - **Resources**
        - ibm_is_placement_group
    - **DataSources**
        - ibm_is_placement_group
        - ibm_is_placement_groups

ENHANCEMENTS

* Add name validation at missing places for VPC resources ([3051](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3051))

* Remove computed from encryption_key and added encryption_type for VPC volume ([3057](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3057))

* Set direction which was showing a change on import for VPC security group rule resource ([3063](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3063))

* Add `description` argument support for IAM policy mangament resources ([3095](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3095))

* Add crn for is_vpc resource and datasource in default sg and network acl ([3096](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3096))

* Ignore 410-gone status during deletion of resource instance ([3096](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3096))


BUGFIXES

* doc update: security group rule import update ([3066](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3066))

* Fix satellite route example issue ([3073](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3073))

* Fix docs for DNS Custom resolver ([2997](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2997))

* Fix orientation of docs for appID ([3091](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3091))

* Fix encryption crn for volume attachment ([3101](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3101))

* Fix endpoint gateways data source ([3100](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3100))

* Fix function namespace ([3103](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3103))

## 1.31.0 (Sep06, 2021)
FEATURES:

* Support APP ID 
    - **Resources**
        - ibm_appid_idp_cloud_directory
		- ibm_appid_idp_saml
        - ibm_appid_password_regex
        - ibm_appid_theme_text
        - ibm_appid_idp_custom
        - ibm_appid_audit_status
        - ibm_appid_action_url
        - ibm_appid_languages
        - ibm_appid_theme_color
        - ibm_appid_cloud_directory_template
    - **DataSources**
        - ibm_appid_idp_cloud_directory
		- ibm_appid_idp_saml
		- ibm_appid_password_regex 
        - ibm_appid_theme_text
        - ibm_appid_idp_custom
        - ibm_appid_audit_status
        - ibm_appid_action_url
        - ibm_appid_languages
        - ibm_appid_theme_color
        - ibm_appid_cloud_directory_template
        - ibm_appid_idp_saml_metadata
        - ibm_appid_roles
        - ibm_appid_applications

* Support Satellite endpoint/link
    - **Resources**
        - ibm_satellite_link
        - ibm_satellite_endpoint
    - **DataSources**
        - ibm_satellite_link
        - ibm_satellite_endpoint

ENHANCEMENTS

* changed the validator invoked in nacl rule  ([3052](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3052))

* added forcenew for address_prefix_management in ibm_is_vpc ([3025](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3025))

* added support for identifier in is_image datasource ([3012](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3012))

* Add wait for master ready before cluster integrations ([3035](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3035))

* Adding Steering_Policy property ([3061](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3061))

* Add support for vSCSI pi_instance deployment ([3070](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3070))



BUGFIXES

* DocUpdate: add routing table route example ([3046](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3046))

* Documentation for endpoint gateway ([3034](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3034))

* updated cdn and app config resource ([3060](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3060))

* Removed dereferencing in resource group name for is_volume ([3058](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3058))

* Change the instance attribute types ([3053](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3053))

* Increase robustness of openshift login flow ([3035](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3035))

* Remove the references to duplicate KMS Key policies docs ([3055](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3055))


## 1.30.2 (Aug31, 2021)
ENHANCEMENTS

* removed 100GB constraint to support any size (10-250GB) size boot volumes ([3030](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3030))

* upgraded the vpc-go-sdk to use the 0.9.0 version ([3023](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3023))

*  Add crn_token attribute to satellite cluster resource to support remote location ([3032](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3032))


BUGFIXES

* Doc fixes for ibm_database and ibm_container_cluster_config ([3033](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3033))

* Updated argument desc for zones and entitlement ([3016](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3016))

## 1.30.1 (Aug25, 2021)
ENHANCEMENTS

* fix: setting listener policy rule values ([2964](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2964))

* fix for updating serviceIds and IBMIds in read method ([2987](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2987))

* Include the IBM ID in the ibm_iam_user_profile data source ([2940](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2940))

* Added new San Paolo (SAO) MZR support for cos ([3022](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3022))

* Added retry logic to get location information ([3020](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/3020))


BUGFIXES

* Argument reference description for MFA setting is incorrect ([2984](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2984))

* Doc fixes for APPID Management resources and datasources ([3001](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3001))

* fixed the header for is_instance_volume_attachment doc ([3005](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/3005))

## 1.30.0 (Aug18, 2021)
FEATURES:

* Support ATracker 
    - **Resources**
        - ibm_atracker_target
		- ibm_atracker_route
    - **DataSources**
        - ibm_atracker_targets
		- ibm_atracker_routes
		- ibm_atracker_endpoints 
    
* Support APP ID
    - **Resources**
        -ibm_appid_application
	    - ibm_appid_token_config
        - ibm_appid_application_roles
        - ibm_appid_application_scopes
        - ibm_appid_redirect_urls
        - ibm_appid_role
    - **DataSources**
        - ibm_appid_application
        - ibm_appid_application_roles
        - ibm_appid_application_scopes
        - ibm_appid_token_config
        - ibm_appid_redirect_urls
        - ibm_appid_role

* Support KMS Policies 
    - **Resources**
        - ibm_kms_key_policies
    - **DataSources**
        - ibm_kms_key_policies
    
* Support VPC Instance Template 
    - **DataSources**
        - ibm_is_instance_template

* Support NLB DNS for cluster and Satellite Location
    - **DataSources**
        - ibm_container_nlb_dns
        - ibm_satellite_location_nlb_dns
   

ENHANCEMENTS

* Support access management tags in subnet resource and datasource ([2778](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2778))

* Added pagination support for endpoint gateway targets datasource ([2904](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2904))

* Support lb cookie session persistence ([#2884](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2884))

* Expose catalog type for catalog management ([#2924](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2924))

* Support COS Object SQL URL for Cloud Object Storage ([#2934](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2934))

* Migrate ibm_resource_group datasource to platform SDK and support other computed attributes ([2936](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2936))

* Migrate authorisation policy to platform-go-SDK ([#2926](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2926))

* Filter instances datasource by instance group ([#2947](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2947)) 

* Enable filter for security group by VPC ([#2982](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2982)) 

* Deprecate creation of KMS policies from ibm_kms_key resource ([#2832](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2832))

BUGFIXES

* Session timeout to 15 minutes and terraform cannot refresh the token ([2892](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2892))
* Added missing rule_id in tcp/udp and vpcId in sg filtering on ibm_is_vpc datasource ([2855](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2855))
* Return error statement if security group with name is not found in datasource ([#2932](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2932))
* Return error during Subnet detroy failure ([#2779](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2779))
* Update pi_instance virtual cores and processors after checking the capability ([2939](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2939))
* Reduce delay from 1m to 10s on cmr_order resource ([#2944](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2944))
* ibm_is_floating_ip target change results in a failure ([#2911](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2911))
* Added a new resource check on name update which was tainting the is_instance resource ([#2970](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2970))
* Fix Error handling during `waitForVpcClusterIngressAvailable` and fix condition check while setting `disable_public_service_endpoint` ([2971](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2971))
* VPE gateway creation suppresses error message ([#2923](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2923))
* Doc fix catalog management ([2965](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2965))
* Doc fix for DNS Forwarding rule ([2919](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2919))
* Remove location validaton on ibm_satellite_location resource ([2931](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2931)) 
* Set Id immediately after Creation in ibm_database resource ([2928](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2928))
* Fix: delete topic does not fail when topic does not exist ([#2976](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2976)) 


## 1.29.0 (Jul30, 2021)
FEATURES:
* Support VPC datasource
    - **DataSources**
        - ibm_is_vpcs
* Support Private DNS Custom Resolver
    - **DataSources**
        - ibm_dns_custom_resolvers
        - ibm_dns_custom_resolver_forwarding_rules
    - **Resources**
        - ibm_dns_custom_resolver
        - ibm_dns_custom_resolver_forwarding_rule
        - ibm_dns_custom_resolver_location
* Support Hyper Protect Crypto Service
     - **DataSources**
        - ibm_hpcs
    - **Resources**
        - ibm_hpcs

ENHANCEMENTS

* Support `status_reasons` attribute in ibm_is_instance ([2900](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2900))

* Support pagination for dedicated hosts datasources ([2906](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2906))

* Support Limit Feature for Keys ([2538](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2538))

* Support `address` in `ips` attribute for ibm_is_virtual_endpoint_gateway ([2913](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2913))

* Support `worker_pool_id` attribute in container worker pool ([2910](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2910))


BUGFIXES
* Fix doc format for ibm_cos_bucket ([2891](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2891))
* Fix Creating a resource group with IBMCLOUD_VISIBILITY=private seems to require a call to public IAM API ([2890](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2890))
* Fix the tainteffects ([2862](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2862))
* Fix ibm_satellite_cluster resource - kube_version upgrade, provider never returns ([2827](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2827))
* Fix subnet destroy fails immediately after cluster destroy unless delay added ([2779](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2779))
* Fix return error when topic creation failed ([2912](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2912))
* Fix deprecated field in ibm_container_worker_pool doc example ([2915](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2915))
* Update example of Hyper Protect DBaaS ([2917](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2917))

## 1.28.0 (Jul16, 2021)
FEATURES:
* Support VPC address prefixes
    - **DataSources**
        - ibm_is_vpc_address_prefixes
* Support VPC Instance Volume attachment
    - **Resources**
        - ibm_is_instance_volume_attachment
    - **DataSources**
        - ibm_is_instance_volume_attachment
        - ibm_is_instance_volume_attachments
* Support VPC Snapshots
    - **Resources**
        - ibm_is_snapshot
    - **DataSources**
        - ibm_is_snapshot
        - ibm_is_snapshots

ENHANCEMENTS

* Support to provision an instance from `instance_template` and instance boot volume from `snapshot` ([2672](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2672))

* Support Worker Pool Taints for IBM Cloud Clusters ([#2862](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2862))

* Clean up the reference for Gen1 code from all VPC resources and datasources

BUGFIXES
* ibm_is_security_group datasource does not find existing SGs ([#2868](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2868))
* Apply for ibm_iam_access_group_members fails with 404 when using IBMCLOUD_VISIBILITY=private ([#2828](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2828))
* Removed IPV6 References in vpc resources ([2697](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2697))
* IBM Cloud Provider - ibm_satellite_location silent correction ([2724](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2724))
* fix: set required cluster_name_id attribute ([2836](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2836))
* Doc fix for IAM policies and COS bucket ([2861](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2861))


## 1.27.2 (Jul09, 2021)

ENHANCEMENTS
* Support `vpc_name` attribute in ibm_is_subnet datasource ([#2783](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2783))

* Support `authentication_key` in Direct Link Gateway resource and datasource ([#2792](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2792))

* Support `pi_affinity_policy`, `pi_affinity_volume` and `pi_affinity_instance` for PI Volume ([#2800](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2800))

* Support `pi_migratable` argument for ibm_pi_instance resource and deprecate `migratable` ([#2801](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2801))

* Support `pi_storage_type` argument for ibm_pi_instance resource ([#2797](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2797))

* Added a fix to skip instance volumes setting in ibm_is_instance resource ([#2798](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2798))

BUGFIXES
* Fix the documentation for ibm_is_lb_listener ([#2790](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2790))

* Fix the documentation for ibm_cr_namespaces and ibm_cr_retention_policy ([#2821](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2821))

* Fix the documentation for ibm_kms_key_alias ([#2825](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2825))

* Remove the guides folder from website docs ([#2833](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2833))

* Fix: Panics on Import ([#2585](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2585))

* Fix host script for GCP ([#2802](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2802))

* Fix: Resource group crash for no default resource group ([#2809](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2809))

* Fix Detach targets before deleting security group ([#2723](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2723))

* Fix ibm_database datasource for different instances with same name ([#2817](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2817))

* Fix: pi_public_network datasource crash ([#2801](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2801))

* Fix: added failure check and tainting on failure for ibm_is_instance ([#2812](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2812))

## 1.27.1 (Jun27, 2021)
BUGFIXES

* Add retries on reading a IAM Policy ([#2788](https://github.com/IBM-Cloudterraform-provider-ibm/issues/2788))

## 1.27.0 (Jun25, 2021)
FEATURES:
* Support VPC network rules
    - **Resources**
        - ibm_is_network_acl_rule
    - **DataSources**
        - ibm_is_network_acl_rule
* Support CIS filter
    - **Resources**
        - ibm_cis_filter
    - **DataSources**
        - ibm_cis_filters

ENHANCEMENTS

* Support `install_plan`, `channel` and `wait_until_successful` arguments in ibm_cm_offering_instance resource ([#2745](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2745))

* Support provisioning of VPC images from source_volume ([#2682](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2682))

* Support filters to filter the VPC resources sshkeys, subnets, images, subnet, instances

* Support `hard_quota` argument for ibm_cos_bucket resource ([#2756](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2756))

* Add support creating gre tunnel connections on a transit gateway ([#2700](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2700))

* Support service roles in iam_role_actions ([#2746](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2746))

* Add retry to download the cluster config ([#2743](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2743))

* Add save to file feature to apikeys resource ([#2775](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2775))

* Migrate: Service Id resource to Platform-go-SDK ([#2560](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2560))

BUGFIXES
* Bug fix for Retrieving Policy with Key ([#2730](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2730))

* Fixed resource_tag crn validation issue ([#2749](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2749))

* Fix the updated of kube version for ROKS cluster ([#2754](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2754))

* Fix Documentation error for ibm_cis_rate_limit match.request.methods ([#2764](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2764))

* Fix ibm_cis_rate_limit Error: cis_id or zone_id not passed but it was passed in ([#2770](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2770))

* Fix the doc formats for VPC, Calssic Infrastructure services 

* Add doc link to role definition in IAM policy resources ([#2751](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2751))

* Plugin crashes if import is attempted without API Key being set ([#2729](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2729))

* wrong resource name in the doc link for ibm_iam_api_key ([#2736](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2736))

* Couldn't able to delete the Service Policy ([#2703](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2703))


## 1.26.2 (Jun15, 2021)
BUGFIXES
* Fix: Rollback the ibm_resource_tag resource ([#2718](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2718))

* Fixed typo that prevented host_attach_script from completing ([#2715](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2715))

* Fix InvalidBucketState: Versioning cannot be enabled for a bucket with expiration lifecycle actions ([#2727](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2727))

* Fix the empty ca_certificate for download of network cluster config ([#2732](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2732))

## 1.26.1 (Jun11, 2021)
Note : Don't use this release we have an issue with checksum errors for users when attempting to download the plugin. We have a new release v1.26.2

BUGFIXES
* Fix: Rollback the ibm_resource_tag resource ([#2718](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2718))

* Fixed typo that prevented host_attach_script from completing ([#2715](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2715))

* Fix InvalidBucketState: Versioning cannot be enabled for a bucket with expiration lifecycle actions ([#2727](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2727))

* Fix the empty ca_certificate for download of network cluster config ([#2732](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2732))

## 1.26.0 (Jun04, 2021)
FEATURES:
* Support AppConfiguration
    - **Resources**
        - ibm_app_config_environment
        - ibm_app_config_feature
    - **DataSources**
        - ibm_app_config_environment
        - ibm_app_config_environments
        - ibm_app_config_feature
        - ibm_app_config_features
* Support VPC instance group membership
    - **Resources**
        - ibm_is_instance_group_membership
    - **DataSources**
        - ibm_is_instance_group_memberships
        - ibm_is_instance_group_membership
* Support VPC instance group manager action
    - **Resources**
        - ibm_is_instance_group_manager_action
    - **DataSources**
        - ibm_is_instance_group_manager_action
        - ibm_is_instance_group_manager_actions
* Support Satellite Cluster
    - **Resources**
        - ibm_satellite_cluster
        - ibm_satellite_cluster_worker_pool
    - **DataSources**
        - ibm_satellite_cluster
        - ibm_satellite_cluster_worker_pool
    
* Support VPC Operating System
    - **DataSources**
        - ibm_is_operating_system
        - ibm_is_operating_systems

ENHANCEMENTS
* Update catalog management offering instance ([#2628](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2628))

* support object versioning feature for COS ([#2664](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2664))

* Support data volume in instance group ([#2673](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2673))

* Added VMWare Host Provider to Satellite Host Attach Script ([#2688](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2688))

* Support: tags in resource_instance datasource ([#2691](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2691))

* Add dbaas example ([#2683](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2683))


BUGFIXES
* Fix: Remove resource from statefile for pending_reclamation state  ([#2643](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2643))

* ibm_container_cluster_config data source is broken in 1.25 release on windows ([#2651](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2651))

* ibm_is_vpn_gateways data source documentation is wrong ([#2629](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2629))

* Documentation/example covering how to encrypt boot volume on ibm_is_instance ([#2577](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2577))

* Container Registry Documentation Error ([#2685](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2685))

*  Create Key With Writer Role Assignment ([#2619](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2619))

* Fix the NotFound Checks if the resource is still provisioning ([#2613](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2613))

* A link should be fixed in ibm_kms_key ([#2704](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2704))

* DocUpdate: Add Timeout Blocks to vpc cluster ([#2706](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2706))

* ibm_schematics_workspace data source unable to parse env_values  ([#2708](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2708))


## 1.25.0 (May18, 2021)
FEATURES:
* Support Resource Tag Management
    - **Resources**
        - ibm_resource_tag
    - **DataSources**
        - ibm_resource_tag
* Support VPC dedicated host disk management
    - **Resources**
        - ibm_is_dedicated_host_disk_management
    - **DataSources**
        - ibm_is_dedicated_host_disk
        - ibm_is_dedicated_host_disks
* Support IAM User API Key
    - **Resources**
        - ibm_iam_api_key
    - **DataSources**
        - ibm_iam_api_key
* Support VPC endpoint target gateways
    - **DataSources**
        - ibm_is_endpoint_gateway_targets
* Support VPC security group target management
    - **Resources**
        - ibm_is_security_group_target
    - **DataSources**
        - ibm_is_security_group_target
        - ibm_is_security_group_targets
* Support COS Bucket Object
    - **Resources**
        - ibm_cos_bucket_object
    - **DataSources**
        - ibm_cos_bucket_object
* Support container Registry Retention Policy
    - **Resources**
        - ibm_cr_retention_policy
	

ENHANCEMENTS
* Add the capabilities for offering speeds which provides the bmetered and unmetered
billing options for that offering speed ([#2584](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2584))

* Support for dedicated host/group in instance template and instance ([#2579](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2579))

* location, managed_from and resource_group_id mark these attibutes in ibm_satellite_location either DiffSuppress or throw back an error ([#2567](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2567))

* Add max allowed sessions to account settings ([#2610](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2610))

* Add vcpus and memory to instance profiles data source ([#2492](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2492))

* Support creating IAM policies with operator ([#2533](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2533))

* Mark ibm_container_cluster_config `admin_certificate`, `ca_certificate`, `token` attributes as sensitive ([#2622](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2622))

* Bump up go version: 1.16 ([#2600](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2600))


BUGFIXES
* Fix the deletion of instance group due to load balancer status ([#2547](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2547))

* No way to make a vpc_address_prefix default on vpc when using manual address_prefix_management ([#2282](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2282))

* ibm_is_image data source silently fails if the image is not available ([#2587](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2587))

* ibm_is_image data source does not provide a warning if the image is deprecated ([#2588](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2588))

* update azure script to grow root volume group ([#2621](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2621))

* instance group manager policy synchronization ([#2635](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2635))

* Fix the db task timeout ([#2607](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2607))

* KMS keys created with endpoint_type = "public" regardless of the actual setting ([#2482](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2482))

* Error finding VLAN order: couldn't find resource (21 retries) ([#2613](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2613))

* Terraform import for routingtable fails ([#2580](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2580))


## 1.24.0 (May04, 2021)
FEATURES:
* Support VPC instance disk management
    - **Resources**
        - ibm_is_instance_disk_management
    - **DataSources**
        - ibm_is_instance_disk
        - ibm_is_instance_disks
	
ENHANCEMENTS
* Support resize of VPC instance ([#2448](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2448))
* Support Load balancer Parameter based routing ([#2518](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2516))
* Support horizontal scaling on database with new arguments node_count, node_memory_allocation_mb, node_disk_allocation_mb, node_cpu_allocation_count ([#2313](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2313))
* Support request_metrics_enabled for COS Bucket metric monitoring ([#2530](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2530))
* Support virtual endpoint gateway as target to subnet reserved IP ([#2521](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2521))

BUGFIXES
* Creating ibm_pi_key fails everytime with context deadline exceeded ([#2527](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2527))
* Fix diff on resource key parameters ([#2182](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2182))
* Fails to create PTR records causing Terraform crash ([#2535](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2535))
* Fix crash for VPC instance group manager ([#2554](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2554))
* VPC network ACL rule ICMP does not set type ([#2559](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2559))
* Conflict with exec.image and exec.code/exec.code_path (can't use custom docker images) ([#2556](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2556))

## 1.23.2 (Apr20, 2021)
ENHANCEMENTS
* Add support for COS retention policy ([#1880](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1880))
* Add support for private_address for VPN gateway ([#2282](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2382))
* List all certificates in a certificate manager instance ([#2358](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2358))
* Enhance description for attribute reference ([#2475](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2475))
*  Add support for regional ca-tor COS bucket ([#2483](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2483))
BUGFIXES
* Fix the broken links for classic infrastructure bare metal ([#2481](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2481))
* Fix cis primary certificate crash ([#2490](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2490))
* Fix ibm_satellite_location: cannot specify resource group ([#2499](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2499))
* Fix ibm_satellite_location resource doesn't work correctly to ensure that resource is created / deleted appropriately ([#2497](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2497))
* Fix invalid example for ibm_iam_account_settings ([#2484](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2484))
* Fix the documentiaon for VPC reserved IP ([#2512](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2512))

## 1.23.1 (Apr07, 2021)
ENHANCEMENTS
* Add support to retry the update of patch version ([#2379](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2379))
* Add gateway_connection argument for VPC VPN gateway Connection ([2270](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2270))

BUGFIXES
* Fix the crash for resource key ([#2462](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2462))
* Change the order to place to use billing_order ([#554](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/554))

## 1.23.0 (Apr02, 2021)
FEATURES:
* Support Catalog Management
    - **Resources**
        - ibm_cm_offering_instance
        - ibm_cm_catalog
        - ibm_cm_offering
        - ibm_cm_version
    - **DataSources**
        - ibm_cm_catalog
        - ibm_cm_offering
	    - ibm_cm_version
		- ibm_cm_offering_instance
* Support IAM Account Management
    - **Resources**
        - ibm_iam_account_settings
    - **DataSources**
    	- ibm_iam_account_settings

* Support Enterprise Management
    - **Resources**
        - ibm_enterprise
        - ibm_enterprise_account_group
        - ibm_enterprise_account
    - **DataSources**
    	- ibm_enterprises
        - ibm_enterprise_account_groups
        - ibm_enterprise_accounts

BUGFIXES
* Fix the provision of classic Infrastructure VM to apply sshkeys, imageid and script ([#2448](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2448))
* Fix documentation updates ([#2443](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2443))
* Fix Dedicated host with status 'failed' throws error during destroy ([#2443](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2446))
* Fix while creating a DL Connect gateway do not wait for gateway to be provisioned for few providers ([#2458](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2458))

## 1.22.0 (Mar30, 2021)

FEATURES:

* Support VPC dedicated hosts
    - **Resources**
        - ibm_is_dedicated_host
		- ibm_is_dedicated_host_group
    - **DataSources**
        - ibm_is_dedicated_host
		- ibm_is_dedicated_hosts
		- ibm_is_dedicated_host_profile
		- ibm_is_dedicated_host_profiles
		- ibm_is_dedicated_host_group
		- ibm_is_dedicated_host_groups
* Support VPC reserved IP
     - **Resources**
        - ibm_is_subnet_reserved_ip
     - **DataSources**
        - ibm_is_subnet_reserved_ip
        - ibm_is_subnet_reserved_ips
* Support Push Notification chrome web
    - **Resources**
        - ibm_pn_application_chrome
    - **DataSources**
        - ibm_pn_application_chrome
        
* Support Key Management Alias and Rings
    - **Resources**
        - ibm_kms_key_alias
        - ibm_kms_key_rings
    - **DataSources**
        - ibm_kms_key_rings

* Support for reading secrets from IBM Cloud Secrets Manager
    - **DataSources**
        - ibm_secrets_manager_secrets
        - ibm_secrets_manager_secret

* Support Schematics
    - **Resources**
        - ibm_schematics_workspace
        - ibm_schematics_action
        - ibm_schematics_job
    - **DataSources**
        - ibm_schematics_action
        - ibm_schematics_job

* Support Observability
     - **Resources**
        - ibm_ob_logging
        - ibm_ob_monitoring

* Support for Satellite 
    - **Resources**
        - ibm_satellite_location
        - ibm_satellite_host
    - **DataSources**
        - ibm_satellite_location
        - ibm_satellite_attach_host_script

* Support for CIS Cache setting
    - **Datasource**
        - ibm_cis_cache_settings

PROVIDER

* Support `visibility` argument to control the visibility to IBM Cloud endpoint.

* `generation` argument is depreated. By default the provider targets to IBM Cloud VPC Infrastructure.


ENHANCEMENTS

* Support added DH group 19 and sha 512 for IKE and IPSec Policy([#2361](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2361))

* Support `delegate_vpc` action for VPC routing table ([#2355](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2355))

* Support tags for IBM Cloud VPC security group ([#2353](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2353))

* Support for renaming of default Network ACL, Security Group and Routing Table ([#2216](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2353))

* Support for allow control of Security Groups on Load Balancer ([#2324](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2324))

* Support recover the public key from an SSH key in IBM Cloud VPC ([#2388](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2388))

* Support ibm_is_vpc_routing_table_route to accept a VPN connection ID as next_hop ([#2270](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2270))

* Support to provision classic Infrastructure Virtual instance using quoteID ([#2433](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2433))

* Support `serve_stale_content` for CIS Cache settings ([#2219](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2219))

* Support filtering of subnets based on metro for IKS kubernetes Cluster ([#2403](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2403))

* Support `alias` argument to filter the keys in ibm_kms_key and ibm_kms_keys datasources and aliases attribute ([#2293](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2293))

* Support `key_ring_id` in ibm_kms_key resource and datasources ([#2378](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2378))

BUGFIXES

* Fix increase in panics while refreshing resources for cos_bucket ([#2373](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2373))

* Fix Cloud Function action runtimes version ([#2424](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2424))

* Fix COS buckets allow modifying key_protect after creation, but they should not ([#2310](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2310))


## 1.21.2 (Mar15, 2021)

PROVIDER:
* Updgrade Terraform SDK to v2

FEATURES:

* Support checksum argument for VPC Images ([#2227](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2161))

* Support iam_id argument cross Account iam_service_policy ([#2331](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2331))

* Support tags argument for VPC subnet ([#2321](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2332))

* Support tags argument for VPC network acl ([#2343](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2343))

* Add resource schema timeouts for classic infrastructure compute VM ([#2291](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2291))

* Support accept_proxy_protocol argument for vpc loadbalancer listener ([#2325](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2325))


BUGFIXES

* Fix logging not supported for VPC Network Loadbalancer ([#2332](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2332))

* Fix addons not being enabled post-cluster creation ([#2346](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2346))

* Fix ibm_iam_user_policy data source produces no results ([#2312](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2312))


## 1.21.1 (Mar03, 2021)

FEATURES:

* Support sort argument for IAM service policies and IAM user policies ([#2227](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2227))

* Support default_routing_table attribute for VPC resource and datasource ([#2286](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2286))

* Support logging argument for VPC load balancer ([#2228](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2228))

* Add transactionID for IAM authentication error messages ([#2304](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2304))

BUGFIXES

* Fix the provision of instance template with boot volume ([#2205](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2205))

* Fix ibm_resource_key is not tainted by change to role ([#2182](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2182))

* Fix ibm_cis, ibm_database, ibm_resource_instance not tainted by change to resource groupID ([#2297](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2297))

* Fix ibm_container_addons resource not detecting version change ([#2295](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2295))

* Fix error when trying to use data source ibm_container_cluster for existing lite IKS ([#2300](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2300))

## 1.21.0 (Feb12, 2021)

FEATURES:

* Support datasource for VPC volume profiles `ibm_is_volume_profile`, `ibm_is_volume_profiles`

* Support datasource to list power instance catalog images `ibm_pi_catalog_images`

ENHANCEMENTS

* Support `lunid` attribute for block classic block storage ([#1491]https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1491)

* Support `HTTP_COOKIE` session_affinity for lbass ([#2218](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2218)) 

* Support `auto_delete_volume` argument to delete data volumes of VPC instance [#646](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/646)) 

* Support `wait_till` argument for ibm_containar_cluster to control the behaviour of waiting for cluster ([#2232]https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2232)

* Enable retries on authnetication failures ([#2248]https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2248)


BUGFIXES

* Fix the nil pointer exception on lbs of vpc_cluster destroy ([#2226](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/22276)) 

* Fix the nil pointer exception on is_instance_group ([#2247](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2247)) 

* Fix the nil pointer on iam_service_api_key ([#2259](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2259)) 

* Fix the validation for LB listener policy rule ([#2257](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2257)) 

* Fix the patch_version update for kubernetes clusters ([#2217](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2217)


## 1.20.1 (Jan27, 2021)

BUGFIXES
Fix the regression issue provisioning of contianer clusters ([#2206](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2206)) 

## 1.20.0 (Jan25, 2021)

FEATURES:

* Support directlink provider gateway resource `ibm_dl_provider_gateway`

* Support directlink provider gateways, ports datasource `ibm_dl_provider_gateways`, `ibm_dl_provider_ports`

ENHANCEMENTS

* Support provision file storage size from 10TB to 13TB ([#2158](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2158))

* Support for pod-subnet and service-subnet for `ibm_container_cluster` resource ([#1196](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1196))

* Support the ability to retrieve the instances in a specific VPC ([#1961](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1961))

* Support for patch update for cluster worker nodes ([#1978](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1978)) 

* Support architecture attribute for VPC instance profiles ([#2002](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2002)) 

BUGFIXES

* Fix the nil pointer exception for cos bucket import scenario ([#2151](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2151))

* Fix Transit gateway Connection creation fails for cross account ([#2170](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2170))

* Fix Terraform crash when subnet not found ([#2058](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2058))

* Fix VPC LB creation with count greater than 1 ([#2168](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2168))

# 1.19.0 (Jan08, 2020)

FEATURES:
* Support Contianer Registry resource and datasource `ibm_cr_namespace`, `ibm_cr_namespaces` ([#2119](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2119))

* Support reset APIkey for cluster `ibm_container_api_key_reset` ([#2118](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2118))

* Support APIkey for serviceID `ibm_iam_service_api_key` ([#666](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/666))

ENHANCEMENTS:

* Move next_hop from optional to required for `ibm_is_vpc_routing_table_route` ([#2141](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2141))

* Support jp-osa endpoints for `ibm_cos_bucket` [#2149](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2149))

* Support crn in target attribute for `ibm_is_virtual_endpoint_gateway` ([#2147](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2147))

## 1.18.0 (Dec22, 2020)

FEATURES:

* Support CIS Certificate resources `ibm_cis_certificate_order`, `ibm_cis_certificate_upload`

* Support CIS Certificate datasources `ibm_cis_certificates`, `ibm_cis_custom_certificates`

* Support CIS DNS Records import and export `ibm_cis_dns_records_import`, `ibm_cis_dns_records`

* Support virtual private endpoint gateways `ibm_is_virtual_endpoint_gateway`, `ibm_is_virtual_endpoint_gateway_ip` resources and `ibm_is_virtual_endpoint_gateways`, `ibm_is_virtual_endpoint_gateway_ips`, `ibm_is_virtual_endpoint_gateway` datasources


ENHANCEMENTS:

* resource: Support `labels` argument and updates for kubernetes clusters ([#2109](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2109)) 

* resource: Support `namespace` argument and `persistence` and `status` attributes ([#2097](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2097)) 



BUGFIXES:

* Fix an ibm_resource_key that is removed outside of Terraform is not being recreated ([#2125](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2125))

* Fix cluster addon fails on apply after timeout ([#2129](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2129))


## 1.17.0 (Dec10, 2020)

FEATURES:

* Support CIS WAF resources `ibm_cis_range_app`, `ibm_cis_waf_package`, `ibm_cis_waf_group`, `ibm_cis_waf_rule`

* Support CIS WAF datasources `ibm_cis_range_apps`, `ibm_cis_waf_packages`, `ibm_cis_waf_groups`, `ibm_cis_waf_rules`

ENHANCEMENTS:

* resource: Support `force_delete` argument for ibm_cos_bucket ([#2017](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2017)) 

* resource: Move `bgp_base_cidr` as optiona argument ([#2087](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2087)) 

* resource: Support `expire_rule` argument for ibm_cos_bucket ([#1590](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1590))

* resource: Add validate function for resoure_instance_id argument to ibm_cos_bucket ([#2103](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2103))

* resource: Support Route and Profile based VPN gateways ([#2094](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2094))


BUGFIXES:

* Fix users not found when adding to access group ([#2034](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2034))

* Set the `instance_id` in ibm_kms_key ([#2106](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2106))

* Fix multiple cis_domain leads into inconsistency ([#2086](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2086))


## 1.16.1 (Dec02, 2020)

BUGFIXES:

* Fix issue when trying to delete a ibm_container_alb_cert ([#2067](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2067))


## 1.16.0 (Nov30, 2020)

FEATURES:

* Support VPC Routing Table `ibm_is_vpc_routing_table` and VPC Routing Table Route  `ibm_is_vpc_routing_table_route` resources ([#1395](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1395))

* Support `ibm_is_vpc_default_routing_table`, `ibm_is_vpc_routing_tables` and `ibm_is_vpc_routing_table_routes` datasources ([#1395](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1395))


ENHANCEMENTS:

* resource: Extend CIS firewall resource to support `access_rules` and `ua_rules` ([#2025](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2025)) 

* resource: Support anti-spoofing `allow_ip_spoofing` for ibm_is_instance ([#1396](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1396)) 

* resource: Support `routing_table` and `ip_version` agruments for ibm_is_subnet ([#1395](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1395))

* data: Support `policies` attribute for ibm_kms_keys and ibm_kms_key ([#1928](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1928))

* resource: Support `number_of_invited_users` and `invited_users` attribute for ibm_iam_user_invite ([#2053](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2053))

BUGFIXES:

* Fix the upgrade of kube_version for master and worker nodes of cluster ([#1952](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1952))

* Fix issue when trying to provision a new ibm_container_alb_cert ([#2067](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2067))

## 1.15.0 (Nov24, 2020)

FEATURES:

* Support for subnet network interface attachment `ibm_is_subnet_network_acl_attachment` resource ([#1941]

* Support `ibm_cis_routing` resource ([#1991](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1991))

* Support `ibm_cis_cache_settings` resource ([#1995](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1995))

* Support `ibm_cis_global_load_balancers` datasource ([#1981](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1981))

* Supoort `ibm_cis_custom_page`resource and `ibm_cis_custom_pages` datasoruce ([#1997](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1997))

* Support `ibm_dns_glb`, `ibm_dns_glb_monitor`, `ibm_dns_glb_pool` resource for IBM Cloud PDNS Service ([#1887](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1887))

* Support `ibm_dns_glbs`, `ibm_dns_glb_monitors`, `ibm_dns_glb_pools` datasources for IBM Cloud PDNS Service ([#1887](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1887))

ENHANCEMENTS:

* resource: Support `public_ip` attribute in ibm_pi_network_port resource ([#1930](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1930)) 

* resource: Support encrypted images `encrypted_data_key` and `encryption_key` arguments in ibm_is_image resource ([#1938](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1938)) 

* resource: Support archive rule `archive_rule` for ibm_cos_bucket ([#1950](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1950)) 

* resource: Support Polcies for ibm_kms_key ([#1928](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1928))


* data: Support `list_bounded_services` argument for ibm_container_cluster ([#2051](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2051))

BUGFIXES:

* Fix provision of cloud funciton resources ([#837](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/837))

* Fix the destroy of ibm_pi_instance wait ([#2047](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/2047))

## 1.14.0 (Oct28, 2020)

FEATURES:

* Support for subnet network interface attachment `ibm_is_subnet_network_acl_attachment` resource ([#1941](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1941))

* Support for CIS tls settigns `ibm_cis_tls_settings` resource ([#1954](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1954))

* Support `ibm_cis_origin_pools` datasource ([#1959](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1959))

ENHANCEMENTS:

* resource: Support additional domain settings (max_upload, cipher, minify, security_header, mobile_redirect, challenge_ttl, dnssec, browser_check) for `ibm_cis_domain_settings` ([#1939](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1939))

* resource: Support `expiration_date` argument to ibm_kms_key resource ([#1967](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1967))

* resource: Support `wait_for_worker_update` argument to `ibm_container_cluster` and `ibm_container_voc_cluster` to control the upgrade of worker nodes ([#1969](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1969))

* data: Support `cert_file_path` attribute to `ibm_database` ([#1985](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1985))

BUGFIXES

* Remove forcenew for IS instance group ([#1951](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1951))

* Support resource instance's parameter to be an array type ([#1953](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1953))

* Fix the tags attachemnt for ibm_databse resource ([#1971](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1971))

* Fix changing allowed_ip to a list of IPs to nothing leads to an error when configuring a COS bucket ([#1661](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1661))

* Fix the update of VPC worker nodes kube verison ([#1952](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1952))


## 1.13.1 (Oct07, 2020)

ENHANCEMENTS:

* resource: Support endpoint_type argument and endpoint environmental variable for COS bucket ([#1945](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1945))

* resource: Support Direct Link Connect Type([#1927](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1927))

* doc: Update supported parameters for Event Streams([#1946](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1946))

BUGFIXES

* Fix the nil pointer exception for transist gateway delete ([#1943](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1943))

## 1.13.0 (Oct01, 2020)

FEATURES:

**VPC NLB Feature**: 
* Support for provisioning of NLB load balancers ([#1937](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1937))
    * data/ibm_is_lb_profiles
    * data/ibm_is_lbs

**CIS Edge Functions**: 
* Support for CIS Edge Functions ([#1873](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1873))
    * resource/ibm_cis_edge_functions_actions
    * resource/ibm_cis_edge_functions_trigger
    * data/ibm_cis_edge_functions_actions
    * data/ibm_cis_edge_functions_triggers


ENHANCEMENTS:

* datasource: Support `pools` attribute for is_lb datasource ([#1895](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1895))

* resource: Support of renew certificate in certificates manager ([#1909](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1909))

* resource: Support update of parameters for resource instance ([#1705](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1705))

* resource: Support for NLB load balancers in VPC ([#1937](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1937))

* resource: Migrate HPCS endpoints to cloud.ibm.com domain ([#1932](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1932))

* resource: Support retry of VPC instance to recover from a perpetual "starting" or "stopping" state by using "force_recovery_time" argument ([#1934](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1934))

* resource: Support customer health check request headers for Cloud Internet Services ([#1844](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1844))

* resource: Support ibm_iam_service_id (data / resouce) should return iam_id ([#1820](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1820))

* resource: Support ICD Service endpoint doesn't exist for region: "che01" ([#1894](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1894))

BUGFIXES

* Fix the instance template destroy error ([#1886](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1886))

* Fix the delete of ibm_cdn resource ([#1925](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1925))

* Fix the provision of free cluster ([#1901](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1901))

* Fix the ibm_container_addons not working on other resource_group !=Default ([#1920](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1920))

* Fix the crash of ibm_is_subnet datasource with empty identifier ([#1933](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1933))

* Fix Instance Group/AutoScale Max count should be 1000 not 100 ([#1889](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1889))

* Fix ibm_pi_instance not failing on ERROR state ([#1879](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1879))

## 1.12.0 (Sep14, 2020)

FEATURES:

**VPC Flow Logs**: 
* Support for IBM Cloud VPC Flow Logs ([#1356](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1356))
    * resource/ibm_is_flow_logs
    * data/ibm_is_flow_logs

**VPC Auto Scale**: 
* Support for IBM Cloud VPC Auto Scale ([#1357](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1357))
    * resource/ibm_is_instance_group
    * resource/ibm_is_instance_group_manager
    * resource/ibm_is_instance_group_manager_policy
    * resource/ibm_is_instance_template
    * data/ibm_is_instance_group
    * data/ibm_is_instance_group_manager
    * data/ibm_is_instance_group_managers
    * data/ibm_is_instance_group_manager_policies
    * data/ibm_is_instance_group_manager_policy
    * data/ibm_is_instance_templates
    * data/ibm_is_instance_profiles
    * data/ibm_is_instance_profile

**Power Instance**:
* Support for IBM Cloud Power Instance network port attachement and snapshot ([#1867](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1867))
    * resource/ibm_pi_snapshot
    * resource/ibm_pi_network_port_attach

**Cluster Addons**:
* Support for addons for container cluster ([#721](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/721))
    * resource/ibm_container_addons
    * data/ibm_container_addons

* data/ibm_is_lb: Support for ibm_is_lb [#1849](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1849))

* data/ibm_container_alb: Support for ibm_container_alb [#1850](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1849))

* data/ibm_container_alb_cert: Support for ibm_container_alb_cert [#1850](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1849))

* data/ibm_container_bind_service: Support for ibm_container_bind_service [#1850](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1849))


ENHANCEMENTS:

* resource: Support key protect configuraton for Container clusters ([#673] (https://github.com/IBM-Cloud/terraform-provider-ibm/issues/673))

* resource: Support delete of PVC Storage for Container clusters([#1847] (https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1847))

BUGFIXES

* Fix the diff on classsic VM instance ([#1828] (https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1828))

## 1.11.2 (Sep 01, 2020)

ENHANCEMENTS:
* resource: Support ProxyFromEnvironment for honouring for HPCS service.
* resource: Support auto scaling for IBM Cloud database service.


## 1.11.1 (Aug 26, 2020)

ENHANCEMENTS:

* resource: Assign IP address to a VSI on provisioning for VPC ([#1830](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1830))

* resource: Support kr-seo region for database ([#1831](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1831))

BUGFIXES

* Fix ibm_is_vpc datasource for Gen1 ([#1834](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1834))

* Fix provision of ibm_pi_instance ([#1833](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1833))

* Fix provision of ibm_pi_network_port ([#1823](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1823))

## 1.11.0 (Aug 24, 2020)

FEATURES:

* data/ibm_is_public_gateway: Support for ibm_is_public_gateway ([#1745](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1745))
* data/ibm_is_floating_ip: Support for ibm_is_floating_ip ([#1794](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1794))

ENHANCEMENTS:

* resource: Allow configuration of the key used to encrypt IBM Cloud Databases backups ([#1761] (https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1761))

* Support customer managed volume encryption for VPC Nextgen ([#1673] (https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1673))

* Support virtual cores capability for power instance instance ([#1798] (https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1798))

* Support for interconnecting two ibm cloud functions by target_url ([#1526] (https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1526))

*Support for cross account Transist Gateway ([#1021] (https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1021))

BUGFIXES

* ibm_tg_gateway delete is not complete when it has reported deleted ([#1783] (https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1783))

* Fix the IAM IP address restriction for invited user ([#1780] (https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1780))

* VPC instance datasource failing to fetch the correct instance ([#1801] (https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1801))

* Fix ibm_is_network_acl Validate name input field of ACL Rules ([#1262] (https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1262))

* Fix for datasource ibm_cis_domain only retrieve 20 domains ([#1804] (https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1804))

* Fix iam_access_group_policy with addition of account_management doesn't apply ([#1551] (https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1551))

* Fix error received - Current user does not have access to team directory ([#1536] (https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1536))

## 1.10.0 (Aug 06, 2020)

FEATURES:

**Transist Gateway**: 
* Support for Trasist Gateway Service ([#1021](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1021))
    * resource/ibm_tg_gateway
    * resource/ibm_tg_connection
    * data/ibm_tg_gateway
    * data/ibm_tg_gateways
    * data/ibm_tg_locations
    * data/ibm_tg_location

**DirectLink Gateway**: 
* Support for DirectLink Gateway Service ([#1349](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1349))
    * resource/ibm_dl_gateway
    * resource/ibm_dl_virtual_connection
    * data/ibm_dl_gateways
    * data/ibm_dl_offering_speeds
    * data/ibm_dl_port
    * data/ibm_dl_ports
    * data/ibm_dl_gateway
    * data/ibm_dl_locations
    * data/ibm_dl_routers

**CloudFunction**		
* Support for Cloud Function Namespace ([#682](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/682))
    * resource/ibm_function_namespace
    * data/ibm_function_namespace

**KMS (keyprotect/hpcs crypto)**
* Support for Key management (key protect/HPCS Crypto Service) ([#1353](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1353))
    * resource/ibm_kms_key
    * data/ibm_kms_key
    * data/ibm_kms_keys

**IAM User Setting**
* Support for IAM User Management settings ([#1780](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1780))
    * resource/ibm_iam_user_settings
    * data/ibm_iam_users
    * data/ibm_iam_user_profile

**Event Stream**
* Support for IBM Event Stream Topic ([#1781](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1781))
    * resource/ibm_event_streams_topic
    * data/ibm_event_streams_topic

* data/ibm_is_security_group: Support for ibm_is_security_group ([#1223](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1223))

* data/ibm_container_worker_pool: Support for ibm_container_worker_pool ([#1751](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1751))

* data/ibm_container_vpc_cluster_worker_pool: Support for ibm_container_vpc_cluster_worker_pool ([#1773](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1773))

* data/ibm_container_vpc_cluster_alb: Support for ibm_container_vpc_cluster_alb  ([#1775](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1775))

* data/ibm_certificate_manager_certificate: Support for ibm_certificate_manager_certificate  ([#1679](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1679))

ENHANCEMENTS:

* resource: Support configure geo routes in cis_glb ([#985] (https://github.com/IBM-Cloud/terraform-provider-ibm/issues/985))

* data: Retrieve icd disk encryption details ([#1742] (https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1742))

* resource: auto_renew_enabled support for ibm_certificate_manager_order ([#1657] (https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1657))

* resource: ibm_container_alb_cert destroy not synchronous ([#1712] (https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1712))


BUGFIXES

* Fix ibm_container_vpc_worker_pool resource forces a replace if resource_group_id not set ([#1748] (https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1748))

* Fix IBM Container VPC Cluster with no default worker pool crashes on destroy ([#1733] (https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1733))

## 1.9.0 (July 22, 2020)

FEATURES:

* data/ibm_is_instances: Support for ibm_is_instances ([#1454](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1454))
* data/ibm_is_instance: Support for ibm_is_instance ([#1454](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1454))

ENHANCEMENTS:

* resource: Support ibm_function_action, ibm_function_package, ibm_function_trigger, ibm_function_rule resources for IAM and CF based namespace ([#837](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/837))
**Note** - The provider level argument `function_namespace` is deprecated.The namespace is a required argument part of the function resources. The users need to update the templates to add the `namespace` argument to function resources.

* resource: Support update of adding additional zones to VPC cluster and worker pool resource ([#1546] (https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1546))

* Support update_all_workers flag to control the update of workers for VPC clusters ([#1681] (https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1681))

* resource: Support extension attribute for ibm_resource_instance ([#1686] (https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1686))

* resource: Support dashboard_url attributes for ibm_resource_instance ([#1682] (https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1682))

* resource: Support for update of key_protect_key parameter in ibm_database ([#1622] (https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1622))

* Support for resource synchronization of private dns permitted network ([#1674] (https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1674))

* data/ibm_resource_instance: Support guid attribute for ibm_resource_instance datasource ([#1724] (https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1724))

* resource: Support label argument for default worker pool ibm_container_cluster  ([#775] (https://github.com/IBM-Cloud/terraform-provider-ibm/issues/775))

BUGFIXES

* ibm_cis_domain_settings does not allow for Standard plans ([#1623] (https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1623))

* Fix the update of attachment of public gateway to VPC subnet ([#1626] (https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1626))

* Fix RHCOS via ibm_pi_instance timeout waiting for networ ([#1620] (https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1620))

* Fix Gateway enabled cluster recreated on every apply ([#1706] (https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1706))


## 1.8.1 (June 30, 2020)

ENHANCEMENTS:

* datasource: Support for aggregation of VPC images([#1580](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1580))

BUGFIXES

* resource: Fix the destroy of virtual instance with volumes
* resource: Fix the mutex of pdns resource records([#1601](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1601))

## 1.8.0 (June 23, 2020)

FEATURES:

* New Datasource: ([ibm_iam_roles](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1433))
* New Datasource: ([ibm_iam_role_actions](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1433))

ENAHANCEMENTS:

*resource: Support for provisioning ROKS on VPC Gen2 ([#1437](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1437))
*resource: Support for intergating firewall IP, activity tracker, metric monitoring to COS bucket ([#1487](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1487))
*data: Add new attribute `output_json` for ibm_scheamtics_output ([#1413](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1413))

BUGFIXES
*resource: Fix the list of network rules attached to network acl([#1547](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1547))

## 1.7.1 (June 11, 2020)

ENHANCEMENTS:

* resource/ibm_cis_domain_settings: Support additional domain settings ([#1475](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1475))
* resource/resource_ibm_certificate_manager_order: Added key_algorithm to order certificate in CMS ([#1512](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1512))
* data/ibm_is_vpc: Add zone name to data source vpc.subnets outputs ([#1450](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1450))

BUG FIXES:

* docs: Add documentation for is_instance and is_security_groups( [#1522](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1522))
* data/ibm_is_vpc: Regression on vpc_source_addresses support ([#1530](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1530))
* data/ibm_container_vpc_cluster: container_vpc_cluster fails for OCP VPC cluster with ALB error ([#1528](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1528))
* resource/ibm_iam_access_group_dynamic_rule : Fix dynamic Rule returning wrong resource ([#1535](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1535))
* resource/ibm_is_image: Fix the nil pointer exception for is image ([#1540](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1540))

## 1.7.0 (May 28, 2020)

ENHANCEMENTS:

* resource/ibm_cis_rate_limit: Support for CIS Rate Limiting ( [#1271](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1466))
* data/ibm_cis_rate_limit: Support for CIS Rate Limiting ( [#1271](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1466))

BUG FIXES:

* resource/ibm_is_security_group_rule: Gen1-Security Group Rule fix: allow 'Any' type for ICMP, TCP, UDP( [#1499](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1499))
* resource/ibm_dns_resource_record : Changes to lock resource record id and zone id ( [#1430](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1490))
* resource/ibm_is_vpc: Resource level Timeout updation and docs for vpc resources (is_vpc, is_vpc_route, is_vpn_gateway, is_vpn_gateway_connection )( [#1442](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1442))
* resource/ibm_is_vpn_gateway: Fix for deletion of VPN gateway( [#1495](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1495))
* resource/ibm_private_dns: Fix for provisioning of private dns resource records( [#1476](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1476 ))
* data/source_ibm_is_subnets: Fix for ibm_is_subnets output duplicates( [#1500](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1500 ))


## 1.6.0 (May 20, 2020)
FEATURES:

* New Resource: ([ibm_iam_custom_role](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1433))
* New Datasource: ([ibm_dns_zones](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/958))
* New Datasource: ([ibm_dns_permitted_networks](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/958))
* New Datasource: ([ibm_dns_resource_records](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/958))

ENAHANCEMENTS:

*resource: Adopt custom roles to IAM policies ([#1433](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1433))

## 1.5.3 (May 19, 2020)

ENAHANCEMENTS:

* resource :  Support for pi_pin_policy argument  in ibm_pi_instance ([#1469](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1469))
* resource :  Support for wait_till_albs argument  in ibm_container_workerpool_zone_attachment ([#1463](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1463))
* data : Support for state_store_json attribute in ibm_schematics_state ([#1411](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1411))

BUG FIXES:

* resource : Fix nil pointer if apikey not given for VPC ([#1427](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1427))
* data : CMS issuance_info update ([#1277](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1277))

## 1.5.2 (May 07, 2020)

ENAHANCEMENTS:

* resource : Support for entitlement argument for IKS Classic ROKS cluster (ibm_container_cluster) and worker pool(ibm_container_worker_pool)([#1350](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1350))

* resource : Support for source_resource_group_id and target_resource_group_id([#1364](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1364))

BUG FIXES:

* resource : Error deleting instance with data volume ([#1412](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1412))
* resource : Add force_new true for cidr argument of ibm_is_address_prefix ([#1416](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1416))
* resource : Fix import of ibm_container_cluster ([#1360](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1360))

## 1.5.1 (May 04, 2020)
BUG FIXES:

* resource : Fix VPC subnets created in incorrect resource group([#1398](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1398))

## 1.5.0 (April 29, 2020)
FEATURES:

* New Resource: ([ibm_is_lb_listener_policy_rule](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1147) )
* New Datasource: ([ibm_certificate_manager_certificates](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1277) )

ENAHANCEMENTS:

* resource : Support for auto-generate client_id and client_id for API gateway endpoint([#1390](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1390))
* resource: Support point_in_time_recovery_time and point_in_time_recovery_deployment_id arguments for ICD database([#1259](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1259))
* resource: Support for pending_reclamination for database and CIS instances ([#1242](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1242))


BUG FIXES:

* resource : Fix VPC Load Balancer resource ID is appended to Pool/Listener/Listener Policy ID ([#1359](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1359))
* resource : Fix domainID for CIS firewall resource ([#1201](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1201))
* resource : Fix the update of private dns resource record TTL ([#1331](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1331))
* resource : ibm_container_worker_pool_zone_attachment should wait for ALBs to finish in a new zone ([#1372](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1372))


## 1.4.0 (April 16, 2020)
  
NOTE :  For creating either vpc-classic (generation=1) or vpc-Gen2 (generation=1) IKS cluster, generation parameter needs to be set either in provider block or export via environment variable “IC_GENERATION”. By default the generation value is 2. 

FEATURES:

* New Resource: ([Terraform support for DNS service (beta service ) ibm_dns_zone, ibm_dns_permitted_network, ibm_dns_resource_record](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/958 ))
* New Resource: ([ibm_cis_firewall (lockdown)](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1201) )
* New Resource: ([ibm_lb_listener_policy](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1147) )

ENAHANCEMENTS:

* resource : Add support for resource group argument in ibm_is_network_acl ([#1265](https://github.com/IBM-Cloud/terraform-provider-ibm/issue/1265))
* resource : Support for IKS on Gen-2  (beta service) ([#1321](https://github.com/IBM-Cloud/terraform-provider-ibm/issue/1321))
* resource : Update functionality support for cis resources  ([#1180](https://github.com/IBM-Cloud/terraform-provider-ibm/issue/1180))
* resource : Add support for crn attribute for is_vpc   ([#1315](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1315) )
* data : Add support for crn attribute for is_vpc   ([#1317](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1317))

BUG FIXES:

* resource :  Fix the nil pointer exception for ibm_is_lb_listener resource ([#1289](https://github.com/IBM-Cloud/terraform-provider-ibm/issue/1289))


## 1.3.0 (April 02, 2020)

FEATURES:

* New Resource: ([ibm_iam_access_group_dynamic_rule](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/691))
* New Resource: ([ibm_api_gateway_endpoint](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1247))
* New Resource: ([ibm_api_gateway_endpoint_subscription](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1247))
* New DataSource: ([ibm_iam_access_group](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/953))
* New DataSource: ([ibm_api_gateway](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1247))

BUG FIXES:

* resource : Fix the destroy of cloudantnosqldb service([#1242](https://github.com/IBM-Cloud/terraform-provider-ibm/issue/1242)) 
* resource : Fix the ICD service endpoint for osl01([#1158](https://github.com/IBM-Cloud/terraform-provider-ibm/issue/1158)) 


## 1.2.6 (March 26, 2020)

ENHANCEMENTS:

* resource : Added support for cse_source_addresses  attribute for ibm_is_vpc  ([#1165](https://github.com/IBM-Cloud/terraform-provider-ibm/issue/1165)) 
* data : Added support for cse_source_addresses  attribute for ibm_is_vpc ([#1165](https://github.com/IBM-Cloud/terraform-provider-ibm/issue/1165)) 
* resource: Added support for new storage class smart for COS bucket  ([#1184](https://github.com/IBM-Cloud/terraform-provider-ibm/issue/1184))
* resource:  Allow deletion of non-existing resources like is_vpc, is_subnet, is_vpc_address_prefix and is_instance  ([#1229](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1229))
* resource:  Added support for force_delete argument for ibm_kp_key ([#1214](https://github.com/IBM-Cloud/terraform-provider-ibm/issue/1214))


## 1.2.5 (March 19, 2020)

ENHANCEMENTS:

* Provider : Adapt IAM access resources to v2 version ([#1183](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1183)) 
* resource: Added support for GUID attribute for ibm_cis and ibm_database ([#1169](https://github.com/IBM-Cloud/terraform-provider-ibm/issue/1169)) 
* data: Added support for GUID attribute for ibm_cis and ibm_database ([#1169](https://github.com/IBM-Cloud/terraform-provider-ibm/issue/1169))

BUG FIXES:

* resources : Updated the status string for `ibm_resource_instance, ibm_database and ibm_cis` to be inline with resource controller API changes ([#1190](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1190)) 
* resource/ibm_compute_bare_metal: Fix the order of provisioning of `bare metal` for processor capacity restriction type and SAP servers ([#1189](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1189)) 
* resource/ibm_resource_instance: Fix the order of provisioning of `block chain` platform service ([#1186](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1186)) 
* resource/ibm_container_cluster: Fix the force new for deprecated `billing` argument([#1187](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1187))  


## 1.2.4 (March 11, 2020)

ENHANCEMENTS:

* Provider: Added new parameter `zone` to support power virtual resources and data sources to work in multi-zone environment ([#1141](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1141))
* resource/ibm_pi_volume: Updated the list of volume types for power virtual volume ([#1149](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1149))
* resource/ibm_container_vpc_cluster : Added support for `ingress_hostname` and `ingress_secret` attributes ([#1167](https://github.com/IBM-Cloud/terraform-provider-ibm/issue/1167))
* data/ibm_container_vpc_cluster : Added support for `ingress_hostname` and `ingress_secret` attributes ([#1167](https://github.com/IBM-Cloud/terraform-provider-ibm/issue/1167))
* resource/ibm_is_floating_ip : Handle the case when floating IP is deleted manually ([#1160](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1160))

BUG FIXES:

* resources : Handle the case where the resource might be already deleted (manually) for ibm_iam_access_policies, ibm_iam_authorization_policies, ibm_iam_service_policies ([#1162](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1162))
* resource/ibm_is_inetwork_acl: Fix the order of creation of network acl ([#1123](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1123))
* resource/ibm_container_vpc_cluster: Added new attribute `wait_till` to control the cluster creation. Now user can control the cluster creation until master is ready / any one worker node is ready / ingress_hostname is  
  assigned.  ([#1143](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1143))
* resource/ibm_pi_instance: Fix the timeout configuration for create ([#1178](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1178))
* doc/ibm_cis_ip_addresses : Fix the description of data source ([#1178](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1178))

## 1.2.3 (March 03, 2020)

BUG FIXES:

* data/ibm_container_cluster_config : Fix the error to download the cluster config for VPC clusters ([#1150](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1150))

## 1.2.2 (February 26, 2020)

ENHANCEMENTS:

* resource/ibm_is_vpc: Improved error message for VPC creation ([#1106](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1106))
* resource/ibm_is_ssh_key: Improved error message for VPC SSH Key creation ([#1105](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1105))
* resource/ibm_container_cluster : Added gateway feature support for IKS clusters. This feature helps to create a cluster with a gateway worker pool of two gateway worker nodes that are connected to public and private VLANs to provide limited public access, and a compute worker pool of compute worker nodes that are connected to the private VLAN only. 
([#1125](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1125))
* data/ibm_conatiner_cluster_config : Extended the data source to provide additional attribute like admin_key, admin_certificate, ca_certificate, host and token. This attributes helps to connect to other providers like Kubernetes and Helm without loading cluster config file. ([#895](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/895))

BUG FIXES:

* doc/ibm_certificate_manager_order: Changed the type of rotate_key from string to bool ([#1110](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1110))
* resource/ibm_is_instance: Fix for updating security group for primary network interface for vpc instance. Now users can add or delete security groups([#1078](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1078))
* doc/ibm_resource_key : Provided an example in the docs as a workaround to create credentials using serviceID parameter ([#1121](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1121))
* resource/ibm_is_network_acl : Fix for crash during the update of rules. Fix for the order of rules creation. Now users can add or delete rules for network_acl ([#1117](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1117))
* resource/ibm_is_public_gateway : Added support for resource group and tags parameters ([#1102](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1102))
* resource/ibm_is_floating_ip : Added support for tags parameters ([#1131](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1131))
* resource/ibm_database : Parameters remote_leader_id, key_protect_instance and key_protect_key can’t be updated after creation. ([#1111](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1111))
* example/ibm-key-protect : Updated example to create an authorisation policy between COS and Key Protect instance([#1133](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1133))
* resource/ibm_resource_group: Removed suppression of error during deletion ([#1108](https://github.com/IBM-Cloud/terraform-provider-ibm/pull/1108))
* resource/ibm_iam_user_invite : Fix for inviting user from IBM Cloud lite account. ([#1114](https://github.com/IBM-Cloud/terraform-provider-ibm/issues/1114))
