import _ from 'lodash';
import { combineReducers } from 'redux';
import { fromJS } from 'immutable';

import {
  DEFAULT_CLUSTER_CONFIG,
} from './cluster-config';

import {
  awsActionTypes,
  clusterReadyActionTypes,
  configActionTypes,
  dirtyActionTypes,
  eventErrorsActionTypes,
  loadFactsActionTypes,
  navActionTypes,
  restoreActionTypes,
  serverActionTypes,
  sequenceActionTypes,
  commitPhases,
} from './actions';

const UNLOADED_AWS_VALUE = {
  inFly: false,
  value: [],
  error: null,
};

const DEFAULT_AWS = {
  availableVpcs: UNLOADED_AWS_VALUE,
  availableVpcSubnets: UNLOADED_AWS_VALUE,
  availableSsh: UNLOADED_AWS_VALUE,
  availableRegions: UNLOADED_AWS_VALUE,
  availableKms: UNLOADED_AWS_VALUE,
  availableR53Zones: UNLOADED_AWS_VALUE,
  createdKms: UNLOADED_AWS_VALUE,
  domainInfo: UNLOADED_AWS_VALUE,
  subnets: UNLOADED_AWS_VALUE,
  validateSubnets: UNLOADED_AWS_VALUE,
  destroy: UNLOADED_AWS_VALUE,
};

// setIn({...}, 'a.b.c', 'd')
function setIn(object, path, value) {
  const array = _.isString(path) ? path.split('.') : path;
  return fromJS(object).setIn(array, value).toJS();
}

const reducersTogether = combineReducers({

  // State machine associated with server submissions
  // Should not be saved or restored
  commitState: (state, action) => {
    if (!state) {
      return {
        phase: commitPhases.IDLE,
        response: null,
      };
    }

    switch(action.type) {
    case serverActionTypes.COMMIT_REQUESTED:
      return {
        phase: commitPhases.REQUESTED,
      };
    case serverActionTypes.COMMIT_SENT:
      return {
        phase: commitPhases.WAITING,
      };
    case serverActionTypes.COMMIT_SUCCESSFUL:
      console.log('COMMIT SUCCESSFUL: provisioner has been configured');
      console.dir(action.payload);
      return {
        phase: commitPhases.SUCCEEDED,
        response: action.payload,
      };
    case serverActionTypes.COMMIT_FAILED:
      console.log('COMMIT FAILED');
      console.dir(action.payload);
      return {
        phase: commitPhases.FAILED,
        response: action.payload,
      };
    case serverActionTypes.COMMIT_RESET:
      if (state.phase !== commitPhases.SUCCEEDED &&
          state.phase !== commitPhases.FAILED &&
          state.phase !== commitPhases.IDLE) {
        throw Error('attempt to reset a working server connection');
      }

      return {
        phase: commitPhases.IDLE,
      };
    default:
      return state;
    }
  },

  aws: (state, action) => {
    if (!state) {
      return DEFAULT_AWS;
    }

    switch (action.type) {
    case awsActionTypes.SET:
      Object.keys(action.payload).forEach(k => {
        if (!DEFAULT_AWS.hasOwnProperty(k)) {
          throw Error(`attempt to set aws property ${k} missing from defaults`);
        }
      });
      return Object.assign({}, state, action.payload);
    default:
      return state;
    }
  },

  // The user's intentions for their cluster. Should be
  // saveable/restorable
  clusterConfig: (state, action) => {
    if (!state) {
      return DEFAULT_CLUSTER_CONFIG;
    }

    switch (action.type) {
    case configActionTypes.RESET:
      return {};
    case configActionTypes.SET:
      Object.keys(action.payload).forEach(k => {
        if (!DEFAULT_CLUSTER_CONFIG.hasOwnProperty(k)) {
          throw Error(`attempt to set cluster property ${k} missing from defaults`);
        }
      });
      return Object.assign({}, state, action.payload);

    case configActionTypes.BATCH_SET_IN: {
      let object = fromJS(state);
      action.payload.forEach(batch => {
        const [path, value] = batch;
        const array = _.isString(path) ? path.split('.') : path;
        object = object.setIn(array, value);
      });
      return object.toJS();
    }
    case configActionTypes.SET_IN:
      return setIn(state, action.payload.path, action.payload.value);

    case configActionTypes.MERGE:
      return fromJS(state).mergeDeep(action.payload).toJS();

    case configActionTypes.APPEND: {
      const length = _.get(state, action.payload.path).length;
      const path = `${action.payload.path}.${length}`;
      return setIn(state, path, action.payload.value);
    }

    case configActionTypes.REMOVE_AT: {
      const {path, index} = action.payload;
      const array = _.isString(path) ? path.split('.') : path;
      array.push(index.toString());
      // TODO: (kans) delete all the other stuff too
      const invalidArray = ['error'].concat(array);
      const asyncArray = ['error_async'].concat(array);
      const arrays = [array, asyncArray, invalidArray];
      return fromJS(state).withMutations(map => {
        arrays.forEach(a => {
          if (map.getIn(a)) {
            map = map.deleteIn(a);
          }
          return map;
        });
      }).toJS();
    }
    default:
      return state;
    }
  },

  // Errors resulting from server states, user uploads, network issues,
  // or other transient phenomena. Should not be saved or restored.
  eventErrors: (state, action) => {
    if (!state) {
      return {};
    }

    switch (action.type) {
    case eventErrorsActionTypes.ERROR:
      return Object.assign({}, state, {
        [action.payload.name]: action.payload.error,
      });
    default:
      return state;
    }
  },

  // Facts the server knows at load time, that we have to get asynchronously.
  // Should not change value across restores. (so never save or restore it)
  serverFacts: (state, action) => {
    if (state === undefined) {
      return {
        loaded: false,
        error: null,
        awsRegions: null,
      };
    }

    switch (action.type) {
    case loadFactsActionTypes.LOADED:
      return {
        loaded: true,
        error: null,
        awsRegions: action.payload.awsRegions,
      };
    case loadFactsActionTypes.ERROR:
      return {
        loaded: true,
        error: action.payload,
        awsRegions: null,
      };
    default:
      return state;
    }
  },

  // The current location as reported by react-router
  // Should not be saved or restored.
  path: (state, action) => {
    if (state === undefined) {
      return window.location.pathname;
    }

    switch (action.type) {
    case navActionTypes.LOCATION_CHANGE:
      return action.payload.pathname;
    default:
      return state;
    }
  },

  // The status of the cluster. Should be preserved across restores.
  cluster: (state, action) => {
    if (!state) {
      return {
        loaded: false,
        ready: false,
        error: null,
        status: null,
      };
    }

    switch (action.type) {
    case clusterReadyActionTypes.NOT_READY:
      return Object.assign({}, state, {
        loaded: true,
        ready: false,
        status: null,
      });
    case clusterReadyActionTypes.STATUS:
      return Object.assign({}, state, {
        loaded: true,
        ready: true,
        status: action.payload,
        error: null,
      });
    case clusterReadyActionTypes.ERROR:
      return Object.assign({}, state, {
        loaded: true,
        error: action.payload,
      });
    default:
      return state;
    }
  },

  // Stores the "dirtiness" of UI fields. Should be saved and restored
  dirty: (state, action) => {
    // this isArray check is just to prevent errors with old progress files & dev mode
    if (!state || _.isArray(state)) {
      return {};
    }

    switch (action.type) {
    case configActionTypes.REMOVE_AT: {
      const {path, index} = action.payload;
      const array = _.isString(path) ? path.split('.') : path;
      array.push(index.toString());
      return fromJS(state).deleteIn(array).toJS();
    }
    case dirtyActionTypes.ADD: {
      // {awsTags: [{key: true}]
      const split = action.payload.split('.')
        .map(s => {
          const int = parseInt(s, 10);
          return isNaN(int) ? s : int;
        });
      // mirror structure of clusterconfig for storing field dirtiness
      const obj = _.cloneDeep(state);
      let s = split[0];
      let current = obj;
      for (let i = 1; i < split.length; ++i) {
        if (!_.has(current, s) || current[s] === null) {
          current[s] = _.isInteger(split[i]) ? [] : {};
        }
        current = current[s];
        s = split[i];
      }
      current[s] = true;
      return obj;
    }
    case dirtyActionTypes.CLEAN: {
      const array = _.isString(action.payload) ? action.payload.split('.') : action.payload;
      return fromJS(state).deleteIn(array).toJS();
    }
    default:
      return state;
    }
  },

  // A value guaranteed to monotonically increase. Should be saved/restored.
  sequence: (state, action) => {
    if (state === undefined) {
      return 0;
    }

    switch (action.type) {
    case sequenceActionTypes.INCREMENT:
      return state + 1;
    default:
      return state;
    }
  },
});

function filterClusterConfig(cc={}) {
  Object.keys(cc).forEach(k => {
    if (!DEFAULT_CLUSTER_CONFIG.hasOwnProperty(k)) {
      console.error(`Removed clusterConfig.${k} because it's not in defaults.`);
      delete cc[k];
      return;
    }
    let defaultType;
    let restoreType;
    try {
      defaultType = Object.getPrototypeOf(DEFAULT_CLUSTER_CONFIG[k]);
    } catch (unused) {
      defaultType = typeof DEFAULT_CLUSTER_CONFIG[k];
    }
    try {
      restoreType = Object.getPrototypeOf(cc[k]);
    } catch (unused) {
      restoreType = typeof cc[k];
    }
    if (defaultType !== restoreType) {
      console.error(`Removed clusterConfig.${k}. Should be ${defaultType} but was ${restoreType}`);
      delete cc[k];
      return;
    }
    if (cc[k].inFly) {
      console.debug(`Set clusterConfig.${k}.inFly to false`);
      cc[k].inFly = false;
    }
  });
  delete cc.inFly;
  return cc;
}

// Shim for tests
export const reducer = (state, action) => {
  let restored;
  switch (action.type) {
  case restoreActionTypes.RESTORE_STATE:
    restored = Object.assign({}, action.payload, {
      serverFacts: state.serverFacts,
      cluster: state.cluster,
    });
    restored.aws = _.defaults(restored.aws, DEFAULT_AWS);
    if (restored.clusterConfig) {
      restored.clusterConfig = _.defaults(filterClusterConfig(restored.clusterConfig), DEFAULT_CLUSTER_CONFIG);
    }
    break;
  default:
    restored = state;
  }

  return reducersTogether(restored, action);
};

// Preserves only the savable bits of state
export const savable = (state) => {
  const {dirty, clusterConfig, sequence} = state;
  return {dirty, clusterConfig, sequence};
};
