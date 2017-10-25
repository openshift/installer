const keysPageCommands = {
  test(json) {
    this.selectOption(`option[value=${json.aws_ssh}]`);
  },
};

module.exports = {
  commands: [keysPageCommands],
  elements: {},
};
