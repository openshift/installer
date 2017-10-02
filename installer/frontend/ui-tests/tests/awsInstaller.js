const fs = require('fs');
const path = require('path');

const log = require('../utils/log');
const wizard = require('../utils/wizard');
const tfvarsUtil = require('../utils/terraformTfvars');

// Expected Terraform tfvars file variables
// TODO: Confusingly, this is also used as input to the tests
const jsonPath = path.join(__dirname, '..', '..', '__tests__', 'examples', 'aws.json');
// eslint-disable-next-line no-sync
const json = JSON.parse(fs.readFileSync(jsonPath, 'utf8')).variables;

const testPage = (page, nextInitiallyDisabled) => wizard.testPage(page, 'aws', json, nextInitiallyDisabled);

const REQUIRED_ENV_VARS = ['AWS_ACCESS_KEY_ID', 'AWS_SECRET_ACCESS_KEY', 'TF_VAR_tectonic_license_path', 'TF_VAR_tectonic_pull_secret_path'];

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

  'AWS: Platform': client => {
    const platformPage = client.page.platformPage();
    platformPage.test('@awsGUI');
    platformPage.expect.element(wizard.nextStep).to.have.attribute('class').which.not.contains('disabled');
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
    submitPage.expect.element('a[href="/terraform/assets"]').to.be.visible.before(120000);
    client.getCookie('tectonic-installer', result => {
      tfvarsUtil.returnTerraformTfvars(client.launch_url, result.value, (err, actualJson) => {
        if (err) {
          return client.assert.fail(err);
        }
        tfvarsUtil.assertDeepEqual(client, actualJson, json);
      });
    });
  },
};
