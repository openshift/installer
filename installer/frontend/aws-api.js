const postJSON = url => (body, creds) => {
  const opts = {
    credentials: 'same-origin',
    method: 'POST',
    headers: {
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      'Tectonic-AccessKeyID': creds.AccessKeyID,
      'Tectonic-SecretAccessKey': creds.SecretAccessKey,
      'Tectonic-Region': creds.Region,
    },
  };
  // Edge's fetch() doesn't allow empty headers
  if (creds.SessionToken) {
    opts.headers['Tectonic-SessionToken'] = creds.SessionToken;
  }
  // Don't stringify null, false, undefined, etc.
  if (body) {
    opts.body = JSON.stringify(body);
  }
  return fetch(url, opts).then(response => response.ok ?
    response.json() :
    response.text().then(t => Promise.reject(t))
  );
};

const TFPostJSON = url => (body = {}, creds = {}, platform = 'aws') => {
  const opts = {
    credentials: 'same-origin',
    method: 'POST',
  };
  body.platform = platform;
  if (platform === 'aws') {
    body.credentials = {
      AWSAccessKeyID: creds.AccessKeyID,
      AWSSecretAccessKey: creds.SecretAccessKey,
      AWSSessionToken: creds.SessionToken,
      AWSRegion: creds.Region,
    };
  }
  opts.body = JSON.stringify(body);
  return fetch(url, opts).then(response => response.ok ?
    response.text() :
    response.text().then(t => Promise.reject(t))
  );
};

export const getRegions = postJSON('/aws/regions');
export const getSsh = postJSON('/aws/ssh-key-pairs');
export const getIamRoles = postJSON('/aws/iam-roles');
export const getVpcs = postJSON('/aws/vpcs');
export const getVpcSubnets = postJSON('/aws/vpcs/subnets');
export const getZones = postJSON('/aws/zones');
export const getDomainInfo = postJSON('/aws/domain');
export const getDefaultSubnets = postJSON('/aws/default-subnets');
export const validateSubnets = postJSON('/aws/subnets/validate');
export const TFDestroy = TFPostJSON('/terraform/destroy');
