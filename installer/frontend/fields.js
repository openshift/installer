import { Field } from './form';
import { CLUSTER_NAME, PLATFORM_TYPE } from './cluster-config';
import { AWS_TF } from './platforms';

const FIELDS = {};
export default FIELDS;

FIELDS[CLUSTER_NAME] = new Field(CLUSTER_NAME, {
  default: '',
  validator: (s = '', CC) => {
    switch (CC[PLATFORM_TYPE]) {
    case AWS_TF:
      // http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-using-console-create-stack-parameters.html
      // http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-elasticloadbalancingv2-loadbalancer.html#cfn-elasticloadbalancingv2-loadbalancer-name
      if (s.length === 0 || s.length > 28) {
        return 'Value must be between 1 and 28 characters';
      }
      if (s.toLowerCase() !== s) {
        return 'Value must be lower case.';
      }
      if (!/^[a-zA-Z][-a-zA-Z0-9]*$/.test(s)) {
        return 'Value must be a valid AWS Stack Name: [a-zA-Z][-a-zA-Z0-9]*';
      }
      if (s.endsWith('-')) {
        return 'Value must not end with -';
      }
      break;
    default:
      if (s.length === 0 || s.length > 253) {
        return 'Value must be between 1 and 253 characters';
      }

      if (s.toLowerCase() !== s) {
        return 'Value must be lower case.';
      }

      // TODO: (ggreer) this falsely accepts "blah.-blah"
      // lower case alphanumeric characters, -, and ., with alpha beginning and end
      if (!/^([a-z0-9][a-z0-9.-]*)?[a-z0-9]$/.test(s)) {
        return 'Value must be alphanumeric [a-z0-9.-], beginning & ending with alphanumeric. Please refer to http://kubernetes.io/docs/user-guide/identifiers/.';
      }

      for (const t of s.split('.')) {
        // each segment no more than 63 characters
        if (t.length > 63) {
          return 'No segment between dots can be more than 63 characters.';
        }
      }
    }
  },
});
