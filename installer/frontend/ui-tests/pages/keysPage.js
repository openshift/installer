const installerInput = require('../utils/awsInstallerInput');
const inputJson = installerInput.buildExpectedJson();

const keysPageCommands = {
  selectSshKeys() {
    return this
      .waitForElementPresent('@sshKeys', 10000)
      .click('@sshKeys')
      .waitForElementPresent('@nextStep', 10000)
      .click('@nextStep');

  },
};

module.exports = {
  url: '',
  commands: [keysPageCommands],
  elements: {
    sshKeys: {
      selector: 'option[value=' + inputJson.tectonic_aws_ssh_key + ']',
    },
    nextStep: {
      selector:'//*[text()[contains(.,"Next Step")]]',
      locateStrategy: 'xpath',
    },
  },
};
