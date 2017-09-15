const pageCommands = {
  test() {
    this.expect.element('@optionProvisioned').to.be.selected;
    this.expect.element('@address').to.not.be.present;

    this.selectOption('@optionSelfHosted');
    this.expect.element('@address').to.not.be.present;

    this
      .selectOption('@optionExternal')
      .setField('@address', 'abc')
      .expectValidationErrorContains('Invalid format')
      .setField('@address', 'example.com:1234')
      .expectNoValidationError();

    this.selectOption('@optionProvisioned');
    this.expect.element('@address').to.not.be.present;
  },
};

module.exports = {
  commands: [pageCommands],
  elements: {
    address: 'input[type=text]#externalETCDClient',
    optionProvisioned: 'input[type=radio]#provisioned',
    optionSelfHosted: 'input[type=radio]#selfHosted',
    optionExternal: 'input[type=radio]#external',
  },
};
