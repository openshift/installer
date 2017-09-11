import _ from 'lodash';
import classNames from 'classnames';
import React from 'react';
import { connect } from 'react-redux';
import { saveAs } from 'file-saver';

import { Alert } from './alert';
import { DropdownInline } from './ui';
import { AWS_DomainValidation } from './aws-domain-validation';
import { commitPhases } from '../actions';
import { TFDestroy } from '../aws-actions';
import { CLUSTER_NAME, PLATFORM_TYPE, getTectonicDomain } from '../cluster-config';
import { AWS_TF, BARE_METAL_TF } from '../platforms';
import { commitToServer, observeClusterStatus } from '../server';

const ProgressBar = ({progress, isActive}) => <div className="progress-bar-wrap">
  <div className={`progress-bar progress-bar--${isActive ? 'active' : 'stalled'}`} style={{width: `${progress * 100}%`}}></div>
</div>;

// Estimate the Terraform action progress based on the log output. The intention is to replace this in the future with
// progress information provided by the backend.
const estimateTerraformProgress = terraform => {
  const {output = ''} = terraform;
  if (output === '') {
    return 0;
  }

  let done = output.match(/.*: Creation complete/g) || [];

  // Ignore resources from the bootkube, flannel-vxlan and tectonic modules because they all complete very quickly
  done = done.filter(c => !/module\.(bootkube|flannel-vxlan|tectonic)/.test(c));

  // Approximate number of AWS resources we expect Terraform to create
  const total = 82;

  // We have some output, but are not finished, so don't show the progress as either completely empty or full
  return _.clamp(done.length / total, 0.01, 0.99);
};

const Step = ({pending, done, error, cancel, children, substep}) => {
  const progressClasses = classNames('wiz-launch-progress__step', {
    'wiz-launch-progress__step--substep': substep,
    'wiz-pending-fg': pending,
    'wiz-error-fg': error,
    'wiz-success-fg': done && !error,
    'wiz-cancel-fg': !done && !error && cancel,
    'wiz-running-fg': !done && !error && !cancel && !pending,
  });
  const iconClasses = classNames('fa', 'fa-fw', {
    'fa-circle-o': pending,
    'fa-check-circle': done && !error,
    'fa-ban': error || (cancel && !done),
    'fa-spin fa-circle-o-notch': !done && !error && !cancel && !pending,
  });

  return <div className={progressClasses}>
    {!substep && <i className={iconClasses} style={{margin: '0 3px 0 -3px'}}></i>}{children}
  </div>;
};

const Btn = ({isError, onClick, title}) => <button
  className={`btn btn-flat ${isError ? 'btn-warning' : 'btn-info'}`}
  onClick={onClick}
  style={{marginRight: 15}}
>{title}</button>;

const stateToProps = ({cluster, clusterConfig, commitState}) => {
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
    commitState,
    isAWS: clusterConfig[PLATFORM_TYPE] === AWS_TF,
    isBareMetal: clusterConfig[PLATFORM_TYPE] === BARE_METAL_TF,
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

      this.state = {
        services: [],
        showLogs: null,
        tectonicProgress: 0,
        terraformProgress: 0,
        xhrError: null,
      };

      this.destroy = this.destroy.bind(this);
      this.retry = this.retry.bind(this);
      this.startOver = this.startOver.bind(this);
    }

    isOutputSame (terraform) {
      const oldOutput = _.get(this.props, 'terraform.output', '');
      const newOutput = _.get(terraform, 'output', '');
      // Output can be 10MB, so compare lengths first. Don't trust the JS engine to be smart.
      return newOutput.length === oldOutput.length && newOutput === oldOutput;
    }

    componentDidMount () {
      if (this.outputNode) {
        this.outputNode.scrollTop = this.outputNode.scrollHeight;
      }
    }

    updateStatus ({tectonic, terraform}) {
      if (terraform.action === 'apply') {
        const services = (tectonic.isEtcdSelfHosted ? [{key: 'etcd', name: 'Etcd'}] : []).concat([
          {key: 'kubernetes', name: 'Kubernetes'},
          {key: 'identity', name: 'Tectonic Identity'},
          {key: 'ingress', name: 'Tectonic Ingress Controller'},
          {key: 'console', name: 'Tectonic Console'},
          {key: 'tectonicSystem', name: 'other Tectonic services'},
        ]);
        this.setState({services});

        const tectonicSucceeded = services.filter(s => _.get(tectonic[s.key], 'success')).length;
        const tectonicProgress = tectonicSucceeded / services.length;

        // Don't let progress bars go backwards (e.g. if a service flips back to incomplete)
        if (tectonicProgress > this.state.tectonicProgress) {
          this.setState({tectonicProgress});
        }

        if (this.props.isAWS && (this.state.terraformProgress === 0 || !this.isOutputSame(terraform))) {
          const terraformProgress = estimateTerraformProgress(terraform);
          if (terraformProgress > this.state.terraformProgress) {
            this.setState({terraformProgress});
          }
        }
      }
    }

    componentWillMount () {
      this.updateStatus(this.props);
    }

    componentWillReceiveProps (nextProps) {
      this.updateStatus(nextProps);
    }

    componentWillUpdate ({terraform}) {
      if (this.isOutputSame(terraform) || this.state.showLogs === false) {
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

    startOver () {
      // eslint-disable-next-line no-alert
      if (window.config.devMode || window.confirm('Do you really want to start over?')) {
        window.reset();
      }
    }

    retry () {
      // eslint-disable-next-line no-alert
      if (window.config.devMode || window.confirm('Are you sure you want to re-run terraform apply?')) {
        this.setState({terraformProgress: 0});
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
      const {clusterName, commitState, isAWS, isBareMetal, terraform, tectonic, tectonicDomain} = this.props;
      const commitPhase = _.get(commitState, 'phase');
      const {action, tfError, output, statusMsg} = terraform;
      const state = this.state;
      const showLogs = state.showLogs === null ? statusMsg !== 'success' : state.showLogs;
      const isApply = action === 'apply';
      const isTFRunning = statusMsg === 'running';
      const isTFSuccess = !isTFRunning && !tfError;
      const isApplySuccess = isApply && isTFSuccess;
      const isDestroySuccess = action === 'destroy' && isTFSuccess;

      const saveLog = () => saveAs(new Blob([output], {type: 'text/plain'}), `tectonic-${clusterName}.log`);

      const consoleSubsteps = [];

      if (isApply && tectonic.console) {
        if (isAWS) {
          consoleSubsteps.push(<AWS_DomainValidation key="domain" />);
        }

        const dnsReady = tectonic.console.success || ((tectonic.console.message || '').search('no such host') === -1);
        consoleSubsteps.push(
          <Step pending={isTFRunning} done={dnsReady && !isTFRunning} cancel={tfError} key="dnsReady">
            Resolving <a href={`https://${tectonicDomain}`} target="_blank">{tectonicDomain}</a>
          </Step>
        );

        const anyFailed = state.services.some(s => tectonic[s.key].failed);
        const allDone = state.services.every(s => tectonic[s.key].success);

        consoleSubsteps.push(
          <Step pending={isTFRunning} done={allDone} error={anyFailed} cancel={tfError} key="tectonicReady">
            {anyFailed
              ? <span><span className="wiz-running-fg">Starting Tectonic:</span> [Failure] Tectonic Console Error</span>
              : 'Starting Tectonic'}
            <div style={{marginLeft: 22}}>
              {isApplySuccess && !allDone && <div>
                <ProgressBar progress={state.tectonicProgress} isActive={!anyFailed} />
                {state.services.map(({key, name}) => {
                  const {failed, message, success} = tectonic[key];
                  return <Step done={success} error={failed} key={key} substep={true}>
                    {failed ? <span><b>{name}</b>{message && `: ${message}`}</span> : name}
                  </Step>;
                })}
              </div>
              }
              {anyFailed &&
                <Alert severity="error" noIcon>
                  <h4>Tectonic Console Failed</h4>
                  <p>To debug, save logs and download "kubeconfig" for troubleshooting:</p>
                  <code>
                    $ export KUBECONFIG=/path/to/kubeconfig{'\n'}
                    {_(tectonic).filter('failed').map(({namespace, pod}) => `$ kubectl --namespace=${namespace} logs ${pod}\n`).value()}
                  </code>
                  <a href="/tectonic/kubeconfig" download>
                    <Btn isError={true} title="Download kubeconfig" />
                  </a>
                  <p>To start over, destroy cluster to clear anything that may have been created.</p>
                  <Btn isError={true} onClick={this.destroy} title="Destroy Cluster and Start Over" />
                </Alert>
              }
            </div>
          </Step>
        );
      }

      const tfTitle = `${isApply ? 'Applying' : 'Destroying'} Terraform`;

      return <div>
        {!isBareMetal &&
          <p>
            Kubernetes is starting up. We're committing your cluster details.
            Grab some tea and sit tight. This process can take up to 20 minutes.
            Status updates will appear below.
          </p>
        }
        {isBareMetal &&
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
          <div className="col-xs-12 wiz-launch-progress">
            <Step done={statusMsg === 'success'} error={tfError}>
              {tfError
                ? <span><span className="wiz-running-fg">{tfTitle}:</span> [Failure] Terraform {action} failed</span>
                : tfTitle}
              {output && !isApplySuccess &&
                <div className="pull-right" style={{fontSize: '13px'}}>
                  <a onClick={() => this.setState({showLogs: !showLogs})}>
                    {showLogs
                      ? <span><i className="fa fa-angle-up"></i>&nbsp;&nbsp;Hide logs</span>
                      : <span><i className="fa fa-angle-down"></i>&nbsp;&nbsp;Show logs</span>}
                  </a>
                  <span className="spacer"></span>
                  <a onClick={saveLog}>
                    <i className="fa fa-download"></i>&nbsp;&nbsp;Save logs
                  </a>
                </div>
              }
            </Step>
            <div style={{marginLeft: 22}}>
              {isAWS && isApply && statusMsg !== 'success' && <ProgressBar progress={state.terraformProgress} isActive={isTFRunning} />}
              {showLogs && output && !isApplySuccess &&
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
              {state.xhrError && <Alert severity="error">{state.xhrError}</Alert>}
              {commitPhase === commitPhases.FAILED && <Alert severity="error">{commitState.response}</Alert>}
              {tfError && <Alert severity="error">{tfError.toString()}</Alert>}
              {tfError && !isTFRunning &&
                <Alert severity="error" noIcon>
                  <b>{_.startCase(action)} Failed</b>. Your installation is blocked. To continue:
                  <ol style={{fontSize: 13, paddingLeft: 30, paddingTop: 10, paddingBottom: 10}}>
                    <li><a onClick={saveLog}>Save Terraform logs</a> and <a href="/terraform/assets" download>download assets</a> for debugging purposes.</li>
                    <li>Destroy your cluster to clear anything that may have been created. Or,</li>
                    <li>Reapply Terraform.</li>
                  </ol>
                  <Btn isError={true} onClick={this.destroy} title="Destroy Cluster and Start Over" />
                  <Btn isError={true} onClick={this.retry} title="Retry Terraform Apply" />
                </Alert>
              }
              {isDestroySuccess &&
                <Alert noIcon>
                  <b style={{fontWeight: '600'}}>Destroy Succeeded</b>
                  <p>To continue, make a fresh start with Tectonic Installer, or simply close the browser tab to quit.</p>
                  <Btn isError={false} onClick={this.startOver} title="Start Over" />
                  <Btn isError={false} onClick={this.retry} title="Retry Terraform Apply" />
                </Alert>
              }
              {isApplySuccess &&
                <div className="wiz-launch-progress__help">
                  You can save Terraform logs, or destroy your cluster if you changed your mind:&nbsp;
                  <DropdownInline
                    header="Terraform Actions"
                    items={[
                      ['Destroy Cluster', this.destroy],
                      ['Retry Terraform Apply', this.retry],
                      ['Save Terraform Logs', saveLog],
                    ]}
                  />
                </div>
              }
            </div>
            {consoleSubsteps}
          </div>
        </div>
        <br />
        {!isTFRunning && !isDestroySuccess &&
          <div className="row">
            <div className="col-xs-12">
              <a href="/terraform/assets" download>
                <button className={classNames('btn btn-primary wiz-giant-button')}>
                  <i className="fa fa-download"></i>&nbsp;&nbsp;Download assets
                </button>
              </a>
            </div>
          </div>
        }
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
