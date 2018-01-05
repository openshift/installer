/* eslint-env jest */
jest.unmock('../reducer');

import { reducer, savable } from '../reducer';
import { restoreActionTypes } from '../actions';

const initialState = reducer(undefined, {type: 'Some Initial Action'});

describe('reducer', () => {
  it('saves and restores the initial state faithfully', () => {
    const saved = savable(initialState);
    const restoredState = reducer(initialState, {
      type: restoreActionTypes.RESTORE_STATE,
      payload: saved,
    });
    expect(initialState).toEqual(restoredState);
  });
});
