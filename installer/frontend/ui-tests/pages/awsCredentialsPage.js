const awsCredentialsPageCommands = {
  test(json) {
    const regionOption = `select#awsRegion option[value=${json.tectonic_aws_region}]`;

    if (process.env.AWS_SESSION_TOKEN) {
      this
        .click('@awsCredentialSessionTokenOption')
        .expect.element('@awsCredentialSessionTokenOption').to.be.selected;

      this
        .setField('@awsAccessKey', process.env.AWS_ACCESS_KEY_ID)
        .setField('@secretAccesskey', process.env.AWS_SECRET_ACCESS_KEY)
        .setField('@sessionToken', process.env.AWS_SESSION_TOKEN);
    } else {
      this
        .setField('@awsAccessKey', process.env.AWS_ACCESS_KEY_ID)
        .setField('@secretAccesskey', process.env.AWS_SECRET_ACCESS_KEY);
    }
    return this
      .waitForElementPresent(regionOption, 100000)
      .click(regionOption)
      .expect.element(regionOption).to.be.selected;
  },
};

module.exports = {
  commands: [awsCredentialsPageCommands],
  elements: {
    awsAccessKey: 'input#accessKeyId',
    secretAccesskey: 'input#secretAccessKey',
    sessionToken: 'input#awsSessionToken',
    awsCredentialSessionTokenOption: {
      selector: '(//*[@id="sts_enabled"])[2]',
      locateStrategy: 'xpath',
    },
  },
};
