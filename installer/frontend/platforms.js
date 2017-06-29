export const AWS_TF = 'aws-tf';
export const AZURE = 'azure';
export const BARE_METAL_TF = 'bare-metal-tf';
export const OPENSTACK = 'openstack';

export const isSupported = platform => [AWS_TF, BARE_METAL_TF].includes(platform);

export const PLATFORM_NAMES = {
  [AWS_TF]: 'Amazon Web Services',
  [AZURE]: 'Microsoft Azure',
  [BARE_METAL_TF]: 'Bare Metal',
  [OPENSTACK]: 'Openstack',
};

export const isTerraform = p => [AWS_TF, AZURE, BARE_METAL_TF, OPENSTACK].includes(p);

export const OptGroups = [
  ['Graphical Installer (default)', AWS_TF, BARE_METAL_TF],
  ['Advanced Installer', AZURE, OPENSTACK],
];

const platOrder = [
  AWS_TF, BARE_METAL_TF, AZURE, OPENSTACK,
];

export const SELECTED_PLATFORMS = (window.config ? window.config.platforms : []).sort((a, b) =>
  platOrder.indexOf(a) - platOrder.indexOf(b)
);

export const DOCS = {
  [BARE_METAL_TF]: 'https://coreos.com/tectonic/docs/latest/install/bare-metal/index.html',
  [AWS_TF]: 'https://coreos.com/tectonic/docs/latest/install/aws/index.html',
  [AZURE]: 'https://coreos.com/tectonic/docs/latest/install/azure/azure-terraform.html',
  [OPENSTACK]: 'https://coreos.com/tectonic/docs/latest/install/openstack/openstack-terraform.html',
};
