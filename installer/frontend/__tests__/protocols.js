/* eslint-env jest, node */
/* eslint no-sync: "off" */

import fs from 'fs';
import path from 'path';

import { reducer } from '../reducer';
import { restoreActions } from '../actions';
import { commitToServer } from '../server';
import '../components/aws-cloud-credentials';
import '../components/aws-cluster-info';
import '../components/aws-submit-keys';
import '../components/aws-vpc';
import '../components/bm-credentials';
import '../components/bm-hostname';
import '../components/bm-matchbox';
import '../components/bm-nodeforms';
import '../components/bm-sshkeys';
import '../components/certificate-authority';
import '../components/cluster-type';
import '../components/nodes';
import '../components/users';

const structureOnly = (obj) => {
  const toString = Object.prototype.toString;

  if (toString.call(obj) === '[object Object]') {
    const ret = {};
    Object.keys(obj).forEach(k => {
      ret[k] = structureOnly(obj[k]);
    });
    return ret;
  }

  if (Array.isArray(obj)) {
    return [];
  }

  return '';
};

const initialState = reducer(undefined, {type: 'Some Initial Action'});

const tests = [
  {
    description: 'works with bare metal',
    expected: 'metal-out.json',
    state: 'metal-in.json',
  },
  {
    description: 'works with AWS',
    expected: 'aws-custom-vpc-out.json',
    state: 'aws-custom-vpc-in.json',
  },
  {
    description: 'works with AWS (existing VPC)',
    expected: 'aws-existing-vpc-out.json',
    state: 'aws-existing-vpc-in.json',
  },
];

let dispatch;

beforeEach(() => {
  dispatch = jest.fn();
});

const readExample = filename => {
  let json;
  try {
    json = JSON.parse(fs.readFileSync(path.resolve(__dirname, `examples/${filename}`), 'utf8'));
  } catch (e) {
    console.warn(`${filename} is not json`);
    throw e;
  }
  return json;
};

/* eslint-disable max-nested-callbacks */
describe('progress file example', () => {
  tests.forEach(t => {
    it(t.description, () => {
      const restored = reducer(initialState, restoreActions.restore(readExample(t.state)));

      global.fetch = jest.fn(() => Promise.resolve({
        ok: true,
        blob: () => Promise.resolve('success'),
        text: () => Promise.reject('failed'),
        json: () => Promise.resolve({}),
      }));

      // TODO: wait for this action to finish & then check expected state
      commitToServer(false, false)(dispatch, () => restored);

      expect(fetch.mock.calls.length).toBe(1);
      const body = JSON.parse(fetch.mock.calls[0][1].body);

      expect(readExample(t.expected)).toEqual(body);
    });
  });
});
/* eslint-enable max-nested-callbacks */
