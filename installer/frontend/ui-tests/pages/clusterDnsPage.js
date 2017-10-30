const pageCommands = {
  test(json) {
    this
      .setField('@controllerDomain', '%')
      .expectValidationErrorContains('Invalid domain name')
      .setField('@controllerDomain', '')
      .expectValidationErrorContains('This field is required')
      .setField('@controllerDomain', json.controllerDomain)
      .expectNoValidationError()
      .setField('@tectonicDomain', '%')
      .expectValidationErrorContains('Invalid domain name')
      .setField('@tectonicDomain', '')
      .expectValidationErrorContains('This field is required')
      .setField('@tectonicDomain', json.tectonicDomain)
      .expectNoValidationError();
  },
};

module.exports = {
  commands: [pageCommands],
  elements: {
    controllerDomain: 'input#controllerDomain',
    tectonicDomain: 'input#tectonicDomain',
  },
};
