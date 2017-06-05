const fs = require('fs');
const path = require('path');

const readAwsTestDataJson = () => {
  const fileName = path.join(__dirname, '../','../', '__tests__', 'examples', 'aws.json');
    //eslint-disable-next-line no-sync
  const awsTestDatajson = JSON.parse(fs.readFileSync(fileName, 'utf8'));
  return awsTestDatajson.variables;
};

const cluster = name => {
  return name + new Date().getTime().toString();
};

/** Returns expected json. This json is used to prep the data required for the test */

const buildExpectedJson = (opt) => {
  const options = {} || opt;
  const json = readAwsTestDataJson();

  json.tectonic_cluster_name = cluster("aws-test-");
  json.tectonic_dns_name = cluster("aws-test-");
  delete json.tectonic_admin_password_hash;
  delete json.tectonic_aws_extra_tags;
  delete json.tectonic_aws_external_private_zone;

  Object.keys(options).map(key => {
    json[key] = options[key];
  });
  return json;
};

const clusterTagKey = () => {
  const json = buildExpectedJson();
  return json.tectonic_aws_extra_tags[0] + (new Date().getTime().toString());
};
const sshKeys = () => {
  const json = buildExpectedJson();
  json.tectonic_aws_ssh_key = 'jenkins';
  return json.tectonic_aws_ssh_key;
};
const clusterSubdomainName = () => {
  const json = buildExpectedJson();
  return json.tectonic_dns_name;
};
const clusterBaseDomainName = () => {
  const json = buildExpectedJson();
  return json.tectonic_base_domain;
};
const adminEmail = () => {
  const json = buildExpectedJson();
  return json.tectonic_admin_email;
};
const adminPassword = () => {
  const json = buildExpectedJson();
  json.tectonic_admin_password_hash = 'password';
  return json.tectonic_admin_password_hash;
};
const createCoreosCredentials = (Path, env_var) => {
  //eslint-disable-next-line no-sync
  fs.writeFileSync(Path, env_var);
};

module.exports = {readAwsTestDataJson,buildExpectedJson,clusterTagKey,sshKeys,clusterSubdomainName,clusterBaseDomainName,adminEmail,adminPassword,createCoreosCredentials };
