import _ from 'lodash';
import semver from 'semver';

// Compares two version strings, of the form \d+(.\d+)*
// Suitable for use with Array.sort()

const VERSION_MATCH = /^(\d+)(?:\.(\d+))*$/;

export const keyToAlg = (privateKey) => {
  const trimmed = privateKey.trim();
  const result = trimmed.match(/^-----BEGIN ([A-Z]{2,10}) PRIVATE KEY-----/);
  const keyType = result[1];
  if (keyType === 'EC') {
    return 'ECDSA';
  }
  return keyType;
};

export const compareVersions = (v1, v2) => {
  if (!v1.match(VERSION_MATCH)) {
    throw Error('compareVersions requires two version-formatted strings');
  }

  if (!v2.match(VERSION_MATCH)) {
    throw Error('compareVersions requires two version-formatted strings');
  }

  const v1Parts = v1.split('.').map(x => parseInt(x, 10));
  const v2Parts = v2.split('.').map(x => parseInt(x, 10));
  const minLen = Math.min(v1Parts.length, v2Parts.length);

  for (let i = 0; i < minLen; i++) {
    if (v1Parts[i] < v2Parts[i]) {
      return -1;
    }
    if (v1Parts[i] > v2Parts[i]) {
      return 1;
    }
  }

  if (v1Parts.length < v2Parts.length) {
    return -1;
  }

  if (v1Parts.length > v2Parts.length) {
    return 1;
  }

  return 0;
};

const toPath = base => {
  return field => `${base}.${_.isArray(field) ? field.join('.') : field}`;
};

export const toError = toPath('error');
export const toInFly = toPath('inFly');
export const toExtraData = toPath('extra');
export const toExtraDataError = toPath('extraError');
export const toExtraDataInFly = toPath('extraInFly');

export const isReleaseVersion = () => {
  return GIT_RELEASE_TAG === 'unknown' ? false : semver.valid(GIT_RELEASE_TAG) && !GIT_RELEASE_TAG.includes('-rc.');
};
