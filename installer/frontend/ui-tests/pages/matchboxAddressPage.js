const installerInput = require('../utils/bareMetalInstallerInput');
const inputJson = installerInput.buildExpectedJson();

const matchboxPageCommands = {
  enterMatchBoxEndPoints() {
    return this
      .setValue('@matchboxHTTP', inputJson.tectonic_metal_matchbox_http_url.replace(/^http?:\/\//i, ''))
      .setValue('@matchboxRPC', inputJson.tectonic_metal_matchbox_rpc_endpoint)
      .waitForElementPresent('@nextStep', 6000).click('@nextStep');
  },
};

module.exports = {
  url: '',
  commands: [matchboxPageCommands],
  elements: {
    matchboxHTTP: {
      selector: 'input#matchboxHTTP',
    },
    matchboxRPC: {
      selector: 'input#matchboxRPC',
    },
    nextStepp: {
      selector: '//*[text()[contains(.,"Next Step")]]',
      locateStrategy: 'xpath',
    },
    nextStep: {
      selector: '//*[@class="withtooltip"]/button[not(contains(@class, "btn-primary disabled"))]',
      locateStrategy: 'xpath',
    },
  },
};
