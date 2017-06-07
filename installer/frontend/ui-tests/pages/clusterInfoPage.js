const installerInput = require('../utils/installerInput');
const coreOSLicensePath = require('path').resolve(__dirname, '..') +
  '/tectonic-license.txt';
const configPath = require('path').resolve(__dirname, '..') + '/config.json';

const clusterInfoPageCommands = {
  enterClusterInfo(clusterName) {

    installerInput.createCoreosCredentials(coreOSLicensePath, process.env.TECTONIC_LICENSE);
    installerInput.createCoreosCredentials(configPath, process.env.PULL_SECRET);
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
