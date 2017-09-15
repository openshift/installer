const pageCommands = {
  test() {
    this.expect.element('@optionGenerate').to.be.selected;
    this.expect.element('@cert').to.not.be.present;
    this.expect.element('@key').to.not.be.present;

    this.selectOption('@optionUseOwn');
    this.expect.element('@cert').to.be.visible;
    this.expect.element('@key').to.be.visible;

    this.selectOption('@optionGenerate');
    this.expect.element('@cert').to.not.be.present;
    this.expect.element('@key').to.not.be.present;
  },
};

module.exports = {
  commands: [pageCommands],
  elements: {
    cert: 'textarea#caCertificate',
    key: 'textarea#caPrivateKey',
    optionGenerate: '.wiz-radio-group:nth-child(1) input[type=radio]',
    optionUseOwn: '.wiz-radio-group:nth-child(2) input[type=radio]',
  },
};
