/* eslint-env jest */
jest.unmock('../reducer');

import { reducer, savable } from '../reducer';
import { restoreActions } from '../actions';

const initialState = reducer(undefined, {type: 'Some Initial Action'});

describe('reducer', () => {
  it('saves and restores the initial state faithfully', () => {
    const saved = savable(initialState);
    const restoredState = reducer(initialState, restoreActions.restore(saved));
    expect(initialState).toEqual(restoredState);
  });
});
