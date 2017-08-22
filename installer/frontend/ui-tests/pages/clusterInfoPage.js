const fs = require('fs');
const path = require('path');

const clusterInfoPageCommands = {
  enterClusterInfo(clusterName) {
    const parentDir = path.resolve(__dirname, '..');
    const coreOSLicensePath = path.join(parentDir, 'tectonic-license.txt');
    const configPath = path.join(parentDir, 'config.json');

    /* eslint-disable no-sync */
    const tectonic_license = fs.readFileSync(process.env.TF_VAR_tectonic_license_path, 'utf8');
    const pull_secret = fs.readFileSync(process.env.TF_VAR_tectonic_pull_secret_path, 'utf8');
    fs.writeFileSync(coreOSLicensePath, tectonic_license);
    fs.writeFileSync(configPath, pull_secret);
    /* eslint-enable no-sync */

    return this
      .setValue('@name', clusterName)
      .setValue('@coreOSLicenseUpload', coreOSLicensePath)
      .setValue('@pullSecretUpload', configPath)
      .click('@nextStep');
  },
};

module.exports = {
  url: '',
  commands: [clusterInfoPageCommands],
  elements: {
    name: {
      selector: 'input#clusterName',
    },
    coreOSLicenseUpload: {
      selector: '//*[text()[contains(.,"tectonic-license.txt")]]/input[@type="file"]',
      locateStrategy: 'xpath',
    },
    pullSecretUpload: {
      selector: '//*[text()[contains(.,"config.json")]]/input[@type="file"]',
      locateStrategy: 'xpath',
    },
    nextStep: {
      selector: '//*[text()[contains(.,"Next Step")]]',
      locateStrategy: 'xpath',
    },
  },
};
