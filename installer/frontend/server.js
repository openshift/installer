import _ from 'lodash';

import {
  DRY_RUN,
  PLATFORM_TYPE,
  PULL_SECRET,
  RETRY,
  TECTONIC_LICENSE,
  getTectonicDomain,
  toAWS_TF,
  toBaremetal_TF,
} from './cluster-config';
import { clusterReadyActions, configActions, loadFactsActions, serverActions, FORMS } from './actions';
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
      dispatch(clusterReadyActions.notReady());
      return;
    }
    if (response.ok) {
      return response.json().then(json => dispatch(clusterReadyActions.status(json)));
    }
    return response.text().then(err => dispatch(clusterReadyActions.error(err)));
  }, payload => {
    if (payload instanceof TypeError) {
      payload = `${payload.message}. Is the installer running?`;
    }
    return dispatch(clusterReadyActions.error(payload));
  })
    .catch(err => console.error(err) || err);
};

const platformToFunc = {
  [AWS_TF]: toAWS_TF,
  [BARE_METAL_TF]: toBaremetal_TF,
};

let observeInterval;

// An action creator that builds a server message, calls fetch on that message, fires the appropriate actions
export const commitToServer = (dryRun = false, retry = false) => (dispatch, getState) => {
  setIn(DRY_RUN, dryRun, dispatch);
  setIn(RETRY, retry, dispatch);
  dispatch(serverActions.requested);

  const state = getState();

  const f = platformToFunc[state.clusterConfig[PLATFORM_TYPE]];
  if (!_.isFunction(f)) {
    throw Error(`unknown platform type "${state.clusterConfig[PLATFORM_TYPE]}"`);
  }

  const body = f(state, FORMS);
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
          return dispatch(serverActions.successful(payload));
        }) :
        response.text().then(payload => dispatch(serverActions.failed(payload)))
      , payload => dispatch(serverActions.failed(payload)))
    .catch(err => console.error(err));

  return dispatch(serverActions.sent);
};

// One-time fetch of initial data from server, followed by firing appropriate actions
// Guaranteed not to reject.
export const loadFacts = (dispatch) => {
  return fetchJSON('/tectonic/facts', {retries: 5})
    .then(
      facts => {
        addIn(TECTONIC_LICENSE, facts.license, dispatch);
        addIn(PULL_SECRET, facts.pullSecret, dispatch);
        dispatch(loadFactsActions.loaded({
          awsRegions: _.map(facts.amis, 'name'),
          buildTime: facts.buildTime,
          version: facts.tectonicVersion,
        }));
      },
      err => dispatch(loadFactsActions.error(err))
    )
    .catch(err => console.error(err));
};
