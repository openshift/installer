import _ from 'lodash';
import classNames from 'classnames';
import React from 'react';
import { connect } from 'react-redux';
import { saveAs } from 'file-saver';

import { Alert } from './alert';
import { WaitingLi, LoadingLi } from './ui';
import { AWS_DomainValidation } from './aws-domain-validation';
import { ResetButton } from './reset-button';
import { TFDestroy } from '../aws-actions';
import { CLUSTER_NAME, PLATFORM_TYPE } from '../cluster-config';
import { AWS_TF, BARE_METAL_TF } from '../platforms';
import { commitToServer, observeClusterStatus } from '../server';

const stateToProps = ({cluster, clusterConfig}) => {
  const status = cluster.status || {terraform: {}};
  return {
    terraform: {
      action: status.terraform.action,
      tfError: status.terraform.error,
      output: status.terraform.output,
      outputBlob: new Blob([status.terraform.output], {type: 'text/plain'}),
      statusMsg: status.terraform.status ? status.terraform.status.toLowerCase() : '',
    },
    clusterName: clusterConfig[CLUSTER_NAME],
    platformType: clusterConfig[PLATFORM_TYPE],
    tectonic: status.tectonic || {},
  };
};

const dispatchToProps = dispatch => ({
  TFDestroy: () => dispatch(TFDestroy()).then(() => dispatch(observeClusterStatus)),
  TFRetry: () => dispatch(commitToServer(false, true)),
});

export const TF_PowerOn = connect(stateToProps, dispatchToProps)(
class TF_PowerOn extends React.Component {
  constructor (props) {
    super(props);
    this.state = {showLogs: null};
  }

  componentDidMount () {
    if (this.outputNode) {
      this.outputNode.scrollTop = this.outputNode.scrollHeight;
    }
  }

  componentWillUpdate ({output}) {
    if (output === this.props.output || this.state.showLogs === false) {
      this.shouldScroll = false;
      return;
    }

    const node = this.outputNode;
    if (!node) {
      // outputNode will exist once componentDidUpdate fires. Scroll to the bottom at that time.
      this.shouldScroll = true;
      return;
    }
    this.shouldScroll = node.scrollHeight - node.clientHeight <= node.scrollTop + 20;
  }

  componentDidUpdate () {
    if (this.shouldScroll && this.outputNode) {
      this.outputNode.scrollTop = this.outputNode.scrollHeight;
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
      this.props.TFDestroy().catch(xhrError => this.setState({xhrError}));
    }
  }

  render () {
    const {clusterName, platformType, terraform, tectonic} = this.props;
    const {action, tfError, output, outputBlob, statusMsg} = terraform;
    const state = this.state;
    const showLogs = state.showLogs === null ? statusMsg !== 'success' : state.showLogs;
    const terraformRunning = statusMsg === 'running';

    const consoleSubsteps = [];

    if (action === 'apply' && tectonic.console) {
      let msg = <span>Resolving <a href={`https://${tectonic.console.instance}`} target="_blank">{tectonic.console.instance}</a></span>;
      const dnsReady = (tectonic.console.message || '').search('no such host') === -1;
      if (platformType === AWS_TF) {
        consoleSubsteps.push(<AWS_DomainValidation key="domain" />);
      }

      consoleSubsteps.push(
        <WaitingLi done={dnsReady} key='dnsReady'>
          <span title={msg}>{msg}</span>
        </WaitingLi>
      );

      const components = [
        {key: 'kubernetes', name: 'Kubernetes'},
        {key: 'identity', name: 'Tectonic Identity'},
        {key: 'ingress', name: 'Tectonic Ingress Controller'},
        {key: 'console', name: 'Tectonic Console'},
        {key: 'tectonic', name: 'other Tectonic services'},
      ];

      if (tectonic.hasetcd) {
        components.splice(1, 0, {key: 'etcd', name: 'Etcd'});
      }

      let allDone = true, anyFailed = false;

      for (let c of components) {
        anyFailed = anyFailed || tectonic[c.key].failed;
        allDone = !anyFailed && allDone &&
          (c.key === 'console' ? tectonic[c.key].ready : tectonic[c.key].success);
      }

      let tectonicSubsteps = [];

      if (anyFailed || !allDone) {
        for (let c of components) {
          msg = `Starting ${c.name}`;
          let ready = false;
          if (c.key === 'console') {
            ready = tectonic[c.key].ready;
          } else {
            ready = tectonic[c.key].success;
          }
          tectonicSubsteps.push(
            <LoadingLi done={ready} error={tectonic[c.key].failed} key={c.key + 'Ready'}>
              <span title={msg}>{msg}</span>
            </LoadingLi>
          );
        }
      }

      msg = 'Starting Tectonic';
      consoleSubsteps.push(
        <WaitingLi done={allDone} error={anyFailed} key="tectonicReady">
          <span title={msg}>{msg}</span><br />
          <ul>
            {tectonicSubsteps}
          </ul>
        </WaitingLi>
      );
    }
    let platformMsg = <p>
      Kubernetes is starting up. We're committing your cluster details.
      Grab some tea and sit tight. This process can take up to 20 minutes.
      Status updates will appear below.
    </p>;
    if (platformType === BARE_METAL_TF) {
      platformMsg = <div>
        <div className="wiz-herotext">
          <i className="fa fa-power-off wiz-herotext-icon"></i> Power on the nodes
        </div>
        <div className="form-group">
          After powering up, your nodes will provision themselves automatically.
          This process can take up to 30 minutes, while the following happens.
        </div>
        <div className="form-group">
          <ul>
            <li>Container Linux is downloaded and installed to disk (about 200 MB)</li>
            <li>Cluster software is downloaded (about 500 MB)</li>
            <li>One or two reboots may occur</li>
          </ul>
        </div>
      </div>;
    }

    const tfButtonClasses = classNames("btn btn-flat", {disabled: terraformRunning, 'btn-warning': tfError, 'btn-info': !tfError});
    const tfButtons = <div className="row">
      <div className="col-xs-12">
        <button className={tfButtonClasses} onClick={() => this.destroy()}>
          <i className="fa fa-trash"></i>&nbsp;&nbsp;Destroy Cluster
        </button>&nbsp;&nbsp;&nbsp;&nbsp;
        <button className={tfButtonClasses} onClick={() => this.retry()}>
          <i className="fa fa-exclamation-triangle"></i>&nbsp;&nbsp;Retry Terraform Apply
        </button>
      </div>
    </div>;

    return <div>
      { platformMsg }
      <hr />
      <div className="row">
        <div className="col-xs-12">
          <div className="wiz-launch__progress">
            <ul className="wiz-launch-progress">
              <WaitingLi done={statusMsg === 'success'} error={tfError}>
                Terraform {action} {statusMsg}
                {output && <div className="pull-right" style={{fontSize: "13px"}}>
                  <a onClick={() => this.setState({showLogs: !showLogs})}>
                    { showLogs ? <span><i className="fa fa-angle-up"></i>&nbsp;&nbsp;Hide logs</span>
                               : <span><i className="fa fa-angle-down"></i>&nbsp;&nbsp;Show logs</span> }
                  </a>
                  <span className="spacer"></span>
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
                          <div className="log-contents">{output}</div>
                        </div>
                      </div>
                    </div>
                  </div>
                }
              </WaitingLi>
              <li style={{paddingLeft: 20, listStyle: 'none'}}>
              { state.xhrError &&
                <div className="row">
                  <div className="col-xs-12">
                    <Alert severity="error">{state.xhrError}</Alert>
                  </div>
                </div>
              }
              { tfError && <Alert severity="error">{tfError.toString()}</Alert> }
              { !terraformRunning && tfError &&
                <Alert severity="error" noIcon>
                  <b>{_.startCase(action)} Failed</b>. Your installation is blocked. To continue:
                  <ol style={{ paddingLeft: 30, paddingTop: 10, paddingBottom: 10 }}>
                    <li>Save your logs for debugging purposes.</li>
                    <li>Destroy your cluster to clear anything that may have been created.</li>
                    <li>Reapply Terraform.</li>
                  </ol>
                  {tfButtons}
                </Alert>
              }
              { !terraformRunning && !tfError &&
                <Alert severity="info" noIcon>
                  <b>{_.startCase(action)} Succeeded</b>.
                  <p>
                    If you've changed your mind, you can {action === 'apply' ? 'destroy' : 'reapply'} your cluster.
                  </p>
                  {tfButtons}
                </Alert>
              }
              </li>
              { consoleSubsteps }
            </ul>
          </div>
        </div>
      </div>
      <br />
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

TF_PowerOn.canNavigateForward = ({cluster}) => {
  ready = ready || (_.get(cluster, 'status.tectonic.console.ready') === true
    && _.get(cluster, 'status.status') !== 'running');
  return ready;
};
