const clusterDnsPageCommands = {
  test(json) {
    return this
      .setValue('@controllerDomain', json.tectonic_metal_controller_domain)
      .setValue('@tectonicDomain', json.tectonic_metal_ingress_domain);
  },
};

module.exports = {
  commands: [clusterDnsPageCommands],
  elements: {
    controllerDomain: 'input#controllerDomain',
    tectonicDomain: 'input#tectonicDomain',
  },
};
