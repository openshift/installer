import _ from 'lodash';
import React from 'react';
import { connect } from 'react-redux';

import { PLATFORM_TYPE } from '../cluster-config';
import { Select, WithClusterConfig } from './ui';
import { TectonicGA } from '../tectonic-ga';

import {
  DOCS,
  PLATFORM_NAMES,
  SELECTED_PLATFORMS,
  isSupported,
  OptGroups,
} from '../platforms';

const SelectPlatform = props => {
  const options = [];
  _.each(OptGroups, optgroup => {
    const [name, ...group] = optgroup;
    const platforms = _.filter(group, p => SELECTED_PLATFORMS.includes(p));
    if (platforms.length) {
      options.push(<optgroup label={name} key={name}>{
        platforms.map(p => <option value={p} key={p}>{PLATFORM_NAMES[p]}</option>)
      }
      </optgroup>);
    }
  });
  return <Select id={PLATFORM_TYPE} {...props}
    onValue={(value) => {TectonicGA.sendEvent('Platform Changed', 'user input', value); props.onValue(value);}}>
    {options}
  </Select>;
};

export const ClusterType = connect(
  ({clusterConfig}) => {
    const platform = clusterConfig[PLATFORM_TYPE];
    return {
      platform,
      platformName: PLATFORM_NAMES[platform],
      supportedPlatform: isSupported(platform),
    };
  },
)(({platform, platformName, supportedPlatform}) => <div>
  <div className="row form-group">
    <div className="col-xs-12">
      Select an installation path from the options below.
    </div>
  </div>

  <div className="row form-group">
    <div className="col-xs-3">
      <label htmlFor={PLATFORM_TYPE}>
        Platform
      </label>
    </div>
    <div className="col-xs-9">
      <WithClusterConfig field={PLATFORM_TYPE}>
        <SelectPlatform />
      </WithClusterConfig>
      { supportedPlatform && <p className="text-muted">
        Use the graphical installer to input cluster details, this is best for demos and your first Tectonic cluster.
        &nbsp;&nbsp;<a href={DOCS[platform]} target="_blank" >{platformName} documentation&nbsp;&nbsp;<i className="fa fa-external-link" /></a>
      </p>
      }
      { !supportedPlatform && <div>
        <p>
          Use the documentation and the Terraform CLI to install a cluster with specific infrastructure use-cases.
          This method is designed for automation and doesn't use the graphical insaller.
        </p>
        <p>
          <a href={DOCS[platform]} target="_blank">
            <button className="btn btn-primary">{platformName.split("(Alpha)")[0]} Docs&nbsp;&nbsp;<i className="fa fa-external-link" /></button>
          </a>
        </p>
      </div>
      }
    </div>
  </div>
  <div className="row">
    <div className="col-xs-12">
      <p className="text-muted" style={{marginTop: 60}}>
        CoreOS collects data about your Tectonic cluster for billing purposes. See the <a href="https://coreos.com/data-policy/tectonic" onClick={TectonicGA.sendDocsEvent} target="_blank">data policy</a> for details.
       </p>
    </div>
  </div>
</div>);

ClusterType.canNavigateForward = ({clusterConfig}) => isSupported(clusterConfig[PLATFORM_TYPE]);
