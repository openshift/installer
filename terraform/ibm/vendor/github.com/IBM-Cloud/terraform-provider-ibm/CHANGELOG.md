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


