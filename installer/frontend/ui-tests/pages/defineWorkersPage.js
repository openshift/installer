const defineWorkersPageCommands = {
  test(json) {
    this
      .setField('@workers0', json.tectonic_metal_worker_macs[0])
      .setField('@hosts0', json.tectonic_metal_worker_domains[0])
      .click('@addMore')
      .setField('@workers1', json.tectonic_metal_worker_macs[1])
      .setField('@hosts1', json.tectonic_metal_worker_domains[1]);
  },
};

module.exports = {
  commands: [defineWorkersPageCommands],
  elements: {
    workers0: 'input[id="workers.0.mac"]',
    hosts0: 'input[id="workers.0.host"]',
    workers1: 'input[id="workers.1.mac"]',
    hosts1: 'input[id="workers.1.host"]',
    addMore: {
      selector: '//*[text()[contains(.,"Add More")]]',
      locateStrategy: 'xpath',
    },
  },
};
