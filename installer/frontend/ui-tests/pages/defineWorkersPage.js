const installerInput = require('../utils/bareMetalInstallerInput');
const inputJson = installerInput.buildExpectedJson();

const defineWorkersPageCommands = {
  enterWorkersDnsNames() {
    return this
      .setValue('@workers0', inputJson.tectonic_metal_worker_macs[0])
      .setValue('@hosts0', inputJson.tectonic_metal_worker_domains[0])
      .click('@addMore')
      .setValue('@workers1', inputJson.tectonic_metal_worker_macs[1])
      .setValue('@hosts1', inputJson.tectonic_metal_worker_domains[1])
      .click('@nextStep');
  },
};

module.exports = {
  url: '',
  commands: [defineWorkersPageCommands],
  elements: {
    workers0: {
      selector: 'input[id="workers.0.mac"]',
    },
    hosts0: {
      selector: 'input[id="workers.0.host"]',
    },
    workers1: {
      selector: 'input[id="workers.1.mac"]',
    },
    hosts1: {
      selector: 'input[id="workers.1.host"]',
    },
    addMore: {
      selector:'//*[text()[contains(.,"Add More")]]',
      locateStrategy: 'xpath',
    },
    nextStep: {
      selector:'//*[text()[contains(.,"Next Step")]]',
      locateStrategy: 'xpath',
    },
  },
};
