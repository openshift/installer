const awsCredentialsPageCommands = {
  test(json) {
    const regionOption = `select#awsRegion option[value=${json.tectonic_aws_region}]`;

    this.expect.element('@sessionTokenFalse').to.be.selected;

    if (process.env.AWS_SESSION_TOKEN) {
      this
        .setField('@sessionTokenTrue')
        .setField('@awsAccessKey', process.env.AWS_ACCESS_KEY_ID)
        .setField('@secretAccesskey', process.env.AWS_SECRET_ACCESS_KEY)
        .setField('@sessionToken', process.env.AWS_SESSION_TOKEN);
    } else {
      this
        .setField('@awsAccessKey', process.env.AWS_ACCESS_KEY_ID)
        .setField('@secretAccesskey', process.env.AWS_SECRET_ACCESS_KEY);
    }

    this.expect.element(regionOption).to.be.visible.before(60000);
    return this.selectOption(regionOption);
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
