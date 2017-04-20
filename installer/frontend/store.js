import { createStore, applyMiddleware } from 'redux';
import thunkMiddleware from 'redux-thunk';

import { reducer } from './reducer';

export const store = applyMiddleware(thunkMiddleware)(createStore)(reducer);
export const dispatch = (...args) => store.dispatch(...args);
