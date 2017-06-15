const installerInput = require('../utils/installerInput');
const adminEmail = installerInput.adminEmail();
const adminPassword = installerInput.adminPassword();


const consoleLoginPageCommands = {
  enterLoginCredentails() {
    return this
      .setValue('@email', adminEmail)
      .setValue('@password', adminPassword)
      .setValue('@confirmPassword',adminPassword)
      .click('@nextStep');
  },
};

module.exports = {
  url: '',
  commands: [consoleLoginPageCommands],
  elements: {
    email: {
      selector: 'input[id=adminEmail]',
    },
    password: {
      selector: 'input[id=adminPassword]',
    },
    confirmPassword: {
      selector: 'input[id=adminPassword2]',
    },
    nextStep: {
      selector:'//*[text()[contains(.,"Next Step")]]',
      locateStrategy: 'xpath',
    },
  },
};
