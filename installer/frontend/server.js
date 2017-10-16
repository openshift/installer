import _ from 'lodash';

import { getTectonicDomain, toAWS_TF, toBaremetal_TF, DRY_RUN, PULL_SECRET, RETRY, TECTONIC_LICENSE } from './cluster-config';
import { clusterReadyActionTypes, configActions, loadFactsActionTypes, serverActionTypes, FORMS } from './actions';
import { savable } from './reducer';
import { AWS_TF, BARE_METAL_TF } from './platforms';

const { addIn, setIn } = configActions;

// Either return parsable JSON, or fail (and assume returned text is an error message)
const fetchJSON = (url, opts, ...args) => {
  opts = opts || {};
  opts.credentials = 'same-origin';
  return fetch(url, opts, ...args).then(response => {
    if (response.ok) {
      return response.json();
    }

    if (opts.retries > 0) {
      opts.retries--;
      return fetchJSON(url, opts, ...args);
    }
    return response.text().then(Promise.reject);
  });
};

// Poll server for cluster status.
// Guaranteed not to reject
const {NOT_READY, STATUS, ERROR} = clusterReadyActionTypes;

export const observeClusterStatus = (dispatch, getState) => {
  const cc = getState().clusterConfig;
  const tectonicDomain = getTectonicDomain(cc);
  const opts = {
    credentials: 'same-origin',
    body: JSON.stringify({tectonicDomain}),
    method: 'POST',
    headers: {
      'Accept': 'application/json',
      'Content-Type': 'application/json',
    },
  };

  return fetch('/tectonic/status', opts).then(response => {
    if (response.status === 404) {
      dispatch({type: NOT_READY});
      return;
    }
    if (response.ok) {
      return response.json().then(payload => dispatch({type: STATUS, payload}));
    }
    return response.text().then(payload => dispatch({type: ERROR, payload}));
  }, payload => {
    if (payload instanceof TypeError) {
      payload = `${payload.message}. Is the installer running?`;
    }
    return dispatch({type: ERROR, payload});
  })
    .catch(err => console.error(err) || err);
};

const platformToFunc = {
  [AWS_TF]: toAWS_TF,
  [BARE_METAL_TF]: toBaremetal_TF,
};

let observeInterval;

// An action creator that builds a server message, calls fetch on that message, fires the appropriate actions
export const commitToServer = (dryRun = false, retry = false, opts = {}) => (dispatch, getState) => {
  setIn(DRY_RUN, dryRun, dispatch);
  setIn(RETRY, retry, dispatch);

  const {COMMIT_REQUESTED, COMMIT_FAILED, COMMIT_SUCCESSFUL, COMMIT_SENT} = serverActionTypes;

  dispatch({type: COMMIT_REQUESTED});

  const state = getState();
  const request = Object.assign({}, state.clusterConfig, {progress: savable(state)});

  const f = platformToFunc[request.platformType];
  if (!_.isFunction(f)) {
    throw Error(`unknown platform type "${request.platformType}"`);
  }

  const body = f(request, FORMS, opts);
  fetch('/terraform/apply', {
    credentials: 'same-origin',
    method: 'POST',
    body: JSON.stringify(body),
  })
    .then(
      response => response.ok ?
        response.blob().then(payload => {
          observeClusterStatus(dispatch, getState);
          if (!observeInterval) {
            observeInterval = setInterval(() => observeClusterStatus(dispatch, getState), 10000);
          }
          return dispatch({payload, type: COMMIT_SUCCESSFUL});
        }) :
        response.text().then(payload => dispatch({payload, type: COMMIT_FAILED}))
      , payload => dispatch({payload, type: COMMIT_FAILED}))
    .catch(err => console.error(err));

  return dispatch({
    type: COMMIT_SENT,
    payload: body,
  });
};

// One-time fetch of initial data from server, followed by firing appropriate actions
// Guaranteed not to reject.
export const loadFacts = (dispatch) => {
  return fetchJSON('/tectonic/facts', {retries: 5})
    .then(m => {
      addIn(TECTONIC_LICENSE, m.license, dispatch);
      addIn(PULL_SECRET, m.pullSecret, dispatch);
      dispatch({
        type: loadFactsActionTypes.LOADED,
        payload: {awsRegions: _.map(m.amis, 'name')},
      });
    },
    err => {
      dispatch({
        type: loadFactsActionTypes.ERROR,
        payload: err,
      });
    }).catch(err => console.error(err));
};
