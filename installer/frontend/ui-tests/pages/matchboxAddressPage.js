const matchboxPageCommands = {
  test(json) {
    return this
      .setField('@matchboxHTTP', json.tectonic_metal_matchbox_http_url.replace(/^https?:\/\//i, ''))
      .setField('@matchboxRPC', json.tectonic_metal_matchbox_rpc_endpoint);
  },
};

module.exports = {
  commands: [matchboxPageCommands],
  elements: {
    matchboxHTTP: 'input#matchboxHTTP',
    matchboxRPC: 'input#matchboxRPC',
  },
};
