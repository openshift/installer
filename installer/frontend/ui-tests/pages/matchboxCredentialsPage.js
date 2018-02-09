const pageCommands = {
  test (json) {
    this.testFileTextCombo(
      'input[type=file]#matchboxCA',
      'textarea#matchboxCA',
      json.matchboxCA,
      json.matchboxClientKey,
      'Invalid certificate'
    );
    this.testFileTextCombo(
      'input[type=file]#matchboxClientCert',
      'textarea#matchboxClientCert',
      json.matchboxClientCert,
      json.matchboxClientKey,
      'Invalid certificate'
    );
    this.testFileTextCombo(
      'input[type=file]#matchboxClientKey',
      'textarea#matchboxClientKey',
      json.matchboxClientKey,
      json.matchboxCA,
      'Invalid private key'
    );
  },
};

module.exports = {
  commands: [pageCommands],
  elements: {},
};
