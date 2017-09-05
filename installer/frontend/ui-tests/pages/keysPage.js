const keysPageCommands = {
  test(json) {
    const sshKeys = `option[value=${json.tectonic_aws_ssh_key}]`;
    return this
      .waitForElementPresent(sshKeys)
      .click(sshKeys);
  },
};

module.exports = {
  commands: [keysPageCommands],
  elements: {},
};
