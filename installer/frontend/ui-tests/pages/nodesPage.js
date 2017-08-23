const installerInput = require('../utils/awsInstallerInput');

const nodesPageCommands = {
  test() {
    this.click('@etcdOption');
  },
};

module.exports = {
  url: '',
  commands: [nodesPageCommands],
  elements: {
    etcdOption: {
      selector: `#${installerInput.etcdOption()}`,
    },
  },
};
