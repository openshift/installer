import _ from 'lodash';
import Baby from 'babyparse';
import classNames from 'classnames';
import React from 'react';
import { connect } from 'react-redux';

import { Alert } from './alert';
import { configActionTypes } from '../actions';
import { validate } from '../validate';
import { readFile } from '../readfile';
import {
  BM_MASTERS,
  BM_MASTERS_COUNT,
  BM_WORKERS,
  BM_WORKERS_COUNT,
} from '../cluster-config';

import { WithClusterConfig, NumberInput, Input } from './ui';

const TOO_MANY_MASTERS = 9;
const TOO_MANY_WORKERS = 1000;

const validateMasters = v => validate.int({min: 1, max: TOO_MANY_MASTERS})(v) || validate.isOdd(v);
const validateWorkers = validate.int({min: 1, max: TOO_MANY_WORKERS});

const countBy = (collection, f) => {
  const ret = new Map();
  collection.forEach(obj => {
    const k = f(obj);
    const v = ret.get(k) || 0;
    ret.set(k, v + 1);
  });

  return ret;
};

// label is "Master" or "Worker"
// index is index of the node
const macFieldID = (label, index) => `nodetable:${label}:${index}:mac`;
const nameFieldID = (label, index) => `nodetable:${label}:${index}:name`;

// TODO (ggreer) make a real modal with real modal classes
class BulkUpload extends React.Component {
  constructor (props) {
    super(props);
    this.state = {
      name: null,
      macCol: 0,
      nameCol: 1,
      csv: null,
    };
  }

  handleUpload (e) {
    const blob = e.target.files.item(0);
    readFile(blob)
    .then(result => {
      const csv = Baby.parse(result, {delimiter: ','});
      this.setState({
        name: blob.name,
        macCol: 0,
        nameCol: 1,
        csv,
      });
    })
    .catch((msg) => {
      console.error(msg);
    });
  }

  handleSelectMACColumn (e) {
    this.setState({
      macCol: parseInt(e.target.value, 10),
    });
  }

  handleSelectNameColumn (e) {
    this.setState({
      nameCol: parseInt(e.target.value, 10),
    });
  }

  cancel () {
    this.props.close();
  }

  handleDone () {
    const {nameCol, macCol, csv} = this.state;
    const rows = csv.data.slice(1).filter(row => {
      // BabyParse will append a single [""] row to a well-formed CSV,
      // the following happens to fix that, and forgive other
      // possible CSV weirdnesses.
      return row.length > Math.max(nameCol, macCol);
    });
    const nodes = rows.map(row => {
      return {
        name: row[nameCol],
        mac: row[macCol],
      };
    });
    this.props.updateNodes(nodes, nodes.length);
    this.props.close();
  }


  render () {
    const { csv, name, nameCol, macCol } = this.state;

    let body;
    if (!csv) {
      body = <div>
        <div>
          Select a CSV file to populate the node addresses
        </div>
        <div className="wiz-minimodal__body">
          <input type="file" onChange={e => this.handleUpload(e)} />
          <div className="wiz-upload-csv-settings">
            <p>After uploading, you can select which columns correspond to the required data.</p>
          </div>
        </div>
      </div>;
    } else if (csv.errors.length) {
      body = <Alert severity="error">
        Error parsing CSV:
        <ul>
          {csv.errors.map((e, i) => <li key={i}>{e.message} on line {e.row}.</li>)}
        </ul>
      </Alert>;
    } else {
      const options = csv.data[0].map((txt, ix) => {
        return <option value={ix} key={`${ix}:${txt}`}>{txt}</option>;
      });

      body = <div>
        <div className="row">
          <div className="col-xs-3">
            <label>CSV File</label>
          </div>
          <div className="col-xs-6">
            {name}
          </div>
          <div className="col-xs-3">
            <a onClick={() => this.cancel()}>change file</a>
          </div>
        </div>
        <div className="wiz-minimodal__body">
          <div className="wiz-upload-csv-settings">
            <div>Choose the CSV Column that matches each input</div>
            <div className="row wiz-minimodal__controlblock">
              <div className="col-xs-3">
                <label htmlFor="mac-column">Mac Address</label>
              </div>
              <div className="col-xs-6">
                <select id="mac-column"
                        onChange={e => this.handleSelectMACColumn(e)}
                        defaultValue={macCol}>
                  {options}
                </select>
              </div>
            </div>
            <div className="row wiz-minimodal__controlblock">
              <div className="col-xs-3">
                <label htmlFor="name-column">Node Name</label>
              </div>
              <div className="col-xs-6">
                <select id="name-column"
                        onChange={e => this.handleSelectNameColumn(e)}
                        defaultValue={nameCol}>
                  {options}
                </select>
              </div>
            </div>
          </div>
        </div>
      </div>;
    }

    const doneClasses = classNames('btn btn-primary', {disabled: !csv});

    return (
      <div className="wiz-minimodal">
        {body}
        <div className="wiz-minimodal__actions">
          <button type="button" className={doneClasses} onClick={e => this.handleDone(e)}>Done</button>
          <button className="btn btn-link" onClick={() => this.cancel()}>Cancel</button>
        </div>
      </div>
    );
  }
}

// Table of input fields for manual entry of node information
const NodeTable = ({count, theseNodes, allNodes, label, updateNodes}) => {
  const updatedNodes = theseNodes.slice();
  const nodeElems = [];
  for (let i = 0; i < count; i++) {
    const node = theseNodes[i] || {mac: '', name: ''};
    const macOnInput = (mac) => {
      const newNode = Object.assign({}, node, {mac});
      updatedNodes[i] = newNode;
      updateNodes(updatedNodes, count);
    };
    const nameOnInput = (name) => {
      const newNode = Object.assign({}, node, {name});
      updatedNodes[i] = newNode;
      updateNodes(updatedNodes, count);
    };
    const startprops = i > 0 ? {} : {autoFocus: true};

    const duplicateMACs = allNodes.filter(n => n.mac && n.mac === node.mac);
    const duplicateNames = allNodes.filter(n => n.name && n.name === node.name);

    nodeElems.push(
      <div className="row wiz-minitable__row" key={i}>
        <div className="col-xs-3">
          <span className="wiz-minitable__label">{label} {i + 1}:</span>
        </div>
        <div className="col-xs-4">
          <Input
              id={macFieldID(label, i)}
              className="wiz-node-field"
              forceDirty={duplicateMACs.length > 1}
              invalid={duplicateMACs.length > 1 ?
                       'This MAC address is already in use by another node' :
                       validate.MAC(node.mac)}
              value={node.mac}
              placeholder="MAC address"
              onValue={macOnInput}
              {...startprops} />
        </div>
        <div className="col-xs-4">
          <Input
              id={nameFieldID(label, i)}
              className="wiz-node-field"
              forceDirty={duplicateNames.length > 1}
              invalid={duplicateNames.length > 1 ?
                       'This name is already in use by another node' :
                       validate.host(node.name)}
              value={node.name}
              placeholder="node.domain.com"
              onValue={nameOnInput} />
        </div>
      </div>
    );
  }

  return (
    <div>
      <div className="row wiz-minitable__header">
        <div className="col-xs-3">Profile Name</div>
        <div className="col-xs-4">MAC Address</div>
        <div className="col-xs-4">Domain Name</div>
      </div>
      {nodeElems}
    </div>
  );
};

class NodeForm extends React.Component {
  constructor() {
    super();
    this.state = {bulkUpload: false};
  }

  render() {
    if (this.state.bulkUpload) {
      return <BulkUpload close={() => this.setState({bulkUpload: false})} updateNodes={this.props.updateNodes} />;
    }

    return (
      <div>
        <div className="form-group">
          <a onClick={() => this.setState({bulkUpload: true})}>
          <span className="fa fa-upload"></span> Bulk Upload Addresses</a>
        </div>
        <NodeTable {...this.props} />
      </div>
    );
  }
}

const CONTROLLERS_FILE = 'CONTROLLERS_FILE';

export const BM_Controllers = connect(
  ({clusterConfig}) => {
    const masters = clusterConfig[BM_MASTERS];
    return {
      theseNodes: masters,
      allNodes: masters.concat(clusterConfig[BM_WORKERS]),
      count: parseInt(clusterConfig[BM_MASTERS_COUNT], 10),
    };
  },
  (dispatch) => {
    return {
      updateNodes: (nodes, count) => {
        dispatch({
          type: configActionTypes.SET_MASTERS_LIST,
          payload: {nodes, count},
        });
      },
    };
  }
)(({count, theseNodes, allNodes, updateNodes}) => {
  if (count > TOO_MANY_MASTERS) {
    count = TOO_MANY_MASTERS;
  }
  if (count < 1 || !_.isInteger(count)) {
    count = 1;
  }
  return (
    <div>
      <div className="form-group">
        Master nodes run essential cluster services and don't run end-user apps.
      </div>
      <div className="form-group">
        <div className="row">
          <div className="col-xs-3">
            <label htmlFor={BM_MASTERS_COUNT}>Masters</label>
          </div>
          <div className="col-xs-9">
            <WithClusterConfig field={BM_MASTERS_COUNT} validator={validateMasters}>
              <NumberInput
                id={BM_MASTERS_COUNT}
                className="wiz-super-short-input"
                min="1"
                max={TOO_MANY_MASTERS} />
            </WithClusterConfig>
            <p className="text-muted">An odd number of masters is required.</p>
          </div>
        </div>
      </div>
      <div className="form-group">
        Enter the MAC addresses of the nodes you'd like to use as masters,
        and the host names you'll use to refer to them.
      </div>
      <div className="form-group">
        <NodeForm count={count}
                  theseNodes={theseNodes}
                  allNodes={allNodes}
                  label="Master"
                  file={CONTROLLERS_FILE}
                  updateNodes={updateNodes} />
      </div>
    </div>
  );
});
BM_Controllers.canNavigateForward = ({clusterConfig}) => {
  if (validateMasters(clusterConfig[BM_MASTERS_COUNT])) {
    return false;
  }
  const masters = clusterConfig[BM_MASTERS];
  const mastersOkSet = masters.filter((m) => {
    return m && !validate.MAC(m.mac) && !validate.host(m.name);
  });

  if (mastersOkSet.length < parseInt(clusterConfig[BM_MASTERS_COUNT], 10)) {
    return false;
  }

  // In order to prevent weird lockouts and invalidation at a distance,
  // the deduplicate validation for controllers and workers isn't
  // symmetric. In particular, the Controllers form is valid if it
  // contains duplicates of Worker but not if the masters group has
  // duplicates within itself.
  const nameCounts = countBy(masters, n => n.name);
  const macCounts = countBy(masters, n => n.mac);
  for (let i = 0; i < masters.length; i++) {
    const masterI = masters[i];
    if (nameCounts.get(masterI.name) > 1) {
      return false;
    }
    if (macCounts.get(masterI.mac) > 1) {
      return false;
    }
  }

  return true;
};

const WORKERS_FILE = 'WORKERS_FILE';

export const BM_Workers = connect(
  ({clusterConfig}) => {
    const workers = clusterConfig[BM_WORKERS];
    return {
      theseNodes: workers,
      allNodes: workers.concat(clusterConfig[BM_MASTERS]),
      count: parseInt(clusterConfig[BM_WORKERS_COUNT], 10),
    };
  },
  (dispatch) => {
    return {
      updateNodes: (nodes, count) => {
        dispatch({
          type: configActionTypes.SET_WORKERS_LIST,
          payload: {nodes, count},
        });
      },
    };
  }
)(({count, theseNodes, allNodes, updateNodes}) => {
  return (
    <div>
      <div className="form-group">
        Worker nodes run your end-user's apps. The cluster software automatically shares load
        between these nodes.
      </div>
      <div className="form-group">
        <div className="row">
          <div className="col-xs-3">
            <label htmlFor={BM_WORKERS_COUNT}>Workers</label>
          </div>
          <div className="col-xs-9">
            <WithClusterConfig field={BM_WORKERS_COUNT} validator={validateWorkers}>
              <NumberInput
                id={BM_WORKERS_COUNT}
                className="wiz-super-short-input"
                min="1"
                max={TOO_MANY_WORKERS} />
            </WithClusterConfig>
            <p className="text-muted">Workers can be added and removed at any time.</p>
          </div>
        </div>
      </div>
      <div className="form-group">
        Enter the MAC addresses of the nodes you'd like to use as workers,
        and the host names you'll use to refer to them.
      </div>
      <div className="form-group">
        <NodeForm count={count}
                  theseNodes={theseNodes}
                  allNodes={allNodes}
                  label="Worker"
                  file={WORKERS_FILE}
                  updateNodes={updateNodes} />
      </div>
    </div>
  );
});
BM_Workers.canNavigateForward = ({clusterConfig}) => {
  if (validateWorkers(clusterConfig[BM_WORKERS_COUNT])) {
    return false;
  }
  const workers = clusterConfig[BM_WORKERS];
  const workersOk = workers.filter((m) => {
    return m && !validate.MAC(m.mac) && !validate.host(m.name);
  });

  let workersExpected = parseInt(clusterConfig[BM_WORKERS_COUNT], 10);
  if (isNaN(workersExpected)) {
    workersExpected = 3;
  }

  if (workersOk.length < workersExpected) {
    return false;
  }

  // The worker form is invalid if workers have the same mac or name
  // as other workers, or if they have the same mac or name as
  // controller nodes.
  const allNodes = workers.concat(clusterConfig[BM_MASTERS]);
  const nameCounts = countBy(allNodes, n => n.name);
  const macCounts = countBy(allNodes, n => n.mac);
  for (let i = 0; i < workers.length; i++) {
    if (nameCounts.get(workers[i].name) > 1) {
      return false;
    }
    if (macCounts.get(workers[i].mac) > 1) {
      return false;
    }
  }

  return true;
};
