const installerInput = require('../utils/installerInput');
const tfvarsUtil = require('../utils/terraformTfvars');

module.exports = {
  after (client) {
    client.end();
  },

  'Tectonic Installer Aws Test': (client) => {
    const options = {
      "tectonic_aws_region": "us-west-1",
    };
    const expectedJson = installerInput.buildExpectedJson(options);
    const platformPage = client.page.platformPage();
    const awsCredentialsPage = client.page.awsCredentialsPage();
    const clusterInfoPage = client.page.clusterInfoPage();
    const certificateAuthorityPage = client.page.certificateAuthorityPage();
    const keysPage = client.page.keysPage();
    const nodesPage = client.page.nodesPage();
    const networkingPage = client.page.networkingPage();
    const consoleLoginPage = client.page.consoleLoginPage();
    const submitPage = client.page.submitPage();

    platformPage.navigate(client.launch_url)
      .selectPlatform();
    awsCredentialsPage.enterAwsCredentials()
      .waitForElementPresent(awsCredentialsPage.el('@region', expectedJson.tectonic_aws_region), 60000)
      .click(awsCredentialsPage.el('@region', expectedJson.tectonic_aws_region))
      .nextStep();

    clusterInfoPage.enterClusterInfo(expectedJson.tectonic_cluster_name);
    certificateAuthorityPage.click('@nextStep');

    keysPage.selectSshKeys();
    nodesPage.waitForElementVisible('@nextStep', 10000).click('@nextStep');
    networkingPage.provideNetworkingDetails();
    consoleLoginPage.enterLoginCredentails();
    submitPage.click('@manuallyBoot').click('@manuallyBoot');
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
