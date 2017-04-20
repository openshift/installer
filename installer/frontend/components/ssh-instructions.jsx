import React from 'react';

export const SSHInstructions = ({controllerIP}) => <div className="form-group">
  <p>
    Almost done! Next, run the commands below to start Tectonic. You may need to use ssh's -i flag to specify your ssh key. e.g. <code>ssh -i my_key.pem …</code>
  </p>
  <pre className="wiz-shell-example">{
`ssh core@${controllerIP} 'sudo systemctl start bootkube'`}
  </pre>
  <p>
    If Tectonic Console doesn't come up in 10 minutes, check the logs or contact support.
  </p>
  <pre className="wiz-shell-example">{`ssh core@${controllerIP} 'journalctl -u bootkube -f'`}
  </pre>
</div>;

