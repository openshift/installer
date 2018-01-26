const testExternalEtcd = page => {
  const address = 'input[type=text]#externalETCDClient';
  const optionExternal = 'input[type=radio]#external';
  const optionProvisioned = 'input[type=radio]#provisioned';

  page.expect.element(optionProvisioned).to.be.selected;
  page.expect.element(address).to.not.be.present;

  page
    .selectOption(optionExternal)
    .setField(address, 'https://example.com')
    .expectValidationErrorContains('Invalid format')
    .setField(address, 'example.com:1234')
    .expectValidationErrorContains('Invalid format')
    .setField(address, 'example.com')
    .expectNoValidationError()
    .selectOption(optionProvisioned)
    .expect.element(address).to.not.be.present;
};

const pageCommands = {
  test () {
    testExternalEtcd(this);
  },
};

module.exports = {
  commands: [pageCommands],
  elements: {},
  testExternalEtcd,
};
