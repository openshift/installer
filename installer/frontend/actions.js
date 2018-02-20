import _ from 'lodash';

import { DEFAULT_CLUSTER_CONFIG } from './cluster-config';

export const awsActionTypes = {
  SET: 'AWS_SET',
};
export const awsActions = {
  error: (key, error) => awsActions.set(key, {error, inFly: false, value: []}),
  loaded: (key, value) => awsActions.set(key, {error: null, inFly: false, value}),
  set: (key, data) => ({payload: {[key]: data}, type: awsActionTypes.SET}),
};

export const configActionTypes = {
  ADD_IN: 'CONFIG_ACTION_ADD_IN',
  APPEND: 'CONFIG_ACTION_APPEND',
  REMOVE_FIELD_LIST_ROW: 'CONFIG_ACTION_REMOVE_FIELD_LIST_ROW',
  SET: 'CONFIG_ACTION_SIMPLE_SET',
  SET_IN: 'CONFIG_ACTION_SET_IN',
  BATCH_SET_IN: 'CONFIG_ACTION_BATCH_SET_IN',
  RESET: 'CONFIG_ACTION_RESET',
};

export const clusterReadyActionTypes = {
  ERROR: 'CLUSTER_READY_ERROR',
  NOT_READY: 'CLUSTER_READY_NOT_READY',
  STATUS: 'CLUSTER_READY_STATUS',
};
export const clusterReadyActions = {
  error: payload => ({payload, type: clusterReadyActionTypes.ERROR}),
  notReady: () => ({type: clusterReadyActionTypes.NOT_READY}),
  status: payload => ({payload, type: clusterReadyActionTypes.STATUS}),
};

export const dirtyActionTypes = {
  ADD: 'DIRTY_ADD',
  CLEAN: 'DIRTY_CLEAN',
};
export const dirtyActions = {
  add: field => ({type: dirtyActionTypes.ADD, payload: field}),
  clean: field => ({type: dirtyActionTypes.CLEAN, payload: field}),
};

export const eventErrorsActionTypes = {
  ERROR: 'EVENT_ERRORS_ERROR',
};
export const eventErrorsActions = {
  error: payload => ({payload, type: eventErrorsActionTypes.ERROR}),
};

export const loadFactsActionTypes = {
  ERROR: 'LOAD_FACTS_ERROR',
  LOADED: 'LOAD_FACTS_LOADED',
};
export const loadFactsActions = {
  error: payload => ({payload, type: loadFactsActionTypes.ERROR}),
  loaded: payload => ({payload, type: loadFactsActionTypes.LOADED}),
};

export const restoreActionTypes = {
  RESTORE_STATE: 'RESTORE_RESTORE_STATE',
};
export const restoreActions = {
  restore: payload => ({payload, type: restoreActionTypes.RESTORE_STATE}),
};

export const serverActionTypes = {
  COMMIT_REQUESTED: 'SERVER_COMMIT_REQUESTED',
  COMMIT_SENT: 'SERVER_COMMIT_SENT',
  COMMIT_SUCCESSFUL: 'SERVER_COMMIT_SUCCESSFUL',
  COMMIT_FAILED: 'SERVER_COMMIT_FAILED',
};
export const serverActions = {
  requested: () => ({type: serverActionTypes.COMMIT_REQUESTED}),
  sent: () => ({type: serverActionTypes.COMMIT_SENT}),
  successful: payload => ({payload, type: serverActionTypes.COMMIT_SUCCESSFUL}),
  failed: payload => ({payload, type: serverActionTypes.COMMIT_FAILED}),
};

// Commit state machine:
//   IDLE|FAILED -> REQUESTED -> WAITING -> SUCCEEDED|FAILED
export const commitPhases = {
  IDLE: 'COMMIT_IDLE',
  REQUESTED: 'COMMIT_REQUESTED',
  WAITING: 'COMMIT_WAITING',
  SUCCEEDED: 'COMMIT_SUCCEEDED',
  FAILED: 'COMMIT_FAILED',
};

export const FIELDS = {};
export const FIELD_TO_DEPS = {};
export const FORMS = {};

const getField = name => {
  if (!FIELDS[name]) {
    throw new Error(`Field ${name} not found`);
  }
  return FIELDS[name];
};

// TODO (ggreer) standardize on order of params. is dispatch first or last?
export const configActions = {
  addIn: (path, value, dispatch) => dispatch({payload: {path, value}, type: configActionTypes.ADD_IN}),
  set: (payload, dispatch) => dispatch({type: configActionTypes.SET, payload}),
  setIn: (path, value, dispatch) => dispatch({payload: {path, value}, type: configActionTypes.SET_IN}),
  batchSetIn: (dispatch, payload) => dispatch({payload, type: configActionTypes.BATCH_SET_IN}),

  // TODO: (kans) move below to form actions...
  removeField: (fieldListId, index) => (dispatch, getState) => {
    const fieldList = getField(fieldListId);
    dispatch({payload: {fieldListId, index}, type: configActionTypes.REMOVE_FIELD_LIST_ROW});
    fieldList.validate(dispatch, getState);
  },
  appendField: fieldListId => (dispatch, getState) => {
    const fieldList = getField(fieldListId);
    const value = _.mapValues(fieldList.rowFields, 'default');
    dispatch({payload: {path: fieldListId, value}, type: configActionTypes.APPEND});
    fieldList.validate(dispatch, getState);
  },
  refreshExtraData: fieldName => (dispatch, getState) => {
    const field = getField(fieldName);
    field.getExtraStuff(dispatch, getState);
  },
  updateField: (fieldName, inputValue) => (dispatch, getState) => {
    const [name, ...split] = fieldName.split('.');
    const field = getField(name);
    return field.update(dispatch, inputValue, getState, split);
  },
};

export const __deleteEverything__ = () => {
  [FIELDS, FIELD_TO_DEPS, FORMS, DEFAULT_CLUSTER_CONFIG]
    .forEach(o => _.keys(o).forEach(k => delete o[k]));

  ['error', 'inFly', 'extra'].forEach(k => DEFAULT_CLUSTER_CONFIG[k] = {});

  return {type: configActionTypes.RESET};
};

export const validateFields = async (ids, getState, dispatch, updatedId, isNow) => {
  const unvisitedIds = ids;

  // Just shake the array really hard until all the nodes fall out...
  while (unvisitedIds.length) {
    // All the fields that have already had their dependencies validated
    const toVisit = unvisitedIds.filter(id => !_.intersection(unvisitedIds, FIELDS[id].dependencies).length);
    if (!toVisit.length) {
      throw new Error(`Unresolvable fields: ${unvisitedIds}`);
    }
    await Promise.all(toVisit.map(
      id => FIELDS[id].getExtraStuff(dispatch, getState, isNow)
        .then(() => FIELDS[id].validate(dispatch, getState, updatedId, isNow))
    ));
    _.pullAll(unvisitedIds, toVisit);
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

    _.each(f.dependencies, d => addDep(f.id, d));
  });

  // HACK to avoid figuring out the "correct" order in FIELD_TO_DEPS
  // ... peers can have deps on the same branch
  _.each(form.dependencies, d => addDep(form.id, d));
};
