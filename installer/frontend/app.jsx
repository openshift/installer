// Polyfills, imported for side effects
import 'babel-polyfill'; // Required for `new Promise()`
import 'whatwg-fetch'; // Required for `fetch()`

import _ from 'lodash';
import React from 'react';
import ReactDom from 'react-dom';
import { Provider } from 'react-redux';
import { Router } from 'react-router-dom';
import createHistory from 'history/createBrowserHistory';
import Cookie from 'js-cookie';

import { navChange, restoreActionTypes, validateAllFields } from './actions';
import { trail } from './trail';
import { TectonicGA } from './tectonic-ga';
import { savable } from './reducer';
import { loadFacts, observeClusterStatus } from './server';
import { store, dispatch } from './store';
import { Base } from './components/base';
import { clusterReadyActionTypes } from './actions';

const history = createHistory();

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

const setLocation = path => store.dispatch(navChange(path));

const fixLocation = () => {
  const state = store.getState();
  const t = trail(state);
  const fixed = t.fixPath(state.path, state);
  if (fixed !== state.path) {
    setLocation(fixed);
  }
  if (fixed !== window.location.pathname) {
    history.push(fixed);
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
  history.listen(({pathname, state}) => {
    // Process next step / previous step navigation trigger if present
    if (state && (state.next || state.previous)) {
      const storeState = store.getState();
      const t = trail(storeState);
      const currentPage = t.pageByPath.get(storeState.path);
      const nextPage = state.next ? t.nextFrom(currentPage) : t.previousFrom(currentPage);
      setLocation(_.get(nextPage, 'path'));
      return;
    }

    setLocation(pathname);
    TectonicGA.sendPageView(pathname);
  });

  store.subscribe(fixLocation);

  ReactDom.render(
    <Provider store={store}>
      <Router history={history}>
        <Base />
      </Router>
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
