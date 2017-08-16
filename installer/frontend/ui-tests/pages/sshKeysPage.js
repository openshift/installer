const fs = require('fs');
const path = require('path');

const installerInput = require('../utils/bareMetalInstallerInput');
const inputJson = installerInput.buildExpectedJson();

const sshKeysPageCommands = {
  enterPublicKey() {
    const parentDir = path.resolve(__dirname, '..');
    const sshKeysPath = path.join(parentDir, 'ssh-keys.txt');

    /* eslint-disable no-sync */
    fs.writeFileSync(sshKeysPath, inputJson.tectonic_ssh_authorized_key);

    return this
      .setValue('@publicKey', sshKeysPath)
      .click('@nextStep');
  },
};

module.exports = {
  url: '',
  commands: [sshKeysPageCommands],
  elements: {
    publicKey: {
      selector: '(//*[text()="Upload"]/input[@type="file"])',
      locateStrategy: 'xpath',
    },
    nextStep: {
      selector: '//*[text()[contains(.,"Next Step")]]',
      locateStrategy: 'xpath',
    },
  },
};
