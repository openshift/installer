const awsCredentialsPageCommands = {
  test(json) {
    const regionOption = `select#awsRegion option[value=${json.awsRegion}]`;

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

    const testInvalidCredentials = () => {
      this.setField('@awsAccessKey', '12345678901234567890');
      this.expect.element('@alertError').text.to.contain('not able to validate the provided access credentials');
      this.setField('@awsAccessKey', process.env.AWS_ACCESS_KEY_ID);
      this.expect.element('@alertError').to.not.be.present;
    };

    // Test that the error for invalid credentials is correctly shown both before and after a region is selected
    testInvalidCredentials();
    this.expect.element(regionOption).to.be.visible.before(60000);
    this.selectOption(regionOption);
    testInvalidCredentials();
  },
};

module.exports = {
  commands: [awsCredentialsPageCommands],
  elements: {
    alertError: '.alert-error',
    awsAccessKey: 'input#accessKeyId',
    secretAccesskey: 'input#secretAccessKey',
    sessionToken: 'input#awsSessionToken',
    sessionTokenTrue: 'input#stsEnabledTrue',
    sessionTokenFalse: 'input#stsEnabledFalse',
  },
};
