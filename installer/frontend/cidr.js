import { validate } from './validate';

// Get the number of IP addresses in a CIDR range
export const cidrSize = cidr => {
  if (validate.CIDR(cidr)) {
    return undefined;
  }
  const [, bits] = cidr.split('/');

  // JavaScript's bit shifting only works on signed 32bit ints so <<31 would be negative :(
  return Math.pow(2, 32 - parseInt(bits, 10));
};

// Convert an IPv4 string to an integer
const ip2int = ip => {
  const [a, b, c, d] = ip.split('.').map(i => parseInt(i, 10));
  return a * (256 ** 3) + b * (256 ** 2) + c * 256 + d;
};

// Convert an integer to an IPv4 string
const int2ip = int => [0, 1, 2, 3].map(i => Math.floor(int / (256 ** i)) % 256).reverse().join('.');

// Get the first IP address in a CIDR range
export const cidrStart = cidr => {
  if (validate.CIDR(cidr)) {
    return undefined;
  }
  const ip = cidr.split('/')[0];
  const size = cidrSize(cidr);
  return int2ip(size * Math.floor(ip2int(ip) / size));
};

// Get the last IP address in a CIDR range
export const cidrEnd = cidr => {
  if (validate.CIDR(cidr)) {
    return undefined;
  }
  return int2ip(ip2int(cidrStart(cidr)) + cidrSize(cidr) - 1);
};
