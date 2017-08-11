const installerInput = require('../utils/bareMetalInstallerInput');
const inputJson = installerInput.buildExpectedJson();

const defineMastersPageCommands = {
  enterMastersDnsNames() {
    return this
      .setValue('@masters0', inputJson.tectonic_metal_controller_macs[0])
      .setValue('@hosts0', inputJson.tectonic_metal_controller_domains[0])
      .click('@nextStep');
  },
};

module.exports = {
  url: '',
  commands: [defineMastersPageCommands],
  elements: {
    masters0: {
      selector: 'input[id="masters.0.mac"]',
    },
    hosts0: {
      selector: 'input[id="masters.0.host"]',
    },
    nextStep: {
      selector: '//*[text()[contains(.,"Next Step")]]',
      locateStrategy: 'xpath',
    },
  },
};
