import _ from 'lodash';

export const compose = (...validators) => {
  return s => {
    for (const v of validators) {
      const err = v(s);
      if (err) {
        return err;
      }
    }
  };
};

export const validate = {
  nonEmpty: function (s) {
    if (s && ('' + s).trim().length > 0) {
      return;
    }

    return 'This field is required. You must provide a value.';
  },

  certificate: (s) => {
    const trimmed = (s || '').trim();
    if (trimmed.length >= 55 && trimmed.match(/^-----BEGIN CERTIFICATE-----[^]*-----END CERTIFICATE-----$/)) {
      return;
    }

    return 'Invalid certificate. Check your certificate for copy/paste errors.';
  },

  privateKey: (s) => {
    const trimmed = (s || '').trim();
    if (trimmed.length >= 63 && trimmed.match(/^-----BEGIN [A-Z]{2,10} PRIVATE KEY-----[^]*-----END [A-Z]{2,10} PRIVATE KEY-----$/)) {
      return;
    }

    return 'Invalid private key. Check your key for copy/paste errors';
  },

  email: (s = '') => {
    const errMsg = validate.nonEmpty(s);
    if (errMsg) {
      return errMsg;
    }
    const [name, domain] = s.split('@');
    // No whitespace allowed
    if (validate.nonEmpty(name) || /\s/g.test(name) || !domain || validate.domainName(domain)) {
      return 'Invalid email address.';
    }
  },

  MAC: (s = '') => {
    if (!s.length) {
      return;
    }

    if (s.trim() !== s) {
      return 'Leading/trailing whitespace not allowed.';
    }

    // We want to accept everything that golang net.ParseMAC will
    // see https://golang.org/src/net/mac.go?s=1054:1106#L28
    const error = 'Invalid MAC address.';

    if (s.match(/^([a-fA-F0-9]{2}:)+([a-fA-F0-9]{2})$/)) {
      if (s.length === '01:23:45:67:89:ab'.length ||
        s.length === '01:23:45:67:89:ab:cd:ef'.length ||
        s.length === '01:23:45:67:89:ab:cd:ef:00:00:01:23:45:67:89:ab:cd:ef:00:00'.length) {
        return;
      }

      return error;
    }

    if (s.match(/^([a-fA-F0-9]{2}-)+([a-fA-F0-9]{2})$/)) {
      if (s.length === '01-23-45-67-89-ab'.length ||
        s.length === '01-23-45-67-89-ab-cd-ef'.length ||
        s.length === '01-23-45-67-89-ab-cd-ef-00-00-01-23-45-67-89-ab-cd-ef-00-00'.length) {
        return;
      }

      return error;
    }

    if (s.match(/^([a-fA-F0-9]{4}\.)+([a-fA-F0-9]{4})$/)) {
      if (s.length === '0123.4567.89ab'.length ||
        s.length === '0123.4567.89ab.cdef'.length ||
        s.length === '0123.4567.89ab.cdef.0000.0123.4567.89ab.cdef.0000'.length) {
        return error;
      }
    }

    return error;
  },

  IP: (s) => {
    // four octets of decimal numbers in the valid range. This allows
    // screwy IPs like 0.0.0.0 and 127.0.0.1
    const matched = (s || '').match(/^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})$/);
    if (matched && matched.slice(1).every(oct => parseInt(oct, 10) < 256)) {
      return;
    }

    return 'Invalid IP address.';
  },

  domainName: (s) => {
    const split = (s || '').split('.');
    if (split.slice(-1)[0] === '') {
      // Trailing dot is ok
      split.splice(-1, 1);
    }
    if (split.every(l => l.match(/^[a-zA-Z0-9-]{1,63}$/))) {
      return;
    }

    return 'Invalid domain name.';
  },

  host: (s) => {
    // either valid IP address or domain name
    if (!validate.IP(s) || !validate.domainName(s)) {
      return;
    }

    return 'Invalid format. You must provide a domain name or IP address.';
  },

  port: (s = '') => {
    const errMsg = 'Invalid port value. You must provide a valid port number.';
    if (!s.match(/^[0-9]+$/)) {
      return errMsg;
    }
    if (parseInt(s, 10) > 65535) {
      return errMsg;
    }
    return;
  },

  hostPort: (s) => {
    const [host, port] = (s || '').split(':', 2);
    if (!host || !port) {
      return 'Invalid format. You must use <host>:<port> format.';
    }

    if (validate.IP(host) && validate.domainName(host)) {
      return 'Invalid format. Host must be a domain name or an IP address.';
    }

    if (validate.port(port)) {
      return 'Invalid port value. You must provide a valid port number.';
    }

    return;
  },

  SSHKey: (s) => {
    const err = validate.nonEmpty(s);
    if (err) {
      return err;
    }
    // Don't let users hang themselves
    if (s.match(/-{5}BEGIN [\w-]+ PRIVATE KEY-{5}/)) {
      return 'Private key detected! Please paste your public key.';
    }
    const [pubKey, ...extraKeys] = _.trimEnd(s).split('\n');
    if (extraKeys.length) {
      return 'Invalid SSH pubkey. Did you paste multiple keys?';
    }
    const errMsg = 'Invalid SSH pubkey. Check your key for copy/paste errors.';
    const [type, base64] = pubKey.split(/[\s]+/);
    if (!type || !type.match(/^[\w-]+$/)) {
      return errMsg;
    }
    if (!base64 || !base64.match(/^[A-Za-z0-9+/]+={0,2}$/)) {
      return errMsg;
    }
    return;
  },

  schema: (schema) => {
    return (value) => {
      for (const k of Object.keys(schema)) {
        const validity = schema[k](value[k]);
        if (validity) {
          return validity;
        }
      }

      return;
    };
  },

  int: ({min, max}) => {
    return (value) => {
      if (!/^(?:[-+]?(?:0|[1-9][0-9]*))$/.test(value)) {
        return 'Invalid format. Please enter a number.';
      }

      if (min !== undefined && value < min) {
        return `Invalid value. Provide a value greater than ${min - 1}.`;
      }

      if (max !== undefined && value > max) {
        return `Invalid value. Provide a value less than ${max + 1}.`;
      }

      return;
    };
  },

  isOdd: v => {
    if (_.toNumber(v) % 2 !== 1) {
      return 'Invalid value. Must be odd.';
    }
  },

  CIDR: value => {
    let err = validate.nonEmpty(value);
    if (err) {
      return err;
    }
    const split = value.split('/');
    if (split.length === 1) {
      return 'Invalid value. You must provide a CIDR netmask (eg, /24).';
    }
    if (split.length !== 2) {
      return 'Invalid value.';
    }
    const ip = split[0];
    err = validate.IP(ip);
    if (err){
      return err;
    }
    let mask = split[1];
    // Digits only. Can't start with a zero (unless it is zero).
    if (!mask.match(/^[1-9]\d?$/) && mask !== '0') {
      return 'Invalid netmask size.';
    }
    mask = parseInt(mask, 10);
    if (isNaN(mask) || mask > 32 || mask < 0) {
      return 'Invalid netmask size. Must be 0-32.';
    }
    if (ip.startsWith('172.17.')) {
      return 'Overlaps with default Docker Bridge subnet (172.17.0.0/16)';
    }
    return;
  },

  AWSsubnetCIDR: value => {
    const err = validate.CIDR(value);
    if (err) {
      return err;
    }
    const mask = value.split('/')[1];
    if (mask > 28 || mask < 16) {
      // http://docs.aws.amazon.com/AmazonVPC/latest/UserGuide/VPC_Subnets.html#vpc-sizing-ipv4
      return 'AWS subnets must be between /16 and /28.';
    }
    return;
  },

  someSelected: (fields, deselectedFields) => {
    if (!deselectedFields) {
      return;
    }
    for (const f of fields) {
      if (!deselectedFields[f]) {
        return;
      }
    }
    return 'At least one field must be selected';
  },
};
