import _ from 'lodash';
import React from 'react';
import { connect } from 'react-redux';

import { validate } from '../validate';
import { Connect, DocsA, ExternalLinkIcon, Select } from './ui';
import { Field, Form } from '../form';
import { PLATFORM_TYPE, PLATFORM_FORM } from '../cluster-config';
import { TectonicGA } from '../tectonic-ga';
import { AWS_TF, BARE_METAL_TF, DOCS, PLATFORM_NAMES, isSupported, optGroups } from '../platforms';

const ErrorComponent = connect(({clusterConfig}) => ({platform: clusterConfig[PLATFORM_TYPE]}))(
  ({error, platform}) => {
    const platformName = PLATFORM_NAMES[platform];
    if (error) {
      return <p>
        Use the documentation and the Terraform CLI to install a cluster with specific infrastructure use-cases.
        This method is designed for automation and doesn't use the graphical installer.
        <br />
        <DocsA path={DOCS[platform]}>
          <button className="btn btn-primary" style={{marginTop: 8}}>{platformName && platformName.split('(Alpha)')[0]} Docs<ExternalLinkIcon /></button>
        </DocsA>
      </p>;
    }
    return <p className="text-muted">
      Use the graphical installer to input cluster details, this is best for demos and your first Tectonic cluster.
      &nbsp;&nbsp;{platform === BARE_METAL_TF
        ? <span><br />{platformName} <DocsA path="/install/bare-metal/requirements.html">requirements<ExternalLinkIcon /></DocsA> and <DocsA path={DOCS[platform]}>install guide<ExternalLinkIcon /></DocsA>.</span>
        : <DocsA path={DOCS[platform]}>{platformName} documentation<ExternalLinkIcon /></DocsA>}
    </p>;
  });

const isEnabled = platform => _.get(window.config, 'platforms', []).includes(platform);

const platformForm = new Form(PLATFORM_FORM, [
  new Field(PLATFORM_TYPE, {
    default: _.find([AWS_TF, BARE_METAL_TF], isEnabled) || _.get(window.config, 'platforms[0]'),
    validator: validate.nonEmpty,
  }),
], {
  validator: (data, cc) => {
    if (!isSupported(cc[PLATFORM_TYPE])) {
      return 'unused';
    }
  },
});

const platformOptions = [];
_.each(optGroups, optgroup => {
  const [name, ...group] = optgroup;
  const platforms = _.filter(group, isEnabled);
  if (platforms.length) {
    platformOptions.push(<optgroup label={name} key={name}>{
      platforms.map(p => <option value={p} key={p}>{PLATFORM_NAMES[p]}</option>)
    }
    </optgroup>);
  }
});

export const ClusterType = () => <div>
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
      <Connect field={PLATFORM_TYPE}>
        <Select onValue={(value) => TectonicGA.sendEvent('Platform Changed', 'user input', value, value)}>
          {platformOptions}
        </Select>
      </Connect>
      <platformForm.Errors ErrorComponent={ErrorComponent} />
    </div>
  </div>
</div>;

ClusterType.canNavigateForward = platformForm.canNavigateForward;
