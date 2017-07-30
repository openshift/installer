const _ = require('lodash');

const logger = logs => {
  console.log('==== BEGIN BROWSER LOGS ====');
  _.each(logs, log => {
    const { level, message } = log;
    const messageStr = _.isArray(message) ? message.join(' ') : message;

    switch (level) {
    case 'DEBUG':
      console.log(level, messageStr);
      break;
    case 'SEVERE':
      console.warn(level, messageStr);
      break;
    case 'INFO':
    default:
      console.info(level, messageStr);
    }
  });
  console.log('==== END BROWSER LOGS ====');
};

module.exports = {
  logger,
};
