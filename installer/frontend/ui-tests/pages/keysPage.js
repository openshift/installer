const installerInput = require('../utils/installerInput');

const sshKey = installerInput.sshKeys();

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
      selector: 'option[value=' + sshKey + ']',
    },
    nextStep: {
      selector:'//*[text()[contains(.,"Next Step")]]',
      locateStrategy: 'xpath',
    },
  },
};
