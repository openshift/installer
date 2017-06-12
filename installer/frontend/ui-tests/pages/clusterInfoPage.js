const fs = require('fs');
const installerInput = require('../utils/installerInput');
const coreOSLicensePath = require('path').resolve(__dirname, '..') +
  '/tectonic-license.txt';
const configPath = require('path').resolve(__dirname, '..') + '/config.json';
  //eslint-disable-next-line no-sync
const tectonic_license = fs.readFileSync(process.env.TF_VAR_tectonic_license_path, 'utf8');
  //eslint-disable-next-line no-sync
const pull_secret = fs.readFileSync(process.env.TF_VAR_tectonic_pull_secret_path, 'utf8');


const clusterInfoPageCommands = {
  enterClusterInfo(clusterName) {

    installerInput.createCoreosCredentials(coreOSLicensePath, tectonic_license);
    installerInput.createCoreosCredentials(configPath, pull_secret);
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
      selector: 'input[id=clusterName]',
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
      selector:'//*[text()[contains(.,"Next Step")]]',
      locateStrategy: 'xpath',
    },
  },
};
