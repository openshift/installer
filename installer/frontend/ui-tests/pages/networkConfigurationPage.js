const installerInput = require('../utils/bareMetalInstallerInput');
const inputJson = installerInput.buildExpectedJson();

const networkConfigurationPageCommands = {
  enterCIDRs() {
    return this
      .setField('#podCIDR', inputJson.tectonic_cluster_cidr, true)
      .setField('#serviceCIDR', inputJson.tectonic_service_cidr, true)
      .click('@nextStep');
  },
};

module.exports = {
  url: '',
  commands: [networkConfigurationPageCommands],
  elements: {
    nextStep: {
      selector: '//*[text()[contains(.,"Next Step")]]',
      locateStrategy: 'xpath',
    },
  },
};
