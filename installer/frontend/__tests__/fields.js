/* eslint-env jest, node */

// monkey patch node's console :-/
console.debug = console.debug || console.info;

import _ from 'lodash';
import { __deleteEverything__, configActions, configActionTypes } from '../actions';
import { Field, Form } from '../form';
import { store } from '../store';
import { DEFAULT_CLUSTER_CONFIG } from '../cluster-config';
import { toError } from '../utils';
// toIgnore, toAsyncError, toExtraData, toInFly, toExtraDataInFly, toExtraDataError

const getCC = (path, f) => _.get(store.getState().clusterConfig, f ? f(path) : path);

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
  store.dispatch({type: configActionTypes.SET, payload: DEFAULT_CLUSTER_CONFIG});

  expect(store.getState().clusterConfig[name]).toEqual('a');
  store.dispatch(configActions.updateField(name, 'b'));
  expect(store.getState().clusterConfig[name]).toEqual('b');
  done();
});

test('test field dependency validator is called', done => {
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

  store.dispatch({type: configActionTypes.SET, payload: DEFAULT_CLUSTER_CONFIG});
  store.dispatch(configActions.updateField(aName, 'b'));
});

test('tests form validator is called', done => {
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

  store.dispatch({type: configActionTypes.SET, payload: DEFAULT_CLUSTER_CONFIG});
  store.dispatch(configActions.updateField(fieldName, 'b'));
});

test('tests sync invalidation', done => {
  expect.assertions(3);

  const invalid = 'is invalid';
  const fieldName = 'aField';
  const field = new Field(fieldName, {
    default: 'a',
    validator: value => value === 'b' && invalid,
  });

  new Form('aForm', [field]);

  store.dispatch({type: configActionTypes.SET, payload: DEFAULT_CLUSTER_CONFIG});
  expect(getCC(fieldName, toError)).toEqual(undefined);

  store.dispatch(configActions.updateField(fieldName, 'b'));
  expect(getCC(fieldName, toError)).toEqual(invalid);

  store.dispatch(configActions.updateField(fieldName, 'a'));
  expect(getCC(fieldName, toError)).toEqual(undefined);

  done();
});
