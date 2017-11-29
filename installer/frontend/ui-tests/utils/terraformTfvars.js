const deep = require('deep-diff').diff;
const fs = require('fs');
const JSZip = require('jszip');
const path = require('path');
const request = require('request');

const diffTfvars = (client, assetsZip, expected) => {
  JSZip.loadAsync(assetsZip).then(zip => {
    zip.file(/tfvars$/)[0].async('string').then(tfvars => {
      const diff = deep(JSON.parse(tfvars), expected);
      if (diff !== undefined) {
        client.assert.fail(
          'The following terraform.tfvars attributes differ from their expected value: ' +
          diff.map(d => `\n  ${d.path.join('.')} (expected: ${d.rhs}, got: ${d.lhs})`)
        );
      }
    });
  });
};

const testManualBoot = (client, expectedFilename) => {
  client.page.submitPage()
    .click('@manuallyBoot')
    .expect.element('a[href="/terraform/assets"]').to.be.visible.before(120000);

  client.getCookie('tectonic-installer', ({value}) => {
    const options = {
      url: `${client.launch_url}/terraform/assets`,
      method: 'GET',
      encoding: null,
      headers: {'Cookie': `tectonic-installer=${value}`},
    };
    request(options, (err, res, assetsZip) => {
      if (err) {
        return client.assert.fail(err);
      }
      if (res.statusCode !== 200 || res.headers['content-type'] !== 'application/zip') {
        return client.assert.fail('Terraform get assets API call failed');
      }

      // eslint-disable-next-line no-sync
      const expected = JSON.parse(fs.readFileSync(path.join(__dirname, '..', 'output', expectedFilename), 'utf8'));

      diffTfvars(client, assetsZip, expected);
    });
  });
};

const jsonDir = path.join(__dirname, '..', '..', '__tests__', 'examples');
// eslint-disable-next-line no-sync
const loadJson = filename => JSON.parse(fs.readFileSync(path.join(jsonDir, filename), 'utf8'));

module.exports = {
  loadJson,
  testManualBoot,
};
