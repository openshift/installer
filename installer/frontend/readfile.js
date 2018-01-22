export const readFile = (blob) => {
  const reader = new FileReader();
  const ret = new Promise((resolve, reject) => {
    reader.onerror = function () {
      reject(this.error);
    };
    reader.onabort = function () {
      reject('ABORTED');
    };
    reader.onload = function () {
      resolve(this.result);
    };
  });
  reader.readAsText(blob);
  return ret;
};
