const awsCredentialsPageCommands = {
  enterAwsCredentials(region) {
    const regionOption = `select#awsRegion option[value=${region}]`;
    this
      .setField('@awsAccessKey', process.env.AWS_ACCESS_KEY_ID)
      .setField('@secretAccesskey', process.env.AWS_SECRET_ACCESS_KEY)
      .waitForElementPresent(regionOption, 60000)
      .click(regionOption)
      .expect.element(regionOption).to.be.selected;
    return this.click('@nextStep');
  },
};

module.exports = {
  url: '',
  commands: [awsCredentialsPageCommands],
  elements: {
    awsAccessKey: {
      selector: 'input#accessKeyId',
    },
    secretAccesskey: {
      selector: 'input#secretAccessKey',
    },
    nextStep: {
      selector: '//*[text()[contains(.,"Next Step")]]',
      locateStrategy: 'xpath',
    },
  },
};
