const _ = require('lodash');
const fs = require('fs');

const clusterInfoPageCommands = {
  test (json, platform) {
    this.setField('@name', 'a%$#b');
    if (platform === 'aws-tf') {
      this.expectValidationErrorContains('must be a lower case AWS Stack Name');
    }
    if (platform === 'bare-metal-tf') {
      this.expectValidationErrorContains('must be alphanumeric');
    }

    this.setField('@name', json.clusterName);
    this.expectNoValidationError();

    /* eslint-disable no-sync */
    this.setFileField('@licenseUpload', fs.readFileSync(process.env.TF_VAR_tectonic_license_path, 'utf8'));
    this.setFileField('@pullSecretUpload', fs.readFileSync(process.env.TF_VAR_tectonic_pull_secret_path, 'utf8'));
    /* eslint-enable no-sync */

    if (platform === 'aws-tf' && !_.isEmpty(json.awsTags)) {
      this
        .setField('input[id="awsTags.0.key"]', 'abc')
        .setField('input[id="awsTags.0.value"]', 'abc')
        .expectNoValidationError()
        .setField('input[id="awsTags.0.key"]', '')
        .expectValidationErrorContains('Both fields are required')
        .setField('input[id="awsTags.0.key"]', 'abc')
        .expectNoValidationError()
        .setField('input[id="awsTags.0.value"]', '')
        .expectValidationErrorContains('Both fields are required')
        .setField('input[id="awsTags.0.key"]', '')
        .expectNoValidationError()
        .click('.fa-plus-circle')
        .click('.fa-plus-circle')
        .setField('input[id="awsTags.1.key"]', 'abc')
        .setField('input[id="awsTags.2.key"]', 'abc')
        .expectValidationErrorContains('Tag keys must be unique')
        .click('.fa-minus-circle')
        .click('.fa-minus-circle')
        .setField('input[id="awsTags.0.key"]', json.awsTags[0].key)
        .setField('input[id="awsTags.0.value"]', json.awsTags[0].value)
        .expectNoValidationError();

      for (let i = 1; i < json.awsTags.length; i++) {
        this
          .click('.fa-plus-circle')
          .setField('input[id="awsTags.1.key"]', json.awsTags[1].key)
          .setField('input[id="awsTags.1.value"]', json.awsTags[1].value)
          .expectNoValidationError();
      }
    }
  },
};

module.exports = {
  commands: [clusterInfoPageCommands],
  elements: {
    name: 'input#clusterName',
    licenseUpload: 'input[type="file"]#tectonicLicense',
    pullSecretUpload: 'input[type="file"]#pullSecret',
  },
};
