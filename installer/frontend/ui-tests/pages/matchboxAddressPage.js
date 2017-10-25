const pageCommands = {
  test(json) {
    this
      .setField('@matchboxHTTP', 'abc')
      .expectValidationErrorContains('Invalid format')
      .setField('@matchboxHTTP', json.matchboxHTTP)
      .expectNoValidationError()
      .setField('@matchboxRPC', 'abc')
      .expectValidationErrorContains('Invalid format')
      .setField('@matchboxRPC', json.matchboxRPC)
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
