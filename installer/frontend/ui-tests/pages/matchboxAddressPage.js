const pageCommands = {
  test(json) {
    this
      .setField('@matchboxHTTP', 'abc')
      .expectValidationErrorContains('Invalid format')
      .setField('@matchboxHTTP', json.tectonic_metal_matchbox_http_url.replace(/^https?:\/\//i, ''))
      .expectNoValidationError()
      .setField('@matchboxRPC', 'abc')
      .expectValidationErrorContains('Invalid format')
      .setField('@matchboxRPC', json.tectonic_metal_matchbox_rpc_endpoint)
      .expectNoValidationError();
  },
};

module.exports = {
  commands: [pageCommands],
  elements: {
    matchboxHTTP: 'input#matchboxHTTP',
    matchboxRPC: 'input#matchboxRPC',
  },
};
