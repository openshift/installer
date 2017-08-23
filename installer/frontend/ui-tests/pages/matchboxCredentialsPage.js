const fs = require('fs');
const path = require('path');

const installerInput = require('../utils/bareMetalInstallerInput');
const inputJson = installerInput.buildExpectedJson();

const matchBoxCredentialsPageCommands = {
  enterMatchBoxCredentials() {
    const parentDir = path.resolve(__dirname, '..');
    const caCertPath = path.join(parentDir, 'ca-cert.txt');
    const clientCertPath = path.join(parentDir, 'client-cert.txt');
    const clientKeyPath = path.join(parentDir, 'client-key.txt');

    /* eslint-disable no-sync */
    fs.writeFileSync(caCertPath, inputJson.tectonic_metal_matchbox_ca);
    fs.writeFileSync(clientCertPath, inputJson.tectonic_metal_matchbox_client_cert);
    fs.writeFileSync(clientKeyPath, inputJson.tectonic_metal_matchbox_client_key);
    /* eslint-enable no-sync */
    return this
      .setValue('@caCertificate', caCertPath)
      .setValue('@clientCertificate', clientCertPath)
      .setValue('@clientKey', clientKeyPath)
      .waitForElementVisible('@nextStep', 6000).click('@nextStep');
  },
};

module.exports = {
  url: '',
  commands: [matchBoxCredentialsPageCommands],
  elements: {
    name: {
      selector: 'input#clusterName',
    },
    caCertificate: {
      selector: '(//*[text()="Upload"]/input[@type="file"])[1]',
      locateStrategy: 'xpath',
    },
    clientCertificate: {
      selector: '(//*[text()="Upload"]/input[@type="file"])[2]',
      locateStrategy: 'xpath',
    },
    clientKey: {
      selector: '(//*[text()="Upload"]/input[@type="file"])[3]',
      locateStrategy: 'xpath',
    },
    nextStep: {
      selector: '//*[@class="withtooltip"]/button[not(contains(@class, "btn-primary disabled"))]',
      locateStrategy: 'xpath',
    },
  },
};
