const fs = require('fs');
const path = require('path');

const bareMetalJsonPath = path.join(__dirname, '..', '..', '__tests__', 'examples', 'metal.json');
// eslint-disable-next-line no-sync
const bareMetalTestDataJson = JSON.parse(fs.readFileSync(bareMetalJsonPath, 'utf8'));

const json = bareMetalTestDataJson.variables;

/** Returns expected json. This json is used to prep the data required for the test */
const buildExpectedJson = () => {
  delete json.tectonic_admin_password_hash;
  return json;
};

module.exports = {
  buildExpectedJson,
};
