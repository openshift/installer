import _ from 'lodash';

import { BARE_METAL } from './platforms';
import { keyToAlg } from './utils';

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
export const AWS_KMS = 'aws_kms';
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
export const SSH_AUTHORIZED_KEYS = 'sshAuthorizedKeys';
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
export const LICENSING = 'licensing';

export const toVPCSubnet = (region, subnets, deselected) => {
  const vpcSubnets = [];
  _.each(subnets, (v, availabilityZone) => {
    if (!availabilityZone.startsWith(region) || deselected && deselected[availabilityZone]) {
      return;
    }
    if (!v) {
      return;
    }
    const key = v.startsWith('subnet-') ? 'id' : 'instanceCIDR';
    vpcSubnets.push({availabilityZone, [key]: v});
  });
  return vpcSubnets;
};

const getZoneDomain = (cc) => {
  if (cc[PLATFORM_TYPE] === BARE_METAL) {
    throw new Error("Can't get base domain for bare metal!");
  }
  // TODO: if we ever change toExtraData()'s key, this breaks
  return _.get(cc, ['extra', AWS_HOSTED_ZONE_ID, 'zoneToName', cc[AWS_HOSTED_ZONE_ID]]);
};

export const getControllerDomain = (cc) => {
  if (cc[PLATFORM_TYPE] === BARE_METAL) {
    return cc[CONTROLLER_DOMAIN];
  }
  return `${cc[CLUSTER_SUBDOMAIN]}-k8s.${getZoneDomain(cc)}`;
};

export const getTectonicDomain = (cc) => {
  if (cc[PLATFORM_TYPE] === BARE_METAL) {
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
  [AWS_KMS]: '',
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
  [SSH_AUTHORIZED_KEYS]: [
    {
      id: 'initial-key',
      value: '',
    },
  ],
  [STS_ENABLED]: false,
  [TECTONIC_LICENSE]: '',
  [UPDATER]: {
    server: 'https://tectonic.update.core-os.net',
    channel: 'tectonic-1.5',
    appID: '6bc7b986-4654-4a0f-94b3-84ce6feb1db4',
  },
  [UPDATER_ENABLED]: false,
  [POD_CIDR]: "10.2.0.0/16",
  [SERVICE_CIDR]: "10.3.0.0/24",
};

export const toBaremetal = cc => {
  // TODO: (ggreer) send CLUSTER_NAME to backend in bare metal
  const ret = {
    clusterKind: 'tectonic-metal',
    dryRun: cc[DRY_RUN],
    cluster: {
      matchboxCA: cc[BM_MATCHBOX_CA],
      matchboxClientCert: cc[BM_MATCHBOX_CLIENT_CERT],
      matchboxClientKey: cc[BM_MATCHBOX_CLIENT_KEY],
      matchboxRPC: cc[BM_MATCHBOX_RPC],
      matchboxHTTP: cc[BM_MATCHBOX_HTTP],
      channel: cc[CHANNEL_TO_USE] || 'UNKNOWN',
      externalETCDClient: cc[EXTERNAL_ETCD_ENABLED] ? cc[EXTERNAL_ETCD_CLIENT] : '',
      version: cc[BM_OS_TO_USE],
      controllerDomain: getControllerDomain(cc),
      tectonicDomain: getTectonicDomain(cc),
      controllers: cc[BM_MASTERS].map(({mac, name}) => {
        return {mac, name};
      }),
      workers: cc[BM_WORKERS].map(({mac, name}) => {
        return {mac, name};
      }),
      sshAuthorizedKeys: cc[SSH_AUTHORIZED_KEYS].map(k => k.key).filter(k => k && k.length),
      podCIDR: cc[POD_CIDR],
      serviceCIDR: cc[SERVICE_CIDR],
      tectonic: {
        license: cc[TECTONIC_LICENSE],
        dockercfg: cc[PULL_SECRET],
        ingressKind: 'HostPort',
        identityAdminUser: cc[ADMIN_EMAIL],
        identityAdminPassword: window.btoa(cc[ADMIN_PASSWORD]),
        updater: {
          enabled: cc[UPDATER_ENABLED],
          server: cc[UPDATER].server,
          channel: cc[UPDATER].channel,
          appID: cc[UPDATER].appID,
        },
      },
    },
  };

  if (cc[CA_TYPE] === 'owned') {
    ret.cluster.caCertificate = cc[CA_CERTIFICATE];
    ret.cluster.caPrivateKey = cc[CA_PRIVATE_KEY];
  }

  return ret;
};

export const toAWS = (cc, FORMS) => {
  const controllerDomain = getControllerDomain(cc);
  const tectonicDomain = getTectonicDomain(cc);
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
    controllerSubnets = toVPCSubnet(region, cc[AWS_CONTROLLER_SUBNET_IDS], cc[DESELECTED_FIELDS][AWS_SUBNETS]);
    workerSubnets = toVPCSubnet(region, cc[AWS_WORKER_SUBNET_IDS], cc[DESELECTED_FIELDS][AWS_SUBNETS]);
  }

  const ret = {
    clusterKind: 'tectonic-aws',
    dryRun: cc[DRY_RUN],
    cluster: {
      accessKeyID: cc[AWS_ACCESS_KEY_ID],
      secretAccessKey: cc[AWS_SECRET_ACCESS_KEY],
      cloudForm: {
        channel: cc[CHANNEL_TO_USE],
        clusterName: cc[CLUSTER_NAME],
        elbScheme: cc[AWS_CREATE_VPC] === 'VPC_PRIVATE' ? "internal" : "internet-facing",
        controllerDomain: controllerDomain,
        externalETCDClient: cc[EXTERNAL_ETCD_ENABLED] ? cc[EXTERNAL_ETCD_CLIENT] : '',
        tectonicDomain: tectonicDomain,
        region: cc[AWS_REGION],
        hostedZoneID: cc[AWS_HOSTED_ZONE_ID],
        kmsKeyARN: cc[AWS_KMS],
        keyName: cc[AWS_SSH],
        etcdCount: etcds[NUMBER_OF_INSTANCES],
        etcdInstanceType: etcds[INSTANCE_TYPE],
        etcdRootVolumeSize: etcds[STORAGE_SIZE_IN_GIB],
        etcdRootVolumeType: etcds[STORAGE_TYPE],
        etcdRootVolumeIOPS: etcds[STORAGE_TYPE] === 'io1' ? etcds[STORAGE_IOPS] : undefined,
        controllerCount: controllers[NUMBER_OF_INSTANCES],
        controllerInstanceType: controllers[INSTANCE_TYPE],
        controllerRootVolumeSize: controllers[STORAGE_SIZE_IN_GIB],
        controllerRootVolumeType: controllers[STORAGE_TYPE],
        controllerRootVolumeIOPS: controllers[STORAGE_TYPE] === 'io1' ? controllers[STORAGE_IOPS] : undefined,
        workerCount: workers[NUMBER_OF_INSTANCES],
        workerInstanceType: workers[INSTANCE_TYPE],
        workerRootVolumeSize: workers[STORAGE_SIZE_IN_GIB],
        workerRootVolumeType: workers[STORAGE_TYPE],
        workerRootVolumeIOPS: workers[STORAGE_TYPE] === 'io1' ? workers[STORAGE_IOPS] : undefined,
        podCIDR: cc[POD_CIDR],
        serviceCIDR: cc[SERVICE_CIDR],
        tags: _.filter(cc[AWS_TAGS], ({key, value}) => key && value), // don't send empty tags
        controllerSubnets,
        workerSubnets,
      },
      tectonic: {
        license: cc[TECTONIC_LICENSE],
        dockercfg: cc[PULL_SECRET],
        ingressKind: 'NodePort',
        identityAdminUser: cc[ADMIN_EMAIL],
        identityAdminPassword: window.btoa(cc[ADMIN_PASSWORD]),
        updater: {
          enabled: cc[UPDATER_ENABLED],
          server: cc[UPDATER].server,
          channel: cc[UPDATER].channel,
          appID: cc[UPDATER].appID,
        },
      },
    },
  };
  const { cluster } = ret;

  if (cc[STS_ENABLED]) {
    cluster.sessionToken = cc[AWS_SESSION_TOKEN];
  }
  if (cc[AWS_CREATE_VPC] === 'VPC_CREATE') {
    cluster.cloudForm.vpcCIDR = cc[AWS_VPC_CIDR];
  } else {
    cluster.cloudForm.vpcID = cc[AWS_VPC_ID];
  }

  if (cc[CA_TYPE] === 'owned') {
    cluster.caCertificate = cc[CA_CERTIFICATE];
    cluster.caPrivateKey = cc[CA_PRIVATE_KEY];
  }
  return ret;
};

const toSubnetObj = (subnets, key) => _(subnets)
  .keyBy(o => o.availabilityZone)
  .mapValues(v => v[key])
  .value();


export const toAWS_TF = (cc, FORMS) => {
  const controllers = FORMS[AWS_CONTROLLERS].getData(cc);
  const etcds = FORMS[AWS_ETCDS].getData(cc);
  const workers = FORMS[AWS_WORKERS].getData(cc);

  const region = cc[AWS_REGION];
  let controllerSubnets;
  let workerSubnets;

  if (cc[AWS_CREATE_VPC] === 'VPC_CREATE') {
    controllerSubnets = toSubnetObj(
      toVPCSubnet(region, cc[AWS_CONTROLLER_SUBNETS], cc[DESELECTED_FIELDS][AWS_SUBNETS]),
      'instanceCIDR'
    );
    workerSubnets = toSubnetObj(
      toVPCSubnet(region, cc[AWS_WORKER_SUBNETS], cc[DESELECTED_FIELDS][AWS_SUBNETS]),
      'instanceCIDR'
    );
  } else {
    controllerSubnets = toSubnetObj(
      toVPCSubnet(region, cc[AWS_CONTROLLER_SUBNET_IDS], cc[DESELECTED_FIELDS][AWS_SUBNETS]),
      'id'
    );
    workerSubnets = toSubnetObj(
      toVPCSubnet(region, cc[AWS_WORKER_SUBNET_IDS], cc[DESELECTED_FIELDS][AWS_SUBNETS]),
      'id'
    );
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
    adminPassword: window.btoa(cc[ADMIN_PASSWORD]),
    credentials: {
      AWSAccessKeyID: cc[AWS_ACCESS_KEY_ID],
      AWSSecretAccessKey: cc[AWS_SECRET_ACCESS_KEY],
      AWSRegion: cc[AWS_REGION],
    },
    variables: {
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
  } else {
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
