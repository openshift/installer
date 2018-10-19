# OpenShift Installer Tests

## Running smoke tests locally

1. Expose environment variables

To run the smoke tests locally you need to set the following
environment variables:
``` bash
AWS_ACCESS_KEY_ID
AWS_SECRET_ACCESS_KEY
PULL_SECRET_PATH
```

Optionally you can also set:
```bash
DOMAIN
AWS_REGION
```

2. Launch the tests
Once the environment variables are set, run `./tests/run.sh aws`.

If you already have a cluster running, follow [smoke/README.md](./smoke/README.md).
