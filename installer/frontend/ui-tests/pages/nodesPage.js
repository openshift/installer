const installerInput = require('../utils/installerInput');

module.exports = {
  url: '',
  elements: {
    etcdOption: {
      selector:`#${installerInput.etcdOption()}`,
    },
    nextStep: {
      selector:'//*[text()[contains(.,"Next Step")]]',
      locateStrategy: 'xpath',
    },
  },
};
