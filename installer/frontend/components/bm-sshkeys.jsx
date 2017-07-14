import React from 'react';

import { validate } from '../validate';
import { BM_SSH_KEY, SSH_AUTHORIZED_KEY } from '../cluster-config';
import { Field, Form } from '../form';

import { Alert } from './alert';
import { Connect, FileArea } from './ui';

const keyPlaceholder = `Paste your SSH public key here. It's often found in ~/.ssh/id_rsa.pub.

ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDxL
...
you@example.com
`;

const sshKeyForm = new Form(BM_SSH_KEY, [
  new Field(SSH_AUTHORIZED_KEY, {
    default: '',
    validator: validate.SSHKey,
  }),
]);

export const BM_SSHKeys = () => <div>
  <div className="form-group">
    Access to the nodes is intended for use by admins.
    End users can run applications using the API, CLI, and Tectonic
    Console. They typically don’t need SSH access.
  </div>
  <Alert severity="info">
    The public key below will be added to all machines in this cluster.
    This key must be on this machine's ssh-agent, as it will be used when configuring your nodes.
  </Alert>
  <div className="row form-group">
    <div className="col-xs-3">
      <label htmlFor="ssh-key">Public Key</label>
    </div>
    <div className="col-xs-9">
      <div className="wiz-ssh-key-container">
        <Connect field={SSH_AUTHORIZED_KEY}>
          <FileArea
            className="wiz-ssh-key-container__input"
            placeholder={keyPlaceholder}
            autoFocus />
        </Connect>
      </div>
    </div>
  </div>
  <sshKeyForm.Errors />
</div>;

BM_SSHKeys.canNavigateForward = sshKeyForm.canNavigateForward;
