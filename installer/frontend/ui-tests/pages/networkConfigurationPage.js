const installerInput = require('../utils/bareMetalInstallerInput');
const inputJson = installerInput.buildExpectedJson();

const networkConfigurationPageCommands = {
  enterCIDRs() {
    return this
      .clearValue('#podCIDR')
      .setValue('#podCIDR', inputJson.tectonic_cluster_cidr)
      .assert.value('#podCIDR', inputJson.tectonic_cluster_cidr)
      .clearValue('#serviceCIDR')
      .setValue('#serviceCIDR', inputJson.tectonic_service_cidr)
      .assert.value('#serviceCIDR', inputJson.tectonic_service_cidr)
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
