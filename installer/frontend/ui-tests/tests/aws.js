const log = require('../utils/log');
const wizard = require('../utils/wizard');
const tfvarsUtil = require('../utils/terraformTfvars');

const REQUIRED_ENV_VARS = ['AWS_ACCESS_KEY_ID', 'AWS_SECRET_ACCESS_KEY', 'TF_VAR_tectonic_license_path', 'TF_VAR_tectonic_pull_secret_path'];

// Expects an input file <prefix>.progress to exist
const steps = (prefix, expectedOutputFilePath, ignoredKeys) => {
  // Test input cluster config
  const cc = tfvarsUtil.loadJson(`${prefix}.progress`).clusterConfig;

  const testPage = (page, nextInitiallyDisabled) => wizard.testPage(page, 'aws-tf', cc, nextInitiallyDisabled);

  return {
    [`${prefix}: Platform`]: client => {
      const platformPage = client.page.platformPage();
      platformPage.test('@awsGUI');
      platformPage.expect.element(wizard.nextStep).to.have.attribute('class').which.not.contains('disabled');
      platformPage.click(wizard.nextStep);
    },

    [`${prefix}: AWS Credentials`]: ({page}) => testPage(page.awsCredentialsPage()),
    [`${prefix}: Cluster Info`]: ({page}) => testPage(page.clusterInfoPage()),
    [`${prefix}: Certificate Authority`]: ({page}) => testPage(page.certificateAuthorityPage(), false),
    [`${prefix}: SSH Key`]: ({page}) => testPage(page.keysPage()),
    [`${prefix}: Define Nodes`]: ({page}) => testPage(page.nodesPage(), false),
    [`${prefix}: Networking`]: ({page}) => testPage(page.networkingPage()),
    [`${prefix}: Console Login`]: ({page}) => testPage(page.consoleLoginPage()),

    [`${prefix}: Manual Boot`]: client => tfvarsUtil.testManualBoot(client, expectedOutputFilePath, ignoredKeys),
  };
};

const toExport = {
  before (client) {
    const missing = REQUIRED_ENV_VARS.filter(ev => !process.env[ev]);
    if (missing.length) {
      console.error(`Missing environment variables: ${missing.join(', ')}.\n`);
      process.exit(1);
    }
    client.url(client.launch_url);
  },

  after (client) {
    client.getLog('browser', log.logger);
    client.end();
  },
};

const ignoredKeys = [
  'tectonic_admin_email',
  'tectonic_admin_password',
  'tectonic_aws_ssh_key',
  'tectonic_license_path',
  'tectonic_pull_secret_path',
  'tectonic_stats_url',
];

module.exports = Object.assign(
  toExport,
  steps('aws', '../../../../tests/smoke/aws/vars/aws.tfvars.json', ignoredKeys),
  steps('aws-custom-vpc', '../output/aws-custom-vpc.tfvars')
);
