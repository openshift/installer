/* eslint-env jest, node */

// monkey patch node's console :-/
console.debug = console.debug || console.info;

import _ from 'lodash';
import { __deleteEverything__, configActions, configActionTypes } from '../actions';
import { Field, Form } from '../form';
import { store } from '../store';
import { DEFAULT_CLUSTER_CONFIG } from '../cluster-config';
import { toError, toAsyncError } from '../utils';
// TODO: (kans) test these things
// toIgnore, , toExtraData, toInFly, toExtraDataInFly, toExtraDataError

const expectCC = (path, expected, f) => {
  const value = _.get(store.getState().clusterConfig, f ? f(path) : path);
  expect(value).toEqual(expected);
};

const resetCC = () => store.dispatch({
  type: configActionTypes.SET, payload: DEFAULT_CLUSTER_CONFIG,
});

const updateField = (field, value) => store.dispatch(configActions.updateField(field, value));


beforeEach(() => store.dispatch(__deleteEverything__()));

test('updates a Field', () => {
  expect.assertions(4);

  const name = 'aField';
  const aField = new Field(name, {
    default: 'a',
    validator: (value, cc) => {
      expect(value).toEqual('b');
      expect(cc[name]).toEqual('b');
    },
  });

  new Form('aForm', [aField]);
  resetCC();

  expectCC(name, 'a');
  updateField(name, 'b');
  expectCC(name, 'b');
});

test('field dependency validator is called', done => {
  expect.assertions(1);

  const aName = 'aField';

  const aField = new Field(aName, {default: 'a'});
  const bField = new Field('bField', {
    default: 'c',
    dependencies: [aName],
    validator: (value, cc) => {
      expect(cc[aName]).toEqual('b');
      done();
    },
  });

  new Form('aForm', [aField, bField]);

  resetCC();
  updateField(aName, 'b');
});

test('form validator is called', done => {
  expect.assertions(3);

  const fieldName = 'aField';

  new Form('aForm', [new Field(fieldName, {default: 'a'})], {
    validator: (value, cc, oldValue) => {
      expectCC(fieldName, 'b');
      expect(value[fieldName]).toEqual('b');
      expect(oldValue[fieldName]).toEqual('a');
      done();
    },
  });

  resetCC();
  updateField(fieldName, 'b');
});

test('sync invalidation', () => {
  expect.assertions(3);

  const invalid = 'is invalid';
  const fieldName = 'aField';
  const field = new Field(fieldName, {
    default: 'a',
    validator: value => value === 'b' && invalid,
  });

  new Form('aForm', [field]);

  resetCC();
  expectCC(fieldName, undefined, toError);

  updateField(fieldName, 'b');
  expectCC(fieldName, invalid, toError);

  updateField(fieldName, 'a');
  expectCC(fieldName, undefined, toError);
});


test('async invalidation', async done => {
  expect.assertions(5);

  const invalid = 'is invalid';
  const fieldName = 'aField';

  const field = new Field(fieldName, {
    default: 'a',
    asyncValidator: (dispatch, getState, value, oldValue, isNow) => {
      expect(isNow()).toEqual(true);
      return value === 'b' && invalid;
    },
  });

  new Form('aForm', [field]);

  expectCC(fieldName, undefined, toAsyncError);

  await updateField(fieldName, 'b');
  expectCC(fieldName, invalid, toAsyncError);

  await updateField(fieldName, 'a');
  expectCC(fieldName, undefined, toAsyncError);

  done();
});
