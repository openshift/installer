import React from 'react';
import classNames from 'classnames';
import { connect } from 'react-redux';

import { configActionTypes, eventErrorsActionTypes } from '../actions';
import { compareVersions } from '../utils';
import { validate } from '../validate';
import { TectonicGA } from '../tectonic-ga';

import { WithClusterConfig, Input } from './ui';
import {
  DEFAULT_CLUSTER_CONFIG,
  BM_MATCHBOX_HTTP,
  BM_MATCHBOX_RPC,
  BM_OS_TO_USE,
} from '../cluster-config';

const COREOS_VERSIONS_ERROR_NAME = 'COREOS_VERSIONS_ERROR_NAME';

const stateToProps = ({eventErrors, clusterConfig}) => ({
  http: clusterConfig[BM_MATCHBOX_HTTP],
  osToUse: clusterConfig[BM_OS_TO_USE],
  versionError: eventErrors[COREOS_VERSIONS_ERROR_NAME],
  // TODO: (ggreer) total hack
  bootCfgInfly: clusterConfig.bootCfgInfly,
});

const dispatchToProps = dispatch => ({
  handleHTTP: value => {
    dispatch({
      type: eventErrorsActionTypes.ERROR,
      payload: {
        name: COREOS_VERSIONS_ERROR_NAME,
        error: null,
      },
    });
    dispatch({
      type: configActionTypes.SET,
      payload: {matchboxHTTP: value},
    });
  },
  getOsToUse: matchboxHTTP => {
    if (validate.hostPort(matchboxHTTP)) {
      return;
    }

    dispatch({
      type: configActionTypes.SET,
      payload: {bootCfgInfly: true},
    });

    const endpointURL = `http://${matchboxHTTP}/assets/coreos`;
    return fetch(`/containerlinux/images/matchbox?endpoint=${encodeURIComponent(endpointURL)}`)
    .then(response => {
      if (response.ok) {
        return response.json();
      }

      return response.text().then(txt => Promise.reject(txt));
    })
    .then(value => {
      const available = value.coreos;
      if (!available || available.length === 0) {
        return Promise.reject(`could not find any coreos images at ${endpointURL}`);
      }

      const useVersion = available.map(v => v.version).sort(compareVersions).pop();

      dispatch({
        type: configActionTypes.SET,
        payload: {[BM_OS_TO_USE]: useVersion},
      });

      return Promise.resolve(true);
    })
    .catch(err => {
      dispatch({
        type: eventErrorsActionTypes.ERROR,
        payload: {
          name: COREOS_VERSIONS_ERROR_NAME,
          error: err,
        },
      });
      dispatch({
        type: configActionTypes.SET,
        payload: {[BM_OS_TO_USE]: DEFAULT_CLUSTER_CONFIG[BM_OS_TO_USE]},
      });
    })
    .then(() => dispatch({type: configActionTypes.SET, payload: {bootCfgInfly: false}}));
  },
});

export const BM_Matchbox = connect(stateToProps, dispatchToProps)(
class Matchbox extends React.Component {
  componentDidMount() {
    if (!this.props.osToUse) {
      this.props.getOsToUse(this.props.http);
    }
  }

  componentWillUnmount() {
    clearTimeout(this.timeout);
  }

  onBootCfg (v, ms=1400) {
    this.props.handleHTTP(v);
    if (this.timeout) {
      clearTimeout(this.timeout);
    }
    this.timeout = setTimeout(() => this.props.getOsToUse(v), ms);
  }

  render () {
    const {http, osToUse, versionError, bootCfgInfly} = this.props;
    const osError = osToUse ? false : versionError;
    const iClassNames = classNames('fa', 'fa-refresh', {
      'fa-spin': bootCfgInfly,
    });

    return <div className="row form-group">
      <div className="col-sm-12">
        <div className="form-group">Matchbox will provision nodes during network boot. Enter the matchbox endpoints.</div>
        <div className="form-group">To find your matchbox endpoints, follow the instructions in the {' '}
          <a href="https://coreos.com/tectonic/docs/latest/install/bare-metal/index.html" onClick={TectonicGA.sendDocsEvent} target="_blank">
            Tectonic Deploy Documentation
          </a>.
        </div>
        <div className="form-group">
          <div className="row">
            <div className="col-xs-3">
              <label htmlFor={BM_MATCHBOX_HTTP}>HTTP address</label>
            </div>
            <div className="col-xs-9">
              <Input id={BM_MATCHBOX_HTTP}
                     className="wiz-inline-field wiz-inline-field--protocol"
                     autoFocus="true"
                     prefix={<span className="input__prefix--protocol">http://</span>}
                     placeholder="matchbox.example.com:8080"
                     forceDirty={!!osError}
                     invalid={osError || validate.hostPort(http)}
                     value={http}
                     onValue={v => this.onBootCfg(v)}>
                <button className="btn btn-default" disabled={!!validate.hostPort(http)} onClick={() => this.onBootCfg(http, 0)}>
                  <i className={iClassNames}></i>
                </button>
              </Input>
              <p className="text-muted">Hostname and port of matchbox HTTP endpoint</p>
            </div>
          </div>
          <div className="row">
            <div className="col-xs-3">
              <label htmlFor={BM_MATCHBOX_RPC}>API address</label>
            </div>
            <div className="col-xs-9">
              <WithClusterConfig field={BM_MATCHBOX_RPC} validator={validate.hostPort}>
                <Input id={BM_MATCHBOX_RPC}
                       className="wiz-inline-field wiz-inline-field--protocol"
                       prefix={<span className="input__prefix--protocol">https://</span>}
                       placeholder="matchbox.example.com:8081"/>
              </WithClusterConfig>
              <p className="text-muted">Hostname and port of matchbox API endpoint</p>
            </div>
          </div>
        </div>
      </div>
    </div>;
  }
});

BM_Matchbox.canNavigateForward = ({clusterConfig}) => {
  return !validate.hostPort(clusterConfig[BM_MATCHBOX_HTTP]) &&
         !validate.nonEmpty(clusterConfig[BM_OS_TO_USE]) &&
         !validate.hostPort(clusterConfig[BM_MATCHBOX_RPC]);
};
