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

const getCC = (path, f) => _.get(store.getState().clusterConfig, f ? f(path) : path);
const expectCC = (path, f, value) => expect(getCC(path, f)).toEqual(value);
const resetCC = () => store.dispatch({
  type: configActionTypes.SET, payload: DEFAULT_CLUSTER_CONFIG,
});

beforeEach(() => store.dispatch(__deleteEverything__()));

test('updates a Field', done => {
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

  expect(store.getState().clusterConfig[name]).toEqual('a');
  store.dispatch(configActions.updateField(name, 'b'));
  expect(store.getState().clusterConfig[name]).toEqual('b');
  done();
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
  store.dispatch(configActions.updateField(aName, 'b'));
});

test('form validator is called', done => {
  expect.assertions(3);

  const fieldName = 'aField';

  new Form('aForm', [new Field(fieldName, {default: 'a'})], {
    validator: (value, cc, oldValue) => {
      expect(value[fieldName]).toEqual('b');
      expect(cc[fieldName]).toEqual('b');
      expect(oldValue[fieldName]).toEqual('a');
      done();
    },
  });

  resetCC();
  store.dispatch(configActions.updateField(fieldName, 'b'));
});

test('sync invalidation', done => {
  expect.assertions(3);

  const invalid = 'is invalid';
  const fieldName = 'aField';
  const field = new Field(fieldName, {
    default: 'a',
    validator: value => value === 'b' && invalid,
  });

  new Form('aForm', [field]);

  resetCC();
  expect(getCC(fieldName, toError)).toEqual(undefined);

  store.dispatch(configActions.updateField(fieldName, 'b'));
  expect(getCC(fieldName, toError)).toEqual(invalid);

  store.dispatch(configActions.updateField(fieldName, 'a'));
  expect(getCC(fieldName, toError)).toEqual(undefined);

  done();
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

  expectCC(fieldName, toAsyncError, undefined);

  await store.dispatch(configActions.updateField(fieldName, 'b'));
  expectCC(fieldName, toAsyncError, invalid);

  await store.dispatch(configActions.updateField(fieldName, 'a'));
  expectCC(fieldName, toAsyncError, undefined);

  done();
});
