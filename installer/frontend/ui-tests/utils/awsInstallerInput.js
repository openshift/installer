const fs = require('fs');
const path = require('path');

const awsJsonPath = path.join(__dirname, '..', '..', '__tests__', 'examples', 'aws.json');
// eslint-disable-next-line no-sync
const awsTestDatajson = JSON.parse(fs.readFileSync(awsJsonPath, 'utf8'));
const awsProgressPath = path.join(__dirname, '..', '..', '__tests__', 'examples', 'tectonic-aws.progress');
// eslint-disable-next-line no-sync
const awsProgressData = JSON.parse(fs.readFileSync(awsProgressPath, 'utf8'));
const json = awsTestDatajson.variables;
/** Returns expected json. This json is used to prep the data required for the test */
const buildExpectedJson = () => {
  delete json.tectonic_admin_password_hash;
  delete json.tectonic_aws_extra_tags;
  delete json.tectonic_aws_external_private_zone;
  delete json.tectonic_aws_private_endpoints;

  return json;
};

const etcdOption = () => {
  return awsProgressData.clusterConfig.etcdOption;
};

module.exports = {
  buildExpectedJson,
  etcdOption,
};
