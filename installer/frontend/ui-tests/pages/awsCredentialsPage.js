const awsCredentialsPageCommands = {
  test(json) {
    const regionOption = `select#awsRegion option[value=${json.tectonic_aws_region}]`;
    return this
      .setField('@awsAccessKey', process.env.AWS_ACCESS_KEY_ID)
      .setField('@secretAccesskey', process.env.AWS_SECRET_ACCESS_KEY)
      .waitForElementPresent(regionOption, 60000)
      .click(regionOption)
      .expect.element(regionOption).to.be.selected;
  },
};

module.exports = {
  commands: [awsCredentialsPageCommands],
  elements: {
    awsAccessKey: 'input#accessKeyId',
    secretAccesskey: 'input#secretAccessKey',
  },
};
