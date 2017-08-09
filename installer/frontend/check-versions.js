#!/usr/bin/env node

/* eslint-env node */
const semver = require('semver');
const engines = require('./package').engines;

if (!semver.satisfies(process.version, engines.node)) {
  console.error(`Need node ${engines.node} but you have ${process.version}`);
  process.exit(1);
}

// Yarn puts its version info in npm_config_user_agent env var O_o
const userAgent = process.env.npm_config_user_agent;
const result = /^yarn\/([\w.]+)/.exec(userAgent);
if (!result) {
  console.error(`Unknown yarn user agent: ${userAgent}`);
}
if (!semver.satisfies(result[1], engines.yarn)) {
  console.error(`Need yarn ${engines.yarn} but you have ${result[1]}`);
  process.exit(1);
}
