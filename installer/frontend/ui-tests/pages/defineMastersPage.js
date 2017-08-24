const defineMastersPageCommands = {
  test(json) {
    return this
      .setValue('@masters0', json.tectonic_metal_controller_macs[0])
      .setValue('@hosts0', json.tectonic_metal_controller_domains[0]);
  },
};

module.exports = {
  commands: [defineMastersPageCommands],
  elements: {
    masters0: {
      selector: 'input[id="masters.0.mac"]',
    },
    hosts0: {
      selector: 'input[id="masters.0.host"]',
    },
  },
};
