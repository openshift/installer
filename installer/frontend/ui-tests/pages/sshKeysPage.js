const fs = require('fs');
const path = require('path');

const pageCommands = {
  test(json) {
    const parentDir = path.resolve(__dirname, '..');
    const sshKeyPath = path.join(parentDir, 'ssh-keys.txt');

    /* eslint-disable no-sync */
    fs.writeFileSync(sshKeyPath, json.tectonic_ssh_authorized_key);

    this
      .setField('@key', 'abc')
      .expectValidationErrorContains('Invalid SSH pubkey')
      .setValue('@upload', sshKeyPath)
      .expectNoValidationError();
  },
};

module.exports = {
  commands: [pageCommands],
  elements: {
    key: 'textarea#sshAuthorizedKey',
    upload: {
      selector: '(//*[text()="Upload"]/input[@type="file"])',
      locateStrategy: 'xpath',
    },
  },
};
