const fs = require('fs');
const path = require('path');

const clusterInfoPageCommands = {
  test(json, platform) {
    const parentDir = path.resolve(__dirname, '..');
    const licensePath = path.join(parentDir, 'tectonic-license.txt');
    const configPath = path.join(parentDir, 'config.json');

    /* eslint-disable no-sync */
    const tectonic_license = fs.readFileSync(process.env.TF_VAR_tectonic_license_path, 'utf8');
    const pull_secret = fs.readFileSync(process.env.TF_VAR_tectonic_pull_secret_path, 'utf8');
    fs.writeFileSync(licensePath, tectonic_license);
    fs.writeFileSync(configPath, pull_secret);
    /* eslint-enable no-sync */

    this.setField('@name', 'a%$#b');
    if (platform === 'aws') {
      this.expectValidationErrorContains('must be a valid AWS Stack Name');
    }
    if (platform === 'metal') {
      this.expectValidationErrorContains('must be alphanumeric');
    }

    this.setField('@name', json.tectonic_cluster_name);
    this.expectNoValidationError();

    this.setValue('@licenseUpload', licensePath);
    this.setValue('@pullSecretUpload', configPath);

    if (platform === 'aws') {
      this
        .setField('input[id="awsTags.0.key"]', '')
        .setField('input[id="awsTags.0.value"]', 'testing')
        .expectValidationErrorContains('Both fields are required')
        .setField('input[id="awsTags.0.key"]', 'test_tag')
        .setField('input[id="awsTags.0.value"]', '')
        .expectValidationErrorContains('Both fields are required')
        .setField('input[id="awsTags.0.key"]', 'test_tag')
        .setField('input[id="awsTags.0.value"]', 'testing')
        .expectNoValidationError();
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
