import _ from 'lodash';

import { DEFAULT_CLUSTER_CONFIG } from './cluster-config';
export const awsActionTypes = {
  SET: 'AWS_SET',
};

export const configActionTypes = {
  APPEND: 'CONFIG_ACTION_APPEND',
  REMOVE_AT: 'CONFIG_ACTION_REMOVE_AT',
  SET: 'CONFIG_ACTION_SIMPLE_SET',
  SET_IN: 'CONFIG_ACTION_SET_IN',
  BATCH_SET_IN: 'CONFIG_ACTION_BATCH_SET_IN',
  MERGE: 'CONFIG_ACTION_MERGE',
  RESET: 'CONFIG_ACTION_RESET',
};

export const clusterReadyActionTypes = {
  ERROR: 'clusterReadyActionTypes.ERROR',
  STATUS: 'clusterReadyActionTypes.CLUSTER_STATUS',
  NOT_READY: 'clusterReadyActionTypes.NOT_READY',
};

export const dirtyActionTypes = {
  ADD: 'DIRTY_ADD',
  CLEAN: 'DIRTY_CLEAN',
};

export const eventErrorsActionTypes = {
  ERROR: 'EVENT_ERRORS_ERROR',
};

export const loadFactsActionTypes = {
  LOADED: 'LOAD_FACTS_LOADED',
  ERROR: 'LOAD_FACTS_ERROR',
};

export const navActionTypes = {
  LOCATION_CHANGE: 'NAV_ACTION_LOCATION_CHANGE',
};

export const navChange = pathname => ({type: navActionTypes.LOCATION_CHANGE, payload: {pathname}});

export const restoreActionTypes = {
  RESTORE_STATE: 'RESTORE_RESTORE_STATE',
};

export const sequenceActionTypes = {
  INCREMENT: 'SEQUENCE_INCREMENT',
};

export const serverActionTypes = {
  COMMIT_REQUESTED: 'COMMIT_REQUESTED',
  COMMIT_SENT: 'COMMIT_SENT',
  COMMIT_SUCCESSFUL: 'COMMIT_SUCCESSFUL',
  COMMIT_FAILED: 'COMMIT_FAILED',
  COMMIT_RESET: 'COMMIT_RESET',
};

// Commit state machine
//
// IDLE|FAILED -> REQUESTED -> WAITING -> SUCCEEDED|FAILED

export const commitPhases = {
  IDLE: 'COMMIT_IDLE',
  REQUESTED: 'COMMIT_REQUESTED',
  WAITING: 'COMMIT_WAITING',
  SUCCEEDED: 'COMMIT_SUCCEEDED',
  FAILED: 'COMMIT_FAILED',
};
const FIELDS = {};
const FIELD_TO_DEPS = {};
const FIELD_TO_FORM = {};
export const FORMS = {};

// TODO (ggreer) standardize on order of params. is dispatch first or last?
export const configActions = {
  setIn: (path, value, dispatch) => dispatch({payload: {path, value}, type: configActionTypes.SET_IN}),
  batchSetIn: (dispatch, payload) => {
    const values = dispatch({payload, type: configActionTypes.BATCH_SET_IN});
    payload.splice(0, payload.length - 1);
    return values;
  },
  append: (path, value, dispatch) => dispatch({payload: {path, value}, type: configActionTypes.APPEND}),
  removeAt: (path, index, dispatch) => dispatch({payload: {path, index}, type: configActionTypes.REMOVE_AT}),
  merge: payload => dispatch => dispatch({payload, type: configActionTypes.MERGE}),
  // TODO: (kans) move below to form actions...
  removeField: (fieldName, i) => (dispatch, getState) => {
    const field = FIELDS[fieldName];
    if (!field) {
      throw new Error(`${fieldName} has no field for removing...`);
    }
    field.remove(dispatch, i, getState);
  },
  appendField: fieldName => (dispatch, getState) => {
    const field = FIELDS[fieldName];
    if (!field) {
      throw new Error(`${fieldName} has no field for appending...`);
    }
    field.append(dispatch, getState);
  },
  refreshExtraData: fieldName => (dispatch, getState) => {
    const field = FIELDS[fieldName];
    if (!field) {
      throw new Error(`${fieldName} has no field`);
    }
    field.getExtraStuff(dispatch, getState().clusterConfig, FIELDS, 0);
  },
  updateField: (fieldName, inputValue) => (dispatch, getState) => {
    const [name, ...split] = fieldName.split('.');
    const field = FIELDS[name];
    if (!field) {
      throw new Error(`${name} has no field`);
    }

    return field.update(dispatch, inputValue, getState, FIELDS, FIELD_TO_DEPS, split);
  },
};

export const __deleteEverything__ = () => {
  [FIELDS, FIELD_TO_DEPS, FIELD_TO_FORM, FORMS, DEFAULT_CLUSTER_CONFIG]
    .forEach(o => _.keys(o).forEach(k => delete o[k]));

  ['error', 'error_async', 'ignore', 'inFly', 'extra']
    .forEach(k => DEFAULT_CLUSTER_CONFIG[k] = {});

  return {type: configActionTypes.RESET};
};

export const validateAllFields = cb => async (dispatch, getState) => {
  const initialCC = getState().clusterConfig;
  const unvisitedFields = new Set(_.values(FIELDS));
  const visitedNames = new Set();

  const visit = async field => {
    const { id } = field;

    visitedNames.add(id);
    unvisitedFields.delete(field);

    const { clusterConfig } = getState();
    // we must update ignores before validation... because validation depends on it
    field.ignoreWhen(dispatch, clusterConfig);
    // TODO: (kans) this is bad
    await field.getExtraStuff(dispatch, clusterConfig, FIELDS, 0);
    await field.validate(dispatch, getState, initialCC, 0);
  };

  // Just shake the array really hard until all the nodes fall out...
  while (unvisitedFields.size > 0) {
    const toVisit = [];

    unvisitedFields.forEach(f => {
      const unresolvedDeps = f.dependencies.filter(d => !visitedNames.has(d)).length;
      if (unresolvedDeps) {
        return;
      }
      toVisit.push(f);
    });

    if (!toVisit.length && unvisitedFields.size > 0) {
      throw new Error(`unresolvable fields: ${unvisitedFields.toJSON().map(f => f.name).join(' ')}`);
    }
    // TODO: (kans) use promise.All here for speeds
    for (let dep of toVisit) {
      await visit(dep);
    }
  }
  if (_.isFunction(cb)) {
    cb();
  }
};

const addDep = (field, dep) => {
  if (!FIELD_TO_DEPS[dep]) {
    FIELD_TO_DEPS[dep] = [field];
    return;
  }
  FIELD_TO_DEPS[dep].push(field);
};

export const registerForm = (form, fields) => {
  const formName = form.id;

  if (FORMS[formName]) {
    throw new Error(`form ${formName} already exists`);
  }

  FORMS[formName] = form;

  // forms can be fields too
  FIELDS[form.id] = form;

  _.each(fields, f => {
    const fieldName = f.id;
    FIELD_TO_FORM[fieldName] = form;
    if (f.isForm) {
      return;
    }
    if (!fieldName) {
      throw new Error(`form ${formName}: field has no name!`);
    }
    if (DEFAULT_CLUSTER_CONFIG[fieldName]) {
      throw new Error(`form ${formName}: field ${fieldName} already exists`);
    }

    DEFAULT_CLUSTER_CONFIG[fieldName] = f.default;
    FIELDS[fieldName] = f;

    _.each(f.dependencies, d => addDep(f, d));
  });

  // HACK to avoid figuring out the "correct" order in FIELD_TO_DEPS
  // ... peers can have deps on the same branch
  _.each(form.dependencies, d => addDep(form, d));
};
