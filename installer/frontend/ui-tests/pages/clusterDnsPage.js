const installerInput = require('../utils/bareMetalInstallerInput');
const inputJson = installerInput.buildExpectedJson();

const clusterDnsPageCommands = {
  enterDnsNames() {
    return this
      .setValue('@controllerDomain', inputJson.tectonic_metal_controller_domain)
      .setValue('@tectonicDomain', inputJson.tectonic_metal_ingress_domain)
      .click('@nextStep');
  },
};

module.exports = {
  url: '',
  commands: [clusterDnsPageCommands],
  elements: {
    controllerDomain: {
      selector: 'input[id=controllerDomain]',
    },
    tectonicDomain: {
      selector: 'input[id=tectonicDomain]',
    },
    nextStep: {
      selector: '//*[text()[contains(.,"Next Step")]]',
      locateStrategy: 'xpath',
    },
  },
};
