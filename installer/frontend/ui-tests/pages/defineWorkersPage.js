const pageCommands = {
  test(json) {
    this.click('@deleteIcon0');
    this.expect.element('@alertError').text.to.contain('At least 1 Worker is required');

    this
      .click('@addMore')
      .setField('@mac0', 'abc')
      .expectValidationErrorContains('Invalid MAC address')
      .setField('@mac0', json.workers[0].mac)
      .setField('@hosts0', '%')
      .expectValidationErrorContains('Invalid format')
      .setField('@hosts0', json.workers[0].host)
      .expectNoValidationError()
      .click('@addMore')
      .setField('@mac1', json.workers[0].mac)
      .expectValidationErrorContains('MACs must be unique')
      .setField('@mac1', json.workers[1].mac)
      .setField('@hosts1', json.workers[0].host)
      .expectValidationErrorContains('Hostnames must be unique')
      .setField('@hosts1', json.workers[1].host)
      .expectNoValidationError();
  },
};

module.exports = {
  commands: [pageCommands],
  elements: {
    mac0: 'input[id="workers.0.mac"]',
    hosts0: 'input[id="workers.0.host"]',
    deleteIcon0: '.row:nth-child(1) i.fa-minus-circle',
    mac1: 'input[id="workers.1.mac"]',
    hosts1: 'input[id="workers.1.host"]',
    addMore: {
      selector: '//*[text()[contains(.,"Add More")]]',
      locateStrategy: 'xpath',
    },
    alertError: '.alert-error',
  },
};
