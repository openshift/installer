import _ from 'lodash';
import React from 'react';
import { connect } from 'react-redux';

import { compose, validate } from '../validate';
import { getDefaultSubnets, getZones, getVpcs, getVpcSubnets, validateSubnets } from '../aws-actions';
import {
  AsyncSelect,
  Radio,
  Select,
  Deselect,
  DeselectField,
  Connect,
  Input,
  Selector,
  ToggleButton,
} from './ui';
import { Alert } from './alert';
import { configActions } from '../actions';
import { AWS_DomainValidation } from './aws-domain-validation';
import { TectonicGA } from '../tectonic-ga';
import { KubernetesCIDRs } from './k8s-cidrs';
import { CIDRRow } from './cidr';
import { Field, Form } from '../form';
import { toError } from '../utils';

import {
  AWS_CONTROLLER_SUBNETS,
  AWS_CONTROLLER_SUBNET_IDS,
  AWS_CREATE_VPC,
  AWS_HOSTED_ZONE_ID,
  AWS_REGION,
  AWS_REGION_FORM,
  AWS_SPLIT_DNS,
  AWS_SUBNETS,
  AWS_VPC_CIDR,
  AWS_VPC_FORM,
  AWS_VPC_ID,
  AWS_WORKER_SUBNETS,
  AWS_WORKER_SUBNET_IDS,
  CLUSTER_NAME,
  CLUSTER_SUBDOMAIN,
  DESELECTED_FIELDS,
  POD_CIDR,
  SERVICE_CIDR,
  SPLIT_DNS_ON,
  SPLIT_DNS_OPTIONS,
  VPC_CREATE ,
  VPC_PRIVATE,
  VPC_PUBLIC,
  getAwsZoneDomain,
  selectedSubnets,
} from '../cluster-config';

const AWS_ADVANCED_NETWORKING = 'awsAdvancedNetworking';
const DEFAULT_AWS_VPC_CIDR = '10.0.0.0/16';

const {setIn} = configActions;

const validateVPC = async (data, cc, updatedId, dispatch) => {
  const hostedZoneID = data[AWS_HOSTED_ZONE_ID];
  const privateZone = _.get(cc, ['extra', AWS_HOSTED_ZONE_ID, 'privateZones', hostedZoneID]);
  if (privateZone && hostedZoneID && data[AWS_CREATE_VPC] !== VPC_PRIVATE) {
    return 'Private Route 53 Zones must use an existing private VPC.';
  }

  const isCreate = cc[AWS_CREATE_VPC] === VPC_CREATE;
  const awsVpcId = cc[AWS_VPC_ID];

  if (!isCreate && !awsVpcId) {
    // User hasn't selected a VPC yet. Don't try to validate.
    return;
  }

  // Prevent unnecessary calls to validate API by only continuing if a field relevant to VPC validation has changed.
  // However, we always continue if updatedId is undefined, since that probably means the form has just been loaded and
  // we are doing initial validation.
  if (updatedId && ![
    AWS_CONTROLLER_SUBNETS,
    AWS_CONTROLLER_SUBNET_IDS,
    AWS_CREATE_VPC,
    AWS_REGION,
    AWS_VPC_CIDR,
    AWS_VPC_ID,
    AWS_WORKER_SUBNETS,
    AWS_WORKER_SUBNET_IDS,
    DESELECTED_FIELDS,
    POD_CIDR,
    SERVICE_CIDR,
  ].includes(updatedId)) {
    return _.get(cc, toError(AWS_VPC_FORM));
  }

  const getSubnets = subnets => {
    const selected = selectedSubnets(cc, subnets);
    return _.map(selected, (v, k) => isCreate ? {availabilityZone: k, instanceCIDR: v} : {availabilityZone: k, id: v});
  };

  const controllerSubnets = getSubnets(cc[isCreate ? AWS_CONTROLLER_SUBNETS : AWS_CONTROLLER_SUBNET_IDS]);
  const workerSubnets = getSubnets(cc[isCreate ? AWS_WORKER_SUBNETS : AWS_WORKER_SUBNET_IDS]);

  const isPrivate = cc[AWS_CREATE_VPC] === VPC_PRIVATE;
  const network = {
    privateSubnets: isPrivate ? _.uniqWith([...controllerSubnets, ...workerSubnets], _.isEqual) : workerSubnets,
    publicSubnets: isPrivate ? [] : controllerSubnets,
    podCIDR: cc[POD_CIDR],
    serviceCIDR: cc[SERVICE_CIDR],
  };

  // No need to validate if no subnets have been set yet
  if (_.isEmpty(network.privateSubnets) && _.isEmpty(network.publicSubnets)) {
    return;
  }

  if (isCreate) {
    network.vpcCIDR = cc[AWS_VPC_CIDR];
  } else {
    network.awsVpcId = awsVpcId;
  }

  try {
    return await dispatch(validateSubnets(network))
      .then(json => json.valid ? undefined : json.message);
  } catch (e) {
    return e.message || e.toString();
  }
};

const vpcInfoForm = new Form(AWS_VPC_FORM, [
  new Field(AWS_ADVANCED_NETWORKING, {default: false}),
  new Field(AWS_CONTROLLER_SUBNETS, {
    default: {},
    dependencies: [AWS_REGION_FORM],
    getExtraStuff: (dispatch, isNow) => dispatch(getDefaultSubnets(null, null, isNow)),
  }),
  new Field(AWS_CONTROLLER_SUBNET_IDS, {default: {}}),
  new Field(AWS_CREATE_VPC, {default: VPC_CREATE}),
  new Field(AWS_HOSTED_ZONE_ID, {
    default: '',
    dependencies: [AWS_REGION_FORM],
    validator: (value, cc) => {
      const empty = validate.nonEmpty(value);
      if (empty) {
        return empty;
      }
      if (!getAwsZoneDomain(cc)) {
        return 'Unknown zone ID.';
      }
    },
    getExtraStuff: (dispatch, isNow) => dispatch(getZones(null, null, isNow))
      .then(zones => {
        if (!isNow()) {
          return;
        }
        const zoneToName = {};
        const privateZones = {};
        const options = zones.map(({label, value}) => {
          const id = value.split('/hostedzone/')[1];
          if (label.endsWith('(private)')) {
            privateZones[id] = true;
          }
          zoneToName[id] = label.split(' ')[0];
          return {
            label,
            value: id,
          };
        });
        return {options: _.sortBy(options, 'label'), zoneToName, privateZones};
      }),
  }),
  new Field(AWS_SPLIT_DNS, {default: SPLIT_DNS_ON}),
  new Field(AWS_VPC_CIDR, {default: DEFAULT_AWS_VPC_CIDR, validator: validate.AWSsubnetCIDR}),
  new Field(AWS_VPC_ID, {default: '', ignoreWhen: cc => cc[AWS_CREATE_VPC] === VPC_CREATE}),
  new Field(AWS_WORKER_SUBNETS, {default: {}}),
  new Field(AWS_WORKER_SUBNET_IDS, {default: {}}),
  new Field(CLUSTER_SUBDOMAIN, {default: '', validator: compose(validate.nonEmpty, validate.domainName)}),
  new Field(DESELECTED_FIELDS, {default: {}}),
], {
  dependencies: [POD_CIDR, SERVICE_CIDR],
  validator: validateVPC,
});

const SubnetSelect = ({field, name, subnets, disabled, fieldName}) => <div className="row form-group">
  <div className="col-xs-4">
    <Deselect field={fieldName} />
    <label htmlFor={`${DESELECTED_FIELDS}.${fieldName}`}>{name}</label>
  </div>
  <div className="col-xs-6">
    <Connect field={field}>
      <Select disabled={disabled}>
        <option disabled value="">Select a subnet</option>
        {_.filter(subnets, ({availabilityZone}) => availabilityZone === name)
          .map(({id, instanceCIDR}) => <option value={id} key={instanceCIDR}>{instanceCIDR} ({id})</option>)
        }
      </Select>
    </Connect>
  </div>
</div>;

const stateToProps = ({aws, clusterConfig}) => {
  // populate subnet selection with all available azs ... many to many :(
  const azs = new Set();
  const availableVpcSubnets = aws.availableVpcSubnets.value;
  _.each(availableVpcSubnets.public, v => {
    azs.add(v.availabilityZone);
  });
  _.each(availableVpcSubnets.private, v => {
    azs.add(v.availabilityZone);
  });

  return {
    azs: new Array(...azs).sort(),
    availableVpcs: aws.availableVpcs,
    availableVpcSubnets: aws.availableVpcSubnets,
    awsWorkerSubnets: clusterConfig[AWS_WORKER_SUBNETS],
    awsControllerSubnets: clusterConfig[AWS_CONTROLLER_SUBNETS],
    awsCreateVpc: clusterConfig[AWS_CREATE_VPC] === VPC_CREATE,
    awsVpcId: clusterConfig[AWS_VPC_ID],
    clusterName: clusterConfig[CLUSTER_NAME],
    clusterSubdomain: clusterConfig[CLUSTER_SUBDOMAIN],
    internalCluster: clusterConfig[AWS_CREATE_VPC] === VPC_PRIVATE,
    advanced: clusterConfig[AWS_ADVANCED_NETWORKING],
  };
};

const dispatchToProps = dispatch => ({
  clearControllerSubnets: () => setIn(AWS_CONTROLLER_SUBNET_IDS, {}, dispatch),
  clearWorkerSubnets: () => setIn(AWS_WORKER_SUBNET_IDS, {}, dispatch),
  getVpcSubnets: vpcID => dispatch(getVpcSubnets({vpcID})),
  getVpcs: () => dispatch(getVpcs()),
});

export const AWS_VPC = connect(stateToProps, dispatchToProps)(props => {
  const {availableVpcs, awsCreateVpc, availableVpcSubnets, awsVpcId, clusterName, clusterSubdomain, internalCluster, advanced} = props;

  let controllerSubnets;
  let workerSubnets;
  if (awsCreateVpc) {
    controllerSubnets = _.map(props.awsControllerSubnets, (subnet, az) => {
      const fieldName = `${AWS_SUBNETS}.${az}`;
      return <DeselectField key={az} field={fieldName}>
        <CIDRRow
          field={`${AWS_CONTROLLER_SUBNETS}.${az}`}
          fieldName={fieldName}
          name={az}
          placeholder="10.0.0.0/24"
          validator={validate.AWSsubnetCIDR}
        />
      </DeselectField>;
    });
    workerSubnets = _.map(props.awsWorkerSubnets, (subnet, az) => {
      const fieldName = `${AWS_SUBNETS}.${az}`;
      return <DeselectField key={az} field={fieldName}>
        <CIDRRow
          field={`${AWS_WORKER_SUBNETS}.${az}`}
          fieldName={fieldName}
          name={az}
          placeholder="10.0.0.0/24"
          validator={validate.AWSsubnetCIDR}
        />
      </DeselectField>;
    });
  } else if (awsVpcId) {
    const availableControllerSubnets = internalCluster ? availableVpcSubnets.value.private : availableVpcSubnets.value.public;
    if (_.size(availableControllerSubnets)) {
      controllerSubnets = _.map(props.azs, az => {
        const fieldName = `${AWS_SUBNETS}.${az}`;
        return <DeselectField key={az} field={fieldName}>
          <SubnetSelect
            field={`${AWS_CONTROLLER_SUBNET_IDS}.${az}`}
            name={az}
            fieldName={fieldName}
            key={az}
            subnets={availableControllerSubnets}
          />
        </DeselectField>;
      });
    } else if (!availableVpcSubnets.inFly) {
      controllerSubnets = <Alert>{awsVpcId} has no {internalCluster ? 'private' : 'public'} subnets. Please create some using the AWS console.</Alert>;
    }
    if (_.size(availableVpcSubnets.value.private)) {
      workerSubnets = _.map(props.azs, az => {
        const fieldName = `${AWS_SUBNETS}.${az}`;
        return <DeselectField key={az} field={fieldName}>
          <SubnetSelect
            field={`${AWS_WORKER_SUBNET_IDS}.${az}`}
            name={az}
            fieldName={fieldName}
            key={az}
            subnets={availableVpcSubnets.value.private}
          />
        </DeselectField>;
      });
    } else if (!availableVpcSubnets.inFly) {
      workerSubnets = <Alert>{awsVpcId} has no private subnets. Please create some using the AWS console.</Alert>;
    }
  }

  return <div>
    <div className="row form-group">
      <div className="col-xs-12">
        <div className="wiz-radio-group">
          <div className="radio wiz-radio-group__radio">
            <label>
              <Connect field={AWS_CREATE_VPC}>
                <Radio name={AWS_CREATE_VPC} value={VPC_CREATE} />
              </Connect>
              Create a new VPC (Public)
            </label>&nbsp;(default)
            <p className="text-muted wiz-help-text">Launch into a new VPC with subnet defaults.</p>
          </div>
        </div>
        <div className="wiz-radio-group">
          <div className="radio wiz-radio-group__radio">
            <label>
              <Connect field={AWS_CREATE_VPC}>
                <Radio name={AWS_CREATE_VPC} value={VPC_PUBLIC} onChange={() => props.clearControllerSubnets()} />
              </Connect>
              Use an existing VPC (Public)
            </label>
            <p className="text-muted wiz-help-text">
              {/* eslint-disable react/jsx-no-target-blank */}
              Useful for installing beside existing resources. Your VPC must be <a href="https://coreos.com/tectonic/docs/latest/install/aws/requirements.html#using-an-existing-vpc" onClick={() => TectonicGA.sendDocsEvent('aws-tf')} rel="noopener" target="_blank">set up correctly</a>.
              {/* eslint-enable react/jsx-no-target-blank */}
            </p>
          </div>
        </div>
        <div className="wiz-radio-group">
          <div className="radio wiz-radio-group__radio">
            <label>
              <Connect field={AWS_CREATE_VPC}>
                <Radio name={AWS_CREATE_VPC} value={VPC_PRIVATE} onChange={() => props.clearControllerSubnets()} />
              </Connect>
              Use an existing VPC (Private)
            </label>
            <p className="text-muted wiz-help-text">
              {/* eslint-disable react/jsx-no-target-blank */}
              Useful for installing beside existing resources. Your VPC must be <a href="https://coreos.com/tectonic/docs/latest/install/aws/requirements.html#using-an-existing-vpc" onClick={() => TectonicGA.sendDocsEvent('aws-tf')} rel="noopener" target="_blank">set up correctly</a>.
              {/* eslint-enable react/jsx-no-target-blank */}
            </p>
          </div>
        </div>
      </div>
    </div>

    <hr />

    <p className="text-muted">
      Please select a Route 53 hosted zone. For more information, see AWS Route 53 docs on <a target="_blank" href="https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/AboutHZWorkingWith.html" rel="noopener noreferrer">Working with Hosted Zones</a>.
    </p>
    <div className="row form-group">
      <div className="col-xs-2">
        <label htmlFor="r53Zone">DNS</label>
      </div>
      <div className="col-xs-10">
        <div className="row">
          <div className="col-xs-4" style={{paddingRight: 0}}>
            <Connect field={CLUSTER_SUBDOMAIN} getDefault={() => clusterSubdomain || clusterName}>
              <Input placeholder="subdomain" />
            </Connect>
          </div>
          <div className="col-xs-8">
            <Connect field={AWS_HOSTED_ZONE_ID}>
              <Selector refreshBtn={true} disabledValue="Please select domain" />
            </Connect>
          </div>
        </div>
      </div>
    </div>
    {!internalCluster &&
      <div className="row form-group">
        <div className="col-xs-offset-2 col-xs-10">
          <Connect field={AWS_SPLIT_DNS}>
            <Select>
              {_.map(SPLIT_DNS_OPTIONS, ((k, v) => <option value={v} key={k}>{k}</option>))}
            </Select>
          </Connect>
          <p className="text-muted wiz-help-text">
            See AWS <a href="https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/hosted-zones-private.html" rel="noopener noreferrer" target="_blank">Split-View DNS documentation&nbsp;<i className="fa fa-external-link" /></a>
          </p>
        </div>
      </div>
    }

    <vpcInfoForm.Errors />
    <AWS_DomainValidation />
    <hr />

    {awsCreateVpc && <Connect field={AWS_ADVANCED_NETWORKING}>
      <ToggleButton className="btn btn-default">Advanced Settings</ToggleButton>
    </Connect>
    }
    {(advanced || !awsCreateVpc) && <div>
      {internalCluster && <Alert>
        You must be on a VPN with access to the target VPC. The cluster will have no public endpoints.
      </Alert>}

      {awsCreateVpc &&
        <div>
          <br />
          <Alert>
            The installer will create your EC2 instances within the following CIDR ranges.
            <br /><br />
            Safe defaults have been chosen for you.
            If you make changes, the ranges must not overlap and subnets must be within the VPC CIDR.
          </Alert>
          <div className="row form-group">
            <div className="col-xs-12">
              Specify a range of IPv4 addresses for the VPC in the form of a <a href="https://tools.ietf.org/html/rfc4632" rel="noopener noreferrer" target="_blank">CIDR block</a>. Safe defaults have been chosen for you.
            </div>
          </div>
          <CIDRRow name="CIDR Block" field={AWS_VPC_CIDR} placeholder={DEFAULT_AWS_VPC_CIDR} />
        </div>
      }
      {!awsCreateVpc &&
        <div className="row">
          <div className="col-xs-3">
            <label htmlFor="r53Zone">VPC</label>
          </div>
          <div className="col-xs-9">
            <div className="radio wiz-radio-group__radio">
              <Connect field={AWS_VPC_ID}>
                <AsyncSelect
                  id={AWS_VPC_ID}
                  availableValues={availableVpcs}
                  disabledValue="Please select a VPC"
                  onRefresh={() => {
                    props.getVpcs();
                    if (awsVpcId) {
                      props.getVpcSubnets(awsVpcId);
                    }
                  }}
                  onChange={vpcID => {
                    if (vpcID !== awsVpcId) {
                      props.clearControllerSubnets();
                      props.clearWorkerSubnets();
                    }
                    props.getVpcSubnets(vpcID);
                  }}
                />
              </Connect>
            </div>
          </div>
        </div>
      }

      {(controllerSubnets || workerSubnets) && <hr />}
      {controllerSubnets && <div className="row form-group">
        <div className="col-xs-12">
          <h4>Masters</h4>
          {controllerSubnets}
        </div>
      </div>
      }
      {workerSubnets && <div className="row form-group">
        <div className="col-xs-12">
          <h4>Workers</h4>
          {workerSubnets}
        </div>
      </div>
      }
      <hr />
      <KubernetesCIDRs autoFocus={false} />
    </div>
    }
  </div>;
});

AWS_VPC.canNavigateForward = ({clusterConfig}) => {
  if (!vpcInfoForm.canNavigateForward({clusterConfig}) || !KubernetesCIDRs.canNavigateForward({clusterConfig})) {
    return false;
  }

  const deselectedSubnets = clusterConfig[DESELECTED_FIELDS][AWS_SUBNETS];
  const isSelected = field => !validate.someSelected(_.keys(clusterConfig[field]), deselectedSubnets);

  if (clusterConfig[AWS_CREATE_VPC] === VPC_CREATE) {
    // The subnet CIDR fields are dynamically generated, so their validators won't automatically invalidate the form
    return _.every(clusterConfig[AWS_CONTROLLER_SUBNETS], subnet => !validate.AWSsubnetCIDR(subnet)) &&
      _.every(clusterConfig[AWS_WORKER_SUBNETS], subnet => !validate.AWSsubnetCIDR(subnet)) &&
      isSelected(AWS_CONTROLLER_SUBNETS) &&
      isSelected(AWS_WORKER_SUBNETS);
  }
  return isSelected(AWS_CONTROLLER_SUBNET_IDS) && isSelected(AWS_WORKER_SUBNET_IDS);
};
