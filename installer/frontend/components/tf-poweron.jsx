import _ from 'lodash';
import classNames from 'classnames';
import React from 'react';
import { connect } from 'react-redux';
import { saveAs } from 'file-saver';

import { Alert } from './alert';
import { WaitingLi } from './ui';
import { AWS_DomainValidation } from './aws-domain-validation';
import { ResetButton } from './reset-button';
import { TFDestroy } from '../aws-actions';
import { CLUSTER_NAME, PLATFORM_TYPE, getTectonicDomain } from '../cluster-config';
import { AWS_TF, BARE_METAL_TF } from '../platforms';
import { commitToServer, observeClusterStatus } from '../server';

const stateToProps = ({cluster, clusterConfig}) => {
  const status = cluster.status || {terraform: {}};
  const { terraform, tectonic } = status;
  return {
    terraform: {
      action: terraform.action,
      tfError: terraform.error,
      output: terraform.output,
      statusMsg: terraform.status ? terraform.status.toLowerCase() : '',
    },
    clusterName: clusterConfig[CLUSTER_NAME],
    platformType: clusterConfig[PLATFORM_TYPE],
    tectonic: tectonic || {},
    tectonicDomain: getTectonicDomain(clusterConfig),
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

  componentWillUpdate ({terraform}) {
    const output = _.get(terraform, 'output', '');
    const oldOutput = _.get(this.props, 'terraform.output', '');
    // Output can be 10MB, so compare lengths first. Don't trust the JS engine to be smart.
    const sameOutput = output.length === oldOutput.length && output === oldOutput;
    if (sameOutput || this.state.showLogs === false) {
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
    const {clusterName, platformType, terraform, tectonic, tectonicDomain} = this.props;
    const {action, tfError, output, statusMsg} = terraform;
    const state = this.state;
    const showLogs = state.showLogs === null ? statusMsg !== 'success' : state.showLogs;
    const terraformRunning = statusMsg === 'running';

    const consoleSubsteps = [];

    if (action === 'apply' && tectonic.console) {
      if (platformType === AWS_TF) {
        consoleSubsteps.push(<AWS_DomainValidation key="domain" />);
      }

      const dnsReady = tectonic.console.success || ((tectonic.console.message || '').search('no such host') === -1);
      consoleSubsteps.push(
        <WaitingLi done={dnsReady && !terraformRunning} key='dnsReady'>
          Resolving <a href={`https://${tectonicDomain}`} target="_blank">{tectonicDomain}</a>
        </WaitingLi>
      );

      const services = [
        {key: 'etcd', name: 'Etcd'},
        {key: 'kubernetes', name: 'Kubernetes'},
        {key: 'identity', name: 'Tectonic Identity'},
        {key: 'ingress', name: 'Tectonic Ingress Controller'},
        {key: 'console', name: 'Tectonic Console'},
        {key: 'tectonicSystem', name: 'other Tectonic services'},
      ];

      if (!tectonic.isEtcdSelfHosted) {
        services.shift(0);
      }

      const anyFailed = _.some(services, s => tectonic[s.key].failed);
      const allDone = _.every(services, s => tectonic[s.key].success);

      let tectonicSubsteps = null;
      if ((!allDone || anyFailed) && !terraformRunning) {
        tectonicSubsteps = _.map(services, service => <WaitingLi done={tectonic[service.key].success} error={tectonic[service.key].failed} key={service.key} substep={true}>
          Starting {service.name}
        </WaitingLi>);
      }

      consoleSubsteps.push(
        <WaitingLi done={allDone} error={anyFailed} key="tectonicReady">
          Starting Tectonic
          <br />
          <ul className="service-launch-progress__steps">{ tectonicSubsteps }</ul>
        </WaitingLi>
      );
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
      { platformType !== BARE_METAL_TF &&
        <p>
          Kubernetes is starting up. We're committing your cluster details.
          Grab some tea and sit tight. This process can take up to 20 minutes.
          Status updates will appear below.
        </p>
      }
      { platformType === BARE_METAL_TF &&
        <div>
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
        </div>
      }
      <hr />
      <div className="row">
        <div className="col-xs-12">
          <ul className="wiz-launch-progress">
            <WaitingLi done={statusMsg === 'success'} error={tfError}>
              Terraform {action} {statusMsg}
              {output && <div className="pull-right" style={{fontSize: "13px"}}>
                <a onClick={() => this.setState({showLogs: !showLogs})}>
                  { showLogs ? <span><i className="fa fa-angle-up"></i>&nbsp;&nbsp;Hide logs</span>
                              : <span><i className="fa fa-angle-down"></i>&nbsp;&nbsp;Show logs</span> }
                </a>
                <span className="spacer"></span>
                <a onClick={() => saveAs(new Blob([output], {type: 'text/plain'}), `tectonic-${clusterName}.log`)}>
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
// we never reset this value...
let ready = false;

TF_PowerOn.canNavigateForward = ({cluster}) => {
  const {tectonic, terraform} = _.get(cluster, 'status') || {};
  if (_.get(terraform, 'action') === 'destroy') {
    return false;
  }
  ready = ready || (
    _.get(tectonic, 'console.success') === true
    && _.get(tectonic, 'tectonicSystem.success') === true
    && _.toLower(_.get(terraform, 'status')) !== 'running');

  return ready;
};
