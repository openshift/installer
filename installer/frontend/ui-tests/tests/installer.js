const installerInput = require('../utils/installerInput');
const tfvarsUtil = require('../utils/terraformTfvars');

const logger = logs => logs.forEach(log => {
  const { message, level } = log;
  switch (level) {
  case `DEBUG`:
    console.log("browser:", level, message);
    break;
  case `SEVERE`:
    console.warn("browser:", level, message);
    break;
  case `INFO`:
  default:
    console.info("browser:", level, message);
  }
});

module.exports = {
  after (client) {
    client.end();
  },

  'Tectonic Installer Aws Test': (client) => {
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
    client.getLog('browser', logger);

    awsCredentialsPage.enterAwsCredentials()
      .waitForElementPresent(awsCredentialsPage.el('@region', expectedJson.tectonic_aws_region), 60000)
      .click(awsCredentialsPage.el('@region', expectedJson.tectonic_aws_region))
      .nextStep();
    client.getLog('browser', logger);

    clusterInfoPage.enterClusterInfo(expectedJson.tectonic_cluster_name);
    certificateAuthorityPage.click('@nextStep');
    client.getLog('browser', logger);

    keysPage.selectSshKeys();
    nodesPage.waitForElementVisible('@nextStep', 10000).click('@nextStep');
    networkingPage.provideNetworkingDetails();
    consoleLoginPage.enterLoginCredentails();
    submitPage.click('@manuallyBoot').click('@manuallyBoot');
    client.getLog('browser', logger);
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
