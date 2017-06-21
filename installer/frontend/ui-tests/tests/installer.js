const _ = require('lodash');

const installerInput = require('../utils/installerInput');
const tfvarsUtil = require('../utils/terraformTfvars');

const logger = logs => {
  console.log('==== BEGIN BROWSER LOGS ====');
  _.each(logs, log => {
    const { level, message } = log;
    const messageStr = _.isArray(message) ? message.join(" ") : message;

    switch (level) {
    case `DEBUG`:
      console.log(level, messageStr);
      break;
    case `SEVERE`:
      console.warn(level, messageStr);
      break;
    case `INFO`:
    default:
      console.info(level, messageStr);
    }
  });
  console.log('==== END BROWSER LOGS ====');
};

const REQUIRED_ENV_VARS = ["AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "TF_VAR_tectonic_license_path", "TF_VAR_tectonic_pull_secret_path"];

module.exports = {
  after (client) {
    client.getLog('browser', logger);
    client.end();
  },

  'Tectonic Installer AWS Test': (client) => {
    const missing = _.filter(REQUIRED_ENV_VARS, ev => !process.env[ev]);
    if (missing.length) {
      console.error(`Missing environment variables: ${missing.join(', ')}.\n`);
      process.exit(1);
    }

    const expectedJson = installerInput.buildExpectedJson();
    const platformPage = client.page.platformPage();
    const awsCredentialsPage = client.page.awsCredentialsPage();
    const clusterInfoPage = client.page.clusterInfoPage();
    const certificateAuthorityPage = client.page.certificateAuthorityPage();
    const keysPage = client.page.keysPage();
    const nodesPage = client.page.nodesPage();
    const networkingPage = client.page.networkingPage();
    const consoleLoginPage = client.page.consoleLoginPage();
    const submitPage = client.page.submitPage();

    platformPage.navigate(client.launch_url).selectPlatform();
    awsCredentialsPage.enterAwsCredentials()
      .waitForElementPresent(awsCredentialsPage.el('@region', expectedJson.tectonic_aws_region), 60000)
      .click(awsCredentialsPage.el('@region', expectedJson.tectonic_aws_region))
      .nextStep();

    clusterInfoPage.enterClusterInfo(expectedJson.tectonic_cluster_name);
    certificateAuthorityPage.click('@nextStep');

    keysPage.selectSshKeys();
    nodesPage.click('@etcdOption')
      .waitForElementVisible('@nextStep', 10000)
      .click('@nextStep');
    networkingPage.provideNetworkingDetails();
    consoleLoginPage.enterLoginCredentails();
    submitPage.click('@manuallyBoot');
    client.pause(10000);
    client.getCookie('tectonic-installer', result => {
      tfvarsUtil.returnTerraformTfvars(client.launch_url, result.value, (err, actualJson) => {
        if (err) {
          return client.assert.fail(err);
        }
        const msg = tfvarsUtil.compareJson(actualJson, expectedJson);
        if (msg) {
          return client.assert.fail(msg);
        }
      });
    });
  },
};
