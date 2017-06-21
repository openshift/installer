const fs = require('fs');
const path = require('path');

const awsJsonPath = path.join(__dirname, '..', '..', '__tests__', 'examples', 'aws.json');
// eslint-disable-next-line no-sync
const awsTestDatajson = JSON.parse(fs.readFileSync(awsJsonPath, 'utf8'));
const awsProgressPath = path.join(__dirname, '..', '..', '__tests__', 'examples', 'tectonic-aws.progress');
// eslint-disable-next-line no-sync
const awsProgressData = JSON.parse(fs.readFileSync(awsProgressPath, 'utf8'));

/** Returns expected json. This json is used to prep the data required for the test */
const buildExpectedJson = () => {
  const tfvars = awsTestDatajson.variables;
  const clusterName = `aws-test-${new Date().getTime().toString()}`;
  tfvars.tectonic_cluster_name = clusterName;
  tfvars.tectonic_dns_name = clusterName;
  delete tfvars.tectonic_admin_password_hash;
  delete tfvars.tectonic_aws_extra_tags;
  delete tfvars.tectonic_aws_external_private_zone;

  return tfvars;
};

const clusterTagKey = () => {
  const json = buildExpectedJson();
  return json.tectonic_aws_extra_tags[0] + (new Date().getTime().toString());
};

const sshKeys = () => {
  const json = buildExpectedJson();
  json.tectonic_aws_ssh_key = 'tectonic-jenkins';
  return json.tectonic_aws_ssh_key;
};

const etcdOption = () => {
  return awsProgressData.clusterConfig.etcdOption;
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

module.exports = {
  buildExpectedJson,
  clusterTagKey,
  sshKeys,
  clusterSubdomainName,
  clusterBaseDomainName,
  adminEmail,
  adminPassword,
  etcdOption,
};
