const awsCredentialsPageCommands = {
  enterAwsCredentials(region) {
    const regionOption = `select#awsRegion option[value=${region}]`;
    return this
      .setValue('@awsAccessKey', process.env.AWS_ACCESS_KEY_ID)
      .setValue('@secretAccesskey', process.env.AWS_SECRET_ACCESS_KEY)
      .waitForElementPresent(regionOption, 60000)
      .click(regionOption)
      .click('@nextStep');
  },
};

module.exports = {
  url: '',
  commands: [awsCredentialsPageCommands],
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
