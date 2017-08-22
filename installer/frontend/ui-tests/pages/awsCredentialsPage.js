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
    nextStep: {
      selector: '//*[text()[contains(.,"Next Step")]]',
      locateStrategy: 'xpath',
    },
  },
};
