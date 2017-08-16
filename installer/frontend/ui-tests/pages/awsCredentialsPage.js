const util = require('util');

const awsCredentialsPageCommands = {
  enterAwsCredentials() {
    return this
      .setValue('@awsAccessKey', process.env.AWS_ACCESS_KEY_ID)
      .setValue('@secretAccesskey', process.env.AWS_SECRET_ACCESS_KEY);
  },
};

module.exports = {
  url: '',
  commands: [
    awsCredentialsPageCommands, {
      el: function(elementName, data) {
        const element = this.elements[elementName.slice(1)];
        return util.format(element.selector, data);
      },
    }, {
      nextStep() {
        return this
          .click('@nextStep');
      },
    },
  ],
  elements: {
    awsAccessKey: {
      selector: 'input[id=accessKeyId]',
    },
    secretAccesskey: {
      selector: 'input[id=secretAccessKey]',
    },
    region: "option[value='%s']",
    nextStep: {
      selector: '//*[text()[contains(.,"Next Step")]]',
      locateStrategy: 'xpath',
    },
  },
};
