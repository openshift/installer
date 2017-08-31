const clusterDnsPageCommands = {
  test(json) {
    return this
      .setField('@controllerDomain', json.tectonic_metal_controller_domain)
      .setField('@tectonicDomain', json.tectonic_metal_ingress_domain);
  },
};

module.exports = {
  commands: [clusterDnsPageCommands],
  elements: {
    controllerDomain: 'input#controllerDomain',
    tectonicDomain: 'input#tectonicDomain',
  },
};
