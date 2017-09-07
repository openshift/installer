const request = require('request');
const JSZip = require('jszip');
const deep = require('deep-diff').diff;

const getAssets = (launchUrl, cookie, callback) => {
  const options = {
    url: launchUrl + '/terraform/assets',
    method: 'GET',
    encoding: null,
    headers: {
      'Cookie': 'tectonic-installer=' + cookie,
    },
  };
  request(options, (err, res, body) => {
    if (err) {
      return callback(err, res);
    }
    if (res.statusCode !== 200 || res.headers['content-type'] !== 'application/zip') {
      return callback({
        statusCode: res.statusCode,
        contentType: res.headers['content-type'],
      }, res);
    }
    return callback(null, res, body);
  });
};

const getTerraformTfvars = (response, callback) => {
  let fileName;
  JSZip.loadAsync(response).then(zip => {
    Object.keys(zip.files).forEach(key => {
      if (/tfvars$/.test(key)) {
        fileName = key;
      }
    });
    zip.file(fileName).async('string').then(callback);
  });
};

const returnRequiredTerraformTfvars = (terraformTfvars) => {
  const json = JSON.parse(terraformTfvars);
  const extraTfvars = [
    'tectonic_license_path',
    'tectonic_pull_secret_path',
    'tectonic_kube_apiserver_service_ip',
    'tectonic_kube_dns_service_ip',
    'tectonic_kube_etcd_service_ip',
  ];
  extraTfvars.forEach(key => {
    delete json[key];
  });
  return json;
};

const returnTerraformTfvars = (launchUrl, cookie, callback) => {
  getAssets(launchUrl, cookie, (err, res, terraformAssestsResponse) => {
    if (err) {
      return callback(err);
    }
    if (res.statusCode !== 200 || res.headers['content-type'] !== 'application/zip' ) {
      return callback('Terraform get assets api call failed', res);
    }
    getTerraformTfvars(terraformAssestsResponse, (terraformTfvars) => {
      const actualJson = returnRequiredTerraformTfvars(terraformTfvars);
      callback(null, actualJson);
    });
  });
};

const assertDeepEqual = (client, actual, expected) => {
  // The password hash will be different every time, so can't diff
  delete actual.tectonic_admin_password_hash;
  delete expected.tectonic_admin_password_hash;

  const diff = deep(actual, expected);
  if (diff !== undefined) {
    client.assert.fail(
      'The following terraform.tfvars attributes differ from their expected value: ' +
      diff.map(({path, lhs, rhs}) => `${path} (expected: ${rhs}, got: ${lhs})`).join(', '),
    );
  }
};

module.exports = {
  returnTerraformTfvars,
  assertDeepEqual,
};
