const log = require('../utils/log');
const wizard = require('../utils/wizard');
const tfvarsUtil = require('../utils/terraformTfvars');

// Test input .progress file
const input = tfvarsUtil.loadJson('tectonic-baremetal.progress').clusterConfig;

// Expected Terraform tfvars file variables
const expected = tfvarsUtil.loadJson('metal.json').variables;

const testPage = (page, nextInitiallyDisabled) => wizard.testPage(page, 'metal', input, nextInitiallyDisabled);

const REQUIRED_ENV_VARS = ['TF_VAR_tectonic_license_path', 'TF_VAR_tectonic_pull_secret_path'];

module.exports = {
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

  'BM: Platform': client => {
    const platformPage = client.page.platformPage();
    platformPage.test('@metalGUI');
    platformPage.expect.element(wizard.nextStep).to.have.attribute('class').which.not.contains('disabled');
    platformPage.click(wizard.nextStep);
  },

  'BM: Cluster Info': ({page}) => testPage(page.clusterInfoPage()),
  'BM: Cluster DNS': ({page}) => testPage(page.clusterDnsPage()),
  'BM: Certificate Authority': ({page}) => testPage(page.certificateAuthorityPage(), false),
  'BM: Matchbox Address': ({page}) => testPage(page.matchboxAddressPage()),
  'BM: Matchbox Credentials': ({page}) => testPage(page.matchboxCredentialsPage()),
  'BM: Define Masters': ({page}) => testPage(page.defineMastersPage()),
  'BM: Define Workers': ({page}) => testPage(page.defineWorkersPage()),
  'BM: Network Configuration': ({page}) => testPage(page.networkConfigurationPage(), false),
  'BM: etcd Connection': ({page}) => testPage(page.etcdConnectionPage(), false),
  'BM: SSH Key': ({page}) => testPage(page.sshKeysPage()),
  'BM: Console Login': ({page}) => testPage(page.consoleLoginPage()),

  'BM: Submit': client => {
    const submitPage = client.page.submitPage();
    submitPage.click('@manuallyBoot');
    submitPage.expect.element('a[href="/terraform/assets"]').to.be.visible.before(120000);
    client.getCookie('tectonic-installer', result => {
      tfvarsUtil.returnTerraformTfvars(client.launch_url, result.value, (err, actualJson) => {
        if (err) {
          return client.assert.fail(err);
        }
        tfvarsUtil.assertDeepEqual(client, actualJson, expected);
      });
    });
  },
};
