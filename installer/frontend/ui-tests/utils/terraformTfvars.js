const request = require('request');
const JSZip = require("jszip");
const deep = require('deep-diff').diff;

const getAssets = (launchUrl,cookie, callback) => {
  const options = {
    url: launchUrl + '/terraform/assets',
    method: 'GET',
    encoding: null,
    headers: {
      'Cookie': "tectonic-installer=" + cookie,
    },
  };
  request(options, (err, res, body) => {
    if (err || res.statusCode !== 200 || res.headers['content-type'] !== 'application/zip') {
      return callback(err || {
        statusCode: res.statusCode,
        contentType: res.headers['content-type'],
      }, null);
    }
    callback(null, res, body);
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
    zip.file(fileName).async("string").then(callback);
  });
};

const returnRequiredTerraformTfvars = (terraformTfvars) => {
  const json = JSON.parse(terraformTfvars);
  const extraTfvars = ['tectonic_admin_password_hash', 'tectonic_license_path', 'tectonic_pull_secret_path', 'tectonic_kube_apiserver_service_ip',
    'tectonic_kube_dns_service_ip', 'tectonic_kube_etcd_service_ip', 'tectonic_aws_etcd_ec2_type',
    'tectonic_aws_etcd_root_volume_size', 'tectonic_aws_etcd_root_volume_type', 'tectonic_etcd_count'];
  extraTfvars.forEach(key => {
    delete json[key];
  });
  return json;
};

const returnTerraformTfvars = (launchUrl, cookie, callback) => {
  getAssets(launchUrl,cookie, (err, res, terraformAssestsResponse) => {
    if (err !== null || res.statusCode !== 200 || res.headers['content-type'] !== 'application/zip' ) {
      return callback(new Error("Terraform get assets api call failed"),null);
    }
    getTerraformTfvars(terraformAssestsResponse, (terraformTfvars) => {
      const actualJson = returnRequiredTerraformTfvars(terraformTfvars);
      callback(null, actualJson);
    });
  });
};

const compareJson = (actualJson,expectedJson) => {
  let msg = '';
  const diff = deep(actualJson,expectedJson);
  if (typeof diff !== 'undefined'){
    diff.forEach(key => {
      msg = msg + "\n" + "TerraformTfvar:" + key.path + " ||"+" actualValue:" + key.lhs + " ||"
      +" expectedValue:"+ key.rhs;
    });
    msg = msg + "\nThe above Json attributes are not matching.";
  }
  return msg;
};
module.exports = { getAssets, getTerraformTfvars, returnRequiredTerraformTfvars, returnTerraformTfvars , compareJson };
