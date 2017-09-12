const fs = require('fs');
const path = require('path');

const sshKeysPageCommands = {
  test(json) {
    const parentDir = path.resolve(__dirname, '..');
    const sshKeysPath = path.join(parentDir, 'ssh-keys.txt');

    /* eslint-disable no-sync */
    fs.writeFileSync(sshKeysPath, json.tectonic_ssh_authorized_key);

    this.setValue('@publicKey', sshKeysPath);
  },
};

module.exports = {
  commands: [sshKeysPageCommands],
  elements: {
    publicKey: {
      selector: '(//*[text()="Upload"]/input[@type="file"])',
      locateStrategy: 'xpath',
    },
  },
};
