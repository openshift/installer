const fs = require('fs');
const path = require('path');

const clusterInfoPageCommands = {
  test(json, platform) {
    const parentDir = path.resolve(__dirname, '..');
    const coreOSLicensePath = path.join(parentDir, 'tectonic-license.txt');
    const configPath = path.join(parentDir, 'config.json');

    /* eslint-disable no-sync */
    const tectonic_license = fs.readFileSync(process.env.TF_VAR_tectonic_license_path, 'utf8');
    const pull_secret = fs.readFileSync(process.env.TF_VAR_tectonic_pull_secret_path, 'utf8');
    fs.writeFileSync(coreOSLicensePath, tectonic_license);
    fs.writeFileSync(configPath, pull_secret);
    /* eslint-enable no-sync */

    this
      .setField('@name', json.tectonic_cluster_name)
      .setValue('@coreOSLicenseUpload', coreOSLicensePath)
      .setValue('@pullSecretUpload', configPath);
    if (platform === 'aws') {
      this
        .click('.list-add')
        .setField('input[id="awsTags.1.key"]', 'test_tag')
        .setField('input[id="awsTags.1.value"]', 'testing');
    }
    return this;
  },
};

module.exports = {
  commands: [clusterInfoPageCommands],
  elements: {
    name: 'input#clusterName',
    coreOSLicenseUpload: {
      selector: '//*[text()[contains(.,"tectonic-license.txt")]]/input[@type="file"]',
      locateStrategy: 'xpath',
    },
    pullSecretUpload: {
      selector: '//*[text()[contains(.,"config.json")]]/input[@type="file"]',
      locateStrategy: 'xpath',
    },
  },
};
