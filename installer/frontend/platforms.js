import _ from 'lodash';

export const AWS = 'aws';
export const AWS_TF = 'aws-tf';
export const AZURE = 'azure';
export const BARE_METAL = 'bare-metal';
export const BARE_METAL_TF = 'bare-metal-tf';
export const OPENSTACK = 'openstack';

export const isSupported = platform => [AWS_TF, BARE_METAL_TF].includes(platform);

export const isEnabled = platform => _.get(window.config, 'platforms', []).includes(platform);

export const PLATFORM_NAMES = {
  [AWS]: 'Amazon Web Services',
  [AWS_TF]: 'Amazon Web Services',
  [AZURE]: 'Microsoft Azure',
  [BARE_METAL]: 'Bare Metal',
  [BARE_METAL_TF]: 'Bare Metal',
  [OPENSTACK]: 'OpenStack',
};

export const optGroups = [
  ['Graphical Installer (default)', AWS_TF, BARE_METAL_TF],
  ['Advanced Installer', AWS, BARE_METAL, AZURE, OPENSTACK],
];

export const DOCS = {
  [AWS]: 'https://coreos.com/tectonic/docs/latest/install/aws/aws-terraform.html',
  [AWS_TF]: 'https://coreos.com/tectonic/docs/latest/install/aws/index.html',
  [AZURE]: 'https://coreos.com/tectonic/docs/latest/install/azure/azure-terraform.html',
  [BARE_METAL]: 'https://coreos.com/tectonic/docs/latest/install/bare-metal/metal-terraform.html',
  [BARE_METAL_TF]: 'https://coreos.com/tectonic/docs/latest/install/bare-metal/index.html',
  [OPENSTACK]: 'https://coreos.com/tectonic/docs/latest/install/openstack/openstack-terraform.html',
};
