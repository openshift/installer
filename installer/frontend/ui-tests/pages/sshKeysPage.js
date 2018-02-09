const pageCommands = {
  test (json) {
    this.testFileTextCombo(
      'input[type=file]#sshAuthorizedKey',
      'textarea#sshAuthorizedKey',
      json.sshAuthorizedKey,
      'abc',
      'Invalid SSH pubkey'
    );
  },
};

module.exports = {
  commands: [pageCommands],
  elements: {},
};
