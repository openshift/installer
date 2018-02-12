import classNames from 'classnames';
import React from 'react';
import { connect } from 'react-redux';

import { commitToServer } from '../server';
import { commitPhases } from '../actions';
import { withNav } from '../nav';
import { DocsA } from './ui';

export const SubmitDefinition = withNav(connect(
  state => ({
    phase: state.commitState.phase,
    response: state.commitState.response,
    ready: state.cluster.ready,
  }),
  {onFinish: commitToServer}
)(({navPrevious, phase, response, ready, onFinish}) => {
  const inProgress = (phase === commitPhases.REQUESTED ||
                      phase === commitPhases.WAITING ||
                      phase === commitPhases.SUCCEEDED && !ready);

  let btn = <button className="btn btn-primary wiz-giant-button" onClick={() => onFinish(false)}>Submit</button>;
  if (inProgress) {
    btn = <button className="btn btn-primary wiz-giant-button disabled">
      <i className="fa fa-spin fa-circle-o-notch"></i> Checking Your Definition...
    </button>;
  }

  const errorMessage = response ? response.toString() : '';
  const errorClasses = classNames('wiz-error-message', {
    hidden: phase !== commitPhases.FAILED,
  });

  return (
    <div>
      <p>
        Congratulations! Your cluster has been defined and will be submitted to Terraform. After submission, the definition cannot be updated. Go <a onClick={inProgress ? undefined : navPrevious} className={inProgress ? 'disabled' : undefined}>back</a> to update or make changes.
      </p>
      <p>
        You'll be able to download your <DocsA path="/admin/assets-zip.html">assets zip file</DocsA> after the definition is submitted.
      </p>
      <br />
      <div className="wiz-giant-button-container">
        {btn}
      </div>
      <p>
        <b>Advanced mode: </b>
        <a id="manualBoot" onClick={() => onFinish(true)}>Manually boot</a> your own cluster. Validate configuration, generate assets, but don't create the cluster.
      </p>
      <div className={errorClasses}>
        {errorMessage}
      </div>
    </div>
  );
}));
