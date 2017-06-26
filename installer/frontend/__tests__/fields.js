/* eslint-env jest, node */

// monkey patch node's console :-/
console.debug = console.debug || console.info;

import { __deleteEverything__, configActions, configActionTypes } from '../actions';
import { Field, Form } from '../form';
import { store } from '../store';
import { DEFAULT_CLUSTER_CONFIG } from '../cluster-config';


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

test('tests two Fields', done => {
  expect.assertions(1);

  const aName = 'aField';

  const aField = new Field(aName, {default: 'a'});
  const bField = new Field('bField', {
    default: 'c',
    dependencies: [aName],
    validator: (value, cc) => {
      expect(cc.aField).toEqual('b');
      done();
    },
  });

  new Form('aForm', [aField, bField]);

  store.dispatch({type: configActionTypes.SET, payload: DEFAULT_CLUSTER_CONFIG});
  store.dispatch(configActions.updateField(aName, 'b'));
});
