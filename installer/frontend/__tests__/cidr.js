/* eslint-env jest, node */

import { cidrEnd, cidrSize, cidrStart } from '../cidr';

/* eslint-disable no-multi-spaces */
const tests = [
  {cidr: '1.2.3.4/0',   size: 2 ** 32,   start: '0.0.0.0', end: '255.255.255.255'},
  {cidr: '1.2.3.4/1',   size: 2 ** 31,   start: '0.0.0.0', end: '127.255.255.255'},
  {cidr: '1.2.3.4/2',   size: 2 ** 30,   start: '0.0.0.0', end: '63.255.255.255'},
  {cidr: '1.2.3.4/8',   size: 2 ** 24,   start: '1.0.0.0', end: '1.255.255.255'},
  {cidr: '1.2.3.4/16',  size: 2 ** 16,   start: '1.2.0.0', end: '1.2.255.255'},
  {cidr: '1.2.3.4/24',  size: 2 ** 8,    start: '1.2.3.0', end: '1.2.3.255'},
  {cidr: '1.2.3.4/29',  size: 8,         start: '1.2.3.0', end: '1.2.3.7'},
  {cidr: '1.2.3.4/30',  size: 4,         start: '1.2.3.4', end: '1.2.3.7'},
  {cidr: '1.2.3.4/31',  size: 2,         start: '1.2.3.4', end: '1.2.3.5'},
  {cidr: '1.2.3.4/32',  size: 1,         start: '1.2.3.4', end: '1.2.3.4'},
  {cidr: '1.2.3.4/33',  size: undefined, start: undefined, end: undefined},
  {cidr: '1.2.3.4/-1',  size: undefined, start: undefined, end: undefined},
  {cidr: '1.2.3.4/abc', size: undefined, start: undefined, end: undefined},

  {cidr: '0.0.0.0/1',         size: 2 ** 31, start: '0.0.0.0',   end: '127.255.255.255'},
  {cidr: '127.255.255.255/1', size: 2 ** 31, start: '0.0.0.0',   end: '127.255.255.255'},
  {cidr: '128.0.0.0/1',       size: 2 ** 31, start: '128.0.0.0', end: '255.255.255.255'},
  {cidr: '255.255.255.255/1', size: 2 ** 31, start: '128.0.0.0', end: '255.255.255.255'},
];
/* eslint-enable no-multi-spaces */

describe('CIDR helpers', () => {
  tests.forEach(t => {
    it(`Get size for ${t.cidr}`, () => expect(cidrSize(t.cidr)).toBe(t.size));
  });

  tests.forEach(t => {
    it(`Get start address for ${t.cidr}`, () => expect(cidrStart(t.cidr)).toBe(t.start));
  });

  tests.forEach(t => {
    it(`Get end address for ${t.cidr}`, () => expect(cidrEnd(t.cidr)).toBe(t.end));
  });
});
