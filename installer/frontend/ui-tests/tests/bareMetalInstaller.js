const log = require('../utils/log');
const wizard = require('../utils/wizard');
const installerInput = require('../utils/bareMetalInstallerInput');
const tfvarsUtil = require('../utils/terraformTfvars');

const json = installerInput.buildExpectedJson();
const testPage = (page, nextInitiallyDisabled) => wizard.testPage(page, json, nextInitiallyDisabled);

const REQUIRED_ENV_VARS = ['TF_VAR_tectonic_license_path', 'TF_VAR_tectonic_pull_secret_path'];

module.exports = {
  after (client) {
    client.getLog('browser', log.logger);
    client.end();
  },

  'Tectonic Installer BareMetal Test': (client) => {
    const missing = REQUIRED_ENV_VARS.filter(ev => !process.env[ev]);
    if (missing.length) {
      console.error(`Missing environment variables: ${missing.join(', ')}.\n`);
      process.exit(1);
    }

    const platformPage = client.page.platformPage();
    platformPage.navigate(client.launch_url);
    platformPage.test('@bareMetalPlatform');
    platformPage.expect.element(wizard.nextStep).to.not.have.attribute('class').which.contains('disabled');
    platformPage.click(wizard.nextStep);

    testPage(client.page.clusterInfoPage());
    testPage(client.page.clusterDnsPage());
    testPage(client.page.certificateAuthorityPage(), false);
    testPage(client.page.matchboxAddressPage());
    testPage(client.page.matchboxCredentialsPage());
    testPage(client.page.networkConfigurationPage(), false);
    testPage(client.page.defineMastersPage());
    testPage(client.page.defineWorkersPage());
    testPage(client.page.etcdConnectionPage(), false);
    testPage(client.page.sshKeysPage());
    testPage(client.page.consoleLoginPage());

    client.page.submitPage().click('@manuallyBoot');
    client.pause(10000);
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
