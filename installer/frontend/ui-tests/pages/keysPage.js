const keysPageCommands = {
  test(json) {
    return this.selectOption(`option[value=${json.tectonic_aws_ssh_key}]`);
  },
};

module.exports = {
  commands: [keysPageCommands],
  elements: {},
};
