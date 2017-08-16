const installerInput = require('../utils/bareMetalInstallerInput');
const inputJson = installerInput.buildExpectedJson();

const networkingPageCommands = {
  provideNetworkingDetails() {
    return this
      .waitForElementPresent('@domain', 10000)
      .click('@domain')
      .click('@advanced')
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
  commands: [networkingPageCommands],
  elements: {
    advanced: {
      selector: '//*[text()[contains(.,"Advanced Settings")]]',
      locateStrategy: 'xpath',
    },
    domain: {
      selector: 'option[value=Z1ILIMNSJGTMO2]',
    },
    nextStep: {
      selector: '//*[text()[contains(.,"Next Step")]]',
      locateStrategy: 'xpath',
    },
  },
};
