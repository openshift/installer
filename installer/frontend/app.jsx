// Polyfills, imported for side effects
import 'babel-polyfill'; // Required for `new Promise()`
import 'whatwg-fetch'; // Required for `fetch()`

import _ from 'lodash';
import React from 'react';
import ReactDom from 'react-dom';
import { Provider } from 'react-redux';
import { browserHistory, Router } from 'react-router';
import Cookie from 'js-cookie';

import { navActionTypes, restoreActionTypes, validateAllFields } from './actions';
import { trail, getAllRoutes } from './trail';
import { TectonicGA } from './tectonic-ga';
import { savable } from './reducer';
import { loadFacts, observeClusterStatus } from './server';
import { store, dispatch } from './store';
import { Base } from './components/base';
import { clusterReadyActionTypes } from './actions';

const saveState = () => {
  const state = store.getState();
  const data = savable(state);
  sessionStorage.setItem('state', JSON.stringify(data));
};

window.reset = () => {
  Cookie.remove('tectonic-installer');
  window.removeEventListener('beforeunload', saveState);
  sessionStorage.clear();
  fetch('/cluster/done', {method: 'POST', credentials: 'same-origin'})
    .catch(() => undefined) // We don't really care if this completes - we're done here!
    .then(() => window.location = '/');
};

export const navigateNext = () => {
  const state = store.getState();
  const t = trail(state);
  const currentPage = t.pageByPath.get(state.path);
  const nextPage = t.nextFrom(currentPage);
  browserHistory.push(nextPage.path);
};

const fixLocation = () => {
  const state = store.getState();
  const t = trail(state);
  const fixed = t.fixPath(state.path, state);
  if (fixed !== state.path) {
    browserHistory.push(fixed);
  }
};

store.subscribe(_.debounce(saveState, 5000));
window.addEventListener('beforeunload', saveState);

// Stuff we need to load before we can run anything
loadFacts(dispatch);

try {
  const state = JSON.parse(sessionStorage.getItem('state'));
  dispatch({
    type: restoreActionTypes.RESTORE_STATE,
    payload: state,
  });
  console.debug('Restored state from sessionStorage.');
} catch (e) {
  console.error(`Error restoring state from sessionStorage: ${e.message || e.toString()}`);
}

store.dispatch(validateAllFields(() => {
  TectonicGA.initialize();

  try {
    observeClusterStatus(dispatch, store.getState)
    .then(res => {
      if (res && res.type === clusterReadyActionTypes.STATUS) {
        setInterval(() => observeClusterStatus(dispatch, store.getState), 10000);
      }
      fixLocation();
    });
  } catch (e) {
    console.error(`Error restoring state from sessionStorage: ${e.message || e.toString()}`);
  }

  // Because route.onEnter doesn't monitor the state, there is a race condition where
  //   - onEnter runs on URL, decides that URL is in trail(state(time0))
  //   - state updates to state(time1) based on server events, etc
  //   - base renders URL, which is not in trail(state(time1))
  //   - mass hysteria
  //
  // As a result, we shuffle the state before the router sees it.
  browserHistory.listen(location => {
    store.dispatch({
      type: navActionTypes.LOCATION_CHANGE,
      payload: location,
    });
    TectonicGA.sendPageView(location.pathname);
  });

  store.subscribe(fixLocation);

  // Set up routing

  const routes = {
    path: '/',
    component: Base,
    childRoutes: getAllRoutes(),
  };

  ReactDom.render(
    <Provider store={store}>
      <Router history={browserHistory} routes={routes} />
    </Provider>,
    document.getElementById('application')
  );
}));

window.onerror = (message, source, lineno, colno, optError={}) => {
  try {
    const e = `${message} ${source} ${lineno} ${colno}`;
    TectonicGA.sendError(e, optError.stack);
  } catch(err) {
    try {
      // eslint-disable-next-line no-console
      console.error(err);
    } catch (ignored) {
      // ignore
    }
  }
};
