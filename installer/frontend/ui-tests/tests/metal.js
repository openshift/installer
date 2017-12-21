const log = require('../utils/log');
const wizard = require('../utils/wizard');
const tfvarsUtil = require('../utils/terraformTfvars');

const REQUIRED_ENV_VARS = ['TF_VAR_tectonic_license_path', 'TF_VAR_tectonic_pull_secret_path'];

// Expects an input file <prefix>.progress to exist
const steps = (prefix, expectedOutputFilePath, ignoredKeys) => {
  // Test input cluster config
  const cc = tfvarsUtil.loadJson(`${prefix}.progress`).clusterConfig;

  const testPage = (page, nextInitiallyDisabled) => wizard.testPage(page, 'bare-metal-tf', cc, nextInitiallyDisabled);

  return {
    [`${prefix}: Platform`]: client => {
      const platformPage = client.page.platformPage();
      platformPage.test('@metalGUI');
      platformPage.expect.element(wizard.nextStep).to.have.attribute('class').which.not.contains('disabled');
      platformPage.click(wizard.nextStep);
    },

    [`${prefix}: Cluster Info`]: ({page}) => testPage(page.clusterInfoPage()),
    [`${prefix}: Cluster DNS`]: ({page}) => testPage(page.clusterDnsPage()),
    [`${prefix}: Certificate Authority`]: ({page}) => testPage(page.certificateAuthorityPage(), false),
    [`${prefix}: Matchbox Address`]: ({page}) => testPage(page.matchboxAddressPage()),
    [`${prefix}: Matchbox Credentials`]: ({page}) => testPage(page.matchboxCredentialsPage()),
    [`${prefix}: Define Masters`]: ({page}) => testPage(page.defineMastersPage()),
    [`${prefix}: Define Workers`]: ({page}) => testPage(page.defineWorkersPage()),
    [`${prefix}: Network Configuration`]: ({page}) => testPage(page.networkConfigurationPage(), false),
    [`${prefix}: etcd Connection`]: ({page}) => testPage(page.etcdConnectionPage(), false),
    [`${prefix}: SSH Key`]: ({page}) => testPage(page.sshKeysPage()),
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
  'tectonic_license_path',
  'tectonic_pull_secret_path',
  'tectonic_stats_url',
  'tectonic_update_app_id',
  'tectonic_update_channel',
  'tectonic_update_server',
];

module.exports = Object.assign(
  toExport,
  steps('metal', '../output/metal.tfvars'),
  steps('metal-smoke', '../../../../tests/smoke/bare-metal/vars/metal.tfvars.json', ignoredKeys)
);
