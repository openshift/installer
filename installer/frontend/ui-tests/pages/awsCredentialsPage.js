const awsCredentialsPageCommands = {
  test(json) {
    const regionOption = `select#awsRegion option[value=${json.tectonic_aws_region}]`;

    this.expect.element('@sessionTokenFalse').to.be.selected;

    if (process.env.AWS_SESSION_TOKEN) {
      this
        .setField('@sessionTokenTrue')
        .setField('@sessionToken', process.env.AWS_SESSION_TOKEN);
    }

    this
      .setField('@awsAccessKey', 'abc')
      .expectValidationErrorContains('AWS key IDs are at least 20 characters')
      .setField('@awsAccessKey', process.env.AWS_ACCESS_KEY_ID)
      .expectNoValidationError()
      .setField('@secretAccesskey', 'abc')
      .expectValidationErrorContains('AWS secrets are at least 40 characters')
      .setField('@secretAccesskey', process.env.AWS_SECRET_ACCESS_KEY)
      .expectNoValidationError();

    this.expect.element(regionOption).to.be.visible.before(60000);
    this.selectOption(regionOption);
  },
};

module.exports = {
  commands: [awsCredentialsPageCommands],
  elements: {
    awsAccessKey: 'input#accessKeyId',
    secretAccesskey: 'input#secretAccessKey',
    sessionToken: 'input#awsSessionToken',
    sessionTokenTrue: 'input#stsEnabledTrue',
    sessionTokenFalse: 'input#stsEnabledFalse',
  },
};
