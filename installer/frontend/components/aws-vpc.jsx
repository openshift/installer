import _ from 'lodash';
import React from 'react';
import { connect } from 'react-redux';

import { compose, validate } from '../validate';
import * as awsActions from '../aws-actions';
import {
  AsyncSelect,
  Radio,
  Select,
  Deselect,
  DeselectField,
  WithClusterConfig,
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
import { CIDR } from './cidr';
import { Field, Form } from '../form';

const { setIn } = configActions;

const AWS_ADVANCED_NETWORKING = 'awsAdvancedNetworking';

import {
  AWS_CONTROLLER_SUBNET_IDS,
  AWS_CONTROLLER_SUBNETS,
  AWS_CREATE_VPC,
  AWS_HOSTED_ZONE_ID,
  AWS_REGION,
  AWS_REGION_FORM,
  AWS_SUBNETS,
  AWS_SPLIT_DNS,
  AWS_VPC_CIDR,
  AWS_VPC_ID,
  AWS_VPC_FORM,
  AWS_WORKER_SUBNET_IDS,
  AWS_WORKER_SUBNETS,
  CLUSTER_NAME,
  CLUSTER_SUBDOMAIN,
  DESELECTED_FIELDS,
  POD_CIDR,
  SERVICE_CIDR,
  toVPCSubnetID,
  SPLIT_DNS_ON,
  // SPLIT_DNS_OPTIONS,
} from '../cluster-config';

const vpcInfoForm = new Form(AWS_VPC_FORM, [
  new Field(AWS_SPLIT_DNS, {default: SPLIT_DNS_ON}),
  new Field(AWS_CREATE_VPC, {
    default: 'VPC_CREATE',
  }),
  new Field(CLUSTER_SUBDOMAIN, {
    default: '',
    validator: compose(validate.nonEmpty, validate.domainName),
  }),
  new Field(AWS_ADVANCED_NETWORKING, {
    default: false,
  }),
  new Field(AWS_HOSTED_ZONE_ID, {
    default: '',
    dependencies: [AWS_REGION_FORM],
    validator: (value, clusterConfig, oldValue, extraData) => {
      const empty = validate.nonEmpty(value);
      if (empty) {
        return empty;
      }

      if (!extraData || !extraData.zoneToName[value]) {
        return 'Unknown zone id.';
      }
    },
    getExtraStuff: (dispatch, isNow) => dispatch(awsActions.getZones(null, null, isNow))
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
], {
  validator: (data, cc) => {
    const hostedZoneID = data[AWS_HOSTED_ZONE_ID];
    const privateZone = _.get(cc, ['extra', AWS_HOSTED_ZONE_ID, 'privateZones', hostedZoneID]);
    if (!privateZone || !hostedZoneID) {
      return;
    }
    if (privateZone && data[AWS_CREATE_VPC] !== 'VPC_PRIVATE') {
      return 'Private Route 53 Zones must use an existing private VPC.';
    }
  },
});

const SubnetSelect = ({field, name, subnets, asyncValidator, disabled, fieldName}) => <div className="row form-group">
  <div className="col-xs-3">
    <Deselect field={fieldName} />
    <label htmlFor={`${DESELECTED_FIELDS}.${fieldName}`}>{name}</label>
  </div>
  <div className="col-xs-6">
    <WithClusterConfig field={field} asyncValidator={asyncValidator}>
      <Select disabled={disabled}>
        <option disabled>Select a subnet</option>
        {_.filter(subnets, ({availabilityZone}) => availabilityZone === name)
          .map(({id, instanceCIDR}) => <option value={id} key={instanceCIDR}>{instanceCIDR} ({id})</option>)
        }
      </Select>
    </WithClusterConfig>
  </div>
</div>;

/* ggreer: This component (along with its usage) is commented-out because it will soon be needed,
 * but not in this release. Matt & I didn't realize that until after we'd written it.
 * ESLint will fail if you have an unused variable, so we commented it out.
 */
// const IP = ({field, name, disabled, placeholder}) => <div className="row form-group">
//   <div className="col-xs-4">
//     <label htmlFor={field}>{name}</label>
//   </div>
//   <div className="col-xs-8">
//     <WithClusterConfig field={field} validator={validate.IP}>
//       <Input placeholder={placeholder} id={field} disabled={disabled} />
//     </WithClusterConfig>
//   </div>
// </div>;

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
    region: clusterConfig[AWS_REGION],
    awsWorkerSubnets: clusterConfig[AWS_WORKER_SUBNETS],
    awsWorkerSubnetIds: clusterConfig[AWS_WORKER_SUBNET_IDS],
    awsControllerSubnets: clusterConfig[AWS_CONTROLLER_SUBNETS],
    awsControllerSubnetIds: clusterConfig[AWS_CONTROLLER_SUBNET_IDS],
    awsCreateVpc: clusterConfig[AWS_CREATE_VPC] === 'VPC_CREATE',
    awsVpcId: clusterConfig[AWS_VPC_ID],
    awsVpcCIDR: clusterConfig[AWS_VPC_CIDR],
    clusterName: clusterConfig[CLUSTER_NAME],
    clusterSubdomain: clusterConfig[CLUSTER_SUBDOMAIN],
    internalCluster: clusterConfig[AWS_CREATE_VPC] === 'VPC_PRIVATE',
    podCIDR: clusterConfig[POD_CIDR],
    serviceCIDR: clusterConfig[SERVICE_CIDR],
    advanced: clusterConfig[AWS_ADVANCED_NETWORKING],
    privateZone: _.get(clusterConfig, ['extra', AWS_HOSTED_ZONE_ID, 'privateZones', clusterConfig[AWS_HOSTED_ZONE_ID]]),
  };
};

const dispatchToProps = dispatch => ({
  getVpcs: () => dispatch(awsActions.getVpcs()),
  getVpcSubnets: vpcID => dispatch(awsActions.getVpcSubnets({vpcID})),
  validate: body => dispatch(awsActions.validateSubnets(body)),
  reset: () => {
    setIn(AWS_CONTROLLER_SUBNET_IDS, {}, dispatch);
    setIn(AWS_WORKER_SUBNET_IDS, {}, dispatch);
  },
  getDefaultSubnets: () => dispatch(awsActions.getDefaultSubnets()),
  setSubdomain: subdomain => setIn(CLUSTER_SUBDOMAIN, subdomain, dispatch),
});

export const AWS_VPC = connect(stateToProps, dispatchToProps)(
  class AWS_VPCComponent extends React.Component {
    validateVPC () {
      const { awsCreateVpc, awsVpcCIDR, awsVpcId, region, internalCluster, serviceCIDR, podCIDR } = this.props;
      let controllerSubnets, privateSubnets;
      if (awsCreateVpc) {
        controllerSubnets = toVPCSubnetID(region, this.props.awsControllerSubnets);
        privateSubnets = toVPCSubnetID(region, this.props.awsWorkerSubnets);
      } else {
        if (!awsVpcId) {
          // User hasn't selected a VPC yet. Don't try to validate.
          return Promise.resolve();
        }
        controllerSubnets = toVPCSubnetID(region, this.props.awsControllerSubnetIds);
        privateSubnets = toVPCSubnetID(region, this.props.awsWorkerSubnetIds);
      }

      let publicSubnets;
      if (internalCluster) {
        publicSubnets = [];
        privateSubnets.push(...controllerSubnets);
        privateSubnets = _.uniqWith(privateSubnets, _.isEqual);
      } else {
        publicSubnets = controllerSubnets;
      }

      const network = { publicSubnets, privateSubnets, podCIDR, serviceCIDR };
      if (awsCreateVpc) {
        network.vpcCIDR = awsVpcCIDR;
      } else {
        network.awsVpcId = awsVpcId;
      }
      return this.props.validate(network).then(json => {
        if (!json.valid) {
          return Promise.reject(json.message);
        }
      });
    }

    componentDidMount () {
      if (_.size(this.props.awsControllerSubnets) && _.size(this.props.awsWorkerSubnets)) {
        return;
      }
      this.props.getDefaultSubnets();
    }

    render () {
      const { availableVpcs, awsCreateVpc, availableVpcSubnets, awsVpcId, clusterName, clusterSubdomain, internalCluster, advanced } = this.props;

      let controllerSubnets;
      let workerSubnets;
      if (awsCreateVpc) {
        controllerSubnets = _.map(this.props.awsControllerSubnets, (subnet, az) => {
          const fieldName = `${AWS_SUBNETS}.${az}`;
          return <DeselectField key={az} field={fieldName}>
            <CIDR field={`${AWS_CONTROLLER_SUBNETS}.${az}`} name={az} fieldName={fieldName} placeholder="10.0.0.0/24" autoFocus={az.endsWith('a')} />
          </DeselectField>;
        });
        workerSubnets = _.map(this.props.awsWorkerSubnets, (subnet, az) => {
          const fieldName = `${AWS_SUBNETS}.${az}`;
          return <DeselectField key={az} field={fieldName}>
            <CIDR field={`${AWS_WORKER_SUBNETS}.${az}`} name={az} fieldName={fieldName} placeholder="10.0.0.0/24" />
          </DeselectField>;
        });
      } else if (awsVpcId) {
        const availableControllerSubnets = internalCluster ? availableVpcSubnets.value.private : availableVpcSubnets.value.public;
        if (_.size(availableControllerSubnets)) {
          controllerSubnets = _.map(this.props.azs, az => {
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
          workerSubnets = _.map(this.props.azs, az => {
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
                    <Radio name={AWS_CREATE_VPC} value="VPC_CREATE" />
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
                    <Radio name={AWS_CREATE_VPC} value="VPC_PUBLIC" />
                  </Connect>
                Use an existing VPC (Public)
                </label>
                <p className="text-muted wiz-help-text">
                Useful for installing beside existing resources. Your VPC must be <a href="https://coreos.com/tectonic/docs/latest/install/aws/requirements.html#using-an-existing-vpc" onClick={() => TectonicGA.sendDocsEvent('aws-tf')} target="_blank">set up correctly</a>.
                </p>
              </div>
            </div>
            <div className="wiz-radio-group">
              <div className="radio wiz-radio-group__radio">
                <label>
                  <Connect field={AWS_CREATE_VPC}>
                    <Radio name={AWS_CREATE_VPC} value="VPC_PRIVATE" />
                  </Connect>
                Use an existing VPC (Private)
                </label>
                <p className="text-muted wiz-help-text">
                Useful for installing beside existing resources. Your VPC must be <a href="https://coreos.com/tectonic/docs/latest/install/aws/requirements.html#using-an-existing-vpc" onClick={() => TectonicGA.sendDocsEvent('aws-tf')} target="_blank">set up correctly</a>.
                </p>
              </div>
            </div>
          </div>
        </div>

        <hr />

        <p className="text-muted">
        Please select a Route 53 hosted zone. For more information, see AWS Route 53 docs on <a target="_blank" href="https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/AboutHZWorkingWith.html">Working with Hosted Zones</a>.
        </p>
        <div className="row form-group">
          <div className="col-xs-3">
            <label htmlFor="r53Zone">DNS</label>
          </div>
          <div className="col-xs-9">
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
        {/* privateZone &&
          <div className="row form-group">
            <div className="col-xs-offset-3 col-xs-9">
              <Connect field={AWS_SPLIT_DNS}>
                <Select>
                  {_.map(SPLIT_DNS_OPTIONS, ((k, v) => <option value={v} key={k}>{k}</option>))}
                </Select>
              </Connect>
              <p className="text-muted wiz-help-text">
                See AWS <a href="https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/hosted-zones-private.html"
                  target="_blank">Split-View DNS documentation&nbsp;<i className="fa fa-external-link" /></a>
              </p>
            </div>
          </div>
        */}

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
                  Specify a range of IPv4 addresses for the VPC in the form of a <a href="https://tools.ietf.org/html/rfc4632" target="_blank">CIDR block</a>. Safe defaults have been chosen for you.
                </div>
              </div>
              <CIDR name="CIDR block" field={AWS_VPC_CIDR} placeholder="10.0.0.0/16" />
            </div>
          }
          {!awsCreateVpc &&
            <div className="row">
              <div className="col-xs-3">
                <label htmlFor="r53Zone">VPC</label>
              </div>
              <div className="col-xs-9">
                <div className="radio wiz-radio-group__radio">
                  <WithClusterConfig field={AWS_VPC_ID} asyncValidator={() => this.validateVPC()}>
                    <AsyncSelect
                      id={AWS_VPC_ID}
                      availableValues={availableVpcs}
                      disabledValue="Please select a VPC"
                      onRefresh={() => {
                        this.props.getVpcs();
                        if (awsVpcId) {
                          this.props.getVpcSubnets(awsVpcId);
                        }
                      }}
                      onChange={vpcID => {
                        if (vpcID !== awsVpcId) {
                          this.props.reset();
                        }
                        this.props.getVpcSubnets(vpcID);
                      }}
                    />
                  </WithClusterConfig>
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
          <KubernetesCIDRs validator={validate.AWSsubnetCIDR} />
        </div>
        }
      </div>;
    }
  });

AWS_VPC.canNavigateForward = ({clusterConfig}) => {
  if (!vpcInfoForm.canNavigateForward({clusterConfig})) {
    return false;
  }

  if (clusterConfig[AWS_CREATE_VPC] === 'VPC_CREATE') {
    const workerSubnets = clusterConfig[AWS_WORKER_SUBNETS];
    const controllerSubnets = clusterConfig[AWS_CONTROLLER_SUBNETS];
    return !validate.AWSsubnetCIDR(clusterConfig[AWS_VPC_CIDR]) &&
           _.every(controllerSubnets, subnet => !validate.AWSsubnetCIDR(subnet)) &&
           _.every(workerSubnets, subnet => !validate.AWSsubnetCIDR(subnet)) &&
          !validate.AWSsubnetCIDR(clusterConfig[POD_CIDR]) &&
          !validate.AWSsubnetCIDR(clusterConfig[SERVICE_CIDR]) &&
          !validate.someSelected(_.keys(clusterConfig[AWS_CONTROLLER_SUBNETS]), clusterConfig[DESELECTED_FIELDS][AWS_SUBNETS]) &&
          !validate.someSelected(_.keys(clusterConfig[AWS_WORKER_SUBNETS]), clusterConfig[DESELECTED_FIELDS][AWS_SUBNETS]);
  }

  return _.size(clusterConfig[AWS_CONTROLLER_SUBNET_IDS]) > 0 &&
         _.size(clusterConfig[AWS_WORKER_SUBNET_IDS]) > 0 &&
         _.some(clusterConfig[AWS_CONTROLLER_SUBNET_IDS], id => !validate.nonEmpty(id)) &&
         _.some(clusterConfig[AWS_WORKER_SUBNET_IDS], id => !validate.nonEmpty(id)) &&
         !validate.someSelected(_.keys(clusterConfig[AWS_CONTROLLER_SUBNET_IDS]), clusterConfig[DESELECTED_FIELDS][AWS_SUBNETS]) &&
         !validate.someSelected(_.keys(clusterConfig[AWS_WORKER_SUBNET_IDS]), clusterConfig[DESELECTED_FIELDS][AWS_SUBNETS]) &&
         !validate.AWSsubnetCIDR(clusterConfig[POD_CIDR]) &&
         !validate.AWSsubnetCIDR(clusterConfig[SERVICE_CIDR]);
};
