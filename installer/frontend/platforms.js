export const AWS = 'aws';
export const AWS_TF = 'aws-tf';
export const AZURE = 'azure';
export const BARE_METAL = 'bare-metal';
export const OPENSTACK = 'openstack';

export const isSupported = platform => [AWS, BARE_METAL].includes(platform) ||
  (platform === AWS_TF && window.config.devMode);

export const PLATFORM_NAMES = {
  [AWS]: 'Amazon Web Services',
  [BARE_METAL]: 'Bare Metal',
  [AWS_TF]: 'Amazon Web Services (Terraform)',
  [AZURE]: 'Microsoft Azure (Terraform)',
  [OPENSTACK]: 'Openstack (Terraform)',
};

export const isTerraform = p => [AWS_TF, AZURE, OPENSTACK].includes(p);

export const OptGroups = [
  ['Graphical Installer (default)', AWS, BARE_METAL],
  ['Advanced Installer', AWS_TF, AZURE, OPENSTACK],
];

const platOrder = [
  AWS, BARE_METAL, AWS_TF, AZURE, OPENSTACK,
];

export const SELECTED_PLATFORMS = (window.config ? window.config.platforms : []).sort((a, b) =>
  platOrder.indexOf(a) - platOrder.indexOf(b)
);

export const DOCS = {
  [AWS]: 'https://coreos.com/tectonic/docs/latest/install/aws/index.html',
  [BARE_METAL]: 'https://coreos.com/tectonic/docs/latest/install/bare-metal/index.html',
  [AWS_TF]: 'https://coreos.com/tectonic/docs/latest/install/aws/aws-terraform.html',
  [AZURE]: 'https://coreos.com/tectonic/docs/latest/install/azure/azure-terraform.html',
  [OPENSTACK]: 'https://coreos.com/tectonic/docs/latest/install/openstack/openstack-terraform.html',
};
