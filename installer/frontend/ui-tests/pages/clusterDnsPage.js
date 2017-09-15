const pageCommands = {
  test(json) {
    this
      .setField('@controllerDomain', '%')
      .expectValidationErrorContains('Invalid domain name')
      .setField('@controllerDomain', '')
      .expectValidationErrorContains('This field is required')
      .setField('@controllerDomain', json.tectonic_metal_controller_domain)
      .expectNoValidationError()
      .setField('@tectonicDomain', '%')
      .expectValidationErrorContains('Invalid domain name')
      .setField('@tectonicDomain', '')
      .expectValidationErrorContains('This field is required')
      .setField('@tectonicDomain', json.tectonic_metal_ingress_domain)
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
