/* eslint-env jest, node */
/* eslint no-sync: "off" */

import fs from 'fs';
import path from 'path';

import { reducer } from '../reducer';
import { restoreActionTypes } from '../actions';
import { commitToServer } from '../server';
import '../components/aws-cluster-info';
import '../components/aws-define-nodes';
import '../components/aws-vpc';
import '../components/etcd';
import '../components/bm-sshkeys';
import '../components/bm-nodeforms';
import '../components/bm-hostname';
import '../components/cluster-type';

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
    description: 'works with baremetal',
    jsonPath: 'metal.json',
    progressPath: 'tectonic-baremetal.progress',
  },
  {
    description: 'works with aws',
    jsonPath: 'aws.json',
    progressPath: 'tectonic-aws.progress',
  },
  {
    description: 'works with aws (existing subnets)',
    jsonPath: 'aws-vpc.json',
    progressPath: 'tectonic-aws-vpc.progress',
  },
];

let dispatch;

beforeEach(() => {
  dispatch = jest.fn();
});


const readExample = example => {
  let json;
  try {
    json = JSON.parse(fs.readFileSync(path.resolve(__dirname, `examples/${example}`), 'utf8'));
  } catch (e) {
    console.warn(`${example} is not json`);
    throw e;
  }
  return json;
};


/* eslint-disable max-nested-callbacks */
describe('progress file example', () => {
  tests.forEach(t => {
    const example = readExample(t.jsonPath);
    const payload = readExample(t.progressPath);

    it(t.description, () => {
      const restored = reducer(initialState, {
        type: restoreActionTypes.RESTORE_STATE, payload,
      });

      global.fetch = jest.fn(() => Promise.resolve({
        ok: true,
        blob: () => Promise.resolve('success'),
        text: () => Promise.reject('failed'),
        json: () => Promise.resolve({}),
      }));

      // TODO: wait for this action to finish & then check expected state
      commitToServer(false, false, {salt: '$2a$12$96LR7NxL/T7LaijR0fxl3.'})(dispatch, () => restored);

      expect(fetch.mock.calls.length).toBe(1);
      const body = JSON.parse(fetch.mock.calls[0][1].body);

      expect(example).toEqual(body);
    });
  });
});
/* eslint-enable max-nested-callbacks */
