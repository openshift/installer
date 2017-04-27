import _ from 'lodash';
import classNames from 'classnames';
import React from 'react';
import { connect } from 'react-redux';
import { saveAs } from 'file-saver';

import { Alert } from './alert';
import { WaitingLi } from './ui';
import { AWS_DomainValidation } from './aws-domain-validation';
import { ResetButton } from './reset-button';
import { CLUSTER_NAME } from '../cluster-config';
import { TFDestroy } from '../aws-actions';
import { commitToServer, observeClusterStatus } from '../server';

const stateToProps = ({cluster, clusterConfig}) => {
  const status = cluster.status || {};
  return {
    action: status.action,
    clusterName: clusterConfig[CLUSTER_NAME],
    error: status.error,
    output: status.output,
    outputBlob: new Blob([status.output], {type: 'text/plain'}),
    statusMsg: status.status ? status.status.toLowerCase() : '',
    tectonicConsole: status.tectonicConsole || {},
  };
};

const dispatchToProps = dispatch => ({
  TFDestroy: () => dispatch(TFDestroy()).then(() => dispatch(observeClusterStatus)),
  TFRetry: () => dispatch(commitToServer(false, true)),
});

export const AWS_TF_PowerOn = connect(stateToProps, dispatchToProps)(
class AWS_TF_PowerOn extends React.Component {
  constructor (props) {
    super(props);
    this.state = {showLogs: null};
  }
  componentDidMount () {
    if (this.outputNode) {
      this.outputNode.scrollTop = this.outputNode.scrollHeight - this.outputNode.clientHeight;
    }
  }

  componentWillUpdate ({output}) {
    const node = this.outputNode;
    if (!node || output === this.props.output) {
      this.shouldScroll = false;
      return;
    }

    this.shouldScroll = node.scrollHeight - node.clientHeight <= node.scrollTop;
  }

  componentDidUpdate () {
    if (this.shouldScroll && this.outputNode) {
      this.outputNode.scrollTop = this.outputNode.scrollHeight - this.outputNode.clientHeight;
    }
  }

  retry () {
    // eslint-disable-next-line no-alert
    if (window.config.devMode || window.confirm('Are you sure you want to re-run terraform apply?')) {
      this.props.TFRetry();
    }
  }

  destroy () {
    // eslint-disable-next-line no-alert
    if (window.config.devMode || window.confirm('Are you sure you want to destroy your cluster?')) {
      this.props.TFDestroy().catch(tfError => this.setState({tfError}));
    }
  }

  render () {
    const {action, clusterName, error, output, outputBlob, statusMsg, tectonicConsole} = this.props;
    const state = this.state;
    const showLogs = state.showLogs === null ? statusMsg !== 'success' : state.showLogs;
    const terraformRunning = statusMsg === 'running';

    const consoleSubsteps = [];
    if (action === 'apply') {
      let msg = `Resolving ${tectonicConsole.instance}`;
      const dnsReady = (tectonicConsole.message || '').search('no such host') === -1;
      consoleSubsteps.push(<AWS_DomainValidation key="domain" />);
      consoleSubsteps.push(
        <WaitingLi done={dnsReady} key="dns" substep={true}>
          <span title={msg}>{msg}</span>
        </WaitingLi>
      );
      msg = `Starting Tectonic console`;
      consoleSubsteps.push(
        <WaitingLi done={tectonicConsole.ready} key="consoleReady" substep={true}>
          <span title={msg}>{msg}</span>
        </WaitingLi>
      );
    }

    return <div>
      <div className="row">
        <div className="col-xs-12">
          Kubernetes is starting up. We're commiting your cluster details.
          Grab some tea and sit tight. This process can take up to 20 minutes.
          Status updates will appear below.
        </div>
      </div>
      <hr />
      <div className="row">
        <div className="col-xs-12">
          <div className="wiz-launch__progress">
            <ul className="wiz-launch-progress">
              <WaitingLi done={statusMsg === 'success'} error={error}>
                Terraform {action} {statusMsg}
                {output && <div className="pull-right">
                  <a className="spacer-right" onClick={() => this.setState({showLogs: !showLogs})}>
                    { showLogs ? <span><i className="fa fa-angle-up"></i>&nbsp;&nbsp;Hide logs</span>
                               : <span><i className="fa fa-angle-down"></i>&nbsp;&nbsp;Show logs</span> }
                  </a>
                  <a onClick={() => saveAs(outputBlob, `tectonic-${clusterName}.log`)}>
                    <i className="fa fa-download"></i>&nbsp;&nbsp;Save log
                  </a>
                </div>}
                { showLogs && output &&
                  <div className="log-pane">
                    <div className="log-pane__header">
                      <div className="log-pane__header__message">Terraform logs</div>
                    </div>
                    <div className="log-pane__body">
                      <div className="log-area">
                        <div className="log-scroll-pane" ref={node => this.outputNode = node}>
                          <div className="log-contents">
                            <div className="log-contents__block">{output}</div>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                }
              </WaitingLi>
              { consoleSubsteps }
            </ul>
          </div>
        </div>
      </div>
      <div className="row">
        <div className="col-xs-6">
          <button className={classNames("btn btn-default", {disabled: terraformRunning})} onClick={() => this.retry()}>
            <i className="fa fa-exclamation-triangle"></i>&nbsp;&nbsp;Retry Terraform apply
          </button>
        </div>
        <div className="col-xs-6">
          <button className={classNames("btn btn-default pull-right", {disabled: terraformRunning})} onClick={() => this.destroy()}>
            <i className="fa fa-trash"></i>&nbsp;&nbsp;Destroy cluster
          </button>
        </div>
      </div>
      { state.tfError &&
        <div className="row">
          <div className="col-xs-12">
           <Alert severity="error">{state.tfError}</Alert>
          </div>
        </div>
      }
      <hr />
      <div className="row">
        <div className="col-xs-12">
        { error && <Alert severity="error">{error.toString()}</Alert> }
        <br />
        </div>
      </div>
      <div className="row">
        <div className="col-xs-12">
          <a href="/terraform/assets" download>
            <button className={classNames("btn btn-primary wiz-giant-button pull-right", {disabled: terraformRunning})}>
              <i className="fa fa-download"></i>&nbsp;&nbsp;Download assets
            </button>
          </a>
          <ResetButton />
        </div>
      </div>
    </div>;
  }
});

// horrible hack to prevent a flapping cluster from redirecting user from connect page back to poweron page
let ready = false;

AWS_TF_PowerOn.canNavigateForward = ({cluster}) => {
  ready = ready || (_.get(cluster, 'status.tectonicConsole.ready') === true
    && _.get(cluster, 'status.status') !== 'running');
  return ready;
};
