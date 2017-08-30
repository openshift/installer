const log = require('../utils/log');
const wizard = require('../utils/wizard');
const installerInput = require('../utils/awsInstallerInput');
const tfvarsUtil = require('../utils/terraformTfvars');

const json = installerInput.buildExpectedJson();
const testPage = (page, nextInitiallyDisabled) => wizard.testPage(page, json, nextInitiallyDisabled);

const REQUIRED_ENV_VARS = ['AWS_ACCESS_KEY_ID', 'AWS_SECRET_ACCESS_KEY', 'TF_VAR_tectonic_license_path', 'TF_VAR_tectonic_pull_secret_path'];

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

  'AWS: Platform': client => {
    const platformPage = client.page.platformPage();
    platformPage.navigate(client.launch_url);
    platformPage.test('@awsPlatform');
    platformPage.expect.element(wizard.nextStep).to.not.have.attribute('class').which.contains('disabled');
    platformPage.click(wizard.nextStep);
  },

  'AWS: AWS Credentials': ({page}) => testPage(page.awsCredentialsPage()),
  'AWS: Cluster Info': ({page}) => testPage(page.clusterInfoPage()),
  'AWS: Certificate Authority': ({page}) => testPage(page.certificateAuthorityPage(), false),
  'AWS: SSH Key': ({page}) => testPage(page.keysPage()),
  'AWS: Define Nodes': ({page}) => testPage(page.nodesPage(), false),
  'AWS: Networking': ({page}) => testPage(page.networkingPage()),
  'AWS: Console Login': ({page}) => testPage(page.consoleLoginPage()),

  'AWS: Submit': client => {
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
