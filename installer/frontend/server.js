import _ from 'lodash';

import { getTectonicDomain, toAWS_TF, toBaremetal_TF, DRY_RUN, RETRY } from './cluster-config';
import { clusterReadyActionTypes, configActions, loadFactsActionTypes, serverActionTypes, FORMS } from './actions';
import { savable } from './reducer';
import { AWS_TF, BARE_METAL_TF } from './platforms';

const { setIn } = configActions;

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
  [AWS_TF]: {
    f: toAWS_TF,
  },
  [BARE_METAL_TF]: {
    f: toBaremetal_TF,
  },
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

  const obj = _.get(platformToFunc, request.platformType);
  if (!_.isFunction(obj.f)) {
    throw Error(`unknown platform type "${request.platformType}"`);
  }

  const body = obj.f(request, FORMS, opts);
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


// One-time fetch of AMIs from server, followed by firing appropriate actions
// Guaranteed not to reject.
const getAMIs = (dispatch) => {
  return fetchJSON('/containerlinux/images/amis', { retries: 5 })
    .then(m => {
      const awsRegions = m.map(({name}) => {
        return {label: name, value: name};
      });
      dispatch({
        type: loadFactsActionTypes.LOADED,
        payload: {awsRegions},
      });
    },
    err => {
      dispatch({
        type: loadFactsActionTypes.ERROR,
        payload: err,
      });
    }).catch(err => console.error(err));
};

// One-time fetch of facts from server. Abstracts getAMIs.
// Guaranteed not to reject.
export const loadFacts = (dispatch) => {
  if (_.includes(window.config.platforms, AWS_TF)) {
    return getAMIs(dispatch);
  }
  dispatch({type: loadFactsActionTypes.LOADED, payload: {}});
  return Promise.resolve();
};
