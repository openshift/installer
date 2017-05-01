/* eslint-env jest */
jest.unmock('../reducer');

import { reducer, savable } from '../reducer';
import { configActionTypes, restoreActionTypes } from '../actions';
import { SSH_AUTHORIZED_KEYS } from '../cluster-config';

const initialState = reducer(undefined, {type: 'Some Initial Action'});

describe('reducer', () => {
  it('saves and restores the initial state faithfully', () => {
    const saved = savable(initialState);
    const restoredState = reducer(initialState, {
      type: restoreActionTypes.RESTORE_STATE,
      payload: saved,
    });
    const initialJSON = JSON.stringify(initialState);
    const restoredJSON = JSON.stringify(restoredState);
    expect(initialJSON).toEqual(restoredJSON);
  });

  it('allows creation and deletion of ssh keys', () => {
    const moreKeys = reducer(initialState, {
      type: configActionTypes.SET_SSH_AUTHORIZED_KEYS,
      payload: {
        index: 0,
        value: 'KEY ZERO',
      },
    });

    const lessKeys = reducer(moreKeys, {
      type: configActionTypes.REMOVE_SSH_AUTHORIZED_KEYS,
      payload: 0,
    });

    expect(moreKeys.clusterConfig[SSH_AUTHORIZED_KEYS]).toEqual(['KEY ZERO']);
    expect(lessKeys.clusterConfig[SSH_AUTHORIZED_KEYS]).toEqual([]);
  });
});
