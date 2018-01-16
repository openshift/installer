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

import { clusterReadyActionTypes, FIELDS, restoreActionTypes, validateFields } from './actions';
import { PLATFORM_TYPE } from './cluster-config';
import { trail } from './trail';
import { TectonicGA } from './tectonic-ga';
import { savable } from './reducer';
import { loadFacts, observeClusterStatus } from './server';
import { store, dispatch } from './store';
import { Base } from './components/base';

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

store.subscribe(_.debounce(saveState, 5000));
window.addEventListener('beforeunload', saveState);

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

history.listen(({pathname, state}) => {
  // Process next step / previous step navigation trigger if present
  if (state && (state.next || state.previous)) {
    const storeState = store.getState();
    const t = trail(storeState);
    const currentPage = t.pageByPath.get(history.location.pathname);
    const nextPage = state.next ? t.nextFrom(currentPage) : t.previousFrom(currentPage);
    if (state.next) {
      TectonicGA.sendEvent('Page Navigation Next', 'click', 'next on', storeState.clusterConfig[PLATFORM_TYPE]);
    }
    history.replace(_.get(nextPage, 'path'));
    return;
  }
  TectonicGA.sendPageView(pathname);
});

ReactDom.render(
  <Provider store={store}>
    <Router history={history}>
      <Base />
    </Router>
  </Provider>,
  document.getElementById('application')
);

// Stuff we need to load before we can run anything
loadFacts(dispatch)
  .then(() => validateFields(_.keys(FIELDS), store.getState, dispatch))
  .then(() => {
    TectonicGA.initialize();
    return observeClusterStatus(dispatch, store.getState);
  })
  .then(res => {
    if (res && res.type === clusterReadyActionTypes.STATUS) {
      setInterval(() => observeClusterStatus(dispatch, store.getState), 10000);
    }
  });

window.onerror = (message, source, lineno, colno, optError = {}) => {
  try {
    const e = `${message} ${source} ${lineno} ${colno}`;
    TectonicGA.sendError(e, optError.stack);
  } catch (err) {
    try {
      // eslint-disable-next-line no-console
      console.error(err);
    } catch (ignored) {
      // ignore
    }
  }
};
