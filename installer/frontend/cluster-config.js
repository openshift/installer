import _ from 'lodash';

import { BARE_METAL_TF } from './platforms';
import { keyToAlg } from './utils';

// TODO: (ggreer) clean up key names. Warning: Doing this will break progress files.
export const AWS_ACCESS_KEY_ID = 'awsAccessKeyId';
export const AWS_SUBNETS = 'awsSubnets';
export const AWS_CONTROLLER_SUBNETS = 'awsControllerSubnets';
export const AWS_CONTROLLER_SUBNET_IDS = 'awsControllerSubnetIds';
export const DESELECTED_FIELDS = 'deselectedFields';
export const AWS_DOMAIN = 'awsDomain';
export const AWS_HOSTED_ZONE_ID = 'awsHostedZoneId';
export const AWS_SPLIT_DNS = 'awsSplitDNS';
export const AWS_REGION = 'awsRegion';
export const AWS_SECRET_ACCESS_KEY = 'awsSecretAccessKey';
export const AWS_SESSION_TOKEN = 'awsSessionToken';
export const AWS_SSH = 'aws_ssh';
export const AWS_TAGS = 'awsTags';

export const AWS_CREATE_VPC = 'awsCreateVpc';
export const AWS_VPC_CIDR = 'awsVpcCIDR';
export const AWS_VPC_ID = 'awsVpcId';

export const AWS_WORKER_SUBNETS = 'awsWorkerSubnets';
export const AWS_WORKER_SUBNET_IDS = 'awsWorkerSubnetIds';

export const BM_MATCHBOX_CA = 'matchboxCA';
export const BM_MATCHBOX_CLIENT_CERT = 'matchboxClientCert';
export const BM_MATCHBOX_CLIENT_KEY = 'matchboxClientKey';
export const BM_MATCHBOX_HTTP = 'matchboxHTTP';
export const BM_MATCHBOX_RPC = 'matchboxRPC';
export const BM_MASTERS = 'masters';
export const BM_OS_TO_USE = 'osToUse';
export const BM_TECTONIC_DOMAIN = 'tectonicDomain';
export const BM_WORKERS = 'workers';

export const CA_CERTIFICATE = 'caCertificate';
export const CA_PRIVATE_KEY = 'caPrivateKey';
export const CA_TYPE = 'caType';
export const CA_TYPES = {SELF_SIGNED: 'self-signed', OWNED: 'owned'};
export const CLUSTER_NAME = 'clusterName';
export const CLUSTER_SUBDOMAIN = 'clusterSubdomain';
export const CONTROLLER_DOMAIN = 'controllerDomain';
export const EXTERNAL_ETCD_CLIENT = 'externalETCDClient';

export const ETCD_OPTION = 'etcdOption';

export const DRY_RUN = 'dryRun';
export const PLATFORM_TYPE = 'platformType';
export const PULL_SECRET = 'pullSecret';
export const SSH_AUTHORIZED_KEY = 'sshAuthorizedKey';
export const STS_ENABLED = 'sts_enabled';
export const TECTONIC_LICENSE = 'tectonicLicense';
export const ADMIN_EMAIL = 'adminEmail';
export const ADMIN_PASSWORD = 'adminPassword';
export const ADMIN_PASSWORD2 = 'adminPassword2';

// Networking
export const POD_CIDR = 'podCIDR';
export const SERVICE_CIDR = 'serviceCIDR';

export const IAM_ROLE = 'iamRole';
export const NUMBER_OF_INSTANCES = 'numberOfInstances';
export const INSTANCE_TYPE = 'instanceType';
export const STORAGE_SIZE_IN_GIB = 'storageSizeInGiB';
export const STORAGE_TYPE = 'storageType';
export const STORAGE_IOPS = 'storageIOPS';

export const RETRY = 'retry';

// FORMS:
export const AWS_CREDS = 'AWSCreds';
export const AWS_ETCDS = 'aws_etcds';
export const AWS_VPC_FORM = 'aws_vpc';
export const AWS_CONTROLLERS = 'aws_controllers';
export const AWS_CLUSTER_INFO = 'aws_clusterInfo';
export const AWS_WORKERS = 'aws_workers';
export const AWS_REGION_FORM = 'aws_regionForm';
export const BM_SSH_KEY = 'bm_sshKey';
export const CREDS = 'creds';
export const LICENSING = 'licensing';
export const PLATFORM_FORM = 'platform';

export const SPLIT_DNS_ON = 'on';
export const SPLIT_DNS_OFF = 'off';
export const SPLIT_DNS_OPTIONS = {
  [SPLIT_DNS_ON]: 'Create an additional Route 53 private zone (default).',
  [SPLIT_DNS_OFF]: 'Do not create a private zone.',
};

const EXTERNAL = 'external';
const PROVISIONED = 'provisioned';
export const ETCD_OPTIONS = { EXTERNAL, PROVISIONED };

// String that would be an invalid IAM role name
export const IAM_ROLE_CREATE_OPTION = '%create%';

export const toVPCSubnet = (region, subnets, deselected) => {
  const vpcSubnets = {};
  _.each(subnets, (v, availabilityZone) => {
    if (!availabilityZone.startsWith(region) || deselected && deselected[availabilityZone]) {
      return;
    }
    if (!v) {
      return;
    }
    vpcSubnets[availabilityZone] = v;
  });
  return vpcSubnets;
};

export const toVPCSubnetID = (region, subnets, deselected) => {
  const vpcSubnets = [];
  _.each(subnets, (v, availabilityZone) => {
    if (!availabilityZone.startsWith(region) || deselected && deselected[availabilityZone]) {
      return;
    }
    if (!v) {
      return;
    }
    vpcSubnets.push(v);
  });
  return vpcSubnets;
};

const getZoneDomain = (cc) => {
  if (cc[PLATFORM_TYPE] === BARE_METAL_TF) {
    throw new Error("Can't get base domain for bare metal!");
  }
  return _.get(cc, ['extra', AWS_HOSTED_ZONE_ID, 'zoneToName', cc[AWS_HOSTED_ZONE_ID]]);
};

export const getControllerDomain = (cc) => {
  if (cc[PLATFORM_TYPE] === BARE_METAL_TF) {
    return cc[CONTROLLER_DOMAIN];
  }
  return `${cc[CLUSTER_SUBDOMAIN]}-k8s.${getZoneDomain(cc)}`;
};

export const getTectonicDomain = (cc) => {
  if (cc[PLATFORM_TYPE] === BARE_METAL_TF) {
    return cc[BM_TECTONIC_DOMAIN];
  }
  return cc[CLUSTER_SUBDOMAIN] + (cc[CLUSTER_SUBDOMAIN].endsWith('.') ? '' : '.') + getZoneDomain(cc);
};

export const DEFAULT_CLUSTER_CONFIG = {
  error: {}, // to store validation errors
  error_async: {}, // to store async validation errors
  ignore: {}, // to store validation errors
  inFly: {}, // to store inFly
  extra: {}, // extraneous, non-value data for this field
  bootCfgInfly: false, // TODO (ggreer): total hack. clean up after release
  [AWS_VPC_ID]: '',
  [AWS_CONTROLLER_SUBNET_IDS]: {},
  [DESELECTED_FIELDS]: {},
  [AWS_WORKER_SUBNET_IDS]: {},
  [BM_MATCHBOX_HTTP]: '',
  [BM_OS_TO_USE]: '',
  [DRY_RUN]: false,
  [RETRY]: false, // whether we're retrying a terraform apply
};

export const toAWS_TF = (cc, FORMS) => {
  const controllers = FORMS[AWS_CONTROLLERS].getData(cc);
  const etcds = FORMS[AWS_ETCDS].getData(cc);
  const workers = FORMS[AWS_WORKERS].getData(cc);

  const region = cc[AWS_REGION];
  let controllerSubnets;
  let workerSubnets;

  if (cc[AWS_CREATE_VPC] === 'VPC_CREATE') {
    controllerSubnets = toVPCSubnet(region, cc[AWS_CONTROLLER_SUBNETS], cc[DESELECTED_FIELDS][AWS_SUBNETS]);
    workerSubnets = toVPCSubnet(region, cc[AWS_WORKER_SUBNETS], cc[DESELECTED_FIELDS][AWS_SUBNETS]);
  } else {
    controllerSubnets = toVPCSubnetID(region, cc[AWS_CONTROLLER_SUBNET_IDS], cc[DESELECTED_FIELDS][AWS_SUBNETS]);
    workerSubnets = toVPCSubnetID(region, cc[AWS_WORKER_SUBNET_IDS], cc[DESELECTED_FIELDS][AWS_SUBNETS]);
  }

  const extraTags = {};
  _.each(cc[AWS_TAGS], ({key, value}) => {
    if (key && value) {
      extraTags[key] = value;
    }
  });

  const ret = {
    clusterKind: 'tectonic-aws-tf',
    dryRun: cc[DRY_RUN],
    platform: 'aws',
    license: cc[TECTONIC_LICENSE],
    pullSecret: cc[PULL_SECRET],
    retry: cc[RETRY],
    credentials: {
      AWSAccessKeyID: cc[AWS_ACCESS_KEY_ID],
      AWSSecretAccessKey: cc[AWS_SECRET_ACCESS_KEY],
    },
    variables: {
      // eslint-disable-next-line no-sync
      tectonic_admin_password: cc[ADMIN_PASSWORD],
      tectonic_aws_region: cc[AWS_REGION],
      tectonic_admin_email: cc[ADMIN_EMAIL],
      tectonic_aws_master_ec2_type: controllers[INSTANCE_TYPE],
      tectonic_aws_master_iam_role_name: controllers[IAM_ROLE] === IAM_ROLE_CREATE_OPTION ? undefined : controllers[IAM_ROLE],
      tectonic_aws_master_root_volume_iops: controllers[STORAGE_TYPE] === 'io1' ? controllers[STORAGE_IOPS] : undefined,
      tectonic_aws_master_root_volume_size: controllers[STORAGE_SIZE_IN_GIB],
      tectonic_aws_master_root_volume_type: controllers[STORAGE_TYPE],
      tectonic_aws_worker_ec2_type: workers[INSTANCE_TYPE],
      tectonic_aws_worker_iam_role_name: workers[IAM_ROLE] === IAM_ROLE_CREATE_OPTION ? undefined : workers[IAM_ROLE],
      tectonic_aws_worker_root_volume_iops: workers[STORAGE_TYPE] === 'io1' ? controllers[STORAGE_IOPS] : undefined,
      tectonic_aws_worker_root_volume_size: workers[STORAGE_SIZE_IN_GIB],
      tectonic_aws_worker_root_volume_type: workers[STORAGE_TYPE],
      tectonic_aws_ssh_key: cc[AWS_SSH],
      tectonic_base_domain: getZoneDomain(cc),
      tectonic_cluster_cidr: cc[POD_CIDR],
      tectonic_cluster_name: cc[CLUSTER_NAME],
      tectonic_master_count: controllers[NUMBER_OF_INSTANCES],
      tectonic_service_cidr: cc[SERVICE_CIDR],
      tectonic_worker_count: workers[NUMBER_OF_INSTANCES],
      // TODO: shouldn't hostedZoneID be specified somewhere?
      tectonic_dns_name: cc[CLUSTER_SUBDOMAIN],
    },
  };

  if (cc[ETCD_OPTION] === EXTERNAL) {
    ret.variables.tectonic_etcd_servers = [cc[EXTERNAL_ETCD_CLIENT]];
  } else if (cc[ETCD_OPTION] === PROVISIONED) {
    ret.variables.tectonic_aws_etcd_ec2_type = etcds[INSTANCE_TYPE];
    ret.variables.tectonic_aws_etcd_root_volume_iops = etcds[STORAGE_TYPE] === 'io1' ? etcds[STORAGE_IOPS] : undefined;
    ret.variables.tectonic_aws_etcd_root_volume_size = etcds[STORAGE_SIZE_IN_GIB];
    ret.variables.tectonic_aws_etcd_root_volume_type = etcds[STORAGE_TYPE];
    ret.variables.tectonic_etcd_count = etcds[NUMBER_OF_INSTANCES];
  }

  if (_.size(extraTags) > 0) {
    ret.variables.tectonic_aws_extra_tags = extraTags;
  }

  if (cc[STS_ENABLED]) {
    ret.credentials.AWSSessionToken = cc[AWS_SESSION_TOKEN];
  }
  if (cc[AWS_CREATE_VPC] === 'VPC_CREATE') {
    ret.variables.tectonic_aws_vpc_cidr_block = cc[AWS_VPC_CIDR];
    ret.variables.tectonic_aws_master_custom_subnets = controllerSubnets;
    ret.variables.tectonic_aws_worker_custom_subnets = workerSubnets;
  } else {
    ret.variables.tectonic_aws_external_vpc_id = cc[AWS_VPC_ID];
    ret.variables.tectonic_aws_external_master_subnet_ids = controllerSubnets;
    ret.variables.tectonic_aws_external_worker_subnet_ids = workerSubnets;
    ret.variables.tectonic_aws_public_endpoints = cc[AWS_CREATE_VPC] !== 'VPC_PRIVATE';
  }

  if (cc[AWS_CREATE_VPC] !== 'VPC_PRIVATE' && cc[AWS_SPLIT_DNS] === SPLIT_DNS_OFF) {
    ret.variables.tectonic_aws_private_endpoints = false;
  }

  if (cc[CA_TYPE] === CA_TYPES.OWNED) {
    ret.variables.tectonic_ca_cert = cc[CA_CERTIFICATE];
    ret.variables.tectonic_ca_key = cc[CA_PRIVATE_KEY];
    ret.variables.tectonic_ca_key_alg = keyToAlg(cc[CA_PRIVATE_KEY]);
  }
  return ret;
};

export const toBaremetal_TF = (cc, FORMS) => {
  const sshKey = FORMS[BM_SSH_KEY].getData(cc);
  const masters = cc[BM_MASTERS];
  const workers = cc[BM_WORKERS];

  const ret = {
    clusterKind: 'tectonic-metal',
    dryRun: cc[DRY_RUN],
    platform: 'metal',
    license: cc[TECTONIC_LICENSE],
    pullSecret: cc[PULL_SECRET],
    retry: cc[RETRY],
    variables: {
      // eslint-disable-next-line no-sync
      tectonic_admin_password: cc[ADMIN_PASSWORD],
      tectonic_cluster_name: cc[CLUSTER_NAME],
      tectonic_admin_email: cc[ADMIN_EMAIL],
      tectonic_container_linux_version: cc[BM_OS_TO_USE],
      tectonic_metal_ingress_domain: getTectonicDomain(cc),
      tectonic_metal_controller_domain: getControllerDomain(cc),
      tectonic_metal_controller_domains: masters.map(({host}) => host),
      tectonic_metal_controller_names: masters.map(({host}) => host.split('.')[0]),
      tectonic_metal_controller_macs: masters.map(({mac}) => mac),
      tectonic_metal_worker_domains: workers.map(({host}) => host),
      tectonic_metal_worker_names: workers.map(({host}) => host.split('.')[0]),
      tectonic_metal_worker_macs: workers.map(({mac}) => mac),
      tectonic_metal_matchbox_http_url: `http://${cc[BM_MATCHBOX_HTTP]}`,
      tectonic_metal_matchbox_rpc_endpoint: cc[BM_MATCHBOX_RPC],
      tectonic_metal_matchbox_ca: cc[BM_MATCHBOX_CA],
      tectonic_metal_matchbox_client_cert: cc[BM_MATCHBOX_CLIENT_CERT],
      tectonic_metal_matchbox_client_key: cc[BM_MATCHBOX_CLIENT_KEY],
      tectonic_ssh_authorized_key: sshKey[SSH_AUTHORIZED_KEY],
      tectonic_cluster_cidr: cc[POD_CIDR],
      tectonic_service_cidr: cc[SERVICE_CIDR],
      tectonic_dns_name: cc[CLUSTER_SUBDOMAIN],
      tectonic_base_domain: 'unused',
    },
  };

  if (cc[ETCD_OPTION] === EXTERNAL) {
    ret.variables.tectonic_etcd_servers = [cc[EXTERNAL_ETCD_CLIENT]];
  }

  if (cc[CA_TYPE] === CA_TYPES.OWNED) {
    ret.variables.tectonic_ca_cert = cc[CA_CERTIFICATE];
    ret.variables.tectonic_ca_key = cc[CA_PRIVATE_KEY];
    ret.variables.tectonic_ca_key_alg = keyToAlg(cc[CA_PRIVATE_KEY]);
  }

  return ret;
};
