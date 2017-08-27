const matchboxPageCommands = {
  test(json) {
    return this
      .setValue('@matchboxHTTP', json.tectonic_metal_matchbox_http_url.replace(/^https?:\/\//i, ''))
      .setValue('@matchboxRPC', json.tectonic_metal_matchbox_rpc_endpoint);
  },
};

module.exports = {
  commands: [matchboxPageCommands],
  elements: {
    matchboxHTTP: 'input#matchboxHTTP',
    matchboxRPC: 'input#matchboxRPC',
  },
};
