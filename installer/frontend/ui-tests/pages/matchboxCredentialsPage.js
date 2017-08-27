const fs = require('fs');
const path = require('path');

const matchBoxCredentialsPageCommands = {
  test(json) {
    const parentDir = path.resolve(__dirname, '..');
    const caCertPath = path.join(parentDir, 'ca-cert.txt');
    const clientCertPath = path.join(parentDir, 'client-cert.txt');
    const clientKeyPath = path.join(parentDir, 'client-key.txt');

    /* eslint-disable no-sync */
    fs.writeFileSync(caCertPath, json.tectonic_metal_matchbox_ca);
    fs.writeFileSync(clientCertPath, json.tectonic_metal_matchbox_client_cert);
    fs.writeFileSync(clientKeyPath, json.tectonic_metal_matchbox_client_key);
    /* eslint-enable no-sync */
    return this
      .setValue('@caCertificate', caCertPath)
      .setValue('@clientCertificate', clientCertPath)
      .setValue('@clientKey', clientKeyPath);
  },
};

module.exports = {
  commands: [matchBoxCredentialsPageCommands],
  elements: {
    name: 'input#clusterName',
    caCertificate: {
      selector: '(//*[text()="Upload"]/input[@type="file"])[1]',
      locateStrategy: 'xpath',
    },
    clientCertificate: {
      selector: '(//*[text()="Upload"]/input[@type="file"])[2]',
      locateStrategy: 'xpath',
    },
    clientKey: {
      selector: '(//*[text()="Upload"]/input[@type="file"])[3]',
      locateStrategy: 'xpath',
    },
  },
};
