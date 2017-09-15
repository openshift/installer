const pageCommands = {
  test(json) {
    this.click('@deleteIcon0');
    this.expect.element('@alertError').text.to.contain('At least 1 Master is required');

    this
      .click('@addMore')
      .setField('@mac0', 'abc')
      .expectValidationErrorContains('Invalid MAC address')
      .setField('@mac0', json.tectonic_metal_controller_macs[0])
      .setField('@hosts0', '%')
      .expectValidationErrorContains('Invalid format')
      .setField('@hosts0', json.tectonic_metal_controller_domains[0])
      .expectNoValidationError()
      .click('@addMore')
      .setField('@mac1', json.tectonic_metal_controller_macs[0])
      .expectValidationErrorContains('MACs must be unique')
      .click('@deleteIcon1')
      .click('@addMore')
      .setField('@hosts1', json.tectonic_metal_controller_domains[0])
      .expectValidationErrorContains('Hostnames must be unique')
      .click('@deleteIcon1')
      .expectNoValidationError();
  },
};

module.exports = {
  commands: [pageCommands],
  elements: {
    mac0: 'input[id="masters.0.mac"]',
    hosts0: 'input[id="masters.0.host"]',
    deleteIcon0: '.row:nth-child(1) i.fa-minus-circle',
    mac1: 'input[id="masters.1.mac"]',
    hosts1: 'input[id="masters.1.host"]',
    deleteIcon1: '.row:nth-child(2) i.fa-minus-circle',
    addMore: {
      selector: '//*[text()[contains(.,"Add More")]]',
      locateStrategy: 'xpath',
    },
    alertError: '.alert-error',
  },
};
