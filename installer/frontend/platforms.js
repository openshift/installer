export const AWS = 'aws';
export const AWS_TF = 'aws-tf';
export const AZURE = 'azure';
export const BARE_METAL = 'bare-metal';
export const OPENSTACK = 'openstack';

export const isSupported = platform => [AWS, AWS_TF, BARE_METAL].includes(platform);

export const PLATFORM_NAMES = {
  [AWS]: 'Amazon Web Services (CloudFormation, legacy v1.5)',
  [AWS_TF]: 'Amazon Web Services (Terraform)',
  [AZURE]: 'Microsoft Azure (Terraform)',
  [BARE_METAL]: 'Bare Metal',
  [OPENSTACK]: 'Openstack (Terraform)',
};

export const isTerraform = p => [AWS_TF, AZURE, OPENSTACK].includes(p);

export const OptGroups = [
  ['Graphical Installer (default)', AWS_TF, AWS, BARE_METAL],
  ['Advanced Installer', AZURE, OPENSTACK],
];

const platOrder = [
  AWS_TF, AWS, BARE_METAL, AZURE, OPENSTACK,
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
