const log = require('../utils/log');
const wizard = require('../utils/wizard');
const installerInput = require('../utils/bareMetalInstallerInput');
const tfvarsUtil = require('../utils/terraformTfvars');

const json = installerInput.buildExpectedJson();
const testPage = (page, nextInitiallyDisabled) => wizard.testPage(page, json, nextInitiallyDisabled);

const REQUIRED_ENV_VARS = ['TF_VAR_tectonic_license_path', 'TF_VAR_tectonic_pull_secret_path'];

module.exports = {
  before () {
    const missing = REQUIRED_ENV_VARS.filter(ev => !process.env[ev]);
    if (missing.length) {
      console.error(`Missing environment variables: ${missing.join(', ')}.\n`);
      process.exit(1);
    }
  },

  after (client) {
    client.getLog('browser', log.logger);
    client.end();
  },

  'BM: Platform': client => {
    const platformPage = client.page.platformPage();
    platformPage.navigate(client.launch_url);
    platformPage.test('@bareMetalPlatform');
    platformPage.expect.element(wizard.nextStep).to.not.have.attribute('class').which.contains('disabled');
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
    submitPage.waitForElementPresent('a[href="/terraform/assets"]', 10000);
    client.getCookie('tectonic-installer', result => {
      tfvarsUtil.returnTerraformTfvars(client.launch_url, result.value, (err, actualJson) => {
        if (err) {
          return client.assert.fail(err);
        }
        const msg = tfvarsUtil.compareJson(actualJson, json);
        if (msg) {
          return client.assert.fail(msg);
        }
      });
    });
  },
};
