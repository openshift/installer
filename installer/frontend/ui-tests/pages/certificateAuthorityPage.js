const pageCommands = {
  test (json) {
    const certFile = 'input[type="file"]#caCertificate';
    const certText = 'textarea#caCertificate';
    const keyFile = 'input[type="file"]#caPrivateKey';
    const keyText = 'textarea#caPrivateKey';
    const customCaFields = [certFile, certText, keyFile, keyText];

    this.expect.element('@optionGenerate').to.be.selected;
    customCaFields.forEach(el => this.expect.element(el).to.not.be.present);

    this.selectOption('@optionUseOwn');
    customCaFields.forEach(el => this.expect.element(el).to.be.present);
    this.expect.element(certText).to.be.visible;
    this.expect.element(keyText).to.be.visible;

    if (json.caCertificate && json.caPrivateKey) {
      this.testFileTextCombo(certFile, certText, json.caCertificate, '-----abc-----', 'Invalid certificate');
      this.testFileTextCombo(keyFile, keyText, json.caPrivateKey, '-----abc------', 'Invalid private key');
    } else {
      this.selectOption('@optionGenerate');
      customCaFields.forEach(el => this.expect.element(el).to.not.be.present);
    }
  },
};

module.exports = {
  commands: [pageCommands],
  elements: {
    optionGenerate: '.wiz-radio-group:nth-child(1) input[type=radio]',
    optionUseOwn: '.wiz-radio-group:nth-child(2) input[type=radio]',
  },
};
