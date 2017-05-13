import _ from 'lodash';
import bcrypt from 'bcryptjs';

import { BARE_METAL_TF } from './platforms';
import { keyToAlg } from './utils';

const bcryptCost = 12;

let defaultPlatformType = '';
try {
  defaultPlatformType = window.config.platforms[0];
} catch (unused) {
  // So tests pass
}

// TODO: (ggreer) clean up key names. Warning: Doing this will break progress files.
export const AWS_ACCESS_KEY_ID = 'awsAccessKeyId';
export const AWS_SUBNETS = 'awsSubnets';
export const AWS_CONTROLLER_SUBNETS = 'awsControllerSubnets';
export const AWS_CONTROLLER_SUBNET_IDS = 'awsControllerSubnetIds';
export const DESELECTED_FIELDS = 'deselectedFields';
export const AWS_DOMAIN = 'awsDomain';
export const AWS_HOSTED_ZONE_ID = 'awsHostedZoneId';
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
export const BM_MASTERS_COUNT = 'mastersCount';
export const BM_OS_TO_USE = 'osToUse';
export const BM_TECTONIC_DOMAIN = 'tectonicDomain';
export const BM_WORKERS = 'workers';
export const BM_WORKERS_COUNT = 'workersCount';

export const CA_CERTIFICATE = 'caCertificate';
export const CA_PRIVATE_KEY = 'caPrivateKey';
export const CA_TYPE = 'caType';
export const CHANNEL_TO_USE = 'channelToUse';
export const CLUSTER_NAME = 'clusterName';
export const CLUSTER_SUBDOMAIN = 'clusterSubdomain';
export const CONTROLLER_DOMAIN = 'controllerDomain';
export const EXTERNAL_ETCD_CLIENT = 'externalETCDClient';
export const EXTERNAL_ETCD_ENABLED = 'externalETCDEnabled';
export const DRY_RUN = 'dryRun';
export const ENTITLEMENTS = 'entitlements';
export const PLATFORM_TYPE = 'platformType';
export const PULL_SECRET = 'pullSecret';
export const SSH_AUTHORIZED_KEY = 'sshAuthorizedKey';
export const STS_ENABLED = 'sts_enabled';
export const TECTONIC_LICENSE = 'tectonicLicense';
export const UPDATER = 'updater';
export const UPDATER_ENABLED = 'updater_enabled';
export const ADMIN_EMAIL = 'adminEmail';
export const ADMIN_PASSWORD = 'adminPassword';

// Networking
export const POD_CIDR = 'podCIDR';
export const SERVICE_CIDR = 'serviceCIDR';

export const NUMBER_OF_INSTANCES = 'numberOfInstances';
export const INSTANCE_TYPE = 'instanceType';
export const STORAGE_SIZE_IN_GIB = 'storageSizeInGiB';
export const STORAGE_TYPE = 'storageType';
export const STORAGE_IOPS = 'storageIOPS';

export const RETRY = 'retry';

// FORMS:
export const AWS_ETCDS = 'aws_etcds';
export const AWS_VPC_FORM = 'aws_vpc';
export const AWS_CONTROLLERS = 'aws_controllers';
export const AWS_CLUSTER_INFO = 'aws_clusterInfo';
export const AWS_WORKERS = 'aws_workers';
export const BM_SSH_KEY = 'bm_sshKey';
export const LICENSING = 'licensing';

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
  // TODO: if we ever change toExtraData()'s key, this breaks
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
  const tectonicDomain = cc[CLUSTER_SUBDOMAIN] + (cc[CLUSTER_SUBDOMAIN].endsWith('.') ? '' : '.') + getZoneDomain(cc);
  return tectonicDomain;
};

export const DEFAULT_CLUSTER_CONFIG = {
  error: {}, // to store validation errors
  error_async: {}, // to store async validation errors
  ignore: {}, // to store validation errors
  inFly: {}, // to store inFly
  extra: {}, // extraneous, non-value data for this field
  bootCfgInfly: false, // TODO (ggreer): total hack. clean up after release
  [ADMIN_EMAIL]: '',
  [ADMIN_PASSWORD]: '',
  [AWS_ACCESS_KEY_ID]: '',
  [AWS_REGION]: '',
  [AWS_SECRET_ACCESS_KEY]: '',
  [AWS_SESSION_TOKEN]: '',
  [AWS_SSH]: '',
  [AWS_VPC_ID]: '',
  [AWS_VPC_CIDR]: '10.0.0.0/16',
  [AWS_CONTROLLER_SUBNETS]: {},
  [AWS_CONTROLLER_SUBNET_IDS]: {},
  [DESELECTED_FIELDS]: {},
  [AWS_WORKER_SUBNETS]: {},
  [AWS_WORKER_SUBNET_IDS]: {},
  [BM_MATCHBOX_CA]: '',
  [BM_MATCHBOX_CLIENT_CERT]: '',
  [BM_MATCHBOX_CLIENT_KEY]: '',
  [BM_MATCHBOX_HTTP]: '',
  [BM_MATCHBOX_RPC]: '',
  [BM_MASTERS]: [],
  [BM_MASTERS_COUNT]: 1,
  [BM_OS_TO_USE]: '',
  [BM_TECTONIC_DOMAIN]: '',
  [BM_WORKERS]: [],
  [BM_WORKERS_COUNT]: 1,
  [CA_CERTIFICATE]: '',
  [CA_PRIVATE_KEY]: '',
  [CA_TYPE]: 'self-signed',
  [CLUSTER_NAME]: '',
  [CONTROLLER_DOMAIN]: '',
  [DRY_RUN]: false,
  [ENTITLEMENTS]: {},
  [PLATFORM_TYPE]: defaultPlatformType,
  [PULL_SECRET]: '',
  [RETRY]: false, // whether we're retrying a terraform apply
  [STS_ENABLED]: false,
  [TECTONIC_LICENSE]: '',
  [UPDATER]: {
    server: 'https://tectonic.update.core-os.net',
    channel: 'tectonic-1.6',
    appID: '6bc7b986-4654-4a0f-94b3-84ce6feb1db4',
  },
  [UPDATER_ENABLED]: false,
  [POD_CIDR]: "10.2.0.0/16",
  [SERVICE_CIDR]: "10.3.0.0/16",
};


export const toAWS_TF = (cc, FORMS, opts={}) => {
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
    if(key && value) {
      extraTags[key] = value;
    }
  });

  const ret = {
    clusterKind: 'tectonic-aws-tf',
    dryRun: cc[DRY_RUN],
    platform: "aws",
    license: cc[TECTONIC_LICENSE],
    pullSecret: cc[PULL_SECRET],
    credentials: {
      AWSAccessKeyID: cc[AWS_ACCESS_KEY_ID],
      AWSSecretAccessKey: cc[AWS_SECRET_ACCESS_KEY],
    },
    variables: {
      // eslint-disable-next-line no-sync
      tectonic_admin_password_hash: bcrypt.hashSync(cc[ADMIN_PASSWORD], opts.salt || bcrypt.genSaltSync(bcryptCost)),
      tectonic_aws_region: cc[AWS_REGION],
      tectonic_admin_email: cc[ADMIN_EMAIL],
      tectonic_aws_master_ec2_type: controllers[INSTANCE_TYPE],
      tectonic_aws_master_root_volume_iops: controllers[STORAGE_TYPE] === 'io1' ? controllers[STORAGE_IOPS] : undefined,
      tectonic_aws_master_root_volume_size: controllers[STORAGE_SIZE_IN_GIB],
      tectonic_aws_master_root_volume_type: controllers[STORAGE_TYPE],
      tectonic_aws_worker_ec2_type: workers[INSTANCE_TYPE],
      tectonic_aws_worker_root_volume_iops: workers[STORAGE_TYPE] === 'io1' ? controllers[STORAGE_IOPS] : undefined,
      tectonic_aws_worker_root_volume_size: workers[STORAGE_SIZE_IN_GIB],
      tectonic_aws_worker_root_volume_type: workers[STORAGE_TYPE],
      tectonic_aws_ssh_key: cc[AWS_SSH],
      tectonic_base_domain: getZoneDomain(cc),
      tectonic_cl_channel: cc[CHANNEL_TO_USE],
      tectonic_cluster_cidr: cc[POD_CIDR],
      tectonic_cluster_name: cc[CLUSTER_NAME],
      tectonic_master_count: controllers[NUMBER_OF_INSTANCES],
      tectonic_service_cidr: cc[SERVICE_CIDR],
      tectonic_worker_count: workers[NUMBER_OF_INSTANCES],
      // TODO: shouldn't hostedZoneID be specified somewhere?
      tectonic_dns_name: cc[CLUSTER_SUBDOMAIN],
      tectonic_experimental: cc[UPDATER_ENABLED],
    },
  };

  if (cc[EXTERNAL_ETCD_ENABLED]) {
    ret.variables.tectonic_etcd_servers = [cc[EXTERNAL_ETCD_CLIENT]];
  } else if (!cc[UPDATER_ENABLED]) {
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
    ret.variables.tectonic_aws_external_vpc_public = cc[AWS_CREATE_VPC] !== 'VPC_PRIVATE';
  }
  if (cc[CA_TYPE] === 'owned') {
    ret.variables.tectonic_ca_cert = cc[CA_CERTIFICATE];
    ret.variables.tectonic_ca_key = cc[CA_PRIVATE_KEY];
    ret.variables.tectonic_ca_key_alg = keyToAlg(cc[CA_PRIVATE_KEY]);
  }

  return ret;
};

export const toBaremetal_TF = (cc, FORMS, opts={}) => {
  const sshKey = FORMS[BM_SSH_KEY].getData(cc);

  const ret = {
    clusterKind: 'tectonic-metal',
    dryRun: cc[DRY_RUN],
    platform: 'metal',
    license: cc[TECTONIC_LICENSE],
    pullSecret: cc[PULL_SECRET],
    variables: {
      // eslint-disable-next-line no-sync
      tectonic_admin_password_hash: bcrypt.hashSync(cc[ADMIN_PASSWORD], opts.salt || bcrypt.genSaltSync(bcryptCost)),
      tectonic_cluster_name: cc[CLUSTER_NAME],
      tectonic_cl_channel: cc[CHANNEL_TO_USE],
      tectonic_admin_email: cc[ADMIN_EMAIL],
      tectonic_metal_cl_version: cc[BM_OS_TO_USE],
      tectonic_metal_ingress_domain: getTectonicDomain(cc),
      tectonic_metal_controller_domain: getControllerDomain(cc),
      tectonic_metal_controller_domains: cc[BM_MASTERS].map(({name}) => name),
      tectonic_metal_controller_names: cc[BM_MASTERS].map(({name}) => name.split('.')[0]),
      tectonic_metal_controller_macs: cc[BM_MASTERS].map(({mac}) => mac),
      tectonic_metal_worker_domains: cc[BM_WORKERS].map(({name}) => name),
      tectonic_metal_worker_names: cc[BM_WORKERS].map(({name}) => name.split('.')[0]),
      tectonic_metal_worker_macs: cc[BM_WORKERS].map(({mac}) => mac),
      tectonic_metal_matchbox_http_url: `http://${cc[BM_MATCHBOX_HTTP]}`,
      tectonic_metal_matchbox_rpc_endpoint: cc[BM_MATCHBOX_RPC],
      tectonic_metal_matchbox_ca: cc[BM_MATCHBOX_CA],
      tectonic_metal_matchbox_client_cert: cc[BM_MATCHBOX_CLIENT_CERT],
      tectonic_metal_matchbox_client_key: cc[BM_MATCHBOX_CLIENT_KEY],
      tectonic_ssh_authorized_key: sshKey[SSH_AUTHORIZED_KEY],
      tectonic_cluster_cidr: cc[POD_CIDR],
      tectonic_service_cidr: cc[SERVICE_CIDR],
      tectonic_dns_name: cc[CLUSTER_SUBDOMAIN],
      tectonic_experimental: cc[UPDATER_ENABLED],
      tectonic_base_domain: 'unused',
    },
  };

  if (cc[EXTERNAL_ETCD_ENABLED]) {
    ret.variables.tectonic_etcd_servers = [cc[EXTERNAL_ETCD_CLIENT]];
  }

  if (cc[CA_TYPE] === 'owned') {
    ret.variables.tectonic_ca_cert = cc[CA_CERTIFICATE];
    ret.variables.tectonic_ca_key = cc[CA_PRIVATE_KEY];
    ret.variables.tectonic_ca_key_alg = keyToAlg(cc[CA_PRIVATE_KEY]);
  }

  return ret;
};
