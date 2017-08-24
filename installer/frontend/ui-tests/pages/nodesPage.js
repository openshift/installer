const installerInput = require('../utils/awsInstallerInput');

const nodesPageCommands = {
  test() {
    this.click('@etcdOption');
  },
};

module.exports = {
  commands: [nodesPageCommands],
  elements: {
    etcdOption: `#${installerInput.etcdOption()}`,
  },
};
